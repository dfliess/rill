package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/act/trigger"
	"go.uber.org/zap/exp/zapslog"
)

// Act plane defaults. The two schemas live in the same Act Postgres database but are kept apart: DBOS owns its
// workflow-engine tables, and the product schema holds Act's own run/approval/action/trigger tables.
const (
	defaultActDBOSSchema    = "act_dbos"
	defaultActProductSchema = "act_product"
	defaultActDispatch      = 10 * time.Second
	defaultActCloseTimeout  = 10 * time.Second
)

// ActConfig configures the in-process Act plane. Act is enabled only when PostgresDSN is set; with it empty the
// runtime serves exactly as before — no Act run/approval handlers, no DBOS worker, no trigger observer — so a
// deployment that does not provision an Act Postgres is unaffected (Kairos ADR-0012: minimal, gated wiring).
type ActConfig struct {
	// PostgresDSN is the connection string for Act's Postgres. It is the single gate: empty => Act disabled. The DBOS
	// engine opens its own connection to it; the run store and trigger store share one pool over it.
	PostgresDSN string
	// DBOSSchema holds the DBOS engine's workflow tables. Defaults to act_dbos.
	DBOSSchema string
	// ProductSchema holds Act's run/approval/action/trigger tables. Defaults to act_product.
	ProductSchema string
	// DispatchInterval is how often the trigger dispatcher polls each instance's outbox. Defaults to 10s.
	DispatchInterval time.Duration
	// ApplicationVersion pins the DBOS worker version. DBOS only recovers workflows that match it, so pinning it
	// across a rolling deploy keeps a new worker from recovering an incompatible in-flight run. Empty lets the
	// executor default; tests set a unique value to isolate their runs on the shared DBOS schema.
	ApplicationVersion string
}

// BootstrapAct wires and starts the in-process Act plane on this server, then returns an io.Closer that drains it.
// It is intended to be called once during startup, after NewServer and before serving, only when Act is configured
// (a non-empty PostgresDSN); the caller gates on that and defers the returned Close.
//
// It builds, in order: a shared Postgres pool for the product plane; the run store (product schema, migrated); the
// in-process DBOS executor over a production SessionRunner (catalog provider + claims-bound session factory that
// fails closed on nil claims + a fail-closed action sink); the trigger store (migrated); the execution observer,
// registered on the runtime so alert/report executions feed the trigger outbox; and a polling loop that drains that
// outbox per instance. Finally it calls ConfigureAct so the run/approval handlers go live. On any error it tears
// down whatever it already built and returns.
func (s *Server) BootstrapAct(ctx context.Context, cfg ActConfig) (io.Closer, error) {
	if cfg.PostgresDSN == "" {
		// Defensive: the caller gates on this, but never build a half-wired plane from an empty DSN.
		return nil, errors.New("act: BootstrapAct requires a Postgres DSN")
	}
	dbosSchema := cfg.DBOSSchema
	if dbosSchema == "" {
		dbosSchema = defaultActDBOSSchema
	}
	productSchema := cfg.ProductSchema
	if productSchema == "" {
		productSchema = defaultActProductSchema
	}
	dispatch := cfg.DispatchInterval
	if dispatch <= 0 {
		dispatch = defaultActDispatch
	}

	// The act components take a *slog.Logger; the runtime speaks zap. Bridge the server's logger so Act logs land in
	// the same stream, at the same levels, as the rest of the runtime.
	logger := slog.New(zapslog.NewHandler(s.logger.Core(), nil))

	// One pool backs both product stores (run store + trigger store); the DBOS engine opens its own connection. The
	// pool is created against the startup ctx but outlives it — its lifetime is owned by the returned closer.
	pool, err := pgxpool.New(context.WithoutCancel(ctx), cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("act: create postgres pool: %w", err)
	}
	// fail closes the pool (and any already-built teardown) before returning a construction error.
	fail := func(err error) (io.Closer, error) {
		pool.Close()
		return nil, err
	}

	runStore, err := act.NewPostgresRunStoreWithPool(pool, productSchema)
	if err != nil {
		return fail(fmt.Errorf("act: build run store: %w", err))
	}
	if err := runStore.Migrate(ctx); err != nil {
		return fail(fmt.Errorf("act: migrate run store: %w", err))
	}

	// The production runner: the catalog resolves agent snapshots, and the session factory binds each run to its
	// initiator's claims (failing closed on nil). Its ActionSink (the lean notify write) is retained only as the
	// fallback for a deployment without a gateway; with the gateway wired below it is never reached (the run's action
	// phase goes through the gateway).
	runner := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(s.runtime),
		Sessions: act.NewSessionFactory(s.runtime, s.activity),
		Actions:  act.NewNotifyActionSink(s.runtime.Email),
	}

	// The action gateway routes a run's proposed write through policy, the action ledger, server-side secret
	// resolution and an idempotent MCP write. The loop captures the structured tool call and the captured proposer
	// surfaces it; the gateway's ledger is the run store itself (its agent_actions table lives in the product schema),
	// so a run's lifecycle and its actions share one Postgres. This is the governed action path: a run investigates,
	// proposes a structured tool call, a human approves the exact arguments, and the gateway executes it exactly once.
	executor, err := act.NewDBOSExecutor(ctx, act.Config{
		DatabaseURL:        cfg.PostgresDSN,
		DatabaseSchema:     dbosSchema,
		ApplicationVersion: cfg.ApplicationVersion,
		Runner:             runner,
		Store:              runStore,
		Gateway:            act.NewMCPGateway(s.runtime, runStore, logger),
		Proposer: act.NewCapturedProposer(),
		Logger:   logger,
	})
	if err != nil {
		return fail(fmt.Errorf("act: build executor: %w", err))
	}

	// Trigger plane: the store's outbox tables live alongside the run tables in the product schema.
	triggerStore := trigger.NewStore(pool, productSchema)
	if err := triggerStore.Migrate(ctx); err != nil {
		executor.Close(defaultActCloseTimeout)
		return fail(fmt.Errorf("act: migrate trigger store: %w", err))
	}

	// The observer turns persisted alert/report executions into trigger-outbox events; the dispatcher drains the
	// outbox and starts runs on the executor. The observer must be registered before reconciliation produces
	// executions, which BootstrapAct's call site (startup, before serving) guarantees.
	observer := trigger.NewObserver(triggerStore, logger)
	s.runtime.SetExecutionObserver(observer)
	dispatcher := trigger.NewDispatcher(triggerStore, trigger.NewRuntimeTriggerCatalog(s.runtime), executor, trigger.DispatcherOptions{Logger: logger})

	// The polling loop runs under its own context, cancelled by Close, so shutdown stops dispatching regardless of
	// the startup ctx's fate.
	loopCtx, loopCancel := context.WithCancel(context.Background())
	loopDone := make(chan struct{})
	go func() {
		defer close(loopDone)
		runDispatchLoop(loopCtx, s, dispatcher, dispatch, logger)
	}()

	// Handlers go live only now that the store and executor both exist and share the run store, so a run started
	// through the executor is immediately readable through the run store.
	s.ConfigureAct(runStore, executor)

	logger.Info("act plane started",
		slog.String("dbos_schema", dbosSchema),
		slog.String("product_schema", productSchema),
		slog.Duration("dispatch_interval", dispatch),
	)

	return &actPlane{
		logger:       logger,
		server:       s,
		pool:         pool,
		executor:     executor,
		observer:     observer,
		loopCancel:   loopCancel,
		loopDone:     loopDone,
		closeTimeout: defaultActCloseTimeout,
	}, nil
}

// runDispatchLoop polls every instance's trigger outbox on an interval until ctx is cancelled. Dispatch is idempotent
// per (trigger, event), so a failed tick is retried on the next one; an error is logged and never stops the loop.
func runDispatchLoop(ctx context.Context, s *Server, dispatcher *trigger.Dispatcher, interval time.Duration, logger *slog.Logger) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			instances, err := s.runtime.Instances(ctx)
			if err != nil {
				logger.Error("act dispatch: list instances", slog.Any("error", err))
				continue
			}
			for _, inst := range instances {
				// Enqueue any due schedule ticks first, so this same iteration's Dispatch drains them: a schedule
				// trigger has no upstream execution feeding the outbox, so the ticker is what makes it fire.
				if _, err := dispatcher.EnqueueDueSchedules(ctx, inst.ID); err != nil {
					logger.Error("act schedule tick", slog.String("instance_id", inst.ID), slog.Any("error", err))
				}
				if _, err := dispatcher.Dispatch(ctx, inst.ID); err != nil {
					logger.Error("act dispatch", slog.String("instance_id", inst.ID), slog.Any("error", err))
				}
			}
		}
	}
}

// actPlane holds the running Act components so BootstrapAct's caller can tear them all down with one Close.
type actPlane struct {
	logger       *slog.Logger
	server       *Server
	pool         *pgxpool.Pool
	executor     *act.DBOSExecutor
	observer     *trigger.Observer
	loopCancel   context.CancelFunc
	loopDone     chan struct{}
	closeTimeout time.Duration
}

var _ io.Closer = (*actPlane)(nil)

// Close drains the Act plane in dependency order: stop feeding new work (clear the observer, stop the dispatch loop),
// then drain the executor's worker, then close the observer's publisher and the shared pool. It is safe to call once.
func (p *actPlane) Close() error {
	// Clear the observer first so no further execution is enqueued while we tear down, then stop the dispatcher loop.
	p.server.runtime.SetExecutionObserver(nil)
	p.loopCancel()
	<-p.loopDone

	// Drain the DBOS worker (bounded), then close the observer's publisher and release the shared pool.
	p.executor.Close(p.closeTimeout)
	p.observer.Close()
	p.pool.Close()
	p.logger.Info("act plane stopped")
	return nil
}

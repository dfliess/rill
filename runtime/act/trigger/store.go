package trigger

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Dispatch outcomes recorded in the ledger. A slot is first "reserved" (freezing the agent to run before the run is
// started), then "dispatched" once the run exists; "deduplicated" is a window-collapsed event that started no run.
const (
	outcomeReserved     = "reserved"
	outcomeDispatched   = "dispatched"
	outcomeDeduplicated = "deduplicated"
)

// Store is the durable backing for the trigger path: the outbox the observer writes to and the dispatch ledger the
// dispatcher writes to. It lives in the Platform Postgres, in its own schema, alongside (but separate from) the DBOS
// system tables (§15.2).
type Store struct {
	pool   *pgxpool.Pool
	schema string // already-sanitized, quoted identifier
}

// NewStore returns a Store over pool, keeping all tables in the given schema. The schema name is quoted as an
// identifier, so it is safe to interpolate into the DDL/DML below.
func NewStore(pool *pgxpool.Pool, schema string) *Store {
	return &Store{pool: pool, schema: pgx.Identifier{schema}.Sanitize()}
}

// OutboxEvent is a pending event claimed from the outbox for dispatch.
type OutboxEvent struct {
	ID            int64
	InstanceID    string
	EventID       string
	EventType     string
	ResourceKind  string
	ResourceName  string
	ExecutionTime time.Time
	OccurredOn    time.Time
	Status        string
}

// Migrate creates the schema and tables if they do not exist. It is idempotent and safe to call on every startup.
func (s *Store) Migrate(ctx context.Context) error {
	stmts := []string{
		fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, s.schema),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.agent_trigger_outbox (
			id             bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			instance_id    text NOT NULL,
			event_id       text NOT NULL,
			event_type     text NOT NULL,
			resource_kind  text NOT NULL,
			resource_name  text NOT NULL,
			execution_time timestamptz NOT NULL,
			occurred_on    timestamptz NOT NULL,
			status         text NOT NULL,
			payload        jsonb NOT NULL,
			payload_hash   text NOT NULL,
			schema_version integer NOT NULL,
			state          text NOT NULL DEFAULT 'pending',
			available_on   timestamptz NOT NULL DEFAULT now(),
			leased_until   timestamptz,
			attempts       integer NOT NULL DEFAULT 0,
			last_error     text,
			created_on     timestamptz NOT NULL DEFAULT now(),
			updated_on     timestamptz NOT NULL DEFAULT now(),
			UNIQUE (instance_id, event_id)
		)`, s.schema),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS agent_trigger_outbox_pending_idx
			ON %s.agent_trigger_outbox (instance_id, available_on) WHERE state = 'pending'`, s.schema),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.agent_trigger_events (
			id             bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			instance_id    text NOT NULL,
			trigger_id     text NOT NULL,
			event_id       text NOT NULL,
			event_type     text NOT NULL,
			agent_name     text NOT NULL,
			agent_run_id   text,
			outcome        text NOT NULL,
			dedup_window_seconds integer NOT NULL DEFAULT 0,
			created_on     timestamptz NOT NULL DEFAULT now(),
			UNIQUE (instance_id, trigger_id, event_id)
		)`, s.schema),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS agent_trigger_events_dispatched_idx
			ON %s.agent_trigger_events (instance_id, trigger_id, created_on) WHERE outcome = 'dispatched'`, s.schema),
	}
	for _, stmt := range stmts {
		if _, err := s.pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("act/trigger: migrate: %w", err)
		}
	}
	return nil
}

// EnqueueEvent publishes a derived event to the outbox. It is idempotent: a re-published event (same instance and
// event_id) is a no-op, which is what makes at-least-once delivery from the source safe (§9.6). It reports whether a
// new row was inserted.
func (s *Store) EnqueueEvent(ctx context.Context, e Event) (bool, error) {
	payload, err := json.Marshal(e.Envelope())
	if err != nil {
		return false, fmt.Errorf("act/trigger: marshal envelope: %w", err)
	}
	sum := sha256.Sum256(payload)
	q := fmt.Sprintf(`INSERT INTO %s.agent_trigger_outbox
		(instance_id, event_id, event_type, resource_kind, resource_name, execution_time, occurred_on, status, payload, payload_hash, schema_version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb, $10, $11)
		ON CONFLICT (instance_id, event_id) DO NOTHING`, s.schema)
	tag, err := s.pool.Exec(ctx, q,
		e.InstanceID, e.EventID, e.EventType, e.ResourceKind, e.ResourceName,
		e.ExecutionTime, e.OccurredOn, e.Status, string(payload), hex.EncodeToString(sum[:]), schemaVersion)
	if err != nil {
		return false, fmt.Errorf("act/trigger: enqueue event: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

// ClaimPending leases up to limit pending events for an instance for leaseSeconds, so a crashed dispatcher's rows
// become available again once the lease expires (at-least-once). FOR UPDATE SKIP LOCKED lets concurrent pollers not
// step on each other, though shadow mode runs a single worker (§29).
func (s *Store) ClaimPending(ctx context.Context, instanceID string, limit, leaseSeconds int) ([]OutboxEvent, error) {
	q := fmt.Sprintf(`UPDATE %s.agent_trigger_outbox
		SET leased_until = now() + make_interval(secs => $2), attempts = attempts + 1, updated_on = now()
		WHERE id IN (
			SELECT id FROM %s.agent_trigger_outbox
			WHERE instance_id = $1 AND state = 'pending' AND available_on <= now()
			  AND (leased_until IS NULL OR leased_until < now())
			ORDER BY id
			LIMIT $3
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, instance_id, event_id, event_type, resource_kind, resource_name, execution_time, occurred_on, status`,
		s.schema, s.schema)
	rows, err := s.pool.Query(ctx, q, instanceID, leaseSeconds, limit)
	if err != nil {
		return nil, fmt.Errorf("act/trigger: claim pending: %w", err)
	}
	defer rows.Close()

	var out []OutboxEvent
	for rows.Next() {
		var e OutboxEvent
		if err := rows.Scan(&e.ID, &e.InstanceID, &e.EventID, &e.EventType, &e.ResourceKind, &e.ResourceName,
			&e.ExecutionTime, &e.OccurredOn, &e.Status); err != nil {
			return nil, fmt.Errorf("act/trigger: scan pending: %w", err)
		}
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("act/trigger: iterate pending: %w", err)
	}
	return out, nil
}

// MarkDone marks an outbox event as fully dispatched. It will not be claimed again.
func (s *Store) MarkDone(ctx context.Context, id int64) error {
	q := fmt.Sprintf(`UPDATE %s.agent_trigger_outbox SET state = 'done', leased_until = NULL, updated_on = now() WHERE id = $1`, s.schema)
	if _, err := s.pool.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("act/trigger: mark done: %w", err)
	}
	return nil
}

// MarkFailed releases the lease and schedules a retry after backoffSeconds, recording the error. The event stays
// pending so a later poll re-attempts it (at-least-once); nothing is dropped silently.
func (s *Store) MarkFailed(ctx context.Context, id int64, cause string, backoffSeconds int) error {
	q := fmt.Sprintf(`UPDATE %s.agent_trigger_outbox
		SET leased_until = NULL, available_on = now() + make_interval(secs => $2), last_error = $3, updated_on = now()
		WHERE id = $1`, s.schema)
	if _, err := s.pool.Exec(ctx, q, id, backoffSeconds, cause); err != nil {
		return fmt.Errorf("act/trigger: mark failed: %w", err)
	}
	return nil
}

// HasRecentDispatch reports whether a run was already dispatched for (instance, trigger) at or after since. It is
// the window-deduplication query: within the window, a further event collapses onto the existing run instead of
// starting a second one (§9.2).
func (s *Store) HasRecentDispatch(ctx context.Context, instanceID, triggerID string, since time.Time) (bool, error) {
	q := fmt.Sprintf(`SELECT EXISTS (
		SELECT 1 FROM %s.agent_trigger_events
		WHERE instance_id = $1 AND trigger_id = $2 AND outcome = 'dispatched' AND created_on >= $3
	)`, s.schema)
	var exists bool
	if err := s.pool.QueryRow(ctx, q, instanceID, triggerID, since).Scan(&exists); err != nil {
		return false, fmt.Errorf("act/trigger: check recent dispatch: %w", err)
	}
	return exists, nil
}

// RecordDispatch writes the dispatch ledger row for (instance, trigger, event). The unique (instance_id,
// trigger_id, event_id) constraint makes it the point where a replayed event is deduplicated by identity: it reports
// whether this call inserted the row (true) or a prior one already claimed the slot (false). runID may be empty when
// the outcome is "deduplicated" (window dedup produced no run).
func (s *Store) RecordDispatch(ctx context.Context, instanceID, triggerID, eventID, eventType, agentName, runID, outcome string, dedupWindowSeconds int) (bool, error) {
	q := fmt.Sprintf(`INSERT INTO %s.agent_trigger_events
		(instance_id, trigger_id, event_id, event_type, agent_name, agent_run_id, outcome, dedup_window_seconds)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), $7, $8)
		ON CONFLICT (instance_id, trigger_id, event_id) DO NOTHING`, s.schema)
	tag, err := s.pool.Exec(ctx, q, instanceID, triggerID, eventID, eventType, agentName, runID, outcome, dedupWindowSeconds)
	if err != nil {
		return false, fmt.Errorf("act/trigger: record dispatch: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

// Reservation is the frozen dispatch identity for one (instance, trigger, event). AgentName is the agent captured at
// reservation time, which the dispatcher must use to start the run rather than the trigger's current agent: it is
// what keeps the composed run id stable across retries even if the trigger is re-pointed to a different agent.
type Reservation struct {
	AgentName string // the agent frozen at reservation; use this, not the trigger's live agent
	RunID     string // the started run, empty until MarkDispatched records it
	Outcome   string // reserved | dispatched
	New       bool   // true only on the call that first reserved the slot
}

// ReserveDispatch reserves the (instance, trigger, event) slot before a run is started, freezing agentName so the
// run's identity derives from the stable (instance, trigger, event) plus the agent captured here — never the
// trigger's live, mutable agent. On the first call it inserts a 'reserved' row (New=true); on any later call, a
// retry after a crash between reserve and start/mark, it returns the stored row so the caller reuses the frozen
// agent and the recorded run id. This is what stops recovery from starting a second run when the trigger's agent
// changed in between (§9.6, §12).
func (s *Store) ReserveDispatch(ctx context.Context, instanceID, triggerID, eventID, eventType, agentName string, dedupWindowSeconds int) (Reservation, error) {
	q := fmt.Sprintf(`WITH ins AS (
			INSERT INTO %s.agent_trigger_events
				(instance_id, trigger_id, event_id, event_type, agent_name, outcome, dedup_window_seconds)
			VALUES ($1, $2, $3, $4, $5, '%s', $6)
			ON CONFLICT (instance_id, trigger_id, event_id) DO NOTHING
			RETURNING agent_name, agent_run_id, outcome
		)
		SELECT agent_name, COALESCE(agent_run_id, ''), outcome, true FROM ins
		UNION ALL
		SELECT agent_name, COALESCE(agent_run_id, ''), outcome, false
			FROM %s.agent_trigger_events
			WHERE instance_id = $1 AND trigger_id = $2 AND event_id = $3 AND NOT EXISTS (SELECT 1 FROM ins)`,
		s.schema, outcomeReserved, s.schema)
	var r Reservation
	if err := s.pool.QueryRow(ctx, q, instanceID, triggerID, eventID, eventType, agentName, dedupWindowSeconds).
		Scan(&r.AgentName, &r.RunID, &r.Outcome, &r.New); err != nil {
		return Reservation{}, fmt.Errorf("act/trigger: reserve dispatch: %w", err)
	}
	return r, nil
}

// MarkDispatched records the started run id on a reserved slot, moving it from 'reserved' to 'dispatched'. It only
// advances a reserved row, so it is idempotent: a repeat (the slot already dispatched) updates nothing.
func (s *Store) MarkDispatched(ctx context.Context, instanceID, triggerID, eventID, runID string) error {
	q := fmt.Sprintf(`UPDATE %s.agent_trigger_events
		SET outcome = '%s', agent_run_id = $4
		WHERE instance_id = $1 AND trigger_id = $2 AND event_id = $3 AND outcome = '%s'`,
		s.schema, outcomeDispatched, outcomeReserved)
	if _, err := s.pool.Exec(ctx, q, instanceID, triggerID, eventID, runID); err != nil {
		return fmt.Errorf("act/trigger: mark dispatched: %w", err)
	}
	return nil
}

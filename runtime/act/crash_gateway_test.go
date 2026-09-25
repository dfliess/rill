package act_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/stretchr/testify/require"
)

// The crash tests below kill a real worker process on the path production runs: the segmented loop through the Tool
// Gateway, with the production proposer and a Postgres ledger. The crash tests used to exercise a simulated path that
// production never took, and that is how kairos-cloud#170 hid behind a green suite.

const gatewayCrashInstance, gatewayCrashAgent = "inst", "crash-agent"

// crashSegmentRunner proposes ACT_ACTIONS governed writes on the first segment and closes once their results come
// back, marking every segment in the step log. ACT_CRASH_AT=segment kills the process inside the first segment, after
// the snapshot was checkpointed and before the segment was.
type crashSegmentRunner struct {
	log     func(step, detail string)
	actions int
	crashAt string
}

func (r *crashSegmentRunner) LoadSnapshot(_ context.Context, _, agentName string) (*ai.AgentSnapshot, error) {
	r.log("load_snapshot", agentName)
	return &ai.AgentSnapshot{Name: agentName, Instructions: "crash on the gateway path"}, nil
}

func (r *crashSegmentRunner) RunSegment(_ context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	if len(in.Resume) > 0 {
		r.log("segment", "resumed")
		return act.RunSegmentResult{Response: "Hecho.", SessionID: in.SessionID}, nil
	}
	r.log("segment", "first")
	if r.crashAt == "segment" {
		killSelf()
	}
	proposed := make([]*ai.ProposedAction, r.actions)
	for i := range proposed {
		proposed[i] = &ai.ProposedAction{
			ToolCallID: fmt.Sprintf("tc-%d", i+1),
			Tool:       "jira.create_issue",
			Connector:  "jira_ops",
			Args:       map[string]any{"summary": fmt.Sprintf("ticket %d", i+1)},
			Summary:    fmt.Sprintf("crear ticket %d", i+1),
		}
	}
	return act.RunSegmentResult{Response: "Propongo tickets.", SessionID: "sess-crash", Proposed: proposed}, nil
}

// crashActionExecutor stands in for the MCP server: it marks every write in the step log, and ACT_CRASH_AT=execute
// kills the process inside the write, after the gateway claimed it and before it could record the outcome.
type crashActionExecutor struct {
	log     func(step, detail string)
	crashAt string
}

func (x *crashActionExecutor) Execute(_ context.Context, req act.ExecuteRequest) (act.ExecuteResult, error) {
	x.log("execute", req.ToolCallID)
	if x.crashAt == "execute" {
		killSelf()
	}
	return act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-" + req.ToolCallID}, nil
}

func (x *crashActionExecutor) Verify(context.Context, act.VerifyRequest) (act.VerifyResult, error) {
	return act.VerifyResult{Confirmed: true}, nil
}

// killSelf hard-kills the process and blocks until the signal lands: SIGKILL is asynchronous, and returning would let
// the caller run code the crash was meant to prevent.
func killSelf() {
	_ = syscall.Kill(os.Getpid(), syscall.SIGKILL)
	time.Sleep(time.Hour)
}

// runGatewayCrashWorker is the worker subprocess for the gateway-path crash tests. It always exits the process.
//
//	start          start the run, then do what ACT_THEN says and wait to be killed by a crash hook.
//	recover        let DBOS recover the run, approve every tool call in ACT_APPROVE as it comes up, print the result.
func runGatewayCrashWorker() {
	key := os.Getenv("ACT_RUN_ID")
	var runID string // set by Start, or found by its key when recovering
	logPath := os.Getenv("ACT_ACTION_LOG")
	mark := func(step, detail string) { markStep(logPath, step, detail) }
	fail := func(what string, err error) {
		fmt.Fprintln(os.Stderr, what+":", err)
		os.Exit(3)
	}

	ctx := context.Background()
	store, err := act.NewPostgresRunStore(ctx, act.StoreConfig{DatabaseURL: os.Getenv("ACT_DSN"), Schema: os.Getenv("ACT_STORE_SCHEMA")})
	if err != nil {
		fail("store", err)
	}
	if err := store.Migrate(ctx); err != nil {
		fail("migrate", err)
	}
	actions, _ := strconv.Atoi(os.Getenv("ACT_ACTIONS"))
	if actions == 0 {
		actions = 1
	}
	timeout, _ := time.ParseDuration(os.Getenv("ACT_GATEWAY_TIMEOUT"))
	crashAt := os.Getenv("ACT_CRASH_AT")

	// Build the executor on a CANCELLABLE context, as the runtime does (graceful.WithCancelOnTerminate), so the
	// park-sigterm mode can reproduce production's shutdown order.
	hostCtx, stopHost := context.WithCancel(context.Background())
	defer stopHost()
	e, err := act.NewDBOSExecutor(hostCtx, act.Config{
		DatabaseURL:        os.Getenv("ACT_DSN"),
		DatabaseSchema:     os.Getenv("ACT_SCHEMA"),
		ApplicationVersion: os.Getenv("ACT_APP_VERSION"),
		Runner:             &crashSegmentRunner{log: mark, actions: actions, crashAt: crashAt},
		Store:              store,
		Gateway: &act.Gateway{
			Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative)),
			Ledger:   store,
			Executor: &crashActionExecutor{log: mark, crashAt: crashAt},
			Timeout:  timeout,
		},
		Proposer: act.NewCapturedProposer(),
		Logger:   slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
	})
	if err != nil {
		fail("executor", err)
	}

	approve := func(toolCallID string) {
		approvalID := act.ApprovalIDForToolCall(runID, toolCallID)
		deadline := time.Now().Add(20 * time.Second)
		for {
			a, err := store.GetApproval(ctx, gatewayCrashInstance, approvalID)
			if err == nil && a.Status == act.ApprovalStatusPending {
				break
			}
			if time.Now().After(deadline) {
				fail("await approval "+toolCallID, fmt.Errorf("not pending in time (last error: %w)", err))
			}
			time.Sleep(50 * time.Millisecond)
		}
		if err := e.Resume(ctx, runID, act.ApprovalApproved, toolCallID); err != nil {
			fail("resume "+toolCallID, err)
		}
	}

	switch os.Getenv("ACT_MODE") {
	case "start":
		runID, err = e.Start(ctx, act.AgentRunInput{
			InstanceID: gatewayCrashInstance, AgentName: gatewayCrashAgent, Prompt: "alerta", IdempotencyKey: key,
			Actor: act.Actor{Subject: "user:alice"},
		})
		if err != nil {
			fail("start", err)
		}
		switch os.Getenv("ACT_THEN") {
		case "approve-first":
			approve("tc-1")
		case "park-kill", "park-graceful", "park-sigterm":
			// Wait until the run is genuinely parked on its approval, then go down the way the mode says.
			deadline := time.Now().Add(20 * time.Second)
			for {
				run, err := store.GetRun(ctx, gatewayCrashInstance, runID)
				if err == nil && run.Status == act.RunStatusWaitingApproval {
					break
				}
				if time.Now().After(deadline) {
					fail("await waiting_approval", fmt.Errorf("not parked in time (last error: %w)", err))
				}
				time.Sleep(50 * time.Millisecond)
			}
			switch os.Getenv("ACT_THEN") {
			case "park-kill":
				killSelf()
			case "park-sigterm":
				// The signal cancels the host first; the runtime only closes Act after it stops serving.
				stopHost()
				time.Sleep(2 * time.Second)
			}
			e.Close(5 * time.Second)
			os.Exit(0)
		case "approve-first-then-crash":
			// Decide one action of the batch, wait until its write is recorded, then die with the rest undecided.
			approve("tc-1")
			deadline := time.Now().Add(20 * time.Second)
			for {
				a, err := store.GetAction(ctx, gatewayCrashInstance, runID, "tc-1")
				if err == nil && (a.Status == act.ActionSucceeded || a.Status == act.ActionVerified) {
					break
				}
				if time.Now().After(deadline) {
					fail("await tc-1 executed", fmt.Errorf("not executed in time (last error: %w)", err))
				}
				time.Sleep(50 * time.Millisecond)
			}
			killSelf()
		}
		time.Sleep(time.Minute)
		fail("start", fmt.Errorf("the crash hook never fired"))
	case "recover":
		if runID, err = gatewayCrashRunID(ctx, store, key); err != nil {
			fail("find run", err)
		}
		for _, tc := range strings.Fields(os.Getenv("ACT_APPROVE")) {
			approve(tc)
		}
		res, err := e.Result(runID)
		if err != nil {
			fail("recover result", err)
		}
		fmt.Printf("RECOVERED status=%s ref=%s\n", res.Status, res.ActionRef)
		os.Exit(0)
	}
	fail("mode", fmt.Errorf("unknown ACT_MODE %q", os.Getenv("ACT_MODE")))
}

// gatewayCrashRunID finds the run the start worker created under key: the run ID is an opaque UUID, so a process that
// did not start the run finds it by its key (ADR-0034).
func gatewayCrashRunID(ctx context.Context, store *act.PostgresRunStore, key string) (string, error) {
	run, err := store.LatestRunInKeyFamily(ctx, gatewayCrashInstance, gatewayCrashAgent, key)
	if err != nil {
		return "", err
	}
	return run.RunID, nil
}

// gatewayCrash runs one start worker, then one recover worker, and returns the recover worker's output and the step
// log. killed says whether the start worker is meant to die (a crash) or exit cleanly (a shutdown); between, when set,
// runs after the start worker ends and before recovery, with the run's ID.
func gatewayCrash(t *testing.T, killed bool, startEnv, recoverEnv map[string]string, between func(runID string)) (string, stepLog) {
	t.Helper()
	dsn, schema := requirePostgres(t)
	storeSchema := "act_crash_gw_" + strings.ReplaceAll(uuid.NewString()[:8], "-", "")
	key := "gw-crash-" + uuid.NewString()
	logPath := filepath.Join(t.TempDir(), "steps.log")
	base := map[string]string{
		"ACT_WORKER_PATH": "gateway", "ACT_RUN_ID": key, "ACT_DSN": dsn, "ACT_SCHEMA": schema,
		"ACT_APP_VERSION": "act-crash-gw-" + uuid.NewString(), "ACT_STORE_SCHEMA": storeSchema, "ACT_ACTION_LOG": logPath,
	}
	with := func(extra map[string]string) map[string]string {
		env := map[string]string{}
		for k, v := range base {
			env[k] = v
		}
		for k, v := range extra {
			env[k] = v
		}
		return env
	}

	out, err := workerCmd(with(startEnv)).CombinedOutput()
	if killed {
		require.Error(t, err, "the start worker should have been killed; output:\n%s", out)
	} else {
		require.NoError(t, err, "the start worker should have shut down cleanly; output:\n%s", out)
	}
	if between != nil {
		store, err := act.NewPostgresRunStore(t.Context(), act.StoreConfig{DatabaseURL: dsn, Schema: storeSchema})
		require.NoError(t, err)
		runID, err := gatewayCrashRunID(t.Context(), store, key)
		store.Close()
		require.NoError(t, err)
		between(runID)
	}
	out2, err := workerCmd(with(recoverEnv)).CombinedOutput()
	require.NoError(t, err, "the recover worker should finish cleanly; output:\n%s", out2)
	return string(out2), readSteps(t, logPath)
}

// TestGatewayRunRecoversFromCrashMidSegment: a worker killed inside the model loop, before the segment checkpointed,
// is recovered by a fresh process that re-runs only that segment and carries the run to its end.
func TestGatewayRunRecoversFromCrashMidSegment(t *testing.T) {
	out, steps := gatewayCrash(t, true,
		map[string]string{"ACT_MODE": "start", "ACT_CRASH_AT": "segment"},
		map[string]string{"ACT_MODE": "recover", "ACT_APPROVE": "tc-1"}, nil)

	require.Contains(t, out, "RECOVERED status=succeeded")
	require.Equal(t, 1, steps.count("load_snapshot"), "the checkpointed snapshot step must replay, not re-run")
	require.Equal(t, 1, steps.count("execute"), "the write runs exactly once")
	require.NotEqual(t, steps.pid("load_snapshot"), steps.pid("execute"), "recovery crossed a process boundary")
}

// TestGatewayRunCrashMidWriteNeverReissuesIt: a worker killed inside the external write leaves an action whose
// outcome nobody knows. The recovered run must not issue the write again; it waits out the dead attempt's lease and
// hands the action to a human as indeterminate (kairos-cloud#170). The lease has a fixed 60s margin, so this test
// takes about a minute and runs in parallel with the rest.
func TestGatewayRunCrashMidWriteNeverReissuesIt(t *testing.T) {
	t.Parallel()
	out, steps := gatewayCrash(t, true,
		map[string]string{"ACT_MODE": "start", "ACT_CRASH_AT": "execute", "ACT_THEN": "approve-first", "ACT_GATEWAY_TIMEOUT": "1s"},
		map[string]string{"ACT_MODE": "recover", "ACT_GATEWAY_TIMEOUT": "1s"}, nil)

	require.Contains(t, out, "RECOVERED status=failed", "an unconfirmed write ends the run for a human to resolve")
	require.Equal(t, 1, steps.count("execute"), "the interrupted write must never be issued again")
}

// TestGatewayBatchRecoversWithOneActionDecided: a batch of two actions where the first was approved and written before
// the crash. The recovered run must not write the first again, and must still ask for and perform the second.
func TestGatewayBatchRecoversWithOneActionDecided(t *testing.T) {
	out, steps := gatewayCrash(t, true,
		map[string]string{"ACT_MODE": "start", "ACT_ACTIONS": "2", "ACT_THEN": "approve-first-then-crash"},
		map[string]string{"ACT_MODE": "recover", "ACT_ACTIONS": "2", "ACT_APPROVE": "tc-2"}, nil)

	require.Contains(t, out, "RECOVERED status=succeeded")
	var tc1, tc2 int
	for _, s := range steps.steps {
		if s.step == "execute" {
			switch {
			case s.pid == steps.pid("load_snapshot"):
				tc1++
			default:
				tc2++
			}
		}
	}
	require.Equal(t, 1, tc1, "tc-1 was written once, before the crash, and never again")
	require.Equal(t, 1, tc2, "tc-2 was written once, after recovery")
}

package act_test

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/stretchr/testify/require"
)

// TestMain lets this test binary re-exec itself as a DBOS worker subprocess (ACT_TEST_WORKER=1), which the crash
// test uses to kill a worker mid-run and recover it in a fresh process. The env-var check runs before any test, so
// a worker subprocess never runs the test suite (and never forks further workers).
func TestMain(m *testing.M) {
	if os.Getenv("ACT_TEST_WORKER") == "1" {
		runCrashWorker() // never returns
	}
	os.Exit(m.Run())
}

// fileRunner is a durable-step-instrumented act.Runner for the crash test. It needs no runtime: it returns a fixed
// snapshot and appends one observable line to a log file per step, so the parent can see which steps ran, in which
// process. It carries the crash and hand-off hooks the subprocess pattern needs.
type fileRunner struct {
	logPath  string
	snapshot *ai.AgentSnapshot

	// runAgentSentinel, if set, makes RunSegment block until the file exists. The start worker writes it only after
	// it has sent the approval, so run_agent completes (and checkpoints) with the approval already buffered.
	runAgentSentinel string
	// crashAt, if "apply_action", SIGKILLs the process at the start of ApplyAction: after the earlier steps are
	// checkpointed but before this action's effect, which is exactly the window the recovery test targets.
	crashAt string
}

var _ act.Runner = (*fileRunner)(nil)

func (r *fileRunner) LoadSnapshot(_ context.Context, _, agentName string) (*ai.AgentSnapshot, error) {
	r.mark("load_snapshot", agentName)
	return r.snapshot, nil
}

func (r *fileRunner) RunSegment(_ context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	r.mark("run_agent", in.Snapshot.Name)
	if r.runAgentSentinel != "" {
		deadline := time.Now().Add(20 * time.Second)
		for time.Now().Before(deadline) {
			if _, err := os.Stat(r.runAgentSentinel); err == nil {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	return act.RunSegmentResult{Response: "respuesta simulada", SessionID: "sess-" + in.Snapshot.Name}, nil
}

func (r *fileRunner) ApplyAction(_ context.Context, req act.ActionRequest) (act.ActionResult, error) {
	if r.crashAt == "apply_action" {
		// Hard crash: unrecoverable in-process, so DBOS re-runs this uncheckpointed step on the next worker.
		_ = syscall.Kill(os.Getpid(), syscall.SIGKILL)
	}
	r.mark("apply_action", req.RunID)
	return act.ActionResult{Ref: "ticket-" + req.RunID}, nil
}

func (r *fileRunner) mark(step, detail string) {
	line := fmt.Sprintf("ts=%s pid=%d step=%s detail=%s\n", time.Now().Format(time.RFC3339Nano), os.Getpid(), step, detail)
	f, err := os.OpenFile(r.logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mark:", err)
		return
	}
	defer f.Close()
	_, _ = f.WriteString(line)
}

// runCrashWorker is the entry point of a worker subprocess. In "start" mode it launches a run, approves it, and is
// SIGKILLed inside apply_action. In "recover" mode it launches (recovering the pending run) and waits for it to
// finish. It always exits the process.
func runCrashWorker() {
	const instanceID, agentName = "inst", "crash-agent"
	key := os.Getenv("ACT_RUN_ID")
	// Start derives the durable workflow ID from (instance, agent, idempotency key); Resume/Result address the run
	// by that derived ID. Build it from the same parts here, since the recover worker never calls Start.
	workflowID := act.ComposeRunID(instanceID, agentName, key)
	runner := &fileRunner{
		logPath:          os.Getenv("ACT_ACTION_LOG"),
		runAgentSentinel: os.Getenv("ACT_RUN_AGENT_SENTINEL"),
		crashAt:          os.Getenv("ACT_CRASH_AT"),
		snapshot:         &ai.AgentSnapshot{Name: agentName, Instructions: "crash spike"},
	}

	// The park modes wire a run store so the worker can wait for the run to actually reach waiting_approval (a real
	// handshake) before killing it, rather than guessing with a sleep. Other modes leave it nil (the workflow guards
	// every store write on non-nil), so their behaviour is unchanged.
	var store act.RunStore
	if schema := os.Getenv("ACT_STORE_SCHEMA"); schema != "" {
		ps, err := act.NewPostgresRunStore(context.Background(), act.StoreConfig{DatabaseURL: os.Getenv("ACT_DSN"), Schema: schema})
		if err != nil {
			fmt.Fprintln(os.Stderr, "worker store init:", err)
			os.Exit(3)
		}
		if err := ps.Migrate(context.Background()); err != nil {
			fmt.Fprintln(os.Stderr, "worker store migrate:", err)
			os.Exit(3)
		}
		store = ps
	}

	e, err := act.NewDBOSExecutor(context.Background(), act.Config{
		DatabaseURL:        os.Getenv("ACT_DSN"),
		DatabaseSchema:     os.Getenv("ACT_SCHEMA"),
		ApplicationVersion: os.Getenv("ACT_APP_VERSION"),
		Runner:             runner,
		Store:              store,
		Logger:             slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "worker init:", err)
		os.Exit(3)
	}
	defer e.Close(5 * time.Second)

	switch os.Getenv("ACT_MODE") {
	case "start":
		if _, err := e.Start(context.Background(), act.AgentRunInput{
			InstanceID: instanceID, AgentName: agentName, Prompt: "alerta", IdempotencyKey: key,
		}); err != nil {
			fmt.Fprintln(os.Stderr, "start:", err)
			os.Exit(3)
		}
		// Approve now (durable); the run consumes it when it reaches the approval wait. Then release run_agent, so
		// it only completes once the approval is buffered.
		if err := e.Resume(context.Background(), workflowID, act.ApprovalApproved, ""); err != nil {
			fmt.Fprintln(os.Stderr, "resume:", err)
			os.Exit(3)
		}
		if err := os.WriteFile(runner.runAgentSentinel, []byte("go"), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "sentinel:", err)
			os.Exit(3)
		}
		// Expect to be SIGKILLed inside apply_action before this returns.
		res, err := e.Result(workflowID)
		fmt.Printf("UNEXPECTED completion status=%v err=%v\n", res.Status, err)
		os.Exit(0)
	case "start-park":
		// Start the run and wait until it has actually reached the approval wait, then die hard WITHOUT approving. This
		// leaves the workflow PENDING (parked mid-Recv), which is exactly the state a hard-killed deploy leaves an
		// approval-waiting run in. The handshake (not a sleep) is what makes the crash land while it is suspended.
		if _, err := e.Start(context.Background(), act.AgentRunInput{
			InstanceID: instanceID, AgentName: agentName, Prompt: "alerta", IdempotencyKey: key,
		}); err != nil {
			fmt.Fprintln(os.Stderr, "start:", err)
			os.Exit(3)
		}
		if err := awaitWaitingApproval(store, instanceID, workflowID); err != nil {
			fmt.Fprintln(os.Stderr, "await waiting_approval:", err)
			os.Exit(3)
		}
		// No graceful DBOS shutdown, so nothing terminalizes the parked workflow. SIGKILL is asynchronous, so block
		// instead of exiting: an os.Exit here would race the kernel and could exit 0 cleanly, defeating the crash.
		_ = syscall.Kill(os.Getpid(), syscall.SIGKILL)
		time.Sleep(time.Hour)
	case "start-park-graceful":
		// Like start-park, but instead of a hard kill it does a GRACEFUL DBOS shutdown while the run is parked on the
		// approval Recv, mirroring a deploy that drains the executor. This is the path that must NOT terminalize the
		// parked workflow if it is to recover on the next process.
		if _, err := e.Start(context.Background(), act.AgentRunInput{
			InstanceID: instanceID, AgentName: agentName, Prompt: "alerta", IdempotencyKey: key,
		}); err != nil {
			fmt.Fprintln(os.Stderr, "start:", err)
			os.Exit(3)
		}
		if err := awaitWaitingApproval(store, instanceID, workflowID); err != nil {
			fmt.Fprintln(os.Stderr, "await waiting_approval:", err)
			os.Exit(3)
		}
		e.Close(5 * time.Second)
		os.Exit(0)
	case "recover-approve":
		// Building the executor above already launched DBOS recovery, which re-runs the parked workflow and re-enters
		// its approval Recv. Deliver the (durable) approval and wait for the recovered run to finish.
		if err := e.Resume(context.Background(), workflowID, act.ApprovalApproved, ""); err != nil {
			fmt.Fprintln(os.Stderr, "resume:", err)
			os.Exit(3)
		}
		res, err := e.Result(workflowID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "recover result:", err)
			os.Exit(4)
		}
		fmt.Printf("RECOVERED status=%s ref=%s\n", res.Status, res.ActionRef)
		os.Exit(0)
	case "recover":
		res, err := e.Result(workflowID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "recover result:", err)
			os.Exit(4)
		}
		fmt.Printf("RECOVERED status=%s ref=%s\n", res.Status, res.ActionRef)
		os.Exit(0)
	default:
		fmt.Fprintln(os.Stderr, "unknown ACT_MODE")
		os.Exit(5)
	}
}

// TestExecutorRecoversAfterCrash kills a worker between the checkpointed steps and the external write, then
// recovers in a fresh process. It asserts the completed steps do NOT re-execute and the run finishes with the
// action performed exactly once.
func TestExecutorRecoversAfterCrash(t *testing.T) {
	dsn, schema := requirePostgres(t)

	version := "act-crash-" + uuid.NewString()
	runID := "crash-run-" + uuid.NewString()
	dir := t.TempDir()
	logPath := filepath.Join(dir, "actions.log")
	sentinel := filepath.Join(dir, "approval.sentinel")

	// Phase 1: worker A starts the run, approves it, and is SIGKILLed at the start of apply_action.
	a := workerCmd(map[string]string{
		"ACT_MODE": "start", "ACT_RUN_ID": runID,
		"ACT_DSN": dsn, "ACT_SCHEMA": schema, "ACT_APP_VERSION": version,
		"ACT_ACTION_LOG": logPath, "ACT_RUN_AGENT_SENTINEL": sentinel, "ACT_CRASH_AT": "apply_action",
	})
	out, err := a.CombinedOutput()
	require.Error(t, err, "worker A should have been killed, not exit cleanly; output:\n%s", out)

	steps1 := readSteps(t, logPath)
	require.Equal(t, 1, steps1.count("load_snapshot"), "load_snapshot should have run once in worker A")
	require.Equal(t, 1, steps1.count("run_agent"), "run_agent should have run once in worker A")
	require.Equal(t, 0, steps1.count("apply_action"), "apply_action must NOT have run before the crash")

	// Phase 2: worker B recovers. The completed steps must replay from their checkpoints (not re-execute), and the
	// run must finish with the action performed exactly once.
	b := workerCmd(map[string]string{
		"ACT_MODE": "recover", "ACT_RUN_ID": runID,
		"ACT_DSN": dsn, "ACT_SCHEMA": schema, "ACT_APP_VERSION": version,
		"ACT_ACTION_LOG": logPath,
	})
	out2, err2 := b.CombinedOutput()
	require.NoError(t, err2, "worker B should finish cleanly; output:\n%s", out2)
	require.Contains(t, string(out2), "RECOVERED status=succeeded", "recovered run should succeed; output:\n%s", out2)

	steps2 := readSteps(t, logPath)
	require.Equal(t, 1, steps2.count("load_snapshot"), "completed step load_snapshot was re-executed after recovery")
	require.Equal(t, 1, steps2.count("run_agent"), "completed step run_agent was re-executed after recovery")
	require.Equal(t, 1, steps2.count("apply_action"), "apply_action must run exactly once, in worker B")

	// The initial steps ran in worker A; the action ran in worker B: recovery crossed a process boundary.
	require.NotEqual(t, steps2.pid("run_agent"), steps2.pid("apply_action"),
		"run_agent and apply_action should have executed in different processes")
}

// TestExecutorRecoversApprovalWaitAfterCrash proves the durability property that matters for a run parked on a human
// approval: a worker hard-killed while the run waits on the approval Recv leaves the workflow PENDING, and a fresh
// worker recovers it, re-enters the wait, and finishes the run once the approval is delivered. This is the crash-kill
// baseline; the production bug is the graceful shutdown path, which instead terminalizes the parked workflow.
func TestExecutorRecoversApprovalWaitAfterCrash(t *testing.T) {
	dsn, schema := requirePostgres(t)

	version := "act-park-" + uuid.NewString()
	runID := "park-run-" + uuid.NewString()
	storeSchema := "act_park_store_" + uuid.NewString()[:8]

	// Phase 1: worker A starts the run and is SIGKILLed while it is parked on the approval Recv (never approved).
	a := workerCmd(map[string]string{
		"ACT_MODE": "start-park", "ACT_RUN_ID": runID,
		"ACT_DSN": dsn, "ACT_SCHEMA": schema, "ACT_APP_VERSION": version, "ACT_STORE_SCHEMA": storeSchema,
	})
	out, err := a.CombinedOutput()
	require.Error(t, err, "worker A should have been SIGKILLed while parked on the approval wait; output:\n%s", out)

	// Phase 2: worker B recovers the PENDING workflow, delivers the approval, and the run finishes succeeded. It wires
	// the SAME store as worker A so the workflow replays the identical step sequence (the store-emit steps included) —
	// a different store presence between the two would trip DBOS's determinism check.
	b := workerCmd(map[string]string{
		"ACT_MODE": "recover-approve", "ACT_RUN_ID": runID,
		"ACT_DSN": dsn, "ACT_SCHEMA": schema, "ACT_APP_VERSION": version, "ACT_STORE_SCHEMA": storeSchema,
	})
	out2, err2 := b.CombinedOutput()
	require.NoError(t, err2, "worker B should recover the parked run and finish cleanly; output:\n%s", out2)
	require.Contains(t, string(out2), "RECOVERED status=succeeded",
		"a run parked on the approval Recv must recover after a crash and finish once approved; output:\n%s", out2)
}

// TestExecutorRecoversApprovalWaitAfterGracefulShutdown is the production-shaped version of the crash-kill baseline: a
// worker drains via a GRACEFUL DBOS shutdown while a run is parked on the approval Recv (what a deploy does), then a
// fresh worker must recover the run and finish it once approved.
//
// This is the case that used to be lost. A graceful dbos.Shutdown cancels the parked workflow's context and the
// approval Recv returns context.Canceled; until DBOS v1.1 returning that error told the engine the run was
// unrecoverable, so every unsigned approval was written off as ERROR and recovery (which only re-runs PENDING
// workflows) never picked it up again. Since v1.1 the shutdown cancellation carries its own cause and a run unwinding
// under it skips the outcome write, so the row keeps its PENDING (upstream #423).
//
// That is what makes this test worth keeping rather than deleting with the workaround it once guarded: the behaviour
// it asserts now lives in a dependency, so this is where an upgrade that regressed it would be caught. It covers the
// whole chain, not one decision: the row survives the drain, recovery replays the checkpointed steps instead of
// re-running them, the re-entered Recv accepts a decision delivered afterwards, and the run finishes with its action
// performed. Its sibling TestExecutorRecoversApprovalWaitAfterCrash is the same assertion for a hard kill, which
// always worked; the pair separates "recovery is configured right" from "the shutdown path preserves the run".
func TestExecutorRecoversApprovalWaitAfterGracefulShutdown(t *testing.T) {
	dsn, schema := requirePostgres(t)

	version := "act-park-graceful-" + uuid.NewString()
	runID := "park-graceful-run-" + uuid.NewString()
	storeSchema := "act_park_graceful_store_" + uuid.NewString()[:8]

	// Phase 1: worker A starts the run and gracefully shuts down while it is parked on the approval Recv.
	a := workerCmd(map[string]string{
		"ACT_MODE": "start-park-graceful", "ACT_RUN_ID": runID,
		"ACT_DSN": dsn, "ACT_SCHEMA": schema, "ACT_APP_VERSION": version, "ACT_STORE_SCHEMA": storeSchema,
	})
	out, err := a.CombinedOutput()
	require.NoError(t, err, "worker A should shut down gracefully; output:\n%s", out)

	// Phase 2: worker B recovers the run and finishes it once the approval is delivered (same store as worker A so the
	// replayed step sequence matches).
	b := workerCmd(map[string]string{
		"ACT_MODE": "recover-approve", "ACT_RUN_ID": runID,
		"ACT_DSN": dsn, "ACT_SCHEMA": schema, "ACT_APP_VERSION": version, "ACT_STORE_SCHEMA": storeSchema,
	})
	out2, err2 := b.CombinedOutput()
	require.NoError(t, err2, "worker B should recover the parked run and finish cleanly; output:\n%s", out2)
	require.Contains(t, string(out2), "RECOVERED status=succeeded",
		"a run parked on approval must survive a graceful shutdown and finish once approved; output:\n%s", out2)
}

// awaitWaitingApproval blocks until the run reaches waiting_approval in the store, so a park worker can crash while the
// workflow is genuinely suspended on the approval Recv rather than at some earlier (or later) step.
func awaitWaitingApproval(store act.RunStore, instanceID, runID string) error {
	if store == nil {
		return fmt.Errorf("park worker requires a run store (set ACT_STORE_SCHEMA)")
	}
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		run, err := store.GetRun(context.Background(), instanceID, runID)
		if err == nil && run.Status == act.RunStatusWaitingApproval {
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("run %q did not reach waiting_approval within the deadline", runID)
}

// workerCmd re-execs this test binary as a DBOS worker subprocess with the given environment.
func workerCmd(env map[string]string) *exec.Cmd {
	cmd := exec.Command(os.Args[0])
	cmd.Env = append(os.Environ(), "ACT_TEST_WORKER=1")
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	return cmd
}

// stepLog is the parsed action log: one entry per step invocation.
type stepLog struct {
	steps []stepEntry
}

type stepEntry struct {
	pid  string
	step string
}

func (l stepLog) count(step string) int {
	n := 0
	for _, e := range l.steps {
		if e.step == step {
			n++
		}
	}
	return n
}

func (l stepLog) pid(step string) string {
	for _, e := range l.steps {
		if e.step == step {
			return e.pid
		}
	}
	return ""
}

func readSteps(t *testing.T, path string) stepLog {
	t.Helper()
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	var log stepLog
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := map[string]string{}
		for _, tok := range strings.Fields(sc.Text()) {
			if k, v, ok := strings.Cut(tok, "="); ok {
				fields[k] = v
			}
		}
		if fields["step"] != "" {
			log.steps = append(log.steps, stepEntry{pid: fields["pid"], step: fields["step"]})
		}
	}
	require.NoError(t, sc.Err())
	return log
}

package act_test

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// TestMain lets this test binary re-exec itself as a DBOS worker subprocess (ACT_TEST_WORKER=1), which the crash
// tests use to kill a worker mid-run and recover it in a fresh process. The env-var check runs before any test, so
// a worker subprocess never runs the test suite (and never forks further workers).
func TestMain(m *testing.M) {
	if os.Getenv("ACT_TEST_WORKER") == "1" {
		runGatewayCrashWorker() // never returns
	}
	os.Exit(m.Run())
}

// markStep appends one observable line to the step log, so the parent can see which steps ran, in which process.
func markStep(logPath, step, detail string) {
	line := fmt.Sprintf("ts=%s pid=%d step=%s detail=%s\n", time.Now().Format(time.RFC3339Nano), os.Getpid(), step, detail)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mark:", err)
		return
	}
	defer f.Close()
	_, _ = f.WriteString(line)
}

// TestParkedRunRecoversAfterCrash: a worker hard-killed while a run waits on its approval leaves the workflow PENDING,
// and a fresh worker recovers it, re-enters the wait, and finishes the run once the approval is delivered. This is
// the baseline; its siblings below are the shutdown paths that used to lose the run.
func TestParkedRunRecoversAfterCrash(t *testing.T) {
	out, steps := gatewayCrash(t, true,
		map[string]string{"ACT_MODE": "start", "ACT_THEN": "park-kill"},
		map[string]string{"ACT_MODE": "recover", "ACT_APPROVE": "tc-1"}, nil)
	require.Contains(t, out, "RECOVERED status=succeeded")
	require.Equal(t, 1, steps.count("execute"))
}

// TestParkedRunRecoversAfterGracefulShutdown: a worker that drains with a graceful DBOS shutdown while a run is parked
// (what a deploy does) must leave the run for the next process.
//
// This is the case that used to be lost: until DBOS v1.1 the Recv's cancellation told the engine the run was
// unrecoverable, and every unsigned approval was written off as ERROR (kairos-cloud#135, upstream #423). The behaviour
// now lives in the dependency, so this is where an upgrade that regressed it would be caught.
func TestParkedRunRecoversAfterGracefulShutdown(t *testing.T) {
	out, steps := gatewayCrash(t, false,
		map[string]string{"ACT_MODE": "start", "ACT_THEN": "park-graceful"},
		map[string]string{"ACT_MODE": "recover", "ACT_APPROVE": "tc-1"}, nil)
	require.Contains(t, out, "RECOVERED status=succeeded")
	require.Equal(t, 1, steps.count("execute"))
}

// TestParkedRunRecoversAfterHostCancelledShutdown reproduces how the runtime really stops: the host's context is
// cancelled by the signal first, and the executor is closed afterwards. DBOS spares a parked run only when its OWN
// Shutdown caused the cancellation; a host context that dies first delivers a plain context.Canceled and the run is
// durably cancelled. The executor detaches from the caller's cancellation to close that gap, and this test says so.
// It reads what the engine recorded between the two workers, so a regression names itself.
func TestParkedRunRecoversAfterHostCancelledShutdown(t *testing.T) {
	dsn, schema := requirePostgres(t)
	out, steps := gatewayCrash(t, false,
		map[string]string{"ACT_MODE": "start", "ACT_THEN": "park-sigterm"},
		map[string]string{"ACT_MODE": "recover", "ACT_APPROVE": "tc-1"},
		func(runID string) {
			require.Equal(t, "PENDING", workflowStatus(t, dsn, schema, runID),
				"the runtime's shutdown durably cancelled a parked run instead of leaving it for recovery")
		})
	require.Contains(t, out, "RECOVERED status=succeeded")
	require.Equal(t, 1, steps.count("execute"))
}

// workflowStatus reads a workflow's status straight from the DBOS system tables. It reaches past the executor on
// purpose: the whole question these shutdown tests ask is what the ENGINE recorded, which no Act-level API exposes.
func workflowStatus(t *testing.T, dsn, schema, workflowID string) string {
	t.Helper()
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	defer pool.Close()

	var status string
	err = pool.QueryRow(context.Background(),
		fmt.Sprintf("SELECT status FROM %s.workflow_status WHERE workflow_uuid = $1", schema), workflowID).Scan(&status)
	require.NoError(t, err, "no workflow row for %q", workflowID)
	return status
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

package act_test

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// The invariant these tests defend: a run's status and its approvals are two halves of one fact, and a reader must
// never see them disagree. A run that says waiting_approval must have something pending to decide; a run that reached
// a terminal state must have nothing.
//
// Production violated it in one direction: 5 of 8 runs waiting for approval had every approval cancelled, with no
// decider, in batches matching process restarts. Nobody could act on them and nothing explained why
// (kairos-cloud#135). The cause was two writes that could land independently, one of them on a context that outlived
// the shutdown the other died to.

// TestShutdownDuringApprovalWaitLeavesRunDecidable reproduces that production state directly: a run parked on a human
// approval while the process drains, which is what every deploy does. The run legitimately stays waiting_approval (it
// was not decided and did not fail), so its approvals MUST stay pending — otherwise the run is visible, undecidable
// and unexplained, forever.
//
// Before the fix this failed on the segmented path only: the unsegmented one returned the interrupt without sweeping,
// while the segmented one swept, so which path a run took decided whether a deploy stranded it.
func TestShutdownDuringApprovalWaitLeavesRunDecidable(t *testing.T) {
	store := newRunStore(t)
	gateway := &act.Gateway{
		Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative)),
		Ledger:   store,
		Executor: &fakeExecutor{},
	}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: &batchSegmentRunner{n: 2}, Store: store, Gateway: gateway, Proposer: batchProposer{},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	require.NoError(t, err)

	const inst = "inst-shutdown-decidable"
	ctx := t.Context()
	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "shutdown while parked",
		IdempotencyKey: "shutdown-parked-" + uuid.NewString(),
		Actor:          act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// Park the run: wait until both approvals exist and the run is actually waiting on the first one. Waiting on the
	// store rather than sleeping is what makes the shutdown land while the workflow is suspended mid-Recv.
	for _, tc := range []string{"tc-1", "tc-2"} {
		require.Eventually(t, func() bool {
			_, gErr := store.GetApproval(ctx, inst, act.ApprovalIDForToolCall(runID, tc))
			return gErr == nil
		}, 15*time.Second, 50*time.Millisecond)
	}
	require.Eventually(t, func() bool {
		run, gErr := store.GetRun(ctx, inst, runID)
		return gErr == nil && run.Status == act.RunStatusWaitingApproval
	}, 15*time.Second, 50*time.Millisecond)

	// Drain the executor, exactly as a deploy does. This cancels the workflow context out from under the approval wait.
	e.Close(10 * time.Second)

	// The sweep the bug performed ran on a context that outlived the shutdown, so it could land after Close returned:
	// assert the approvals stay pending over a window, not just at one instant.
	require.Never(t, func() bool {
		for _, tc := range []string{"tc-1", "tc-2"} {
			a, gErr := store.GetApproval(ctx, inst, act.ApprovalIDForToolCall(runID, tc))
			if gErr != nil {
				return true
			}
			if a.Status != act.ApprovalStatusPending {
				return true
			}
		}
		return false
	}, 3*time.Second, 200*time.Millisecond,
		"a shutdown must not withdraw the approvals of a run it leaves waiting: that is what stranded runs in production")

	// And the two halves must agree: the run still says it is waiting, and there is genuinely something to decide.
	run, err := store.GetRun(ctx, inst, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusWaitingApproval, run.Status)

	pending, err := store.ListApprovals(ctx, act.ListApprovalsFilter{InstanceID: inst, RunID: runID, Status: act.ApprovalStatusPending})
	require.NoError(t, err)
	require.NotEmpty(t, pending, "a run reading waiting_approval must have a pending approval a human can act on")
}

// TestCloseRunBindsTerminalStatusAndApprovalWithdrawal covers the other direction of the invariant at its source: the
// store write that marks a run terminal is the same write that withdraws its approvals, so no caller can perform one
// without the other and no reader can observe them apart.
func TestCloseRunBindsTerminalStatusAndApprovalWithdrawal(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const inst = "inst-close-run"
	runID := act.ComposeRunID(inst, "triage", "close-run")
	seedRun(t, store, inst, runID)
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: inst, RunID: runID, Status: act.RunStatusRunning, EventType: act.EventTypeRunning,
	}))
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: inst, RunID: runID, Status: act.RunStatusWaitingApproval, EventType: act.EventTypeWaitingApproval,
	}))
	approvalID := act.ApprovalIDForToolCall(runID, "tc-1")
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: approvalID, RunID: runID, InstanceID: inst,
		ToolName: "create_issue", ToolCallID: "tc-1",
		ArgsHash: act.HashArgs(`{"title":"x"}`), CanonicalArgs: `{"title":"x"}`,
	}))

	cancelled, err := store.CloseRun(ctx, act.RunTransition{
		InstanceID: inst, RunID: runID, Status: act.RunStatusFailed, EventType: act.EventTypeFailed,
		Error: "simulated",
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), cancelled)

	run, err := store.GetRun(ctx, inst, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusFailed, run.Status)

	approval, err := store.GetApproval(ctx, inst, approvalID)
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusCancelled, approval.Status)
	require.Empty(t, approval.DecidedBy, "a withdrawal is not a human decision, so it names no decider")

	// A non-terminal status is refused rather than silently withdrawing the approvals of a run that is still going.
	_, err = store.CloseRun(ctx, act.RunTransition{
		InstanceID: inst, RunID: runID, Status: act.RunStatusRunning, EventType: act.EventTypeRunning,
	})
	require.Error(t, err)

	// Closing an already-closed run stays safe to repeat: the transition is a no-op and nothing is left pending.
	_, err = store.CloseRun(ctx, act.RunTransition{
		InstanceID: inst, RunID: runID, Status: act.RunStatusCancelled, EventType: act.EventTypeCancelled,
	})
	require.NoError(t, err)
	run, err = store.GetRun(ctx, inst, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusFailed, run.Status, "a terminal run never transitions again")
}

// TestMigrateRepairsStrandedRuns covers the startup repair for the runs already written by the bug. The repair has to
// be safe to re-run on every boot, so it is asserted from both sides: it closes a stranded run and explains itself on
// the timeline, and it leaves a live one alone.
func TestMigrateRepairsStrandedRuns(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const inst = "inst-stranded"

	// The damaged shape: waiting_approval, approvals all withdrawn with no decider, untouched for hours.
	stranded := act.ComposeRunID(inst, "triage", "stranded")
	seedRun(t, store, inst, stranded)
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: inst, RunID: stranded, Status: act.RunStatusWaitingApproval, EventType: act.EventTypeWaitingApproval,
	}))
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: act.ApprovalIDForToolCall(stranded, "tc-1"), RunID: stranded, InstanceID: inst,
		ToolName: "create_issue", ToolCallID: "tc-1",
		ArgsHash: act.HashArgs(`{"title":"x"}`), CanonicalArgs: `{"title":"x"}`,
	}))
	_, err := store.CancelPendingApprovals(ctx, inst, stranded)
	require.NoError(t, err)
	require.NoError(t, store.BackdateRunForTest(ctx, inst, stranded, 3*time.Hour))
	require.NoError(t, store.BackdateApprovalDecisionsForTest(ctx, inst, stranded, 3*time.Hour))

	// A run parked legitimately: old enough to pass the age check, but with a decision still outstanding.
	live := act.ComposeRunID(inst, "triage", "live")
	seedRun(t, store, inst, live)
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: inst, RunID: live, Status: act.RunStatusWaitingApproval, EventType: act.EventTypeWaitingApproval,
	}))
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: act.ApprovalIDForToolCall(live, "tc-1"), RunID: live, InstanceID: inst,
		ToolName: "create_issue", ToolCallID: "tc-1",
		ArgsHash: act.HashArgs(`{"title":"y"}`), CanonicalArgs: `{"title":"y"}`,
	}))
	require.NoError(t, store.BackdateRunForTest(ctx, inst, live, 3*time.Hour))

	// A run that reached waiting_approval before any approval row was written: a create that raced a crash, not a
	// stranded run. It must be left alone, since the approval may still be on its way.
	bare := act.ComposeRunID(inst, "triage", "bare")
	seedRun(t, store, inst, bare)
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: inst, RunID: bare, Status: act.RunStatusWaitingApproval, EventType: act.EventTypeWaitingApproval,
	}))
	require.NoError(t, store.BackdateRunForTest(ctx, inst, bare, 3*time.Hour))

	require.NoError(t, store.Migrate(ctx))

	got, err := store.GetRun(ctx, inst, stranded)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusCancelled, got.Status, "a run with nothing left to decide must not stay in the inbox")

	// Closing it silently would repeat the original sin of changing state with no trace: the timeline must say why.
	events, err := store.ListRunEvents(ctx, inst, stranded, 0, 0)
	require.NoError(t, err)
	var explained bool
	for _, ev := range events {
		if ev.EventType == act.EventTypeCancelled {
			explained = ev.Payload["reason"] != nil
		}
	}
	require.True(t, explained, "the repair must record a cancelled event carrying its reason")

	for _, untouched := range []string{live, bare} {
		got, err = store.GetRun(ctx, inst, untouched)
		require.NoError(t, err)
		require.Equal(t, act.RunStatusWaitingApproval, got.Status,
			"the repair must only close runs with no decision outstanding")
	}

	// Idempotent: a second boot finds nothing new and appends no second event.
	require.NoError(t, store.Migrate(ctx))
	after, err := store.ListRunEvents(ctx, inst, stranded, 0, 0)
	require.NoError(t, err)
	require.Len(t, after, len(events), "re-running the repair must not append a second cancelled event")
}

// TestMigrateLeavesRecentlyActiveRunsAlone pins the age guard, which is the whole reason the repair is safe to run on
// every startup: a live run passes through "waiting_approval with every approval decided" for the instants between a
// decision and the resumed transition, and closing it there would kill a run mid-flight.
func TestMigrateLeavesRecentlyActiveRunsAlone(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const inst = "inst-recent"
	runID := act.ComposeRunID(inst, "triage", "recent")
	seedRun(t, store, inst, runID)
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: inst, RunID: runID, Status: act.RunStatusWaitingApproval, EventType: act.EventTypeWaitingApproval,
	}))
	approvalID := act.ApprovalIDForToolCall(runID, "tc-1")
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: approvalID, RunID: runID, InstanceID: inst,
		ToolName: "create_issue", ToolCallID: "tc-1",
		ArgsHash: act.HashArgs(`{"title":"z"}`), CanonicalArgs: `{"title":"z"}`,
	}))
	// Decided just now: the run is between the decision and its resumed transition.
	_, err := store.ResolveApproval(ctx, inst, approvalID, act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)

	require.NoError(t, store.Migrate(ctx))

	got, err := store.GetRun(ctx, inst, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusWaitingApproval, got.Status,
		"a run that was active moments ago is mid-flight, not stranded")
}

// TestMigrateLeavesLongParkedRunApprovedJustNowAlone is the case an age check on the RUN row silently gets wrong, and
// the reason the predicate reads staleness off the approvals instead. A run can sit waiting for days — so its
// updated_on is ancient — and then be approved one second before a replica boots. Judged by the run row it looks
// stranded; judged by the decision it is a run that just came alive. Cancelling it here would kill live work, and the
// symptom would look exactly like the bug being repaired.
func TestMigrateLeavesLongParkedRunApprovedJustNowAlone(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const inst = "inst-long-parked"
	runID := act.ComposeRunID(inst, "triage", "long-parked")
	seedRun(t, store, inst, runID)
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: inst, RunID: runID, Status: act.RunStatusWaitingApproval, EventType: act.EventTypeWaitingApproval,
	}))
	// An earlier approval that really was withdrawn by the process, so this run carries the bug's exact fingerprint:
	// every other clause of the predicate says "stranded", and only the age of the last decision says otherwise.
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: act.ApprovalIDForToolCall(runID, "tc-2"), RunID: runID, InstanceID: inst,
		ToolName: "create_issue", ToolCallID: "tc-2",
		ArgsHash: act.HashArgs(`{"title":"withdrawn"}`), CanonicalArgs: `{"title":"withdrawn"}`,
	}))
	_, err := store.CancelPendingApprovals(ctx, inst, runID)
	require.NoError(t, err)
	require.NoError(t, store.BackdateApprovalDecisionsForTest(ctx, inst, runID, 72*time.Hour))
	require.NoError(t, store.BackdateRunForTest(ctx, inst, runID, 72*time.Hour))

	// The run is parked again on a new approval, and the decision lands right now.
	approvalID := act.ApprovalIDForToolCall(runID, "tc-1")
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: approvalID, RunID: runID, InstanceID: inst,
		ToolName: "create_issue", ToolCallID: "tc-1",
		ArgsHash: act.HashArgs(`{"title":"parked"}`), CanonicalArgs: `{"title":"parked"}`,
	}))
	_, err = store.ResolveApproval(ctx, inst, approvalID, act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)

	require.NoError(t, store.Migrate(ctx))

	got, err := store.GetRun(ctx, inst, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusWaitingApproval, got.Status,
		"a run approved a moment ago is alive, however long it had been parked before that")
}

// TestMigrateLeavesFullyDecidedBatchesAlone pins the other half of the predicate: the repair targets the fingerprint
// of THIS bug (an approval withdrawn by the process, with no decider), not merely "a run whose approvals are all
// resolved". A batch a human decided in full is a different situation with a different cause, and closing it here
// would be guessing.
func TestMigrateLeavesFullyDecidedBatchesAlone(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const inst = "inst-decided"
	runID := act.ComposeRunID(inst, "triage", "decided")
	seedRun(t, store, inst, runID)
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: inst, RunID: runID, Status: act.RunStatusWaitingApproval, EventType: act.EventTypeWaitingApproval,
	}))
	for i, decision := range []string{act.ApprovalStatusApproved, act.ApprovalStatusDenied} {
		tc := fmt.Sprintf("tc-%d", i+1)
		require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
			ApprovalID: act.ApprovalIDForToolCall(runID, tc), RunID: runID, InstanceID: inst,
			ToolName: "create_issue", ToolCallID: tc,
			ArgsHash: act.HashArgs(`{"title":"d"}`), CanonicalArgs: `{"title":"d"}`,
		}))
		_, err := store.ResolveApproval(ctx, inst, act.ApprovalIDForToolCall(runID, tc), decision, "admin:bob")
		require.NoError(t, err)
	}
	require.NoError(t, store.BackdateRunForTest(ctx, inst, runID, 3*time.Hour))
	require.NoError(t, store.BackdateApprovalDecisionsForTest(ctx, inst, runID, 3*time.Hour))

	require.NoError(t, store.Migrate(ctx))

	got, err := store.GetRun(ctx, inst, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusWaitingApproval, got.Status,
		"every approval decided by a human is not this bug's fingerprint, so the repair must not touch it")
}

// TestMigrateWithdrawsApprovalsLeftOnTerminalRuns covers the inconsistency from the opposite side: an approval still
// pending on a run that is already over. A DBOS step checkpointed by the previous binary replays its recorded output
// and so skips the body that now withdraws approvals, which makes this reachable during the very deploy that ships
// the fix.
func TestMigrateWithdrawsApprovalsLeftOnTerminalRuns(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const inst = "inst-terminal-pending"
	runID := act.ComposeRunID(inst, "triage", "terminal-pending")
	seedRun(t, store, inst, runID)
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: inst, RunID: runID, Status: act.RunStatusWaitingApproval, EventType: act.EventTypeWaitingApproval,
	}))
	approvalID := act.ApprovalIDForToolCall(runID, "tc-1")
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: approvalID, RunID: runID, InstanceID: inst,
		ToolName: "create_issue", ToolCallID: "tc-1",
		ArgsHash: act.HashArgs(`{"title":"t"}`), CanonicalArgs: `{"title":"t"}`,
	}))
	// Mark the run terminal the way the old binary did: status only, approvals untouched.
	require.NoError(t, store.ForceTerminalWithoutCloseForTest(ctx, inst, runID, act.RunStatusFailed))
	require.NoError(t, store.BackdateRunForTest(ctx, inst, runID, 3*time.Hour))

	require.NoError(t, store.Migrate(ctx))

	approval, err := store.GetApproval(ctx, inst, approvalID)
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusCancelled, approval.Status,
		"a finished run must not keep offering a decision that can no longer change anything")
}

// TestRecordTransitionRefusesTerminalStatus pins the door CloseRun exists to close. Leaving RecordTransition able to
// write a terminal status would keep the invariant a matter of picking the right method name, which is how the
// original bug got written in the first place.
func TestRecordTransitionRefusesTerminalStatus(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const inst = "inst-refuse-terminal"
	runID := act.ComposeRunID(inst, "triage", "refuse-terminal")
	seedRun(t, store, inst, runID)

	for _, status := range []act.RunStatus{act.RunStatusSucceeded, act.RunStatusFailed, act.RunStatusRejected, act.RunStatusCancelled} {
		err := store.RecordTransition(ctx, act.RunTransition{
			InstanceID: inst, RunID: runID, Status: status, EventType: act.EventTypeFailed,
		})
		require.Error(t, err, "RecordTransition must refuse the terminal status %q", status)
		require.Contains(t, err.Error(), "CloseRun")
	}
}

// TestCreateApprovalRefusesTerminalRun closes the race in the other direction: CloseRun withdraws what is pending and
// commits, and an approval insert already in flight lands afterwards, leaving a finished run with a live approval.
func TestCreateApprovalRefusesTerminalRun(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const inst = "inst-approval-after-close"
	runID := act.ComposeRunID(inst, "triage", "approval-after-close")
	seedRun(t, store, inst, runID)
	_, err := store.CloseRun(ctx, act.RunTransition{
		InstanceID: inst, RunID: runID, Status: act.RunStatusCancelled, EventType: act.EventTypeCancelled,
	})
	require.NoError(t, err)

	err = store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: act.ApprovalIDForToolCall(runID, "tc-1"), RunID: runID, InstanceID: inst,
		ToolName: "create_issue", ToolCallID: "tc-1",
		ArgsHash: act.HashArgs(`{"title":"late"}`), CanonicalArgs: `{"title":"late"}`,
	})
	require.Error(t, err, "an approval must not attach to a run that is already over")

	pending, err := store.ListApprovals(ctx, act.ListApprovalsFilter{InstanceID: inst, RunID: runID, Status: act.ApprovalStatusPending})
	require.NoError(t, err)
	require.Empty(t, pending)
}

// TestNoExpirationVocabularyRemains guards the cleanup that removed the approval-expiry mechanism's leftovers. They
// were dead statements that still read as a live feature, and they cost a wrong diagnosis of the bug above: the
// expired->cancelled normalization looked like a plausible cause of the cancelled approvals until the dates ruled it
// out. Keeping the vocabulary out of the store keeps the next reader from making the same detour.
// It asserts against the SOURCE, not against behaviour, deliberately: dead vocabulary costs nothing at runtime and so
// cannot be caught by exercising the store. What it costs is a reader's time, and the only way to pin that is to check
// that the words are gone from the files a reader reads.
func TestNoExpirationVocabularyRemains(t *testing.T) {
	require.False(t, act.RunStatus("expired").IsTerminal(), "expired is not part of the run vocabulary")

	// The word survives legitimately in two places, so neither is banned: the execution-lease code, where a lease
	// really does expire, and the one ALTER TABLE that DROPS the retired column, which must stay for any database
	// that has not run it yet. What is banned is vocabulary that implies a live expiry mechanism.
	banned := regexp.MustCompile(`status\s*=\s*'expired'|run\.expired|gives expiration|expired_at|ExpiresOn`)
	for _, name := range []string{"pgstore.go", "store.go", "dbos_executor.go", "executor.go", "README.md"} {
		body, err := os.ReadFile(name)
		require.NoError(t, err)
		require.NotRegexp(t, banned, string(body),
			"%s still refers to the removed approval-expiry mechanism; that leftover cost one wrong diagnosis already", name)
	}
}

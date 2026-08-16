package server_test

import (
	"testing"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/pkg/ratelimit"
	"github.com/rilldata/rill/runtime/server"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Launching an agent by hand used to demand a prompt, which meant retyping — or guessing — text already written in
// the agent's YAML, and getting it wrong meant the manual run did something subtly different from the scheduled one.
// Choosing the agent already says what to run: the task is in its instructions (the system message), and the trigger
// it declares carries the user turn that goes with them.

// newActServerWithFiles is newActServer with a caller-supplied project, so a test can declare agents with triggers.
func newActServerWithFiles(t *testing.T, files map[string]string) (*server.Server, *fakeAgentExecutor, string) {
	t.Helper()
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{Files: files})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)

	srv, err := server.NewServer(context.Background(), &server.Options{}, rt, zap.NewNop(), ratelimit.NewNoop(), activity.NewNoopClient())
	require.NoError(t, err)

	store := newActStore(t)
	exec := &fakeAgentExecutor{store: store}
	srv.ConfigureAct(store, exec)
	return srv, exec, instanceID
}

const triggerPrompt = "Compila el resumen de la ultima semana completa y propone enviarlo."

func scheduledAgentFiles() map[string]string {
	return map[string]string{
		"rill.yaml": "features:\n  agents: true\n",
		"weekly.yaml": `
type: agent
display_name: Weekly Summary
instructions: "Sos el analista que resume la semana."
tools: []
limits:
  max_steps: 3
triggers:
  - source:
      kind: schedule
      cron: "0 8 * * 1"
    input:
      prompt: "` + triggerPrompt + `"
`,
	}
}

// TestStartRunWithoutPromptUsesTheTriggersOwnPrompt is the point of the change: an empty prompt means "do now what
// you would do on Monday", sending the very text the schedule sends.
func TestStartRunWithoutPromptUsesTheTriggersOwnPrompt(t *testing.T) {
	srv, exec, instanceID := newActServerWithFiles(t, scheduledAgentFiles())
	ctx := testCtx()

	_, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{
		InstanceId: instanceID, Name: "weekly", IdempotencyKey: "k1",
	})
	require.NoError(t, err)
	require.Equal(t, triggerPrompt, exec.lastStart().Prompt)

	// The text is reused; the provenance is not. The run is manual because a person launched it, and an audit trail
	// that let a caller borrow "schedule" would be forgeable.
	require.Equal(t, "manual", exec.lastStart().Trigger)
}

// TestStartRunWithPromptReplacesRatherThanAppends pins the injection-safe half of the design. The trigger prompt is
// written by the agent's author; this text is written by whoever can launch the agent. Appending one to the other
// would put a launcher's words beside an author's inside a single user turn, where either can pass for the other.
// Substituting keeps them from ever mixing, and the agent's instructions stay untouched in the system message.
func TestStartRunWithPromptReplacesRatherThanAppends(t *testing.T) {
	srv, exec, instanceID := newActServerWithFiles(t, scheduledAgentFiles())
	ctx := testCtx()

	const mine = "Centrate en devoluciones."
	_, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{
		InstanceId: instanceID, Name: "weekly", IdempotencyKey: "k2", Prompt: mine,
	})
	require.NoError(t, err)

	sent := exec.lastStart().Prompt
	require.Equal(t, mine, sent)
	require.NotContains(t, sent, triggerPrompt, "the caller's prompt must replace the trigger's, never be concatenated with it")
}

// TestStartRunWithoutPromptNeedsOneWhenTheAgentHasNoTrigger covers the case where there is genuinely no text to
// reuse: an agent meant to be launched by hand. Refusing with a reason beats sending an empty user turn, which the
// instructions are not written to answer.
func TestStartRunWithoutPromptNeedsOneWhenTheAgentHasNoTrigger(t *testing.T) {
	srv, _, _, instanceID := newActServer(t)
	ctx := testCtx()

	_, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{
		InstanceId: instanceID, Name: "triage", IdempotencyKey: "k3",
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Contains(t, status.Convert(err).Message(), "neither a prompt of its own nor a trigger")
}

// TestStartRunWithoutPromptRefusesToGuessBetweenTriggers covers the ambiguous case. Two triggers ask the agent for
// different things on different occasions, so there is no "the" usual task to fall back on. Picking one silently
// would run the wrong job and look like it worked.
func TestStartRunWithoutPromptRefusesToGuessBetweenTriggers(t *testing.T) {
	files := map[string]string{
		"rill.yaml": "features:\n  agents: true\n",
		"dual.yaml": `
type: agent
display_name: Two Jobs
instructions: "Sos un analista polivalente."
tools: []
limits:
  max_steps: 3
triggers:
  - source:
      kind: schedule
      cron: "0 8 * * 1"
    input:
      prompt: "Resumen semanal."
  - source:
      kind: schedule
      cron: "0 8 1 * *"
    input:
      prompt: "Cierre mensual."
`,
	}
	srv, _, instanceID := newActServerWithFiles(t, files)
	ctx := testCtx()

	_, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{
		InstanceId: instanceID, Name: "dual", IdempotencyKey: "k4",
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Contains(t, status.Convert(err).Message(), "several triggers")
}

// TestStartRunUsesTheAgentsOwnPromptWhenItHasNoTrigger covers the agent built to be launched by hand. It has no
// trigger to borrow a user turn from, and demanding one from whoever presses the button would be the same friction
// the trigger fallback exists to remove: retyping a request the author already knows how to word. Declaring it on
// the agent means choosing the agent is enough to run it.
func TestStartRunUsesTheAgentsOwnPromptWhenItHasNoTrigger(t *testing.T) {
	const ownPrompt = "Revisa las devoluciones del ultimo mes y propone un ticket."
	files := map[string]string{
		"rill.yaml": "features:\n  agents: true\n",
		"manual.yaml": `
type: agent
display_name: Manual Only
instructions: "Sos el analista de devoluciones."
prompt: "` + ownPrompt + `"
tools: []
limits:
  max_steps: 3
`,
	}
	srv, exec, instanceID := newActServerWithFiles(t, files)
	ctx := testCtx()

	_, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{
		InstanceId: instanceID, Name: "manual", IdempotencyKey: "k5",
	})
	require.NoError(t, err)
	require.Equal(t, ownPrompt, exec.lastStart().Prompt)
}

// TestAgentsOwnPromptWinsOverItsTriggers is the tie-breaker. An agent with several triggers is otherwise refused,
// because each asks for something different and picking one would silently run the wrong job; declaring a prompt
// says what a manual launch means without touching what the triggers do on their own occasions.
func TestAgentsOwnPromptWinsOverItsTriggers(t *testing.T) {
	const ownPrompt = "Corre la revision completa."
	files := map[string]string{
		"rill.yaml": "features:\n  agents: true\n",
		"both.yaml": `
type: agent
display_name: Both
instructions: "Sos un analista polivalente."
prompt: "` + ownPrompt + `"
tools: []
limits:
  max_steps: 3
triggers:
  - source:
      kind: schedule
      cron: "0 8 * * 1"
    input:
      prompt: "Resumen semanal."
  - source:
      kind: schedule
      cron: "0 8 1 * *"
    input:
      prompt: "Cierre mensual."
`,
	}
	srv, exec, instanceID := newActServerWithFiles(t, files)
	ctx := testCtx()

	_, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{
		InstanceId: instanceID, Name: "both", IdempotencyKey: "k6",
	})
	require.NoError(t, err)
	require.Equal(t, ownPrompt, exec.lastStart().Prompt,
		"a declared prompt settles what a manual launch means, where the triggers alone are ambiguous")
}

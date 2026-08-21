package act_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/stretchr/testify/require"
)

// TestSessionFactoryFailsClosedOnNilClaims is the security-critical case: the production factory must refuse a run
// that arrives without an initiator identity (nil claims) rather than fabricate a privileged session. This is the
// fail-closed half of the Actor contract — an automatic trigger whose service principal is not yet resolved (§17.3)
// must not escalate to the worker's authority. It errors before touching the runtime, so no session is opened.
func TestSessionFactoryFailsClosedOnNilClaims(t *testing.T) {
	rt, _ := newInstanceWithAgent(t, "Investiga la alerta.")
	factory := act.NewSessionFactory(rt, activity.NewNoopClient())

	s, release, err := factory(t.Context(), "any-instance", "", "", nil)
	require.ErrorIs(t, err, act.ErrNoInitiatorClaims)
	require.Nil(t, s)
	require.Nil(t, release)
}

// TestSessionFactoryOpensSessionWithClaims confirms the happy path: given the initiator's claims the factory opens a
// session bound to them. It passes the claims through faithfully — even a local-dev SkipChecks identity, which
// reflects the runtime's own auth posture rather than a fabricated escalation; only nil is refused.
func TestSessionFactoryOpensSessionWithClaims(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga la alerta.")
	factory := act.NewSessionFactory(rt, activity.NewNoopClient())

	s, release, err := factory(t.Context(), instanceID, "", "", &runtime.SecurityClaims{UserID: "u-1", SkipChecks: true})
	require.NoError(t, err)
	require.NotNil(t, s)
	require.NotNil(t, release)
	release()
}

// TestRunSegmentRoutesTheAgentsOwnConnector is the regression test for the defect this routing fixes: an agent's
// `model.connector` used to be parsed, validated by the reconciler and carried on the snapshot, and then dropped on the
// floor at execution, so every agent ran on rill.yaml's ai_connector no matter what it declared. The run here declares
// "gemini_lite"; before the fix the factory received "" and the assertion below failed. It asserts on the snapshot the
// segment was bound to, which is what makes a resumed segment stay on the model it started on.
func TestRunSegmentRoutesTheAgentsOwnConnector(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga la alerta.")

	var got string
	var calls int
	runner := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(rt),
		Sessions: func(ctx context.Context, instanceID, sessionID, modelConnector string, claims *runtime.SecurityClaims) (*ai.Session, func(), error) {
			got = modelConnector
			calls++
			// Refuse to open a session: the assertion is about what the runner passes down, and failing here keeps the
			// test off the model/tool loop, which needs a scripted LLM this case has no use for.
			return nil, nil, errors.New("stop after routing")
		},
	}

	_, err := runner.RunSegment(t.Context(), act.RunSegmentInput{
		InstanceID: instanceID,
		Claims:     &runtime.SecurityClaims{UserID: "u-1", SkipChecks: true},
		Snapshot:   &ai.AgentSnapshot{Name: "triage", ModelConnector: "gemini_lite"},
		Prompt:     "Investiga la alerta.",
	})
	require.Error(t, err)
	require.Equal(t, 1, calls)
	require.Equal(t, "gemini_lite", got)
}

// TestRunSegmentKeepsTheInstanceDefaultWhenTheAgentNamesNoConnector pins the additive half of the change: an agent with
// no model block must still resolve its model from rill.yaml's ai_connector, exactly as before. An empty string is the
// signal for that fallback, so the runner must pass it through rather than substituting a default of its own.
func TestRunSegmentKeepsTheInstanceDefaultWhenTheAgentNamesNoConnector(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga la alerta.")

	got := "unset"
	runner := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(rt),
		Sessions: func(ctx context.Context, instanceID, sessionID, modelConnector string, claims *runtime.SecurityClaims) (*ai.Session, func(), error) {
			got = modelConnector
			return nil, nil, errors.New("stop after routing")
		},
	}

	_, err := runner.RunSegment(t.Context(), act.RunSegmentInput{
		InstanceID: instanceID,
		Claims:     &runtime.SecurityClaims{UserID: "u-1", SkipChecks: true},
		Snapshot:   &ai.AgentSnapshot{Name: "triage"},
		Prompt:     "Investiga la alerta.",
	})
	require.Error(t, err)
	require.Empty(t, got)
}

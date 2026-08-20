package act_test

import (
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

	s, release, err := factory(t.Context(), "any-instance", "", nil, &ai.AgentSnapshot{})
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

	s, release, err := factory(t.Context(), instanceID, "", &runtime.SecurityClaims{UserID: "u-1", SkipChecks: true}, &ai.AgentSnapshot{
		ModelConnector: "openai", ModelDriver: "openai", ModelProperties: map[string]any{},
	})
	require.NoError(t, err)
	require.NotNil(t, s)
	require.NotNil(t, release)
	release()
}

func TestSessionFactoryFailsClosedOnLegacyModelSnapshot(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga la alerta.")
	factory := act.NewSessionFactory(rt, activity.NewNoopClient())

	s, release, err := factory(t.Context(), instanceID, "", &runtime.SecurityClaims{UserID: "u-1", SkipChecks: true}, &ai.AgentSnapshot{
		// Old checkpoints may carry the declared name but not the resolved driver/properties added by this release.
		ModelConnector: "deepseek",
	})
	require.ErrorIs(t, err, act.ErrIncompleteModelSnapshot)
	require.Nil(t, s)
	require.Nil(t, release)
}

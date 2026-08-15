package act

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// errKillSwitch is a switch whose state cannot be established.
type errKillSwitch struct{}

func (errKillSwitch) Disabled(context.Context, KillScope) (bool, string, error) {
	return false, "", errors.New("switch unreachable")
}

func TestMatchKillSpec(t *testing.T) {
	scope := KillScope{
		InstanceID: "inst",
		AgentName:  "gestion-cobranza",
		Tool:       "mcp.crm.create_crm_task",
		Connector:  "crm",
	}

	t.Run("unset lets everything through", func(t *testing.T) {
		killed, reason, err := matchKillSpec("", scope)
		require.NoError(t, err)
		require.False(t, killed)
		require.Empty(t, reason)
	})

	t.Run("star kills the whole project", func(t *testing.T) {
		killed, _, err := matchKillSpec("*", scope)
		require.NoError(t, err)
		require.True(t, killed)
	})

	t.Run("by agent", func(t *testing.T) {
		killed, reason, err := matchKillSpec("agent:gestion-cobranza", scope)
		require.NoError(t, err)
		require.True(t, killed)
		require.Contains(t, reason, "gestion-cobranza")

		killed, _, err = matchKillSpec("agent:otro", scope)
		require.NoError(t, err)
		require.False(t, killed)
	})

	t.Run("by connector", func(t *testing.T) {
		killed, _, err := matchKillSpec("connector:crm", scope)
		require.NoError(t, err)
		require.True(t, killed)
	})

	// The operator copies a tool name out of the audit trail, where it may appear either way. Both must work,
	// or the switch fails exactly when someone is depending on it.
	t.Run("by tool, raw or effective name", func(t *testing.T) {
		killed, _, err := matchKillSpec("tool:create_crm_task", scope)
		require.NoError(t, err)
		require.True(t, killed, "the raw name the MCP server exposes must match")

		killed, _, err = matchKillSpec("tool:mcp.crm.create_crm_task", scope)
		require.NoError(t, err)
		require.True(t, killed, "the effective name the ledger records must match too")

		killed, _, err = matchKillSpec("tool:delete_account", scope)
		require.NoError(t, err)
		require.False(t, killed)
	})

	t.Run("a list matches on any entry, and tolerates spacing", func(t *testing.T) {
		killed, _, err := matchKillSpec("agent:otro , connector:jira , agent:gestion-cobranza", scope)
		require.NoError(t, err)
		require.True(t, killed)

		killed, _, err = matchKillSpec("agent:otro,connector:jira", scope)
		require.NoError(t, err)
		require.False(t, killed)
	})

	// A typo must not read as "nothing is killed": the operator setting this is mid-incident and will believe
	// the writes stopped. The gateway turns the error into an engaged switch.
	t.Run("a malformed entry errors, which the gateway reads as engaged", func(t *testing.T) {
		for _, spec := range []string{"gestion-cobranza", "agent:", "agente:x", "tool", "  :x"} {
			_, _, err := matchKillSpec(spec, scope)
			require.Error(t, err, "spec %q must not pass silently", spec)
		}
	})
}

// TestGatewayKillSwitchFailsClosed pins the contract the switch depends on: an error from the switch stops the
// action. If this ever inverts, every malformed kill spec becomes a silent allow.
func TestGatewayKillSwitchFailsClosed(t *testing.T) {
	g := &Gateway{KillSwitch: errKillSwitch{}}
	killed, reason := g.killDisabled(t.Context(), KillScope{InstanceID: "inst"})
	require.True(t, killed)
	require.NotEmpty(t, reason)

	// And a gateway with no switch wired proceeds normally, which is what every test double relies on.
	none := &Gateway{}
	killed, _ = none.killDisabled(t.Context(), KillScope{InstanceID: "inst"})
	require.False(t, killed)
}

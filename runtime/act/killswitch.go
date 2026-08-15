package act

import (
	"context"
	"fmt"
	"strings"

	"github.com/rilldata/rill/runtime"
)

// KillVariable is the project variable that engages the kill switch. It is read from the project's variables
// rather than from a table of our own because that surface already exists and is the one an operator can reach
// without git and without a redeploy: `rill env set KILL_ACT ...`, per project and per environment.
//
// This is the fine-grained half of stopping Act. The coarse half is the `agents` feature flag, which closes the
// whole Act API for a project: nobody can even read past runs. The kill switch instead stops actions while
// leaving the history readable, which is what an incident actually calls for — halt the writes, keep the audit
// trail. Use the flag to turn the pillar off, this to stop it acting.
const KillVariable = "KILL_ACT"

// killAll is the value that engages the switch for everything in the project.
const killAll = "*"

// variableKillSwitch implements KillSwitch by reading KillVariable from the project's variables.
type variableKillSwitch struct {
	rt *runtime.Runtime
}

// NewVariableKillSwitch returns the production KillSwitch: the one that reads the project variable.
func NewVariableKillSwitch(rt *runtime.Runtime) KillSwitch {
	return &variableKillSwitch{rt: rt}
}

// Disabled reports whether the scope is killed. An unreadable instance is an error, which the gateway treats as
// engaged (fail closed), so a switch whose state cannot be established never lets an action through.
func (k *variableKillSwitch) Disabled(ctx context.Context, scope KillScope) (bool, string, error) {
	inst, err := k.rt.Instance(ctx, scope.InstanceID)
	if err != nil {
		return false, "", fmt.Errorf("act: read kill switch: %w", err)
	}
	// Lowercased keys so the variable matches however the operator cased it when setting it.
	spec := inst.ResolveVariables(true)[strings.ToLower(KillVariable)]
	return matchKillSpec(spec, scope)
}

// matchKillSpec decides whether a kill spec covers a scope. It is separate from the runtime lookup so the
// grammar can be tested on its own.
//
// The spec is a comma-separated list. "*" covers everything in the project; otherwise each entry is
// "agent:<name>", "tool:<name>" or "connector:<name>". A tool matches on the raw name the MCP server exposes
// (create_crm_task) as well as on the effective one the ledger records (mcp.crm.create_crm_task), because an
// operator reaching for this is usually copying a name out of the audit trail and should not have to know
// which form they are holding.
//
// A malformed entry is an error, and an error means engaged. That is deliberate: a typo in a kill switch that
// silently did nothing would be the single worst failure this code could have, since the operator typing it is
// mid-incident and will believe the writes have stopped.
func matchKillSpec(spec string, scope KillScope) (bool, string, error) {
	for _, entry := range strings.Split(spec, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if entry == killAll {
			return true, "act is disabled for this project", nil
		}
		kind, value, found := strings.Cut(entry, ":")
		value = strings.TrimSpace(value)
		if !found || value == "" {
			return false, "", fmt.Errorf("act: invalid %s entry %q: want %q, \"agent:<name>\", \"tool:<name>\" or \"connector:<name>\"", KillVariable, entry, killAll)
		}
		switch strings.TrimSpace(kind) {
		case "agent":
			if scope.AgentName == value {
				return true, fmt.Sprintf("agent %q is disabled", value), nil
			}
		case "tool":
			if scope.Tool == value || rawToolName(scope.Tool, scope.Connector) == value {
				return true, fmt.Sprintf("tool %q is disabled", value), nil
			}
		case "connector":
			if scope.Connector == value {
				return true, fmt.Sprintf("connector %q is disabled", value), nil
			}
		default:
			return false, "", fmt.Errorf("act: invalid %s scope %q in %q: want \"agent\", \"tool\" or \"connector\"", KillVariable, kind, entry)
		}
	}
	return false, "", nil
}

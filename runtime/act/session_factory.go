package act

import (
	"context"
	"errors"

	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/pkg/activity"
)

// ErrNoInitiatorClaims is returned by the production SessionFactory when a run reaches the model/tool loop without
// the initiating actor's SecurityClaims. It is the fail-closed half of the Actor contract (executor.go): a run must
// execute with exactly the authority of whoever (or whatever) started it, so a nil-claims run — e.g. an automatic
// trigger whose service principal's claims are not yet resolved (§17.3) — must never fall back to a privileged
// session. The run fails instead of escalating.
var ErrNoInitiatorClaims = errors.New("act: run has no initiator claims; refusing to open a privileged session (fail closed)")

// ErrIncompleteModelSnapshot protects upgrades from silently rerouting a durable run. Snapshots checkpointed by a
// build from before per-agent connector capture have no ModelDriver (and may have no effective ModelConnector). They
// cannot be reconstructed deterministically after a restart, so they must fail closed and be started again rather
// than falling through to whichever instance-default model happens to be live after the deploy.
var ErrIncompleteModelSnapshot = errors.New("act: run has no durable model connector snapshot; start a new run")

// NewSessionFactory returns the production SessionFactory: it opens an ai.Session bound to the run's initiating
// claims, so every tool CheckAccess in the agent loop resolves against the initiator's authority and never the
// worker's. It differs from the test factory in one security-critical way: on nil claims it does NOT fabricate a
// SkipChecks session, it fails closed (ErrNoInitiatorClaims). The claims are passed through exactly as received —
// including a legitimately SkipChecks-bearing local-dev identity, which reflects the runtime's own auth posture
// rather than a fabricated escalation. The session's model resolves from the immutable connector snapshot captured
// at run start; only that connector's current secrets are read live.
func NewSessionFactory(rt *runtime.Runtime, ac *activity.Client) SessionFactory {
	// One Runner is built up front and reused across sessions: it only registers the tool set and holds no per-run
	// state, so sharing it avoids re-registering tools on every run open.
	runner := ai.NewRunner(rt, ac)
	return func(ctx context.Context, instanceID, sessionID string, claims *runtime.SecurityClaims, snapshot *ai.AgentSnapshot) (*ai.Session, func(), error) {
		// Fail closed: a run with no initiator identity (a not-yet-wired automatic trigger, §17.3) must not run under
		// a privileged session. Refuse rather than escalate to the worker's own authority.
		if claims == nil {
			return nil, nil, ErrNoInitiatorClaims
		}
		if snapshot == nil {
			return nil, nil, errors.New("act: run has no agent snapshot")
		}
		if snapshot.ModelConnector == "" || snapshot.ModelDriver == "" {
			return nil, nil, ErrIncompleteModelSnapshot
		}
		// sessionID empty opens a fresh session (first segment); non-empty reopens that session so a resumed segment
		// continues the same conversation. Runner.Session loads the existing AISession and its messages (ordered by
		// index) when SessionID is set, which is the durable substrate the resumed segment reconstructs from.
		s, err := runner.Session(ctx, &ai.SessionOptions{
			InstanceID:    instanceID,
			SessionID:     sessionID,
			Claims:        claims,
			UserAgent:     "act-runtime",
			LLMConnector:  snapshot.ModelConnector,
			LLMDriver:     snapshot.ModelDriver,
			LLMProperties: snapshot.ModelProperties,
		})
		if err != nil {
			return nil, nil, err
		}
		// The session holds no resources of its own that outlive the run step, so the release is a no-op; the run
		// step's context governs the underlying work.
		return s, func() {}, nil
	}
}

// denyActionSink is a fail-closed ActionSink for deployments that run agents (read/investigate) but have not wired a
// real action gateway. Reaching it means a run was approved for an external write with no execution path configured;
// it refuses rather than silently dropping the write or performing an unmediated one.
type denyActionSink struct{}

var _ ActionSink = denyActionSink{}

func (denyActionSink) Apply(context.Context, ActionRequest) (ActionResult, error) {
	return ActionResult{}, errors.New("act: no action sink configured; refusing to perform an external write (fail closed)")
}

// NewDenyActionSink returns a fail-closed ActionSink: it errors on any Apply. It is the safe default for a runtime
// that mediates agent runs but has no action gateway wired, keeping a nil sink (which would panic in the run step)
// out of the executor while guaranteeing no unmediated external write happens.
func NewDenyActionSink() ActionSink { return denyActionSink{} }

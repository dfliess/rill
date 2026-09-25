package act

import (
	"context"
	"errors"
	"fmt"

	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
)

// ErrNoInitiatorClaims is returned by the production SessionFactory when a run reaches the model/tool loop without
// the initiating actor's SecurityClaims. It is the fail-closed half of the Actor contract (executor.go): a run must
// execute with exactly the authority of whoever (or whatever) started it, so a nil-claims run (an automatic trigger
// that declares no identity, or whose user resolved to no attributes, §17.3) must never fall back to a privileged
// session. The run fails instead of escalating.
var ErrNoInitiatorClaims = errors.New("act: run has no initiator claims; refusing to open a privileged session (fail closed)")

// NewSessionFactory returns the production SessionFactory: it opens an ai.Session bound to the run's initiating
// claims, so every tool CheckAccess in the agent loop resolves against the initiator's authority and never the
// worker's. It differs from the test factory in one security-critical way: on nil claims it does NOT fabricate a
// SkipChecks session, it fails closed (ErrNoInitiatorClaims). The claims are passed through exactly as received —
// including a legitimately SkipChecks-bearing local-dev identity, which reflects the runtime's own auth posture
// rather than a fabricated escalation. The session's model resolves from the agent's own connector when it names one,
// and from the instance's configured connector otherwise.
func NewSessionFactory(rt *runtime.Runtime, ac *activity.Client) SessionFactory {
	// One Runner is built up front and reused across sessions: it only registers the tool set and holds no per-run
	// state, so sharing it avoids re-registering tools on every run open.
	runner := ai.NewRunner(rt, ac)
	return func(ctx context.Context, instanceID, sessionID, modelConnector string, claims *runtime.SecurityClaims) (*ai.Session, func(), error) {
		// Fail closed: a run with no initiator identity (an automatic trigger that declares no identity, §17.3) must not run under
		// a privileged session. Refuse rather than escalate to the worker's own authority.
		if claims == nil {
			return nil, nil, ErrNoInitiatorClaims
		}
		// sessionID empty opens a fresh session (first segment); non-empty reopens that session so a resumed segment
		// continues the same conversation. Runner.Session loads the existing AISession and its messages (ordered by
		// index) when SessionID is set, which is the durable substrate the resumed segment reconstructs from.
		s, err := runner.Session(ctx, &ai.SessionOptions{
			InstanceID: instanceID,
			SessionID:  sessionID,
			Claims:     claims,
			UserAgent:  "act-runtime",
		})
		if err != nil {
			return nil, nil, err
		}
		// Route the loop onto the agent's own connector. Without this the session resolves the model from the instance
		// (rill.yaml's ai_connector), which makes an agent's `model.connector` a decorative field: validated by the
		// reconciler, carried on the snapshot, and then ignored at execution. Overriding here rather than in
		// runtime/ai keeps the routing in Act's own code and leaves upstream's session untouched. An agent that names
		// no connector keeps the instance default, so this is additive: nothing changes for agents without a model
		// block. The model itself is not chosen here — it comes from the named connector's own `model:` config, which
		// is how Rill models connections everywhere else.
		if modelConnector != "" {
			s.SetLLM(func(ctx context.Context) (drivers.AIService, func(), error) {
				conn, release, err := rt.AcquireHandle(ctx, instanceID, modelConnector)
				if err != nil {
					return nil, nil, fmt.Errorf("act: agent's model connector %q: %w", modelConnector, err)
				}
				llm, ok := conn.AsAI(instanceID)
				if !ok {
					release()
					return nil, nil, fmt.Errorf("act: agent's model connector %q is not a valid AI service", modelConnector)
				}
				return llm, release, nil
			})
		}
		// The session holds no resources of its own that outlive the run step, so the release is a no-op; the run
		// step's context governs the underlying work.
		return s, func() {}, nil
	}
}

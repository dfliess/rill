package act

import (
	"context"
	"fmt"

	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/ai"
)

// SessionFactory opens an ai.Session for an instance, bound to the given security claims, and returns a release
// function. The executor calls it inside the run_segment step with the initiating actor's claims, so the session (and
// thus every tool CheckAccess) runs with the initiator's authority, not the worker's. sessionID selects the session:
// empty opens a fresh one (the first segment of a run), non-empty reopens that session so a resumed segment continues
// the same conversation from its persisted message tree. Production wires it to the instance's configured model;
// tests wire it to a scripted model.
type SessionFactory func(ctx context.Context, instanceID, sessionID string, claims *runtime.SecurityClaims) (*ai.Session, func(), error)

// ActionSink performs the run's external write. In the spike it records an observable effect; in production it is
// the tool gateway that talks to MCP/action systems (§16). It must be idempotent per run: DBOS steps are
// at-least-once, so ApplyAction may be invoked more than once for the same run after a mid-step crash.
type ActionSink interface {
	Apply(ctx context.Context, req ActionRequest) (ActionResult, error)
}

// SessionRunner is the production-shaped Runner: it resolves snapshots from an AgentDefinitionProvider, runs the
// existing runtime/ai dynamic-agent loop over a Session, and delegates the write to an ActionSink. It holds no
// durable state of its own; all durability comes from the DBOS steps the executor wraps these calls in.
type SessionRunner struct {
	Provider ai.AgentDefinitionProvider
	Sessions SessionFactory
	Actions  ActionSink
}

var _ Runner = (*SessionRunner)(nil)

func (r *SessionRunner) LoadSnapshot(ctx context.Context, instanceID, agentName string) (*ai.AgentSnapshot, error) {
	return r.Provider.GetAgent(ctx, instanceID, agentName)
}

// RunSegment runs one segment of the agent's model/tool loop from the already-captured snapshot, on a session bound to
// the initiator's claims. It runs the snapshot directly rather than re-resolving the agent by name, so the loop is
// bound to the definition captured at run start; the session's claims come from the actor that started the run, so the
// tool intersection reflects the initiator's authority.
//
// A fresh segment (in.SessionID empty) opens a new session; a resumed segment reopens in.SessionID so the loop
// reconstructs the conversation from its persisted message tree and injects in.Resume before the model runs. The
// session is flushed at the end of every segment so the next segment (a separate durable step, possibly after a crash
// or a long approval wait) reopens a fully persisted tree. The opened session ID is returned even on a run/flush error
// (the trace still exists and is worth linking); it is empty only when the session could not be opened at all.
func (r *SessionRunner) RunSegment(ctx context.Context, in RunSegmentInput) (RunSegmentResult, error) {
	s, release, err := r.Sessions(ctx, in.InstanceID, in.SessionID, in.Claims)
	if err != nil {
		return RunSegmentResult{}, err
	}
	defer release()

	// Capture the session ID up front: the segment executes inside this AISession, and the executor links the run to it.
	sessionID := s.ID()

	// A fresh DynamicAgent per segment: its captured-proposal slice must start empty, and the pause/resume state lives
	// in the session tree and in.Resume, not in the agent struct. With multiple proposals per turn this is even more
	// important — a leftover slice from a prior segment would miscount the actions in the next one.
	agent := &ai.DynamicAgent{Snapshot: in.Snapshot}
	// Flush the session trace even when the segment errors: a failed run is exactly the one whose trace matters for
	// auditing. The flush is best-effort, so a flush error only surfaces when the segment itself succeeded and must not
	// mask its own error.
	res, runErr := agent.Run(ctx, s, in.Prompt, in.Resume)
	if flushErr := s.Flush(ctx); flushErr != nil && runErr == nil {
		return RunSegmentResult{SessionID: sessionID}, fmt.Errorf("flush session: %w", flushErr)
	}
	if runErr != nil {
		return RunSegmentResult{SessionID: sessionID}, runErr
	}
	return RunSegmentResult{SessionID: sessionID, Response: res.Response, Proposed: res.Proposed}, nil
}

func (r *SessionRunner) ApplyAction(ctx context.Context, req ActionRequest) (ActionResult, error) {
	return r.Actions.Apply(ctx, req)
}

// ProposeInput is what the executor hands a Proposer to derive the action a finished run proposes.
type ProposeInput struct {
	InstanceID string
	RunID      string
	AgentName  string
	Actor      Actor
	// Snapshot is the immutable agent definition the run executed from; a proposer may consult its tools/instructions.
	Snapshot *ai.AgentSnapshot
	// Response is the agent's final model response, the text from which a proposed tool call is drawn.
	Response string
	// Captured is the structured write the runtime/ai loop already captured during the run (nil for a pure
	// investigation). The production proposer surfaces it directly; a text-parsing proposer may ignore it.
	Captured *ai.ProposedAction
}

// Proposer derives the concrete action a run proposes from its final response and snapshot, or ok=false when the run
// proposed no external action (a pure, read-only investigation ends without a write). It is the boundary between the
// model/tool loop and the Tool Gateway: in production the runtime/ai loop emits a structured tool call and this seam
// surfaces it; for this frente it lets the gateway/ledger/executor path be wired and tested without the loop refactor
// or the outbound MCP client. It must be deterministic given its input, because the executor runs it inside a durable
// step whose result is checkpointed and replayed rather than recomputed.
type Proposer interface {
	Propose(ctx context.Context, in ProposeInput) (proposal ToolProposal, ok bool, err error)
}

package act_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// fakeExecutor is an in-memory ActionExecutor: it records how many times it was asked to Execute/Verify and with what
// resolved arguments, and returns a scripted result. It is the stand-in for the outbound MCP client the parallel
// frente builds, so the gateway can be exercised without a real connector.
type fakeExecutor struct {
	mu sync.Mutex

	result  act.ExecuteResult
	execErr error
	verify  act.VerifyResult

	executeCalls int
	verifyCalls  int
	lastArgs     map[string]any
	lastIdemKey  string
}

var _ act.ActionExecutor = (*fakeExecutor)(nil)

func (f *fakeExecutor) Execute(_ context.Context, req act.ExecuteRequest) (act.ExecuteResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.executeCalls++
	f.lastArgs = req.Args
	f.lastIdemKey = req.IdempotencyKey
	if f.execErr != nil {
		return act.ExecuteResult{}, f.execErr
	}
	return f.result, nil
}

func (f *fakeExecutor) Verify(_ context.Context, _ act.VerifyRequest) (act.VerifyResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.verifyCalls++
	return f.verify, nil
}

func (f *fakeExecutor) executeCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.executeCalls
}

// testSecrets is an act_test-visible SecretResolver (the internal one lives in the act package's own test file).
type testSecrets map[string]string

func (m testSecrets) Resolve(_ context.Context, name string) (string, error) {
	v, ok := m[name]
	if !ok {
		return "", errors.New("unknown secret")
	}
	return v, nil
}

// objectSchema builds a permissive object schema requiring the named string properties, for a tool descriptor.
func objectSchema(required ...string) *jsonschema.Schema {
	props := make(map[string]*jsonschema.Schema, len(required))
	for _, r := range required {
		props[r] = &jsonschema.Schema{Type: "string"}
	}
	return &jsonschema.Schema{Type: "object", Properties: props, Required: required}
}

// createIssueDescriptor is a verifiable, idempotent-native write tool used across the gateway tests.
func createIssueDescriptor(class act.ToolClass) act.ToolDescriptor {
	return act.ToolDescriptor{
		Name:        "jira.create_issue",
		Connector:   "jira_ops",
		Version:     "v1",
		Class:       class,
		InputSchema: objectSchema("summary"),
		Verifiable:  true,
	}
}

// newGateway builds a gateway over a fresh Postgres ledger, the given executor, a one-tool registry and an optional
// secret resolver. It returns the gateway and the ledger so a test can assert the persisted action.
func newGateway(t *testing.T, exec act.ActionExecutor, desc act.ToolDescriptor, secrets act.SecretResolver) (*act.Gateway, *act.PostgresRunStore) {
	t.Helper()
	ledger := newRunStore(t)
	g := &act.Gateway{
		Registry: act.NewMapToolRegistry(desc),
		Ledger:   ledger,
		Executor: exec,
		Secrets:  secrets,
	}
	return g, ledger
}

func createIssueProposal() act.ToolProposal {
	return act.ToolProposal{
		ToolCallID: "call-1",
		Tool:       "jira.create_issue",
		Connector:  "jira_ops",
		Args:       map[string]any{"summary": "coste alto"},
		Summary:    "crear ticket P2",
	}
}

// approveAndExecute drives a proposal through propose → approve → execute directly against the gateway (not the
// workflow), returning the execute report. It is the unit-level equivalent of the run's action phase.
func approveAndExecute(t *testing.T, g *act.Gateway, instanceID, runID string, proposal act.ToolProposal) (act.Authorization, act.ExecuteReport) {
	t.Helper()
	ctx := t.Context()
	auth, err := g.Propose(ctx, act.ProposeActionInput{
		InstanceID: instanceID, RunID: runID, AgentName: "triage", Actor: act.Actor{Subject: "user:alice"}, Proposal: proposal,
	})
	require.NoError(t, err)
	require.Equal(t, act.PolicyApprovalRequired, auth.Decision)

	require.NoError(t, g.RecordApproved(ctx, act.ApprovedInput{
		InstanceID: instanceID, RunID: runID, ToolCallID: proposal.ToolCallID, ArgsHash: auth.ArgsHash, DecidedBy: "admin:bob",
	}))

	report, err := g.Execute(ctx, act.ExecuteActionInput{
		InstanceID: instanceID, RunID: runID, AgentName: "triage", Proposal: proposal,
	})
	require.NoError(t, err)
	return auth, report
}

// TestAuthorizeUnknownToolDenies verifies the pure authorization denies a tool with no descriptor (fail closed): it
// cannot be validated or classified, so it is never presented as executable.
func TestAuthorizeUnknownToolDenies(t *testing.T) {
	auth := act.Authorize(act.AuthorizeInput{
		RunID:    "run-1",
		Proposal: act.ToolProposal{ToolCallID: "call-1", Tool: "jira.delete_project", Args: map[string]any{}},
		Found:    false,
	})
	require.Equal(t, act.PolicyDeny, auth.Decision)
	require.NotEmpty(t, auth.ArgsHash, "even a denied proposal gets an identity for the ledger")
}

// TestAuthorizeInvalidArgsDenies verifies arguments that violate the tool schema are denied before any effect.
func TestAuthorizeInvalidArgsDenies(t *testing.T) {
	auth := act.Authorize(act.AuthorizeInput{
		RunID:      "run-1",
		Proposal:   act.ToolProposal{ToolCallID: "call-1", Tool: "jira.create_issue", Args: map[string]any{"wrong": "field"}},
		Descriptor: createIssueDescriptor(act.ClassIdempotentNative),
		Found:      true,
	})
	require.Equal(t, act.PolicyDeny, auth.Decision)
	require.Contains(t, auth.Reason, "invalid arguments")
}

// TestGatewayProposeRecordsLedger verifies Propose records the action proposed then approval_pending, and returns the
// approval-required decision, so the run knows to suspend and the inbox can see the proposal.
func TestGatewayProposeRecordsLedger(t *testing.T) {
	exec := &fakeExecutor{}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil)
	const inst, run = "inst-propose", "run-1"

	auth, err := g.Propose(t.Context(), act.ProposeActionInput{
		InstanceID: inst, RunID: run, AgentName: "triage", Actor: act.Actor{Subject: "user:alice"}, Proposal: createIssueProposal(),
	})
	require.NoError(t, err)
	require.Equal(t, act.PolicyApprovalRequired, auth.Decision)
	require.NotEmpty(t, auth.IdempotencyKey)

	action, err := ledger.GetAction(t.Context(), inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionApprovalPending, action.Status)
	require.Equal(t, auth.ArgsHash, action.ArgsHash)
	require.Equal(t, act.PolicyApprovalRequired, action.PolicyDecision)
}

// TestGatewayUnknownToolRejectedAndAudited verifies a proposed tool the registry does not know is denied and recorded
// as policy_rejected — rejected, but still auditable (§6.5, §18.3).
func TestGatewayUnknownToolRejectedAndAudited(t *testing.T) {
	exec := &fakeExecutor{}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil)
	const inst, run = "inst-unknown", "run-1"

	auth, err := g.Propose(t.Context(), act.ProposeActionInput{
		InstanceID: inst, RunID: run, AgentName: "triage", Actor: act.Actor{Subject: "user:alice"},
		Proposal: act.ToolProposal{ToolCallID: "call-x", Tool: "jira.delete_project", Args: map[string]any{}},
	})
	require.NoError(t, err)
	require.Equal(t, act.PolicyDeny, auth.Decision)

	action, err := ledger.GetAction(t.Context(), inst, run, "call-x")
	require.NoError(t, err)
	require.Equal(t, act.ActionPolicyRejected, action.Status)
	require.Equal(t, 0, exec.executeCount(), "a denied tool never reaches the executor")
}

// TestGatewayApproveExecuteSucceeds is the happy path: an approved action executes once and is recorded succeeded with
// its external reference, and (being verifiable) is confirmed to verified.
func TestGatewayApproveExecuteSucceeds(t *testing.T) {
	exec := &fakeExecutor{
		result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-123", Message: "created PROJ-123"},
		verify: act.VerifyResult{Confirmed: true, Detail: "exists"},
	}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil)
	const inst, run = "inst-happy", "run-1"
	proposal := createIssueProposal()

	auth, report := approveAndExecute(t, g, inst, run, proposal)
	require.Equal(t, act.OutcomeSucceeded, report.Outcome)
	require.Equal(t, "PROJ-123", report.ExternalReference)
	require.True(t, auth.Verifiable)

	// The idempotency key the executor received is the one the gateway derived.
	require.Equal(t, auth.IdempotencyKey, exec.lastIdemKey)

	// Verify the action, then check the ledger reached verified.
	_, err := g.Verify(t.Context(), act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.NoError(t, err)
	action, err := ledger.GetAction(t.Context(), inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionVerified, action.Status)
	require.Equal(t, "PROJ-123", action.ExternalReference)
	require.Equal(t, "created PROJ-123", action.RedactedMessage, "the tool's text result is persisted on the ledger for audit and durable feedback")
}

// TestGatewayExecuteIsIdempotent proves the core §12 guarantee: a second Execute after success does NOT re-drive the
// external write. The ledger short-circuit returns the recorded outcome and the executor is called exactly once.
func TestGatewayExecuteIsIdempotent(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"}}
	g, _ := newGateway(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil)
	const inst, run = "inst-idem", "run-1"
	proposal := createIssueProposal()

	approveAndExecute(t, g, inst, run, proposal)

	// Re-execute (a replay after the checkpoint was lost): must not re-issue the write.
	report, err := g.Execute(t.Context(), act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeSucceeded, report.Outcome)
	require.Equal(t, "PROJ-1", report.ExternalReference)
	require.Equal(t, 1, exec.executeCount(), "the external write ran exactly once despite two Execute calls")
}

// TestGatewayIndeterminateNotRetried proves an uncertain result on a non-idempotent tool goes to indeterminate and is
// never retried automatically (§12): a second Execute does not call the executor again.
func TestGatewayIndeterminateNotRetried(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeIndeterminate}}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassNonIdempotent), nil)
	const inst, run = "inst-indeterminate", "run-1"
	proposal := createIssueProposal()
	ctx := t.Context()

	auth, err := g.Propose(ctx, act.ProposeActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Actor: act.Actor{Subject: "u"}, Proposal: proposal})
	require.NoError(t, err)
	require.NoError(t, g.RecordApproved(ctx, act.ApprovedInput{InstanceID: inst, RunID: run, ToolCallID: "call-1", ArgsHash: auth.ArgsHash, DecidedBy: "admin"}))

	report, err := g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeIndeterminate, report.Outcome)

	action, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionIndeterminate, action.Status)

	// A retry must not re-execute: an indeterminate non-idempotent write waits for a human.
	report, err = g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeIndeterminate, report.Outcome)
	require.Equal(t, 1, exec.executeCount(), "an indeterminate action is never auto-retried")
}

// TestGatewayApprovalBoundToArgsHash proves the approval binds to the exact arguments (§11.2): approving a different
// args hash than the action holds is refused, forcing a fresh approval for a changed argument.
func TestGatewayApprovalBoundToArgsHash(t *testing.T) {
	exec := &fakeExecutor{}
	g, _ := newGateway(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil)
	const inst, run = "inst-bind", "run-1"
	proposal := createIssueProposal()
	ctx := t.Context()

	auth, err := g.Propose(ctx, act.ProposeActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Actor: act.Actor{Subject: "u"}, Proposal: proposal})
	require.NoError(t, err)

	// Approving a different hash (as if an argument changed after the human saw it) is refused.
	err = g.RecordApproved(ctx, act.ApprovedInput{InstanceID: inst, RunID: run, ToolCallID: "call-1", ArgsHash: act.HashArgs("different"), DecidedBy: "admin"})
	require.ErrorIs(t, err, act.ErrActionTransitionInvalid)

	// The exact hash is accepted.
	require.NoError(t, g.RecordApproved(ctx, act.ApprovedInput{InstanceID: inst, RunID: run, ToolCallID: "call-1", ArgsHash: auth.ArgsHash, DecidedBy: "admin"}))
}

// TestGatewayResolvesSecretsAndRedacts proves the two-sided secret guarantee (§16.2): the executor receives the
// resolved plaintext, while nothing the gateway persists (args or result) contains the secret in the clear.
func TestGatewayResolvesSecretsAndRedacts(t *testing.T) {
	const secret = "SEKRET-xyz789"
	exec := &fakeExecutor{result: act.ExecuteResult{
		Outcome:           act.OutcomeSucceeded,
		ExternalReference: "PROJ-9",
		Message:           "created with token " + secret, // a misbehaving connector echoes the secret back
		Output:            map[string]any{"echo": secret},
	}}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassIdempotentNative), testSecrets{"jira_token": secret})
	const inst, run = "inst-secret", "run-1"
	proposal := act.ToolProposal{
		ToolCallID: "call-1", Tool: "jira.create_issue", Connector: "jira_ops",
		Args:    map[string]any{"summary": "coste alto", "auth": "{{ secret.jira_token }}"},
		Summary: "crear ticket",
	}

	approveAndExecute(t, g, inst, run, proposal)

	// The executor received the resolved secret.
	require.Equal(t, secret, exec.lastArgs["auth"])

	// The persisted action contains no plaintext secret anywhere: args keep the reference, result is scrubbed.
	action, err := ledger.GetAction(t.Context(), inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, "{{ secret.jira_token }}", action.RedactedArgs["auth"], "stored args keep the reference, not the value")
	require.Equal(t, "[REDACTED]", action.RedactedResult["echo"], "a leaked secret in the result is scrubbed")

	blob, err := json.Marshal(action)
	require.NoError(t, err)
	require.False(t, strings.Contains(string(blob), secret), "no plaintext secret survives anywhere on the ledger row")
}

// TestGatewayExecutesFrozenLedgerArgs proves the write uses the arguments frozen at propose, not the caller's live map
// (§17.1): a Proposer that mutates its own map after approval cannot change what executes, and the executor receives
// the ledger's JSON-normalized copy — proving the source is the frozen ledger row, not any caller-held proposal map.
func TestGatewayExecutesFrozenLedgerArgs(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"}}
	// A permissive object schema so an argument carrying a number (to expose JSON normalization) still validates.
	desc := act.ToolDescriptor{
		Name: "jira.create_issue", Connector: "jira_ops", Version: "v1",
		Class: act.ClassIdempotentNative, InputSchema: &jsonschema.Schema{Type: "object"},
	}
	g, ledger := newGateway(t, exec, desc, nil)
	const inst, run = "inst-freeze", "run-1"
	ctx := t.Context()

	// The caller holds this map and keeps mutating it after proposing (a buggy or hostile Proposer). count is an int.
	callerArgs := map[string]any{"summary": "coste alto", "count": 5}
	proposal := act.ToolProposal{ToolCallID: "call-1", Tool: "jira.create_issue", Connector: "jira_ops", Args: callerArgs, Summary: "crear ticket"}

	auth, err := g.Propose(ctx, act.ProposeActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Actor: act.Actor{Subject: "u"}, Proposal: proposal})
	require.NoError(t, err)
	require.NoError(t, g.RecordApproved(ctx, act.ApprovedInput{InstanceID: inst, RunID: run, ToolCallID: "call-1", ArgsHash: auth.ArgsHash, DecidedBy: "admin"}))

	// Mutate the caller's map after approval: it must not change what runs.
	callerArgs["summary"] = "TAMPERED"
	callerArgs["count"] = 999

	// The run re-derives the same proposal deterministically, so Execute is driven with the original arguments (a fresh
	// map that still hashes to the approved value). count is an int here too.
	execProposal := act.ToolProposal{ToolCallID: "call-1", Tool: "jira.create_issue", Connector: "jira_ops", Args: map[string]any{"summary": "coste alto", "count": 5}, Summary: "crear ticket"}
	report, err := g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: execProposal})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeSucceeded, report.Outcome)

	// The executor received the frozen arguments: the original summary (never TAMPERED), and count as the ledger's
	// JSON-normalized float64 — proving the executed args came from the frozen ledger row, not any caller-held map.
	require.Equal(t, "coste alto", exec.lastArgs["summary"])
	require.EqualValues(t, 5, exec.lastArgs["count"])
	require.IsType(t, float64(0), exec.lastArgs["count"], "executed args come from the JSON-round-tripped ledger copy")

	// The persisted args also reflect the frozen snapshot, not the mutation.
	action, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, "coste alto", action.RedactedArgs["summary"])
	require.Equal(t, auth.ArgsHash, action.ArgsHash)
}

// createIssueDescriptorWithAuth is createIssueDescriptor plus a connector credential reference, so the gateway adds
// the bearer token to its redaction set.
func createIssueDescriptorWithAuth(class act.ToolClass, authSecret string) act.ToolDescriptor {
	d := createIssueDescriptor(class)
	d.AuthSecret = authSecret
	return d
}

// TestGatewayScrubsConnectorTokenAndRejectsLeakedReference proves the connector bearer token is redacted everywhere a
// result is persisted, and an external reference that echoes it is dropped rather than stored (§17.1). A hostile server
// returns the credential in the output, the message and as the "handle"; none of it survives on the ledger.
func TestGatewayScrubsConnectorTokenAndRejectsLeakedReference(t *testing.T) {
	const token = "BEARER-tok-99887766"
	exec := &fakeExecutor{result: act.ExecuteResult{
		Outcome:           act.OutcomeSucceeded,
		ExternalReference: token,
		Message:           "created with " + token,
		Output:            map[string]any{"echo": token},
	}}
	g, ledger := newGateway(t, exec, createIssueDescriptorWithAuth(act.ClassIdempotentNative, "jira_bearer"), testSecrets{"jira_bearer": token})
	const inst, run = "inst-token", "run-1"
	proposal := createIssueProposal()

	_, report := approveAndExecute(t, g, inst, run, proposal)
	require.Equal(t, act.OutcomeSucceeded, report.Outcome)
	require.Empty(t, report.ExternalReference, "a reference that echoes the connector token is dropped")

	action, err := ledger.GetAction(t.Context(), inst, run, "call-1")
	require.NoError(t, err)
	require.Empty(t, action.ExternalReference)
	require.Equal(t, "[REDACTED]", action.RedactedResult["echo"])
	blob, err := json.Marshal(action)
	require.NoError(t, err)
	require.False(t, strings.Contains(string(blob), token), "the connector token never survives on the ledger")
}

// TestGatewayRejectsMalformedExternalReference proves a handle that does not match the expected reference shape is
// dropped, while the confirmed write still succeeds: a rejected reference does not undo the effect (§17.1).
func TestGatewayRejectsMalformedExternalReference(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{
		Outcome:           act.OutcomeSucceeded,
		ExternalReference: "line1\nline2 with spaces and \x00 control",
	}}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil)
	const inst, run = "inst-badref", "run-1"

	_, report := approveAndExecute(t, g, inst, run, createIssueProposal())
	require.Equal(t, act.OutcomeSucceeded, report.Outcome)
	require.Empty(t, report.ExternalReference, "a malformed reference is dropped, but the write still succeeds")

	action, err := ledger.GetAction(t.Context(), inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionSucceeded, action.Status)
	require.Empty(t, action.ExternalReference)
}

// TestGatewayInterruptedNonIdempotentGoesIndeterminate proves recovery safety (§12): an action found still executing
// on a non-idempotent tool (a crash mid-write) is marked indeterminate rather than re-issued, because its effect
// cannot be confirmed.
func TestGatewayInterruptedNonIdempotentGoesIndeterminate(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "X"}}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassNonIdempotent), nil)
	const inst, run = "inst-interrupted", "run-1"
	proposal := createIssueProposal()
	ctx := t.Context()

	auth, err := g.Propose(ctx, act.ProposeActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Actor: act.Actor{Subject: "u"}, Proposal: proposal})
	require.NoError(t, err)
	require.NoError(t, g.RecordApproved(ctx, act.ApprovedInput{InstanceID: inst, RunID: run, ToolCallID: "call-1", ArgsHash: auth.ArgsHash, DecidedBy: "admin"}))
	// Simulate a crash mid-write: the action reached executing but never recorded a terminal outcome.
	_, err = ledger.TransitionAction(ctx, act.ActionTransition{InstanceID: inst, RunID: run, ToolCallID: "call-1", Status: act.ActionExecuting})
	require.NoError(t, err)

	report, err := g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeIndeterminate, report.Outcome)
	require.Equal(t, 0, exec.executeCount(), "a non-idempotent interrupted write is not re-issued")

	action, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionIndeterminate, action.Status)
}

// revocableRegistry wraps a fixed descriptor set and can be told, mid-run, to stop returning a tool — standing in for
// an admin removing it from the agent's manifest during an approval wait. It is used to exercise reauthorization.
type revocableRegistry struct {
	inner   act.ToolRegistry
	mu      sync.Mutex
	revoked bool
}

func (r *revocableRegistry) Lookup(ctx context.Context, instanceID, agentName, tool string) (act.ToolDescriptor, bool, error) {
	r.mu.Lock()
	revoked := r.revoked
	r.mu.Unlock()
	if revoked {
		return act.ToolDescriptor{}, false, nil
	}
	return r.inner.Lookup(ctx, instanceID, agentName, tool)
}

func (r *revocableRegistry) revoke() {
	r.mu.Lock()
	r.revoked = true
	r.mu.Unlock()
}

// TestGatewayClaimIsExclusive proves the single-effect guarantee under concurrency (§12): two Execute calls race on
// one approved action, but the atomic approved->executing claim lets exactly one perform the external write. The
// loser either reconciles to the winner's recorded outcome (if the winner finalized first) or gets a retryable
// ErrLeaseHeld (if the winner's lease is still live), so the executor runs exactly once.
func TestGatewayClaimIsExclusive(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"}}
	g, _ := newGateway(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil)
	const inst, run = "inst-claim", "run-1"
	proposal := createIssueProposal()
	ctx := t.Context()

	auth, err := g.Propose(ctx, act.ProposeActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Actor: act.Actor{Subject: "u"}, Proposal: proposal})
	require.NoError(t, err)
	require.NoError(t, g.RecordApproved(ctx, act.ApprovedInput{InstanceID: inst, RunID: run, ToolCallID: "call-1", ArgsHash: auth.ArgsHash, DecidedBy: "admin"}))

	type execResult struct {
		report act.ExecuteReport
		err    error
	}
	var wg sync.WaitGroup
	results := make([]execResult, 2)
	for i := range results {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r, err := g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
			results[i] = execResult{r, err}
		}(i)
	}
	wg.Wait()

	require.Equal(t, 1, exec.executeCount(), "only the claim winner performs the external write")
	// Under at-least-once delivery, one caller wins the claim and reports succeeded; the other either reconciles
	// to the winner's outcome (if it ran after the winner finalized) or gets a retryable ErrLeaseHeld (if it ran
	// while the winner's lease was still live). Both cases preserve the single-effect guarantee.
	var succeeded int
	for _, r := range results {
		if r.err != nil {
			require.ErrorIs(t, r.err, act.ErrLeaseHeld, "the loser's only allowed error is a live-lease retry signal")
		} else {
			require.Equal(t, act.OutcomeSucceeded, r.report.Outcome)
			succeeded++
		}
	}
	require.GreaterOrEqual(t, succeeded, 1, "at least the claim winner reports succeeded")
}

// TestGatewayReauthorizationDenies proves an approved action is re-checked against the registry immediately before the
// write (§17.2): a tool revoked during the approval wait refuses execution, records policy_rejected, and never reaches
// the executor.
func TestGatewayReauthorizationDenies(t *testing.T) {
	exec := &fakeExecutor{}
	ledger := newRunStore(t)
	reg := &revocableRegistry{inner: act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative))}
	g := &act.Gateway{Registry: reg, Ledger: ledger, Executor: exec}
	const inst, run = "inst-reauth", "run-1"
	proposal := createIssueProposal()
	ctx := t.Context()

	auth, err := g.Propose(ctx, act.ProposeActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Actor: act.Actor{Subject: "u"}, Proposal: proposal})
	require.NoError(t, err)
	require.NoError(t, g.RecordApproved(ctx, act.ApprovedInput{InstanceID: inst, RunID: run, ToolCallID: "call-1", ArgsHash: auth.ArgsHash, DecidedBy: "admin"}))

	reg.revoke() // the admin pulls the tool from the manifest while the approval was pending

	report, err := g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeFailed, report.Outcome)
	require.Equal(t, 0, exec.executeCount(), "a revoked tool never reaches the executor")

	action, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionPolicyRejected, action.Status)
}

// TestGatewayRebindMismatchRefuses proves the write is rebound to exactly what was approved (§11.2, §17.1): executing a
// proposal whose arguments changed after approval is refused before any external effect, even though the ledger row is
// approved.
func TestGatewayRebindMismatchRefuses(t *testing.T) {
	exec := &fakeExecutor{}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil)
	const inst, run = "inst-rebind", "run-1"
	proposal := createIssueProposal()
	ctx := t.Context()

	auth, err := g.Propose(ctx, act.ProposeActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Actor: act.Actor{Subject: "u"}, Proposal: proposal})
	require.NoError(t, err)
	require.NoError(t, g.RecordApproved(ctx, act.ApprovedInput{InstanceID: inst, RunID: run, ToolCallID: "call-1", ArgsHash: auth.ArgsHash, DecidedBy: "admin"}))

	// A tampered proposal keeps the same tool call id but changes the arguments the human approved.
	tampered := proposal
	tampered.Args = map[string]any{"summary": "coste MUY alto"}

	report, err := g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: tampered})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeFailed, report.Outcome)
	require.Equal(t, 0, exec.executeCount(), "a proposal that does not match the approved args never executes")

	action, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionPolicyRejected, action.Status)
}

// TestGatewayVerifyRedactsReturnedDetail proves the verification detail returned to the workflow is redacted, not just
// the ledger copy (§16.2, C1): DBOS checkpoints the returned VerifyResult, so a connector token echoed in the detail
// must be scrubbed both on the returned value and on the persisted verification.
func TestGatewayVerifyRedactsReturnedDetail(t *testing.T) {
	const token = "BEARER-verify-11223344"
	exec := &fakeExecutor{
		result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"},
		verify: act.VerifyResult{Confirmed: true, Detail: "confirmed with " + token},
	}
	g, ledger := newGateway(t, exec, createIssueDescriptorWithAuth(act.ClassIdempotentNative, "jira_bearer"), testSecrets{"jira_bearer": token})
	const inst, run = "inst-verify-redact", "run-1"
	proposal := createIssueProposal()
	ctx := t.Context()

	approveAndExecute(t, g, inst, run, proposal)

	res, err := g.Verify(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.NoError(t, err)
	require.True(t, res.Confirmed)
	require.NotContains(t, res.Detail, token, "the detail returned to the workflow (which DBOS checkpoints) is redacted")
	require.Contains(t, res.Detail, "[REDACTED]")

	action, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionVerified, action.Status)
	require.NotContains(t, action.RedactedVerification["detail"], token)
}

// TestGatewayWriteFailsClosedOnUnresolvableConnectorSecret proves the write path fails closed when a declared connector
// credential cannot be resolved for redaction (C1): the gateway cannot guarantee the result would be scrubbed, so it
// refuses the write before any external effect rather than risk persisting a token in the clear.
func TestGatewayWriteFailsClosedOnUnresolvableConnectorSecret(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"}}
	// The descriptor declares a connector credential the resolver does not know: redaction cannot be guaranteed.
	g, ledger := newGateway(t, exec, createIssueDescriptorWithAuth(act.ClassIdempotentNative, "missing_bearer"), testSecrets{"other": "x"})
	const inst, run = "inst-nosecret", "run-1"
	proposal := createIssueProposal()

	_, report := approveAndExecute(t, g, inst, run, proposal)
	require.Equal(t, act.OutcomeFailed, report.Outcome)
	require.Equal(t, 0, exec.executeCount(), "a write whose result could not be redacted never runs")

	action, err := ledger.GetAction(t.Context(), inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionFailed, action.Status)
}

// TestGatewayReauthorizePreservesAutoApprove proves reauthorization re-evaluates with the SAME AutoApprove the proposal
// was authorized under (§17.2, C2): an auto-approved action still permitted under the current policy executes, rather
// than being read as newly approval-required and rejected during the pre-write recheck.
func TestGatewayReauthorizePreservesAutoApprove(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"}}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil)
	const inst, run = "inst-autoapprove", "run-1"
	proposal := createIssueProposal()
	autoApprove := map[string]bool{"jira.create_issue": true}
	ctx := t.Context()

	auth, err := g.Propose(ctx, act.ProposeActionInput{
		InstanceID: inst, RunID: run, AgentName: "triage", Actor: act.Actor{Subject: "u"},
		Proposal: proposal, AutoApprove: autoApprove,
	})
	require.NoError(t, err)
	require.Equal(t, act.PolicyAllow, auth.Decision, "the tool is auto-approved by agent policy")

	report, err := g.Execute(ctx, act.ExecuteActionInput{
		InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal, AutoApprove: autoApprove,
	})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeSucceeded, report.Outcome, "a still-permitted auto-approved action is not rejected at reauthorization")
	require.Equal(t, 1, exec.executeCount())

	action, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionSucceeded, action.Status)
}

// TestGatewayAutoApproveExecutesUnclassified proves a connector's auto-approve posture actually takes effect on the
// production path, where every generic MCP tool is unclassified (the runtime registry cannot establish an idempotency
// class). Auto-approve must skip the human gate even for an unclassified tool: the class governs auto-RETRY of an
// uncertain result (disabled in v1), not whether a human gates the single execution. Without this, approval.auto /
// auto_approve would parse and hash but silently no-op, since the unclassified cap would revert the allow to approval.
// The sibling TestGatewayReauthorizePreservesAutoApprove uses a KNOWN class, so it never exercised this path.
func TestGatewayAutoApproveExecutesUnclassified(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"}}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassUnknown), nil)
	const inst, run = "inst-autoapprove-unknown", "run-1"
	proposal := createIssueProposal()
	autoApprove := map[string]bool{"jira.create_issue": true}
	ctx := t.Context()

	auth, err := g.Propose(ctx, act.ProposeActionInput{
		InstanceID: inst, RunID: run, AgentName: "triage", Actor: act.Actor{Subject: "u"},
		Proposal: proposal, AutoApprove: autoApprove,
	})
	require.NoError(t, err)
	require.Equal(t, act.PolicyAllow, auth.Decision, "an explicitly auto-approved tool is allowed even when unclassified")

	// The ledger recorded approved directly at propose, so Execute proceeds without a human decision.
	action, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionApproved, action.Status)

	report, err := g.Execute(ctx, act.ExecuteActionInput{
		InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal, AutoApprove: autoApprove,
	})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeSucceeded, report.Outcome)
	require.Equal(t, 1, exec.executeCount(), "the auto-approved write executes exactly once, no human gate")

	action, err = ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionSucceeded, action.Status)
}

// TestGatewayCapsOversizedExecutorOutput proves the gateway caps a result's output and message at its own boundary
// before persisting (§16.2, D2): an executor that returns an oversized structured output or message cannot flood the
// ledger, regardless of whether the executor imposed its own cap.
func TestGatewayCapsOversizedExecutorOutput(t *testing.T) {
	big := strings.Repeat("A", 4096)
	exec := &fakeExecutor{result: act.ExecuteResult{
		Outcome:           act.OutcomeSucceeded,
		ExternalReference: "PROJ-1",
		Message:           big,
		Output:            map[string]any{"blob": big},
	}}
	ledger := newRunStore(t)
	g := &act.Gateway{
		Registry:       act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative)),
		Ledger:         ledger,
		Executor:       exec,
		MaxOutputBytes: 256, // a tight cap so the 4 KiB output and message are over it
	}
	const inst, run = "inst-cap", "run-1"
	proposal := createIssueProposal()

	_, report := approveAndExecute(t, g, inst, run, proposal)
	require.Equal(t, act.OutcomeSucceeded, report.Outcome)
	require.LessOrEqual(t, len(report.Message), 256, "the message returned to the workflow is capped")

	action, err := ledger.GetAction(t.Context(), inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, "output exceeded size limit", action.RedactedResult["_dropped"], "an oversized output is dropped, not persisted whole")
	require.LessOrEqual(t, len(action.RedactedMessage), 256, "the persisted text message is capped, not stored whole")
	blob, err := json.Marshal(action)
	require.NoError(t, err)
	require.Less(t, len(blob), 1536, "no oversized field survives on the ledger row (the capped message adds at most MaxOutputBytes)")
}

// TestGatewayKillSwitchDenies proves an engaged kill switch denies the proposal at the gateway and records it
// policy_rejected, with no external effect (§17.2).
func TestGatewayKillSwitchDenies(t *testing.T) {
	exec := &fakeExecutor{}
	ledger := newRunStore(t)
	kill := act.NewMapKillSwitch()
	kill.DisableInstance("inst-kill", "incident response")
	g := &act.Gateway{
		Registry:   act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative)),
		Ledger:     ledger,
		Executor:   exec,
		KillSwitch: kill,
	}

	auth, err := g.Propose(t.Context(), act.ProposeActionInput{
		InstanceID: "inst-kill", RunID: "run-1", AgentName: "triage", Actor: act.Actor{Subject: "u"}, Proposal: createIssueProposal(),
	})
	require.NoError(t, err)
	require.Equal(t, act.PolicyDeny, auth.Decision)
	require.Contains(t, auth.Reason, "incident response")

	action, err := ledger.GetAction(t.Context(), "inst-kill", "run-1", "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionPolicyRejected, action.Status)
	require.Equal(t, 0, exec.executeCount())
}

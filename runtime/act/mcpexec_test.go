package act_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/act/mcpconn"
	"github.com/stretchr/testify/require"
)

// fakeMCPClient is a scripted mcpconn.MCPClient: CallTool returns whatever result/error the test sets, so the
// adapter's outcome classification and reference extraction can be exercised without a live MCP server (the real
// client is covered in package mcpconn).
type fakeMCPClient struct {
	result mcpconn.RemoteToolResult
	err    error

	gotName string
	gotHash string
	gotArgs json.RawMessage
}

func (f *fakeMCPClient) ListTools(context.Context) ([]mcpconn.RemoteTool, error) { return nil, nil }

func (f *fakeMCPClient) CallTool(_ context.Context, name, expectedSchemaHash string, args json.RawMessage) (mcpconn.RemoteToolResult, error) {
	f.gotName, f.gotHash, f.gotArgs = name, expectedSchemaHash, args
	return f.result, f.err
}

func TestMCPExecutor_SuccessExtractsReference(t *testing.T) {
	fake := &fakeMCPClient{result: mcpconn.RemoteToolResult{
		Text:              "created PROJ-123",
		StructuredContent: json.RawMessage(`{"key":"PROJ-123","self":"https://jira/PROJ-123"}`),
	}}
	e := act.NewMCPActionExecutor(fake)

	res, err := e.Execute(context.Background(), act.ExecuteRequest{
		Tool:        "mcp.jira.create_issue",
		ToolVersion: "sha256-abc",
		Args:        map[string]any{"summary": "x"},
	})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeSucceeded, res.Outcome)
	require.Equal(t, "PROJ-123", res.ExternalReference)
	require.Equal(t, "mcp.jira.create_issue", fake.gotName)
	require.Equal(t, "sha256-abc", fake.gotHash, "the action's pinned schema hash is passed through to the client")
	require.JSONEq(t, `{"summary":"x"}`, string(fake.gotArgs))
}

func TestMCPExecutor_WriteToolErrorIsIndeterminate(t *testing.T) {
	// A write tool that reports its own failure may have written before failing: the effect is unknown.
	fake := &fakeMCPClient{result: mcpconn.RemoteToolResult{Text: "timed out after create", IsError: true}}
	e := act.NewMCPActionExecutor(fake)

	res, err := e.Execute(context.Background(), act.ExecuteRequest{Tool: "mcp.jira.create_issue"})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeIndeterminate, res.Outcome, "a write tool's reported failure is not a definitive non-write")
}

func TestMCPExecutor_ReadOnlyToolErrorIsFailed(t *testing.T) {
	// A trusted read-only tool cannot have written, so its reported failure is a safe definitive non-write.
	fake := &fakeMCPClient{result: mcpconn.RemoteToolResult{Text: "permission denied", IsError: true}}
	e := act.NewMCPActionExecutor(fake)

	res, err := e.Execute(context.Background(), act.ExecuteRequest{Tool: "mcp.jira.get_issue", ReadOnly: true})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeFailed, res.Outcome)
}

func TestMCPExecutor_PreSendErrorIsFailed(t *testing.T) {
	// A schema-drift rejection happens before the request is sent: the write definitively did not happen.
	fake := &fakeMCPClient{err: mcpconn.ErrSchemaChanged}
	e := act.NewMCPActionExecutor(fake)

	res, err := e.Execute(context.Background(), act.ExecuteRequest{Tool: "mcp.jira.create_issue"})
	require.Error(t, err)
	require.Equal(t, act.OutcomeFailed, res.Outcome)
}

func TestMCPExecutor_StructuredContentCapped(t *testing.T) {
	// Structured content larger than the gateway's cap is dropped before it is parsed or persisted.
	big := `{"data":"` + strings.Repeat("x", 200) + `"}`
	fake := &fakeMCPClient{result: mcpconn.RemoteToolResult{Text: "ok", StructuredContent: json.RawMessage(big)}}
	e := act.NewMCPActionExecutor(fake)

	res, err := e.Execute(context.Background(), act.ExecuteRequest{Tool: "mcp.jira.create_issue", MaxOutputBytes: 64})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeSucceeded, res.Outcome)
	require.Contains(t, res.Output, "_dropped", "oversized structured content is dropped, not persisted")
}

func TestMCPExecutor_AmbiguousErrorIsBare(t *testing.T) {
	// A transport error after dispatch is ambiguous: the adapter returns a bare error and leaves classification to
	// the gateway (which turns a non-idempotent write into indeterminate rather than retrying).
	fake := &fakeMCPClient{err: errors.New("connection reset")}
	e := act.NewMCPActionExecutor(fake)

	res, err := e.Execute(context.Background(), act.ExecuteRequest{Tool: "mcp.jira.create_issue"})
	require.Error(t, err)
	require.Empty(t, res.Outcome, "an ambiguous error must not assert an outcome; the gateway classifies by tool class")
}

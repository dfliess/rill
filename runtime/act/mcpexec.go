package act

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rilldata/rill/runtime/act/mcpconn"
)

// mcpExecutor is the ActionExecutor that performs a gateway-authorized write by calling an outbound MCP tool. It is
// the join between the two Fase 2 halves: the gateway resolves policy, secrets, idempotency and redaction, then hands
// a fully-prepared ExecuteRequest here, and this turns it into one MCP CallTool on the pinned tool. It never widens
// authority: it can only call a tool the client already discovered and pinned.
type mcpExecutor struct {
	client mcpconn.MCPClient
}

// NewMCPActionExecutor returns an ActionExecutor backed by an outbound MCP client.
func NewMCPActionExecutor(client mcpconn.MCPClient) ActionExecutor {
	return &mcpExecutor{client: client}
}

var _ ActionExecutor = (*mcpExecutor)(nil)

// Execute performs exactly one MCP tool call and classifies its outcome (§12). The classification is what keeps an
// uncertain write from being retried, and it errs toward indeterminate whenever a write may have landed:
//
//   - A pre-dispatch guard error (tool not discovered, schema drift, tool gone) is checked before sess.CallTool ever
//     runs, so the call provably never went out: OutcomeFailed, a definitive non-write.
//   - Any other error — a blocked host, a refused redirect (which fires only after a response), or a transport failure
//     — may have reached the server, so whether the write landed is unknown: return the bare error and let the gateway
//     classify it indeterminate.
//   - A tool-level error (IsError) on a write tool may have written before failing, so it is treated as indeterminate;
//     only a read-only tool's IsError is a safe OutcomeFailed.
//   - A clean result is OutcomeSucceeded, with the external reference extracted best-effort from the structured output.
func (e *mcpExecutor) Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error) {
	if req.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, req.Timeout)
		defer cancel()
	}

	argsJSON, err := json.Marshal(req.Args)
	if err != nil {
		// Marshaling our own already-validated args cannot really fail; if it does the call never went out.
		return ExecuteResult{Outcome: OutcomeFailed}, fmt.Errorf("act: marshal action args: %w", err)
	}

	res, err := e.client.CallTool(ctx, req.Tool, req.ToolVersion, argsJSON)
	if err != nil {
		// Only the guards checked before dispatch prove the request never left the client: report failed so a human
		// knows nothing happened. A blocked host, a refused redirect (ErrRedirectBlocked, post-response) or any other
		// transport error is ambiguous: surface the bare error so the gateway classifies it indeterminate.
		if errors.Is(err, mcpconn.ErrToolNotDiscovered) ||
			errors.Is(err, mcpconn.ErrSchemaChanged) ||
			errors.Is(err, mcpconn.ErrToolGone) {
			return ExecuteResult{Outcome: OutcomeFailed}, fmt.Errorf("act: mcp call rejected before dispatch: %w", err)
		}
		return ExecuteResult{}, fmt.Errorf("act: mcp call failed: %w", err)
	}

	if res.IsError {
		// The tool ran and reported failure. For a write tool it may have written before failing, so the effect is
		// unknown: indeterminate. Only a trusted read-only tool's failure is a safe definitive non-write.
		if req.ReadOnly {
			return ExecuteResult{Outcome: OutcomeFailed, Message: truncate(res.Text, req.MaxOutputBytes)}, nil
		}
		return ExecuteResult{Outcome: OutcomeIndeterminate, Message: truncate(res.Text, req.MaxOutputBytes)}, nil
	}

	out := ExecuteResult{
		Outcome: OutcomeSucceeded,
		Message: truncate(res.Text, req.MaxOutputBytes),
	}
	// Cap the structured content by its serialized size BEFORE parsing it. MaxOutputBytes bounds the text result, but a
	// hostile or runaway server could still return a large structured object; over the cap, drop it with a marker
	// rather than unmarshal and persist an unbounded map (§16.2). The reference is only extracted from content within
	// the cap, so an oversized blob is never scanned either.
	if len(res.StructuredContent) > 0 {
		if req.MaxOutputBytes > 0 && len(res.StructuredContent) > req.MaxOutputBytes {
			out.Output = map[string]any{"_dropped": "structured content exceeded size limit"}
		} else {
			out.ExternalReference = externalReference(res.StructuredContent)
			var m map[string]any
			if err := json.Unmarshal(res.StructuredContent, &m); err == nil {
				out.Output = m
			}
		}
	}
	return out, nil
}

// Verify is best-effort and, for a generic MCP tool, unsupported: MCP has no standard "confirm this write" call, so
// confirmation would require per-connector read logic. Confirmed=false leaves the action unverified (not failed): the
// write already succeeded (§16.1). Certified connectors can supply a verifying executor later.
func (e *mcpExecutor) Verify(_ context.Context, _ VerifyRequest) (VerifyResult, error) {
	return VerifyResult{Confirmed: false, Detail: "verification not supported for generic MCP tools"}, nil
}

// externalReference pulls a durable handle out of a tool's structured output, trying the field names MCP action tools
// commonly use for it. It is best-effort: an empty string just means no handle was surfaced, not a failure.
func externalReference(structured json.RawMessage) string {
	if len(structured) == 0 {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(structured, &m); err != nil {
		return ""
	}
	for _, k := range []string{"key", "id", "reference", "external_reference", "self", "url"} {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// truncate bounds a string to maxBytes (maxBytes <= 0 means no bound), then drops a trailing partial rune so the
// result is always valid UTF-8.
func truncate(s string, maxBytes int) string {
	if maxBytes <= 0 || len(s) <= maxBytes {
		return s
	}
	return strings.ToValidUTF8(s[:maxBytes], "")
}

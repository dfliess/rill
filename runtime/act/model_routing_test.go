package act_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
)

type modelRoutingRequest struct {
	Model          string
	Authorization  string
	SystemMessages []string
}

type modelRoutingServer struct {
	server *httptest.Server
	mu     sync.Mutex
	reqs   []modelRoutingRequest
}

func newModelRoutingServer(t *testing.T, response string) *modelRoutingServer {
	t.Helper()
	r := &modelRoutingServer{}
	r.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/v1/chat/completions" {
			http.Error(w, "unexpected path "+req.URL.Path, http.StatusNotFound)
			return
		}
		var body struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON request: "+err.Error(), http.StatusBadRequest)
			return
		}
		var systemMessages []string
		for _, message := range body.Messages {
			if message.Role == "system" {
				systemMessages = append(systemMessages, message.Content)
			}
		}
		r.mu.Lock()
		r.reqs = append(r.reqs, modelRoutingRequest{Model: body.Model, Authorization: req.Header.Get("Authorization"), SystemMessages: systemMessages})
		r.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id":      "completion-test",
			"object":  "chat.completion",
			"created": 1,
			"model":   body.Model,
			"choices": []any{map[string]any{
				"index":         0,
				"finish_reason": "stop",
				"message": map[string]any{
					"role":    "assistant",
					"content": response,
				},
			}},
			"usage": map[string]any{
				"prompt_tokens":     1,
				"completion_tokens": 1,
				"total_tokens":      2,
			},
		}); err != nil {
			t.Errorf("encode fake completion: %v", err)
		}
	}))
	t.Cleanup(r.server.Close)
	return r
}

func (s *modelRoutingServer) url() string { return s.server.URL + "/v1" }

func (s *modelRoutingServer) requests() []modelRoutingRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]modelRoutingRequest(nil), s.reqs...)
}

func modelRoutingConnectorYAML(baseURL, model string) string {
	return fmt.Sprintf(`
type: connector
driver: openai
api_key: "{{ .env.agent_api_key }}"
base_url: %q
model: %q
`, baseURL, model)
}

func modelRoutingAgentYAML(connector, model string) string {
	return fmt.Sprintf(`
type: agent
instructions: "Answer the request."
model:
  connector: %q
  name: %q
tools: []
`, connector, model)
}

// TestAgentModelRoutingSurvivesSessionReopen proves the production path uses the agent-selected connector and model,
// not the project's default, and keeps that non-secret configuration fixed when a later segment reopens the session.
// Between segments both the agent and its old connector are edited, while the API key rotates: the resumed request
// must still reach the original endpoint/model but authenticate with the new secret.
func TestAgentModelRoutingSurvivesSessionReopen(t *testing.T) {
	original := newModelRoutingServer(t, "original provider")
	forbidden := newModelRoutingServer(t, "wrong provider")

	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Variables: map[string]string{"agent_api_key": "initial-secret"},
		Files: map[string]string{
			"rill.yaml": `ai_connector: default_ai
ai_instructions: "ORIGINAL PROJECT INSTRUCTION"`,
			"connectors/agent_ai.yaml": modelRoutingConnectorYAML(
				original.url(), "connector-original-model",
			),
			"connectors/default_ai.yaml": modelRoutingConnectorYAML(
				forbidden.url(), "default-model",
			),
			"connectors/edited_ai.yaml": modelRoutingConnectorYAML(
				forbidden.url(), "edited-connector-model",
			),
			"triage.yaml": modelRoutingAgentYAML("agent_ai", "agent-original-model"),
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)

	provider := act.NewCatalogAgentProvider(rt)
	snapshot, err := provider.GetAgent(t.Context(), instanceID, "triage")
	require.NoError(t, err)
	require.Equal(t, "agent_ai", snapshot.ModelConnector)
	require.Equal(t, "openai", snapshot.ModelDriver)
	require.Equal(t, "agent-original-model", snapshot.ModelName)
	require.Equal(t, "ORIGINAL PROJECT INSTRUCTION", snapshot.ProjectInstructions)
	require.Equal(t, "agent-original-model", snapshot.ModelProperties["model"])
	require.Equal(t, original.url(), snapshot.ModelProperties["base_url"])
	require.NotContains(t, snapshot.ModelProperties, "api_key")
	checkpoint, err := json.Marshal(snapshot)
	require.NoError(t, err)
	require.NotContains(t, string(checkpoint), "initial-secret")

	runner := &act.SessionRunner{
		Provider: provider,
		Sessions: act.NewSessionFactory(rt, activity.NewNoopClient()),
		Actions:  act.NewDenyActionSink(),
	}
	claims := &runtime.SecurityClaims{UserID: "routing-user", SkipChecks: true}
	first, err := runner.RunSegment(t.Context(), act.RunSegmentInput{
		InstanceID: instanceID,
		Claims:     claims,
		Snapshot:   snapshot,
		Prompt:     "first segment",
	})
	require.NoError(t, err)
	require.NotEmpty(t, first.SessionID)

	// Change both possible live sources of non-secret behaviour: the agent now selects another connector/model, and
	// the old connector itself now points at another endpoint/model. Also rotate only the secret on the instance.
	testruntime.PutFiles(t, rt, instanceID, map[string]string{
		"rill.yaml": `ai_connector: default_ai
ai_instructions: "EDITED PROJECT INSTRUCTION"`,
		"connectors/agent_ai.yaml": modelRoutingConnectorYAML(forbidden.url(), "connector-edited-model"),
		"triage.yaml":              modelRoutingAgentYAML("edited_ai", "agent-edited-model"),
	})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)
	inst, err := rt.Instance(t.Context(), instanceID)
	require.NoError(t, err)
	inst.Variables["agent_api_key"] = "rotated-secret"
	require.NoError(t, rt.EditInstance(t.Context(), inst, false))

	live, err := provider.GetAgent(t.Context(), instanceID, "triage")
	require.NoError(t, err)
	require.Equal(t, "edited_ai", live.ModelConnector)
	require.Equal(t, "agent-edited-model", live.ModelProperties["model"])
	require.Equal(t, "EDITED PROJECT INSTRUCTION", live.ProjectInstructions)

	_, err = runner.RunSegment(t.Context(), act.RunSegmentInput{
		InstanceID: instanceID,
		Claims:     claims,
		Snapshot:   snapshot,
		Prompt:     "first segment",
		SessionID:  first.SessionID,
		Resume: []*ai.InjectedResult{{
			Tool:    "mcp.example.write",
			Message: "approved",
		}},
	})
	require.NoError(t, err)

	require.Empty(t, forbidden.requests(), "neither the project default nor edited connector may serve this run")
	reqs := original.requests()
	require.Len(t, reqs, 2)
	require.Equal(t, "agent-original-model", reqs[0].Model)
	require.Equal(t, "Bearer initial-secret", reqs[0].Authorization)
	require.Contains(t, strings.Join(reqs[0].SystemMessages, "\n"), "ORIGINAL PROJECT INSTRUCTION")
	require.Equal(t, "agent-original-model", reqs[1].Model)
	require.Equal(t, "Bearer rotated-secret", reqs[1].Authorization)
	require.Contains(t, strings.Join(reqs[1].SystemMessages, "\n"), "ORIGINAL PROJECT INSTRUCTION")
	for _, system := range reqs[1].SystemMessages {
		require.NotContains(t, system, "EDITED PROJECT INSTRUCTION")
	}
}

// TestAgentModelRoutingRefusesDriverDrift proves live secret resolution cannot cross provider contracts. If an
// operator replaces the selected OpenAI connector with a Claude connector while a run waits, its new credential must
// never be sent to the endpoint frozen in the old snapshot; the resume fails closed instead.
func TestAgentModelRoutingRefusesDriverDrift(t *testing.T) {
	original := newModelRoutingServer(t, "unused")
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Variables: map[string]string{"agent_api_key": "openai-secret"},
		Files: map[string]string{
			"rill.yaml":                "",
			"connectors/agent_ai.yaml": modelRoutingConnectorYAML(original.url(), "original-model"),
			"triage.yaml":              modelRoutingAgentYAML("agent_ai", "agent-model"),
		},
	})
	provider := act.NewCatalogAgentProvider(rt)
	snapshot, err := provider.GetAgent(t.Context(), instanceID, "triage")
	require.NoError(t, err)

	testruntime.PutFiles(t, rt, instanceID, map[string]string{
		"connectors/agent_ai.yaml": `
type: connector
driver: claude
api_key: "{{ .env.agent_api_key }}"
model: claude-sonnet-4-20250514
`,
	})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)

	_, release, err := rt.AIFromConnectorSnapshot(t.Context(), instanceID, snapshot.ModelConnector, snapshot.ModelDriver, snapshot.ModelProperties)
	require.ErrorContains(t, err, `changed driver from "openai" to "claude"`)
	require.Nil(t, release)
	require.Empty(t, original.requests(), "driver drift must fail before any credential reaches the frozen endpoint")
}

func TestAgentModelRoutingRejectsManagedModelOverride(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"connectors/managed_ai.yaml": `
type: connector
driver: admin
admin_url: https://admin.example.test
access_token: managed-secret
`,
			"triage.yaml": modelRoutingAgentYAML("managed_ai", "provider-specific-model"),
		},
	})

	_, err := act.NewCatalogAgentProvider(rt).GetAgent(t.Context(), instanceID, "triage")
	require.ErrorContains(t, err, "model.name is not supported by the managed AI connector")
}

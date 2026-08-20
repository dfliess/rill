package ai_test

import (
	"testing"

	"github.com/rilldata/rill/runtime/ai"
	"github.com/stretchr/testify/require"
)

// TestStaticAgentProviderDeepCopiesModelProperties protects the snapshot's immutable connector behaviour. Provider
// request bodies commonly contain nested maps and slices, so copying only the top-level map would let either the
// caller that registered the agent or one run mutate every later run's model configuration.
func TestStaticAgentProviderDeepCopiesModelProperties(t *testing.T) {
	original := &ai.AgentSnapshot{
		Name: "triage",
		ModelProperties: map[string]any{
			"model": "original-model",
			"extra_body": map[string]any{
				"thinking": map[string]any{"type": "disabled"},
				"stops":    []any{"one", "two"},
				"typed":    map[string][]bool{"flags": {true, false}},
			},
		},
	}
	provider := ai.NewStaticAgentProvider(original)

	// Mutating the input after registration must not alter the provider's stored definition.
	original.ModelProperties["model"] = "mutated-input"
	original.ModelProperties["extra_body"].(map[string]any)["thinking"].(map[string]any)["type"] = "enabled"
	original.ModelProperties["extra_body"].(map[string]any)["stops"].([]any)[0] = "mutated"
	original.ModelProperties["extra_body"].(map[string]any)["typed"].(map[string][]bool)["flags"][0] = false

	first, err := provider.GetAgent(t.Context(), "instance", "triage")
	require.NoError(t, err)
	require.Equal(t, "original-model", first.ModelProperties["model"])
	require.Equal(t, "disabled", first.ModelProperties["extra_body"].(map[string]any)["thinking"].(map[string]any)["type"])
	require.Equal(t, "one", first.ModelProperties["extra_body"].(map[string]any)["stops"].([]any)[0])
	require.True(t, first.ModelProperties["extra_body"].(map[string]any)["typed"].(map[string][]bool)["flags"][0])

	// Mutating one returned snapshot must not alter the provider or another run's snapshot.
	first.ModelProperties["model"] = "mutated-run"
	first.ModelProperties["extra_body"].(map[string]any)["thinking"].(map[string]any)["type"] = "enabled"
	first.ModelProperties["extra_body"].(map[string]any)["stops"].([]any)[1] = "mutated"
	first.ModelProperties["extra_body"].(map[string]any)["typed"].(map[string][]bool)["flags"][1] = true

	second, err := provider.GetAgent(t.Context(), "instance", "triage")
	require.NoError(t, err)
	require.Equal(t, "original-model", second.ModelProperties["model"])
	require.Equal(t, "disabled", second.ModelProperties["extra_body"].(map[string]any)["thinking"].(map[string]any)["type"])
	require.Equal(t, "two", second.ModelProperties["extra_body"].(map[string]any)["stops"].([]any)[1])
	require.False(t, second.ModelProperties["extra_body"].(map[string]any)["typed"].(map[string][]bool)["flags"][1])
}

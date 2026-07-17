package act

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// mapSecrets is an in-memory SecretResolver for tests: it resolves names it knows and errors on the rest, so a test
// can exercise both the resolution and the unknown-secret fail-closed path.
type mapSecrets map[string]string

func (m mapSecrets) Resolve(_ context.Context, name string) (string, error) {
	v, ok := m[name]
	if !ok {
		return "", errors.New("unknown secret")
	}
	return v, nil
}

// TestSecretRef verifies only a standalone `{{ secret.NAME }}` value is treated as a reference; partial matches and
// non-strings are literal data, so a secret can never be smuggled in as a substring.
func TestSecretRef(t *testing.T) {
	name, ok := SecretRef("{{ secret.jira_token }}")
	require.True(t, ok)
	require.Equal(t, "jira_token", name)

	_, ok = SecretRef("Authorization: {{ secret.jira_token }}") // partial: not a reference
	require.False(t, ok)

	_, ok = SecretRef(42) // non-string
	require.False(t, ok)
}

// TestResolveSecretsReplacesButPreservesOriginal verifies resolution builds a call-only copy with plaintext values,
// leaving the reference form in the original map (which is what gets persisted), and collects the resolved values.
func TestResolveSecretsReplacesButPreservesOriginal(t *testing.T) {
	secrets := mapSecrets{"jira_token": "SEKRET-abc123"}
	args := map[string]any{
		"auth":   "{{ secret.jira_token }}",
		"nested": map[string]any{"header": "{{ secret.jira_token }}"},
		"list":   []any{"literal", "{{ secret.jira_token }}"},
		"plain":  "not a secret",
	}

	resolved, values, err := resolveSecrets(context.Background(), secrets, args)
	require.NoError(t, err)

	// The returned copy carries plaintext for the executor call.
	require.Equal(t, "SEKRET-abc123", resolved["auth"])
	require.Equal(t, "SEKRET-abc123", resolved["nested"].(map[string]any)["header"])
	require.Equal(t, "SEKRET-abc123", resolved["list"].([]any)[1])
	require.Equal(t, "not a secret", resolved["plain"])
	require.Equal(t, []string{"SEKRET-abc123"}, values, "distinct resolved values, once")

	// The original is untouched: the reference form is what the ledger persists (no plaintext).
	require.Equal(t, "{{ secret.jira_token }}", args["auth"])
	require.Equal(t, "{{ secret.jira_token }}", args["nested"].(map[string]any)["header"])
}

// TestResolveSecretsFailsClosed verifies an unknown secret, or a reference with no resolver, is a hard error: the
// action cannot execute with a hole where a credential should be.
func TestResolveSecretsFailsClosed(t *testing.T) {
	_, _, err := resolveSecrets(context.Background(), mapSecrets{}, map[string]any{"auth": "{{ secret.missing }}"})
	require.Error(t, err)

	_, _, err = resolveSecrets(context.Background(), nil, map[string]any{"auth": "{{ secret.jira_token }}"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no resolver")
}

// TestRedactionScrubsSecretValues verifies any secret value echoed back in a string is replaced by the marker, in
// nested maps and slices alike — the last line of defense against a connector leaking a credential into the ledger.
func TestRedactionScrubsSecretValues(t *testing.T) {
	values := []string{"SEKRET-abc123"}

	require.Equal(t, "token is [REDACTED] ok", redactString("token is SEKRET-abc123 ok", values))

	out := RedactMap(map[string]any{
		"echo":          "leaked SEKRET-abc123",
		"nested":        map[string]any{"deep": "SEKRET-abc123"},
		"list":          []any{"SEKRET-abc123", "clean"},
		"num":           1.0,
		"SEKRET-abc123": "value under a leaked key",
	}, values)
	require.Equal(t, "leaked [REDACTED]", out["echo"])
	require.Equal(t, "[REDACTED]", out["nested"].(map[string]any)["deep"])
	require.Equal(t, "[REDACTED]", out["list"].([]any)[0])
	require.Equal(t, "clean", out["list"].([]any)[1])
	require.Equal(t, 1.0, out["num"])
	// A secret placed in a KEY is scrubbed too: the value survives under the redacted key.
	require.Equal(t, "value under a leaked key", out["[REDACTED]"])
	_, leaked := out["SEKRET-abc123"]
	require.False(t, leaked, "a secret must not survive as a map key")
}

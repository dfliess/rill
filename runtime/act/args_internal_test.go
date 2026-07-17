package act

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFreezeArgsIsDeepCopy proves freezeArgs snapshots the argument containers so a later mutation of the caller's map
// (including nested maps and slices) cannot reach the frozen copy the ledger and executor bind to (§17.1).
func TestFreezeArgsIsDeepCopy(t *testing.T) {
	original := map[string]any{
		"summary": "coste alto",
		"labels":  []any{"p2", "cost"},
		"meta":    map[string]any{"owner": "alice"},
	}
	frozen := freezeArgs(original)

	// Mutate every level of the caller's map after freezing.
	original["summary"] = "TAMPERED"
	original["labels"].([]any)[0] = "TAMPERED"
	original["meta"].(map[string]any)["owner"] = "mallory"
	original["extra"] = "injected"

	require.Equal(t, "coste alto", frozen["summary"])
	require.Equal(t, []any{"p2", "cost"}, frozen["labels"])
	require.Equal(t, map[string]any{"owner": "alice"}, frozen["meta"])
	require.NotContains(t, frozen, "extra")
}

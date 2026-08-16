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

// TestIsCanonicalJSONProvesTheFormRatherThanTrustingIt covers the question a renderer has to answer before it can
// display a preimage without ambiguity. The answer must come from re-deriving the form, because free-form text can
// be valid JSON too and would otherwise borrow JSON's backslash rules without obeying them.
func TestIsCanonicalJSONProvesTheFormRatherThanTrustingIt(t *testing.T) {
	canonical, err := CanonicalizeArgs(map[string]any{"summary": "coste alto", "n": 3})
	require.NoError(t, err)
	require.True(t, IsCanonicalJSON(string(canonical)))

	// A large integer survives: UseNumber keeps its digits where float64 would round them and fail a canonical input.
	require.True(t, IsCanonicalJSON(`{"account_id":9007199254740993}`))
	require.True(t, IsCanonicalJSON(`{}`))

	for _, s := range []string{
		`crear ticket P2`,     // free-form proposal text
		`{"b":1,"a":2}`,       // valid JSON, keys not in canonical order
		`{ "a": 1 }`,          // valid JSON, not compact
		`["a"]`,               // valid JSON, not an object
		`{"a":1} trailing`,    // valid prefix, trailing junk
		`{"a":truth}`,         // not JSON at all
		`{"a":"ad\u200cmin"}`, // valid JSON, but with an escape where Marshal emits the character literally
	} {
		require.False(t, IsCanonicalJSON(s), "must not claim canonical JSON for %q", s)
	}
}

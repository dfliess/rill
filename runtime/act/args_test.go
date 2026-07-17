package act_test

import (
	"strings"
	"testing"

	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// TestCanonicalizeArgsStable verifies two argument maps with the same content but different literal key order produce
// identical canonical bytes (and thus the same hash), while a changed value produces a different hash. This is what
// lets an approval bind to arguments regardless of how the model happened to order them (§6.4).
func TestCanonicalizeArgsStable(t *testing.T) {
	a := map[string]any{
		"summary":  "coste alto",
		"priority": "P2",
		"labels":   []any{"cost", "alert"},
		"nested":   map[string]any{"z": 1.0, "a": 2.0},
	}
	b := map[string]any{
		"nested":   map[string]any{"a": 2.0, "z": 1.0},
		"labels":   []any{"cost", "alert"},
		"priority": "P2",
		"summary":  "coste alto",
	}

	hashA, err := act.HashCanonicalArgs(a)
	require.NoError(t, err)
	hashB, err := act.HashCanonicalArgs(b)
	require.NoError(t, err)
	require.Equal(t, hashA, hashB, "key order and map iteration order must not affect the hash")
	require.True(t, strings.HasPrefix(hashA, "sha256:"))

	// A changed value changes the hash: the approval would no longer match.
	c := map[string]any{"summary": "coste alto", "priority": "P1", "labels": []any{"cost", "alert"}, "nested": map[string]any{"z": 1.0, "a": 2.0}}
	hashC, err := act.HashCanonicalArgs(c)
	require.NoError(t, err)
	require.NotEqual(t, hashA, hashC)

	// Array order IS significant: reordering elements changes the hash.
	d := map[string]any{"summary": "coste alto", "priority": "P2", "labels": []any{"alert", "cost"}, "nested": map[string]any{"z": 1.0, "a": 2.0}}
	hashD, err := act.HashCanonicalArgs(d)
	require.NoError(t, err)
	require.NotEqual(t, hashA, hashD)
}

// TestCanonicalizeNilArgs verifies nil and empty maps canonicalize identically, so "no arguments" has one stable hash.
func TestCanonicalizeNilArgs(t *testing.T) {
	hashNil, err := act.HashCanonicalArgs(nil)
	require.NoError(t, err)
	hashEmpty, err := act.HashCanonicalArgs(map[string]any{})
	require.NoError(t, err)
	require.Equal(t, hashNil, hashEmpty)
}

// TestDeriveIdempotencyKey verifies the dedup key is stable for the same (run, tool call, args) and distinct when any
// of them changes (§12): a retry reuses the key so the target system dedups, while a re-approved, modified action gets
// a fresh key so it is a different external effect.
func TestDeriveIdempotencyKey(t *testing.T) {
	const run, call = "run-1", "call-1"
	hashA := act.HashArgs("args-A")
	hashB := act.HashArgs("args-B")

	k1 := act.DeriveIdempotencyKey(run, call, hashA)
	k2 := act.DeriveIdempotencyKey(run, call, hashA)
	require.Equal(t, k1, k2, "same inputs yield the same key (stable under retry)")
	require.True(t, strings.HasPrefix(k1, "act-"))

	require.NotEqual(t, k1, act.DeriveIdempotencyKey(run, call, hashB), "changed args → new key")
	require.NotEqual(t, k1, act.DeriveIdempotencyKey(run, "call-2", hashA), "changed tool call → new key")
	require.NotEqual(t, k1, act.DeriveIdempotencyKey("run-2", call, hashA), "changed run → new key")
}

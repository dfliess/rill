package act_test

import (
	"testing"

	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// TestRetryAttemptKeyRoundTrip pins the shape of the attempt chain. Attempt 0 must be the caller's key UNCHANGED:
// the derived form only ever appears after a failure, so a caller that never fails never sees it, and the key it
// supplied is the key its run is stored under.
func TestRetryAttemptKeyRoundTrip(t *testing.T) {
	const base = "36:f960bb37-9d8a-45df-8c4e-1e18a66f881e/agent-note-20a6380b"

	require.Equal(t, base, act.RetryAttemptKey(base, 0))

	k1 := act.NextRetryAttemptKey(base)
	require.NotEqual(t, base, k1)

	k2 := act.NextRetryAttemptKey(k1)
	require.NotEqual(t, k1, k2)

	// The chain stays inside one family rather than nesting: every attempt decomposes back to the same base.
	for i, k := range []string{base, k1, k2} {
		gotBase, attempt := act.SplitRetryAttemptKey(k)
		require.Equal(t, base, gotBase)
		require.Equal(t, i, attempt)
	}
}

// TestSplitRetryAttemptKeyIgnoresMalformedSuffix keeps a caller's key intact when its tail only looks like an
// attempt. Reading "…#retry-abc" as an attempt of "…" would put two unrelated callers in one family, and the
// second would inherit or skip the first one's run.
func TestSplitRetryAttemptKeyIgnoresMalformedSuffix(t *testing.T) {
	for _, key := range []string{"plain-key", "key#retry-", "key#retry-abc", "key#retry-0", "key#retry--1"} {
		base, attempt := act.SplitRetryAttemptKey(key)
		require.Equal(t, key, base, "key %q must be its own base", key)
		require.Equal(t, 0, attempt)
	}
}

// TestKeyFamilyLikePatternEscapesWildcards is the one that matters for isolation: a caller key is arbitrary text, so
// an unescaped '%' or '_' would turn the family lookup into a wildcard and let one key's retry resolve against a
// different key's runs.
func TestKeyFamilyLikePatternEscapesWildcards(t *testing.T) {
	require.Equal(t, `100\% sure#retry-%`, act.KeyFamilyLikePattern("100% sure"))
	require.Equal(t, `a\_b#retry-%`, act.KeyFamilyLikePattern("a_b"))
	require.Equal(t, `a\\b#retry-%`, act.KeyFamilyLikePattern(`a\b`))
}

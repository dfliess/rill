package act_test

import (
	"testing"

	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// closeRunAs drives a run to a terminal status so a test can set up "what came of this key".
func closeRunAs(t *testing.T, store *act.PostgresRunStore, instanceID, runID string, status act.RunStatus) {
	t.Helper()
	_, err := store.CloseRun(t.Context(), act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: status, EventType: act.EventTypeFailed,
	})
	require.NoError(t, err)
}

// createRunWithKey seeds a run of `agent` under `key`, addressed the way the executor addresses it.
func createRunWithKey(t *testing.T, store *act.PostgresRunStore, instanceID, agent, key string) string {
	t.Helper()
	runID := newRunID()
	require.NoError(t, store.CreateRun(t.Context(), act.NewRun{
		RunID: runID, InstanceID: instanceID, AgentName: agent, IdempotencyKey: key,
	}))
	return runID
}

// TestLatestRunInKeyFamily is the lookup that replaces probing: a caller must be able to ask what came of its key
// without starting anything, because starting an unseen key is what brings a run into being.
func TestLatestRunInKeyFamily(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-key-family"
	const agent = "narrador-home"
	const key = "36:f960bb37-9d8a-45df-8c4e-1e18a66f881e/agent-note-20a6380b"

	// An untouched key: nothing has come of it, so the caller starts under the key it supplied.
	_, err := store.LatestRunInKeyFamily(ctx, instanceID, agent, key)
	require.ErrorIs(t, err, act.ErrRunNotFound)

	first := createRunWithKey(t, store, instanceID, agent, key)
	closeRunAs(t, store, instanceID, first, act.RunStatusFailed)

	got, err := store.LatestRunInKeyFamily(ctx, instanceID, agent, key)
	require.NoError(t, err)
	require.Equal(t, first, got.RunID)
	require.Equal(t, act.RunStatusFailed, got.Status)

	// A retry of that key belongs to the same family and, being newer, is what the caller now sees.
	retryKey := act.NextRetryAttemptKey(key)
	second := createRunWithKey(t, store, instanceID, agent, retryKey)
	got, err = store.LatestRunInKeyFamily(ctx, instanceID, agent, key)
	require.NoError(t, err)
	require.Equal(t, second, got.RunID)

	// And the family is addressable from any of its members, not only from the base.
	got, err = store.LatestRunInKeyFamily(ctx, instanceID, agent, retryKey)
	require.NoError(t, err)
	require.Equal(t, second, got.RunID)
}

// TestLatestRunInKeyFamilyIsScoped keeps the lookup from reaching across the boundaries every other read respects.
// A family that leaked across instances would hand one tenant another tenant's answer.
func TestLatestRunInKeyFamilyIsScoped(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const key = "shared-key"
	createRunWithKey(t, store, "inst-a", "narrador", key)

	_, err := store.LatestRunInKeyFamily(ctx, "inst-b", "narrador", key)
	require.ErrorIs(t, err, act.ErrRunNotFound, "another instance must not see this family")

	_, err = store.LatestRunInKeyFamily(ctx, "inst-a", "otro-agente", key)
	require.ErrorIs(t, err, act.ErrRunNotFound, "another agent must not see this family")
}

// TestLatestRunInKeyFamilyDoesNotWidenOnWildcards is the escaping test with teeth: a caller key is arbitrary text,
// so an unescaped '%' would make one key's family swallow every other key's runs.
func TestLatestRunInKeyFamilyDoesNotWidenOnWildcards(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-wildcard"
	const agent = "narrador"

	// A run under an unrelated key that a wildcard match WOULD pick up.
	createRunWithKey(t, store, instanceID, agent, "prefix-victim#retry-1")

	// The attacker key ends in '%', which as a raw LIKE prefix matches anything starting with "prefix-".
	_, err := store.LatestRunInKeyFamily(ctx, instanceID, agent, "prefix-%")
	require.ErrorIs(t, err, act.ErrRunNotFound, "a '%' in the key must be a literal, not a wildcard")
}

package act

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCanonicalArgsForApprovalRecoversAnOldCheckpoint covers the one case where the preimage does not come straight
// off the Authorization: a run whose propose step was checkpointed by a build from before that field existed. DBOS
// replays the old checkpoint, the field deserializes empty, and without a fallback the store's binding check would
// refuse the approval and kill a legitimate in-flight run mid-deploy.
//
// The fallback rebuilds from the proposal's own arguments and is only allowed to matter when the rebuild hashes back
// to the authorization's hash: a rebuild that matches IS the preimage, whatever produced it. When it does not match,
// it returns nothing and lets the store refuse — the run stops, which is the right outcome when the bytes behind a
// signature cannot be established.
func TestCanonicalArgsForApprovalRecoversAnOldCheckpoint(t *testing.T) {
	args := map[string]any{"summary": "coste alto", "project": "OPS"}
	canonical, err := CanonicalizeArgs(args)
	require.NoError(t, err)

	// The normal path: the gateway's own capture is used as-is.
	fresh := &governedAction{
		proposal: ToolProposal{Args: args},
		auth:     Authorization{ArgsHash: hashBytes(canonical), CanonicalArgs: string(canonical)},
	}
	require.Equal(t, string(canonical), canonicalArgsForApproval(fresh))

	// The replayed checkpoint: no preimage on the authorization, rebuilt from the proposal and accepted because it
	// hashes back to the hash that was checkpointed with it.
	replayed := &governedAction{
		proposal: ToolProposal{Args: args},
		auth:     Authorization{ArgsHash: hashBytes(canonical)},
	}
	require.Equal(t, string(canonical), canonicalArgsForApproval(replayed))

	// A proposal whose arguments no longer produce the checkpointed hash yields nothing rather than a preimage that
	// contradicts the signature. This is the case that must never silently succeed.
	drifted := &governedAction{
		proposal: ToolProposal{Args: map[string]any{"summary": "coste BAJO", "project": "OPS"}},
		auth:     Authorization{ArgsHash: hashBytes(canonical)},
	}
	require.Empty(t, canonicalArgsForApproval(drifted))

	// Arguments that cannot be canonicalized at all take the same safe exit.
	unmarshalable := &governedAction{
		proposal: ToolProposal{Args: map[string]any{"ch": make(chan int)}},
		auth:     Authorization{ArgsHash: hashBytes(canonical)},
	}
	require.Empty(t, canonicalArgsForApproval(unmarshalable))
}

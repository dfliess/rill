package act

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// CanonicalizeArgs renders tool arguments to a stable byte form so their hash is reproducible across the proposal,
// the approval and the execution of one action. encoding/json marshals object keys in sorted order (recursively,
// including nested maps), so a given argument value always encodes to identical bytes regardless of the map's Go
// iteration order; array order is preserved because it is semantically significant. A nil map canonicalizes to the
// empty object, so "no arguments" has one stable representation.
func CanonicalizeArgs(args map[string]any) ([]byte, error) {
	if args == nil {
		return []byte("{}"), nil
	}
	b, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("act: canonicalize args: %w", err)
	}
	return b, nil
}

// freezeArgs returns a deep copy of args so the proposal's hash, its persisted ledger form and the arguments finally
// executed are all bound to one immutable snapshot: a Proposer that keeps and mutates its own map after proposing
// cannot change what was hashed, approved or sent to the target system (§17.1). Only the containers (maps and slices)
// are copied; scalars are immutable in Go and shared safely. A nil map stays nil.
func freezeArgs(args map[string]any) map[string]any {
	if args == nil {
		return nil
	}
	frozen, _ := deepCopyArgValue(args).(map[string]any)
	return frozen
}

// deepCopyArgValue recursively copies the map and slice containers of a JSON-like value, returning scalars unchanged.
func deepCopyArgValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, elem := range t {
			out[k] = deepCopyArgValue(elem)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, elem := range t {
			out[i] = deepCopyArgValue(elem)
		}
		return out
	default:
		return v
	}
}

// HashCanonicalArgs returns the canonical hash of a tool's arguments. This is the value an approval is bound to
// (§6.4, §11.2): a decision is valid only for arguments that hash to exactly this, so changing any argument after
// approval invalidates the decision and forces a fresh request.
func HashCanonicalArgs(args map[string]any) (string, error) {
	b, err := CanonicalizeArgs(args)
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

// hashBytes is the single hashing primitive: a length-stable, prefixed SHA-256 digest. The "sha256:" prefix names
// the algorithm inline so a stored hash stays self-describing if the algorithm ever changes.
func hashBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// DeriveIdempotencyKey builds the stable dedup key for one external write (§12). It folds the run, the tool call and
// the canonical args hash into a single opaque key so it is:
//
//   - stable under retry: the same (run, tool call, args) always yields the same key, so an at-least-once re-issue of
//     the write dedups at the target system rather than creating a second effect;
//   - distinct on change: editing an argument changes the args hash and therefore the key, so a re-approved, modified
//     action is treated as a different external effect and never collides with the original (§11.2).
//
// It carries no argument values in the clear, so it is safe to pass to the target system and to persist.
func DeriveIdempotencyKey(runID, toolCallID, argsHash string) string {
	sum := sha256.Sum256([]byte(runID + "\x00" + toolCallID + "\x00" + argsHash))
	// Half the digest is ample collision resistance for a dedup key and keeps it compact in logs and external headers.
	return "act-" + hex.EncodeToString(sum[:16])
}

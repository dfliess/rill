package act

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
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

// IsCanonicalJSON reports whether s is exactly what CanonicalizeArgs would produce for the object it encodes.
//
// It exists because a renderer has to know whether a preimage is JSON before it can display it faithfully: inside
// json.Marshal's output a literal backslash is already doubled, so escapes are unambiguous as they are, while in
// free-form text they are not and the renderer has to double them itself. Getting that backwards makes two different
// signed texts look identical.
//
// The question is answered by re-deriving rather than by trusting where the bytes came from: only json.Marshal
// produces this exact form, so a round-trip that reproduces s byte for byte IS the proof. UseNumber keeps the
// original digits of a large integer, which float64 would round and which would then fail the comparison for a
// perfectly canonical input.
//
// Domain note: it answers the question for JSON-like values, which is what json.Unmarshal produces and what the MCP
// path carries. A map holding a json.RawMessage whose bytes are already non-canonical would round-trip to different
// bytes and be reported false. That is a limit of the contract rather than a hazard — no caller builds args that way,
// and the direction of the error is the safe one.
func IsCanonicalJSON(s string) bool {
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()
	var v map[string]any
	if err := dec.Decode(&v); err != nil {
		return false
	}
	if dec.More() {
		return false
	}
	round, err := CanonicalizeArgs(v)
	if err != nil {
		return false
	}
	return string(round) == s
}

// Package secretscrub replaces resolved secret values with an opaque marker in strings and JSON-like values. It is the
// last line of defense against a tool result, an error, or a log line echoing a credential back in the clear.
//
// It lives in its own leaf package so the two sides that must scrub — the action gateway (runtime/act) and the
// model/tool loop (runtime/ai) — share one implementation instead of each carrying a copy. The logic cannot live in
// either package: runtime/act imports runtime/ai, so a scrubber in runtime/act is unreachable from runtime/ai (import
// cycle), and a copy in runtime/ai would be a second redactor to keep in sync.
package secretscrub

import "strings"

// Marker replaces any secret value that would otherwise be persisted, returned to the model, or logged. It is
// intentionally opaque: it reveals neither the secret nor its length.
const Marker = "[REDACTED]"

// String replaces every occurrence of every secret value with Marker. An empty secret is skipped, so a connector with
// no credential does not turn every string into the marker.
func String(s string, secretValues []string) string {
	for _, v := range secretValues {
		if v == "" {
			continue
		}
		s = strings.ReplaceAll(s, v, Marker)
	}
	return s
}

// Value walks any JSON-like value and scrubs secret values from every string it contains, returning a copy. Map KEYS
// are scrubbed as well as values: a hostile connector can place a leaked credential in a key as easily as in a value,
// so both run through String.
func Value(v any, secretValues []string) any {
	switch t := v.(type) {
	case string:
		return String(t, secretValues)
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, elem := range t {
			out[String(k, secretValues)] = Value(elem, secretValues)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, elem := range t {
			out[i] = Value(elem, secretValues)
		}
		return out
	default:
		return v
	}
}

// Map returns a copy of m with every secret value scrubbed from every nested string AND every map key. A nil map stays
// nil.
func Map(m map[string]any, secretValues []string) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[String(k, secretValues)] = Value(v, secretValues)
	}
	return out
}

package act

import (
	"context"
	"fmt"
	"regexp"

	"github.com/rilldata/rill/runtime/pkg/secretscrub"
)

// secretRefPattern matches a value that is entirely a `{{ secret.NAME }}` marker. The whole-string anchoring is
// deliberate: only a value the agent author placed as a standalone reference is resolved, so a secret can never be
// smuggled in as a substring of model-controlled text and then expanded.
var secretRefPattern = regexp.MustCompile(`^\{\{\s*secret\.([a-zA-Z0-9_.-]+)\s*\}\}$`)

// SecretResolver resolves a named secret to its plaintext, server-side. The gateway uses the result only to build
// the executor call and to scrub outputs; a resolved value is never persisted to the ledger, never returned to the
// model, and never written to a log (§16.2). Production backs this with the platform secrets manager (ADR-0010);
// tests back it with an in-memory map.
type SecretResolver interface {
	Resolve(ctx context.Context, name string) (string, error)
}

// SecretRef reports the secret name a value references, if the value is a standalone secret reference. A non-string
// value, or a string that only partially matches the marker, is treated as literal data and returns ok=false.
func SecretRef(v any) (name string, ok bool) {
	s, isString := v.(string)
	if !isString {
		return "", false
	}
	m := secretRefPattern.FindStringSubmatch(s)
	if m == nil {
		return "", false
	}
	return m[1], true
}

// resolveSecrets returns a deep copy of args with every secret reference replaced by its resolved plaintext, together
// with the distinct set of resolved values for output scrubbing. The input map is never mutated: the reference form
// (`{{ secret.NAME }}`) stays in the arguments the gateway persists and shows a human, while the resolved form exists
// only in the returned copy handed to the ActionExecutor. A reference to an unknown secret is a hard error — the
// action cannot execute with a hole where a credential should be (fail-closed).
func resolveSecrets(ctx context.Context, r SecretResolver, args map[string]any) (resolved map[string]any, secretValues []string, err error) {
	seen := make(map[string]struct{})
	var values []string
	record := func(v string) {
		if v == "" {
			return
		}
		if _, dup := seen[v]; dup {
			return
		}
		seen[v] = struct{}{}
		values = append(values, v)
	}

	var walk func(v any) (any, error)
	walk = func(v any) (any, error) {
		if name, ok := SecretRef(v); ok {
			if r == nil {
				return nil, fmt.Errorf("act: secret %q referenced but no resolver configured", name)
			}
			secret, resolveErr := r.Resolve(ctx, name)
			if resolveErr != nil {
				return nil, fmt.Errorf("act: resolve secret %q: %w", name, resolveErr)
			}
			record(secret)
			return secret, nil
		}
		switch t := v.(type) {
		case map[string]any:
			out := make(map[string]any, len(t))
			for k, elem := range t {
				sub, walkErr := walk(elem)
				if walkErr != nil {
					return nil, walkErr
				}
				out[k] = sub
			}
			return out, nil
		case []any:
			out := make([]any, len(t))
			for i, elem := range t {
				sub, walkErr := walk(elem)
				if walkErr != nil {
					return nil, walkErr
				}
				out[i] = sub
			}
			return out, nil
		default:
			return v, nil
		}
	}

	if args == nil {
		return nil, nil, nil
	}
	out := make(map[string]any, len(args))
	for k, v := range args {
		sub, walkErr := walk(v)
		if walkErr != nil {
			return nil, nil, walkErr
		}
		out[k] = sub
	}
	return out, values, nil
}

// redactString replaces every occurrence of any secret value with the redaction marker. It is the last line of
// defense: a tool result that echoes a credential back (a misbehaving or malicious connector) must not reach the
// ledger, an approval inbox or a log in the clear (§17.1, §17.2). The scrubbing logic lives in the neutral
// secretscrub package so runtime/ai (which runtime/act imports, precluding the reverse) can share it.
func redactString(s string, secretValues []string) string {
	return secretscrub.String(s, secretValues)
}

// RedactMap returns a copy of m with every secret value scrubbed from every nested string AND every map key. A nil map
// stays nil.
func RedactMap(m map[string]any, secretValues []string) map[string]any {
	return secretscrub.Map(m, secretValues)
}

package act

import (
	"strconv"
	"strings"
)

// retryAttemptSeparator marks a retry attempt of a caller's idempotency key. A run that FAILED is not an answer,
// so a caller that resupplies its key must not inherit it (issue #154); the next attempt gets a derived key, since
// DBOS cannot re-execute a workflow ID that already reached a terminal state.
//
// The separator is deliberately unlikely to appear in a caller's own key: keys sharing a family must be attempts
// of one another, and a caller that literally supplies "…#retry-3" would have its key read as an attempt of
// "…". Callers are not expected to mint keys in this shape, and the API scopes every caller key to its actor
// before this convention applies, so the aliasing cannot cross users.
const retryAttemptSeparator = "#retry-"

// RetryAttemptKey returns the key for the nth retry attempt of base. Attempt 0 is base itself, so the first run of
// a key is stored under exactly the key the caller supplied and the derived form only ever appears after a failure.
func RetryAttemptKey(base string, attempt int) string {
	if attempt <= 0 {
		return base
	}
	return base + retryAttemptSeparator + strconv.Itoa(attempt)
}

// SplitRetryAttemptKey decomposes a key into its family base and attempt number. A key with no attempt suffix is
// attempt 0 of itself.
func SplitRetryAttemptKey(key string) (base string, attempt int) {
	i := strings.LastIndex(key, retryAttemptSeparator)
	if i < 0 {
		return key, 0
	}
	n, err := strconv.Atoi(key[i+len(retryAttemptSeparator):])
	if err != nil || n <= 0 {
		return key, 0
	}
	return key[:i], n
}

// NextRetryAttemptKey returns the key for the attempt after key's own.
func NextRetryAttemptKey(key string) string {
	base, attempt := SplitRetryAttemptKey(key)
	return RetryAttemptKey(base, attempt+1)
}

// KeyFamilyLikePattern builds the SQL LIKE pattern matching the retry attempts of base, for use with ESCAPE '\'.
// The base is escaped because a caller's key is arbitrary text: an unescaped '%' or '_' in it would widen the
// match to other callers' keys.
func KeyFamilyLikePattern(base string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(base) + retryAttemptSeparator + `%`
}

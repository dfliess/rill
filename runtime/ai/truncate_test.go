package ai

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/rilldata/rill/runtime/act/mcpconn"
	"github.com/stretchr/testify/require"
)

// Regression for kairosagentica/rill#14: byte-indexed truncation can split a multi-byte UTF-8
// rune, producing an invalid string that Postgres rejects when the act ledger persists the
// proposal (SQLSTATE 22021, "invalid byte sequence for encoding UTF8: 0xc3 0xe2" — the dangling
// first byte of '×' followed by the first byte of the appended '…').
func TestProposedSummaryValidUTF8OnMultibyteBoundary(t *testing.T) {
	rt := mcpconn.RemoteTool{RawName: "create_ticket"}
	// 119 ASCII bytes + '×' (2 bytes): the 120-byte cut of truncateArgValue lands mid-rune.
	args := map[string]any{
		"title": strings.Repeat("a", 119) + "×33cl en Valencia Centro",
	}
	got := proposedSummary(rt, args)
	require.True(t, utf8.ValidString(got), "proposal summary must be valid UTF-8, got %q", got)
	require.Contains(t, got, "…", "a cut value must be visibly marked as truncated")
}

// TestProposedSummaryKeepsRealisticValuesIntact pins the per-value cap's intent: a realistic one-line value (a
// title, subject or URL of up to 120 bytes) survives whole in the inbox summary, while a body-sized value is cut
// and visibly marked. The exact arguments are persisted separately on the approval; the summary only has to be an
// honest one-liner.
func TestProposedSummaryKeepsRealisticValuesIntact(t *testing.T) {
	rt := mcpconn.RemoteTool{RawName: "create_ticket"}
	title := strings.Repeat("t", 120)
	body := strings.Repeat("b", 500)
	got := proposedSummary(rt, map[string]any{"body": body, "title": title})
	require.Contains(t, got, "title="+title, "a 120-byte value must not be cut")
	require.Contains(t, got, "body="+strings.Repeat("b", 120)+"…", "an over-long value is cut at 120 bytes and marked")
	require.NotContains(t, got, strings.Repeat("b", 121), "no more than the cap of an over-long value may leak into the one-liner")
}

func TestPromptToTitleValidUTF8OnMultibyteBoundary(t *testing.T) {
	// 46 ASCII bytes + 'ú' (2 bytes) …: the 47-byte cut lands mid-rune.
	prompt := strings.Repeat("a", 46) + "úúúú"
	got := promptToTitle(prompt)
	require.True(t, utf8.ValidString(got), "conversation title must be valid UTF-8, got %q", got)
}

func TestTruncateUTF8(t *testing.T) {
	cases := []struct {
		name string
		in   string
		max  int
		want string
	}{
		{"short ascii unchanged", "hola", 40, "hola"},
		{"exact length unchanged", strings.Repeat("a", 40), 40, strings.Repeat("a", 40)},
		{"ascii cut at boundary", strings.Repeat("a", 41), 40, strings.Repeat("a", 40)},
		{"two-byte rune at boundary backs off", strings.Repeat("a", 39) + "××", 40, strings.Repeat("a", 39)},
		{"three-byte rune at boundary backs off", strings.Repeat("a", 38) + "€€", 40, strings.Repeat("a", 38)},
		{"rune ending exactly at boundary kept", strings.Repeat("a", 38) + "×€", 40, strings.Repeat("a", 38) + "×"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := truncateUTF8(c.in, c.max)
			require.Equal(t, c.want, got)
			require.True(t, utf8.ValidString(got))
		})
	}
}

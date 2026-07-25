package ai

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/rilldata/rill/runtime/act/mcpconn"
	"github.com/stretchr/testify/require"
)

// Regression for kairosagentica/rill#14: byte-indexed truncation (s[:40]) can split a multi-byte
// UTF-8 rune, producing an invalid string that Postgres rejects when the act ledger persists the
// proposal (SQLSTATE 22021, "invalid byte sequence for encoding UTF8: 0xc3 0xe2" — the dangling
// first byte of '×' followed by the first byte of the appended '…').
func TestProposedSummaryValidUTF8OnMultibyteBoundary(t *testing.T) {
	rt := mcpconn.RemoteTool{RawName: "create_ticket"}
	// 39 ASCII bytes + '×' (2 bytes): the 40-byte cut lands mid-rune.
	args := map[string]any{
		"title": strings.Repeat("a", 39) + "×33cl en Valencia Centro",
	}
	got := proposedSummary(rt, args)
	require.True(t, utf8.ValidString(got), "proposal summary must be valid UTF-8, got %q", got)
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

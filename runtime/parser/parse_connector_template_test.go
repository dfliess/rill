package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/rilldata/rill/runtime/drivers/openai"
)

func TestOpenAIExtraBodyRejectsEnvironmentTemplates(t *testing.T) {
	repo := makeRepo(t, map[string]string{
		"rill.yaml": "",
		"connectors/deepseek.yaml": `
type: connector
driver: openai
api_key: "{{ .env.DEEPSEEK_API_KEY }}"
extra_body:
  thinking:
    type: disabled
  vendor_token: "{{ .env.VENDOR_TOKEN }}"
`,
	})

	p, err := Parse(context.Background(), repo, "test", "", "duckdb", true)
	require.NoError(t, err)
	require.Len(t, p.Errors, 1)
	require.Contains(t, p.Errors[0].Message, `property "extra_body" does not allow templates`)
}

func TestOpenAIExtraBodyAllowsStaticBehavior(t *testing.T) {
	repo := makeRepo(t, map[string]string{
		"rill.yaml": "",
		"connectors/deepseek.yaml": `
type: connector
driver: openai
api_key: "{{ .env.DEEPSEEK_API_KEY }}"
extra_body:
  thinking:
    type: disabled
`,
	})

	p, err := Parse(context.Background(), repo, "test", "", "duckdb", true)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Empty(t, p.Errors)
}

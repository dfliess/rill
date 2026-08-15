package parser

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestJSONSchemaForRillYAML(t *testing.T) {
	schema, err := JSONSchemaForRillYAML()
	require.NoError(t, err)
	require.NotNil(t, schema)
	require.Equal(t, "Project YAML", schema.Title)
	require.NotEmpty(t, schema.AllOf, "schema should have properties")
}

func TestJSONSchemaForResourceType(t *testing.T) {
	tests := []struct {
		name         string
		resourceType ResourceKind
		wantTitle    string
		wantAllOf    bool
		wantErr      bool
	}{
		{
			name:         "model schema",
			resourceType: ResourceKindModel,
			wantTitle:    "Models YAML",
			wantAllOf:    true,
		},
		{
			name:         "metrics view schema",
			resourceType: ResourceKindMetricsView,
			wantTitle:    "Metrics View YAML",
			wantAllOf:    true,
		},
		{
			name:         "connector schema",
			resourceType: ResourceKindConnector,
			wantTitle:    "Connector YAML",
			wantAllOf:    true,
		},
		{
			name:         "explore schema",
			resourceType: ResourceKindExplore,
			wantTitle:    "Explore Dashboard YAML",
			wantAllOf:    true,
		},
		{
			name:         "canvas schema",
			resourceType: ResourceKindCanvas,
			wantTitle:    "Canvas Dashboard YAML",
			wantAllOf:    true,
		},
		{
			name:         "alert schema",
			resourceType: ResourceKindAlert,
			wantTitle:    "Alert YAML",
			wantAllOf:    true,
		},
		{
			name:         "theme schema",
			resourceType: ResourceKindTheme,
			wantTitle:    "Theme YAML",
			wantAllOf:    true,
		},
		{
			name:         "api schema",
			resourceType: ResourceKindAPI,
			wantTitle:    "API YAML",
			wantAllOf:    true,
		},
		{
			name:         "component schema",
			resourceType: ResourceKindComponent,
			wantTitle:    "Component YAML",
			wantAllOf:    true,
		},
		{
			name:         "agent schema",
			resourceType: ResourceKindAgent,
			wantTitle:    "Agent YAML",
			wantAllOf:    true,
		},
		{
			name:         "agent trigger schema",
			resourceType: ResourceKindAgentTrigger,
			wantTitle:    "Agent Trigger YAML",
			wantAllOf:    true,
		},
		{
			name:         "unsupported resource type",
			resourceType: ResourceKindMigration,
			wantErr:      true,
		},
		{
			name:         "unspecified resource type",
			resourceType: ResourceKindUnspecified,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, err := JSONSchemaForResourceType(tt.resourceType)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, schema)
			require.Equal(t, tt.wantTitle, schema.Title)

			if tt.wantAllOf {
				require.NotEmpty(t, schema.AllOf, "schema should have allOf")
			}
		})
	}
}

// TestAgentSchemaExamplesParse feeds the examples embedded in project.schema.yaml for the agent resource kinds
// through the real parser, so the schema's documented surface cannot drift from what parseAgent and
// parseAgentTrigger accept. The examples are extracted from the embedded schema rather than duplicated here:
// editing an example in the schema is enough to re-test it.
func TestAgentSchemaExamplesParse(t *testing.T) {
	// Decode only the per-definition examples; yaml.Node keeps each example verbatim regardless of its shape.
	var doc struct {
		Definitions map[string]struct {
			Examples []yaml.Node `yaml:"examples"`
		} `yaml:"definitions"`
	}
	require.NoError(t, yaml.Unmarshal([]byte(resourceYAMLSchema), &doc))

	files := map[string]string{`rill.yaml`: ``}
	for defKey, prefix := range map[string]string{"agents": "agent", "agent-triggers": "trigger"} {
		examples := doc.Definitions[defKey].Examples
		require.NotEmpty(t, examples, "definition %q should embed at least one example", defKey)
		for i := range examples {
			content, err := yaml.Marshal(&examples[i])
			require.NoError(t, err)
			files[fmt.Sprintf("agents/example_%s_%d.yaml", prefix, i)] = string(content)

			// A trigger example references an agent by name; add a minimal stub for it so the repo is
			// self-consistent (the parser records the ref either way, but the examples should not dangle).
			var body struct {
				Agent string `yaml:"agent"`
			}
			require.NoError(t, examples[i].Decode(&body))
			if body.Agent != "" {
				if _, ok := files["agents/"+body.Agent+".yaml"]; !ok {
					files["agents/"+body.Agent+".yaml"] = "type: agent\ninstructions: stub\n"
				}
			}
		}
	}

	p, err := Parse(context.Background(), makeRepo(t, files), "", "", "duckdb", true)
	require.NoError(t, err)
	require.Empty(t, p.Errors, "schema examples must parse without errors: %v", p.Errors)
}

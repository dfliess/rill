package ai

import (
	"context"

	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act/mcpconn"
)

// PreloadSkills exposes preloadSkills to the external test package.
var PreloadSkills = preloadSkills

// SetNewMCPClientForTest overrides the outbound MCP client factory and returns a function that restores the
// previous factory. It is the seam an external test uses to point a dynamic agent's connectors at an in-process
// server (granting loopback/http via mcpconn.Options.AllowPrivateNetworks), which production never does.
func SetNewMCPClientForTest(fn func(cfg *mcpconn.ConnectorConfig, token string) (mcpconn.MCPClient, error)) func() {
	prev := newMCPClient
	newMCPClient = fn
	return func() { newMCPClient = prev }
}

// PersonalCanvasPath exposes personalCanvasPath to the external test package.
var PersonalCanvasPath = personalCanvasPath

// PersonalCanvasStatus exposes the status write_personal_canvas reports for a personal canvas to the external test package.
func PersonalCanvasStatus(ctx context.Context, rt *runtime.Runtime, path string) ([]map[string]any, string, error) {
	resources, parseErr, _, err := (&WriteFile{Runtime: rt}).reconcileAndGetStatus(ctx, path, true)
	return resources, parseErr, err
}

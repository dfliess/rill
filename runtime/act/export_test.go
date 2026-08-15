package act

import (
	"context"

	"github.com/rilldata/rill/runtime/ai"
)

// MCPConnectorFromContextForTest exposes mcpConnectorFromContext to the black-box act_test package, so a flow test
// can assert the segmented workflow binds the run's captured connector onto the governed execute step context
// (dbos_executor's withGovernedConnector call). Without that binding the runtime resolver would fall back to the
// live agent definition, undoing the §8.3 freeze the wiring tests pin at the connectorFor level.
func MCPConnectorFromContextForTest(ctx context.Context) (ai.MCPConnector, bool) {
	return mcpConnectorFromContext(ctx)
}

package server

import (
	"strings"
	"testing"
	"time"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestScopeKeyToActorSeparatesUsers(t *testing.T) {
	const key = "agent-note-1f2e3d4c"

	// The point of the whole exercise: one dashboard, one note, two readers, two runs. Sharing one would serve
	// the second reader an answer computed under the first one's access.
	require.NotEqual(t, scopeKeyToActor("user-a", key), scopeKeyToActor("user-b", key))

	// And the same reader must still deduplicate, or every revisit would pay for a fresh completion.
	require.Equal(t, scopeKeyToActor("user-a", key), scopeKeyToActor("user-a", key))
}

func TestScopeKeyToActorIsInjective(t *testing.T) {
	// Without the length prefix these two would both render as "a/b/c" and collide onto one run, handing one
	// user the other's answer. The key is caller-supplied, so it is attacker-shaped: assume it contains anything.
	require.NotEqual(t, scopeKeyToActor("a", "b/c"), scopeKeyToActor("a/b", "c"))

	// An anonymous reader of a public project has no id to scope by, so those readers keep sharing one run.
	require.Equal(t, scopeKeyToActor("", "k"), scopeKeyToActor("", "k"))
	require.NotEqual(t, scopeKeyToActor("", "k"), scopeKeyToActor("someone", "k"))
}

func TestPromptWithDashboardContextPassesThroughWithoutContext(t *testing.T) {
	require.Equal(t, "narra el cuadro", promptWithDashboardContext("narra el cuadro", nil))

	// An empty context carries nothing to scope by, so it must not dress the prompt in an empty preamble.
	require.Equal(t, "narra el cuadro", promptWithDashboardContext("narra el cuadro", &runtimev1.AnalystAgentContext{}))
}

func TestPromptWithDashboardContextRendersState(t *testing.T) {
	start := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)

	got := promptWithDashboardContext("narra el cuadro", &runtimev1.AnalystAgentContext{
		Canvas:    "direccion_canvas",
		TimeStart: timestamppb.New(start),
		TimeEnd:   timestamppb.New(end),
		Measures:  []string{"ingresos", "margen_docente"},
	})

	require.Contains(t, got, "dashboard: direccion_canvas")
	require.Contains(t, got, "2026-06-01T00:00:00Z to 2026-07-01T00:00:00Z")
	require.Contains(t, got, "measures: ingresos, margen_docente")
	// The caller's prompt stays intact and last, so the agent's task is the final thing the model reads.
	require.True(t, strings.HasSuffix(got, "narra el cuadro"))
}

// The client derives its idempotency key from the same context, so an unchanged dashboard must render a byte-identical
// prompt. Map iteration would otherwise reorder the per-metrics-view filters and silently break run reuse.
func TestPromptWithDashboardContextIsDeterministic(t *testing.T) {
	ctx := &runtimev1.AnalystAgentContext{
		Canvas: "direccion_canvas",
		WherePerMetricsView: map[string]*runtimev1.Expression{
			"finanzas_metrics":  exprIdent(t, "campus"),
			"academico_metrics": exprIdent(t, "area_academica"),
			"metas_metrics":     exprIdent(t, "nivel"),
			"ocupacion_metrics": exprIdent(t, "turno"),
		},
	}

	first := promptWithDashboardContext("narra el cuadro", ctx)
	for i := 0; i < 20; i++ {
		require.Equal(t, first, promptWithDashboardContext("narra el cuadro", ctx))
	}

	// Sorted by metrics view name, not by map order.
	require.Less(t, strings.Index(first, "academico_metrics"), strings.Index(first, "finanzas_metrics"))
	require.Less(t, strings.Index(first, "finanzas_metrics"), strings.Index(first, "metas_metrics"))
	require.Less(t, strings.Index(first, "metas_metrics"), strings.Index(first, "ocupacion_metrics"))
}

func exprIdent(t *testing.T, name string) *runtimev1.Expression {
	t.Helper()
	return &runtimev1.Expression{
		Expression: &runtimev1.Expression_Ident{Ident: name},
	}
}

// TestStreamAgentRunEventsGetsAWatchTimeout pins the server-side deadline for the run event stream.
//
// It shipped without an entry in timeoutSelector, so it fell through to the 30-second default meant for unary calls.
// A watch is at its most useful when it has nothing to send — a run parked on a human approval emits no events while
// it waits — so the interceptor cancelled it mid-wait and the timeline read "live events unavailable: context
// deadline exceeded" on exactly the runs someone had open. Asserting it against the default is what makes the
// omission visible: a new streaming method that forgets this fails here rather than in front of a user.
func TestStreamAgentRunEventsGetsAWatchTimeout(t *testing.T) {
	unaryDefault := timeoutSelector("/rill.runtime.v1.AgentService/GetAgentRun")
	streamTimeout := timeoutSelector(runtimev1.AgentService_StreamAgentRunEvents_FullMethodName)

	require.NotEqual(t, unaryDefault, streamTimeout,
		"the run event stream must not inherit the unary default: it is a watch that legitimately sends nothing for long stretches")
	require.GreaterOrEqual(t, streamTimeout, 30*time.Minute,
		"a human deciding an approval takes minutes to hours, so the stream's deadline must be on that scale")

	// The same scale as Rill's own watches, which is where this belongs.
	require.Equal(t, timeoutSelector(runtimev1.RuntimeService_WatchResources_FullMethodName), streamTimeout)
}

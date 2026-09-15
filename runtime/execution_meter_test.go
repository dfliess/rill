package runtime

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/c2h5oh/datasize"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/pkg/email"
	"github.com/rilldata/rill/runtime/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestExecutionMetric(t *testing.T) {
	alert := func(status runtimev1.AssertionStatus, adhoc bool) *runtimev1.Resource {
		return &runtimev1.Resource{
			Meta: &runtimev1.ResourceMeta{Name: &runtimev1.ResourceName{Kind: ResourceKindAlert, Name: "stock-break"}},
			Resource: &runtimev1.Resource_Alert{Alert: &runtimev1.Alert{State: &runtimev1.AlertState{
				ExecutionHistory: []*runtimev1.AlertExecution{
					{Adhoc: adhoc, Result: &runtimev1.AssertionResult{Status: status}},
					// An older execution must not influence the metric: only the one just persisted counts.
					{Result: &runtimev1.AssertionResult{Status: runtimev1.AssertionStatus_ASSERTION_STATUS_PASS}},
				},
			}}},
		}
	}
	report := func(errMsg string, adhoc bool) *runtimev1.Resource {
		return &runtimev1.Resource{
			Meta: &runtimev1.ResourceMeta{Name: &runtimev1.ResourceName{Kind: ResourceKindReport, Name: "weekly-collections"}},
			Resource: &runtimev1.Resource_Report{Report: &runtimev1.Report{State: &runtimev1.ReportState{
				ExecutionHistory: []*runtimev1.ReportExecution{{Adhoc: adhoc, ErrorMessage: errMsg}},
			}}},
		}
	}

	tests := []struct {
		name     string
		res      *runtimev1.Resource
		metric   string
		source   string
		resource string
		adhoc    bool
	}{
		{"alert passed", alert(runtimev1.AssertionStatus_ASSERTION_STATUS_PASS, false), "alert_execution_passed", "alert", "stock-break", false},
		{"alert failed, manual", alert(runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL, true), "alert_execution_failed", "alert", "stock-break", true},
		{"alert error", alert(runtimev1.AssertionStatus_ASSERTION_STATUS_ERROR, false), "alert_execution_error", "alert", "stock-break", false},
		{"alert without result counts as error", alert(runtimev1.AssertionStatus_ASSERTION_STATUS_UNSPECIFIED, false), "alert_execution_error", "alert", "stock-break", false},
		{"report sent, manual", report("", true), "report_execution_sent", "report", "weekly-collections", true},
		{"report error", report("Report run failed: boom", false), "report_execution_error", "report", "weekly-collections", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric, attrs, ok := executionMetric(tt.res)
			require.True(t, ok)
			require.Equal(t, tt.metric, metric)

			got := map[string]any{}
			for _, a := range attrs {
				got[string(a.Key)] = a.Value.AsInterface()
			}
			require.Equal(t, map[string]any{"source": tt.source, "resource_name": tt.resource, "adhoc": tt.adhoc}, got)
		})
	}

	t.Run("no executions", func(t *testing.T) {
		_, _, ok := executionMetric(&runtimev1.Resource{Resource: &runtimev1.Resource_Alert{Alert: &runtimev1.Alert{State: &runtimev1.AlertState{}}}})
		require.False(t, ok)
	})
	t.Run("other resource kinds", func(t *testing.T) {
		_, _, ok := executionMetric(&runtimev1.Resource{Resource: &runtimev1.Resource_Model{Model: &runtimev1.Model{}}})
		require.False(t, ok)
	})
}

func TestObserveExecutionMetersWithInstanceAnnotations(t *testing.T) {
	sink := &captureSink{}
	rt := newMeterTestRuntime(t, activity.NewClient(sink, zap.NewNop()))

	inst := &drivers.Instance{
		Environment:   "test",
		OLAPConnector: "duckdb",
		RepoConnector: "repo",
		Connectors: []*runtimev1.Connector{
			{Type: "file", Name: "repo", Config: must(structpb.NewStruct(map[string]any{"dsn": t.TempDir()}))},
			{Type: "duckdb", Name: "duckdb", Config: must(structpb.NewStruct(map[string]any{}))},
		},
		Annotations: map[string]string{"organization_id": "org-1", "project_id": "proj-1"},
	}
	require.NoError(t, rt.CreateInstance(context.Background(), inst))

	res := &runtimev1.Resource{
		Meta: &runtimev1.ResourceMeta{Name: &runtimev1.ResourceName{Kind: ResourceKindAlert, Name: "stock-break"}},
		Resource: &runtimev1.Resource_Alert{Alert: &runtimev1.Alert{State: &runtimev1.AlertState{
			ExecutionHistory: []*runtimev1.AlertExecution{{Result: &runtimev1.AssertionResult{Status: runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL}}},
		}}},
	}

	// No observer is registered: metering must not depend on Act being wired.
	rt.ObserveExecution(context.Background(), inst.ID, res)

	var metered []activity.Event
	for _, e := range sink.Events() {
		if e.EventName == "alert_execution_failed" {
			metered = append(metered, e)
		}
	}
	require.Len(t, metered, 1)
	e := metered[0]
	require.Equal(t, activity.EventTypeMetric, e.EventType)
	require.Equal(t, float64(1), e.Data["value"])
	require.Equal(t, "alert", e.Data["source"])
	require.Equal(t, "stock-break", e.Data["resource_name"])
	require.Equal(t, "org-1", e.Data["organization_id"])
	require.Equal(t, "proj-1", e.Data["project_id"])
	require.Equal(t, inst.ID, e.Data["instance_id"])
}

// newMeterTestRuntime mirrors newTestRuntime but lets the test observe the activity client.
func newMeterTestRuntime(t *testing.T, ac *activity.Client) *Runtime {
	opts := &Options{
		ConnectionCacheSize:     100,
		MetastoreConnector:      "metastore",
		QueryCacheSizeBytes:     int64(datasize.MB) * 100,
		SecurityEngineCacheSize: 100,
		AllowHostAccess:         true,
		SystemConnectors: []*runtimev1.Connector{{
			Type:   "sqlite",
			Name:   "metastore",
			Config: must(structpb.NewStruct(map[string]any{"dsn": fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())})),
		}},
		ControllerLogBufferCapacity:  10000,
		ControllerLogBufferSizeBytes: int64(datasize.MB * 16),
	}
	rt, err := New(context.Background(), opts, zap.NewNop(), storage.MustNew(t.TempDir(), nil), ac, email.New(email.NewNoopSender()))
	require.NoError(t, err)
	t.Cleanup(func() { rt.Close() })
	return rt
}

type captureSink struct {
	mu     sync.Mutex
	events []activity.Event
}

func (s *captureSink) Emit(e activity.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, e)
	return nil
}

func (s *captureSink) Events() []activity.Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]activity.Event(nil), s.events...)
}

func (s *captureSink) Close() {}

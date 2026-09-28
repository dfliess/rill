package runtime

import (
	"context"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// meterExecution emits one metric per persisted alert or report execution, so consumption can be counted per execution
// instead of being inferred from the queries it ran: a query event carries neither the resource nor the outcome, and
// one execution runs a variable number of queries.
//
// The outcome is part of the metric name ("alert_execution_failed", "report_execution_sent", ...) so that a rollup
// keyed by event name keeps it without a schema change. Which outcomes consume credits is decided downstream, in SQL,
// not here (kairos-cloud ADR-0027).
//
// Skipped report and alert checks are not executions: the reconciler returns before creating one, so they never reach
// this hook. Cancelled executions are discarded by the reconciler and never reach it either.
func (r *Runtime) meterExecution(ctx context.Context, instanceID string, self *runtimev1.Resource) {
	name, attrs, ok := executionMetric(self)
	if !ok {
		return
	}

	inst, err := r.Instance(ctx, instanceID)
	if err != nil {
		// Still emit: an execution without its organization is better than a lost one, and instance_id resolves it later.
		r.Logger.Warn("execution meter: instance not found, emitting without instance annotations", zap.String("instance_id", instanceID), zap.Error(err))
		attrs = append(attrs, attribute.String("instance_id", instanceID))
	} else {
		attrs = append(attrs, instanceAnnotationsToAttribs(inst)...)
	}

	r.activity.RecordMetric(ctx, name, 1, attrs...)
}

// executionMetric maps the latest execution of an alert or report to its metric name and attributes. It reports false
// for any other resource kind and for a resource without executions.
func executionMetric(self *runtimev1.Resource) (string, []attribute.KeyValue, bool) {
	if a := self.GetAlert(); a != nil {
		hist := a.GetState().GetExecutionHistory()
		if len(hist) == 0 {
			return "", nil, false
		}
		e := hist[0]
		outcome := "error"
		switch e.GetResult().GetStatus() {
		case runtimev1.AssertionStatus_ASSERTION_STATUS_PASS:
			outcome = "passed"
		case runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL:
			outcome = "failed"
		}
		return "alert_execution_" + outcome, executionAttrs(RequestSourceAlert, self, e.GetAdhoc()), true
	}

	if rep := self.GetReport(); rep != nil {
		hist := rep.GetState().GetExecutionHistory()
		if len(hist) == 0 {
			return "", nil, false
		}
		e := hist[0]
		outcome := "sent"
		if e.GetErrorMessage() != "" {
			outcome = "error"
		}
		// A report with the AI resolver is priced like an Ask, not like a query (kairos-cloud ADR-0027). Keeping the
		// "report_execution_" prefix lets the existing metering filters pass it through unchanged.
		if rep.GetSpec().GetResolver() == "ai" {
			outcome = "ai_" + outcome
		}
		return "report_execution_" + outcome, executionAttrs(RequestSourceReport, self, e.GetAdhoc()), true
	}

	return "", nil, false
}

func executionAttrs(source RequestSource, self *runtimev1.Resource, adhoc bool) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("source", string(source)),
		// Not "name": the activity client already writes the event name under that key.
		attribute.String("resource_name", self.GetMeta().GetName().GetName()),
		attribute.Bool("adhoc", adhoc),
	}
}

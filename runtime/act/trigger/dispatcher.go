package trigger

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
)

// maxScheduleCatchUp bounds how many cron boundaries one EnqueueDueSchedules scan enqueues per trigger. A steady poll
// enqueues zero or one; this only caps a pathological gap (a large backward-then-forward clock jump), so the ticker
// can never spin enqueuing unbounded historical boundaries.
const maxScheduleCatchUp = 512

// Default dispatcher tuning. Shadow mode runs a single worker, so these stay modest.
const (
	defaultBatchSize           = 100
	defaultLeaseSeconds        = 60
	defaultRetryBackoffSeconds = 30
)

// TriggerDef is the dispatcher's view of one validated AgentTrigger: everything needed to match an event and start
// a run, resolved from AgentTriggerState.ValidSpec.
type TriggerDef struct {
	Name         string // the AgentTrigger resource name; the trigger_id in the dispatch ledger
	Agent        string
	SourceKind   string // "alert" | "report" | "schedule"
	SourceName   string
	Events       []string // subscribed event short-names, e.g. "entered_fail"
	Cron         string   // cron expression for a schedule trigger; empty for alert/report
	ActorSubject string
	ActorService bool
	// ActorAttributes are the run-as user attributes the started run's SecurityClaims are built from (nil => the run
	// carries no claims and fails closed). Resolved from AgentTriggerActor.Attributes.
	ActorAttributes map[string]any
	Prompt          string
	Context         map[string]string
	DedupWindow     time.Duration
}

// TriggerCatalog resolves the valid AgentTriggers for an instance. It is an interface so the dispatcher's
// match/dedup logic can be tested with a fake catalog and needs no running runtime.
type TriggerCatalog interface {
	ListTriggers(ctx context.Context, instanceID string) ([]TriggerDef, error)
}

// Dispatcher polls the outbox, matches events against the catalog's valid triggers, deduplicates by window and, in
// shadow mode, starts a run per matched (trigger, event). It is the only component that touches the AgentExecutor;
// the observer and reconcilers never do.
type Dispatcher struct {
	store    *Store
	catalog  TriggerCatalog
	executor act.AgentExecutor
	logger   *slog.Logger

	batchSize           int
	leaseSeconds        int
	retryBackoffSeconds int
	now                 func() time.Time

	// scanMu guards lastScan, the per-instance high-water mark of the schedule ticker: the last time EnqueueDueSchedules
	// scanned that instance. A scan enqueues every cron boundary in (lastScan, now], so as long as the process keeps
	// polling no boundary is missed regardless of tick jitter. It is in-memory: on restart an instance starts fresh from
	// its first scan (no backfill of boundaries that passed while the process was down).
	scanMu   sync.Mutex
	lastScan map[string]time.Time
}

// DispatcherOptions tunes a Dispatcher. Zero values fall back to the defaults above.
type DispatcherOptions struct {
	BatchSize           int
	LeaseSeconds        int
	RetryBackoffSeconds int
	Logger              *slog.Logger
	// Now overrides the clock used for window deduplication. Tests inject it; production leaves it nil (time.Now).
	Now func() time.Time
}

// NewDispatcher builds a Dispatcher over store, resolving triggers from catalog and starting runs on executor.
func NewDispatcher(store *Store, catalog TriggerCatalog, executor act.AgentExecutor, opts DispatcherOptions) *Dispatcher {
	d := &Dispatcher{
		store:               store,
		catalog:             catalog,
		executor:            executor,
		logger:              opts.Logger,
		batchSize:           opts.BatchSize,
		leaseSeconds:        opts.LeaseSeconds,
		retryBackoffSeconds: opts.RetryBackoffSeconds,
		now:                 opts.Now,
		lastScan:            make(map[string]time.Time),
	}
	if d.logger == nil {
		d.logger = slog.Default()
	}
	if d.batchSize <= 0 {
		d.batchSize = defaultBatchSize
	}
	if d.leaseSeconds <= 0 {
		d.leaseSeconds = defaultLeaseSeconds
	}
	if d.retryBackoffSeconds <= 0 {
		d.retryBackoffSeconds = defaultRetryBackoffSeconds
	}
	if d.now == nil {
		d.now = time.Now
	}
	return d
}

// TODO(act, phase 1): known limit #5 — nothing yet drives Dispatch. The trigger path (ConfigureAct wiring,
// Runtime.SetExecutionObserver, NewDispatcher, and a polling loop that calls Dispatch per instance on an interval)
// is not connected to any running executable; it is exercised only by tests. Wiring it into a long-running process
// (the rill-agent-worker) is a deployment-increment task: until then no automatic trigger actually fires in a
// deployed runtime.

// EnqueueDueSchedules is the schedule ticker: for every schedule-kind trigger of an instance it enqueues a synthetic
// tick event for each cron boundary that fell since this instance's last scan, then Dispatch (called right after)
// drains them exactly like alert/report events. It returns how many new tick events it enqueued.
//
// Schedule triggers have no upstream execution to observe, so this is what feeds the outbox for them. It is
// deliberately stateless on disk: the only state is the in-memory lastScan high-water mark. On the first scan of an
// instance it just records "now" and enqueues nothing, so a freshly (re)started process never backfills boundaries
// that passed while it was down; from then on it enqueues every boundary in (lastScan, now], so no boundary is
// missed while the process runs, whatever the tick jitter. Enqueue is idempotent on (instance, event_id) and the
// event_id is stable per (instance, trigger, boundary), so re-enqueuing a boundary — from an overlapping window or a
// second poller — collapses to one event and, downstream, one run.
func (d *Dispatcher) EnqueueDueSchedules(ctx context.Context, instanceID string) (int, error) {
	now := d.now()

	d.scanMu.Lock()
	last, seen := d.lastScan[instanceID]
	if !seen {
		// First scan for this instance: start the clock, don't fire for boundaries before the process came up.
		d.lastScan[instanceID] = now
		d.scanMu.Unlock()
		return 0, nil
	}
	d.scanMu.Unlock()

	triggers, err := d.catalog.ListTriggers(ctx, instanceID)
	if err != nil {
		return 0, err
	}

	enqueued := 0
	for _, tr := range triggers {
		if tr.SourceKind != resourceKindSchedule || tr.Cron == "" {
			continue
		}
		sched, err := cron.ParseStandard(tr.Cron)
		if err != nil {
			// A trigger only reaches the catalog as a ValidSpec once the reconciler has parsed its cron, so this is
			// defensive: skip this one schedule rather than failing the whole scan (and every other trigger with it).
			d.logger.Warn("act/trigger: skipping schedule with unparseable cron",
				slog.String("trigger", tr.Name), slog.String("cron", tr.Cron), slog.String("error", err.Error()))
			continue
		}
		n := 0
		for t := sched.Next(last); !t.After(now); t = sched.Next(t) {
			ok, err := d.store.EnqueueEvent(ctx, newScheduleEvent(instanceID, tr.Name, t, now))
			if err != nil {
				return enqueued, err
			}
			if ok {
				enqueued++
			}
			if n++; n >= maxScheduleCatchUp {
				d.logger.Warn("act/trigger: schedule catch-up capped; skipping older boundaries",
					slog.String("trigger", tr.Name), slog.Int("cap", maxScheduleCatchUp))
				break
			}
		}
	}

	// Advance the high-water mark only after a fully successful scan: an early return above (catalog or enqueue error)
	// leaves lastScan where it was, so the next tick re-covers the same window and drops no boundary.
	d.scanMu.Lock()
	d.lastScan[instanceID] = now
	d.scanMu.Unlock()
	return enqueued, nil
}

// Dispatch processes one batch of pending events for an instance and returns how many new runs it started. It is
// the unit of work a polling loop (or a test) drives; calling it repeatedly drains the outbox. A run is started at
// most once per (trigger, event) however many times Dispatch sees the event, because the dispatch ledger and the
// executor are both idempotent.
func (d *Dispatcher) Dispatch(ctx context.Context, instanceID string) (int, error) {
	events, err := d.store.ClaimPending(ctx, instanceID, d.batchSize, d.leaseSeconds)
	if err != nil {
		return 0, err
	}
	if len(events) == 0 {
		return 0, nil
	}

	triggers, err := d.catalog.ListTriggers(ctx, instanceID)
	if err != nil {
		// Release the leases so a later poll retries once the catalog is available; drop nothing.
		for _, e := range events {
			if merr := d.store.MarkFailed(ctx, e.ID, err.Error(), d.retryBackoffSeconds); merr != nil {
				d.logger.Warn("act/trigger: failed to release lease", slog.String("error", merr.Error()))
			}
		}
		return 0, err
	}

	started := 0
	for _, e := range events {
		n, err := d.dispatchEvent(ctx, e, triggers)
		if err != nil {
			d.logger.Warn("act/trigger: failed to dispatch event",
				slog.String("event_type", e.EventType), slog.String("event_id", e.EventID), slog.String("error", err.Error()))
			if merr := d.store.MarkFailed(ctx, e.ID, err.Error(), d.retryBackoffSeconds); merr != nil {
				d.logger.Warn("act/trigger: failed to mark event failed", slog.String("error", merr.Error()))
			}
			continue
		}
		started += n
		if err := d.store.MarkDone(ctx, e.ID); err != nil {
			d.logger.Warn("act/trigger: failed to mark event done", slog.String("error", err.Error()))
		}
	}
	return started, nil
}

// dispatchEvent matches one event against every trigger and dispatches the matches, returning how many new runs
// were started.
func (d *Dispatcher) dispatchEvent(ctx context.Context, e OutboxEvent, triggers []TriggerDef) (int, error) {
	short := shortEventName(e.EventType)
	started := 0
	for _, tr := range triggers {
		if !triggerMatches(tr, e, short) {
			continue
		}
		ok, err := d.dispatchOne(ctx, e, tr)
		if err != nil {
			return started, err
		}
		if ok {
			started++
		}
	}
	return started, nil
}

// dispatchOne dispatches one matched (event, trigger): window-deduplicate, reserve the run's identity from the
// stable (instance, trigger, event), then start the run and record it. It returns whether it started a new run.
func (d *Dispatcher) dispatchOne(ctx context.Context, e OutboxEvent, tr TriggerDef) (bool, error) {
	windowSeconds := int(tr.DedupWindow / time.Second)

	// Window deduplication: within the window, a further event collapses onto the existing run (§9.2). This is a
	// best-effort business rule; the reservation's (trigger, event_id) uniqueness below is the hard guarantee.
	if tr.DedupWindow > 0 {
		recent, err := d.store.HasRecentDispatch(ctx, e.InstanceID, tr.Name, d.now().Add(-tr.DedupWindow))
		if err != nil {
			return false, err
		}
		if recent {
			if _, err := d.store.RecordDispatch(ctx, e.InstanceID, tr.Name, e.EventID, e.EventType, tr.Agent, "", outcomeDeduplicated, windowSeconds); err != nil {
				return false, err
			}
			return false, nil
		}
	}

	// Reserve the dispatch slot BEFORE starting the run, freezing the agent to run. The run's identity then derives
	// only from the stable (instance, trigger, event) plus this frozen agent, so a retry after a crash between here
	// and MarkDispatched — even one where the trigger has since been re-pointed to a different agent — reuses the
	// same run rather than starting a second one for the same event (§9.6, §12).
	res, err := d.store.ReserveDispatch(ctx, e.InstanceID, tr.Name, e.EventID, e.EventType, tr.Agent, windowSeconds)
	if err != nil {
		return false, err
	}
	if res.Outcome == outcomeDispatched {
		// A prior delivery already started and recorded the run: OAOO, nothing to do.
		return false, nil
	}

	// Link the run back to the source resource for alert- and report-driven triggers, so the run detail can resolve
	// the alert (or report) that fired it. Use the firing resource (e.ResourceName), not the trigger's configured
	// SourceName: for a named trigger they are identical, but for a wildcard trigger (empty SourceName, ADR-0016) the
	// configured name is blank while e.ResourceName is the exact alert that fired — which is what the run must record.
	// A schedule trigger has no such resource to point at, so it stays empty.
	var triggerRef string
	if tr.SourceKind == act.TriggerAlert || tr.SourceKind == act.TriggerReport {
		triggerRef = e.ResourceName
	}

	// The run's governed identity: build SecurityClaims from the trigger's run-as attributes (mirrors how an alert
	// runs its query as query_for_attributes). Nil when the trigger declares none, in which case the run carries no
	// claims and the session factory fails closed — an automatic run never falls back to a privileged session.
	var claims *runtime.SecurityClaims
	if len(tr.ActorAttributes) > 0 {
		claims = &runtime.SecurityClaims{UserAttributes: tr.ActorAttributes}
	}

	// Start (or, on a retry, attach to) the run using the FROZEN agent, not tr.Agent, so the composed run id is
	// identical across retries. The executor is idempotent on that id, so a replay attaches to the same run.
	key := ComposeIdempotencyKey(tr.Name, e.EventID)
	runID, err := d.executor.Start(ctx, act.AgentRunInput{
		InstanceID:     e.InstanceID,
		AgentName:      res.AgentName,
		Prompt:         renderTriggerInput(tr),
		IdempotencyKey: key,
		// Truthful provenance: the run records its real source kind ("alert"/"report"/"schedule"), so an automatic
		// run never appears as manual in the audit trail (§18.3), and trigger_ref names the exact source resource.
		Trigger:    tr.SourceKind,
		TriggerRef: triggerRef,
		Actor:      act.Actor{Subject: tr.ActorSubject, ServicePrincipal: tr.ActorService, Claims: claims},
	})
	if err != nil {
		return false, fmt.Errorf("act/trigger: start run for trigger %q: %w", tr.Name, err)
	}

	// Record the run id on the reserved slot (reserved -> dispatched; idempotent). Report a newly started run only
	// when this call did the reservation: a retry completing a prior reservation attaches to the existing run.
	if err := d.store.MarkDispatched(ctx, e.InstanceID, tr.Name, e.EventID, runID); err != nil {
		return false, err
	}
	return res.New, nil
}

// triggerMatches reports whether a trigger subscribes to this event: same source kind, the source name matches, and
// the event's short name is in the trigger's allowlist. An empty SourceName is a wildcard (ADR-0016): the trigger
// fires for any resource of that kind (e.g. any alert that enters fail); a non-empty name pins it to that one source.
func triggerMatches(tr TriggerDef, e OutboxEvent, shortEvent string) bool {
	// A schedule trigger is its own source: the ticker emits a tick tagged with the trigger's own name, and the trigger
	// matches only its own ticks. It has no events allowlist (the reconciler forbids one), so the name/event matching
	// used for alert/report does not apply — matching by identity is what keeps two schedules with the same cron apart.
	if tr.SourceKind == resourceKindSchedule {
		return e.ResourceKind == resourceKindSchedule && e.ResourceName == tr.Name
	}
	return tr.SourceKind == e.ResourceKind &&
		(tr.SourceName == "" || tr.SourceName == e.ResourceName) &&
		slices.Contains(tr.Events, shortEvent)
}

// shortEventName strips the "<kind>." namespace from a fully-qualified event type: "alert.entered_fail" ->
// "entered_fail". A trigger subscribes to the short form.
func shortEventName(eventType string) string {
	if i := strings.IndexByte(eventType, '.'); i >= 0 {
		return eventType[i+1:]
	}
	return eventType
}

// renderTriggerInput builds the run's prompt from the trigger definition: the configured prompt followed by the
// trigger's static context (§9.2). The context was resolved into TriggerDef but previously discarded; rendering it
// here is what actually gets it in front of the agent. Keys are emitted in sorted order so the same trigger and
// context always produce the same prompt (reproducible runs).
//
// input.prompt is optional (ADR-0016): an empty prompt is not an error, it just specializes nothing. With no prompt
// and no context this returns "" and the run leans entirely on the agent's own instructions; with context only, it
// returns just the context block. The executor and runner tolerate an empty prompt (a fresh segment opens with the
// system prompt built from the agent's instructions plus an empty user turn).
func renderTriggerInput(tr TriggerDef) string {
	if len(tr.Context) == 0 {
		return tr.Prompt
	}
	keys := make([]string, 0, len(tr.Context))
	for k := range tr.Context {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	var b strings.Builder
	b.WriteString(tr.Prompt)
	if tr.Prompt != "" {
		b.WriteString("\n\n")
	}
	b.WriteString("Context:")
	for _, k := range keys {
		fmt.Fprintf(&b, "\n- %s: %s", k, tr.Context[k])
	}
	return b.String()
}

// ComposeIdempotencyKey builds the executor idempotency key from (trigger, event). It length-prefixes the trigger
// name so two triggers that start the same agent for the same event get distinct runs: the executor namespaces the
// run ID by (instance, agent) but takes this key verbatim, so the key alone must separate triggers.
func ComposeIdempotencyKey(triggerName, eventID string) string {
	return fmt.Sprintf("%d:%s/%s", len(triggerName), triggerName, eventID)
}

// RuntimeTriggerCatalog is the production TriggerCatalog: it lists the valid AgentTrigger resources from the
// runtime's per-instance controller. A trigger with no valid spec (unreconciled or invalid) is skipped, so an
// invalid trigger never fires.
type RuntimeTriggerCatalog struct {
	rt *runtime.Runtime
}

var _ TriggerCatalog = (*RuntimeTriggerCatalog)(nil)

// NewRuntimeTriggerCatalog returns a catalog backed by rt.
func NewRuntimeTriggerCatalog(rt *runtime.Runtime) *RuntimeTriggerCatalog {
	return &RuntimeTriggerCatalog{rt: rt}
}

func (c *RuntimeTriggerCatalog) ListTriggers(ctx context.Context, instanceID string) ([]TriggerDef, error) {
	ctrl, err := c.rt.Controller(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	resources, err := ctrl.List(ctx, runtime.ResourceKindAgentTrigger, "", false)
	if err != nil {
		return nil, err
	}
	defs := make([]TriggerDef, 0, len(resources))
	for _, res := range resources {
		t := res.GetAgentTrigger()
		if t == nil {
			continue
		}
		spec := t.GetState().GetValidSpec()
		if spec == nil {
			continue // only valid triggers dispatch; an invalid one looks absent
		}
		defs = append(defs, triggerDefFromSpec(res.Meta.Name.Name, spec))
	}
	return defs, nil
}

// triggerDefFromSpec maps a validated AgentTriggerSpec to the dispatcher's TriggerDef.
func triggerDefFromSpec(name string, spec *runtimev1.AgentTriggerSpec) TriggerDef {
	def := TriggerDef{
		Name:  name,
		Agent: spec.Agent,
	}
	if src := spec.Source; src != nil {
		def.SourceKind = sourceKindString(src.Kind)
		def.SourceName = src.Name
		def.Events = src.Events
		def.Cron = src.Cron
	}
	if actor := spec.Actor; actor != nil {
		// The audit subject is the named user (mirrors an alert's for.user_*). ServicePrincipal marks a run that acts
		// as an unnamed service/automation identity — explicit attributes, or none — which the UI labels "service
		// principal" rather than showing a user subject.
		switch {
		case actor.UserEmail != "":
			def.ActorSubject = actor.UserEmail
		case actor.UserId != "":
			def.ActorSubject = actor.UserId
		}
		def.ActorService = actor.UserEmail == "" && actor.UserId == ""
		if actor.Attributes != nil {
			def.ActorAttributes = actor.Attributes.AsMap()
		}
	}
	if in := spec.Input; in != nil {
		def.Prompt = in.Prompt
		def.Context = in.Context
	}
	if dd := spec.Deduplication; dd != nil {
		def.DedupWindow = time.Duration(dd.WindowSeconds) * time.Second
	}
	return def
}

// sourceKindString maps a source kind enum to the short string used in events and matching.
func sourceKindString(k runtimev1.AgentTriggerSourceKind) string {
	switch k {
	case runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_ALERT:
		return resourceKindAlert
	case runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_REPORT:
		return resourceKindReport
	case runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_SCHEDULE:
		return resourceKindSchedule
	default:
		return ""
	}
}

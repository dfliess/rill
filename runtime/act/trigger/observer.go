package trigger

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/rilldata/rill/runtime"
)

const (
	// observerBufferSize bounds the in-memory queue between the reconciler goroutine and the background publisher.
	// It is a shadow-mode tunable: large enough to absorb a burst of executions, small enough to bound memory.
	observerBufferSize = 1024
	// observerPublishTimeout bounds one EnqueueEvent so a slow or stuck Postgres cannot wedge the publisher.
	observerPublishTimeout = 10 * time.Second
)

// Observer implements runtime.ExecutionObserver. It derives transition events from a persisted alert/report
// execution and publishes them to the outbox. It is deliberately thin and, crucially, non-blocking: OnExecution
// only derives events (pure, in-memory) and hands them to a bounded in-memory buffer; a background goroutine drains
// that buffer to Postgres on its own context. So a slow or unavailable Postgres can never slow down, block, or fail
// alert/report reconciliation — the observer bounds its own latency, as the seam requires (§9.6).
type Observer struct {
	store  *Store
	logger *slog.Logger

	events chan Event
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// mu guards the closed flag and the send on events. OnExecution takes it for reading (concurrent, and held only
	// for a non-blocking channel send, so it never stalls reconciliation); Close takes it for writing, which waits
	// for any in-flight send to finish and then blocks new ones. This makes shutdown race-free: once Close sets
	// closed, no send can strand an event in a channel the drained publisher will never read.
	mu     sync.RWMutex
	closed bool
}

var _ runtime.ExecutionObserver = (*Observer)(nil)

// NewObserver returns an Observer that publishes to store from a background goroutine. The goroutine runs until
// Close, so a caller that constructs an Observer owns its lifetime and must Close it.
func NewObserver(store *Store, logger *slog.Logger) *Observer {
	if logger == nil {
		logger = slog.Default()
	}
	ctx, cancel := context.WithCancel(context.Background())
	o := &Observer{
		store:  store,
		logger: logger,
		events: make(chan Event, observerBufferSize),
		cancel: cancel,
	}
	o.wg.Add(1)
	go o.publishLoop(ctx)
	return o
}

// TODO(act, phase 1): known limit #2 — there is no transactional outbox or republisher next to the commit of the
// AlertExecution/ReportExecution. The reconciler persists the execution, then (separately) ObserveExecution reaches
// this observer, which now buffers the event in memory before Postgres. A crash between saving the execution and
// the buffered enqueue landing loses the event entirely. The fix is to write the trigger event in the SAME
// transaction that persists the execution (or a republisher that re-derives events from execution history the
// dispatcher has not yet seen), so delivery is at-least-once from the commit, not best-effort from an in-memory
// buffer.

// OnExecution derives the events for one persisted execution and hands each to the background publisher. It runs in
// the reconciler's goroutine, so it does only bounded, non-blocking work: derive (pure, in-memory) and a
// non-blocking channel send. A full buffer drops the event with a log rather than blocking reconciliation; loss is
// acceptable in shadow mode, and §9.6's source re-publish is the backstop.
func (o *Observer) OnExecution(_ context.Context, ev runtime.ExecutionEvent) {
	for _, e := range DeriveEvents(ev) {
		// Decide under the read lock (closed check + non-blocking send), but log outside it: the logger's handler is
		// not ours and could block, and holding the lock across it would let a slow handler stall Close.
		o.mu.RLock()
		var dropped string
		switch {
		case o.closed:
			dropped = "observer closed"
		default:
			select {
			case o.events <- e:
			default:
				dropped = "buffer full"
			}
		}
		o.mu.RUnlock()

		if dropped != "" {
			o.logger.Warn("act/trigger: dropping trigger event: "+dropped,
				slog.String("event_type", e.EventType),
				slog.String("resource", e.ResourceName),
				slog.String("event_id", e.EventID))
		}
	}
}

// publishLoop drains the buffer to Postgres until Close cancels it, then drains what remains so a clean shutdown
// publishes rather than discards already-buffered events.
func (o *Observer) publishLoop(ctx context.Context) {
	defer o.wg.Done()
	for {
		select {
		case <-ctx.Done():
			for {
				select {
				case e := <-o.events:
					o.publish(e)
				default:
					return
				}
			}
		case e := <-o.events:
			o.publish(e)
		}
	}
}

// publish enqueues one event on its own bounded context, independent of the reconciler's, so publishing latency is
// never coupled to reconciliation. A failure or a slow Postgres is logged and dropped; it never propagates.
func (o *Observer) publish(e Event) {
	ctx, cancel := context.WithTimeout(context.Background(), observerPublishTimeout)
	defer cancel()

	inserted, err := o.store.EnqueueEvent(ctx, e)
	if err != nil {
		o.logger.Warn("act/trigger: failed to enqueue trigger event",
			slog.String("event_type", e.EventType),
			slog.String("resource", e.ResourceName),
			slog.String("error", err.Error()))
		return
	}
	if inserted {
		o.logger.Debug("act/trigger: enqueued trigger event",
			slog.String("event_type", e.EventType),
			slog.String("resource", e.ResourceName),
			slog.String("event_id", e.EventID))
	}
}

// Close stops the background publisher and waits for it to finish draining the buffer. Setting closed under the
// write lock first guarantees no OnExecution is mid-send and none can start, so once the publisher drains and exits
// no event can be stranded in the channel. A final drain publishes anything the publisher's own drain and the send
// gate could have left buffered, so a clean shutdown never silently loses an already-accepted event.
func (o *Observer) Close() {
	o.mu.Lock()
	o.closed = true
	o.mu.Unlock()

	o.cancel()
	o.wg.Wait()

	for {
		select {
		case e := <-o.events:
			o.publish(e)
		default:
			return
		}
	}
}

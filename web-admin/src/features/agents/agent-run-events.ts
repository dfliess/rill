import { StreamAgentRunEventsRequest } from "@rilldata/web-common/proto/gen/rill/runtime/v1/agents_pb";
import type { RuntimeClient } from "@rilldata/web-common/runtime-client/v2";
import { readable, type Readable } from "svelte/store";
import type { AgentRunTimelineEvent } from "./types";

export interface AgentRunEventsState {
  events: AgentRunTimelineEvent[];
  /** True while the stream is open and delivering events. */
  streaming: boolean;
  /** Set when the stream could not be consumed; the run detail still works off status polling. */
  error: string | null;
}

const RECONNECT_DELAY_MS = 3_000;
const MAX_STREAM_ATTEMPTS = 3;
/** A stream that stayed open this long was working; its end is a cutoff to resume from, not a failure to count. */
const HEALTHY_STREAM_MS = 20_000;

/**
 * Live event feed for a run, backed by AgentService.StreamAgentRunEvents.
 *
 * The stream is consumed as a ConnectRPC server-streaming call over the same
 * transport as every other request. Each event is normalized to a plain
 * `AgentRunTimelineEvent` so the timeline component never touches proto types.
 * We resume with `after_id = lastId` on reconnect so a dropped connection does
 * not replay or skip events.
 *
 * TODO(#123): The backend serves these events over SSE via a dedicated HTTP
 * handler (see agents.proto). It is not yet verified in this environment that
 * the ConnectRPC server-streaming path is registered alongside that raw SSE
 * endpoint. If it is not, this store surfaces `error` and the run detail
 * degrades gracefully to status-only (GetAgentRun polling) — no events shown.
 * Confirm the streaming handler end-to-end before relying on the live timeline.
 */
export function agentRunEventsStore(
  client: RuntimeClient,
  runId: string,
): Readable<AgentRunEventsState> {
  return readable<AgentRunEventsState>(
    { events: [], streaming: false, error: null },
    (set) => {
      const events: AgentRunTimelineEvent[] = [];
      const controller = new AbortController();
      let stopped = false;
      let lastId = 0n;

      async function consume(): Promise<void> {
        for (let attempt = 0; attempt < MAX_STREAM_ATTEMPTS; attempt++) {
          const openedAt = Date.now();
          try {
            set({ events: [...events], streaming: true, error: null });
            const stream = client.agentService.streamAgentRunEvents(
              new StreamAgentRunEventsRequest({
                instanceId: client.instanceId,
                runId,
                afterId: lastId,
              }),
              { signal: controller.signal },
            );
            for await (const res of stream) {
              if (stopped) return;
              const event = res.event;
              if (!event) continue;
              events.push({
                id: event.id.toString(),
                seq: Number(event.seq),
                eventType: event.eventType,
                status: event.status,
                visibility: event.visibility,
                payload: event.payload ? event.payload.toJson() : undefined,
                createdOn: event.createdOn
                  ? event.createdOn.toDate().toISOString()
                  : undefined,
              });
              lastId = event.id;
              set({ events: [...events], streaming: true, error: null });
            }
            // Stream ended cleanly (the backend closed it, e.g. run finished).
            set({ events: [...events], streaming: false, error: null });
            return;
          } catch (e) {
            if (stopped || controller.signal.aborted) return;
            // The attempt budget is there to stop hammering a backend that is
            // failing, not to put a lifespan on a healthy stream. A connection
            // that stayed up is evidence the backend is fine, so it does not
            // spend from the budget: without this, a run parked on an approval
            // for an afternoon exhausts three server-side cutoffs and the
            // timeline goes dark while the run is still perfectly watchable.
            if (Date.now() - openedAt >= HEALTHY_STREAM_MS) {
              attempt = -1;
              continue;
            }
            const last = attempt === MAX_STREAM_ATTEMPTS - 1;
            if (last) {
              set({
                events: [...events],
                streaming: false,
                error: e instanceof Error ? e.message : String(e),
              });
              return;
            }
            await new Promise((r) => setTimeout(r, RECONNECT_DELAY_MS));
          }
        }
      }

      void consume();

      return () => {
        stopped = true;
        controller.abort();
      };
    },
  );
}

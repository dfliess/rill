import type { PartialMessage } from "@bufbuild/protobuf";
import type {
  V1AnalystAgentContext,
  V1Expression,
} from "@rilldata/web-common/runtime-client";
import type { AnalystAgentContext } from "@rilldata/web-common/proto/gen/rill/runtime/v1/api_pb";
import type { RuntimeClient } from "@rilldata/web-common/runtime-client/v2";
import type { Interval } from "luxon";
import {
  agentServiceGetAgentRun,
  agentServiceStartAgentRun,
} from "@rilldata/web-common/runtime-client/v2/gen/agent-service";
import { runtimeServiceGetConversation } from "@rilldata/web-common/runtime-client/v2/gen/runtime-service";

const POLL_INTERVAL_MS = 1_500;
const POLL_TIMEOUT_MS = 120_000;

const TERMINAL_OK = new Set(["succeeded", "completed"]);
const TERMINAL_ERROR = new Set([
  "failed",
  "errored",
  "error",
  "canceled",
  "cancelled",
]);
// A narrator agent declares no MCP block, so it can never propose a write and should never park on an
// approval. Reaching one of these means the configured agent is not read-only: stop rather than poll
// until the timeout, and say so, because the run is waiting on a human who is not watching this canvas.
const AWAITING_APPROVAL = new Set([
  "waiting_approval",
  "pending_approval",
  "suspended",
]);

export class AgentNoteError extends Error {}

/**
 * Runs the narrator agent and returns its final answer.
 *
 * The run is started with a caller-supplied idempotency key: re-entering the dashboard with the same
 * context attaches to the existing run instead of paying for a second completion, so the note is
 * generated once per (agent, prompt, context) rather than once per page load.
 */
export async function runAgentNote(
  client: RuntimeClient,
  params: {
    agent: string;
    prompt: string;
    idempotencyKey: string;
    dashboardContext?: V1AnalystAgentContext;
  },
  signal?: AbortSignal,
): Promise<string> {
  const started = await agentServiceStartAgentRun(
    client,
    {
      name: params.agent,
      prompt: params.prompt,
      idempotencyKey: params.idempotencyKey,
      trigger: "canvas",
      // The server folds this into the run's prompt, so the note is scoped to the filters and time range the
      // reader is actually looking at. Same object the AI-Chat sends from a canvas.
      //
      // The cast is a type-level formality, not a shape change: the v2 client sends this request through
      // `StartAgentRunRequest.fromJson`, so what it needs on the wire is proto-JSON, which is exactly what
      // the REST client's `V1AnalystAgentContext` already is. The parameter is typed as the decoded message
      // instead, where a filter is a `oneof` (`{ case, value }`) rather than the flat JSON form, so the two
      // types disagree on paper while describing the same bytes.
      dashboardContext: params.dashboardContext as
        | PartialMessage<AnalystAgentContext>
        | undefined,
    },
    { signal },
  );

  const runId = started.runId;
  if (!runId) {
    throw new AgentNoteError("The agent run did not return a run id.");
  }

  const run = await pollUntilSettled(client, runId, signal);

  const status = run.status ?? "";
  if (AWAITING_APPROVAL.has(status)) {
    throw new AgentNoteError(
      `Agent "${params.agent}" paused for approval. A note agent must be read-only: remove its action tools.`,
    );
  }
  if (TERMINAL_ERROR.has(status)) {
    throw new AgentNoteError(run.error || `The agent run ${status}.`);
  }

  const conversationId = run.conversationId;
  if (!conversationId) {
    throw new AgentNoteError("The agent run produced no conversation.");
  }

  const conversation = await runtimeServiceGetConversation(
    client,
    { conversationId },
    { signal },
  );

  const answer = latestAssistantText(conversation?.conversation?.messages);
  if (!answer) {
    throw new AgentNoteError("The agent returned no text.");
  }
  return answer;
}

async function pollUntilSettled(
  client: RuntimeClient,
  runId: string,
  signal?: AbortSignal,
) {
  const deadline = Date.now() + POLL_TIMEOUT_MS;

  for (;;) {
    const res = await agentServiceGetAgentRun(client, { runId }, { signal });
    const run = res?.run;
    const status = run?.status ?? "";

    if (
      run &&
      (TERMINAL_OK.has(status) ||
        TERMINAL_ERROR.has(status) ||
        AWAITING_APPROVAL.has(status))
    ) {
      return run;
    }

    if (Date.now() >= deadline) {
      throw new AgentNoteError(
        "The agent did not finish in time. Try regenerating the note.",
      );
    }

    await sleep(POLL_INTERVAL_MS, signal);
  }
}

/**
 * Returns the text of the last assistant message, which for a read-only run is the agent's final answer.
 * Prefers `contentData` and falls back to the convenience content blocks, since a message carries its
 * text in either shape depending on how the model emitted it.
 */
function latestAssistantText(
  messages: { role?: string; contentData?: string; content?: unknown[] }[] = [],
): string {
  for (let i = messages.length - 1; i >= 0; i--) {
    const msg = messages[i];
    if (msg.role !== "assistant") continue;

    const direct = msg.contentData?.trim();
    if (direct) return direct;

    const blocks = (msg.content ?? [])
      .map((block) => (block as { text?: string })?.text ?? "")
      .join("")
      .trim();
    if (blocks) return blocks;
  }
  return "";
}

function sleep(ms: number, signal?: AbortSignal): Promise<void> {
  return new Promise((resolve, reject) => {
    if (signal?.aborted) {
      reject(signal.reason);
      return;
    }
    const timer = setTimeout(() => {
      signal?.removeEventListener("abort", onAbort);
      resolve();
    }, ms);
    function onAbort() {
      clearTimeout(timer);
      reject(signal?.reason);
    }
    signal?.addEventListener("abort", onAbort, { once: true });
  });
}

/**
 * Reports whether the dashboard context the note will be generated against has settled.
 *
 * The canvas window does not exist until two round trips have landed: the canvas spec, which is what
 * fills `hasTimeSeriesMap` (and until then the range store returns undefined — stores/time-state.ts:91),
 * and then deriveInterval's own query that turns a range like "P12M" into dates. A note that does not
 * wait composes its key WITHOUT the window, runs, and runs a second time under the key the window
 * produces: two completions per load, and the answer left on screen is the one computed without the
 * reader's time range (issue #155).
 *
 * Both halves are needed, one per round trip. The time half alone would let the first render through,
 * because before the spec arrives the map is EMPTY, so "no time series" and "not known yet" look
 * identical. The shape mirrors the markdown component (markdown/util.ts:167-172), which narrates the
 * same dashboards and has the same "with or without time" requirement.
 */
export function noteContextReady(input: {
  /** The canvas spec, once fetched. Undefined while the query is in flight. */
  canvasData?: { metricsViews?: Record<string, unknown> };
  /** One entry per metrics view of the canvas, filled when the spec is PROCESSED, not merely fetched. */
  hasTimeSeriesMap?: Map<string, boolean>;
  /** The resolved window, or undefined while it resolves. */
  interval?: Interval<true>;
}): boolean {
  const { canvasData, hasTimeSeriesMap, interval } = input;
  if (!canvasData || !hasTimeSeriesMap) return false;

  // processSpec is what fills the map, and it skips the spec store's first emission
  // (canvas-entity.ts:565), so `canvasData` alone can precede it. A canvas that declares no metrics view
  // never fills the map, and must not wait for something that will never arrive.
  const declaresMetricsViews =
    Object.keys(canvasData.metricsViews ?? {}).length > 0;
  if (declaresMetricsViews && hasTimeSeriesMap.size === 0) return false;

  // The note names no single metrics view, so it asks what the global time state asks: does ANY of them
  // carry a time dimension (stores/time-state.ts:88).
  const hasTimeSeries = [...hasTimeSeriesMap.values()].some(Boolean);
  return !hasTimeSeries || !!interval;
}

/**
 * Builds the context the note is generated against: the canvas the reader is on, the window they are
 * looking at, and the filters in force per metrics view.
 *
 * Filters with no expressions are dropped rather than sent as empty conditions, so an untouched dashboard
 * and one whose filters were set and cleared again produce the same context, and therefore the same run.
 */
export function buildDashboardContext(input: {
  canvas?: string;
  interval?: Interval<true>;
  filterMap?: Map<string, V1Expression>;
}): V1AnalystAgentContext | undefined {
  const ctx: V1AnalystAgentContext = {};

  if (input.canvas) ctx.canvas = input.canvas;

  if (input.interval?.isValid) {
    ctx.timeStart = input.interval.start.toUTC().toISO() ?? undefined;
    ctx.timeEnd = input.interval.end.toUTC().toISO() ?? undefined;
  }

  const where: NonNullable<V1AnalystAgentContext["wherePerMetricsView"]> = {};
  input.filterMap?.forEach((expr, metricsView) => {
    if (expr?.cond?.exprs?.length) where[metricsView] = expr;
  });
  if (Object.keys(where).length) ctx.wherePerMetricsView = where;

  return Object.keys(ctx).length ? ctx : undefined;
}

/**
 * Serialises the dashboard context into a stable string for the idempotency key.
 *
 * Object key order is not a safe basis for this: the filter map is built by iterating a store, so the entries
 * are sorted here the same way the server sorts them when rendering the prompt. Two visits to an unchanged
 * dashboard must produce the same string, or every page load would pay for a fresh completion.
 */
export function dashboardContextKey(ctx?: V1AnalystAgentContext): string {
  if (!ctx) return "";

  const filters = Object.entries(ctx.wherePerMetricsView ?? {})
    .sort(([a], [b]) => (a < b ? -1 : a > b ? 1 : 0))
    .map(([mv, expr]) => `${mv}=${JSON.stringify(expr)}`)
    .join(";");

  return [
    ctx.canvas ?? "",
    ctx.explore ?? "",
    ctx.timeStart ?? "",
    ctx.timeEnd ?? "",
    (ctx.dimensions ?? []).join(","),
    (ctx.measures ?? []).join(","),
    filters,
  ].join("|");
}

/**
 * Derives the idempotency key for a note. The key folds in everything that would change the answer, so
 * a filter or time-range change produces a new run while a plain revisit reuses the existing one.
 */
export function agentNoteIdempotencyKey(parts: {
  agent: string;
  prompt: string;
  context: string;
  nonce?: number;
}): string {
  const raw = [parts.agent, parts.prompt, parts.context, parts.nonce ?? 0].join(
    "\0",
  );

  // FNV-1a: a short stable digest is enough here; the key only needs to be collision-resistant across
  // one project's canvases, not cryptographically secure.
  let hash = 0x811c9dc5;
  for (let i = 0; i < raw.length; i++) {
    hash ^= raw.charCodeAt(i);
    hash = Math.imul(hash, 0x01000193);
  }
  return `agent-note-${(hash >>> 0).toString(16)}`;
}

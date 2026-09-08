import type {
  V1AnalystAgentContext,
  V1Expression,
} from "@rilldata/web-common/runtime-client";
import { DateTime, Interval } from "luxon";
import { describe, expect, it } from "vitest";
import {
  agentNoteIdempotencyKey,
  buildDashboardContext,
  dashboardContextKey,
  noteContextReady,
} from "./util";

function ctxWithFilters(order: string[]): V1AnalystAgentContext {
  const wherePerMetricsView: NonNullable<
    V1AnalystAgentContext["wherePerMetricsView"]
  > = {};
  for (const mv of order) {
    wherePerMetricsView[mv] = {
      cond: { op: "OPERATION_AND", exprs: [{ ident: `${mv}_dim` }] },
    };
  }
  return { canvas: "direccion_canvas", wherePerMetricsView };
}

function monthInterval(month: number): Interval<true> {
  return Interval.fromDateTimes(
    DateTime.fromObject({ year: 2026, month, day: 1 }, { zone: "utc" }),
    DateTime.fromObject(
      { year: 2026, month: month + 1, day: 1 },
      {
        zone: "utc",
      },
    ),
  ) as Interval<true>;
}

describe("buildDashboardContext", () => {
  it("is undefined when there is nothing to say", () => {
    expect(buildDashboardContext({})).toBeUndefined();
  });

  it("carries the window the reader is looking at, in UTC", () => {
    const ctx = buildDashboardContext({
      canvas: "home_ejecutivo",
      interval: monthInterval(7),
    });

    expect(ctx?.canvas).toBe("home_ejecutivo");
    expect(ctx?.timeStart).toBe("2026-07-01T00:00:00.000Z");
    expect(ctx?.timeEnd).toBe("2026-08-01T00:00:00.000Z");
  });

  it("distinguishes a note pinned to its own window from the canvas one", () => {
    // What a component's `time_filters` buys: the same canvas and filters, a different period, and
    // therefore a different run rather than the canvas-wide note served back.
    const canvasWide = buildDashboardContext({
      canvas: "home_ejecutivo",
      interval: monthInterval(7),
    });
    const pinned = buildDashboardContext({
      canvas: "home_ejecutivo",
      interval: monthInterval(6),
    });

    expect(dashboardContextKey(pinned)).not.toBe(
      dashboardContextKey(canvasWide),
    );
  });

  it("drops filters that carry no expressions", () => {
    // A dashboard whose filters were set and then cleared holds an empty condition. Sending it would read
    // as a different context and pay for a second completion that answers the same question.
    const filterMap = new Map<string, V1Expression>([
      ["sales_metrics", { cond: { op: "OPERATION_AND", exprs: [] } }],
    ]);

    expect(
      buildDashboardContext({ canvas: "home_ejecutivo", filterMap })
        ?.wherePerMetricsView,
    ).toBeUndefined();
  });
});

describe("dashboardContextKey", () => {
  it("is empty without context", () => {
    expect(dashboardContextKey(undefined)).toBe("");
  });

  it("does not depend on the insertion order of the filter map", () => {
    // The filter map is built by iterating a store, so the same dashboard can yield different key orders.
    // If that leaked into the key, an unchanged dashboard would start a fresh (paid) run on every visit.
    const a = ctxWithFilters(["academico_metrics", "finanzas_metrics"]);
    const b = ctxWithFilters(["finanzas_metrics", "academico_metrics"]);
    expect(dashboardContextKey(a)).toBe(dashboardContextKey(b));
  });

  it("changes when the state the reader sees changes", () => {
    const base = dashboardContextKey({ canvas: "direccion_canvas" });
    const filtered = dashboardContextKey(ctxWithFilters(["academico_metrics"]));
    const timed = dashboardContextKey({
      canvas: "direccion_canvas",
      timeStart: "2026-06-01T00:00:00Z",
      timeEnd: "2026-07-01T00:00:00Z",
    });

    expect(filtered).not.toBe(base);
    expect(timed).not.toBe(base);
  });
});

describe("agentNoteIdempotencyKey", () => {
  it("is stable for the same inputs and moves with the nonce", () => {
    const parts = { agent: "narrador", prompt: "narra", context: "ctx" };
    expect(agentNoteIdempotencyKey(parts)).toBe(agentNoteIdempotencyKey(parts));
    expect(agentNoteIdempotencyKey({ ...parts, nonce: 1 })).not.toBe(
      agentNoteIdempotencyKey(parts),
    );
  });

  it("separates runs whose context differs", () => {
    expect(
      agentNoteIdempotencyKey({
        agent: "narrador",
        prompt: "narra",
        context: "campus=norte",
      }),
    ).not.toBe(
      agentNoteIdempotencyKey({
        agent: "narrador",
        prompt: "narra",
        context: "campus=sur",
      }),
    );
  });
});

describe("noteContextReady", () => {
  const window = Interval.fromDateTimes(
    DateTime.utc(2026, 8, 1),
    DateTime.utc(2026, 9, 1),
  ) as Interval<true>;
  const spec = { metricsViews: { sales: {} } };

  // The bug of issue #155: without waiting, the note runs once with no window and again with it, paying
  // for two completions and leaving the windowless answer on screen.
  it("waits while the canvas window is still resolving", () => {
    expect(
      noteContextReady({
        canvasData: spec,
        hasTimeSeriesMap: new Map([["sales", true]]),
        interval: undefined,
      }),
    ).toBe(false);

    expect(
      noteContextReady({
        canvasData: spec,
        hasTimeSeriesMap: new Map([["sales", true]]),
        interval: window,
      }),
    ).toBe(true);
  });

  // The half that is easy to miss: before the spec is processed the map is EMPTY, so asking only "does
  // any metrics view have a time dimension" reads false and lets the premature run through — which is
  // exactly the state the first of the two runs was composed in.
  it("does not mistake an unprocessed spec for a canvas without time", () => {
    expect(
      noteContextReady({
        canvasData: spec,
        hasTimeSeriesMap: new Map(),
        interval: undefined,
      }),
    ).toBe(false);
  });

  // And the mirror image: a canvas whose metrics views carry no time dimension never gets a window, so
  // waiting for one would leave its note unwritten forever.
  it("does not wait for a window a canvas without time will never have", () => {
    expect(
      noteContextReady({
        canvasData: spec,
        hasTimeSeriesMap: new Map([["sales", false]]),
        interval: undefined,
      }),
    ).toBe(true);
  });

  it("does not wait on a canvas that declares no metrics view at all", () => {
    expect(
      noteContextReady({
        canvasData: { metricsViews: {} },
        hasTimeSeriesMap: new Map(),
        interval: undefined,
      }),
    ).toBe(true);
  });

  it("waits while the canvas spec itself is in flight", () => {
    expect(
      noteContextReady({
        canvasData: undefined,
        hasTimeSeriesMap: new Map([["sales", true]]),
        interval: window,
      }),
    ).toBe(false);
  });
});

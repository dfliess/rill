import type { V1AnalystAgentContext } from "@rilldata/web-common/runtime-client";
import { describe, expect, it } from "vitest";
import { agentNoteIdempotencyKey, dashboardContextKey } from "./util";

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

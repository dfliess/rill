import { describe, expect, it } from "vitest";
import { composeRunId, parseRunId, runDetailPath } from "./run-id";

const INSTANCE = "643d3bbff115401a868d50a89758e10f";
const AGENT = "welcome-mcp";

describe("composeRunId / parseRunId", () => {
  it("round-trips a manual-run uuid key", () => {
    const key = "8efeb5ed-83d3-4840-9694-37e45b9a298c";
    const composite = composeRunId(INSTANCE, AGENT, key);
    expect(composite).toBe(
      `32:${INSTANCE}/11:${AGENT}/8efeb5ed-83d3-4840-9694-37e45b9a298c`,
    );
    expect(parseRunId(composite)).toEqual({
      instanceId: INSTANCE,
      agentName: AGENT,
      key,
    });
  });

  it("round-trips a key that itself contains '/'", () => {
    // A trigger-fired run's idempotency key may embed slashes; the length-prefix
    // parse must not split on them.
    const key = "alert/daily-sales/2026-07-17T00:00:00Z";
    const composite = composeRunId(INSTANCE, AGENT, key);
    expect(parseRunId(composite)).toEqual({
      instanceId: INSTANCE,
      agentName: AGENT,
      key,
    });
  });

  it("round-trips when the agent name contains a ':'", () => {
    const agent = "team:reporter";
    const key = "abc";
    const composite = composeRunId(INSTANCE, agent, key);
    expect(parseRunId(composite)).toEqual({
      instanceId: INSTANCE,
      agentName: agent,
      key,
    });
  });

  it("throws on a malformed composite", () => {
    expect(() => parseRunId("not-a-run-id")).toThrow();
    expect(() => parseRunId("999:short/1:a/k")).toThrow();
  });
});

describe("runDetailPath", () => {
  it("uses explicit agentName + idempotencyKey when present", () => {
    expect(
      runDetailPath("acme", "sales", {
        agentName: AGENT,
        idempotencyKey: "k-1",
        runId: "ignored",
      }),
    ).toBe(`/acme/sales/-/agents/${AGENT}/runs/k-1`);
  });

  it("derives agentName + key from a composite runId", () => {
    const key = "8efeb5ed-83d3-4840-9694-37e45b9a298c";
    expect(
      runDetailPath("acme", "sales", {
        runId: composeRunId(INSTANCE, AGENT, key),
      }),
    ).toBe(`/acme/sales/-/agents/${AGENT}/runs/${key}`);
  });

  it("encodes reserved characters in the key segment", () => {
    expect(
      runDetailPath("acme", "sales", {
        agentName: AGENT,
        idempotencyKey: "alert/daily/2026",
      }),
    ).toBe(`/acme/sales/-/agents/${AGENT}/runs/alert%2Fdaily%2F2026`);
  });
});

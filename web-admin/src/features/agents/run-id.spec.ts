import { describe, expect, it } from "vitest";
import { runDetailPath } from "./run-id";

describe("runDetailPath", () => {
  it("addresses the run by its id", () => {
    expect(
      runDetailPath("acme", "sales", {
        runId: "0199a3c2-7e4b-7c1a-9f3e-2b8d5a6c4e10",
      }),
    ).toBe("/acme/sales/-/agents/runs/0199a3c2-7e4b-7c1a-9f3e-2b8d5a6c4e10");
  });

  it("encodes the id, since a run recorded before ADR-0034 has slashes in it", () => {
    expect(
      runDetailPath("acme", "sales", { runId: "32:abc/7:triage/k" }),
    ).toBe("/acme/sales/-/agents/runs/32%3Aabc%2F7%3Atriage%2Fk");
  });
});

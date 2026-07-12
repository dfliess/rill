import { describe, it, expect } from "vitest";
import { resolveReportOpenTarget } from "./report-open-target";

function params(query: string): URLSearchParams {
  return new URL(`https://example.com/open${query}`).searchParams;
}

describe("resolveReportOpenTarget", () => {
  it("opens the AI conversation when session_id is present without a token", () => {
    expect(resolveReportOpenTarget(params("?session_id=abc123"))).toEqual({
      type: "ai",
      sessionId: "abc123",
    });
  });

  it("keeps other query params when routing to the AI conversation", () => {
    expect(
      resolveReportOpenTarget(
        params("?session_id=abc123&execution_time=2026-07-12T00:00:00Z"),
      ),
    ).toEqual({ type: "ai", sessionId: "abc123" });
  });

  it("treats an empty token as no token (recipient/owner mode links)", () => {
    expect(
      resolveReportOpenTarget(params("?session_id=abc123&token=")),
    ).toEqual({ type: "ai", sessionId: "abc123" });
  });

  it("falls back to explore when a magic-link token is present (creator mode)", () => {
    expect(
      resolveReportOpenTarget(params("?session_id=abc123&token=magic")),
    ).toEqual({ type: "explore" });
  });

  it("falls back to explore when session_id is empty", () => {
    expect(resolveReportOpenTarget(params("?session_id="))).toEqual({
      type: "explore",
    });
  });

  it("falls back to explore for a non-AI report link", () => {
    expect(
      resolveReportOpenTarget(
        params("?execution_time=2026-07-12T00:00:00Z&token=magic"),
      ),
    ).toEqual({ type: "explore" });
  });

  it("falls back to explore for an empty query string", () => {
    expect(resolveReportOpenTarget(params(""))).toEqual({ type: "explore" });
  });
});

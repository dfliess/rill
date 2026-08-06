import { describe, expect, it } from "vitest";
import {
  eventDecidedBy,
  isApprovalResolved,
  subjectLabel,
  timestampSortKey,
} from "./utils";

const NAMES = new Map([["f960bb37-9d8a", "Diego Kairos"]]);

describe("subjectLabel", () => {
  it("names a subject it can resolve", () => {
    expect(subjectLabel("f960bb37-9d8a", NAMES)).toBe("Diego Kairos");
  });

  it("falls back to the raw subject", () => {
    // A decider who left the project, or a caller who cannot list members: the
    // audit record still has to render.
    expect(subjectLabel("someone-else", NAMES)).toBe("someone-else");
    expect(subjectLabel("f960bb37-9d8a", undefined)).toBe("f960bb37-9d8a");
  });

  it("renders an empty subject as a dash", () => {
    expect(subjectLabel("", NAMES)).toBe("—");
    expect(subjectLabel(undefined, NAMES)).toBe("—");
  });
});

describe("eventDecidedBy", () => {
  it("reads the decider off a decision event's payload", () => {
    expect(eventDecidedBy({ decided_by: "admin:bob" })).toBe("admin:bob");
  });

  it("returns empty for events that carry no decision", () => {
    expect(eventDecidedBy(undefined)).toBe("");
    expect(eventDecidedBy(null)).toBe("");
    expect(eventDecidedBy({})).toBe("");
    expect(eventDecidedBy({ decided_by: 42 })).toBe("");
    expect(eventDecidedBy("run.resumed")).toBe("");
  });
});

describe("isApprovalResolved", () => {
  it("treats every terminal status as resolved", () => {
    for (const status of ["approved", "denied", "expired", "cancelled"]) {
      expect(isApprovalResolved(status)).toBe(true);
    }
  });

  it("does not resolve a pending or unknown-empty approval", () => {
    expect(isApprovalResolved("pending")).toBe(false);
    expect(isApprovalResolved("")).toBe(false);
    expect(isApprovalResolved(undefined)).toBe(false);
  });
});

describe("timestampSortKey", () => {
  it("orders RFC3339 strings chronologically", () => {
    const keys = ["2026-07-24T22:01:42Z", "2026-07-24T21:59:23Z"].map(
      timestampSortKey,
    );
    expect([...keys].sort()).toEqual([
      "2026-07-24T21:59:23Z",
      "2026-07-24T22:01:42Z",
    ]);
  });

  it("sorts non-string timestamps first rather than throwing", () => {
    expect(timestampSortKey(undefined)).toBe("");
    expect(timestampSortKey({ seconds: 1 })).toBe("");
  });
});

import { describe, expect, it, vi } from "vitest";
import {
  applyBadge,
  countPendingRuns,
  ledgerTotal,
  readPendingUpdate,
  scopeOf,
  type Badging,
} from "./app-badge";

describe("scopeOf", () => {
  it("keys a project the way the admin service does when it sends a push", () => {
    // admin/server/push_notifications.go composes strings.ToLower(org + "/" + project). The two have to agree
    // or a push and the open app would each keep their own entry for the same project and the badge would
    // double-count it.
    expect(scopeOf("Kairos", "Ecommerce")).toBe("kairos/ecommerce");
    expect(scopeOf("kairos", "ecommerce")).toBe("kairos/ecommerce");
  });
});

describe("countPendingRuns", () => {
  it("counts one run, not one per action it proposed", () => {
    // A batch turn leaves several approvals on one run, but it is one pause, one notification and one row in
    // the inbox. It must be one on the icon too.
    expect(
      countPendingRuns([
        { runId: "run-1", canDecide: true },
        { runId: "run-1", canDecide: true },
        { runId: "run-2", canDecide: true },
      ]),
    ).toBe(2);
  });

  it("counts only what this viewer may decide", () => {
    // Authority discriminates by action, so a run can be waiting on someone else entirely.
    expect(
      countPendingRuns([
        { runId: "run-1", canDecide: false },
        { runId: "run-2", canDecide: true },
      ]),
    ).toBe(1);
  });

  it("counts a run whose batch is only partly the viewer's", () => {
    expect(
      countPendingRuns([
        { runId: "run-1", canDecide: false },
        { runId: "run-1", canDecide: true },
      ]),
    ).toBe(1);
  });

  it("is zero for nothing, an empty list and approvals without a run", () => {
    expect(countPendingRuns(undefined)).toBe(0);
    expect(countPendingRuns([])).toBe(0);
    expect(countPendingRuns([{ canDecide: true }])).toBe(0);
  });
});

describe("readPendingUpdate", () => {
  it("reads a count for a project", () => {
    expect(
      readPendingUpdate({ scope: "kairos/ecommerce", count: 3, user: "u1" }),
    ).toEqual({ scope: "kairos/ecommerce", count: 3, user: "u1" });
  });

  it("reads a missing count as zero, because proto3 JSON omits one", () => {
    // This is the payload that clears the badge, so reading it as "no instruction" would leave the icon
    // showing work that is already done.
    expect(
      readPendingUpdate({ scope: "kairos/ecommerce", user: "u1" }),
    ).toEqual({ scope: "kairos/ecommerce", count: 0, user: "u1" });
  });

  it("ignores a payload that carries no count at all", () => {
    // Alerts and reports are read, not resolved: they say nothing about the badge and must leave it alone.
    expect(readPendingUpdate(undefined)).toBeNull();
    expect(readPendingUpdate(null)).toBeNull();
  });

  it("ignores a malformed one rather than putting a number nobody can explain on the icon", () => {
    expect(readPendingUpdate({ count: 2, user: "u1" })).toBeNull();
    expect(readPendingUpdate({ scope: "", count: 2, user: "u1" })).toBeNull();
    expect(
      readPendingUpdate({ scope: "kairos/ecommerce", count: -1, user: "u1" }),
    ).toBeNull();
    expect(
      readPendingUpdate({ scope: "kairos/ecommerce", count: 1.5, user: "u1" }),
    ).toBeNull();
    expect(
      readPendingUpdate({ scope: "kairos/ecommerce", count: "2", user: "u1" }),
    ).toBeNull();
  });

  it("ignores a count that names nobody, rather than crediting it to whoever is here", () => {
    // A browser outlives a session. A count with no owner cannot be attributed, and guessing would mean adding
    // one person's work to another person's icon.
    expect(
      readPendingUpdate({ scope: "kairos/ecommerce", count: 2 }),
    ).toBeNull();
    expect(
      readPendingUpdate({ scope: "kairos/ecommerce", count: 2, user: "" }),
    ).toBeNull();
  });
});

describe("ledgerTotal", () => {
  it("adds up what every project is waiting on", () => {
    expect(
      ledgerTotal(
        { "u1\nkairos/ecommerce": 2, "u1\nkairos/educacion": 1 },
        { scope: "kairos/logistica", count: 3, user: "u1" },
      ),
    ).toBe(6);
  });

  it("replaces the reporting project's count instead of adding to it", () => {
    expect(
      ledgerTotal(
        { "u1\nkairos/ecommerce": 2, "u1\nkairos/educacion": 1 },
        { scope: "kairos/ecommerce", count: 5, user: "u1" },
      ),
    ).toBe(6);
  });

  it("keeps the projects that did not report", () => {
    // The whole reason for a ledger: a push about one project, or the app being open on it, says nothing about
    // the others. Clearing them would hide work that is genuinely waiting.
    expect(
      ledgerTotal(
        { "u1\nkairos/ecommerce": 2, "u1\nkairos/educacion": 1 },
        { scope: "kairos/ecommerce", count: 0, user: "u1" },
      ),
    ).toBe(1);
  });

  it("reaches zero only when every project has reported nothing", () => {
    expect(
      ledgerTotal(
        { "u1\nkairos/ecommerce": 0, "u1\nkairos/educacion": 0 },
        { scope: "kairos/ecommerce", count: 0, user: "u1" },
      ),
    ).toBe(0);
  });

  it("skips an entry that is not a count", () => {
    expect(
      ledgerTotal(
        {
          "u1\nkairos/ecommerce": Number.NaN,
          "u1\nkairos/educacion": 1,
        } as Record<string, number>,
        { scope: "kairos/logistica", count: 2, user: "u1" },
      ),
    ).toBe(3);
  });

  it("leaves out what another account left behind on this browser", () => {
    // Two people signing in on the same browser share one ledger and one icon. Adding up both would show each of
    // them the other's pending work (found by Codex reviewing kairos-cloud#147).
    expect(
      ledgerTotal(
        { "alice\nkairos/ecommerce": 2, "bob\nkairos/educacion": 1 },
        { scope: "kairos/logistica", count: 1, user: "bob" },
      ),
    ).toBe(2);
  });

  it("keeps each account's projects apart even in the same project", () => {
    expect(
      ledgerTotal(
        { "alice\nkairos/ecommerce": 5 },
        { scope: "kairos/ecommerce", count: 1, user: "bob" },
      ),
    ).toBe(1);
  });
});

describe("applyBadge", () => {
  function badgingStub(): Badging & {
    setAppBadge: ReturnType<typeof vi.fn>;
    clearAppBadge: ReturnType<typeof vi.fn>;
  } {
    return {
      setAppBadge: vi.fn().mockResolvedValue(undefined),
      clearAppBadge: vi.fn().mockResolvedValue(undefined),
    };
  }

  it("puts the total on the icon", async () => {
    const badging = badgingStub();
    await applyBadge(badging, 3);
    expect(badging.setAppBadge).toHaveBeenCalledWith(3);
    expect(badging.clearAppBadge).not.toHaveBeenCalled();
  });

  it("takes the badge off when nothing is waiting", async () => {
    const badging = badgingStub();
    await applyBadge(badging, 0);
    expect(badging.clearAppBadge).toHaveBeenCalled();
    expect(badging.setAppBadge).not.toHaveBeenCalled();
  });

  it("does nothing in a browser without the API", async () => {
    // Firefox has none of it, and Safari exposes it only to an installed app. Neither must produce an error.
    await expect(applyBadge({}, 3)).resolves.toBeUndefined();
    await expect(applyBadge({}, 0)).resolves.toBeUndefined();
  });

  it("swallows a refused call", async () => {
    const badging: Badging = {
      setAppBadge: () => Promise.reject(new Error("not allowed")),
    };
    await expect(applyBadge(badging, 3)).resolves.toBeUndefined();
  });
});

import { describe, expect, it } from "vitest";
import { notificationTargetPath } from "./navigate-from-notification";

describe("notificationTargetPath", () => {
  it("keeps the path, query and hash of an absolute URL of this origin", () => {
    // What a push payload actually carries: the admin composes the frontend URL.
    expect(
      notificationTargetPath(
        `${window.location.origin}/acme/ecom/-/alerts/a1/open?execution_time=2026-08-20T08%3A00%3A00Z`,
      ),
    ).toBe(
      "/acme/ecom/-/alerts/a1/open?execution_time=2026-08-20T08%3A00%3A00Z",
    );
  });

  it("accepts a bare path", () => {
    expect(notificationTargetPath("/acme/ecom/-/reports/r1/open")).toBe(
      "/acme/ecom/-/reports/r1/open",
    );
  });

  it("refuses another origin: a tap must never leave the app", () => {
    expect(notificationTargetPath("https://example.com/phish")).toBeUndefined();
    expect(notificationTargetPath("//example.com/phish")).toBeUndefined();
  });

  it("refuses anything that is not a usable URL", () => {
    expect(notificationTargetPath(undefined)).toBeUndefined();
    expect(notificationTargetPath(42)).toBeUndefined();
    expect(notificationTargetPath("")).toBeUndefined();
  });
});

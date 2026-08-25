import { describe, expect, it } from "vitest";
import { listNotificationEntries } from "./notification-history";

describe("listNotificationEntries", () => {
  it("names the category with the message the preference switches already use", () => {
    expect(
      listNotificationEntries([
        { id: "1", category: "alerts" },
        { id: "2", category: "reports" },
        { id: "3", category: "act_approvals" },
      ]).map((e) => e.categoryKey),
    ).toEqual([
      "notifications_category_alerts",
      "notifications_category_reports",
      "notifications_category_act_approvals",
    ]);
  });

  it("keeps the server's order, which is newest first", () => {
    expect(
      listNotificationEntries([
        { id: "nuevo", category: "alerts" },
        { id: "viejo", category: "alerts" },
      ]).map((e) => e.id),
    ).toEqual(["nuevo", "viejo"]);
  });

  it("composes where it came from", () => {
    expect(
      listNotificationEntries([
        { id: "1", category: "alerts", org: "kairos", project: "ecommerce" },
      ])[0].scope,
    ).toBe("kairos/ecommerce");
  });

  it("drops a category this build has no name for", () => {
    // An older frontend meeting a newer server should show one entry fewer, not a row reading "webhooks".
    expect(
      listNotificationEntries([
        { id: "1", category: "webhooks" },
        { id: "2", category: "alerts" },
      ]).map((e) => e.id),
    ).toEqual(["2"]);
  });

  it("drops a row with no id, which nothing could key on", () => {
    expect(listNotificationEntries([{ category: "alerts" }])).toEqual([]);
  });

  it("treats a missing link as no destination, not as the home page", () => {
    // Proto3 JSON omits an empty string, and resolving "" would send the reader somewhere the notification was
    // never about.
    expect(
      listNotificationEntries([{ id: "1", category: "alerts" }])[0].linkPath,
    ).toBeUndefined();
    expect(
      listNotificationEntries([
        { id: "1", category: "alerts", linkPath: "" },
      ])[0].linkPath,
    ).toBeUndefined();
  });

  it("is empty for nothing at all", () => {
    expect(listNotificationEntries(undefined)).toEqual([]);
    expect(listNotificationEntries([])).toEqual([]);
  });
});

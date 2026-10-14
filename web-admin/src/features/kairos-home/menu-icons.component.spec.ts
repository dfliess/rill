import { LayoutDashboard, Wallet } from "lucide-svelte";
import { describe, expect, it } from "vitest";
import { homeMenuIcon } from "./menu-icons";

describe("homeMenuIcon", () => {
  it("maps a known name and falls back to the dashboard icon", () => {
    expect(homeMenuIcon("wallet")).toBe(Wallet);
    expect(homeMenuIcon(undefined)).toBe(LayoutDashboard);
    expect(homeMenuIcon("no-such-icon")).toBe(LayoutDashboard);
  });
});

import { LayoutDashboard, Wallet } from "lucide-svelte";
import { describe, expect, it } from "vitest";
import { homeV2MenuIcon } from "./menu-icons";

describe("homeV2MenuIcon", () => {
  it("maps a known name and falls back to the dashboard icon", () => {
    expect(homeV2MenuIcon("wallet")).toBe(Wallet);
    expect(homeV2MenuIcon(undefined)).toBe(LayoutDashboard);
    expect(homeV2MenuIcon("no-such-icon")).toBe(LayoutDashboard);
  });
});

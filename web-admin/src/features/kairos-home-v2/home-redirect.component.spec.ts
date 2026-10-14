// @vitest-environment jsdom
import { render, screen } from "@testing-library/svelte";
import { tick } from "svelte";
import { beforeEach, describe, expect, it, vi } from "vitest";

// The project's root URL (#231): with the new Home on it opens it, keeping
// the branch prefix and query; with it off it stays the classic Home. Nothing
// is decided before the project's flags have loaded.

const { page, flag, goto } = await vi.hoisted(async () => {
  const { writable } = await import("svelte/store");
  return {
    page: writable({ url: new URL("https://x.test/org/proj@dev?f=1") }),
    flag: writable(true),
    goto: vi.fn(),
  };
});

vi.mock("$app/stores", () => ({ page }));
vi.mock("$app/navigation", () => ({ goto }));
vi.mock("@rilldata/web-common/features/feature-flags", () => ({
  featureFlags: { kairosNewHome: flag, ready: Promise.resolve() },
}));
vi.mock(
  "@rilldata/web-admin/features/kairos-home/KairosProjectHome.svelte",
  async () => ({
    default: (await import("./__fixtures__/BoxStub.svelte")).default,
  }),
);

import RootPage from "../../routes/[organization]/[project]/+page.svelte";

async function settle() {
  await Promise.resolve();
  await tick();
  await tick();
}

beforeEach(() => goto.mockClear());

describe("project root with the new Home", () => {
  it("opens the new Home, keeping the branch and query", async () => {
    flag.set(true);
    render(RootPage);
    await settle();
    expect(goto).toHaveBeenCalledWith("/org/proj@dev/-/home-v2?f=1", {
      replaceState: true,
    });
    expect(screen.queryByTestId("home-box")).toBeNull();
  });

  it("stays the classic Home when the project turns it off", async () => {
    flag.set(false);
    render(RootPage);
    await settle();
    expect(goto).not.toHaveBeenCalled();
    expect(screen.getByTestId("home-box")).toBeTruthy();
  });
});

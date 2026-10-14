// @vitest-environment jsdom
import { fireEvent, render, screen } from "@testing-library/svelte";
import { createRawSnippet } from "svelte";
import { beforeEach, describe, expect, it, vi } from "vitest";

// Smoke test for the Home v2 tab (#231): the page loads with the dashboards
// menu, the home canvas with its Alerts and Actions boxes, and the AI panel
// open on wide screens and one button away on narrow ones. The canvas, the
// chat and the boxes are stubbed: each has its own tests.

// vi.mock factories are hoisted above the imports, so the stores they hand
// out are created in a hoisted block too.
const { page, resources } = await vi.hoisted(async () => {
  const { writable } = await import("svelte/store");
  return {
    page: writable({
      params: { organization: "org", project: "proj" } as Record<
        string,
        string
      >,
    }),
    resources: writable<Record<string, unknown>>({}),
  };
});

vi.mock("$app/stores", () => ({ page }));
vi.mock("$app/navigation", () => ({ goto: vi.fn() }));
vi.mock("@rilldata/web-common/runtime-client/v2", () => ({
  useRuntimeClient: () => ({ instanceId: "test-instance" }),
}));
vi.mock("@rilldata/web-common/runtime-client", async (importOriginal) => ({
  ...(await importOriginal<object>()),
  createRuntimeServiceListResources: () => resources,
}));
vi.mock("@rilldata/web-common/features/feature-flags", async () => {
  const { writable } = await import("svelte/store");
  return { featureFlags: { agents: writable(true), chat: writable(true) } };
});

vi.mock(
  "@rilldata/web-common/features/canvas/CanvasProvider.svelte",
  async () => ({
    default: (
      await import(
        "@rilldata/web-common/features/entity-management/__fixtures__/SlotPassthrough.svelte"
      )
    ).default,
  }),
);
vi.mock(
  "@rilldata/web-common/features/canvas/CanvasDashboardEmbed.svelte",
  async () => ({
    default: (await import("./__fixtures__/CanvasStub.svelte")).default,
  }),
);
vi.mock(
  "@rilldata/web-admin/features/ai/mcp/ConnectClientProvider.svelte",
  async () => ({
    default: (
      await import(
        "@rilldata/web-common/features/entity-management/__fixtures__/SlotPassthrough.svelte"
      )
    ).default,
  }),
);
vi.mock(
  "@rilldata/web-admin/features/kairos-home-v2/HomeV2Chat.svelte",
  async () => ({
    default: (await import("./__fixtures__/ChatStub.svelte")).default,
  }),
);
vi.mock(
  "@rilldata/web-admin/features/kairos-home/HomeRecentAlerts.svelte",
  async () => ({
    default: (await import("./__fixtures__/BoxStub.svelte")).default,
  }),
);
vi.mock(
  "@rilldata/web-admin/features/agents/HomePendingApprovals.svelte",
  async () => ({
    default: (await import("./__fixtures__/BoxStub.svelte")).default,
  }),
);

import Layout from "../../routes/[organization]/[project]/-/home-v2/+layout.svelte";

const canvas = (name: string, displayName: string, home = false) => ({
  meta: { name: { name } },
  canvas: {
    spec: { displayName },
    state: {
      validSpec: home ? { annotations: { kairos_home: "true" } } : {},
    },
  },
});

function setWideScreen(wide: boolean) {
  window.matchMedia = ((query: string) => ({
    matches: wide,
    media: query,
  })) as unknown as typeof window.matchMedia;
}

function renderLayout() {
  return render(Layout, {
    props: {
      children: createRawSnippet(() => ({ render: () => "<span></span>" })),
    },
  });
}

beforeEach(() => {
  page.set({ params: { organization: "org", project: "proj" } });
  resources.set({
    data: {
      resources: [
        canvas("resultados", "Resultados"),
        canvas("home", "Inicio · Facturación", true),
        canvas("caja_y_liquidez", "Caja y Liquidez"),
      ],
    },
  });
});

describe("Home v2 layout", () => {
  it("lists the Home first, then the other canvases by title", () => {
    setWideScreen(true);
    renderLayout();
    const links = screen
      .getAllByRole("link")
      .map((a) => [a.textContent?.trim(), a.getAttribute("href")]);
    expect(links).toEqual([
      ["Home", "/org/proj/-/home-v2"],
      ["Caja y Liquidez", "/org/proj/-/home-v2/caja_y_liquidez"],
      ["Resultados", "/org/proj/-/home-v2/resultados"],
    ]);
  });

  it("opens on the home canvas with its Alerts and Actions boxes", () => {
    setWideScreen(true);
    renderLayout();
    expect(screen.getByTestId("canvas").textContent).toBe("home");
    expect(screen.getAllByTestId("home-box")).toHaveLength(2);
  });

  it("shows the canvas named in the URL, without the Home boxes", () => {
    setWideScreen(true);
    page.set({
      params: { organization: "org", project: "proj", canvas: "resultados" },
    });
    renderLayout();
    expect(screen.getByTestId("canvas").textContent).toBe("resultados");
    expect(screen.queryAllByTestId("home-box")).toHaveLength(0);
    expect(screen.getByTestId("chat").textContent).toBe("resultados");
  });

  it("opens the AI panel on wide screens", async () => {
    setWideScreen(true);
    renderLayout();
    await Promise.resolve();
    expect(screen.getByTestId("chat")).toBeTruthy();
    expect(screen.queryByText("Ask the AI")).toBeNull();
  });

  it("keeps the AI panel one button away on narrow screens", async () => {
    setWideScreen(false);
    renderLayout();
    expect(screen.queryByTestId("chat")).toBeNull();
    await fireEvent.click(screen.getByText("Ask the AI"));
    expect(screen.getByTestId("chat")).toBeTruthy();
  });
});

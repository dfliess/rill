// @vitest-environment jsdom
import { render, screen } from "@testing-library/svelte";
import { createRawSnippet, tick } from "svelte";
import { beforeEach, describe, expect, it, vi } from "vitest";

// Smoke test for the Home (#231): the project URL opens with the dashboards
// menu, the first canvas with the Alerts and Actions boxes under it, and the
// AI panel open on wide screens and waiting for the header toggle on narrow
// ones; the viewer's personal canvases get their own group. The canvas, the
// chat, the boxes and the create dialog are stubbed, and the menu order has
// its own spec (menu.spec.ts).

// The canvas layout utils pull in the theme control, which reads
// matchMedia when its module loads: jsdom has none.
vi.hoisted(() => {
  window.matchMedia = ((query: string) => ({
    matches: false,
    media: query,
    addEventListener: () => {},
    removeEventListener: () => {},
  })) as unknown as typeof window.matchMedia;
});

// vi.mock factories are hoisted above the imports, so the stores they hand
// out are created in a hoisted block too.
const { page, resources, personal, user, personalCanvasesFlag } =
  await vi.hoisted(async () => {
    const { writable } = await import("svelte/store");
    return {
      page: writable({
        params: { organization: "org", project: "proj" } as Record<
          string,
          string
        >,
      }),
      resources: writable<Record<string, unknown>>({}),
      personal: writable<Record<string, unknown>>({}),
      user: writable<Record<string, unknown>>({}),
      personalCanvasesFlag: writable(false),
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
  return {
    featureFlags: {
      agents: writable(true),
      chat: writable(true),
      personalCanvases: personalCanvasesFlag,
    },
  };
});
vi.mock("@rilldata/web-admin/client", () => ({
  createAdminServiceGetCurrentUser: () => user,
}));
vi.mock("@rilldata/web-admin/features/personal-files/selectors", () => ({
  getPersonalFilteredResources: () => personal,
}));
vi.mock(
  "@rilldata/web-admin/features/personal-files/canvas/CreatePersonalCanvasDialog.svelte",
  async () => ({
    default: (await import("./__fixtures__/CreatePersonalCanvasStub.svelte"))
      .default,
  }),
);

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
  "@rilldata/web-admin/features/kairos-home/HomeChat.svelte",
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

import { homeChatActions } from "./chat-state";
import Layout from "../../routes/[organization]/[project]/(home)/+layout.svelte";

const canvas = (
  name: string,
  displayName: string,
  metadata?: Record<string, string>,
  maxWidth?: number,
) => ({
  meta: { name: { name }, metadata },
  canvas: { spec: { displayName }, state: { validSpec: { maxWidth } } },
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
  personalCanvasesFlag.set(false);
  user.set({});
  personal.set({});
  resources.set({
    data: {
      resources: [
        canvas("resultados", "Resultados"),
        canvas("home", "Inicio · Facturación", { kairos_order: "1" }),
        canvas("caja_y_liquidez", "Caja y Liquidez"),
      ],
    },
  });
});

describe("Home layout", () => {
  it("lists the canvases by display name, the ordered one first", () => {
    setWideScreen(true);
    renderLayout();
    const links = screen
      .getAllByRole("link")
      .map((a) => [a.textContent?.trim(), a.getAttribute("href")]);
    expect(links).toEqual([
      ["Inicio · Facturación", "/org/proj"],
      ["Caja y Liquidez", "/org/proj/-/home/caja_y_liquidez"],
      ["Resultados", "/org/proj/-/home/resultados"],
    ]);
  });

  it("opens on the first canvas with the Alerts and Actions boxes", () => {
    setWideScreen(true);
    renderLayout();
    expect(screen.getByTestId("canvas").textContent).toBe("home");
    expect(screen.getAllByTestId("home-box")).toHaveLength(2);
  });

  it("lines the boxes up with the canvas's max width", () => {
    setWideScreen(true);
    const grid = () => screen.getAllByTestId("home-box")[0].parentElement!;
    const { unmount } = renderLayout();
    expect(grid().style.maxWidth).toBe("1200px");
    unmount();
    resources.set({
      data: {
        resources: [canvas("home", "Inicio", { kairos_order: "1" }, 1400)],
      },
    });
    renderLayout();
    expect(grid().style.maxWidth).toBe("1400px");
  });

  it("shows the canvas named in the URL, without the boxes", () => {
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
    await tick();
    expect(screen.getByTestId("chat")).toBeTruthy();
  });

  it("leaves the AI panel to the header toggle on narrow screens", async () => {
    setWideScreen(false);
    renderLayout();
    await tick();
    expect(screen.queryByTestId("chat")).toBeNull();
    homeChatActions.openChat();
    await tick();
    expect(screen.getByTestId("chat")).toBeTruthy();
  });

  describe("personal canvases", () => {
    const mine = {
      meta: { name: { name: "mi_vista" } },
      canvas: {
        spec: { displayName: "Mi vista" },
        // Shared, so only being the viewer's keeps it out of the main list.
        state: {
          validSpec: {
            annotations: { admin_managed: "true", admin_shared: "true" },
          },
        },
      },
    };

    function showPersonal() {
      personalCanvasesFlag.set(true);
      user.set({ data: { user: { id: "u1" } } });
      personal.set({ isPending: false, data: [mine] });
      resources.update((r) => {
        const data = r.data as { resources: unknown[] };
        return { data: { resources: [...data.resources, mine] } };
      });
    }

    it("get their own group with the create button", () => {
      setWideScreen(true);
      showPersonal();
      renderLayout();
      const group = screen.getByRole("navigation", { name: "My dashboards" });
      expect(
        Array.from(group.querySelectorAll("a")).map((a) => [
          a.textContent?.trim() || a.getAttribute("aria-label"),
          a.getAttribute("href"),
        ]),
      ).toEqual([
        ["Mi vista", "/org/proj/-/home/mi_vista"],
        ["Edit Mi vista", "/org/proj/-/personal/mi_vista"],
      ]);
      expect(screen.getByTestId("create-personal").dataset.project).toBe(
        "proj",
      );
      const main = screen.getByRole("navigation", { name: "Dashboards" });
      expect(main.textContent).not.toContain("Mi vista");
    });

    it("open inside the Home, marked in the menu", () => {
      setWideScreen(true);
      showPersonal();
      page.set({
        params: { organization: "org", project: "proj", canvas: "mi_vista" },
      });
      renderLayout();
      expect(screen.getByTestId("canvas").textContent).toBe("mi_vista");
      expect(screen.queryAllByTestId("home-box")).toHaveLength(0);
      const group = screen.getByRole("navigation", { name: "My dashboards" });
      expect(
        group.querySelector('[aria-current="page"]')?.textContent?.trim(),
      ).toBe("Mi vista");
    });

    it("show in the phone selector too", () => {
      setWideScreen(true);
      showPersonal();
      renderLayout();
      const group = document.querySelector("optgroup")!;
      expect(group.label).toBe("My dashboards");
      expect(group.textContent).toContain("Mi vista");
    });

    it("are hidden with the feature off", () => {
      setWideScreen(true);
      showPersonal();
      personalCanvasesFlag.set(false);
      renderLayout();
      expect(
        screen.queryByRole("navigation", { name: "My dashboards" }),
      ).toBeNull();
      expect(screen.queryByTestId("create-personal")).toBeNull();
    });

    it("are hidden with nobody signed in", () => {
      setWideScreen(true);
      showPersonal();
      user.set({ data: {} });
      renderLayout();
      expect(
        screen.queryByRole("navigation", { name: "My dashboards" }),
      ).toBeNull();
      expect(screen.queryByTestId("create-personal")).toBeNull();
    });
  });
});

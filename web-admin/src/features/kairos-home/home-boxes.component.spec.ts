// @vitest-environment jsdom
import { render, screen } from "@testing-library/svelte";
import { writable } from "svelte/store";
import { beforeEach, describe, expect, it, vi } from "vitest";

// Home's Alerts and Actions boxes are always shown (#230): an empty box tells
// the viewer where Monitor and Act will speak up, a missing one hides the
// feature. These cover the three states each box can be in: loading (heading
// only, no premature "nothing here"), loaded and empty, and loaded with items.

const runs = writable<Record<string, unknown>>({});
const approvals = writable<Record<string, unknown>>({});
const resources = writable<Record<string, unknown>>({});

vi.mock("@rilldata/web-common/runtime-client/v2", () => ({
  useRuntimeClient: () => ({ instanceId: "test-instance" }),
}));

vi.mock("@rilldata/web-admin/features/agents/selectors", () => ({
  useAgentRuns: () => runs,
  useAgentApprovals: () => approvals,
  useAgents: () =>
    writable({
      data: { agents: [{ name: "resumen", displayName: "Resumen semanal" }] },
    }),
}));

vi.mock("@rilldata/web-common/runtime-client", async (importOriginal) => ({
  ...(await importOriginal<object>()),
  createRuntimeServiceListResources: () => resources,
}));

// The decision buttons run their own mutations; this file only cares that a
// pending row offers them in place.
vi.mock(
  "@rilldata/web-admin/features/agents/approvals/ApproveDenyButtons.svelte",
  async () => ({
    default: (await import("./__fixtures__/ApproveDenyStub.svelte")).default,
  }),
);

import HomePendingApprovals from "@rilldata/web-admin/features/agents/HomePendingApprovals.svelte";
import HomeRecentAlerts from "./HomeRecentAlerts.svelte";

const props = { organization: "org", project: "proj" };

beforeEach(() => {
  runs.set({});
  approvals.set({});
  resources.set({});
});

describe("HomePendingApprovals", () => {
  it("shows the heading but no empty text while loading", () => {
    render(HomePendingApprovals, { props });
    expect(screen.getByText("Pending approvals")).toBeTruthy();
    expect(
      screen.queryByText("No actions are waiting for your approval."),
    ).toBeNull();
  });

  it("says there is nothing pending once both lists are empty", () => {
    runs.set({ isSuccess: true, data: { runs: [] } });
    approvals.set({ isSuccess: true, data: { approvals: [] } });
    render(HomePendingApprovals, { props });
    expect(screen.getByText("Pending approvals")).toBeTruthy();
    expect(
      screen.getByText("No actions are waiting for your approval."),
    ).toBeTruthy();
    expect(screen.getByText("View all").getAttribute("href")).toBe(
      "/org/proj/-/agents",
    );
  });

  it("drops the empty text when a decidable approval is pending", () => {
    runs.set({
      isSuccess: true,
      data: { runs: [{ runId: "r1", status: "waiting_approval" }] },
    });
    approvals.set({
      isSuccess: true,
      data: {
        approvals: [{ runId: "r1", status: "pending", canDecide: true }],
      },
    });
    render(HomePendingApprovals, { props });
    expect(
      screen.queryByText("No actions are waiting for your approval."),
    ).toBeNull();
  });

  it("names the agent and what it would do, with approve in place", () => {
    runs.set({
      isSuccess: true,
      data: {
        runs: [
          { runId: "r1", agentName: "resumen", status: "waiting_approval" },
        ],
      },
    });
    approvals.set({
      isSuccess: true,
      data: {
        approvals: [
          {
            runId: "r1",
            approvalId: "a1",
            status: "pending",
            canDecide: true,
            canonicalArgs: '{"channel":"#direccion","subject":"Resumen"}',
          },
        ],
      },
    });
    render(HomePendingApprovals, { props });
    expect(screen.getByText("Resumen semanal")).toBeTruthy();
    expect(screen.getByText(/Resumen · #direccion/)).toBeTruthy();
    expect(screen.getByTestId("approve-deny").dataset.approval).toBe("a1");
    expect(screen.getByRole("heading").textContent).toContain("1");
  });
});

describe("HomeRecentAlerts", () => {
  it("shows the heading but no empty text while loading", () => {
    render(HomeRecentAlerts, { props });
    expect(screen.getByText("Recent alerts")).toBeTruthy();
    expect(screen.queryByText("No alert has fired yet.")).toBeNull();
  });

  it("says no alert has fired once resources are loaded", () => {
    resources.set({ isSuccess: true, data: { resources: [] } });
    render(HomeRecentAlerts, { props });
    expect(screen.getByText("Recent alerts")).toBeTruthy();
    expect(screen.getByText("No alert has fired yet.")).toBeTruthy();
    expect(screen.getByText("View all").getAttribute("href")).toBe(
      "/org/proj/-/alerts",
    );
  });

  it("lists a fired alert instead of the empty text", () => {
    resources.set({
      isSuccess: true,
      data: {
        resources: [
          {
            meta: { name: { name: "margen_bajo" } },
            alert: {
              spec: { displayName: "Margen bajo" },
              state: {
                executionHistory: [
                  {
                    finishedOn: "2026-10-13T07:00:00Z",
                    result: { status: "ASSERTION_STATUS_FAIL" },
                  },
                ],
              },
            },
          },
        ],
      },
    });
    render(HomeRecentAlerts, { props });
    expect(screen.getByText("Margen bajo")).toBeTruthy();
    expect(screen.queryByText("No alert has fired yet.")).toBeNull();
  });
});

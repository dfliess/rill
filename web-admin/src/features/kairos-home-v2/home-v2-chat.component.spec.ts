// @vitest-environment jsdom
import { fireEvent, render, screen } from "@testing-library/svelte";
import { describe, expect, it, vi } from "vitest";

// The AI panel's canvas context (#231): the open canvas shows as a removable
// label in the input footer, after @ and /, and opening another canvas brings
// it back. The chat itself is stubbed to a footer that renders the label.

vi.mock(
  "@rilldata/web-common/features/chat/layouts/sidebar/SidebarChat.svelte",
  async () => ({
    default: (await import("./__fixtures__/SidebarChatStub.svelte")).default,
  }),
);

import HomeV2Chat from "./HomeV2Chat.svelte";

const actions = {} as never;

describe("Home v2 chat context", () => {
  it("shows the open canvas in the input footer", () => {
    render(HomeV2Chat, {
      props: { canvasName: "home", canvasTitle: "Inicio", actions },
    });
    expect(screen.getByTestId("input-footer").textContent).toContain("Inicio");
    expect(
      screen.getByRole("button", { name: /todo el proyecto|whole project/i }),
    ).toBeTruthy();
  });

  it("drops the canvas on ✕ and brings it back on another canvas", async () => {
    const { rerender } = render(HomeV2Chat, {
      props: { canvasName: "home", canvasTitle: "Inicio", actions },
    });
    await fireEvent.click(
      screen.getByRole("button", { name: /todo el proyecto|whole project/i }),
    );
    expect(
      screen.queryByRole("button", { name: /todo el proyecto|whole project/i }),
    ).toBeNull();

    await rerender({ canvasName: "caja", canvasTitle: "Caja", actions });
    expect(
      screen.getByRole("button", { name: /todo el proyecto|whole project/i }),
    ).toBeTruthy();
    expect(screen.getByTestId("input-footer").textContent).toContain("Caja");
  });

  it("shows nothing when no canvas is open", () => {
    render(HomeV2Chat, {
      props: { canvasName: undefined, canvasTitle: "", actions },
    });
    expect(screen.getByTestId("input-footer").textContent?.trim()).toBe("");
  });
});

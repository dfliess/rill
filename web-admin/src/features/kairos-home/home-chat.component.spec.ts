// @vitest-environment jsdom
import { fireEvent, render, screen } from "@testing-library/svelte";
import { get, type Readable } from "svelte/store";
import { describe, expect, it, vi } from "vitest";

// The AI panel's canvas context (#231): none until the viewer attaches the
// open canvas from the input footer; once attached it follows the open canvas
// until removed. The chat is stubbed to a footer that renders the label, and
// the config captures the context store the messages are sent with.

const captured = vi.hoisted(() => ({
  context: undefined as Readable<string | undefined> | undefined,
}));

vi.mock(
  "@rilldata/web-common/features/chat/layouts/sidebar/SidebarChat.svelte",
  async () => ({
    default: (await import("./__fixtures__/SidebarChatStub.svelte")).default,
  }),
);
vi.mock("./chat-config", () => ({
  createHomeChatConfig: (context: Readable<string | undefined>) => {
    captured.context = context;
    return {};
  },
}));

import HomeChat from "./HomeChat.svelte";

const actions = {} as never;
const attach = /como contexto|as context/i;
const remove = /todo el proyecto|whole project/i;
const sentContext = () => get(captured.context!);

describe("Home chat context", () => {
  it("attaches nothing by default and offers the open canvas", () => {
    render(HomeChat, {
      props: { canvasName: "home", canvasTitle: "Inicio", actions },
    });
    expect(screen.getByRole("button", { name: attach })).toBeTruthy();
    expect(screen.queryByRole("button", { name: remove })).toBeNull();
    expect(sentContext()).toBeUndefined();
  });

  it("attaches the canvas on click and follows the open canvas", async () => {
    const { rerender } = render(HomeChat, {
      props: { canvasName: "home", canvasTitle: "Inicio", actions },
    });
    await fireEvent.click(screen.getByRole("button", { name: attach }));
    expect(sentContext()).toBe("home");
    expect(screen.getByRole("button", { name: remove })).toBeTruthy();

    await rerender({ canvasName: "caja", canvasTitle: "Caja", actions });
    expect(sentContext()).toBe("caja");
    expect(screen.getByTestId("input-footer").textContent).toContain("Caja");
  });

  it("stays project-wide after removing it, across canvases", async () => {
    const { rerender } = render(HomeChat, {
      props: { canvasName: "home", canvasTitle: "Inicio", actions },
    });
    await fireEvent.click(screen.getByRole("button", { name: attach }));
    await fireEvent.click(screen.getByRole("button", { name: remove }));
    await rerender({ canvasName: "caja", canvasTitle: "Caja", actions });
    expect(sentContext()).toBeUndefined();
    expect(screen.getByRole("button", { name: attach })).toBeTruthy();
  });

  it("shows nothing when no canvas is open", () => {
    render(HomeChat, {
      props: { canvasName: undefined, canvasTitle: "", actions },
    });
    expect(screen.getByTestId("input-footer").textContent?.trim()).toBe("");
  });
});

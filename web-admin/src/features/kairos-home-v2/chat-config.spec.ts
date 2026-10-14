import { get, writable } from "svelte/store";
import { describe, expect, it } from "vitest";
import { createHomeV2ChatConfig } from "./chat-config";

describe("createHomeV2ChatConfig", () => {
  it("sends the open canvas as context while the viewer keeps it", () => {
    const canvas = writable<string | undefined>("caja_y_liquidez");
    const context =
      createHomeV2ChatConfig(canvas).additionalContextStoreGetter();
    expect(get(context)).toEqual({
      analystAgentContext: { canvas: "caja_y_liquidez" },
    });

    // Dropping the chip, or a page with no canvas, leaves a plain project chat.
    canvas.set(undefined);
    expect(get(context)).toEqual({});
  });

  it("follows the canvas the viewer switches to", () => {
    const canvas = writable<string | undefined>("home");
    const context =
      createHomeV2ChatConfig(canvas).additionalContextStoreGetter();
    canvas.set("resultados");
    expect(get(context)).toEqual({
      analystAgentContext: { canvas: "resultados" },
    });
  });
});

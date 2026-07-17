import type { V1Message } from "@rilldata/web-common/runtime-client";
import {
  MessageContentType,
  MessageType,
  ToolName,
} from "@rilldata/web-common/features/chat/core/types.ts";
import { describe, it, expect } from "vitest";
import { transformToBlocks, type TextBlock } from "./block-transform.ts";

// transformToBlocks is a pure function of the message list, so these run without a runtime.
function transform(messages: V1Message[]) {
  return transformToBlocks(messages, false, false);
}

describe("transformToBlocks routing", () => {
  it("renders a plain user text message as a text block, not skipped", () => {
    const userMsg: V1Message = {
      id: "msg-user-1",
      role: "user",
      type: MessageType.TEXT,
      contentType: MessageContentType.TEXT,
      contentData: "hola",
    };

    const blocks = transform([userMsg]);

    const textBlocks = blocks.filter((b): b is TextBlock => b.type === "text");
    expect(textBlocks).toHaveLength(1);
    expect(textBlocks[0].id).toBe("msg-user-1");
    expect(textBlocks[0].message).toBe(userMsg);
  });

  it("renders a dynamic agent's assistant text turn as a text block (its closing answer)", () => {
    // A segmented governed run persists the model's closing as a plain assistant text turn (no tool tag). It must
    // render as a bubble, not be skipped, or the point of the run — the model's answer after the action ran — is lost.
    const assistantMsg: V1Message = {
      id: "msg-assistant-close",
      role: "assistant",
      type: MessageType.TEXT,
      contentType: MessageContentType.TEXT,
      contentData: "Listo, creé el ticket PROJ-42.",
    };

    const blocks = transform([assistantMsg]);

    const textBlocks = blocks.filter((b): b is TextBlock => b.type === "text");
    expect(textBlocks).toHaveLength(1);
    expect(textBlocks[0].id).toBe("msg-assistant-close");
    expect(textBlocks[0].message).toBe(assistantMsg);
  });

  it("hides the injected action-result turn (internal prompt, not a user message)", () => {
    // A governed run injects the executed action's result to feed the model on resume, tagged with ACTION_RESULT. It is
    // an internal prompt, not something the user said, so it must not render as a bubble — the model's closing conveys it.
    const injected: V1Message = {
      id: "msg-injected",
      role: "user",
      tool: ToolName.ACTION_RESULT,
      type: MessageType.TEXT,
      contentType: MessageContentType.TEXT,
      contentData:
        'The action "mcp.mockmcp.HelloWorld" was approved and executed. Result: {...}. Reply to the user with this outcome.',
    };

    const blocks = transform([injected]);

    expect(blocks.filter((b) => b.type === "text")).toHaveLength(0);
  });

  it("still renders an assistant router_agent user turn as text (unchanged)", () => {
    // The assistant chat models a user turn as a router_agent CALL with role=user; it must keep routing to text.
    const routerUserMsg: V1Message = {
      id: "msg-router-1",
      role: "user",
      tool: ToolName.ROUTER_AGENT,
      type: MessageType.CALL,
      contentType: MessageContentType.JSON,
      contentData: JSON.stringify({ prompt: "hola" }),
    };

    const blocks = transform([routerUserMsg]);

    const textBlocks = blocks.filter((b): b is TextBlock => b.type === "text");
    expect(textBlocks).toHaveLength(1);
    expect(textBlocks[0].id).toBe("msg-router-1");
  });
});

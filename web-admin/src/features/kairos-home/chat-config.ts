import type { ChatConfig } from "@rilldata/web-common/features/chat/core/types.ts";
import { projectChat } from "@rilldata/web-common/features/project/chat-context.ts";
import { derived, type Readable } from "svelte/store";

/**
 * Home's AI panel is the project chat (same agent and history as the AI
 * tab), so a conversation survives switching canvases in the menu. The open
 * canvas rides along as context while the viewer keeps it: it tells the
 * analyst what is on screen without tying the panel to one dashboard.
 */
export function createHomeChatConfig(
  canvasContext: Readable<string | undefined>,
): ChatConfig {
  return {
    ...projectChat,
    additionalContextStoreGetter: () =>
      derived(canvasContext, (canvas) =>
        canvas ? { analystAgentContext: { canvas } } : {},
      ),
  };
}

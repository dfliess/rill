<script lang="ts">
  import SidebarChat from "@rilldata/web-common/features/chat/layouts/sidebar/SidebarChat.svelte";
  import type { ChatActions } from "@rilldata/web-common/features/chat/layouts/sidebar/sidebar-store";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { writable } from "svelte/store";
  import { createHomeV2ChatConfig } from "./chat-config";

  let {
    canvasName,
    canvasTitle,
    actions,
  }: {
    canvasName: string | undefined;
    canvasTitle: string;
    actions: ChatActions;
  } = $props();

  // The open canvas is context by default; the viewer can drop it for a
  // project-wide question. Opening another canvas brings it back.
  let keepCanvas = $state(true);
  $effect(() => {
    void canvasName;
    keepCanvas = true;
  });

  const canvasContext = writable<string | undefined>(undefined);
  $effect(() => {
    canvasContext.set(keepCanvas ? canvasName : undefined);
  });

  // Built once: the panel keeps one conversation across canvas switches.
  const config = createHomeV2ChatConfig(canvasContext);
</script>

<div class="flex h-full flex-col">
  {#if canvasName}
    <div
      class="flex items-center gap-x-2 border-b border-l bg-surface-background px-3 py-2 text-xs text-fg-secondary"
    >
      {#if keepCanvas}
        <span
          class="flex min-w-0 items-center gap-x-1 rounded-full border px-2 py-0.5"
        >
          <span class="truncate"
            >{m.home_v2_chat_context({ canvas: canvasTitle })}</span
          >
          <button
            class="shrink-0 hover:text-fg-primary"
            aria-label={m.home_v2_chat_context_remove()}
            onclick={() => (keepCanvas = false)}>✕</button
          >
        </span>
      {:else}
        <button
          class="truncate hover:text-fg-primary"
          onclick={() => (keepCanvas = true)}
        >
          {m.home_v2_chat_context_add({ canvas: canvasTitle })}
        </button>
      {/if}
    </div>
  {/if}
  <div class="min-h-0 flex-1">
    <SidebarChat {config} {actions} surface="dashboard" />
  </div>
</div>

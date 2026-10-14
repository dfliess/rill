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

<!-- The context sits in the input footer, after @ and /, next to what it
     applies to: the next message. -->
{#snippet canvasChip()}
  {#if canvasName}
    {#if keepCanvas}
      <span
        class="flex min-w-0 items-center gap-x-1 rounded-full border px-2 py-0.5 text-xs text-fg-secondary"
      >
        <span class="truncate"
          >{m.home_v2_chat_context({ canvas: canvasTitle })}</span
        >
        <button
          type="button"
          class="shrink-0 hover:text-fg-primary"
          aria-label={m.home_v2_chat_context_remove()}
          onclick={() => (keepCanvas = false)}>✕</button
        >
      </span>
    {:else}
      <button
        type="button"
        class="truncate text-xs text-fg-muted hover:text-fg-primary"
        onclick={() => (keepCanvas = true)}
      >
        {m.home_v2_chat_context_add({ canvas: canvasTitle })}
      </button>
    {/if}
  {/if}
{/snippet}

<div class="h-full">
  <SidebarChat
    {config}
    {actions}
    surface="dashboard"
    inputContext={canvasChip}
  />
</div>

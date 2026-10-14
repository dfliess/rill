<script lang="ts">
  import SidebarChat from "@rilldata/web-common/features/chat/layouts/sidebar/SidebarChat.svelte";
  import type { ChatActions } from "@rilldata/web-common/features/chat/layouts/sidebar/sidebar-store";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { LayoutDashboard, Plus, X } from "lucide-svelte";
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

  // No context until the viewer attaches the open canvas: a Home question is
  // usually about the whole project, and an overview canvas would narrow it.
  // Once attached it follows the open canvas ("what I am looking at") until
  // the viewer removes it (kairos-cloud#231).
  let keepCanvas = $state(false);

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
    <!-- Quiet, like an attachment in Claude or ChatGPT: just the canvas,
         offset from @ and /, the full wording in the tooltip. -->
    {#if keepCanvas}
      <span
        class="ml-2 flex min-w-0 items-center gap-x-1 rounded-md bg-surface-hover px-1.5 py-0.5 text-xs text-fg-muted"
        title={m.home_v2_chat_context({ canvas: canvasTitle })}
      >
        <LayoutDashboard size="12px" class="shrink-0" />
        <span class="truncate">{canvasTitle}</span>
        <button
          type="button"
          class="shrink-0 rounded hover:text-fg-primary"
          aria-label={m.home_v2_chat_context_remove()}
          title={m.home_v2_chat_context_remove()}
          onclick={() => (keepCanvas = false)}><X size="12px" /></button
        >
      </span>
    {:else}
      <button
        type="button"
        class="ml-2 flex min-w-0 items-center gap-x-1 rounded-md border border-dashed px-1.5 py-0.5 text-xs text-fg-muted hover:text-fg-primary"
        title={m.home_v2_chat_context_add({ canvas: canvasTitle })}
        aria-label={m.home_v2_chat_context_add({ canvas: canvasTitle })}
        onclick={() => (keepCanvas = true)}
      >
        <Plus size="12px" class="shrink-0" />
        <span class="truncate">{canvasTitle}</span>
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

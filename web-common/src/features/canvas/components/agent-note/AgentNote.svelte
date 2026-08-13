<script lang="ts">
  import { extractErrorMessage } from "@rilldata/web-common/lib/errors";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import { createQuery } from "@tanstack/svelte-query";
  import DOMPurify from "dompurify";
  import { marked } from "marked";
  import type { AgentNoteCanvasComponent } from "./";
  import { agentNoteIdempotencyKey, runAgentNote } from "./util";

  export let component: AgentNoteCanvasComponent;

  const runtimeClient = useRuntimeClient();

  // Bumped by "regenerate" so the next run gets a fresh idempotency key instead of attaching to the
  // completed one. Without it a re-read would return the same cached answer.
  let nonce = 0;
  let requested = false;

  $: specStore = component?.specStore;
  $: spec = specStore ? $specStore : undefined;
  $: agent = spec?.agent ?? "";
  $: prompt = spec?.prompt ?? "";
  $: autoRun = spec?.auto_run !== false;
  $: configured = agent.trim().length > 0 && prompt.trim().length > 0;

  $: idempotencyKey = agentNoteIdempotencyKey({
    agent,
    prompt,
    context: runtimeClient.instanceId ?? "",
    nonce,
  });

  $: queryOptions = {
    queryKey: ["agent-note", runtimeClient.instanceId, idempotencyKey],
    queryFn: ({ signal }: { signal: AbortSignal }) =>
      runAgentNote(runtimeClient, { agent, prompt, idempotencyKey }, signal),
    enabled: configured && (autoRun || requested),
    // The answer is pinned to its idempotency key, so it never goes stale on its own: a new answer
    // requires a new key, which only "regenerate" or an edited prompt produces.
    staleTime: Infinity,
    gcTime: Infinity,
    retry: false,
    refetchOnWindowFocus: false,
  };
  $: noteQuery = createQuery(queryOptions, queryClient);

  $: note = $noteQuery?.data ?? "";
  $: renderPromise = marked(note || "");
  $: errorMessage = $noteQuery?.error
    ? extractErrorMessage($noteQuery.error)
    : "";

  function regenerate() {
    nonce += 1;
    requested = true;
  }
</script>

<div class="size-full flex flex-col bg-surface-card px-3 py-2 overflow-y-auto">
  {#if !configured}
    <p class="agent-note-hint">{m.canvas_agent_note_unconfigured()}</p>
  {:else if $noteQuery?.isFetching}
    <div class="agent-note-loading" aria-live="polite">
      <span class="agent-note-pulse"></span>
      <span>{m.canvas_agent_note_generating()}</span>
    </div>
  {:else if $noteQuery?.isError}
    <div class="agent-note-error">
      <p>{errorMessage}</p>
      <button type="button" on:click={regenerate}>
        {m.canvas_agent_note_regenerate()}
      </button>
    </div>
  {:else if note}
    <div class="agent-note select-text cursor-text">
      {#await renderPromise then html}
        {@html DOMPurify.sanitize(html)}
      {/await}
    </div>
    <div class="agent-note-footer">
      <button type="button" on:click={regenerate}>
        {m.canvas_agent_note_regenerate()}
      </button>
    </div>
  {:else}
    <div class="agent-note-idle">
      <button type="button" on:click={regenerate}>
        {m.canvas_agent_note_generate()}
      </button>
    </div>
  {/if}
</div>

<style lang="postcss">
  .agent-note {
    @apply text-fg-primary flex-1 min-h-min;
  }
  :global(.agent-note p) {
    font-size: 14px;
    @apply my-1;
  }
  :global(.agent-note ul) {
    @apply list-disc pl-6 my-2;
  }
  :global(.agent-note li) {
    @apply text-sm my-1;
  }
  :global(.agent-note strong) {
    @apply font-medium;
  }

  .agent-note-hint,
  .agent-note-idle {
    @apply text-sm text-fg-muted;
  }
  .agent-note-loading {
    @apply flex items-center gap-x-2 text-sm text-fg-muted;
  }
  .agent-note-pulse {
    @apply size-2 rounded-full bg-primary-solid;
    animation: agent-note-pulse 1.2s ease-in-out infinite;
  }
  @keyframes agent-note-pulse {
    0%,
    100% {
      opacity: 0.3;
    }
    50% {
      opacity: 1;
    }
  }
  .agent-note-error {
    @apply text-sm text-fg-secondary;
  }
  .agent-note-footer {
    @apply pt-1;
  }
  button {
    @apply text-xs text-primary-solid hover:underline;
  }
</style>

<script lang="ts">
  import { extractErrorMessage } from "@rilldata/web-common/lib/errors";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";
  import type { V1AnalystAgentContext } from "@rilldata/web-common/runtime-client";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import { createQuery } from "@tanstack/svelte-query";
  import DOMPurify from "dompurify";
  import { marked } from "marked";
  import type { AgentNoteCanvasComponent } from "./";
  import {
    agentNoteIdempotencyKey,
    dashboardContextKey,
    runAgentNote,
  } from "./util";

  export let component: AgentNoteCanvasComponent;

  const runtimeClient = useRuntimeClient();

  // Bumped by "regenerate" so the next run gets a fresh idempotency key instead of attaching to the
  // completed one. Without it a re-read would return the same cached answer.
  let nonce = 0;
  let requested = false;

  // Collapsed the note is clamped to a few lines with no scrollbar, which is what keeps a dashboard row
  // from turning into a scroll container by default. Expanding does not grow the box (canvas rows are a
  // fixed height) but hands the overflow back to the reader as a scroll.
  let expanded = false;

  $: specStore = component?.specStore;
  $: spec = specStore ? $specStore : undefined;
  $: agent = spec?.agent ?? "";
  $: prompt = spec?.prompt ?? "";
  $: autoRun = spec?.auto_run !== false;
  $: configured = agent.trim().length > 0 && prompt.trim().length > 0;

  // The state the reader is looking at, built from the same canvas stores the AI-Chat reads (see
  // features/canvas/chat-context.ts). The parent entity is to hand here, so there is no store lookup by name.
  $: filterMapStore = component?.parent?.filterManager?.filterMapStore;
  $: intervalStore = component?.parent?.timeManager?.state?.interval;
  $: filterMap = filterMapStore ? $filterMapStore : undefined;
  $: interval = intervalStore ? $intervalStore : undefined;

  $: dashboardContext = ((): V1AnalystAgentContext | undefined => {
    const ctx: V1AnalystAgentContext = {};
    if (component?.parent?.name) ctx.canvas = component.parent.name;
    if (interval?.isValid) {
      ctx.timeStart = interval.start.toUTC().toISO() ?? undefined;
      ctx.timeEnd = interval.end.toUTC().toISO() ?? undefined;
    }
    if (filterMap?.size) {
      const where: NonNullable<V1AnalystAgentContext["wherePerMetricsView"]> =
        {};
      filterMap.forEach((expr, metricsView) => {
        if (expr?.cond?.exprs?.length) where[metricsView] = expr;
      });
      if (Object.keys(where).length) ctx.wherePerMetricsView = where;
    }
    return Object.keys(ctx).length ? ctx : undefined;
  })();

  $: idempotencyKey = agentNoteIdempotencyKey({
    agent,
    prompt,
    context: `${runtimeClient.instanceId ?? ""}|${dashboardContextKey(dashboardContext)}`,
    nonce,
  });

  $: queryOptions = {
    queryKey: ["agent-note", runtimeClient.instanceId, idempotencyKey],
    queryFn: ({ signal }: { signal: AbortSignal }) =>
      runAgentNote(
        runtimeClient,
        { agent, prompt, idempotencyKey, dashboardContext },
        signal,
      ),
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
    expanded = false;
  }

  // A fresh answer starts collapsed, and a note that fits needs no affordance at all: compare the rendered
  // height against the clamp once the DOM has it.
  let noteEl: HTMLDivElement | undefined;
  let overflows = false;
  $: if (noteEl && note) {
    void note;
    overflows = noteEl.scrollHeight - noteEl.clientHeight > 4;
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
    <div
      bind:this={noteEl}
      class="agent-note select-text cursor-text"
      class:is-collapsed={!expanded}
    >
      {#await renderPromise then html}
        {@html DOMPurify.sanitize(html)}
      {/await}
    </div>
    <div class="agent-note-footer">
      {#if overflows || expanded}
        <button type="button" on:click={() => (expanded = !expanded)}>
          {expanded
            ? m.canvas_agent_note_show_less()
            : m.canvas_agent_note_show_more()}
        </button>
      {/if}
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
    @apply text-fg-primary flex-1 min-h-0 overflow-y-auto;
  }
  /* line-clamp gives a clean cut at a line boundary and suppresses the scrollbar, so the collapsed note
     never looks like a scroll container. Expanding swaps it for the overflow the class was hiding. */
  .agent-note.is-collapsed {
    @apply overflow-hidden;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 3;
    line-clamp: 3;
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
    @apply size-2 rounded-full bg-accent-primary;
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
    @apply pt-1 flex gap-x-3 shrink-0;
  }
  button {
    @apply text-xs text-accent-primary-action hover:underline;
  }
</style>

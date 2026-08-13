<script lang="ts">
  import RefreshIcon from "@rilldata/web-common/components/icons/RefreshIcon.svelte";
  import Markdown from "@rilldata/web-common/components/markdown/Markdown.svelte";
  import { extractErrorMessage } from "@rilldata/web-common/lib/errors";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";
  import type { V1AnalystAgentContext } from "@rilldata/web-common/runtime-client";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import { createQuery } from "@tanstack/svelte-query";
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
  $: errorMessage = $noteQuery?.error
    ? extractErrorMessage($noteQuery.error)
    : "";

  function regenerate() {
    nonce += 1;
    requested = true;
  }
</script>

<div
  class="agent-note-root group size-full flex flex-col bg-surface-card px-3 py-2 relative"
>
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
    <!-- Scrolls, like the markdown component: a canvas row is a fixed height and Rill's own text
         component does not dress that up. How much there is to read is governed where it belongs, in the
         agent's instructions and the row's height. -->
    <!-- The shared markdown renderer, the same one the chat uses for an assistant's reply. This component
         does not parse or sanitise anything itself; it owns the box and scrolls, like Rill's markdown
         component does, since a canvas row is a fixed height. -->
    <div class="agent-note select-text cursor-text">
      <Markdown content={note} />
    </div>
    <button
      type="button"
      class="agent-note-regenerate"
      title={m.canvas_agent_note_regenerate()}
      aria-label={m.canvas_agent_note_regenerate()}
      on:click={regenerate}
    >
      <RefreshIcon size="14px" />
    </button>
  {:else}
    <div class="agent-note-idle">
      <button type="button" on:click={regenerate}>
        {m.canvas_agent_note_generate()}
      </button>
    </div>
  {/if}
</div>

<style lang="postcss">
  /* Typography belongs to the shared markdown component; this only owns the box. The right gutter keeps
     the hover icon off the first line. */
  .agent-note {
    @apply flex-1 min-h-0 overflow-y-auto pr-7;
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
  /* Same affordance the canvas toolbars use: absent until the component is hovered, so a note at rest is
     only its text. Focus reveals it too, or it would be unreachable by keyboard. */
  .agent-note-regenerate {
    @apply absolute top-1.5 right-2 p-1 rounded text-fg-secondary opacity-0 transition-opacity;
    @apply hover:text-fg-primary hover:bg-surface-subtle;
  }
  .agent-note-root:hover .agent-note-regenerate,
  .agent-note-regenerate:focus-visible {
    @apply opacity-100;
  }

  .agent-note-error button {
    @apply text-xs text-accent-primary-action hover:underline;
  }
  .agent-note-idle button {
    @apply text-xs text-accent-primary-action hover:underline;
  }
</style>

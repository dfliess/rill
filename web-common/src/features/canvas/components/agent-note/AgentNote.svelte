<script lang="ts">
  import AnimatedDots from "@rilldata/web-common/features/chat/core/messages/AnimatedDots.svelte";
  import RefreshIcon from "@rilldata/web-common/components/icons/RefreshIcon.svelte";
  import Markdown from "@rilldata/web-common/components/markdown/Markdown.svelte";
  import { extractErrorMessage } from "@rilldata/web-common/lib/errors";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import { runtimeServiceGetResource } from "@rilldata/web-common/runtime-client/v2/gen/runtime-service";
  import { ResourceKind } from "@rilldata/web-common/features/entity-management/resource-selectors";
  import { createQuery } from "@tanstack/svelte-query";
  import type { AgentNoteCanvasComponent } from "./";
  import {
    agentNoteIdempotencyKey,
    buildDashboardContext,
    dashboardContextKey,
    noteContextReady,
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
  // Upstream's unified filter component replaced the canvas' `filterManager` with
  // `expressionFilterManager`, whose store keys expressions by metrics view in a plain
  // record rather than a Map. `buildDashboardContext` still takes a Map, so convert here
  // and leave the note's own helper (and its tests) alone.
  $: filterMapStore =
    component?.parent?.expressionFilterManager?.exprByMetricsViewStore;
  $: filterMap = filterMapStore
    ? new Map(Object.entries($filterMapStore))
    : undefined;

  // `time_filters` pins the note to its own window, which is what lets a weekly note sit under monthly
  // cards. The canvas resolves that into a local time state per component, but its usual reader,
  // BaseCanvasComponent's timeAndFilterStore, returns the canvas window before it ever consults the local
  // one unless the component declares a `metrics_view` (see BaseCanvasComponent.ts:271). A note narrates a
  // whole dashboard and so names no single metrics view, hence the choice is made here, on the same
  // condition that component uses.
  $: globalIntervalStore = component?.parent?.timeManager?.state?.interval;
  $: localIntervalStore = component?.localTimeControls?.interval;
  $: globalInterval = globalIntervalStore ? $globalIntervalStore : undefined;
  $: localInterval = localIntervalStore ? $localIntervalStore : undefined;
  $: interval = spec?.time_filters ? localInterval : globalInterval;

  // WAIT FOR THE CONTEXT BEFORE ASKING (issue #155). The canvas window does not exist until two round
  // trips have landed: the canvas spec (which fills hasTimeSeriesMap, and until then rangeStore returns
  // undefined — stores/time-state.ts:91), and then deriveInterval's own query that turns "P12M" into
  // dates. A note that does not wait composes its key WITHOUT the window, runs, and then runs a second
  // time under the key the window produces: two completions per load, and the answer left on screen is
  // the one computed without the reader's time range.
  //
  // The shape is the markdown component's (markdown/util.ts:167-172), which narrates the same dashboards
  // and has the same "with or without time" requirement. Both halves are needed, one per round trip: the
  // time half alone would let the first render through, because before the spec the map is EMPTY, so
  // hasTimeSeries reads false and "no time series" is indistinguishable from "not yet".
  $: canvasSpecStore = component?.parent?.specStore;
  $: canvasData = canvasSpecStore ? $canvasSpecStore?.data : undefined;
  $: hasTimeSeriesStore = component?.parent?.timeManager?.hasTimeSeriesMap;
  $: hasTimeSeriesMap = hasTimeSeriesStore ? $hasTimeSeriesStore : undefined;

  $: contextReady = noteContextReady({
    canvasData,
    hasTimeSeriesMap,
    interval,
  });

  $: dashboardContext = buildDashboardContext({
    canvas: component?.parent?.name,
    interval,
    filterMap,
  });

  // The agent's version belongs in the key. A note is an answer about (dashboard, agent), so editing the agent has
  // to produce a new answer instead of serving the one its previous version wrote. It also gives a note whose agent
  // was misconfigured a way out: fixing the agent moves the key, where a retry of the old one could only fail again.
  // Rill already publishes the version as the resource's spec hash, so reading it costs nothing new in the API.
  $: agentSpecQuery = createQuery(
    {
      queryKey: ["agent-note-spec", runtimeClient.instanceId, agent],
      queryFn: ({ signal }: { signal: AbortSignal }) =>
        runtimeServiceGetResource(
          runtimeClient,
          { name: { kind: ResourceKind.Agent, name: agent } },
          { signal },
        ),
      enabled: !!agent,
      retry: false,
    },
    queryClient,
  );

  // The cast is the same type-level formality the dashboard context needs in util.ts: the v2 client returns proto
  // JSON typed as the REST client's `V1GetResourceResponse`, whose generated `V1Resource` predates the Agent
  // resource and so does not declare the field. The bytes carry it; only the type has not caught up.
  $: agentSpecHash =
    (
      $agentSpecQuery?.data?.resource as
        | { agent?: { state?: { specHash?: string } } }
        | undefined
    )?.agent?.state?.specHash ?? "";

  $: idempotencyKey = agentNoteIdempotencyKey({
    agent,
    prompt,
    context: `${runtimeClient.instanceId ?? ""}|${agentSpecHash}|${dashboardContextKey(dashboardContext)}`,
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
    // Wait for the agent's version before starting anything: generating under an empty hash and again under the
    // real one would pay for two completions and leave the first orphaned under a key nothing asks for again.
    enabled:
      configured &&
      contextReady &&
      $agentSpecQuery?.isSuccess &&
      (autoRun || requested),
    // The answer is pinned to its idempotency key, so it never goes stale on its own: a new answer
    // requires a new key, which only "regenerate" or an edited prompt produces.
    staleTime: Infinity,
    gcTime: Infinity,
    retry: false,
    refetchOnWindowFocus: false,
  };
  $: noteQuery = createQuery(queryOptions, queryClient);

  $: note = $noteQuery?.data ?? "";

  // THE DETAIL GOES TO THE CONSOLE, NEVER TO THE PANEL. A failed run's error
  // arrives verbatim from the runtime, and for a provider failure that string
  // carries a full Go stack trace: `runtime/ai/ai.go` wraps the error as
  // "completion failed: %w (stack: %s)" with `debug.Stack()` inside. Rendering
  // it put ~40 lines of goroutine dump, with absolute source paths from the
  // server that built the binary, into a dashboard panel. Notes are read by
  // business users and are shareable by public link, so the reader gets one
  // sentence and a retry; whoever is debugging opens the console.
  $: if ($noteQuery?.isError) {
    console.error(
      "[agent_note] run failed:",
      extractErrorMessage($noteQuery.error),
    );
  }

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
      <AnimatedDots>{m.canvas_agent_note_generating()}</AnimatedDots>
    </div>
  {:else if $noteQuery?.isError}
    <div class="agent-note-error">
      <p>{m.canvas_agent_note_failed()}</p>
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
  /* The waiting state is Rill's shared one (AnimatedDots), the same affordance the chat's thinking block
     and the canvas's other AI states use, so the trailing dots are animated and never written into the
     string. */
  .agent-note-loading {
    @apply text-sm text-fg-muted;
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

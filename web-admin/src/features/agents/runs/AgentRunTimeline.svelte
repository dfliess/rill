<script lang="ts">
  import { Tag } from "@rilldata/web-common/components/tag";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import { agentRunEventsStore } from "../agent-run-events";
  import {
    agentRunStatusColor,
    agentRunStatusLabel,
    formatDateTime,
    formatJson,
  } from "../utils";

  let { runId }: { runId: string } = $props();

  const runtimeClient = useRuntimeClient();

  // Re-open the stream when the run changes.
  let events = $derived(agentRunEventsStore(runtimeClient, runId));
</script>

<div class="flex flex-col gap-y-3">
  <div class="flex items-center gap-x-2">
    <h2 class="text-fg-primary text-base font-semibold">
      {m.agents_timeline_title()}
    </h2>
    {#if $events.streaming}
      <Tag color="green">{m.agents_timeline_live()}</Tag>
    {/if}
  </div>

  {#if $events.error}
    <!-- The live event stream could not be consumed; the run status above still
         reflects the run's lifecycle via polling. -->
    <div
      class="text-sm text-fg-secondary border border-amber-300 dark:border-amber-400/30 bg-amber-50 dark:bg-amber-500/10 rounded p-3"
    >
      {m.agents_timeline_stream_unavailable()}
      <div class="text-xs text-fg-muted mt-1">{$events.error}</div>
    </div>
  {:else if $events.events.length === 0}
    <div class="text-sm text-fg-secondary">{m.agents_timeline_no_events()}</div>
  {:else}
    <ol class="flex flex-col gap-y-2">
      {#each $events.events as event (event.id)}
        <li class="border rounded p-3 flex flex-col gap-y-1">
          <div class="flex items-center gap-x-2 flex-wrap">
            <span class="text-xs font-mono text-fg-muted">#{event.seq}</span>
            <span class="text-sm font-medium text-fg-primary">
              {event.eventType || m.agents_timeline_event_fallback()}
            </span>
            {#if event.status}
              <Tag color={agentRunStatusColor(event.status)}>
                {agentRunStatusLabel(event.status)}
              </Tag>
            {/if}
            <div class="grow"></div>
            <span class="text-xs text-fg-secondary">
              {formatDateTime(event.createdOn)}
            </span>
          </div>
          {#if event.payload !== undefined}
            <details>
              <summary class="text-xs text-fg-secondary cursor-pointer">
                {m.agents_timeline_payload()}
              </summary>
              <pre
                class="text-xs whitespace-pre-wrap bg-surface-secondary rounded p-2 mt-1 overflow-x-auto">{formatJson(
                  event.payload,
                )}</pre>
            </details>
          {/if}
        </li>
      {/each}
    </ol>
  {/if}
</div>

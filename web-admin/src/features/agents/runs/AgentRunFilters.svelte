<script lang="ts">
  import { Search } from "@rilldata/web-common/components/search";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import type { UrlParamsState } from "web-common/src/lib/store-utils/url-params-state.svelte.ts";
  import { agentRunStatusLabel } from "../utils";

  // The list toolbar: the shared Search input (white, like Alerts/Reports) plus
  // the two server-supported filters (ListAgentRuns accepts agent name and status).
  // Each filter is a URL-backed param so the list is shareable/bookmarkable; the
  // inputs bind through the param's getter/setter (Svelte 5.19+ function binding).
  let {
    agents = [],
    searchParam,
    agentParam,
    statusParam,
  }: {
    agents?: string[];
    searchParam: UrlParamsState<string, string>;
    agentParam: UrlParamsState<string, null>;
    statusParam: UrlParamsState<string, string>;
  } = $props();

  // Default the inbox to the runs that need a human: those waiting on approval.
  const RUN_STATUSES = [
    "waiting_approval",
    "running",
    "queued",
    "succeeded",
    "failed",
    "cancelled",
  ];

  const selectClass =
    "border rounded-lg px-2 py-1.5 text-sm bg-surface-background text-fg-primary h-8";
</script>

<div class="flex flex-wrap gap-x-3 gap-y-2 items-center">
  <div class="min-w-[180px] grow max-w-sm">
    <Search
      placeholder={m.agents_search_placeholder()}
      autofocus={false}
      bind:value={searchParam.getter, searchParam.setter}
      rounded="lg"
      retainValueOnMount
    />
  </div>

  <select
    class={selectClass}
    aria-label={m.agents_field_agent()}
    bind:value={agentParam.getter, agentParam.setter}
  >
    <option value="">{m.agents_field_agent()}: {m.agents_filter_all()}</option>
    {#each agents as agent (agent)}
      <option value={agent}>{agent}</option>
    {/each}
  </select>

  <select
    class={selectClass}
    aria-label={m.agents_field_status()}
    bind:value={statusParam.getter, statusParam.setter}
  >
    <!-- "all" (not "") so the choice differs from the "waiting_approval" default
         and persists in the URL. -->
    <option value="all"
      >{m.agents_field_status()}: {m.agents_filter_all()}</option
    >
    {#each RUN_STATUSES as s (s)}
      <option value={s}>{agentRunStatusLabel(s)}</option>
    {/each}
  </select>
</div>

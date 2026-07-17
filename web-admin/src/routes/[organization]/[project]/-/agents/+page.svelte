<script lang="ts">
  import { page } from "$app/state";
  import ContentContainer from "@rilldata/web-common/components/layout/ContentContainer.svelte";
  import DelayedSpinner from "@rilldata/web-common/features/entity-management/DelayedSpinner.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import AgentRunFilters from "@rilldata/web-admin/features/agents/runs/AgentRunFilters.svelte";
  import AgentRunsTable from "@rilldata/web-admin/features/agents/runs/AgentRunsTable.svelte";
  import StartAgentRunDialog from "@rilldata/web-admin/features/agents/runs/StartAgentRunDialog.svelte";
  import {
    useAgentApprovals,
    useAgentRunsReactive,
    useAgents,
  } from "@rilldata/web-admin/features/agents/selectors";
  import type { AgentApprovalData } from "@rilldata/web-admin/features/agents/types";
  import { UrlParamsState } from "web-common/src/lib/store-utils/url-params-state.svelte.ts";
  import { writable } from "svelte/store";

  const runtimeClient = useRuntimeClient();

  let { organization, project } = $derived(page.params);

  // The 3 filters live in the URL so the list is shareable and bookmarkable.
  // `q` feeds the shared <Search> input, whose `value` is typed `string | number`
  // and rejects `null`, so it is built with a non-null (`string`) default rather
  // than `createStringParam` (whose getter is `string | null`). Same URL semantics
  // as `createStringParam`: an empty value drops the param.
  const searchParam = new UrlParamsState<string, string>(
    "q",
    (v) => (v === "" ? null : v),
    (v) => v ?? "",
    "",
  );
  const agentParam = UrlParamsState.createStringParam("agent");
  // The inbox opens on the runs that need a human: those waiting on approval.
  // A custom serializer keeps that default out of the URL and lets "all" (no
  // server-side status filter) persist.
  const statusParam = new UrlParamsState<string, string>(
    "status",
    (v) => (v === "waiting_approval" ? null : v),
    (v) => v ?? "waiting_approval",
    "waiting_approval",
  );

  const agentsQuery = useAgents(runtimeClient);
  let agentNames = $derived(
    ($agentsQuery.data?.agents ?? []).map((a) => a.name ?? "").filter(Boolean),
  );

  // Agent + status are server-side filters, so bridge the runes filter state into
  // a store the query subscribes to; the search box narrows client-side. Seed the
  // store with the current filters so the first fetch already respects them.
  const runsRequest = writable<{ agentName?: string; status?: string }>({
    agentName: agentParam.value || undefined,
    status: statusParam.value === "all" ? undefined : statusParam.value,
  });
  $effect(() => {
    runsRequest.set({
      agentName: agentParam.value || undefined,
      status: statusParam.value === "all" ? undefined : statusParam.value,
    });
  });
  const runsQuery = useAgentRunsReactive(runtimeClient, runsRequest);

  let term = $derived((searchParam.value ?? "").trim().toLowerCase());
  let runs = $derived(
    ($runsQuery.data?.runs ?? []).filter((r) => {
      if (!term) return true;
      return [r.agentName, r.runId, r.trigger]
        .filter(Boolean)
        .some((v) => v!.toLowerCase().includes(term));
    }),
  );

  // Pending approvals joined to their run, so a waiting run offers approve/deny
  // inline. A run does not carry its approval, so we fetch them separately.
  const approvalsQuery = useAgentApprovals(runtimeClient, {
    status: "pending",
  });
  let approvalsByRun = $derived(
    ($approvalsQuery.data?.approvals ?? []).reduce(
      (map, a: AgentApprovalData) => {
        if (a.runId && !map.has(a.runId)) map.set(a.runId, a);
        return map;
      },
      new Map<string, AgentApprovalData>(),
    ),
  );
</script>

<ContentContainer title={m.agents_page_title()} maxWidth={1200}>
  <div class="flex flex-col gap-y-6">
    <div class="flex flex-wrap items-center justify-between gap-x-4 gap-y-2">
      <p class="text-fg-secondary text-sm">{m.agents_page_description()}</p>
      <StartAgentRunDialog agents={agentNames} {organization} {project} />
    </div>

    <AgentRunFilters
      agents={agentNames}
      {searchParam}
      {agentParam}
      {statusParam}
    />

    {#if $runsQuery.isLoading}
      <div class="m-auto mt-20">
        <DelayedSpinner isLoading size="24px" />
      </div>
    {:else if $runsQuery.isError}
      <div class="text-sm text-red-600">
        {$runsQuery.error?.message ?? m.agents_runs_load_error()}
      </div>
    {:else}
      <AgentRunsTable data={runs} {approvalsByRun} {organization} {project} />
    {/if}
  </div>
</ContentContainer>

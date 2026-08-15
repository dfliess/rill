<script lang="ts">
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import AgentRunsTable from "./runs/AgentRunsTable.svelte";
  import { useAgentApprovals, useAgentRuns } from "./selectors";
  import type { AgentApprovalData } from "./types";

  let { organization, project }: { organization: string; project: string } =
    $props();

  const runtimeClient = useRuntimeClient();

  // Home is a passive landing surface, often left open: poll gently instead of
  // the Act tab's live-operations 5s cadence.
  const HOME_REFETCH_INTERVAL = 15_000;

  // The Home inbox: runs blocked on a human decision. No search or filters here —
  // it is a focused "needs you" list; the Act tab is where you filter and dig in.
  const runsQuery = useAgentRuns(
    runtimeClient,
    { status: "waiting_approval" },
    HOME_REFETCH_INTERVAL,
  );
  let runs = $derived($runsQuery.data?.runs ?? []);

  // Runs blocked on a pending approval. A run waits until a human decides or it is
  // cancelled, so every pending approval is actionable and belongs in the inbox. A
  // batch turn can leave several approvals pending on one run, so the map keeps them
  // all: the row acts on the first and shows how many more are behind it.
  const approvalsQuery = useAgentApprovals(
    runtimeClient,
    { status: "pending" },
    HOME_REFETCH_INTERVAL,
  );
  let approvalsByRun = $derived(
    ($approvalsQuery.data?.approvals ?? []).reduce(
      (map, a: AgentApprovalData) => {
        if (a.runId) map.set(a.runId, [...(map.get(a.runId) ?? []), a]);
        return map;
      },
      new Map<string, AgentApprovalData[]>(),
    ),
  );
  // Home only surfaces approvals the viewer can actually resolve: an inbox of
  // decisions that are not theirs to make is noise, not a call to action. The
  // server resolves can_decide per approval (authority can discriminate by
  // tool); the full (read-only) picture stays available in the Act tab.
  let pendingRuns = $derived(
    runs.filter((r) =>
      (approvalsByRun.get(r.runId ?? "") ?? []).some((a) => !!a.canDecide),
    ),
  );

  let base = $derived(`/${organization}/${project}/-/agents`);
</script>

<!-- Hidden entirely when nothing is waiting, so Home stays quiet until it needs attention. -->
{#if pendingRuns.length > 0}
  <div class="flex flex-col gap-y-4">
    <h2
      class="flex items-center justify-between text-xl font-semibold text-fg-secondary"
    >
      {m.agents_home_pending_heading()}
      <a class="text-sm font-normal text-primary-600" href={base}>
        {m.agents_home_view_all()}
      </a>
    </h2>
    <AgentRunsTable
      data={pendingRuns}
      {approvalsByRun}
      {organization}
      {project}
    />
  </div>
{/if}

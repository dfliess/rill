<script lang="ts">
  import Brain from "@rilldata/web-common/components/icons/Brain.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { timeAgo } from "@rilldata/web-common/lib/time/relative-time";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import HomeBox from "@rilldata/web-admin/features/kairos-home/HomeBox.svelte";
  import ApproveDenyButtons from "./approvals/ApproveDenyButtons.svelte";
  import { runDetailPath } from "./run-id";
  import { useAgentApprovals, useAgentRuns, useAgents } from "./selectors";
  import type { AgentApprovalData, AgentRunData } from "./types";
  import { approvalGlance, formatDateTime } from "./utils";

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

  // Both lists must have answered before "nothing pending" is a fact rather
  // than a loading state.
  let loaded = $derived($runsQuery.isSuccess && $approvalsQuery.isSuccess);

  let base = $derived(`/${organization}/${project}/-/agents`);

  // A glance, not the list: the header counts them all, "View all" has the rest.
  let shown = $derived(pendingRuns.slice(0, 4));

  // Runs carry the agent's slug; the catalog carries its display name.
  const agentsQuery = useAgents(runtimeClient);
  let displayNameByAgent = $derived(
    new Map<string, string>(
      ($agentsQuery.data?.agents ?? []).map((a) => [
        a.name ?? "",
        a.displayName ?? "",
      ]),
    ),
  );

  function startedAgo(run: AgentRunData): string {
    const raw = run.startedOn ?? run.createdOn;
    const date = typeof raw === "string" ? new Date(raw) : undefined;
    return date && !Number.isNaN(date.getTime()) ? timeAgo(date) : "";
  }
</script>

<HomeBox
  title={m.agents_home_pending_heading()}
  count={pendingRuns.length}
  href={base}
  viewAllLabel={m.agents_home_view_all()}
  emptyText={m.agents_home_pending_empty()}
  {loaded}
>
  {#each shown as run (run.runId)}
    {@const approvals = approvalsByRun.get(run.runId ?? "") ?? []}
    <!-- can_decide is per approval, so the row names and acts on the first one
         this viewer may decide (see AgentRunsTable). -->
    {@const decidable = approvals.find((a) => !!a.canDecide)}
    {@const glance = decidable ? approvalGlance(decidable) : ""}
    {@const ago = startedAgo(run)}
    <li class="flex items-center gap-x-3 px-4 py-2.5 hover:bg-surface-hover">
      <span class="shrink-0 text-fg-secondary"><Brain size="15px" /></span>
      <!-- Approve happens in place (#231): the row says what would run (the
           signed argument values, escaped), the server still binds the decision
           to the args hash, and the full arguments are one click away. -->
      <a
        href={runDetailPath(organization, project, run)}
        class="flex min-w-0 grow flex-col"
        title={formatDateTime(run.startedOn ?? run.createdOn)}
      >
        <!-- When on the name's line: the arguments below can run long and
             would push it out of sight. -->
        <span class="flex min-w-0 items-baseline gap-x-2">
          <span class="truncate text-sm font-medium text-fg-primary">
            {displayNameByAgent.get(run.agentName ?? "") || run.agentName}
          </span>
          {#if ago}
            <span class="shrink-0 whitespace-nowrap text-xs text-fg-muted"
              >{ago}</span
            >
          {/if}
        </span>
        <span class="truncate text-xs text-fg-secondary">
          {glance}
          {#if approvals.length > 1}
            · {m.agents_approval_and_n_more({ count: approvals.length - 1 })}
          {/if}
        </span>
      </a>
      {#if decidable}
        <div class="shrink-0">
          <ApproveDenyButtons
            approvalId={decidable.approvalId ?? ""}
            argsHash={decidable.argsHash ?? ""}
            runId={run.runId ?? ""}
            small
          />
        </div>
      {/if}
    </li>
  {/each}
</HomeBox>

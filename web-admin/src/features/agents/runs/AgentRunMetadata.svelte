<script lang="ts">
  import MetadataLabel from "@rilldata/web-admin/features/scheduled-reports/metadata/MetadataLabel.svelte";
  import MetadataValue from "@rilldata/web-admin/features/scheduled-reports/metadata/MetadataValue.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import AgentApprovalCard from "../approvals/AgentApprovalCard.svelte";
  import {
    useAgentApprovals,
    useAgentRun,
    useSubjectNames,
  } from "../selectors";
  import {
    agentTriggerLabel,
    formatDateTime,
    isApprovalPending,
    runActorLabel,
    timestampSortKey,
  } from "../utils";
  import AgentRunStatusChip from "./AgentRunStatusChip.svelte";
  import AgentRunTimeline from "./AgentRunTimeline.svelte";
  import CancelAgentRunButton from "./CancelAgentRunButton.svelte";
  import StartAgentRunDialog from "./StartAgentRunDialog.svelte";

  let {
    organization,
    project,
    runId,
  }: { organization: string; project: string; runId: string } = $props();

  const runtimeClient = useRuntimeClient();

  // These queries key off `runId`, which only changes on navigation to another
  // run. The route remounts this component per run (`{#key run}`), so the query
  // hooks are created once at init with the current id rather than reactively.
  // svelte-ignore state_referenced_locally
  const runQuery = useAgentRun(runtimeClient, runId);
  let run = $derived($runQuery.data?.run);

  // A run's proposed action lives on an approval, not on the run. Fetch the run's
  // approvals so the reviewer can approve/deny the proposal here (folded in from
  // the old standalone approval page). Resolved ones stay on screen: they are the
  // run's audit trail (who decided what, and when).
  // svelte-ignore state_referenced_locally
  const approvalsQuery = useAgentApprovals(runtimeClient, { runId });
  // Pending first (they need a decision), then the resolved ones oldest-first, so a
  // multi-action run reads in the order it happened.
  let approvals = $derived(
    [...($approvalsQuery.data?.approvals ?? [])].sort((a, b) => {
      const aPending = isApprovalPending(a.status);
      if (aPending !== isApprovalPending(b.status)) return aPending ? -1 : 1;
      return timestampSortKey(a.createdOn).localeCompare(
        timestampSortKey(b.createdOn),
      );
    }),
  );

  // Names the run's actor and each approval's decider, instead of a raw user id.
  // svelte-ignore state_referenced_locally
  const subjectNamesQuery = useSubjectNames(organization, project);
  let subjectNames = $derived($subjectNamesQuery.data);
</script>

{#if run}
  <div class="flex flex-col gap-y-9 w-full max-w-full 2xl:max-w-[1200px]">
    <div class="flex flex-col gap-y-2">
      <div class="flex gap-x-2 items-center flex-wrap">
        <h1 class="text-fg-primary text-lg font-bold">{run.agentName}</h1>
        <AgentRunStatusChip status={run.status} />
        <div class="grow"></div>
        <!-- Re-run this agent with a fresh manual prompt. -->
        <StartAgentRunDialog
          agent={run.agentName ?? ""}
          {organization}
          {project}
          label={m.agents_run_new_manual()}
        />
        <!-- An alert-triggered run records the alert that fired it; link back to it. -->
        {#if run.trigger === "alert" && run.triggerRef}
          <a
            class="text-primary-600 text-sm"
            href={`/${organization}/${project}/-/alerts/${run.triggerRef}`}
          >
            {m.agents_run_view_alert()}
          </a>
        {/if}
        <!-- A run executes inside an AI session (Conversation); link to it so the run can be read as a chat. -->
        {#if run.conversationId}
          <a
            class="text-primary-600 text-sm"
            href={`/${organization}/${project}/-/ai/${run.conversationId}`}
          >
            {m.agents_run_view_conversation()}
          </a>
        {/if}
        <CancelAgentRunButton {runId} {run} />
      </div>
      <div class="text-fg-secondary text-xs font-mono">{run.runId}</div>
    </div>

    <!-- Each side-effecting action the run proposed: pending ones await a decision,
         resolved ones record who decided and when. -->
    {#if approvals.length > 0}
      <div class="flex flex-col gap-y-4">
        {#each approvals as approval (approval.approvalId)}
          <AgentApprovalCard {approval} names={subjectNames} />
        {/each}
      </div>
    {/if}

    <div class="flex flex-wrap gap-x-16 gap-y-6">
      <div class="flex flex-col gap-y-3">
        <MetadataLabel>{m.agents_field_trigger()}</MetadataLabel>
        <MetadataValue>{agentTriggerLabel(run.trigger)}</MetadataValue>
      </div>
      <div class="flex flex-col gap-y-3">
        <!-- The identity the run (and any approved action) executes as. -->
        <MetadataLabel>{m.agents_run_runs_as()}</MetadataLabel>
        <MetadataValue>{runActorLabel(run, subjectNames)}</MetadataValue>
      </div>
      <div class="flex flex-col gap-y-3">
        <MetadataLabel>{m.agents_run_spec_version()}</MetadataLabel>
        <MetadataValue>
          <span class="font-mono text-xs">{run.specHash || "—"}</span>
        </MetadataValue>
      </div>
      <div class="flex flex-col gap-y-3">
        <MetadataLabel>{m.agents_field_started()}</MetadataLabel>
        <MetadataValue
          >{formatDateTime(run.startedOn ?? run.createdOn)}</MetadataValue
        >
      </div>
      <div class="flex flex-col gap-y-3">
        <MetadataLabel>{m.agents_field_finished()}</MetadataLabel>
        <MetadataValue>{formatDateTime(run.finishedOn)}</MetadataValue>
      </div>
    </div>

    {#if run.error}
      <div
        class="text-sm text-red-600 dark:text-red-300 border border-red-300 dark:border-red-400/30 bg-red-50 dark:bg-red-500/10 rounded p-3 whitespace-pre-wrap"
      >
        {run.error}
      </div>
    {/if}

    <AgentRunTimeline {runId} names={subjectNames} />
  </div>
{:else if $runQuery.isError}
  <div class="text-sm text-red-600">
    {$runQuery.error?.message ?? m.agents_run_load_error()}
  </div>
{/if}

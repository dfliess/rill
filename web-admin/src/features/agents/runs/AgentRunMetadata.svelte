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
  import * as Tooltip from "@rilldata/web-common/components/tooltip-v2";
  import AgentRunStatusChip from "./AgentRunStatusChip.svelte";
  import AgentRunTimeline from "./AgentRunTimeline.svelte";
  import CancelAgentRunButton from "./CancelAgentRunButton.svelte";
  import StartAgentRunDialog from "./StartAgentRunDialog.svelte";

  let runIdCopied = $state(false);
  async function copyRunId() {
    if (!run?.runId) return;
    await navigator.clipboard.writeText(run.runId);
    runIdCopied = true;
    setTimeout(() => (runIdCopied = false), 1500);
  }

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
      <!-- Run ID hidden from the first plane; available via copy for support. -->
      <div class="flex items-center gap-x-1">
        <Tooltip.Root>
          <Tooltip.Trigger>
            <button
              class="text-fg-muted hover:text-fg-secondary text-xs flex items-center gap-x-1 cursor-pointer"
              onclick={copyRunId}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 16 16"
                fill="currentColor"
                class="size-3.5"
              >
                {#if runIdCopied}
                  <path
                    fill-rule="evenodd"
                    d="M12.416 3.376a.75.75 0 0 1 .208 1.04l-5 7.5a.75.75 0 0 1-1.154.114l-3-3a.75.75 0 0 1 1.06-1.06l2.353 2.353 4.493-6.74a.75.75 0 0 1 1.04-.207Z"
                    clip-rule="evenodd"
                  />
                {:else}
                  <path
                    fill-rule="evenodd"
                    d="M10.986 3H12a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h1.014A2.25 2.25 0 0 1 7.25 1h1.5a2.25 2.25 0 0 1 2.236 2ZM7.25 2.5a.75.75 0 0 0 0 1.5h1.5a.75.75 0 0 0 0-1.5h-1.5Z"
                    clip-rule="evenodd"
                  />
                {/if}
              </svg>
              {runIdCopied ? m.agents_run_id_copied() : m.agents_run_id_copy()}
            </button>
          </Tooltip.Trigger>
          <Tooltip.Content>
            <span class="font-mono text-xs break-all">{run.runId}</span>
          </Tooltip.Content>
        </Tooltip.Root>
      </div>
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

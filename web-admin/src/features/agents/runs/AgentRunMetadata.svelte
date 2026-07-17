<script lang="ts">
  import MetadataLabel from "@rilldata/web-admin/features/scheduled-reports/metadata/MetadataLabel.svelte";
  import MetadataValue from "@rilldata/web-admin/features/scheduled-reports/metadata/MetadataValue.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import AgentApprovalStatusChip from "../approvals/AgentApprovalStatusChip.svelte";
  import ApproveDenyButtons from "../approvals/ApproveDenyButtons.svelte";
  import { useAgentApprovals, useAgentRun } from "../selectors";
  import {
    agentTriggerLabel,
    formatDateTime,
    isApprovalPending,
    parseProposal,
    runActorLabel,
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
  // the old standalone approval page).
  // svelte-ignore state_referenced_locally
  const approvalsQuery = useAgentApprovals(runtimeClient, { runId });
  let pendingApproval = $derived(
    ($approvalsQuery.data?.approvals ?? []).find((a) =>
      isApprovalPending(a.status),
    ),
  );
  let proposal = $derived(parseProposal(pendingApproval?.proposal));
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

    {#if pendingApproval}
      <!-- Proposed side-effecting action awaiting a human decision. -->
      <div
        class="flex flex-col gap-y-4 border border-amber-300 dark:border-amber-400/30 bg-amber-50 dark:bg-amber-500/10 rounded-lg p-4"
      >
        <div class="flex gap-x-2 items-center flex-wrap">
          <h2 class="text-fg-primary text-base font-semibold">
            {m.agents_run_approval_required()}
          </h2>
          <AgentApprovalStatusChip status={pendingApproval.status} />
          <div class="grow"></div>
          <ApproveDenyButtons
            approvalId={pendingApproval.approvalId ?? ""}
            argsHash={pendingApproval.argsHash ?? ""}
          />
        </div>

        <div class="flex flex-wrap gap-x-16 gap-y-4">
          <div class="flex flex-col gap-y-2">
            <MetadataLabel>{m.agents_approval_proposed_action()}</MetadataLabel>
            <MetadataValue>{pendingApproval.toolName || "—"}</MetadataValue>
          </div>
          <div class="flex flex-col gap-y-2">
            <MetadataLabel>{m.agents_field_connector()}</MetadataLabel>
            <MetadataValue>{pendingApproval.connector || "—"}</MetadataValue>
          </div>
          <div class="flex flex-col gap-y-2">
            <MetadataLabel>{m.agents_approval_policy()}</MetadataLabel>
            <MetadataValue>{pendingApproval.policy || "—"}</MetadataValue>
          </div>
          <div class="flex flex-col gap-y-2">
            <MetadataLabel>{m.agents_field_expires()}</MetadataLabel>
            <MetadataValue
              >{formatDateTime(pendingApproval.expiresOn)}</MetadataValue
            >
          </div>
        </div>

        {#if proposal.text}
          <div class="flex flex-col gap-y-2">
            <MetadataLabel
              >{m.agents_approval_proposed_action_args()}</MetadataLabel
            >
            <pre
              class="text-xs text-fg-primary whitespace-pre-wrap bg-surface-secondary rounded p-3 border overflow-x-auto">{proposal.text}</pre>
          </div>
        {/if}
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
        <MetadataValue>{runActorLabel(run)}</MetadataValue>
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

    <AgentRunTimeline {runId} />
  </div>
{:else if $runQuery.isError}
  <div class="text-sm text-red-600">
    {$runQuery.error?.message ?? m.agents_run_load_error()}
  </div>
{/if}

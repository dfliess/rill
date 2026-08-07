<script lang="ts">
  import MetadataLabel from "@rilldata/web-admin/features/scheduled-reports/metadata/MetadataLabel.svelte";
  import MetadataValue from "@rilldata/web-admin/features/scheduled-reports/metadata/MetadataValue.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import type { AgentApprovalData } from "../types";
  import {
    formatDateTime,
    isApprovalPending,
    parseProposal,
    subjectLabel,
  } from "../utils";
  import AgentApprovalStatusChip from "./AgentApprovalStatusChip.svelte";
  import ApproveDenyButtons from "./ApproveDenyButtons.svelte";

  let {
    approval,
    names,
  }: {
    approval: AgentApprovalData;
    /** Subject -> person, to name the decider instead of showing a raw user id. */
    names?: Map<string, string>;
  } = $props();

  let pending = $derived(isApprovalPending(approval.status));
  let proposal = $derived(parseProposal(approval.proposal));
  let showPosition = $derived(
    typeof approval.total === "number" && approval.total > 1,
  );
</script>

<!-- A pending approval is a call to action (amber, with the decision buttons); a
     resolved one is the run's audit record (neutral, no buttons) and stays visible
     so the run keeps showing what was proposed, who decided it and when. -->
<div
  class="flex flex-col gap-y-4 rounded-lg p-4 border {pending
    ? 'border-amber-300 dark:border-amber-400/30 bg-amber-50 dark:bg-amber-500/10'
    : 'bg-surface-subtle'}"
>
  <div class="flex gap-x-2 items-center flex-wrap">
    <h2 class="text-fg-primary text-base font-semibold">
      {pending ? m.agents_run_approval_required() : m.agents_approval_title()}
    </h2>
    <AgentApprovalStatusChip status={approval.status} />
    {#if showPosition}
      <span class="text-xs text-fg-muted">
        {m.agents_approval_position({
          position: approval.position ?? 0,
          total: approval.total ?? 0,
        })}
      </span>
    {/if}
    <div class="grow"></div>
    {#if pending}
      <ApproveDenyButtons
        approvalId={approval.approvalId ?? ""}
        argsHash={approval.argsHash ?? ""}
        runId={approval.runId ?? ""}
      />
    {/if}
  </div>

  <div class="flex flex-wrap gap-x-16 gap-y-4">
    <div class="flex flex-col gap-y-2">
      <MetadataLabel>{m.agents_approval_proposed_action()}</MetadataLabel>
      <MetadataValue>{approval.toolName || "—"}</MetadataValue>
    </div>
    <div class="flex flex-col gap-y-2">
      <MetadataLabel>{m.agents_field_connector()}</MetadataLabel>
      <MetadataValue>{approval.connector || "—"}</MetadataValue>
    </div>
    <div class="flex flex-col gap-y-2">
      <MetadataLabel>{m.agents_approval_policy()}</MetadataLabel>
      <MetadataValue>{approval.policy || "—"}</MetadataValue>
    </div>
    {#if !pending}
      <!-- The audit pair: who decided and when. `decidedBy` is empty for a decision
           no human made (a cancelled approval), so it renders as "—". -->
      <div class="flex flex-col gap-y-2">
        <MetadataLabel>{m.agents_approval_decided_by()}</MetadataLabel>
        <MetadataValue>{subjectLabel(approval.decidedBy, names)}</MetadataValue>
      </div>
      <div class="flex flex-col gap-y-2">
        <MetadataLabel>{m.agents_approval_decided_on()}</MetadataLabel>
        <MetadataValue>{formatDateTime(approval.decidedOn)}</MetadataValue>
      </div>
      <div class="flex flex-col gap-y-2">
        <MetadataLabel>{m.agents_approval_requested_by()}</MetadataLabel>
        <MetadataValue
          >{subjectLabel(approval.requestedBy, names)}</MetadataValue
        >
      </div>
    {/if}
  </div>

  {#if proposal.text}
    <div class="flex flex-col gap-y-2">
      <MetadataLabel>{m.agents_approval_proposed_action_args()}</MetadataLabel>
      <pre
        class="text-xs text-fg-primary whitespace-pre-wrap bg-surface-card rounded p-3 border overflow-x-auto">{proposal.text}</pre>
    </div>
  {/if}
</div>

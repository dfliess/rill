<script lang="ts">
  import { Confirmation } from "@rilldata/web-common/components/alert-dialog";
  import { Button } from "@rilldata/web-common/components/button";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { eventBus } from "@rilldata/web-common/lib/event-bus/event-bus";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import {
    getAgentServiceGetAgentApprovalQueryKey,
    getAgentServiceListAgentApprovalsQueryKey,
    getAgentServiceListAgentRunsQueryKey,
  } from "@rilldata/web-common/runtime-client/v2/gen/agent-service";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { useApproveAgentApproval, useDenyAgentApproval } from "../selectors";

  let {
    approvalId,
    // The hash the approver saw; the server rejects the decision if the proposal
    // changed underneath them.
    argsHash,
  }: { approvalId: string; argsHash: string } = $props();

  const runtimeClient = useRuntimeClient();
  const queryClient = useQueryClient();
  const approve = useApproveAgentApproval(runtimeClient);
  const deny = useDenyAgentApproval(runtimeClient);

  let denyConfirmOpen = $state(false);

  let pending = $derived($approve.isPending || $deny.isPending);

  async function invalidate() {
    await queryClient.invalidateQueries({
      queryKey: getAgentServiceGetAgentApprovalQueryKey(
        runtimeClient.instanceId,
        { approvalId },
      ),
    });
    await queryClient.invalidateQueries({
      queryKey: getAgentServiceListAgentApprovalsQueryKey(
        runtimeClient.instanceId,
      ),
    });
    // A decision flips the run out of waiting_approval, so refresh the run lists
    // (Act tab + Home inbox) too.
    await queryClient.invalidateQueries({
      queryKey: getAgentServiceListAgentRunsQueryKey(runtimeClient.instanceId),
    });
  }

  async function handleApprove() {
    try {
      await $approve.mutateAsync({ approvalId, argsHash });
      await invalidate();
      eventBus.emit("notification", {
        message: m.agents_approval_approved_notification(),
        type: "success",
      });
    } catch (e) {
      // Refresh so a stale/expired approval drops out of the lists on failure.
      await invalidate();
      eventBus.emit("notification", {
        message:
          e instanceof Error ? e.message : m.agents_approval_approve_error(),
        type: "error",
      });
    }
  }

  async function handleDeny() {
    try {
      await $deny.mutateAsync({ approvalId });
      await invalidate();
      eventBus.emit("notification", {
        message: m.agents_approval_denied_notification(),
        type: "success",
      });
    } catch (e) {
      // Refresh so a stale/expired approval drops out of the lists on failure.
      await invalidate();
      eventBus.emit("notification", {
        message:
          e instanceof Error ? e.message : m.agents_approval_deny_error(),
        type: "error",
      });
    }
  }
</script>

<div class="flex gap-x-2">
  <Button type="primary" disabled={pending} onClick={handleApprove}>
    {m.agents_approval_approve()}
  </Button>
  <Button
    type="secondary"
    disabled={pending}
    onClick={() => (denyConfirmOpen = true)}
  >
    {m.agents_approval_deny()}
  </Button>
</div>

<Confirmation
  open={denyConfirmOpen}
  onOpenChange={(open: boolean) => (denyConfirmOpen = open)}
  title={m.agents_approval_deny_confirm_title()}
  description={m.agents_approval_deny_confirm_desc()}
  confirmLabel={m.agents_approval_deny()}
  confirmType="secondary"
  onConfirm={handleDeny}
/>

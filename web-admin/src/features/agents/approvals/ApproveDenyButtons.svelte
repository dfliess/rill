<script lang="ts">
  import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogTitle,
  } from "@rilldata/web-common/components/alert-dialog/index.js";
  import { Button } from "@rilldata/web-common/components/button";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { eventBus } from "@rilldata/web-common/lib/event-bus/event-bus";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import {
    getAgentServiceGetAgentApprovalQueryKey,
    getAgentServiceGetAgentRunQueryKey,
    getAgentServiceListAgentApprovalsQueryKey,
    getAgentServiceListAgentRunsQueryKey,
  } from "@rilldata/web-common/runtime-client/v2/gen/agent-service";
  import { useQueryClient } from "@tanstack/svelte-query";
  import {
    useApproveAgentApproval,
    useCancelAgentRun,
    useDenyAgentApproval,
  } from "../selectors";

  let {
    approvalId,
    // The hash the approver saw; the server rejects the decision if the proposal
    // changed underneath them.
    argsHash,
    runId,
  }: { approvalId: string; argsHash: string; runId: string } = $props();

  const runtimeClient = useRuntimeClient();
  const queryClient = useQueryClient();
  const approve = useApproveAgentApproval(runtimeClient);
  const deny = useDenyAgentApproval(runtimeClient);
  const cancelRun = useCancelAgentRun(runtimeClient);

  let denyConfirmOpen = $state(false);

  let pending = $derived(
    $approve.isPending || $deny.isPending || $cancelRun.isPending,
  );

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
    await queryClient.invalidateQueries({
      queryKey: getAgentServiceListAgentRunsQueryKey(runtimeClient.instanceId),
    });
  }

  async function invalidateRun() {
    await queryClient.invalidateQueries({
      queryKey: getAgentServiceGetAgentRunQueryKey(runtimeClient.instanceId, {
        runId,
      }),
    });
    await queryClient.invalidateQueries({
      queryKey: getAgentServiceListAgentRunsQueryKey(runtimeClient.instanceId),
    });
    await queryClient.invalidateQueries({
      queryKey: getAgentServiceListAgentApprovalsQueryKey(
        runtimeClient.instanceId,
      ),
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
      await invalidate();
      eventBus.emit("notification", {
        message:
          e instanceof Error ? e.message : m.agents_approval_deny_error(),
        type: "error",
      });
    }
  }

  async function handleCancelRun() {
    denyConfirmOpen = false;
    try {
      await $cancelRun.mutateAsync({ runId });
      await invalidateRun();
      eventBus.emit("notification", {
        message: m.agents_run_cancelled_notification(),
        type: "success",
      });
    } catch (e) {
      eventBus.emit("notification", {
        message: e instanceof Error ? e.message : m.agents_run_cancel_error(),
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

<AlertDialog
  open={denyConfirmOpen}
  onOpenChange={(open: boolean) => (denyConfirmOpen = open)}
>
  <AlertDialogContent>
    <AlertDialogTitle>{m.agents_approval_deny_confirm_title()}</AlertDialogTitle
    >
    <AlertDialogDescription>
      {m.agents_approval_deny_confirm_desc()}
    </AlertDialogDescription>
    <AlertDialogFooter>
      <div class="flex w-full items-center">
        <Button
          large
          type="secondary-destructive"
          onClick={handleCancelRun}
          disabled={$cancelRun.isPending}
        >
          {m.agents_approval_deny_cancel_run()}
        </Button>
        <div class="grow"></div>
        <div class="flex gap-x-2">
          <AlertDialogCancel>
            {#snippet child({ props })}
              <Button {...props} large type="secondary">
                {m.agents_approval_deny_close()}
              </Button>
            {/snippet}
          </AlertDialogCancel>
          <AlertDialogAction>
            {#snippet child({ props })}
              <Button {...props} large type="primary" onClick={handleDeny}>
                {m.agents_approval_deny()}
              </Button>
            {/snippet}
          </AlertDialogAction>
        </div>
      </div>
    </AlertDialogFooter>
  </AlertDialogContent>
</AlertDialog>

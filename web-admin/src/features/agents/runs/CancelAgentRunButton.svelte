<script lang="ts">
  import { Button } from "@rilldata/web-common/components/button";
  import { Confirmation } from "@rilldata/web-common/components/alert-dialog";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { eventBus } from "@rilldata/web-common/lib/event-bus/event-bus";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import {
    getAgentServiceGetAgentRunQueryKey,
    getAgentServiceListAgentRunsQueryKey,
  } from "@rilldata/web-common/runtime-client/v2/gen/agent-service";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { useCancelAgentRun } from "../selectors";
  import type { AgentRunData } from "../types";

  let { runId, run }: { runId: string; run: AgentRunData | undefined } =
    $props();

  const runtimeClient = useRuntimeClient();
  const queryClient = useQueryClient();
  const cancelRun = useCancelAgentRun(runtimeClient);

  let confirmOpen = $state(false);

  let finished = $derived(!!run?.finishedOn);
  // The server resolves can_cancel per run (the run's actor, or an EditTrigger
  // holder); the button only renders for callers who can actually use it.
  let canCancel = $derived(!!run?.canCancel);

  async function handleCancel() {
    try {
      await $cancelRun.mutateAsync({ runId });
      await queryClient.invalidateQueries({
        queryKey: getAgentServiceGetAgentRunQueryKey(runtimeClient.instanceId, {
          runId,
        }),
      });
      await queryClient.invalidateQueries({
        queryKey: getAgentServiceListAgentRunsQueryKey(
          runtimeClient.instanceId,
        ),
      });
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

{#if canCancel}
  <Button
    type="secondary"
    disabled={finished || $cancelRun.isPending}
    onClick={() => (confirmOpen = true)}
  >
    {m.agents_run_cancel()}
  </Button>

  <Confirmation
    open={confirmOpen}
    onOpenChange={(open: boolean) => (confirmOpen = open)}
    title={m.agents_run_cancel_confirm_title()}
    description={m.agents_run_cancel_confirm_desc()}
    confirmLabel={m.agents_run_cancel()}
    confirmType="secondary"
    onConfirm={handleCancel}
  />
{/if}

<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button } from "@rilldata/web-common/components/button";
  import {
    Dialog,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
  } from "@rilldata/web-common/components/dialog";
  import Textarea from "@rilldata/web-common/components/forms/Textarea.svelte";
  import { eventBus } from "@rilldata/web-common/lib/event-bus/event-bus";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import { getAgentServiceListAgentRunsQueryKey } from "@rilldata/web-common/runtime-client/v2/gen/agent-service";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { v4 as uuidv4 } from "uuid";
  import { runDetailPath } from "../run-id";
  import { useStartAgentRun } from "../selectors";

  // When `agent` is preset (from a run's detail) the dialog re-runs that agent.
  // When it is empty, `agents` is offered as a picker so a run can be started
  // from the list without a catalog page.
  let {
    agent = "",
    agents = [],
    organization,
    project,
    label = m.agents_run_new(),
  }: {
    agent?: string;
    agents?: string[];
    organization: string;
    project: string;
    label?: string;
  } = $props();

  const runtimeClient = useRuntimeClient();
  const queryClient = useQueryClient();
  const startRun = useStartAgentRun(runtimeClient);

  let open = $state(false);
  let prompt = $state("");
  // Seeded from the preset `agent`; the picker and `reset()` drive it thereafter.
  // svelte-ignore state_referenced_locally
  let selectedAgent = $state(agent);

  let base = $derived(`/${organization}/${project}/-/agents`);
  let showPicker = $derived(!agent && agents.length > 0);
  let targetAgent = $derived(agent || selectedAgent);
  // Choosing an agent already says what to run: the task is in its instructions, and the trigger it declares
  // carries the user turn that goes with them. Leaving this empty runs exactly what the schedule runs, so demanding
  // text here only made the launcher retype — or guess — something already written in the agent's YAML. The server
  // refuses an empty prompt for an agent with no trigger to take one from, which is the one case it is required.
  let canSubmit = $derived(!!targetAgent && !$startRun.isPending);

  function reset() {
    prompt = "";
    selectedAgent = agent;
  }

  async function handleSubmit() {
    if (!canSubmit) return;
    try {
      // A fresh idempotency key per submit so retries after a transient error
      // start a new run rather than colliding with a previous attempt.
      const res = await $startRun.mutateAsync({
        name: targetAgent,
        prompt: prompt.trim(),
        idempotencyKey: uuidv4(),
        trigger: "manual",
      });
      await queryClient.invalidateQueries({
        queryKey: getAgentServiceListAgentRunsQueryKey(
          runtimeClient.instanceId,
        ),
      });
      eventBus.emit("notification", {
        message: m.agents_run_form_started_notification(),
        type: "success",
      });
      open = false;
      reset();
      // Land on the new run so the user can watch it execute; fall back to the
      // runs list if the response did not carry an id.
      await goto(
        res.runId
          ? runDetailPath(organization, project, { runId: res.runId })
          : base,
      );
    } catch (e) {
      eventBus.emit("notification", {
        message: e instanceof Error ? e.message : m.agents_run_form_error(),
        type: "error",
      });
    }
  }
</script>

<Button type="secondary" onClick={() => (open = true)}>
  {label}
</Button>

<Dialog
  bind:open
  onOpenChange={(value) => {
    open = value;
    if (!value) reset();
  }}
>
  <DialogTrigger>
    {#snippet child({ props })}
      <div {...props} class="hidden"></div>
    {/snippet}
  </DialogTrigger>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>{label}</DialogTitle>
    </DialogHeader>
    {#if showPicker}
      <label class="flex flex-col gap-y-1 text-sm text-fg-secondary">
        {m.agents_field_agent()}
        <select
          class="border rounded px-2 py-1.5 text-sm bg-surface-base text-fg-primary"
          bind:value={selectedAgent}
          disabled={$startRun.isPending}
        >
          <option value="" disabled>—</option>
          {#each agents as a (a)}
            <option value={a}>{a}</option>
          {/each}
        </select>
      </label>
    {/if}
    <Textarea
      id="start-agent-run-prompt"
      label={m.agents_run_form_prompt_label()}
      placeholder={m.agents_run_form_prompt_placeholder()}
      bind:value={prompt}
      rows={5}
      disabled={$startRun.isPending}
    />
    <DialogFooter>
      <Button
        type="tertiary"
        disabled={$startRun.isPending}
        onClick={() => {
          open = false;
          reset();
        }}
      >
        {m.agents_run_form_cancel()}
      </Button>
      <Button
        type="primary"
        disabled={!canSubmit}
        loading={$startRun.isPending}
        onClick={handleSubmit}
      >
        {m.agents_run_form_submit()}
      </Button>
    </DialogFooter>
  </DialogContent>
</Dialog>

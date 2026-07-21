<script lang="ts">
  import { goto } from "$app/navigation";
  import Brain from "@rilldata/web-common/components/icons/Brain.svelte";
  import ClockCircle from "@rilldata/web-common/components/icons/ClockCircle.svelte";
  import Conversation from "@rilldata/web-common/components/icons/Conversation.svelte";
  import ResourceListEmptyState from "@rilldata/web-common/features/resources/ResourceListEmptyState.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { timeAgo } from "@rilldata/web-common/lib/time/relative-time";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import ApproveDenyButtons from "../approvals/ApproveDenyButtons.svelte";
  import { runDetailPath } from "../run-id";
  import { useAgents } from "../selectors";
  import type { AgentApprovalData, AgentRunData } from "../types";
  import { agentTriggerLabel, formatDateTime } from "../utils";
  import AgentRunStatusChip from "./AgentRunStatusChip.svelte";

  let {
    data,
    organization,
    project,
    // Pending approvals keyed by run id, so a waiting run offers approve/deny
    // inline without opening the detail. A run's approval does not travel on the
    // run itself.
    approvalsByRun = new Map(),
  }: {
    data: AgentRunData[];
    organization: string;
    project: string;
    approvalsByRun?: Map<string, AgentApprovalData>;
  } = $props();

  // Runs carry the agent's slug (agentName); the catalog carries its human
  // display name. Resolve the friendly name for the card, falling back to the
  // slug for agents without a display_name (e.g. ad-hoc test agents).
  const runtimeClient = useRuntimeClient();
  const agentsQuery = useAgents(runtimeClient);
  let displayNameByAgent = $derived(
    new Map<string, string>(
      ($agentsQuery.data?.agents ?? []).map((a) => [
        a.name ?? "",
        a.displayName ?? "",
      ]),
    ),
  );
  function agentLabel(run: AgentRunData): string {
    return displayNameByAgent.get(run.agentName ?? "") || run.agentName || "—";
  }

  function openRun(run: AgentRunData) {
    void goto(runDetailPath(organization, project, run));
  }

  // Ignore card clicks that land on an inline control (approve/deny, conversation
  // link) so those act on their own without also navigating to the detail.
  function onCardClick(e: MouseEvent, run: AgentRunData) {
    if ((e.target as HTMLElement).closest("button, a, select, input")) return;
    openRun(run);
  }

  function startedLabel(run: AgentRunData): { text: string; title: string } {
    const raw = run.startedOn ?? run.createdOn;
    const title = formatDateTime(raw);
    if (typeof raw !== "string" || raw === "") return { text: "—", title };
    const date = new Date(raw);
    if (Number.isNaN(date.getTime())) return { text: "—", title };
    return { text: timeAgo(date), title };
  }
</script>

{#if data.length === 0}
  <ResourceListEmptyState
    icon={ClockCircle}
    message={m.agents_runs_empty_message()}
    action={m.agents_runs_empty_action()}
  />
{:else}
  <ul class="flex flex-col rounded-lg border divide-y overflow-hidden">
    {#each data as run (run.runId)}
      {@const approval = approvalsByRun.get(run.runId ?? "")}
      {@const started = startedLabel(run)}
      <li>
        <div
          class="flex flex-wrap items-center gap-x-4 gap-y-2 px-4 py-3 cursor-pointer hover:bg-surface-hover"
          role="button"
          tabindex="0"
          onclick={(e) => onCardClick(e, run)}
          onkeydown={(e) => {
            if (e.key === "Enter" || e.key === " ") {
              e.preventDefault();
              openRun(run);
            }
          }}
        >
          <div class="flex items-center gap-x-2.5 min-w-0 grow">
            <span class="shrink-0 text-fg-secondary"><Brain size="15px" /></span
            >
            <div class="flex flex-col gap-y-1 min-w-0">
              <!-- Agent + action share the title line: an approval authorizes
                   this agent to run this specific tool, so the reviewer reads
                   both as one thing they are deciding on. The action only shows
                   when there is a pending approval. -->
              <div class="flex items-center flex-wrap gap-x-2 min-w-0">
                <span
                  class="text-fg-primary text-sm font-semibold truncate min-w-0"
                >
                  {agentLabel(run)}
                </span>
                {#if approval?.toolName}
                  <span class="text-fg-muted" aria-hidden="true">·</span>
                  <span
                    class="font-mono text-xs text-fg-secondary truncate max-w-[16rem]"
                    title={m.agents_approval_proposed_action()}
                  >
                    {approval.toolName}
                  </span>
                {/if}
              </div>
              <div
                class="flex items-center flex-wrap gap-x-2 gap-y-1 text-fg-secondary text-xs"
              >
                <AgentRunStatusChip status={run.status} />
                <span>{agentTriggerLabel(run.trigger)}</span>
                <span aria-hidden="true">·</span>
                <span title={started.title}>{started.text}</span>
              </div>
            </div>
          </div>

          <div class="flex items-center gap-x-2 shrink-0 ml-auto">
            {#if approval}
              <ApproveDenyButtons
                approvalId={approval.approvalId ?? ""}
                argsHash={approval.argsHash ?? ""}
              />
            {/if}
            {#if run.conversationId}
              <a
                class="text-fg-secondary hover:text-primary-600 p-1.5 rounded"
                href={`/${organization}/${project}/-/ai/${run.conversationId}`}
                title={m.agents_run_view_conversation()}
                aria-label={m.agents_run_view_conversation()}
              >
                <Conversation size="16px" />
              </a>
            {/if}
          </div>
        </div>
      </li>
    {/each}
  </ul>
{/if}

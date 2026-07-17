<script lang="ts">
  import { page } from "$app/state";
  import { composeRunId } from "@rilldata/web-admin/features/agents/run-id";
  import AgentRunMetadata from "@rilldata/web-admin/features/agents/runs/AgentRunMetadata.svelte";
  import ContentContainer from "@rilldata/web-common/components/layout/ContentContainer.svelte";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";

  // `agent` and `runId` (the idempotency key) arrive URL-decoded from SvelteKit.
  let { organization, project, agent, runId } = $derived(page.params);

  const client = useRuntimeClient();

  // The API keys off the backend's composite run id. instanceId and agentName are
  // implicit in the path, so we rebuild the composite here rather than carry it in
  // the URL — no API change needed.
  let composite = $derived(composeRunId(client.instanceId, agent, runId));
</script>

<ContentContainer maxWidth={1200}>
  <!-- Remount per run so the detail's queries re-init when navigating between
       runs (e.g. starting a new run from an existing run's detail). -->
  {#key composite}
    <AgentRunMetadata {organization} {project} runId={composite} />
  {/key}
</ContentContainer>

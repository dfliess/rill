<script lang="ts">
  import { page } from "$app/state";
  import ContentContainer from "@rilldata/web-common/components/layout/ContentContainer.svelte";
  import AgentRunMetadata from "@rilldata/web-admin/features/agents/runs/AgentRunMetadata.svelte";

  let { organization, project, run } = $derived(page.params);
</script>

<ContentContainer maxWidth={1200}>
  <!-- Remount per run so the detail's queries re-init when navigating between
       runs (e.g. starting a new run from an existing run's detail). -->
  {#key run}
    <AgentRunMetadata {organization} {project} runId={run} />
  {/key}
</ContentContainer>

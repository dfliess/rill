<script lang="ts">
  let { org, project }: { org: string; project: string } = $props();
</script>

<button
  data-testid="create-personal"
  aria-label="create"
  data-org={org}
  data-project={project}
></button>

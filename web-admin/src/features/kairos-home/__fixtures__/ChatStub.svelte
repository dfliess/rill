<script lang="ts">
  let { canvasName }: { canvasName: string | undefined } = $props();
</script>

<div data-testid="chat">{canvasName ?? ""}</div>

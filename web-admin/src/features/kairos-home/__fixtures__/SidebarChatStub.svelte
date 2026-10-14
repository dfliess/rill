<script lang="ts">
  import type { Snippet } from "svelte";

  let { inputContext }: { inputContext?: Snippet } = $props();
</script>

<div data-testid="input-footer">{@render inputContext?.()}</div>

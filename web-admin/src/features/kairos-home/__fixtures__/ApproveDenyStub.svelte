<script lang="ts">
  let { approvalId, small }: { approvalId: string; small?: boolean } = $props();
</script>

<span data-testid="approve-deny" data-approval={approvalId} data-small={small}
></span>

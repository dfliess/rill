<script lang="ts">
  import { Tag } from "@rilldata/web-common/components/tag";
  import { agentRunStatusColor, agentRunStatusLabel } from "../utils";

  let { status }: { status: string | undefined } = $props();
</script>

<Tag color={agentRunStatusColor(status)}>{agentRunStatusLabel(status)}</Tag>

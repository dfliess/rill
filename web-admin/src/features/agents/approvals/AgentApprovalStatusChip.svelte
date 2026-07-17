<script lang="ts">
  import { Tag } from "@rilldata/web-common/components/tag";
  import { agentApprovalStatusColor, agentApprovalStatusLabel } from "../utils";

  let { status }: { status: string | undefined } = $props();
</script>

<Tag color={agentApprovalStatusColor(status)}>
  {agentApprovalStatusLabel(status)}
</Tag>

<!--
  Keeps the installed app's icon telling the truth while the app is open (kairos-cloud#147).

  A push updates the badge for the project it came from, but a decision taken elsewhere sends nothing: no push is
  emitted when an approval is resolved, on purpose. So the open app recounts what is waiting in the project it is
  on and writes that into the shared ledger, which is also how the number comes down to zero and the badge clears.

  It renders nothing. It is mounted inside the project's RuntimeProvider, which is what scopes the count to the
  right deployment.
-->
<script lang="ts">
  import { usePendingApprovalRunCount } from "@rilldata/web-admin/features/agents/selectors";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import { scopeOf, syncBadge } from "./app-badge";

  let { organization, project }: { organization: string; project: string } =
    $props();

  const runtimeClient = useRuntimeClient();
  const countQuery = usePendingApprovalRunCount(runtimeClient);

  // Only a successful read is worth writing down. A project where Act is off, or a viewer without access to it,
  // answers with an error, and an error is not the same fact as "nothing is waiting for you here": treating it as
  // zero would wipe a badge that other projects are still entitled to.
  $effect(() => {
    if (!$countQuery.isSuccess) return;
    void syncBadge({
      scope: scopeOf(organization, project),
      count: $countQuery.data,
    });
  });
</script>

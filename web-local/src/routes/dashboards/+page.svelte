<script lang="ts">
  import * as m from "@rilldata/web-common/paraglide/messages.js";
  import ContentContainer from "@rilldata/web-common/components/layout/ContentContainer.svelte";
  import DashboardList from "../../features/dashboards/DashboardList.svelte";
</script>

<ContentContainer title={m.nav_dashboards()} maxWidth={1100}>
  <DashboardList showSearch />
</ContentContainer>

<script lang="ts">
  // Feeds the active Org's default Style to the style controller
  // (kairos-cloud#242). While it loads, the controller keeps the default this
  // browser last saw for the Org.
  import { createAdminServiceGetOrganizationStyle } from "@rilldata/web-admin/client";
  import { styleControl } from "@rilldata/web-common/features/kairos-styles/style-control";

  export let organization: string | undefined;

  $: orgStyle = createAdminServiceGetOrganizationStyle(organization ?? "", {
    query: { enabled: !!organization },
  });

  $: if (!organization) {
    styleControl.setOrg(null);
  } else if ($orgStyle.data) {
    styleControl.setOrg(organization, $orgStyle.data.style ?? "");
  } else {
    styleControl.setOrg(organization);
  }
</script>

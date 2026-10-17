<script lang="ts">
  // The Org's default Style (kairos-cloud#242). Members see it unless they
  // pick their own from the avatar menu.
  import {
    createAdminServiceGetOrganizationStyle,
    createAdminServiceUpdateOrganizationStyle,
    getAdminServiceGetOrganizationStyleQueryKey,
  } from "@rilldata/web-admin/client";
  import { getRpcErrorMessage } from "@rilldata/web-admin/components/errors/error-utils";
  import SettingsContainer from "@rilldata/web-admin/features/organizations/settings/SettingsContainer.svelte";
  import Select from "@rilldata/web-common/components/forms/Select.svelte";
  import { styleControl } from "@rilldata/web-common/features/kairos-styles/style-control";
  import {
    DEFAULT_STYLE,
    STYLES,
    type StyleId,
  } from "@rilldata/web-common/features/kairos-styles/styles";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";

  export let organization: string;

  const labels: Record<StyleId, () => string> = {
    "kairos-v1": m.kairos_style_kairos_v1,
    kairos: m.kairos_style_kairos,
    aurora: m.kairos_style_aurora,
    nova: m.kairos_style_nova,
    editorial: m.kairos_style_editorial,
    lumen: m.kairos_style_lumen,
    glass: m.kairos_style_glass,
  };
  const options = STYLES.map((style) => ({
    value: style,
    label: labels[style](),
  }));

  $: orgStyle = createAdminServiceGetOrganizationStyle(organization);
  const updater = createAdminServiceUpdateOrganizationStyle();

  $: current = $orgStyle.data?.style || DEFAULT_STYLE;

  async function onChange(style: string) {
    if (style === current) return;
    await $updater.mutateAsync({ org: organization, data: { style } });
    styleControl.setOrg(organization, style);
    void queryClient.invalidateQueries({
      queryKey: getAdminServiceGetOrganizationStyleQueryKey(organization),
    });
  }
</script>

<SettingsContainer title={m.kairos_style_org_title()}>
  <div class="flex flex-col gap-y-2">
    <div>{m.kairos_style_org_description()}</div>
    <Select
      id="kairos-org-style"
      ariaLabel={m.kairos_style_org_title()}
      value={current}
      {options}
      {onChange}
      disabled={$orgStyle.isPending || $updater.isPending}
      width={240}
    />
    {#if $updater.error}
      <div class="text-red-500 text-sm">
        {getRpcErrorMessage($updater.error)}
      </div>
    {/if}
  </div>
</SettingsContainer>

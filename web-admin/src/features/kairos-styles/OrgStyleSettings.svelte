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
  import Button from "@rilldata/web-common/components/button/Button.svelte";
  import Input from "@rilldata/web-common/components/forms/Input.svelte";
  import Select from "@rilldata/web-common/components/forms/Select.svelte";
  import { styleControl } from "@rilldata/web-common/features/kairos-styles/style-control";
  import {
    DEFAULT_STYLE,
    isStyleId,
    listedStyles,
    type StyleId,
  } from "@rilldata/web-common/features/kairos-styles/styles";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";

  export let organization: string;

  const labels: Record<StyleId, () => string> = {
    kairos: m.kairos_style_kairos,
    aurora: m.kairos_style_aurora,
    nova: m.kairos_style_nova,
    editorial: m.kairos_style_editorial,
    lumen: m.kairos_style_lumen,
    glass: m.kairos_style_glass,
    apetic: m.kairos_style_apetic,
    base: m.kairos_style_base,
  };
  // "Other…" takes an unlisted Style's id as free text (kairos-cloud#244).
  const OTHER = "__other";

  $: orgStyle = createAdminServiceGetOrganizationStyle(organization);
  const updater = createAdminServiceUpdateOrganizationStyle();

  $: current = $orgStyle.data?.style || DEFAULT_STYLE;
  $: options = [
    ...listedStyles(current).map((style) => ({
      value: style as string,
      label: labels[style](),
    })),
    { value: OTHER, label: m.kairos_style_other() },
  ];

  let typing = false;
  let typed = "";
  let unknown = false;

  function onSelect(value: string) {
    typing = value === OTHER;
    unknown = false;
    if (!typing) void onChange(value);
  }

  function applyTyped() {
    const id = typed.trim().toLowerCase();
    unknown = !isStyleId(id);
    if (unknown) return;
    typing = false;
    typed = "";
    void onChange(id);
  }

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
      value={typing ? OTHER : current}
      {options}
      onChange={onSelect}
      disabled={$orgStyle.isPending || $updater.isPending}
      width={240}
    />
    {#if typing}
      <form
        id="kairos-org-style-other-form"
        class="flex items-start gap-x-2"
        on:submit|preventDefault={applyTyped}
      >
        <div class="w-[240px]">
          <Input
            id="kairos-org-style-other"
            placeholder={m.kairos_style_other_placeholder()}
            bind:value={typed}
            onEnter={applyTyped}
            alwaysShowError
            errors={unknown
              ? m.kairos_style_unknown({ id: typed.trim() })
              : undefined}
          />
        </div>
        <Button
          type="secondary"
          submitForm
          form="kairos-org-style-other-form"
          disabled={!typed.trim() || $updater.isPending}
        >
          {m.kairos_style_other_apply()}
        </Button>
      </form>
    {/if}
    {#if $updater.error}
      <div class="text-red-500 text-sm">
        {getRpcErrorMessage($updater.error)}
      </div>
    {/if}
  </div>
</SettingsContainer>

<script lang="ts">
  import * as DropdownMenu from "@rilldata/web-common/components/dropdown-menu";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { styleControl } from "./style-control";
  import { STYLES, type StyleId } from "./styles";

  // Rill Developer has no Org: the fallback is just the default Style.
  export let local = false;

  const { preference } = styleControl;

  const labels: Record<StyleId, () => string> = {
    "kairos-v1": m.kairos_style_kairos_v1,
    kairos: m.kairos_style_kairos,
    aurora: m.kairos_style_aurora,
    nova: m.kairos_style_nova,
    editorial: m.kairos_style_editorial,
    lumen: m.kairos_style_lumen,
    glass: m.kairos_style_glass,
  };
</script>

<DropdownMenu.Sub>
  <DropdownMenu.SubTrigger>{m.kairos_style_label()}</DropdownMenu.SubTrigger>
  <DropdownMenu.SubContent>
    <DropdownMenu.CheckboxItem
      checkRight
      checked={$preference === null}
      onclick={() => styleControl.setPreference(null)}
    >
      {local ? m.kairos_style_default() : m.kairos_style_org_default()}
    </DropdownMenu.CheckboxItem>
    <DropdownMenu.Separator />
    {#each STYLES as style (style)}
      <DropdownMenu.CheckboxItem
        checkRight
        checked={$preference === style}
        onclick={() => styleControl.setPreference(style)}
      >
        {labels[style]()}
      </DropdownMenu.CheckboxItem>
    {/each}
  </DropdownMenu.SubContent>
</DropdownMenu.Sub>

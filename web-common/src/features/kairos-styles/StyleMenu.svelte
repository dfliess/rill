<script lang="ts">
  import * as DropdownMenu from "@rilldata/web-common/components/dropdown-menu";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { styleControl } from "./style-control";
  import { STYLES, type StyleId } from "./styles";

  const { preference } = styleControl;

  // Charts resolve their colors when they render and nothing re-renders them
  // on a Style change, so a switch reloads the page instead of leaving stale
  // series behind.
  function choose(style: StyleId | null) {
    if (style === $preference) return;
    styleControl.setPreference(style);
    window.location.reload();
  }

  const labels: Record<StyleId, () => string> = {
    "kairos-v1": m.kairos_style_kairos_v1,
    kairos: m.kairos_style_kairos,
    aurora: m.kairos_style_aurora,
  };
</script>

<DropdownMenu.Sub>
  <DropdownMenu.SubTrigger>{m.kairos_style_label()}</DropdownMenu.SubTrigger>
  <DropdownMenu.SubContent>
    <DropdownMenu.CheckboxItem
      checkRight
      checked={$preference === null}
      onclick={() => choose(null)}
    >
      {m.kairos_style_org_default()}
    </DropdownMenu.CheckboxItem>
    <DropdownMenu.Separator />
    {#each STYLES as style (style)}
      <DropdownMenu.CheckboxItem
        checkRight
        checked={$preference === style}
        onclick={() => choose(style)}
      >
        {labels[style]()}
      </DropdownMenu.CheckboxItem>
    {/each}
  </DropdownMenu.SubContent>
</DropdownMenu.Sub>

<script lang="ts">
  import * as DropdownMenu from "@rilldata/web-common/components/dropdown-menu";
  import {
    themeControl,
    type ThemeMode,
  } from "@rilldata/web-common/features/themes/theme-control";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { DropdownMenu as DropdownMenuPrimitive } from "bits-ui";
  import { Check, Monitor, Moon, Sun } from "lucide-svelte";
  import { styleControl } from "./style-control";
  import { listedStyles, type StyleId } from "./styles";

  // Rill Developer has no Org: the fallback is just the default Style.
  let { local = false }: { local?: boolean } = $props();

  const { preference, current } = styleControl;
  const { preference: mode } = themeControl;

  // The Mode lives here too, on top: an Appearance is judged in both, so one
  // menu switches both without reopening (it replaced Rill's Mode submenu).
  const modes: { id: ThemeMode; label: () => string; icon: typeof Sun }[] = [
    { id: "light", label: m.theme_light, icon: Sun },
    { id: "dark", label: m.theme_dark, icon: Moon },
    { id: "system", label: m.theme_system, icon: Monitor },
  ];

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

  // "" is the Org's Style (or, local, the default).
  // Unlisted Styles only show while they are the one in use.
  const options: { id: StyleId | ""; label: () => string }[] = $derived([
    {
      id: "",
      label: local ? m.kairos_style_default : m.kairos_style_org_default,
    },
    ...listedStyles($current).map((id) => ({ id, label: labels[id] })),
  ]);
</script>

<DropdownMenu.Sub>
  <DropdownMenu.SubTrigger>{m.kairos_style_label()}</DropdownMenu.SubTrigger>
  <DropdownMenu.SubContent>
    <div
      class="kairos-mode-switch m-1 mb-1.5 grid grid-cols-3 gap-0.5 rounded bg-surface-muted p-0.5"
      role="group"
      aria-label={m.theme_label()}
    >
      {#each modes as { id, label, icon: Icon } (id)}
        <button
          type="button"
          class="flex items-center justify-center gap-x-1 rounded-sm px-2 py-1 text-xs text-fg-secondary hover:text-fg-primary"
          class:selected={$mode === id}
          aria-pressed={$mode === id}
          onclick={themeControl.set[id]}
        >
          <Icon size="12px" />
          {label()}
        </button>
      {/each}
    </div>
    <DropdownMenu.Separator />
    <!-- One choice, so a radio group: the menu stays open to compare, and a
         checkbox item kept its own checked state, so the tick stayed behind
         on the old one. -->
    <DropdownMenuPrimitive.RadioGroup
      bind:value={
        () => $preference ?? "",
        (v) => styleControl.setPreference((v || null) as StyleId | null)
      }
    >
      {#each options as { id, label }, i (id)}
        {#if i === 1}
          <DropdownMenu.Separator />
        {/if}
        <DropdownMenuPrimitive.RadioItem
          value={id}
          closeOnSelect={false}
          class="relative flex cursor-pointer select-none items-center justify-between gap-x-2 rounded-sm px-2 py-1.5 text-xs text-fg-primary outline-none data-[highlighted]:bg-popover-accent data-[highlighted]:text-fg-accent hover:bg-popover-accent"
        >
          {#snippet children({ checked })}
            {label()}
            <span
              class="flex h-3.5 w-3.5 flex-none items-center justify-center"
            >
              {#if checked}<Check class="h-4 w-4" />{/if}
            </span>
          {/snippet}
        </DropdownMenuPrimitive.RadioItem>
      {/each}
    </DropdownMenuPrimitive.RadioGroup>
  </DropdownMenu.SubContent>
</DropdownMenu.Sub>

<style lang="postcss">
  .selected {
    @apply bg-popover text-fg-primary font-semibold shadow-sm;
  }
</style>

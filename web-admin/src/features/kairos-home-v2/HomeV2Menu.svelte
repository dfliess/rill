<script lang="ts">
  import { goto } from "$app/navigation";
  import Resizer from "@rilldata/web-common/layout/Resizer.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { HOME_V2_MENU_WIDTH, homeV2MenuWidth } from "./chat-state";
  import { homeV2MenuIcon } from "./menu-icons";
  import type { HomeV2MenuItem } from "./menu";

  let {
    items,
    active,
  }: { items: HomeV2MenuItem[]; active: string | undefined } = $props();
</script>

<!-- Phones get a selector instead of a column: three columns do not fit. -->
<div class="md:hidden px-4 pt-3">
  <select
    class="w-full rounded border bg-surface-background px-2 py-1.5 text-sm"
    aria-label={m.home_v2_menu_label()}
    value={active}
    onchange={(e) => {
      const item = items.find((i) => i.name === e.currentTarget.value);
      if (item) void goto(item.href);
    }}
  >
    {#each items as item (item.name)}
      <option value={item.name}>{item.title}</option>
    {/each}
  </select>
</div>

<!-- Resizable like the AI panel (drag the edge, double-click to reset), and
     titles wrap instead of truncating: a dashboard's name is how it is found. -->
<!-- The resizer sits on a wrapper: the scrolling list would clip its handle. -->
<div
  class="relative hidden md:flex shrink-0 border-r"
  style:width="{$homeV2MenuWidth}px"
>
  <Resizer
    min={HOME_V2_MENU_WIDTH.MIN}
    max={HOME_V2_MENU_WIDTH.MAX}
    basis={HOME_V2_MENU_WIDTH.BASIS}
    dimension={$homeV2MenuWidth}
    direction="EW"
    side="right"
    onUpdate={(w) => homeV2MenuWidth.set(w)}
  />
  <nav
    class="flex w-full flex-col gap-y-0.5 overflow-y-auto px-2 py-4"
    aria-label={m.home_v2_menu_label()}
  >
    {#each items as item (item.name)}
      {@const Icon = homeV2MenuIcon(item.icon)}
      <a
        href={item.href}
        class="flex items-start gap-x-2.5 rounded px-3 py-2 text-sm leading-snug text-fg-secondary hover:bg-surface-hover"
        class:active={item.name === active}
        aria-current={item.name === active ? "page" : undefined}
      >
        <!-- Pinned to the first line of a wrapped title. -->
        <span
          class="flex h-[1.375em] shrink-0 items-center"
          data-icon={item.icon}><Icon size="16px" /></span
        >
        <span class="min-w-0">{item.title}</span>
      </a>
    {/each}
  </nav>
</div>

<style lang="postcss">
  .active {
    @apply bg-surface-hover font-semibold text-fg-primary;
  }
</style>

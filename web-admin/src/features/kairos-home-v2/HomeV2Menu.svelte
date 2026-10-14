<script lang="ts">
  import { goto } from "$app/navigation";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";

  export type HomeV2MenuItem = { name: string; title: string; href: string };

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

<nav
  class="hidden md:flex flex-col gap-y-0.5 w-56 shrink-0 overflow-y-auto border-r px-2 py-4"
  aria-label={m.home_v2_menu_label()}
>
  {#each items as item (item.name)}
    <a
      href={item.href}
      class="truncate rounded px-3 py-2 text-sm text-fg-secondary hover:bg-surface-hover"
      class:active={item.name === active}
      aria-current={item.name === active ? "page" : undefined}
    >
      {item.title}
    </a>
  {/each}
</nav>

<style lang="postcss">
  .active {
    @apply bg-surface-hover font-semibold text-fg-primary;
  }
</style>

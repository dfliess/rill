<script lang="ts">
  import { goto } from "$app/navigation";
  import CreatePersonalCanvasDialog from "@rilldata/web-admin/features/personal-files/canvas/CreatePersonalCanvasDialog.svelte";
  import Lock from "@rilldata/web-common/components/icons/Lock.svelte";
  import Resizer from "@rilldata/web-common/layout/Resizer.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { Plus } from "lucide-svelte";
  import { HOME_MENU_WIDTH, homeMenuWidth } from "./chat-state";
  import { homeMenuIcon } from "./menu-icons";
  import type { HomeMenuItem } from "./menu";

  let {
    items,
    active,
    personalItems,
    organization,
    project,
  }: {
    items: HomeMenuItem[];
    active: string | undefined;
    // The viewer's personal canvases; undefined hides the group (feature off,
    // or nobody signed in).
    personalItems?: HomeMenuItem[];
    organization: string;
    project: string;
  } = $props();

  let activeHref = $derived(items.find((i) => i.name === active)?.href);
</script>

<!-- Phones get a selector instead of a column: three columns do not fit. -->
<div class="md:hidden px-4 pt-3">
  <select
    class="w-full rounded border bg-surface-background px-2 py-1.5 text-sm"
    aria-label={m.home_menu_label()}
    value={activeHref}
    onchange={(e) => void goto(e.currentTarget.value)}
  >
    {#each items as item (item.href)}
      <option value={item.href}>{item.title}</option>
    {/each}
    {#if personalItems?.length}
      <optgroup label={m.personal_files_canvases_title()}>
        {#each personalItems as item (item.href)}
          <option value={item.href}>{item.title}</option>
        {/each}
      </optgroup>
    {/if}
  </select>
</div>

<!-- Resizable like the AI panel (drag the edge, double-click to reset), and
     titles wrap instead of truncating: a dashboard's name is how it is found. -->
<!-- The resizer sits on a wrapper: the scrolling list would clip its handle. -->
<div
  class="relative hidden md:flex shrink-0 border-r"
  style:width="{$homeMenuWidth}px"
>
  <Resizer
    min={HOME_MENU_WIDTH.MIN}
    max={HOME_MENU_WIDTH.MAX}
    basis={HOME_MENU_WIDTH.BASIS}
    dimension={$homeMenuWidth}
    direction="EW"
    side="right"
    onUpdate={(w) => homeMenuWidth.set(w)}
  />
  <div class="flex w-full flex-col gap-y-5 overflow-y-auto px-2 py-4">
    <nav class="flex flex-col gap-y-0.5" aria-label={m.home_menu_label()}>
      {#each items as item (item.name)}
        {@const Icon = homeMenuIcon(item.icon)}
        <a
          href={item.href}
          class="item"
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

    <!-- The viewer's personal canvases, apart from the project's: only they
         see them. The "+" is the only place to create one (#231). -->
    {#if personalItems}
      <nav
        class="flex flex-col gap-y-0.5"
        aria-label={m.personal_files_canvases_title()}
      >
        <div
          class="flex items-center gap-x-2 px-3 pb-1 text-xs font-semibold text-fg-muted"
        >
          <span class="min-w-0 flex-1">{m.personal_files_canvases_title()}</span
          >
          <CreatePersonalCanvasDialog org={organization} {project}>
            {#snippet trigger(props)}
              <button
                {...props}
                type="button"
                class="rounded p-0.5 hover:bg-surface-hover hover:text-fg-primary"
                aria-label={m.personal_files_create_trigger()}
                title={m.personal_files_create_trigger()}
              >
                <Plus size="14px" />
              </button>
            {/snippet}
          </CreatePersonalCanvasDialog>
        </div>
        {#each personalItems as item (item.name)}
          <a href={item.href} class="item">
            <span class="flex h-[1.375em] shrink-0 items-center"
              ><Lock size="14px" /></span
            >
            <span class="min-w-0">{item.title}</span>
          </a>
        {/each}
      </nav>
    {/if}
  </div>
</div>

<style lang="postcss">
  .item {
    @apply flex items-start gap-x-2.5 rounded px-3 py-2 text-sm leading-snug text-fg-secondary;
  }
  .item:hover {
    @apply bg-surface-hover;
  }
  .active {
    @apply bg-surface-hover font-semibold text-fg-primary;
  }
</style>

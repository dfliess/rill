<script lang="ts">
  import type { Snippet } from "svelte";

  // The frame Home's Monitor and Act summaries share (#231): one card each,
  // same header, same row rhythm, so two of them side by side read as a pair
  // and a grid stretches them to the same height. The rows are the caller's.
  let {
    title,
    count = 0,
    href,
    viewAllLabel,
    emptyText,
    loaded,
    children,
  }: {
    title: string;
    count?: number;
    href: string;
    viewAllLabel: string;
    emptyText: string;
    // "Nothing here" is only a fact once the data has answered.
    loaded: boolean;
    children: Snippet;
  } = $props();
</script>

<!-- Always shown, even when empty: an empty card tells the viewer where
     Monitor and Act will speak up, while a missing one hides the feature. -->
<!-- Framed like a canvas component (outline, shadow, radius) so under a
     canvas they read as part of it. -->
<section
  class="flex min-w-0 flex-col overflow-hidden rounded-sm bg-surface-card shadow-sm outline outline-[1px] outline-border"
>
  <header class="flex items-baseline justify-between gap-x-3 px-4 pb-1 pt-3">
    <h2 class="text-sm font-semibold text-fg-primary">
      {title}
      {#if count > 0}
        <span class="ml-1 font-normal text-fg-muted">{count}</span>
      {/if}
    </h2>
    <a class="shrink-0 text-xs text-primary-600 hover:underline" {href}>
      {viewAllLabel}
    </a>
  </header>
  {#if count > 0}
    <ul class="flex flex-col divide-y pb-1">
      {@render children()}
    </ul>
  {:else if loaded}
    <p
      class="flex flex-1 items-center justify-center px-4 pb-5 pt-3 text-center text-sm text-fg-muted"
    >
      {emptyText}
    </p>
  {/if}
</section>

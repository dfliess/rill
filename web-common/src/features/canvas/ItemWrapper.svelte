<script lang="ts">
  import { getMinHeight, isExpandable } from "./layout-util";

  export let zIndex: number;
  export let type: string | undefined = undefined;
  // When true, the wrapper sizes to its content rather than the type's fixed
  // initial height. Used for the "add widget" filler so it doesn't reserve a
  // full row's worth of empty space below the button.
  export let fitContent = false;

  $: expandable = isExpandable(type);
  $: minHeight = getMinHeight(type) + "px";
</script>

<div
  style:z-index={zIndex}
  class:expandable
  class:fit-content={fitContent}
  style:--min-height={minHeight}
>
  <slot />
</div>

<style lang="postcss">
  div {
    @apply p-2.5 relative pointer-events-none size-full;
    container-type: inline-size;
    container-name: component-container;
  }

  .fit-content {
    @apply h-fit;
  }

  .expandable:not(.fit-content) {
    min-height: var(--row-height);
    /* Cancels the `height: 100%` that `size-full` applies. Blink resolves that
       100% against an auto-sized grid track as `auto`, so the item still grew
       and nobody noticed; WebKit treats it as a definite height, so the item
       stayed at the floor and clipped its content. Visible on every iOS
       browser and in desktop Safari once the row wraps to one column, and not
       reproducible in Chrome at any width. The item still fills a taller row,
       because a grid item stretches by default. */
    height: auto;
  }

  :not(.expandable):not(.fit-content) {
    height: max(var(--row-height), var(--min-height));
    min-height: 100%;
  }
</style>

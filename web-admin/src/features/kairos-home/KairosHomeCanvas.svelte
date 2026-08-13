<script lang="ts">
  import { isKairosHome } from "@rilldata/web-admin/features/dashboards/listing/selectors";
  import CanvasDashboardEmbed from "@rilldata/web-common/features/canvas/CanvasDashboardEmbed.svelte";
  import CanvasProvider from "@rilldata/web-common/features/canvas/CanvasProvider.svelte";
  import { createRuntimeServiceListResources } from "@rilldata/web-common/runtime-client";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";

  const runtimeClient = useRuntimeClient();

  // Canvases marked with the kairos_home annotation render here on the project
  // home; useDashboards excludes them from the dashboards listing. Isolated so
  // canvas state never rewrites the home page URL.
  $: canvasesQuery = createRuntimeServiceListResources(
    runtimeClient,
    {},
    {
      query: { enabled: !!runtimeClient.instanceId },
    },
  );
  $: homeCanvases = ($canvasesQuery.data?.resources ?? [])
    .filter(
      (res) => res.canvas && isKairosHome(res.canvas?.state?.validSpec ?? {}),
    )
    .map((res) => res.meta?.name?.name)
    .filter((name): name is string => !!name)
    .sort();
</script>

{#each homeCanvases as canvasName (canvasName)}
  {#key `${runtimeClient.instanceId}::${canvasName}`}
    <!-- A minimum, not a fixed height: 156px is one 140px card row plus the embed's padding, which is all the
         original numbers strip ever held. Clamping it hid anything a home canvas put in a second row (an
         agent note, say) with no sign that content was there. Below md the canvas stacks the cards into one
         column (RowWrapper container query) and grows on its own. -->
    <div class="kairos-home-strip md:min-h-[156px]">
      <CanvasProvider
        {canvasName}
        instanceId={runtimeClient.instanceId}
        isolated
      >
        <CanvasDashboardEmbed {canvasName} navigationEnabled={false} />
      </CanvasProvider>
    </div>
  {/key}
{/each}

<style>
  /* The embedded canvas paints its own surface and card fills, which blend on
     the white home in light mode but read as mismatched gray boxes in dark.
     On the home strip both inherit the page background; the cards keep their
     outline as the only box treatment. */
  .kairos-home-strip :global(.bg-surface-background) {
    background: transparent;
  }
  .kairos-home-strip :global(.component-card) {
    background: transparent;
  }

  /* The embed insets its content twice horizontally: 8px on the scroll container
     and another 10px on every item. On a canvas page that gutter is the frame,
     but here it makes the strip 18px narrower per side than every other block on
     the home (the ask box, the approvals inbox, the dashboard listing), which
     reads as a misaligned edge rather than as padding. Drop the outer inset and
     pull the row out by one item gutter, so the strip's edge lines up with its
     siblings while the 20px between one card and the next is untouched. The
     vertical padding stays: it is what separates the strip from the heading. */
  .kairos-home-strip :global(#canvas-scroll-container) {
    padding-left: 0;
    padding-right: 0;
  }
  .kairos-home-strip :global(.canvas-row) {
    margin-left: -10px;
    margin-right: -10px;
  }

  /* Below md the canvas stacks the cards into one column; compact them so the
     strip doesn't push the welcome half a screen down. 120px is the floor for
     the card's pixel-positioned text block (title/number/delta/caption).
     !important outranks the ItemWrapper height rule, scoped to this strip. */
  @media (max-width: 767px) {
    .kairos-home-strip :global(.canvas-row > div) {
      height: 120px !important;
    }
  }
</style>

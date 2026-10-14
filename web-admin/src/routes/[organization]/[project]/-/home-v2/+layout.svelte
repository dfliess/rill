<script lang="ts">
  import { page } from "$app/stores";
  import ConnectClientProvider from "@rilldata/web-admin/features/ai/mcp/ConnectClientProvider.svelte";
  import HomePendingApprovals from "@rilldata/web-admin/features/agents/HomePendingApprovals.svelte";
  import HomeRecentAlerts from "@rilldata/web-admin/features/kairos-home/HomeRecentAlerts.svelte";
  import HomeV2Chat from "@rilldata/web-admin/features/kairos-home-v2/HomeV2Chat.svelte";
  import HomeV2Menu from "@rilldata/web-admin/features/kairos-home-v2/HomeV2Menu.svelte";
  import {
    HOME_V2_SIDE_BY_SIDE_QUERY,
    homeV2ChatActions,
    homeV2ChatOpen,
  } from "@rilldata/web-admin/features/kairos-home-v2/chat-state";
  import { homeV2MenuItems } from "@rilldata/web-admin/features/kairos-home-v2/menu";
  import CanvasDashboardEmbed from "@rilldata/web-common/features/canvas/CanvasDashboardEmbed.svelte";
  import CanvasProvider from "@rilldata/web-common/features/canvas/CanvasProvider.svelte";
  import { DEFAULT_DASHBOARD_WIDTH } from "@rilldata/web-common/features/canvas/layout-util";
  import { featureFlags } from "@rilldata/web-common/features/feature-flags";
  import { createRuntimeServiceListResources } from "@rilldata/web-common/runtime-client";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import { onMount, type Snippet } from "svelte";

  let { children }: { children: Snippet } = $props();

  const { agents, chat } = featureFlags;
  const runtimeClient = useRuntimeClient();

  let organization = $derived($page.params.organization);
  let project = $derived($page.params.project);

  // Same unfiltered list the alerts box uses, so one cached request feeds the
  // menu and the alerts.
  const resourcesQuery = createRuntimeServiceListResources(
    runtimeClient,
    {},
    { query: { enabled: !!runtimeClient.instanceId } },
  );

  let menuItems = $derived(
    homeV2MenuItems(
      $resourcesQuery.data?.resources ?? [],
      `/${organization}/${project}/-/home-v2`,
    ),
  );
  // No canvas in the URL: the first one in the menu.
  let selected = $derived($page.params.canvas ?? menuItems[0]?.name);
  let onFirst = $derived(selected === menuItems[0]?.name);
  let canvasMaxWidth = $derived(
    $resourcesQuery.data?.resources?.find(
      (r) => r.canvas && r.meta?.name?.name === selected,
    )?.canvas?.state?.validSpec?.maxWidth || DEFAULT_DASHBOARD_WIDTH,
  );
  let selectedTitle = $derived(
    menuItems.find((i) => i.name === selected)?.title ?? "",
  );

  // Beside the canvas on wide screens; on narrower ones it waits for the
  // header toggle.
  onMount(() => {
    homeV2ChatOpen.set(window.matchMedia(HOME_V2_SIDE_BY_SIDE_QUERY).matches);
  });
</script>

<div class="home-v2">
  <HomeV2Menu items={menuItems} active={selected} />

  <main class="min-w-0 flex-1 overflow-y-auto">
    <!-- A plain block around the canvas: the embed is size-full, which
         against the scrolling column would stretch it to the viewport and push
         the Home boxes below a screen of blank space. Against this auto-height
         block it takes its content's height. -->
    <div>
      {#if selected}
        {#key `${runtimeClient.instanceId}::${selected}`}
          <CanvasProvider
            canvasName={selected}
            instanceId={runtimeClient.instanceId}
            isolated
          >
            <CanvasDashboardEmbed
              canvasName={selected}
              navigationEnabled={false}
            />
          </CanvasProvider>
        {/key}
      {/if}
    </div>

    {#if onFirst}
      <!-- Lined up with the canvas above: its 8px frame, its max width and
           its 10px around every component, so the cards share its edges and
           its 20px gutter. Side by side while each card gets ~24rem, stacked
           below that; that follows the width left for the canvas, not the
           viewport, so opening the AI panel can stack them. -->
      <div class="flex justify-center px-2 pb-8">
        <div
          class="grid w-full grid-cols-[repeat(auto-fit,minmax(min(100%,24rem),1fr))] gap-5 px-2.5 pt-0.5"
          style:max-width="{canvasMaxWidth}px"
        >
          <HomeRecentAlerts {organization} {project} />
          {#if $agents}
            <HomePendingApprovals {organization} {project} />
          {/if}
        </div>
      </div>
    {/if}
  </main>

  {#if $chat && $homeV2ChatOpen}
    <aside class="home-v2-chat">
      <ConnectClientProvider>
        <HomeV2Chat
          canvasName={selected}
          canvasTitle={selectedTitle}
          actions={homeV2ChatActions}
        />
      </ConnectClientProvider>
    </aside>
  {/if}

  <!-- The menu and the canvas live in this layout so they survive navigation
       between canvases; the pages under it are empty. -->
  {@render children()}
</div>

<style lang="postcss">
  .home-v2 {
    @apply flex flex-col bg-surface-background md:flex-row;
    flex: 1;
    min-height: 0;
    overflow: hidden;
  }

  /* Beside the canvas from 1280px (HOME_V2_SIDE_BY_SIDE_QUERY); over it below,
     where three columns leave the canvas too narrow to read. */
  .home-v2-chat {
    @apply fixed inset-y-0 right-0 z-40 shadow-lg;
  }
  @media (min-width: 1280px) {
    .home-v2-chat {
      @apply static z-auto shadow-none;
    }
  }
  /* Rill's 500px panel would squeeze the canvas into one column on a laptop;
     here it starts at 380px until the screen has room for both. */
  @media (min-width: 1280px) and (max-width: 1599px) {
    .home-v2-chat :global(.chat-sidebar) {
      width: min(var(--sidebar-width), 380px);
    }
  }
</style>

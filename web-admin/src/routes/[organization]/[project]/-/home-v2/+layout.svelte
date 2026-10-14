<script lang="ts">
  import { page } from "$app/stores";
  import ConnectClientProvider from "@rilldata/web-admin/features/ai/mcp/ConnectClientProvider.svelte";
  import HomePendingApprovals from "@rilldata/web-admin/features/agents/HomePendingApprovals.svelte";
  import { isKairosHome } from "@rilldata/web-admin/features/dashboards/listing/selectors";
  import HomeRecentAlerts from "@rilldata/web-admin/features/kairos-home/HomeRecentAlerts.svelte";
  import HomeV2Chat from "@rilldata/web-admin/features/kairos-home-v2/HomeV2Chat.svelte";
  import HomeV2Menu from "@rilldata/web-admin/features/kairos-home-v2/HomeV2Menu.svelte";
  import Button from "@rilldata/web-common/components/button/Button.svelte";
  import CanvasDashboardEmbed from "@rilldata/web-common/features/canvas/CanvasDashboardEmbed.svelte";
  import CanvasProvider from "@rilldata/web-common/features/canvas/CanvasProvider.svelte";
  import { createChatActions } from "@rilldata/web-common/features/chat/layouts/sidebar/sidebar-store";
  import { featureFlags } from "@rilldata/web-common/features/feature-flags";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { createRuntimeServiceListResources } from "@rilldata/web-common/runtime-client";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import { onMount, type Snippet } from "svelte";
  import { writable } from "svelte/store";

  let { children }: { children: Snippet } = $props();

  const { agents, chat } = featureFlags;
  const runtimeClient = useRuntimeClient();

  let organization = $derived($page.params.organization);
  let project = $derived($page.params.project);
  let base = $derived(`/${organization}/${project}/-/home-v2`);

  // Same unfiltered list the Home and its alerts box use, so one cached request
  // feeds the menu, the home canvas and the alerts.
  const resourcesQuery = createRuntimeServiceListResources(
    runtimeClient,
    {},
    { query: { enabled: !!runtimeClient.instanceId } },
  );

  let canvases = $derived(
    ($resourcesQuery.data?.resources ?? [])
      .filter((res) => !!res.canvas?.state?.validSpec)
      .map((res) => ({
        name: res.meta?.name?.name ?? "",
        title: res.canvas?.spec?.displayName || res.meta?.name?.name || "",
        home: isKairosHome(res.canvas?.state?.validSpec ?? {}),
      }))
      .filter((c) => !!c.name),
  );
  // The first kairos_home canvas, by name, is the page's default.
  let homeCanvas = $derived(
    canvases
      .filter((c) => c.home)
      .map((c) => c.name)
      .sort()[0],
  );
  let selected = $derived($page.params.canvas ?? homeCanvas);
  let onHome = $derived(selected === homeCanvas);

  let menuItems = $derived([
    ...(homeCanvas
      ? [{ name: homeCanvas, title: m.nav_tab_home(), href: base }]
      : []),
    ...canvases
      .filter((c) => !c.home)
      .sort((a, b) => a.title.localeCompare(b.title))
      .map((c) => ({
        name: c.name,
        title: c.title,
        href: `${base}/${encodeURIComponent(c.name)}`,
      })),
  ]);
  let selectedTitle = $derived(
    menuItems.find((i) => i.name === selected)?.title ?? "",
  );

  // Open on wide screens, a button away on tablets and phones. Its own state:
  // opening it here must not open the side panel on every dashboard.
  const chatOpen = writable(false);
  const chatActions = createChatActions(chatOpen);
  onMount(() => {
    chatOpen.set(window.matchMedia("(min-width: 1024px)").matches);
  });
</script>

<div class="home-v2">
  <HomeV2Menu items={menuItems} active={selected} />

  <main class="min-w-0 flex-1 overflow-y-auto">
    {#if $chat && !$chatOpen}
      <div class="flex justify-end px-4 pt-3">
        <Button type="secondary" compact onClick={() => chatActions.openChat()}>
          {m.home_v2_open_ai()}
        </Button>
      </div>
    {/if}

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

    {#if onHome}
      <div class="grid grid-cols-1 gap-8 px-4 pb-8 md:grid-cols-2 md:px-6">
        <HomeRecentAlerts {organization} {project} />
        {#if $agents}
          <HomePendingApprovals {organization} {project} />
        {/if}
      </div>
    {/if}
  </main>

  {#if $chat && $chatOpen}
    <aside class="home-v2-chat">
      <ConnectClientProvider>
        <HomeV2Chat
          canvasName={selected}
          canvasTitle={selectedTitle}
          actions={chatActions}
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

  /* Beside the canvas on wide screens; over it below lg, where three columns
     leave the canvas too narrow to read. */
  .home-v2-chat {
    @apply fixed inset-y-0 right-0 z-40 shadow-lg;
  }
  @media (min-width: 1024px) {
    .home-v2-chat {
      @apply static z-auto shadow-none;
    }
  }
</style>

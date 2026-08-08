<script lang="ts">
  import { page } from "$app/stores";
  import ContentContainer from "@rilldata/web-common/components/layout/ContentContainer.svelte";
  import DashboardsTable from "@rilldata/web-admin/features/dashboards/listing/DashboardsTable.svelte";
  import HomePendingApprovals from "@rilldata/web-admin/features/agents/HomePendingApprovals.svelte";
  import InlineChat from "@rilldata/web-common/features/chat/layouts/inline/InlineChat.svelte";
  import DelayedContent from "@rilldata/web-common/features/entity-management/DelayedContent.svelte";
  import { featureFlags } from "@rilldata/web-common/features/feature-flags";
  import { createRuntimeServiceGetInstance } from "@rilldata/web-common/runtime-client";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import HomeRecentAlerts from "@rilldata/web-admin/features/kairos-home/HomeRecentAlerts.svelte";
  import KairosHomeCanvas from "@rilldata/web-admin/features/kairos-home/KairosHomeCanvas.svelte";
  import PersonalCanvasesList from "@rilldata/web-admin/features/personal-files/canvas/PersonalCanvasesList.svelte";
  import CreatePersonalCanvasDialog from "@rilldata/web-admin/features/personal-files/canvas/CreatePersonalCanvasDialog.svelte";
  import { getPersonalFilteredResources } from "@rilldata/web-admin/features/personal-files/selectors.ts";
  import { ResourceKind } from "@rilldata/web-common/features/entity-management/resource-selectors.ts";
  import DashboardsTagFilter from "@rilldata/web-admin/features/dashboards/listing/DashboardsTagFilter.svelte";
  import { UrlParamsState } from "web-common/src/lib/store-utils/url-params-state.svelte.ts";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { escapeHtml } from "@rilldata/web-common/lib/i18n";
  import TableToolbarSort from "@rilldata/web-common/components/table-toolbar/TableToolbarSort.svelte";
  import { DashboardTableSortOptions } from "../../../features/dashboards/listing/dashboard-favourites.ts";

  const { chat, personalCanvases, agents } = featureFlags;

  const runtimeClient = useRuntimeClient();

  $: ({
    params: { organization, project },
  } = $page);

  // Query the instance to get the project display name
  $: instanceQuery = createRuntimeServiceGetInstance(runtimeClient, {});
  $: projectDisplayName =
    $instanceQuery.data?.instance?.projectDisplayName || project;
  $: isLoadingDisplayName = $instanceQuery.isLoading;
  $: isErrorDisplayName = $instanceQuery.isError;

  $: personalCanvasesQuery = getPersonalFilteredResources(
    runtimeClient,
    organization,
    project,
    ResourceKind.Canvas,
  );
  $: hasNoPersonalCanvases =
    !$personalCanvasesQuery.isPending &&
    ($personalCanvasesQuery.data?.length ?? 0) === 0;

  const selectedTagsStore = UrlParamsState.createStringArrayParam("tags");
  const sortStore = UrlParamsState.createStringParam(
    "sort",
    DashboardTableSortOptions[0].value,
  );
</script>

<svelte:head>
  <title>{m.home_project_page_title({ projectName: projectDisplayName })}</title
  >
</svelte:head>

<ContentContainer maxWidth={900}>
  <div class="flex flex-col gap-y-8 py-12">
    <!-- Kairos: canvases annotated with kairos_home (executive numbers strip,
         top edge per the executive cockpit layout) -->
    <KairosHomeCanvas />

    <!-- Welcome Section with Chat Input -->
    <div class="flex flex-col gap-y-6">
      <div class="flex flex-col gap-y-4">
        {#if isLoadingDisplayName}
          <DelayedContent visible={isLoadingDisplayName}>
            <div class="h-11 w-96 animate-pulse rounded bg-gray-200"></div>
          </DelayedContent>
        {:else if isErrorDisplayName}
          <h1
            class="text-3xl sm:text-4xl font-semibold text-fg-secondary break-words"
            aria-label={m.home_project_title_label()}
          >
            {@html m.home_welcome_to({
              projectName: `<span class="text-accent-primary-action">${escapeHtml(project)}</span>`,
            })}
          </h1>
        {:else}
          <h1
            class="text-3xl sm:text-4xl font-semibold text-fg-secondary break-words"
            aria-label={m.home_project_title_label()}
          >
            {@html m.home_welcome_to({
              projectName: `<span class="text-accent-primary-action">${escapeHtml(projectDisplayName)}</span>`,
            })}
          </h1>
        {/if}
        <p class="text-lg text-fg-muted">
          {#if $chat}
            {m.home_subtitle_with_chat()}
          {:else}
            {m.home_subtitle_no_chat()}
          {/if}
        </p>
      </div>

      <!-- Chat Input -->
      {#if $chat}
        <div class="w-full">
          <InlineChat noMargin height="110px" />
        </div>
      {/if}
    </div>

    <!-- Act: runs waiting on approval (after chat, before dashboards) -->
    {#if $agents}
      <HomePendingApprovals {organization} {project} />
    {/if}

    <!-- Kairos: recently triggered alerts (awareness feed, after the
         actionable approvals) -->
    <HomeRecentAlerts {organization} {project} />

    <!-- Dashboards Section -->
    {#if $personalCanvases}
      <PersonalCanvasesList org={organization} {project} />
    {/if}

    <div class="flex flex-col gap-y-4">
      <h2
        class="flex flex-row gap-x-2 items-center text-xl font-semibold text-fg-secondary"
      >
        <div class="flex flex-row w-full gap-x-2 items-center grow">
          <span>{m.home_dashboards_heading()}</span>
          <TableToolbarSort
            {sortStore}
            sortOptions={DashboardTableSortOptions}
            size="sm"
            noOutline
          />
          <div class="grow"></div>
          <DashboardsTagFilter align="end" {selectedTagsStore} />
        </div>
        {#if $personalCanvases && hasNoPersonalCanvases}
          <CreatePersonalCanvasDialog org={organization} {project} />
        {/if}
      </h2>
      <DashboardsTable isPreview previewLimit={5} />
    </div>
  </div>
</ContentContainer>

<script lang="ts">
  import * as m from "@rilldata/web-common/paraglide/messages.js";
  import { onNavigate } from "$app/navigation";
  import {
    DashboardBannerID,
    DashboardBannerPriority,
  } from "@rilldata/web-common/components/banner/constants";
  import ErrorPage from "@rilldata/web-common/components/ErrorPage.svelte";
  import { Dashboard } from "@rilldata/web-common/features/dashboards";
  import DashboardBuilding from "@rilldata/web-common/features/dashboards/DashboardBuilding.svelte";
  import { resetSelectedMockUserAfterNavigate } from "@rilldata/web-common/features/dashboards/granular-access-policies/resetSelectedMockUserAfterNavigate";
  import { selectedMockUserStore } from "@rilldata/web-common/features/dashboards/granular-access-policies/stores";
  import DashboardStateManager from "@rilldata/web-common/features/dashboards/state-managers/loaders/DashboardStateManager.svelte";
  import StateManagersProvider from "@rilldata/web-common/features/dashboards/state-managers/StateManagersProvider.svelte";
  import { useProjectParser } from "@rilldata/web-common/features/entity-management/resource-selectors";
  import {
    useExploreWithPolling,
    isExploreReconcilingForFirstTime,
    isExploreErrored,
  } from "@rilldata/web-common/features/explores/selectors";
  import {
    extractErrorStatusCode,
    isNotFoundError,
  } from "@rilldata/web-common/lib/errors";
  import { eventBus } from "@rilldata/web-common/lib/event-bus/event-bus";
  import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";
  import { previewModeStore } from "@rilldata/web-common/layout/preview-mode-store";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import type { PageData } from "./$types";
  import { featureFlags } from "@rilldata/web-common/features/feature-flags.ts";

  const runtimeClient = useRuntimeClient();

  export let data: PageData;
  $: ({ exploreName } = data);

  const { disablePersistentDashboardState } = featureFlags;

  resetSelectedMockUserAfterNavigate(queryClient, runtimeClient);

  $: exploreResource = useExploreWithPolling(runtimeClient, exploreName);

  $: validSpec = $exploreResource.data?.explore?.explore?.state?.validSpec;
  $: metricsViewName = $exploreResource.data?.metricsView?.meta?.name
    ?.name as string;
  $: measures = validSpec?.measures ?? [];

  $: filePaths = [
    ...($exploreResource.data?.explore?.meta?.filePaths ?? []),
    ...($exploreResource.data?.metricsView?.meta?.filePaths ?? []),
  ];

  $: projectParserQuery = useProjectParser(queryClient, runtimeClient, {
    enabled: $selectedMockUserStore?.admin,
  });

  $: hasBanner = !!validSpec?.banner;

  $: if (hasBanner) {
    eventBus.emit("add-banner", {
      id: DashboardBannerID,
      priority: DashboardBannerPriority,
      message: {
        type: "default",
        message: validSpec?.banner ?? "",
        iconType: "alert",
      },
    });
  }

  $: dashboardFileHasParseError =
    $projectParserQuery.data?.projectParser?.state?.parseErrors?.filter(
      (error) => filePaths.includes(error.filePath as string),
    );

  $: isDashboardNotFound =
    !$exploreResource.data &&
    $exploreResource.isError &&
    isNotFoundError($exploreResource.error);

  $: mockUserHasNoAccess =
    $selectedMockUserStore && isNotFoundError($exploreResource.error);

  $: homeHref = $previewModeStore ? "/dashboards" : "/";

  onNavigate(({ from, to }) => {
    const changedDashboard =
      !from || !to || from?.params?.name !== to?.params?.name;
    // Clear out any dashboard banners
    if (hasBanner && changedDashboard) {
      eventBus.emit("remove-banner", DashboardBannerID);
    }
  });
</script>

<svelte:head>
  <title>Rill Developer | {exploreName}</title>
</svelte:head>

{#if $exploreResource.isPending && !$exploreResource.data}
  <DashboardBuilding />
{:else if mockUserHasNoAccess}
  <ErrorPage
    statusCode={extractErrorStatusCode($exploreResource.error)}
    header={m.explore_user_no_access()}
    body={m.explore_security_policy_warning({ email: $selectedMockUserStore?.email ?? "" })}
    href={homeHref}
  />
{:else if isDashboardNotFound}
  <ErrorPage statusCode={404} header={m.explore_dashboard_not_found()} href={homeHref} />
{:else if $exploreResource.isSuccess}
  {#if isExploreReconcilingForFirstTime($exploreResource.data)}
    <DashboardBuilding />
  {:else if isExploreErrored($exploreResource.data)}
    <ErrorPage
      header={m.explore_error_building()}
      body={$exploreResource.data?.explore?.meta?.reconcileError ??
        m.explore_unknown_build_error()}
      href={homeHref}
    />
  {:else if dashboardFileHasParseError && dashboardFileHasParseError.length > 0}
    <ErrorPage
      header={m.explore_error_parsing()}
      body={m.explore_check_yaml_errors()}
      href={homeHref}
    />
  {:else if measures.length === 0 && $selectedMockUserStore !== null}
    <ErrorPage
      statusCode={extractErrorStatusCode($exploreResource.error)}
      header={m.explore_error_fetching()}
      body={m.explore_no_measures()}
      href={homeHref}
    />
  {:else if metricsViewName}
    <div class="h-full overflow-hidden">
      {#key exploreName}
        <StateManagersProvider {metricsViewName} {exploreName}>
          <DashboardStateManager
            {exploreName}
            disableMostRecentDashboardState={$disablePersistentDashboardState}
            disableInitSessionDashboardState={$disablePersistentDashboardState}
          >
            <Dashboard {metricsViewName} {exploreName} />
          </DashboardStateManager>
        </StateManagersProvider>
      {/key}
    </div>
  {/if}
{/if}

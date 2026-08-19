<script lang="ts">
  import {
    createAdminServiceListNotificationPreferences,
    createAdminServiceUpdateNotificationPreferences,
    getAdminServiceListNotificationPreferencesQueryKey,
    type V1ListNotificationPreferencesResponse,
  } from "@rilldata/web-admin/client";
  import SettingsContainer from "@rilldata/web-admin/features/organizations/settings/SettingsContainer.svelte";
  import Label from "@rilldata/web-common/components/forms/Label.svelte";
  import Switch from "@rilldata/web-common/components/forms/Switch.svelte";
  import DelayedCircleOutlineSpinner from "@rilldata/web-common/components/spinner/DelayedCircleOutlineSpinner.svelte";
  import { eventBus } from "@rilldata/web-common/lib/event-bus/event-bus";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";
  import {
    listOrganizationPreferences,
    type NotificationCategory,
    type NotificationPreferences,
  } from "./push-utils";

  const preferencesQuery = createAdminServiceListNotificationPreferences();
  const updatePreferences = createAdminServiceUpdateNotificationPreferences();

  let { data, isSuccess } = $derived($preferencesQuery);

  let organizations = $derived(
    listOrganizationPreferences(data?.organizations),
  );
  // With a single organization there is nothing to tell its switches apart
  // from, so the block goes without a heading.
  let showOrganizationNames = $derived(organizations.length > 1);

  // The switch waiting for the server, so a toggle only spins its own.
  let switching = $state<string | undefined>(undefined);

  let categories = $derived([
    {
      key: "pushAlerts" as NotificationCategory,
      label: m.notifications_category_alerts(),
      description: m.notifications_category_alerts_description(),
    },
    {
      key: "pushReports" as NotificationCategory,
      label: m.notifications_category_reports(),
      description: m.notifications_category_reports_description(),
    },
    {
      key: "pushActApprovals" as NotificationCategory,
      label: m.notifications_category_act_approvals(),
      description: m.notifications_category_act_approvals_description(),
    },
  ]);

  async function toggle(
    org: string,
    preferences: NotificationPreferences,
    category: NotificationCategory,
  ) {
    const next = { ...preferences, [category]: !preferences[category] };
    switching = `${org}/${category}`;
    try {
      const response = await $updatePreferences.mutateAsync({
        org,
        data: { preferences: next },
      });
      // Seed the cache with what the server stored instead of refetching, so
      // the switch never flicks back to its old position on the way.
      queryClient.setQueryData<V1ListNotificationPreferencesResponse>(
        getAdminServiceListNotificationPreferencesQueryKey(),
        (previous) =>
          previous && {
            organizations: (previous.organizations ?? []).map((organization) =>
              organization.org === org
                ? { ...organization, preferences: response.preferences }
                : organization,
            ),
          },
      );
    } catch {
      eventBus.emit("notification", {
        type: "error",
        message: m.notifications_categories_error(),
      });
    } finally {
      switching = undefined;
    }
  }
</script>

<SettingsContainer title={m.notifications_categories_title()}>
  <p>{m.notifications_categories_description()}</p>

  {#if isSuccess && organizations.length === 0}
    <p class="mt-3">{m.notifications_categories_empty()}</p>
  {:else if isSuccess}
    <div class="organizations">
      {#each organizations as { org, label, preferences } (org)}
        <div class="organization">
          {#if showOrganizationNames}
            <h3 class="organization-name">{label}</h3>
          {/if}
          <div class="categories">
            {#each categories as category (category.key)}
              <div class="category">
                <div class="flex flex-col">
                  <Label
                    for="notification-category-{org}-{category.key}"
                    class="font-medium text-fg-primary"
                  >
                    {category.label}
                  </Label>
                  <span>{category.description}</span>
                </div>
                <div class="ml-auto shrink-0">
                  <DelayedCircleOutlineSpinner
                    isLoading={switching === `${org}/${category.key}`}
                  >
                    <Switch
                      id="notification-category-{org}-{category.key}"
                      label={category.label}
                      checked={preferences[category.key]}
                      onclick={() => toggle(org, preferences, category.key)}
                    />
                  </DelayedCircleOutlineSpinner>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</SettingsContainer>

<style lang="postcss">
  .organizations {
    @apply flex flex-col gap-y-5;
  }

  .organization-name {
    @apply text-sm font-medium text-fg-primary;
  }

  .categories {
    @apply flex flex-col gap-y-3 mt-3;
  }

  .category {
    @apply flex flex-row items-center gap-x-4;
  }
</style>

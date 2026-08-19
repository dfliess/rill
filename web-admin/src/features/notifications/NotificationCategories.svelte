<script lang="ts">
  import {
    createAdminServiceGetNotificationPreferences,
    createAdminServiceUpdateNotificationPreferences,
    getAdminServiceGetNotificationPreferencesQueryKey,
  } from "@rilldata/web-admin/client";
  import SettingsContainer from "@rilldata/web-admin/features/organizations/settings/SettingsContainer.svelte";
  import Label from "@rilldata/web-common/components/forms/Label.svelte";
  import Switch from "@rilldata/web-common/components/forms/Switch.svelte";
  import DelayedCircleOutlineSpinner from "@rilldata/web-common/components/spinner/DelayedCircleOutlineSpinner.svelte";
  import { eventBus } from "@rilldata/web-common/lib/event-bus/event-bus";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";
  import {
    readNotificationPreferences,
    type NotificationCategory,
  } from "./push-utils";

  const preferencesQuery = createAdminServiceGetNotificationPreferences();
  const updatePreferences = createAdminServiceUpdateNotificationPreferences();

  let { data, isSuccess } = $derived($preferencesQuery);
  let { isPending } = $derived($updatePreferences);

  let preferences = $derived(readNotificationPreferences(data?.preferences));

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

  async function toggle(category: NotificationCategory) {
    const next = { ...preferences, [category]: !preferences[category] };
    try {
      const response = await $updatePreferences.mutateAsync({
        data: { preferences: next },
      });
      // Seed the cache with what the server stored instead of refetching, so
      // the switch never flicks back to its old position on the way.
      queryClient.setQueryData(
        getAdminServiceGetNotificationPreferencesQueryKey(),
        { preferences: response.preferences },
      );
    } catch {
      eventBus.emit("notification", {
        type: "error",
        message: m.notifications_categories_error(),
      });
    }
  }
</script>

<SettingsContainer title={m.notifications_categories_title()}>
  <p>{m.notifications_categories_description()}</p>

  {#if isSuccess}
    <div class="categories">
      {#each categories as { key, label, description } (key)}
        <div class="category">
          <div class="flex flex-col">
            <Label
              for="notification-category-{key}"
              class="font-medium text-fg-primary"
            >
              {label}
            </Label>
            <span>{description}</span>
          </div>
          <div class="ml-auto shrink-0">
            <DelayedCircleOutlineSpinner isLoading={isPending}>
              <Switch
                id="notification-category-{key}"
                {label}
                checked={preferences[key]}
                onclick={() => toggle(key)}
              />
            </DelayedCircleOutlineSpinner>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</SettingsContainer>

<style lang="postcss">
  .categories {
    @apply flex flex-col gap-y-3 mt-3;
  }

  .category {
    @apply flex flex-row items-center gap-x-4;
  }
</style>

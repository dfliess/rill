<script lang="ts">
  import {
    createAdminServiceDeletePushSubscription,
    createAdminServiceListPushSubscriptions,
    getAdminServiceListPushSubscriptionsQueryKey,
  } from "@rilldata/web-admin/client";
  import SettingsContainer from "@rilldata/web-admin/features/organizations/settings/SettingsContainer.svelte";
  import { Button } from "@rilldata/web-common/components/button";
  import DelayedCircleOutlineSpinner from "@rilldata/web-common/components/spinner/DelayedCircleOutlineSpinner.svelte";
  import { eventBus } from "@rilldata/web-common/lib/event-bus/event-bus";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { getLocale } from "@rilldata/web-common/lib/i18n/gen/runtime";
  import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";
  import { listPushDevices } from "./push-utils";

  let {
    currentEndpoint,
    onRemoveCurrentDevice,
  }: {
    currentEndpoint: string | undefined;
    // Removing the row of the browser looking at the page also has to drop its
    // browser-side subscription, which only the parent knows how to do.
    onRemoveCurrentDevice: () => Promise<void>;
  } = $props();

  const subscriptionsQuery = createAdminServiceListPushSubscriptions();
  const deleteSubscription = createAdminServiceDeletePushSubscription();

  let { data, isSuccess } = $derived($subscriptionsQuery);
  let devices = $derived(listPushDevices(data?.subscriptions, currentEndpoint));
  let removingId = $state<string | undefined>(undefined);

  function formatDate(value: string | undefined): string {
    if (!value) return "";
    return new Date(value).toLocaleDateString(getLocale(), {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  }

  async function remove(id: string, isCurrent: boolean) {
    removingId = id;
    try {
      if (isCurrent) {
        await onRemoveCurrentDevice();
        return;
      }
      await $deleteSubscription.mutateAsync({ id });
      await queryClient.invalidateQueries({
        queryKey: getAdminServiceListPushSubscriptionsQueryKey(),
      });
      eventBus.emit("notification", {
        message: m.notifications_devices_removed_toast(),
      });
    } catch {
      eventBus.emit("notification", {
        type: "error",
        message: m.notifications_devices_remove_error(),
      });
    } finally {
      removingId = undefined;
    }
  }
</script>

<!-- The list itself says what the title promises, so there is no help text: only the empty state, which
     is the one moment the section has nothing to show for itself. -->
<SettingsContainer title={m.notifications_devices_title()}>
  {#if isSuccess}
    {#if devices.length === 0}
      <p>{m.notifications_devices_empty()}</p>
    {:else}
      <ul class="devices">
        {#each devices as { id, label, createdOn, isCurrent } (id)}
          <li class="device">
            <div class="flex flex-col min-w-0">
              <span class="device-name">
                {label || m.notifications_devices_unknown()}
                {#if isCurrent}
                  <span class="device-badge">
                    {m.notifications_devices_this_device()}
                  </span>
                {/if}
              </span>
              {#if createdOn}
                <span>
                  {m.notifications_devices_added_on({
                    date: formatDate(createdOn),
                  })}
                </span>
              {/if}
            </div>
            <div class="ml-auto shrink-0">
              <DelayedCircleOutlineSpinner isLoading={removingId === id}>
                <Button type="text" onClick={() => remove(id, isCurrent)}>
                  {m.notifications_devices_remove()}
                </Button>
              </DelayedCircleOutlineSpinner>
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  {/if}
</SettingsContainer>

<style lang="postcss">
  .devices {
    @apply flex flex-col mt-3;
  }

  .device {
    @apply flex flex-row items-center gap-x-4 py-2 border-b last:border-b-0;
  }

  .device-name {
    @apply flex flex-row items-center gap-x-2 text-fg-primary font-medium truncate;
  }

  .device-badge {
    @apply px-1.5 py-0.5 rounded-sm bg-surface-subtle text-fg-tertiary text-xs font-normal;
  }
</style>

<script lang="ts">
  import {
    createAdminServiceCreatePushSubscription,
    createAdminServiceDeletePushSubscription,
    createAdminServiceGetPushNotificationConfig,
    createAdminServiceListPushSubscriptions,
    getAdminServiceListPushSubscriptionsQueryKey,
  } from "@rilldata/web-admin/client";
  import SettingsContainer from "@rilldata/web-admin/features/organizations/settings/SettingsContainer.svelte";
  import Switch from "@rilldata/web-common/components/forms/Switch.svelte";
  import DelayedCircleOutlineSpinner from "@rilldata/web-common/components/spinner/DelayedCircleOutlineSpinner.svelte";
  import { eventBus } from "@rilldata/web-common/lib/event-bus/event-bus";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";
  import { onMount } from "svelte";
  import NotificationCategories from "./NotificationCategories.svelte";
  import PushDevicesList from "./PushDevicesList.svelte";
  import {
    currentPushEndpoint,
    ensurePushSubscription,
    pushManagerOfActiveWorker,
    unsubscribeThisDevice,
  } from "./push-subscription";
  import {
    browserSupportsPush,
    derivePushSectionState,
    describeUserAgent,
    subscriptionCredentials,
  } from "./push-utils";

  const configQuery = createAdminServiceGetPushNotificationConfig();
  const subscriptionsQuery = createAdminServiceListPushSubscriptions();
  const createSubscription = createAdminServiceCreatePushSubscription();
  const deleteSubscription = createAdminServiceDeletePushSubscription();

  const supported = browserSupportsPush();

  let permission = $state<NotificationPermission>("default");
  // Endpoint of the subscription this browser holds, if any. Read once on
  // mount and kept in step with what this page does to it.
  let currentEndpoint = $state<string | undefined>(undefined);
  let switching = $state(false);

  let vapidPublicKey = $derived($configQuery.data?.vapidPublicKey ?? "");
  let sectionState = $derived(
    derivePushSectionState({
      configLoaded: $configQuery.isSuccess,
      vapidPublicKey,
      supported,
      permission,
    }),
  );

  // The switch is on only when the browser holds a subscription and the server
  // still has the matching row: the device can be removed from another one.
  let enabledHere = $derived(
    !!currentEndpoint &&
      !!$subscriptionsQuery.data?.subscriptions?.some(
        (subscription) => subscription.endpoint === currentEndpoint,
      ),
  );

  onMount(() => {
    if (!supported) return;
    permission = Notification.permission;
    void currentPushEndpoint().then((endpoint) => {
      currentEndpoint = endpoint;
    });
  });

  async function enableHere() {
    switching = true;
    try {
      // Asking has to happen inside the click: browsers refuse a permission
      // prompt that no gesture asked for.
      permission = await Notification.requestPermission();
      if (permission !== "granted") return;

      const manager = await pushManagerOfActiveWorker();
      const subscription = await ensurePushSubscription(
        manager,
        vapidPublicKey,
      );
      const credentials = subscriptionCredentials(subscription);
      if (!credentials) {
        throw new Error("push subscription without encryption keys");
      }

      await $createSubscription.mutateAsync({
        data: {
          ...credentials,
          userAgent: describeUserAgent(navigator.userAgent),
        },
      });
      currentEndpoint = subscription.endpoint;
      await queryClient.invalidateQueries({
        queryKey: getAdminServiceListPushSubscriptionsQueryKey(),
      });
      eventBus.emit("notification", {
        type: "success",
        message: m.notifications_device_enabled_toast(),
      });
    } catch {
      eventBus.emit("notification", {
        type: "error",
        message: m.notifications_device_enable_error(),
      });
    } finally {
      switching = false;
    }
  }

  async function disableHere() {
    switching = true;
    try {
      const registered = $subscriptionsQuery.data?.subscriptions?.find(
        (subscription) => subscription.endpoint === currentEndpoint,
      );
      if (registered?.id) {
        await $deleteSubscription.mutateAsync({ id: registered.id });
      }
      // Dropping the browser-side subscription too, so the push service stops
      // holding an endpoint nothing can deliver to.
      await unsubscribeThisDevice();
      currentEndpoint = undefined;
      await queryClient.invalidateQueries({
        queryKey: getAdminServiceListPushSubscriptionsQueryKey(),
      });
      eventBus.emit("notification", {
        message: m.notifications_device_disabled_toast(),
      });
    } catch {
      eventBus.emit("notification", {
        type: "error",
        message: m.notifications_device_disable_error(),
      });
    } finally {
      switching = false;
    }
  }
</script>

{#if sectionState === "disabled"}
  <p class="text-sm text-fg-tertiary">{m.notifications_unavailable()}</p>
{:else if sectionState !== "loading"}
  <SettingsContainer title={m.notifications_device_title()}>
    <div class="flex flex-row items-center gap-x-4">
      <p>{m.notifications_device_description()}</p>
      <div class="ml-auto shrink-0">
        <DelayedCircleOutlineSpinner isLoading={switching}>
          <Switch
            id="push-on-this-device"
            label={m.notifications_device_title()}
            checked={enabledHere}
            disabled={sectionState !== "ready"}
            onclick={() => (enabledHere ? disableHere() : enableHere())}
          />
        </DelayedCircleOutlineSpinner>
      </div>
    </div>

    {#if sectionState === "denied"}
      <p class="mt-3 text-fg-secondary">{m.notifications_device_denied()}</p>
    {:else if sectionState === "unsupported"}
      <p class="mt-3 text-fg-secondary">
        {m.notifications_device_unsupported()}
      </p>
    {/if}
  </SettingsContainer>

  <NotificationCategories />

  <PushDevicesList {currentEndpoint} onRemoveCurrentDevice={disableHere} />
{/if}

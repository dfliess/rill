<!-- USER NOTIFICATION SETTINGS -->

<script lang="ts">
  import NotificationsSettings from "@rilldata/web-admin/features/notifications/NotificationsSettings.svelte";
  import ContentContainer from "@rilldata/web-common/components/layout/ContentContainer.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
</script>

<ContentContainer title={m.notifications_page_title()} maxWidth={800}>
  <div class="flex flex-col gap-y-6 pt-6">
    <NotificationsSettings />
  </div>
</ContentContainer>

<script lang="ts">
  import { createAdminServiceListNotifications } from "@rilldata/web-admin/client";
  import SettingsContainer from "@rilldata/web-admin/features/organizations/settings/SettingsContainer.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { getLocale } from "@rilldata/web-common/lib/i18n/gen/runtime";
  import { listNotificationEntries } from "./notification-history";

  const historyQuery = createAdminServiceListNotifications();

  let { data, isSuccess, isError } = $derived($historyQuery);
  // Whether the history holds more than this list shows. A list that quietly stops reads as "this is everything".
  let hasMore = $derived(data?.hasMore === true);
  let entries = $derived(listNotificationEntries(data?.notifications));

  // A history is read for "when", so the time is part of the answer, not a detail: two reports on the same
  // morning are otherwise indistinguishable.
  function formatSentOn(value: string | undefined): string {
    if (!value) return "";
    return new Date(value).toLocaleString(getLocale(), {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }
</script>

<!-- No help text: the list says what the title promises. The empty state is the one moment this has nothing to
     show for itself, and the error state is worth spelling out — rendering nothing there would read as "you have
     been notified about nothing", which is a different fact. -->
<SettingsContainer title={m.notifications_history_title()}>
  {#if isError}
    <p>{m.notifications_history_error()}</p>
  {:else if isSuccess}
    {#if entries.length === 0}
      <p>{m.notifications_history_empty()}</p>
    {:else}
      <ul class="entries">
        {#each entries as { id, categoryKey, scope, title, body, linkPath, sentOn } (id)}
          <li class="entry">
            <div class="flex flex-col min-w-0 gap-y-0.5">
              <span class="entry-title">
                {#if linkPath}
                  <a href={linkPath}>{title}</a>
                {:else}
                  {title}
                {/if}
              </span>
              {#if body}
                <span class="truncate">{body}</span>
              {/if}
              <span class="entry-meta">
                <span class="entry-category">{m[categoryKey]()}</span>
                {#if scope}<span>{scope}</span>{/if}
                <span>{formatSentOn(sentOn)}</span>
              </span>
            </div>
          </li>
        {/each}
      </ul>
      {#if hasMore}
        <p class="mt-2">{m.notifications_history_more()}</p>
      {/if}
    {/if}
  {/if}
</SettingsContainer>

<style lang="postcss">
  .entries {
    @apply flex flex-col mt-3;
  }

  .entry {
    @apply flex flex-row items-center gap-x-4 py-2 border-b last:border-b-0;
  }

  .entry-title {
    @apply text-fg-primary font-medium truncate;
  }

  .entry-title a {
    @apply hover:text-primary-600;
  }

  .entry-meta {
    @apply flex flex-row items-center gap-x-2 text-fg-tertiary text-xs;
  }

  .entry-category {
    @apply px-1.5 py-0.5 rounded-sm bg-surface-subtle;
  }
</style>

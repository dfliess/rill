// Pure helpers behind the notification history block (NotificationHistory.svelte).
//
// The history is a record of what concerned you, not an inbox: there is no read state and nothing to dismiss.
// An alert or a report has nothing to resolve, so "mark as read" would be housekeeping done for the software's
// benefit; and the app icon already carries a number, which counts real pending work (kairos-cloud#147).

import type { V1Notification } from "@rilldata/web-admin/client";

// One row of the history.
export type NotificationEntry = {
  id: string;
  // Key of the message naming the category, so the template has no switch of its own.
  categoryKey: NotificationCategoryKey;
  // Where it came from, as "organization/project".
  scope: string;
  title: string;
  body: string;
  // Path to open, or nothing: a row whose link went missing is still a record of what happened, just not a
  // destination.
  linkPath: string | undefined;
  sentOn: string | undefined;
};

export type NotificationCategoryKey =
  | "notifications_category_alerts"
  | "notifications_category_reports"
  | "notifications_category_act_approvals";

// The categories the server sends, mapped to the messages that already name them for the preference switches.
const CATEGORY_KEYS: Record<string, NotificationCategoryKey> = {
  alerts: "notifications_category_alerts",
  reports: "notifications_category_reports",
  act_approvals: "notifications_category_act_approvals",
};

// listNotificationEntries turns the server's history into rows, keeping its ordering (newest first).
//
// A row with no id, or in a category this build has no name for, is dropped rather than rendered as a blank or as
// its raw wire value: an older frontend meeting a newer category should show one entry fewer, not a row reading
// "act_approvals".
export function listNotificationEntries(
  notifications: V1Notification[] | undefined,
): NotificationEntry[] {
  const entries: NotificationEntry[] = [];
  for (const notification of notifications ?? []) {
    const categoryKey = CATEGORY_KEYS[notification.category ?? ""];
    if (!notification.id || !categoryKey) continue;
    entries.push({
      id: notification.id,
      categoryKey,
      scope: [notification.org, notification.project].filter(Boolean).join("/"),
      title: notification.title ?? "",
      body: notification.body ?? "",
      linkPath: notification.linkPath || undefined,
      sentOn: notification.sentOn,
    });
  }
  return entries;
}

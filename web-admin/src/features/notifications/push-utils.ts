// Pure helpers behind the push notification settings section
// (NotificationsSettings.svelte). Everything here is browser-independent so
// it can be unit-tested; the callers pass in what they read from the DOM.

import type {
  V1NotificationPreferences,
  V1OrganizationNotificationPreferences,
  V1PushSubscription,
} from "@rilldata/web-admin/client";

export function browserSupportsPush(): boolean {
  return (
    typeof window !== "undefined" &&
    "serviceWorker" in navigator &&
    "PushManager" in window &&
    "Notification" in window
  );
}

// iOS exposes the Push API only to a web app launched from the Home Screen, never to a Safari tab
// (16.4+). So on iPhone and iPad the missing API is not a browser that cannot do this: it is one asking
// to be installed first, which is an instruction, not a dead end.
export function isIOS(): boolean {
  if (typeof navigator === "undefined") return false;
  // iPadOS 13+ reports itself as a Mac, and the touch points are what tells the two apart.
  return (
    /iPad|iPhone|iPod/.test(navigator.userAgent) ||
    (navigator.userAgent.includes("Macintosh") && navigator.maxTouchPoints > 1)
  );
}

// Whether the page is running as an installed web app rather than inside browser chrome. iOS answers
// through a non-standard flag of its own; everyone else through the display mode.
export function isStandalone(): boolean {
  if (typeof window === "undefined") return false;
  const iosStandalone = (navigator as { standalone?: boolean }).standalone;
  return (
    iosStandalone === true ||
    window.matchMedia?.("(display-mode: standalone)").matches === true
  );
}

// UI state of the push settings section, in order of precedence:
// "loading" until the server config arrives; "disabled" when the deployment
// has no VAPID key configured (the section renders nothing at all);
// "install-required" on an iOS browser tab, where the API only appears once the app is on the Home
// Screen; "unsupported" when the browser lacks the Push API for good; "denied" when the user
// blocked notifications for this site; otherwise "ready".
export type PushSectionState =
  | "loading"
  | "disabled"
  | "install-required"
  | "unsupported"
  | "denied"
  | "ready";

export function derivePushSectionState(args: {
  configLoaded: boolean;
  vapidPublicKey: string;
  supported: boolean;
  permission: NotificationPermission;
  ios: boolean;
  standalone: boolean;
}): PushSectionState {
  if (!args.configLoaded) return "loading";
  if (!args.vapidPublicKey) return "disabled";
  if (!args.supported) {
    // An installed iOS app with no Push API is a version older than 16.4, which installing again will
    // not fix, so that one is genuinely unsupported.
    return args.ios && !args.standalone ? "install-required" : "unsupported";
  }
  if (args.permission === "denied") return "denied";
  return "ready";
}

// The Push API wants `applicationServerKey` as raw bytes, but VAPID public
// keys travel as base64url strings (RFC 7515 §2: "-" and "_" instead of "+"
// and "/", no padding). The buffer is pinned to `ArrayBuffer` because
// `PushManager.subscribe` refuses a view that might sit on a SharedArrayBuffer.
export function urlBase64ToUint8Array(
  base64Url: string,
): Uint8Array<ArrayBuffer> {
  const padding = "=".repeat((4 - (base64Url.length % 4)) % 4);
  const base64 = (base64Url + padding).replace(/-/g, "+").replace(/_/g, "/");
  const raw = atob(base64);
  const output = new Uint8Array(raw.length);
  for (let i = 0; i < raw.length; i++) {
    output[i] = raw.charCodeAt(i);
  }
  return output;
}

// Human-readable device name stored server-side with each subscription and
// shown in the device list. It is data shared across users and locales, so it
// deliberately avoids translatable copy. Order matters in both lists: Edge
// and Opera UAs also contain "Chrome", every Chrome UA contains "Safari",
// and Android UAs contain "Linux".
const UA_BROWSERS: [RegExp, string][] = [
  [/edg(?:e|a|ios)?\//i, "Edge"],
  [/opr\//i, "Opera"],
  [/firefox|fxios/i, "Firefox"],
  [/samsungbrowser/i, "Samsung Internet"],
  [/chrome|crios/i, "Chrome"],
  [/safari/i, "Safari"],
];
const UA_OSES: [RegExp, string][] = [
  [/windows/i, "Windows"],
  [/android/i, "Android"],
  [/iphone|ipad|ipod/i, "iOS"],
  [/mac os x|macintosh/i, "macOS"],
  [/linux/i, "Linux"],
];

export function describeUserAgent(userAgent: string): string {
  const browser = UA_BROWSERS.find(([re]) => re.test(userAgent))?.[1];
  const os = UA_OSES.find(([re]) => re.test(userAgent))?.[1];
  if (browser && os) return `${browser} · ${os}`;
  return browser ?? os ?? userAgent.slice(0, 80);
}

// Extracts the fields CreatePushSubscription needs from a browser
// PushSubscription. Returns null if the browser did not hand out the
// encryption keys; a subscription without them cannot receive payloads.
export function subscriptionCredentials(subscription: {
  toJSON(): PushSubscriptionJSON;
}): { endpoint: string; p256dh: string; auth: string } | null {
  const json = subscription.toJSON();
  const p256dh = json.keys?.p256dh;
  const auth = json.keys?.auth;
  if (!json.endpoint || !p256dh || !auth) return null;
  return { endpoint: json.endpoint, p256dh, auth };
}

// The three per-category switches, keyed by the field they set on the server.
export type NotificationCategory =
  | "pushAlerts"
  | "pushReports"
  | "pushActApprovals";
export type NotificationPreferences = Record<NotificationCategory, boolean>;

// Proto3 JSON omits false, so a flag missing from the response means "off",
// not "unset". Normalizing once keeps the switches from reading `undefined`
// as a third state.
export function readNotificationPreferences(
  preferences: V1NotificationPreferences | undefined,
): NotificationPreferences {
  return {
    pushAlerts: preferences?.pushAlerts ?? false,
    pushReports: preferences?.pushReports ?? false,
    pushActApprovals: preferences?.pushActApprovals ?? false,
  };
}

// One organization's block of category switches. The categories are chosen per
// organization, so a user in several of them gets one block each.
export type OrganizationPreferences = {
  // Organization name, as the update mutation addresses it.
  org: string;
  // What the block is titled with: the display name when the organization has
  // one, its name otherwise.
  label: string;
  preferences: NotificationPreferences;
};

// Turns the server's per-organization preferences into blocks for the settings
// page, keeping the admin's ordering by organization name.
export function listOrganizationPreferences(
  organizations: V1OrganizationNotificationPreferences[] | undefined,
): OrganizationPreferences[] {
  return (organizations ?? [])
    .filter(
      (
        organization,
      ): organization is V1OrganizationNotificationPreferences & {
        org: string;
      } => Boolean(organization.org),
    )
    .map((organization) => ({
      org: organization.org,
      label: organization.orgDisplayName || organization.org,
      preferences: readNotificationPreferences(organization.preferences),
    }));
}

// One row of the device list.
export type PushDevice = {
  id: string;
  // The descriptive name stored when the device subscribed; may be empty for
  // subscriptions created before the browser sent one.
  label: string;
  createdOn: string | undefined;
  isCurrent: boolean;
};

// Turns the server's subscriptions into device rows, with the browser looking
// at the page called out and listed first. Sorting by `createdOn` as a string
// works because the admin emits RFC 3339 in UTC, where lexical and
// chronological order agree.
export function listPushDevices(
  subscriptions: V1PushSubscription[] | undefined,
  currentEndpoint: string | undefined,
): PushDevice[] {
  return (subscriptions ?? [])
    .filter(
      (subscription): subscription is V1PushSubscription & { id: string } =>
        Boolean(subscription.id),
    )
    .map((subscription) => ({
      id: subscription.id,
      label: subscription.userAgent ?? "",
      createdOn: subscription.createdOn,
      isCurrent: !!currentEndpoint && subscription.endpoint === currentEndpoint,
    }))
    .sort((a, b) => {
      if (a.isCurrent !== b.isCurrent) return a.isCurrent ? -1 : 1;
      return (b.createdOn ?? "").localeCompare(a.createdOn ?? "");
    });
}

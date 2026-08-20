// Tapping a notification has to land the reader on the thing it was about. The service worker cannot do
// that on its own in an installed iOS web app, where WindowClient.navigate() is either missing or resolves
// without navigating, so it asks the page to route itself instead (see service-worker.ts).

import { goto } from "$app/navigation";

export const NOTIFICATION_NAVIGATE = "kairos:navigate";

// Reduces a URL from a notification payload to a path this app can route to, or nothing.
// Same-origin only: the payload is ours, but a bug on either side must not turn a notification tap into
// an open redirect.
export function notificationTargetPath(raw: unknown): string | undefined {
  // An empty link is not a destination: resolving it would send the reader to the home page, which is
  // not where the notification was about.
  if (!raw || typeof raw !== "string" || typeof window === "undefined") {
    return undefined;
  }
  let url: URL;
  try {
    url = new URL(raw, window.location.origin);
  } catch {
    return undefined;
  }
  if (url.origin !== window.location.origin) return undefined;
  return url.pathname + url.search + url.hash;
}

// Starts listening for the service worker's navigation requests. Returns the function that stops it.
export function listenForNotificationNavigation(): () => void {
  if (typeof navigator === "undefined" || !("serviceWorker" in navigator)) {
    return () => {};
  }

  const handler = (event: MessageEvent) => {
    const data: unknown = event.data;
    if (typeof data !== "object" || data === null) return;
    if ((data as { type?: unknown }).type !== NOTIFICATION_NAVIGATE) return;
    const path = notificationTargetPath((data as { url?: unknown }).url);
    if (path) void goto(path);
  };

  navigator.serviceWorker.addEventListener("message", handler);
  return () => navigator.serviceWorker.removeEventListener("message", handler);
}

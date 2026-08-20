/// <reference types="@sveltejs/kit" />

import { build, files, version } from "$service-worker";

// Typed locally instead of via `/// <reference lib="webworker" />` with
// no-default-lib: those directives apply to the whole tsc program, and the
// repo-root tsconfig compiles every workspace's .ts files in one program, so
// they silently dropped DOM.Iterable for unrelated files. Only the surface
// this worker actually uses is declared here; everything else it touches
// (caches, fetch, Request, Response, URL) is already in the DOM lib.
type ExtendableEvent = Event & { waitUntil(p: Promise<unknown>): void };
type FetchEvent = ExtendableEvent & {
  readonly request: Request;
  respondWith(r: Promise<Response> | Response): void;
};
type PushEvent = ExtendableEvent & {
  readonly data: { json(): unknown } | null;
};
type NotificationEvent = ExtendableEvent & {
  readonly notification: {
    close(): void;
    readonly data: unknown;
  };
};
type WindowClient = {
  focus(): Promise<unknown>;
  postMessage(message: unknown): void;
};
type ServiceWorkerScope = {
  addEventListener(
    type: "install" | "activate",
    listener: (event: ExtendableEvent) => void,
  ): void;
  addEventListener(type: "fetch", listener: (event: FetchEvent) => void): void;
  addEventListener(type: "push", listener: (event: PushEvent) => void): void;
  addEventListener(
    type: "notificationclick",
    listener: (event: NotificationEvent) => void,
  ): void;
  skipWaiting(): Promise<void>;
  registration: {
    showNotification(
      title: string,
      options?: { body?: string; tag?: string; icon?: string; data?: unknown },
    ): Promise<void>;
  };
  clients: {
    claim(): Promise<void>;
    matchAll(options?: {
      type?: "window";
      includeUncontrolled?: boolean;
    }): Promise<WindowClient[]>;
    openWindow(url: string): Promise<unknown>;
  };
  location: Location;
};

const sw = self as unknown as ServiceWorkerScope;

// Conservative caching strategy for a multi-tenant BI app:
// - HTML is NEVER cached: navigations always hit the network, so deploys
//   propagate on a normal reload and no "new version available" flow is needed.
// - API responses are NEVER intercepted (freshness + tenant isolation).
// - Only immutable build assets and static files are cached, on first use,
//   in a per-build cache that is purged on activate.
const CACHE = `kairos-shell-${version}`;

const OFFLINE_URL = "/offline.html";
const PRECACHE = [OFFLINE_URL, "/pwa-icon-192.png"];

// Cacheable static paths: hashed build assets + files in static/.
const STATIC_PATHS = new Set([...build, ...files]);

sw.addEventListener("install", (event) => {
  event.waitUntil(
    caches
      .open(CACHE)
      .then((cache) => cache.addAll(PRECACHE))
      .then(() => sw.skipWaiting()),
  );
});

sw.addEventListener("activate", (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) =>
        Promise.all(
          keys.filter((k) => k !== CACHE).map((k) => caches.delete(k)),
        ),
      )
      .then(() => sw.clients.claim()),
  );
});

async function cacheFirst(request: Request): Promise<Response> {
  const cache = await caches.open(CACHE);
  const hit = await cache.match(request);
  if (hit) return hit;
  const response = await fetch(request);
  if (response.ok) void cache.put(request, response.clone());
  return response;
}

async function networkWithOfflineFallback(request: Request): Promise<Response> {
  try {
    return await fetch(request);
  } catch {
    const offline = await caches.match(OFFLINE_URL);
    return offline ?? Response.error();
  }
}

sw.addEventListener("fetch", (event) => {
  const { request } = event;
  if (request.method !== "GET") return;

  const url = new URL(request.url);
  // Cross-origin requests (admin/runtime APIs on other subdomains) are
  // never intercepted.
  if (url.origin !== sw.location.origin) return;

  if (request.mode === "navigate") {
    event.respondWith(networkWithOfflineFallback(request));
    return;
  }

  if (STATIC_PATHS.has(url.pathname)) {
    event.respondWith(cacheFirst(request));
  }
  // Anything else (e.g. same-origin API calls) falls through to the network.
});

// Web Push: the admin sends a JSON payload of the shape
// {title, body, link, category, tag}; `link` is an absolute URL into this
// frontend. A malformed or non-JSON payload is dropped silently: showing a
// broken notification (or throwing) would be worse than showing none.
// The page listens for this to route itself; see navigate-from-notification.ts.
const NOTIFICATION_NAVIGATE = "kairos:navigate";

type PushPayload = {
  title: string;
  body?: string;
  link?: string;
  tag?: string;
};

function parsePushPayload(event: PushEvent): PushPayload | null {
  if (!event.data) return null;
  let raw: unknown;
  try {
    raw = event.data.json();
  } catch {
    return null;
  }
  if (typeof raw !== "object" || raw === null) return null;
  const { title, body, link, tag } = raw as Record<string, unknown>;
  if (typeof title !== "string" || !title) return null;
  return {
    title,
    body: typeof body === "string" ? body : undefined,
    link: typeof link === "string" ? link : undefined,
    tag: typeof tag === "string" ? tag : undefined,
  };
}

sw.addEventListener("push", (event) => {
  const payload = parsePushPayload(event);
  if (!payload) return;
  event.waitUntil(
    sw.registration.showNotification(payload.title, {
      body: payload.body,
      tag: payload.tag,
      icon: "/pwa-icon-192.png",
      data: { link: payload.link },
    }),
  );
});

// Focus an existing app window if there is one (matchAll only returns clients of this worker's origin)
// and let it route itself to the notification's link; otherwise open a new window on it.
async function openNotificationLink(link: string | undefined): Promise<void> {
  if (!link) return;

  const windowClients = await sw.clients.matchAll({
    type: "window",
    includeUncontrolled: true,
  });
  const client = windowClients[0];
  if (!client) {
    await sw.clients.openWindow(link);
    return;
  }

  try {
    await client.focus();
  } catch {
    // Focus can be refused (e.g. no transient activation); keep going.
  }
  // The page routes itself. An installed iOS web app has no usable navigate(): it is either missing or
  // resolves without doing anything, so the app came to the foreground still showing whatever page it
  // was on, with no error to fall back from. The page is loaded already, so this is also the faster
  // path everywhere else.
  client.postMessage({ type: NOTIFICATION_NAVIGATE, url: link });
}

sw.addEventListener("notificationclick", (event) => {
  event.notification.close();
  const data = event.notification.data;
  const link =
    typeof data === "object" &&
    data !== null &&
    typeof (data as { link?: unknown }).link === "string"
      ? (data as { link: string }).link
      : undefined;
  event.waitUntil(openNotificationLink(link));
});

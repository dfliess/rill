/// <reference types="@sveltejs/kit" />
/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />

import { build, files, version } from "$service-worker";

const sw = self as unknown as ServiceWorkerGlobalScope;

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

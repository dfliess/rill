/**
 * Per-route mobile viewport control.
 *
 * Rill Cloud (`web-admin`) ships a fixed desktop viewport (`width=1024`) in
 * `app.html`, so screens that are not yet responsive stay usable on phones by
 * zooming out rather than reflowing. As individual surfaces are made
 * responsive, they opt into the true device viewport by adding their SvelteKit
 * route id to {@link MOBILE_READY_ROUTES}. `applyViewportForRoute` is invoked
 * on navigation from the root layout to switch the `<meta name="viewport">`
 * content accordingly.
 */

/**
 * SvelteKit route ids whose screens have been made responsive. Add a route
 * here only once its surface renders without horizontal overflow at 375px (see
 * `tests/mobile-smoke.spec.ts`). Matching is exact, one entry per route id:
 * listing a parent route must not opt in unmigrated children (e.g. listing the
 * project home `/[organization]/[project]` must not affect
 * `/[organization]/[project]/explore/[dashboard]`).
 *
 * `/[organization]/[project]/-/ai/[conversationId]` is intentionally not
 * listed yet: its layout is covered by the `-/ai` surface, and the smoke test
 * has no seeded conversation to visit.
 */
export const MOBILE_READY_ROUTES: string[] = [
  "/[organization]/[project]",
  "/[organization]/[project]/-/ai",
  "/[organization]/[project]/explore/[dashboard]",
  "/[organization]/[project]/canvas/[dashboard]",
];

const MOBILE_VIEWPORT = "width=device-width, initial-scale=1";
// Must match the fallback declared in `app.html` so non-migrated routes are
// restored exactly.
const DESKTOP_VIEWPORT = "width=1024, initial-scale=1.0, user-scalable=yes";

/**
 * Whether the given route id opts into the device viewport. Matches listed
 * routes exactly; children of a listed route stay on the desktop viewport
 * until they are listed themselves.
 */
export function isMobileReady(routeId: string | null): boolean {
  if (!routeId) return false;
  return MOBILE_READY_ROUTES.includes(routeId);
}

/**
 * Sets the `<meta name="viewport">` content for the given route: the device
 * viewport when the route is mobile-ready, otherwise the desktop fallback.
 * No-op during SSR (no `document`).
 */
export function applyViewportForRoute(routeId: string | null): void {
  if (typeof document === "undefined") return;
  const meta = document.querySelector('meta[name="viewport"]');
  if (!meta) return;
  meta.setAttribute(
    "content",
    isMobileReady(routeId) ? MOBILE_VIEWPORT : DESKTOP_VIEWPORT,
  );
}

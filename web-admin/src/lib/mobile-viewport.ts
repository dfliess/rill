/**
 * Mobile smoke manifest.
 *
 * web-admin serves the real device viewport on every route (see `app.html`).
 * The route ids listed here are the surfaces that were explicitly audited for
 * narrow screens and are exercised by `tests/mobile-smoke.spec.ts` at a phone
 * viewport. When a new surface is made responsive, add its route id here and
 * a matching smoke URL there; the suite's coverage guard fails otherwise.
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

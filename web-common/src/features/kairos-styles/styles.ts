// Kairos Styles (kairos-cloud#242): the app's look (surfaces, inks, radius,
// fonts, component shapes), orthogonal to Rill's light/dark mode and to the
// project theme, which only colors dashboard content.
//
// The Style is an attribute on <html>; each Style's CSS lives in this folder,
// with one block per mode. Nothing here touches the mode.

export const STYLES = [
  "kairos",
  "aurora",
  "nova",
  "editorial",
  "lumen",
  "glass",
  "apetic",
  "base",
] as const;

export type StyleId = (typeof STYLES)[number];

// Styles made for one customer (kairos-cloud#244): not listed in the pickers,
// so other customers don't see them, but usable by whoever knows the id (the
// Org settings take it as free text). Not access control: the server accepts
// any Style, and every bundle carries every Style's CSS.
export const UNLISTED_STYLES: readonly StyleId[] = ["apetic", "base"];

/** The Styles a picker lists: the listed ones, plus `active` if unlisted. */
export function listedStyles(active?: string | null): StyleId[] {
  return STYLES.filter((s) => !UNLISTED_STYLES.includes(s) || s === active);
}

export const DEFAULT_STYLE: StyleId = "kairos";

export const STYLE_ATTRIBUTE = "data-kairos-style";

export function isStyleId(value: unknown): value is StyleId {
  return typeof value === "string" && STYLES.includes(value as StyleId);
}

export type StyleSources = {
  // `style` URL parameter of an embed; only read on embed routes.
  embedParam?: string | null;
  // The viewer's own choice, from this browser.
  user?: string | null;
  // The Org's default, or its last known value while it loads.
  org?: string | null;
  isEmbed?: boolean;
};

/**
 * Resolves which Style to paint. App: user, then Org, then the default.
 * Embed: URL parameter, then Org, then the default; the viewer's own choice
 * is ignored so an embed looks the same to every viewer of the host page.
 * Unknown values fall through, so a stale or mistyped id never blanks the app.
 */
export function resolveStyle({
  embedParam,
  user,
  org,
  isEmbed = false,
}: StyleSources): StyleId {
  const candidates = isEmbed ? [embedParam, org] : [user, org];
  return candidates.find(isStyleId) ?? DEFAULT_STYLE;
}

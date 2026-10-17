// Kairos Styles (kairos-cloud#242): the app's look (surfaces, inks, radius,
// fonts, component shapes), orthogonal to Rill's light/dark mode and to the
// project theme, which only colors dashboard content.
//
// The Style is an attribute on <html>; each Style's CSS lives in this folder,
// with one block per mode. Nothing here touches the mode.

export const STYLES = [
  "kairos-v1",
  "kairos",
  "aurora",
  "nova",
  "editorial",
  "lumen",
  "glass",
] as const;

export type StyleId = (typeof STYLES)[number];

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

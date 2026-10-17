import { derived, get, writable } from "svelte/store";
import {
  STYLE_ATTRIBUTE,
  isStyleId,
  resolveStyle,
  type StyleId,
} from "./styles";

// The viewer's own choice: a StyleId, or absent to follow the Org.
const USER_KEY = "kairos:style";
// Last Org default seen on this browser, so a reload paints it before the
// Org request returns (the apps run with ssr=false).
const orgKey = (org: string) => `kairos:style:org:${org}`;

function read(key: string): string | null {
  try {
    return localStorage.getItem(key);
  } catch {
    return null;
  }
}

function write(key: string, value: string | null) {
  try {
    if (value === null) localStorage.removeItem(key);
    else localStorage.setItem(key, value);
  } catch {
    // Storage blocked (private mode, embed sandbox): the choice lasts this page.
  }
}

function isEmbedRoute(): boolean {
  try {
    return window.location.pathname.includes("/-/embed");
  } catch {
    return false;
  }
}

class StyleControl {
  private user = writable<string | null>(read(USER_KEY));
  private org = writable<{ name: string; style: string | null } | null>(null);
  private embedParam = writable<string | null>(null);

  // What the avatar menu shows as selected: a Style, or null for "Org default".
  public preference = derived(this.user, (u) => (isStyleId(u) ? u : null));

  public current = derived(
    [this.user, this.org, this.embedParam],
    ([user, org, embedParam]) =>
      resolveStyle({
        user,
        org: org?.style,
        embedParam,
        isEmbed: isEmbedRoute(),
      }),
  );

  public subscribe = this.current.subscribe;

  constructor() {
    this.current.subscribe(apply);
  }

  /** The viewer picks a Style, or null to follow the Org again. */
  public setPreference(style: StyleId | null) {
    write(USER_KEY, style);
    this.user.set(style);
  }

  /**
   * The active Org and its default. Pass `style: undefined` while the Org is
   * still loading: the cached value for that Org stands in until it arrives.
   */
  public setOrg(name: string | null, style?: string | null) {
    if (!name) {
      this.org.set(null);
      return;
    }
    if (style === undefined) {
      this.org.set({ name, style: read(orgKey(name)) });
      return;
    }
    write(orgKey(name), isStyleId(style) ? style : null);
    this.org.set({ name, style });
  }

  public setEmbedParam(style: string | null) {
    this.embedParam.set(style);
  }

  public get value(): StyleId {
    return get(this.current);
  }
}

function apply(style: StyleId) {
  if (typeof document === "undefined") return;
  document.documentElement.setAttribute(STYLE_ATTRIBUTE, style);
  // The browser chrome (PWA status bar, window title bar) follows the
  // Style's base surface, as theme-control does for the mode.
  const meta = document.querySelector('meta[name="theme-color"]');
  if (!(meta instanceof HTMLMetaElement)) return;
  const base = getComputedStyle(document.documentElement)
    .getPropertyValue("--surface-base")
    .trim();
  if (base) meta.content = base;
}

export const styleControl = new StyleControl();

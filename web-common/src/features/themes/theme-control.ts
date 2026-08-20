import { get, writable } from "svelte/store";
import { sessionStorageStore } from "@rilldata/web-common/lib/store-utils/session-storage";
import { explicitLocalStorageStore } from "@rilldata/web-common/lib/store-utils/local-storage.ts";

export type ThemeMode = "light" | "dark" | "system";

function isEmbedEnvironment(): boolean {
  if (typeof window === "undefined") return false;
  try {
    return window.location.pathname.includes("/-/embed");
  } catch {
    return false;
  }
}

const THEME_LOCAL_STORAGE_KEY = "rill:theme";
const THEME_SESSION_STORAGE_KEY = "rill:embed:theme-mode";

// The surfaces the browser paints its own chrome with: the status bar of an installed app on iOS and the
// window title bar on desktop. They mirror --surface-base for each theme in app.css. This follows the app's
// theme rather than the system's, so a machine in light mode never frames a dark app in a light bar.
const THEME_COLORS: Record<"light" | "dark", string> = {
  light: "#F4F4F1",
  dark: "#232629",
};

function applyThemeColor(scheme: "light" | "dark") {
  const meta = document.querySelector('meta[name="theme-color"]');
  // web-local has no such meta; nothing to keep in step there.
  if (meta instanceof HTMLMetaElement) meta.content = THEME_COLORS[scheme];
}

class ThemeControl {
  public current = writable<"light" | "dark">("light");
  private darkQuery = window.matchMedia("(prefers-color-scheme: dark)");
  // Kairos: oscuro por defecto en los dos, app y embed. Cuando cambiamos el
  // de la app (3f6ee1c39) nos dejamos este, así que un embed sin theme_mode
  // salía claro dentro de un producto oscuro. Quien embeba en una página clara
  // pasa theme_mode=light, igual que hasta ahora tenía que pasar dark.
  private preferenceStore = isEmbedEnvironment()
    ? sessionStorageStore<ThemeMode>(THEME_SESSION_STORAGE_KEY, "dark")
    : explicitLocalStorageStore<ThemeMode>(THEME_LOCAL_STORAGE_KEY, "dark");

  public subscribe = this.current.subscribe;
  public preference = { subscribe: this.preferenceStore.subscribe };

  constructor() {
    this.init();
  }

  init = () => {
    const currentPreference = get(this.preferenceStore);

    if (
      currentPreference === "dark" ||
      (currentPreference === "system" && this.darkQuery.matches)
    ) {
      this.setDark();
    } else {
      // Nothing to undo on a fresh document, but the browser chrome is seeded with the dark default in
      // app.html and has to come back to light for whoever chose it.
      this.removeDark();
    }

    this.darkQuery.addEventListener("change", ({ matches }) => {
      if (get(this.preferenceStore) !== "system") return;

      if (matches) {
        this.setDark();
      } else {
        this.removeDark();
      }
    });
  };

  public set: Record<ThemeMode, () => void> = {
    light: () => {
      this.preferenceStore.set("light");
      this.removeDark();
    },
    dark: () => {
      this.preferenceStore.set("dark");
      this.setDark();
    },
    system: () => {
      this.preferenceStore.set("system");

      if (this.darkQuery.matches) {
        this.setDark();
      } else {
        this.removeDark();
      }
    },
  };

  /**
   * Switches to light mode and returns a function restoring the previous mode,
   * leaving the persisted preference alone. Both the class and `current` have to
   * flip: some components read the class, others subscribe to the store.
   */
  public forceLight(): () => void {
    if (get(this.current) === "light") return () => {};
    this.removeDark();
    return () => this.setDark();
  }

  private setDark() {
    this.current.set("dark");
    document.documentElement.classList.add("dark");
    applyThemeColor("dark");
  }

  private removeDark() {
    this.current.set("light");
    document.documentElement.classList.remove("dark");
    applyThemeColor("light");
  }
}

export const themeControl = new ThemeControl();

/**
 * Returns true if the user needs to select a theme — i.e. no theme has been
 * persisted to localStorage yet. Always returns false in the embed context,
 * which manages its own ephemeral theme preference.
 */
export function isThemeSelectionNeeded(): boolean {
  if (isEmbedEnvironment()) return false;
  try {
    return !localStorage.getItem(THEME_LOCAL_STORAGE_KEY);
  } catch {
    return false;
  }
}

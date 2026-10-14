import { isManagedOrShared } from "@rilldata/web-admin/features/dashboards/listing/selectors";
import type { V1Resource } from "@rilldata/web-common/runtime-client";

export type HomeMenuItem = {
  name: string;
  title: string;
  href: string;
  // `kairos_icon` from the canvas metadata; see menu-icons.ts.
  icon?: string;
  // Personal canvases: the page where their owner edits them.
  editHref?: string;
};

/**
 * The Home menu: the project's canvases, the same ones the Dashboards tab
 * lists, by display name. A canvas's `metadata:` orders it (`kairos_order`,
 * lower first; canvases without it go after, alphabetically), keeps it out of
 * the menu (`kairos_hidden: true`; it stays in the Dashboards tab and by URL)
 * or gives it an icon (`kairos_icon`). The first item is the Home's overview
 * and lives at the project URL; the rest under `-/home/<canvas>`.
 *
 * The viewer's own personal canvases (`personalNames`) are left out: the menu
 * lists them in their own group (homePersonalMenuItems).
 */
export function homeMenuItems(
  resources: V1Resource[],
  projectBase: string,
  personalNames: ReadonlySet<string> = new Set(),
): HomeMenuItem[] {
  const canvases = resources
    .flatMap((res) => {
      const name = res.meta?.name?.name;
      const spec = res.canvas?.state?.validSpec;
      if (!name || !spec) return [];
      if (!isManagedOrShared(spec) || personalNames.has(name)) return [];
      if (res.meta?.metadata?.kairos_hidden === "true") return [];
      const order = Number.parseFloat(res.meta?.metadata?.kairos_order ?? "");
      return [
        {
          name,
          title: res.canvas?.spec?.displayName || name,
          order: Number.isFinite(order) ? order : Infinity,
          icon: res.meta?.metadata?.kairos_icon,
        },
      ];
    })
    .sort((a, b) => a.order - b.order || a.title.localeCompare(b.title));

  return canvases.map((c, i) => ({
    name: c.name,
    title: c.title,
    href:
      i === 0
        ? projectBase
        : `${projectBase}/-/home/${encodeURIComponent(c.name)}`,
    icon: c.icon,
  }));
}

/**
 * The viewer's personal canvases, alphabetically. They open inside the Home
 * like the project's canvases; editing stays in the personal canvas page.
 */
export function homePersonalMenuItems(
  personalCanvases: V1Resource[],
  projectBase: string,
): HomeMenuItem[] {
  return personalCanvases
    .flatMap((res) => {
      const name = res.meta?.name?.name;
      if (!name) return [];
      return [
        {
          name,
          title: res.canvas?.spec?.displayName || name,
          href: `${projectBase}/-/home/${encodeURIComponent(name)}`,
          editHref: `${projectBase}/-/personal/${encodeURIComponent(name)}`,
        },
      ];
    })
    .sort((a, b) => a.title.localeCompare(b.title));
}

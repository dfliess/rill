import type { V1Resource } from "@rilldata/web-common/runtime-client";

export type HomeV2MenuItem = {
  name: string;
  title: string;
  href: string;
  // `kairos_icon` from the canvas metadata; see menu-icons.ts.
  icon?: string;
};

/**
 * The Home v2 menu: every canvas, by display name. A canvas's `metadata:`
 * orders it (`kairos_order`, lower first; canvases without it go after,
 * alphabetically), keeps it out of the menu (`kairos_hidden: true`; it stays
 * in the Dashboards tab and by URL) or gives it an icon (`kairos_icon`). The first item is the page's default and
 * lives at the bare Home v2 URL.
 */
export function homeV2MenuItems(
  resources: V1Resource[],
  base: string,
): HomeV2MenuItem[] {
  const canvases = resources
    .flatMap((res) => {
      const name = res.meta?.name?.name;
      if (!name || !res.canvas?.state?.validSpec) return [];
      if (res.meta?.metadata?.kairos_hidden === "true") return [];
      const order = Number.parseFloat(res.meta?.metadata?.kairos_order ?? "");
      return [
        {
          name,
          title: res.canvas.spec?.displayName || name,
          order: Number.isFinite(order) ? order : Infinity,
          icon: res.meta?.metadata?.kairos_icon,
        },
      ];
    })
    .sort((a, b) => a.order - b.order || a.title.localeCompare(b.title));

  return canvases.map((c, i) => ({
    name: c.name,
    title: c.title,
    href: i === 0 ? base : `${base}/${encodeURIComponent(c.name)}`,
    icon: c.icon,
  }));
}

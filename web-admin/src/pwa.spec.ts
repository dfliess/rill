import { existsSync, readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { describe, expect, it } from "vitest";

// Static assets live in web-common/static (served as the web-admin assets
// dir via `kit.files.assets` in svelte.config.js).
const repoRoot = join(dirname(fileURLToPath(import.meta.url)), "..", "..");
const staticDir = join(repoRoot, "web-common", "static");
const appHtml = readFileSync(
  join(repoRoot, "web-admin", "src", "app.html"),
  "utf-8",
);
const manifestPath = join(staticDir, "manifest.webmanifest");

function pngSize(path: string): { width: number; height: number } {
  const buf = readFileSync(path);
  // IHDR width/height: big-endian uint32 at offsets 16 and 20.
  return { width: buf.readUInt32BE(16), height: buf.readUInt32BE(20) };
}

interface ManifestIcon {
  src: string;
  sizes: string;
  purpose?: string;
}

interface WebManifest {
  name: string;
  start_url: string;
  display: string;
  icons: ManifestIcon[];
}

describe("PWA manifest", () => {
  const manifest = JSON.parse(
    readFileSync(manifestPath, "utf-8"),
  ) as WebManifest;

  it("declares an installable standalone app", () => {
    expect(manifest.name).toBe("Kairos");
    expect(manifest.start_url).toBe("/");
    expect(manifest.display).toBe("standalone");
  });

  it("declares 192 and 512 icons, including maskable variants", () => {
    const sizes = manifest.icons.map((icon) => icon.sizes);
    expect(sizes).toContain("192x192");
    expect(sizes).toContain("512x512");
    const purposes = manifest.icons.map((icon) => icon.purpose);
    expect(purposes).toContain("maskable");
  });

  it("points at real PNG files with the declared dimensions", () => {
    for (const icon of manifest.icons) {
      const path = join(staticDir, icon.src);
      expect(existsSync(path), `${icon.src} missing in web-common/static`).toBe(
        true,
      );
      const [width, height] = icon.sizes.split("x").map(Number);
      expect(pngSize(path)).toEqual({ width, height });
    }
  });

  it("is linked from app.html together with theme-color", () => {
    expect(appHtml).toContain('rel="manifest"');
    expect(appHtml).toContain("manifest.webmanifest");
    expect(appHtml).toContain('name="theme-color"');
  });
});

describe("PWA service worker", () => {
  const serviceWorker = readFileSync(
    join(repoRoot, "web-admin", "src", "service-worker.ts"),
    "utf-8",
  );

  it("has its offline fallback page in static assets", () => {
    expect(serviceWorker).toContain('"/offline.html"');
    expect(existsSync(join(staticDir, "offline.html"))).toBe(true);
  });

  it("never caches navigations or non-GET requests", () => {
    // The offline fallback must be the only response served for navigations;
    // cacheFirst is reserved for static assets.
    expect(serviceWorker).toContain("networkWithOfflineFallback(request)");
    expect(serviceWorker).toContain('request.method !== "GET"');
  });
});

#!/usr/bin/env node
// Kairos Styles smoke check (kairos-cloud#242). scripts/style-guard.js proves
// the names a Style relies on still exist; this proves they still win. It
// compiles app.css and the Styles with the real Tailwind config, renders a
// page carrying the hooked classes in headless Chromium, and compares the
// computed values against what each Style promises, in light and in dark.
// A token that exists but loses the cascade (a specificity tie, a layer
// reorder, an upstream rule that now outranks ours) fails here.
//
// Usage: node scripts/style-smoke.js

import postcss from "postcss";
import postcssImport from "postcss-import";
import tailwindNesting from "tailwindcss/nesting/index.js";
import tailwindcss from "tailwindcss";
import loadConfig from "tailwindcss/loadConfig.js";
import { chromium } from "playwright";

const FIXTURE = `
<header class="kairos-app-header bg-surface-base" id="header">header</header>
<main class="bg-surface-base" id="page">
  <article class="component-card bg-surface-card rounded-sm" id="card">
    <div class="component-header-container"><h1 class="title" id="title">Title</h1></div>
    <div class="big-number">1</div>
  </article>
  <button class="primary" id="button">Button</button>
  <span class="text-accent-primary" id="accent">accent</span>
  <div class="rounded" id="rounded">rounded</div>
</main>
<div id="canvas-pdf-export-view"><article class="component-card" id="pdf-card">pdf</article></div>`;

// What each Style promises. Colors as hex, compared after both sides are
// rendered to RGB; a channel may differ by 2 (color-mix rounding).
const EXPECTED = {
  "kairos-v1": {
    light: { page: "#f4f4f1", card: "#fbfbf9", header: "#f4f4f1", accent: "#2c53b0", cardRadius: "0px", rounded: "4px", font: "Inter", title: "none", pdfEdge: "inset" },
    dark: { page: "#232629", card: "#2b2e32", header: "#232629", accent: "#4b74d0", cardRadius: "0px", rounded: "4px", font: "Inter", title: "none", pdfEdge: "inset" },
  },
  kairos: {
    light: { page: "#f4f4f1", card: "#fbfbf9", header: "#232629", accent: "#2c53b0", cardRadius: "6px", rounded: "4px", button: "4px", font: "Source Sans 3", title: "uppercase", pdfEdge: "inset" },
    dark: { page: "#232629", card: "#2b2e32", header: "#232629", accent: "#4b74d0", cardRadius: "6px", rounded: "4px", button: "4px", font: "Source Sans 3", title: "uppercase", pdfEdge: "inset" },
  },
  aurora: {
    light: { page: "#eef2f5", card: "#ffffff", header: "#ffffff", accent: "#0d9488", cardRadius: "14px", rounded: "8px", button: "9999px", font: "Inter", title: "uppercase", pdfEdge: "inset" },
    dark: { page: "#0a1015", card: "#111a21", header: "#0a1015", accent: "#2dd4bf", cardRadius: "14px", rounded: "8px", button: "9999px", font: "Inter", title: "uppercase", pdfEdge: "inset" },
  },
};

const started = Date.now();
const config = loadConfig("web-common/tailwind.config.ts");
const { css } = await postcss([
  postcssImport(),
  tailwindNesting,
  tailwindcss({ ...config, content: [{ raw: FIXTURE, extension: "html" }] }),
]).process(
  // Styles first: the worst case. They must win on specificity, not on landing after app.css.
  '@import "./web-common/src/features/kairos-styles/styles.css";\n@import "./web-common/src/app.css";\n',
  { from: "style-smoke.css" },
);

// STYLE_SMOKE_CHROMIUM points at a Chromium already on disk when Playwright's
// own build isn't installed (`npx playwright install chromium` otherwise).
const browser = await chromium.launch(
  process.env.STYLE_SMOKE_CHROMIUM ? { executablePath: process.env.STYLE_SMOKE_CHROMIUM } : {},
);
const page = await browser.newPage();
await page.setContent(`<!doctype html><html><head><style>${css}</style></head><body>${FIXTURE}</body></html>`);

const failures = [];
for (const [style, modes] of Object.entries(EXPECTED)) {
  for (const [mode, expected] of Object.entries(modes)) {
    const actual = await page.evaluate(
      ({ style, mode }) => {
        const root = document.documentElement;
        root.setAttribute("data-kairos-style", style);
        root.classList.toggle("dark", mode === "dark");
        // Any CSS color (oklab, color-mix, …) as [r, g, b], through a canvas.
        const ctx = document.createElement("canvas").getContext("2d");
        const rgb = (color) => {
          ctx.clearRect(0, 0, 1, 1);
          ctx.fillStyle = color;
          ctx.fillRect(0, 0, 1, 1);
          return [...ctx.getImageData(0, 0, 1, 1).data.slice(0, 3)];
        };
        const cs = (id) => getComputedStyle(document.getElementById(id));
        return {
          page: rgb(cs("page").backgroundColor),
          card: rgb(cs("card").backgroundColor),
          header: rgb(cs("header").backgroundColor),
          accent: rgb(cs("accent").color),
          cardRadius: cs("card").borderTopLeftRadius,
          rounded: cs("rounded").borderTopLeftRadius,
          button: cs("button").borderTopLeftRadius,
          font: getComputedStyle(document.body).fontFamily.split(",")[0].replace(/["']/g, "").trim(),
          title: cs("title").textTransform,
          // The PDF rasterizes each card alone: its edge must be drawn inside it.
          pdfEdge:
            cs("pdf-card").outlineStyle === "none" &&
            getComputedStyle(document.getElementById("pdf-card"), "::after").borderTopWidth === "1px"
              ? "inset"
              : "outside",
        };
      },
      { style, mode },
    );
    for (const [key, want] of Object.entries(expected)) {
      const got = actual[key];
      const ok = want.startsWith("#")
        ? want.match(/[0-9a-f]{2}/gi).every((c, i) => Math.abs(parseInt(c, 16) - got[i]) <= 2)
        : got === want;
      if (!ok) failures.push(`${style} ${mode} ${key}: expected ${want}, got ${Array.isArray(got) ? `rgb(${got})` : got}`);
    }
  }
}
await browser.close();

const seconds = ((Date.now() - started) / 1000).toFixed(1);
for (const f of failures) console.log(`error: ${f}`);
if (failures.length) {
  console.log(`\nstyle smoke: ${failures.length} failure(s) in ${seconds}s. See the header of scripts/style-smoke.js.`);
  process.exit(1);
}
const checks = Object.values(EXPECTED).flatMap((m) => Object.values(m)).reduce((n, e) => n + Object.keys(e).length, 0);
console.log(`style smoke: ok (${checks} computed values, ${seconds}s)`);

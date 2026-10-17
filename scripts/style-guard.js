#!/usr/bin/env node
// Kairos Styles guard (kairos-cloud#242). The Styles in
// web-common/src/features/kairos-styles/ lean on names upstream owns: its
// tokens and a few of its component classes. A weekly rebase can rename or drop
// any of them without a conflict, and the Style just stops applying. This
// fails when that happens:
//
// 1. A Style sets a token upstream no longer declares (app.css, colors.css),
//    or a --kairos-* token nothing reads.
// 2. A selector targets a class that is not registered in
//    scripts/style-guard-hooks.json, or that its registered component no
//    longer carries; or it uses a Svelte hash class (svelte-xxxx).
// 3. A component rule writes a colour literal instead of var(): inside
//    dashboards that would override the project theme. Token declarations
//    (custom properties) may hold literals, that is what tokens are.
// 4. A Style lacks its light or its dark block, or styles.ts and the CSS
//    disagree on which Styles exist.
// 5. The wiring is gone: the app.css import, a Style file not imported, or
//    the style controller missing from a root layout.
//
// It warns (exit 0) on tokens a Style sets in one mode only, and on upstream
// semantic tokens newer than the baseline, which no Style has decided on yet.
//
// Usage: node scripts/style-guard.js

import { globSync, readFileSync, readdirSync } from "node:fs";
import postcss from "postcss";
import nested from "postcss-nested";

const DIR = "web-common/src/features/kairos-styles";
const config = JSON.parse(
  readFileSync("scripts/style-guard-hooks.json", "utf8"),
);
const errors = [];
const warnings = [];

const read = (path) => readFileSync(path, "utf8");

const upstreamCss = read("web-common/src/app.css") + read("web-common/src/colors.css");
const upstreamTokens = new Set(
  [...upstreamCss.matchAll(/(--[a-z0-9-]+)\s*:/g)].map((m) => m[1]),
);

const cssFiles = readdirSync(DIR)
  .filter((f) => f.endsWith(".css"))
  .map((f) => `${DIR}/${f}`);
const cssText = cssFiles.map(read).join("\n");
const sourceText = globSync([
  "web-common/src/**/*.{ts,svelte}",
  "web-admin/src/**/*.{ts,svelte}",
  "web-local/src/**/*.{ts,svelte}",
])
  .map(read)
  .join("\n");

const COLOR_LITERAL = /#[0-9a-f]{3,8}\b|\b(rgba?|hsla?|oklch|oklab|lab|lch)\(/i;

// Top-level parts of a selector list, ignoring commas inside :is()/:not().
function splitTopLevel(selector, separator) {
  const parts = [];
  let depth = 0;
  let current = "";
  for (const ch of selector) {
    if (ch === "(") depth++;
    if (ch === ")") depth--;
    if (depth === 0 && separator.test(ch)) {
      if (current.trim()) parts.push(current.trim());
      current = "";
    } else current += ch;
  }
  if (current.trim()) parts.push(current.trim());
  return parts;
}

// A rule on :root itself (a token block), not on a descendant.
const isRootOnly = (selector) =>
  splitTopLevel(selector, /,/).every(
    (part) => part.startsWith(":root") && splitTopLevel(part, /[\s>+~]/).length === 1,
  );

const modeOf = (selector) =>
  /:not\(\.dark\)/.test(selector) ? "light" : /\.dark(?![\w-])/.test(selector) ? "dark" : "any";

// Markup hooks must sit in a class attribute or a class: directive. Hooks set
// from a prop (Button's variant) are registered `via: "style"` and must head a
// rule of the component's own <style>.
function hookPresent(cls, hook) {
  const { file, via } = typeof hook === "string" ? { file: hook, via: "markup" } : hook;
  let src;
  try {
    src = read(file);
  } catch {
    return { file, ok: false };
  }
  const word = cls.replace(/-/g, "\\-");
  const patterns =
    via === "style"
      ? [new RegExp(`^\\s*\\.${word}(?![\\w-])[^{]*\\{`, "m")]
      : [
          new RegExp(`class="[^"]*(?<![\\w-])${word}(?![\\w-])[^"]*"`),
          new RegExp(`class:${word}(?![\\w-])`),
        ];
  return { file, ok: patterns.some((re) => re.test(src)) };
}

const styles = new Map(); // id -> { light: Set, dark: Set }
const checkedHooks = new Set();

for (const file of cssFiles) {
  const root = postcss([nested]).process(read(file), { from: file }).root;
  root.walkRules((rule) => {
    const where = `${file}:${rule.source?.start?.line}`;
    const selector = rule.selector;
    if (/\bsvelte-[a-z0-9]+/i.test(selector)) {
      errors.push(`${where}: Svelte hash class in "${selector}"`);
    }

    const rootOnly = isRootOnly(selector);
    const ids = [...selector.matchAll(/data-kairos-style="([^"]+)"/g)].map((m) => m[1]);

    if (!rootOnly) {
      for (const [, cls] of selector.matchAll(/\.([a-zA-Z_][\w-]*)/g)) {
        if (cls === "dark" || checkedHooks.has(cls)) continue;
        checkedHooks.add(cls);
        const hook = config.hooks[cls];
        if (!hook) {
          errors.push(`${where}: class .${cls} is not registered in scripts/style-guard-hooks.json`);
          continue;
        }
        const { file: hookFile, ok } = hookPresent(cls, hook);
        if (!ok) errors.push(`${where}: class .${cls} no longer appears in ${hookFile}`);
      }
    }

    for (const decl of rule.nodes.filter((n) => n.type === "decl")) {
      if (decl.prop.startsWith("--")) {
        if (decl.prop.startsWith("--kairos-")) {
          const reader = new RegExp(`var\\(${decl.prop}\\s*[,)]`);
          if (!reader.test(sourceText) && !reader.test(cssText)) {
            errors.push(`${where}: ${decl.prop} is set but nothing reads it`);
          }
        } else if (!upstreamTokens.has(decl.prop)) {
          errors.push(`${where}: ${decl.prop} is not an upstream token any more (app.css, colors.css)`);
        }
        if (rootOnly) {
          for (const id of ids.length ? ids : ["kairos"]) {
            const entry = styles.get(id) ?? { light: new Set(), dark: new Set() };
            const mode = modeOf(selector);
            if (mode !== "any") entry[mode].add(decl.prop);
            styles.set(id, entry);
          }
        }
      } else if (COLOR_LITERAL.test(decl.value)) {
        errors.push(`${where}: colour literal in ${decl.prop}: ${decl.value} (use a token)`);
      }
    }
  });
}

const declared = [...read(`${DIR}/styles.ts`).matchAll(/STYLES = \[([^\]]+)\]/g)]
  .flatMap((m) => [...m[1].matchAll(/"([^"]+)"/g)].map((s) => s[1]));
for (const id of declared) {
  const entry = styles.get(id);
  if (!entry) {
    errors.push(`Style "${id}" is in styles.ts but has no CSS`);
    continue;
  }
  if (entry.light.size === 0) errors.push(`Style "${id}" has no light block`);
  if (entry.dark.size === 0) errors.push(`Style "${id}" has no dark block`);
  const lightOnly = [...entry.light].filter((t) => !entry.dark.has(t));
  const darkOnly = [...entry.dark].filter((t) => !entry.light.has(t));
  if (lightOnly.length) warnings.push(`Style "${id}" sets in light only (dark keeps Rill's): ${lightOnly.join(" ")}`);
  if (darkOnly.length) warnings.push(`Style "${id}" sets in dark only (light keeps Rill's): ${darkOnly.join(" ")}`);
}
for (const id of styles.keys()) {
  if (!declared.includes(id)) errors.push(`CSS styles "${id}" but styles.ts does not list it`);
}

const baseline = new Set(config.upstreamTokenBaseline);
for (const token of upstreamTokens) {
  if (config.upstreamTokenPrefixes.some((p) => token.startsWith(p)) && !baseline.has(token)) {
    warnings.push(`New upstream token ${token}: decide its value per Style, then add it to the baseline`);
  }
}

if (!read("web-common/src/app.css").includes('@import "./features/kairos-styles/styles.css"')) {
  errors.push("web-common/src/app.css no longer imports the Styles");
}
const index = read(`${DIR}/styles.css`);
for (const file of cssFiles) {
  const name = file.split("/").pop();
  if (name !== "styles.css" && !index.includes(`@import "./${name}"`)) {
    errors.push(`${DIR}/styles.css does not import ${name}`);
  }
}
for (const layout of ["web-admin/src/routes/+layout.svelte", "web-local/src/routes/+layout.svelte"]) {
  if (!read(layout).includes("features/kairos-styles/style-control")) {
    errors.push(`${layout} no longer loads the style controller`);
  }
}

for (const w of warnings) console.log(`warning: ${w}`);
for (const e of errors) console.log(`error: ${e}`);
if (errors.length) {
  console.log(`\nstyle guard: ${errors.length} error(s). See the header of scripts/style-guard.js.`);
  process.exit(1);
}
console.log(`style guard: ok (${declared.length} Styles, ${checkedHooks.size} hooks)`);

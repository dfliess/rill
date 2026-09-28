#!/usr/bin/env node
// Kairos brand guard: the upstream "Rill" brand must not reach Tenant-visible
// surfaces (kairos-cloud ADR-0015). Two checks:
//
// 1. Visible text in web-admin / web-common Svelte markup (text nodes and
//    human-facing attributes, which covers <title>) naming Rill.
// 2. Message catalog values (en.json, es.json) naming Rill.
//
// Unlike the i18n guard, `i18n-ignore` does NOT suppress a finding: upstream
// marks its product name that way, which is how brand leaks slip in. Markup
// inside an `{#if ... !$whiteLabel ...}` branch (up to its `{:else}`) never
// renders in Kairos and is skipped, so the fix for a leak is to gate it there.
// Operator-facing exceptions live in scripts/brand-guard-allowlist.json, each
// with its reason: message keys, and whole Svelte files (an operator surface
// stays one when upstream rewords it).
//
// Usage: node scripts/brand-guard.js

import { globSync, readFileSync } from "node:fs";

const BRAND_RE = /\bRill\b/;
const MESSAGES_DIR = "web-common/src/lib/i18n/messages";
const LOCALES = ["en", "es"];
const SVELTE_GLOBS = ["web-admin/src/**/*.svelte", "web-common/src/**/*.svelte"];
const TEXT_ATTRS = ["placeholder", "title", "aria-label", "alt", "label"];

const allowlist = JSON.parse(
  readFileSync("scripts/brand-guard-allowlist.json", "utf8"),
);
const allowedKeys = new Set(Object.keys(allowlist.catalog));
const allowedFiles = new Set(Object.keys(allowlist.markup));

// Blank a range, keeping newlines so line numbers survive.
function blank(chars, from, to) {
  for (let i = from; i < to; i++) if (chars[i] !== "\n") chars[i] = " ";
}

function stripBlocks(src) {
  return src.replace(/<(script|style)\b[^>]*>[\s\S]*?<\/\1>/gi, (m) =>
    m.replace(/[^\n]/g, " "),
  );
}

// Blank the branches that only render with the white label off: the body of
// `{#if <expr with !$whiteLabel>}` up to its matching `{:else...}` or `{/if}`.
function stripWhiteLabelOff(src) {
  const chars = [...src];
  const tagRe = /\{(#if\b[^}]*|:else\b[^}]*|\/if)\}/g;
  const stack = [];
  for (const m of src.matchAll(tagRe)) {
    const tag = m[1];
    const end = m.index + m[0].length;
    if (tag.startsWith("#if")) {
      stack.push({ gated: /!\s*\$whiteLabel\b/.test(tag), start: end });
    } else if (stack.length > 0) {
      const top = stack[stack.length - 1];
      if (top.gated && top.start !== null) {
        blank(chars, top.start, m.index);
        top.start = null;
      }
      if (tag === "/if") stack.pop();
    }
  }
  return chars.join("");
}

function lineOf(src, index) {
  return src.slice(0, index).split("\n").length;
}

function clean(raw) {
  return raw.replace(/\s+/g, " ").trim();
}

const findings = [];

for (const locale of LOCALES) {
  const file = `${MESSAGES_DIR}/${locale}.json`;
  const messages = JSON.parse(readFileSync(file, "utf8"));
  for (const [key, value] of Object.entries(messages)) {
    const text = typeof value === "string" ? value : JSON.stringify(value);
    if (!BRAND_RE.test(text) || allowedKeys.has(key)) continue;
    findings.push(`${file}  ${key}: ${text}`);
  }
}

const files = SVELTE_GLOBS.flatMap((pattern) => globSync(pattern)).sort();
for (const file of files) {
  if (allowedFiles.has(file)) continue;
  const original = readFileSync(file, "utf8");
  if (!BRAND_RE.test(original)) continue;
  // Svelte expressions are kept on purpose: a literal inside one, such as
  // {title || `${name} - Rill`}, renders just like plain text.
  const markup = stripWhiteLabelOff(stripBlocks(original));

  const hits = [];
  for (const m of markup.matchAll(/>([^<>]+)</g)) {
    if (BRAND_RE.test(m[1])) hits.push({ index: m.index + 1, text: m[1] });
  }
  const attrRe = new RegExp(`\\b(?:${TEXT_ATTRS.join("|")})="([^"]+)"`, "g");
  for (const m of markup.matchAll(attrRe)) {
    if (BRAND_RE.test(m[1])) hits.push({ index: m.index, text: m[1] });
  }

  for (const hit of hits) {
    findings.push(`${file}:${lineOf(markup, hit.index)}  text: ${clean(hit.text)}`);
  }
}

if (findings.length === 0) {
  console.log("brand-guard: no Rill brand on Tenant-visible surfaces.");
  process.exit(0);
}

console.log(`brand-guard: ${findings.length} Rill brand leak(s) (ERROR):`);
for (const f of findings) console.log(`  ${f}`);
console.log(
  "\nGate the markup with `{#if !$whiteLabel}`, override the message in " +
    "en.json + es.json, or, if the surface is operator-facing, add it to " +
    "scripts/brand-guard-allowlist.json with its reason (ADR-0015).",
);
process.exit(1);

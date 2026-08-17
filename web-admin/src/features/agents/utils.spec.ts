import { describe, expect, it } from "vitest";
import {
  approvalPolicyLabel,
  escapeInline,
  escapeInvisible,
  escapePlain,
  escapeReadable,
  escapeToAscii,
  eventDecidedBy,
  parseCanonicalArgs,
  subjectLabel,
  timestampSortKey,
} from "./utils";

const NAMES = new Map([["f960bb37-9d8a", "Diego Kairos"]]);

/**
 * Entries also carry the pre-escape source text, which exact mode renders from.
 * These assertions are about what the table shows, so they compare that.
 */
function shown(e: { key: string; value: string; readable?: string }) {
  return e.readable === undefined
    ? { key: e.key, value: e.value }
    : { key: e.key, value: e.value, readable: e.readable };
}

describe("subjectLabel", () => {
  it("names a subject it can resolve", () => {
    expect(subjectLabel("f960bb37-9d8a", NAMES)).toBe("Diego Kairos");
  });

  it("falls back to the raw subject", () => {
    // A decider who left the project, or a caller who cannot list members: the
    // audit record still has to render.
    expect(subjectLabel("someone-else", NAMES)).toBe("someone-else");
    expect(subjectLabel("f960bb37-9d8a", undefined)).toBe("f960bb37-9d8a");
  });

  it("renders an empty subject as a dash", () => {
    expect(subjectLabel("", NAMES)).toBe("—");
    expect(subjectLabel(undefined, NAMES)).toBe("—");
  });
});

describe("eventDecidedBy", () => {
  it("reads the decider off a decision event's payload", () => {
    expect(eventDecidedBy({ decided_by: "admin:bob" })).toBe("admin:bob");
  });

  it("returns empty for events that carry no decision", () => {
    expect(eventDecidedBy(undefined)).toBe("");
    expect(eventDecidedBy(null)).toBe("");
    expect(eventDecidedBy({})).toBe("");
    expect(eventDecidedBy({ decided_by: 42 })).toBe("");
    expect(eventDecidedBy("run.resumed")).toBe("");
  });
});

describe("timestampSortKey", () => {
  it("orders RFC3339 strings chronologically", () => {
    const keys = ["2026-07-24T22:01:42Z", "2026-07-24T21:59:23Z"].map(
      timestampSortKey,
    );
    expect([...keys].sort()).toEqual([
      "2026-07-24T21:59:23Z",
      "2026-07-24T22:01:42Z",
    ]);
  });

  it("sorts non-string timestamps first rather than throwing", () => {
    expect(timestampSortKey(undefined)).toBe("");
    expect(timestampSortKey({ seconds: 1 })).toBe("");
  });
});

describe("parseCanonicalArgs", () => {
  it("degrades for approvals recorded without a preimage", () => {
    // Rows from before the arguments were persisted: nothing to show, and the
    // card must fall back to the one-line proposal without breaking.
    for (const input of [undefined, ""]) {
      const parsed = parseCanonicalArgs(input);
      expect(parsed.isEmpty).toBe(false);
      expect(parsed.structured).toBe(false);
      expect(parsed.entries.map(shown)).toEqual([]);
      expect(parsed.raw).toBe("");
    }
  });

  it("treats the canonical empty object as an action without arguments", () => {
    const parsed = parseCanonicalArgs("{}");
    expect(parsed.isEmpty).toBe(true);
    expect(parsed.structured).toBe(true);
    expect(parsed.entries.map(shown)).toEqual([]);
    // The raw preimage is still carried: `{}` is what the hash covers.
    expect(parsed.raw).toBe("{}");
  });

  it("decomposes a canonical JSON object into entries in signed order", () => {
    const raw = '{"alpha":{"n":[1,"dos"]},"summary":"coste alto: París"}';
    const parsed = parseCanonicalArgs(raw);
    expect(parsed.structured).toBe(true);
    expect(parsed.isEmpty).toBe(false);
    // Never re-serialized: what the card offers as "the signed bytes" must be
    // the input, untouched.
    expect(parsed.raw).toBe(raw);
    expect(parsed.entries.map((e) => e.key)).toEqual(["alpha", "summary"]);
    // Every value is the exact slice of the signed text, strings included:
    // their quotes stay so a string never reads as a scalar.
    expect(parsed.entries[1]).toMatchObject({
      key: "summary",
      value: '"coste alto: París"',
    });
    expect(parsed.entries[0]).toMatchObject({
      key: "alpha",
      value: '{"n":[1,"dos"]}',
    });
  });

  it("shows scalar values as their exact signed text", () => {
    const parsed = parseCanonicalArgs(
      '{"count":3,"dry_run":false,"note":null}',
    );
    expect(parsed.entries.map(shown)).toEqual([
      { key: "count", value: "3" },
      { key: "dry_run", value: "false" },
      { key: "note", value: "null" },
    ]);
  });

  it("preserves integers beyond JavaScript's double precision digit for digit", () => {
    // JSON.parse would round both: 2^53+1 to ...992 and the 20-digit id to
    // ...567000. The signer must read exactly the digits that were hashed;
    // anything else signs a different account than the one on screen.
    const raw = '{"account_id":9007199254740993,"big":12345678901234567890}';
    const parsed = parseCanonicalArgs(raw);
    expect(parsed.structured).toBe(true);
    expect(parsed.entries.map(shown)).toEqual([
      { key: "account_id", value: "9007199254740993" },
      { key: "big", value: "12345678901234567890" },
    ]);
  });

  it("keeps integer-like keys in the signed order, not JavaScript's index order", () => {
    // Object.entries would surface "9" before "10"; the signed text says
    // otherwise and the table must read in the signed order.
    const raw = '{"10":1,"9":2,"z":3}';
    const parsed = parseCanonicalArgs(raw);
    expect(parsed.entries.map((e) => e.key)).toEqual(["10", "9", "z"]);
    expect(parsed.entries.map((e) => e.value)).toEqual(["1", "2", "3"]);
  });

  it("decodes escapes that hide nothing instead of showing the text twice", () => {
    const parsed = parseCanonicalArgs('{"name":"Jos\\u00e9","city":"París"}');
    // `José` has nothing to hide — no invisible character, no deceiving
    // composition — so its escape was notation and the decoded form IS the
    // value. No second line: the reader would only be comparing one text with
    // itself. The already-literal one reads the same way, as it always did.
    expect(parsed.entries.map(shown)).toEqual([
      { key: "name", value: '"José"' },
      { key: "city", value: '"París"' },
    ]);
  });

  it("decodes a multi-line body rather than stacking \\n above its paragraphs", () => {
    // The case that made this rule necessary: an alert body arrives as a single
    // JSON string, and rendering the signed escapes on top of their own decoding
    // buried the decision under two copies of the same paragraphs.
    const parsed = parseCanonicalArgs('{"details":"Uno.\\n\\nDos."}');
    expect(parsed.entries.map(shown)).toEqual([
      { key: "details", value: '"Uno.\n\nDos."' },
    ]);
  });

  it("is not confused by delimiters inside string values", () => {
    const raw = '{"body":"a,b}c{d","quote":"she said \\"hi\\", bye","next":1}';
    const parsed = parseCanonicalArgs(raw);
    expect(parsed.structured).toBe(true);
    expect(parsed.entries.map(shown)).toEqual([
      { key: "body", value: '"a,b}c{d"' },
      { key: "quote", value: '"she said "hi", bye"' },
      { key: "next", value: "1" },
    ]);
  });

  it("refuses to build a table it cannot fully account for", () => {
    // A truncated object, trailing garbage, or a broken string must fall back
    // to the raw view rather than render a table that omits signed material.
    for (const raw of ['{"a":1', '{"a":1} x', '{"a":"unterminated', '{"a"1}']) {
      const parsed = parseCanonicalArgs(raw);
      expect(parsed.structured).toBe(false);
      expect(parsed.entries.map(shown)).toEqual([]);
      expect(parsed.raw).toBe(raw);
    }
  });

  it("refuses text that only looks like JSON to a slicer", () => {
    // These are the cases a hand-written tokenizer waves through: an unquoted
    // literal, a leading-zero number, mismatched delimiters, garbage inside a
    // container, an invalid escape. None is JSON, so none may be presented as
    // a decomposition of signed arguments — and the simulated path signs
    // arbitrary model text, so this is reachable without any corruption.
    for (const raw of [
      '{"a":truth}',
      '{"a":01}',
      '{"a":[1}}',
      '{"a":{garbage}}',
      '{"a":{"x":"\\q"}}',
    ]) {
      const parsed = parseCanonicalArgs(raw);
      expect(parsed.structured, `expected raw-only for ${raw}`).toBe(false);
      expect(parsed.entries.map(shown)).toEqual([]);
      expect(parsed.raw).toBe(raw);
    }
  });

  it("makes direction-changing and invisible characters visible", () => {
    // Go leaves U+202E literal in the canonical JSON (it is printable Unicode,
    // not a C0 control), and the browser then reverses what follows: the signed
    // filename ends in "fdp.exe" but the eye reads "informeexe.pdf". Whoever
    // writes the argument must not get to choose what the approver sees, so the
    // renderer shows the character instead of obeying it.
    const rlo = "\u202E";
    const parsed = parseCanonicalArgs(
      `{"file":"informe${rlo}fdp.exe","zw":"a\u200Bb"}`,
    );
    expect(parsed.entries[0].value).toBe('"informe\\u202Efdp.exe"');
    expect(parsed.entries[0].value).not.toContain(rlo);
    expect(parsed.entries[1].value).toBe('"a\\u200Bb"');
    // The raw preimage is still the raw preimage: only its rendering is escaped.
    expect(parsed.raw).toContain(rlo);
  });

  it("escapes an argument key, which the proposer controls too", () => {
    // A map key is as attacker-controlled as a value: Go marshals it literally,
    // so a bidi override in the key reorders the row's label just the same.
    const parsed = parseCanonicalArgs(`{"recipient${"\u202E"}txt":"a"}`);
    expect(parsed.entries[0].key).toBe("recipient\\u202Etxt");
    expect(parsed.entries[0].key).not.toContain("\u202E");
  });

  it("covers the invisibles beyond bidi and zero-width", () => {
    // The soft hyphen is the sharpest of these: it renders as nothing at all,
    // so the stored "ops", U+00AD, "dev@example.com" reads as an address it is not.
    expect(escapeInvisible("ops\u00ADdev@example.com")).toBe(
      "ops\\u00ADdev@example.com",
    );
    expect(escapeInvisible("a\u034Fb")).toBe("a\\u034Fb");
    expect(escapeInvisible("a\u180Eb")).toBe("a\\u180Eb");
    // Tag characters are astral and can smuggle a whole hidden ASCII string;
    // they need the u flag to match as single code points, not surrogates, and
    // the braced form so the escape cannot be misread (see the astral test).
    expect(escapeInvisible("a\u{E0041}b")).toBe("a\\u{E0041}b");
  });

  it("escapes the joiners in the signed text and keeps them in the reading aid", () => {
    // ZWJ, ZWNJ and the variation selector are NOT harmless: "admin" and
    // "ad<ZWNJ>min" are different signed values that look identical in Latin
    // text. So the exact rendering escapes them like anything else, and the aid
    // beside it keeps them so emoji and Persian spelling still read normally.
    expect(escapeInvisible("ad\u200Cmin")).toBe("ad\\u200Cmin");
    expect(escapeReadable("ad\u200Cmin")).toBe("ad\u200Cmin");
    expect(escapeInvisible("\u{1F469}\u200D\u{1F4BB}")).toBe(
      "\u{1F469}\\u200D\u{1F4BB}",
    );
    expect(escapeReadable("\u{1F469}\u200D\u{1F4BB}")).toBe(
      "\u{1F469}\u200D\u{1F4BB}",
    );
    // The aid is not a hole: it still escapes everything that hides or reorders.
    expect(escapeReadable("informe\u202Efdp.exe")).toBe(
      "informe\\u202Efdp.exe",
    );
    expect(escapeReadable("ops\u00ADdev@example.com")).toBe(
      "ops\\u00ADdev@example.com",
    );
  });

  it("never renders two different signed texts the same way", () => {
    // The escape is only worth anything if it is injective. Without doubling the
    // backslash, the literal ASCII text `ad\\u200Cmin` and the same word with a
    // real ZWNJ both display as `ad\\u200Cmin`: two different preimages, two
    // different hashes, one appearance — the exact failure the escape exists to
    // prevent, reintroduced by the escape itself.
    const literal = "ad\\u200Cmin";
    const real = "ad\u200Cmin";
    expect(escapeInline(literal)).not.toBe(escapeInline(real));
    expect(escapeInline(literal)).toBe("ad\\\\u200Cmin");
    expect(escapeInline(real)).toBe("ad\\u200Cmin");
    expect(escapeReadable(literal)).not.toBe(escapeReadable(real));
    expect(escapePlain(literal)).not.toBe(escapePlain(real));

    // JSON text is exempt because the encoder already doubled the backslash, so
    // the two are distinguishable without help — and doubling again would
    // misrepresent the signed bytes.
    expect(escapeInvisible('"ad\\\\u200Cmin"')).toBe('"ad\\\\u200Cmin"');
    expect(escapeInvisible('"ad\u200Cmin"')).toBe('"ad\\u200Cmin"');
  });

  it("keeps a line feed but never lets two controls share an appearance", () => {
    // In a pre-wrap block a CR and a LF draw the same break, and a tab is
    // whitespace like any other: three different signed texts, one appearance.
    // The LF survives because it is the one control that renders as itself and
    // the multi-line body is what the block exists to show.
    expect(escapePlain("a\nb")).toBe("a\nb");
    expect(escapePlain("a\rb")).toBe("a\\u000Db");
    expect(escapePlain("a\tb")).toBe("a\\u0009b");
    expect(escapePlain("a\rb")).not.toBe(escapePlain("a\nb"));
    expect(escapePlain("a\u2028b")).toBe("a\\u2028b");
  });

  it("separates a precomposed accent from its decomposed twin", () => {
    // "Jos" + U+00E9 and "Jose" + U+0301 are different signed values that the
    // browser draws as the same "José". The precomposed one — what a keyboard
    // actually produces — is left alone so ordinary Spanish reads normally; the
    // decomposed one, which nobody types by hand, shows its combining mark.
    const precomposed = "Jos\u00E9";
    const decomposed = "Jose\u0301";
    expect(escapeInvisible(precomposed)).toBe("José");
    expect(escapeInvisible(decomposed)).toBe("Jose\\u0301");
    expect(escapeInvisible(precomposed)).not.toBe(escapeInvisible(decomposed));

    // The reading aid keeps both natural, which is what makes scripts that need
    // combining marks in normal spelling still readable beside the exact form.
    expect(escapeReadable(decomposed)).toBe(decomposed);
  });

  it("separates composed twins that draw the same glyphs", () => {
    // No character class catches these: every code point involved is ordinary on
    // its own, and only their composition makes two different signed texts draw
    // identically. NFC stability is what tells them apart.
    const hangulPrecomposed = "\uAC00";
    const hangulJamo = "\u1100\u1161";
    expect(escapeInvisible(hangulPrecomposed)).toBe(hangulPrecomposed);
    expect(escapeInvisible(hangulJamo)).toBe("\u1100\\u1161");
    expect(escapeInvisible(hangulPrecomposed)).not.toBe(
      escapeInvisible(hangulJamo),
    );

    // The Kelvin sign folds to K under NFC and is drawn as one.
    expect(escapeInvisible("K")).toBe("K");
    expect(escapeInvisible("\u212A")).toBe("\\u212A");
  });

  it("leaves scripts whose normal spelling is NFC-stable readable", () => {
    // Escaping combining marks wholesale would have mangled these; the NFC test
    // does not, because their ordinary spelling is already what NFC produces.
    for (const s of ["\u05D0\u05B7", "\u0645\u064E", "revisión de París"]) {
      expect(escapeInvisible(s)).toBe(s);
    }
  });

  it("separates a composition that reaches across another mark", () => {
    // Adjacency is not enough: in "a" + U+0327 + U+0301 every adjacent pair is
    // stable on its own, yet the whole composes to what U+00E1 + U+0327 does.
    // Once the text is known to be unstable, every mark in it becomes visible.
    const acrossMark = "a\u0327\u0301";
    const precomposedFirst = "\u00E1\u0327";
    expect(acrossMark.normalize("NFC")).toBe(precomposedFirst.normalize("NFC"));
    expect(escapeInvisible(acrossMark)).toBe("a\\u0327\\u0301");
    expect(escapeInvisible(acrossMark)).not.toBe(
      escapeInvisible(precomposedFirst),
    );
  });

  it("escapes ASCII spaces at the edges of a one-line identifier", () => {
    // A trailing space leaves no ink: "admin" and "admin " would be one label on
    // screen and two different keys to the target system.
    expect(escapeInline("admin")).toBe("admin");
    expect(escapeInline("admin ")).toBe("admin\\u0020");
    expect(escapeInline(" admin")).toBe("\\u0020admin");
    expect(escapeInline("a b")).toBe("a b");
    expect(escapeInline("")).toBe("");
  });

  it("escapes every non-ASCII code point in exact mode", () => {
    // The escape hatch for what no rule can finish: Cyrillic "а" and Latin "a"
    // need no composition, no invisibility and no trick, and only a rendering
    // confined to printable ASCII can tell them apart.
    expect(escapeToAscii("\u0430")).toBe("\\u0430");
    expect(escapeToAscii("a")).toBe("a");
    expect(escapeToAscii("\u0430")).not.toBe(escapeToAscii("a"));
    expect(escapeToAscii("revisión")).toBe("revisi\\u00F3n");
    expect(escapeToAscii("a\\b")).toBe("a\\\\b");
    expect(escapeToAscii('{"x":1}')).toBe('{"x":1}');
  });

  it("separates a non-breaking space from a plain one", () => {
    expect(escapeInvisible("a b")).toBe("a b");
    expect(escapeInvisible("a\u00A0b")).toBe("a\\u00A0b");
    expect(escapeInvisible("a\u2009b")).toBe("a\\u2009b");
    expect(escapeInvisible("a b")).not.toBe(escapeInvisible("a\u00A0b"));
  });

  it("writes an astral escape unambiguously", () => {
    // A bare \uE0041 could be read as U+E004 followed by "1".
    expect(escapePlain("a\u{E0041}b")).toBe("a\\u{E0041}b");
    expect(escapeInline("a\u{E0041}b")).toBe("a\\u{E0041}b");
  });

  it("pairs an escaped joiner with its natural reading", () => {
    // The signed line shows what is really there; the aid shows how it looks.
    // Seeing both is what lets a reader notice the difference at all. This is
    // the case the decoding rule must NOT collapse: unlike a `\n`, the escape
    // here is the only reason the reader knows the character exists.
    const parsed = parseCanonicalArgs('{"user":"ad\u200Cmin"}');
    expect(parsed.entries[0].value).toBe('"ad\\u200Cmin"');
    expect(parsed.entries[0].readable).toBe("ad\u200Cmin");
  });

  it("keeps the escape when a value hides a character behind an escape it also needs", () => {
    // Both reasons at once: a JSON `\n` (pure notation) and a bidi override
    // (never notation). Decoding the first must not smuggle the second past the
    // reader, so the pair stays.
    const parsed = parseCanonicalArgs('{"file":"a\\nb\u202Ecod.exe"}');
    expect(parsed.entries[0].value).toBe('"a\\nb\\u202Ecod.exe"');
    expect(parsed.entries[0].readable).toBe("a\nb\\u202Ecod.exe");
  });

  it("escapes a control character in a one-line field but not in a value", () => {
    // A real newline in a key does not show as a newline: <dt> collapses it to a
    // space, so "recipient\nadmin" reads as two words with no hint which is the
    // real one. In a multi-line value the newline is the point, so it stays.
    expect(escapeInline("recipient\nadmin")).toBe("recipient\\u000Aadmin");
    expect(escapeInline("a\u2028b")).toBe("a\\u2028b");
    expect(escapeInvisible("linea uno\nlinea dos")).toBe(
      "linea uno\nlinea dos",
    );
  });

  it("escapes invisibles without touching ordinary text", () => {
    expect(escapeInvisible("coste alto: revisión de París")).toBe(
      "coste alto: revisión de París",
    );
    expect(escapeInvisible("")).toBe("");
    expect(escapeInvisible("a\uFEFFb\u2066c")).toBe("a\\uFEFFb\\u2066c");
  });

  it("distinguishes a string from the scalar it imitates", () => {
    // The reason values keep their quotes: a tool with a boolean dry_run does
    // not behave the same when handed the string "false", and an approver who
    // cannot tell them apart is back to signing something they did not read.
    const asString = parseCanonicalArgs('{"enabled":"false"}');
    const asBoolean = parseCanonicalArgs('{"enabled":false}');
    expect(asString.entries[0].value).toBe('"false"');
    expect(asBoolean.entries[0].value).toBe("false");
    expect(asString.entries[0].value).not.toBe(asBoolean.entries[0].value);

    const numeric = parseCanonicalArgs('{"id":"123","other":123}');
    expect(numeric.entries.map((e) => e.value)).toEqual(['"123"', "123"]);
  });

  it("keeps a non-JSON preimage raw-only (the legacy simulated path signs the proposal text)", () => {
    const raw = "crear ticket P2 en Jira";
    const parsed = parseCanonicalArgs(raw);
    expect(parsed.structured).toBe(false);
    expect(parsed.isEmpty).toBe(false);
    expect(parsed.entries.map(shown)).toEqual([]);
    expect(parsed.raw).toBe(raw);
  });

  it("keeps a non-object JSON preimage raw-only", () => {
    for (const raw of ['["a","b"]', '"solo"', "42"]) {
      const parsed = parseCanonicalArgs(raw);
      expect(parsed.structured).toBe(false);
      expect(parsed.entries.map(shown)).toEqual([]);
      expect(parsed.raw).toBe(raw);
    }
  });

  it("survives a huge value without altering the raw bytes", () => {
    const big = "x".repeat(100_000);
    const raw = JSON.stringify({ body: big, title: "t" });
    const parsed = parseCanonicalArgs(raw);
    expect(parsed.structured).toBe(true);
    expect(parsed.raw).toBe(raw);
    expect(parsed.entries.map(shown)).toEqual([
      { key: "body", value: `"${big}"` },
      { key: "title", value: '"t"' },
    ]);
  });
});

describe("approvalPolicyLabel", () => {
  it("localizes the approval_required policy", () => {
    const label = approvalPolicyLabel("approval_required");
    expect(label).toBeTruthy();
    expect(label).not.toBe("approval_required");
  });

  it("falls back to the raw token for unknown policies", () => {
    expect(approvalPolicyLabel("deny")).toBe("deny");
    expect(approvalPolicyLabel(undefined)).toBe("—");
    expect(approvalPolicyLabel("   ")).toBe("—");
  });
});

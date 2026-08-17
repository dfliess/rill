import type { Color } from "@rilldata/web-common/components/tag/types";
import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
import type { AgentRunData } from "./types";

// The run/approval status fields are free-form strings on the proto (the backend
// owns the vocabulary), so we normalize and map the ones we know and fall back to
// a neutral tag for anything unrecognized.
function normalizeStatus(status: string | undefined): string {
  return (status ?? "").trim().toLowerCase();
}

const RUN_STATUS_COLORS: Record<string, Color> = {
  queued: "gray",
  pending: "gray",
  running: "blue",
  in_progress: "blue",
  waiting_approval: "amber",
  pending_approval: "amber",
  suspended: "amber",
  succeeded: "green",
  completed: "green",
  failed: "red",
  errored: "red",
  error: "red",
  cancelled: "gray",
  canceled: "gray",
};

const APPROVAL_STATUS_COLORS: Record<string, Color> = {
  pending: "amber",
  approved: "green",
  denied: "red",
  cancelled: "gray",
  canceled: "gray",
};

export function agentRunStatusColor(status: string | undefined): Color {
  return RUN_STATUS_COLORS[normalizeStatus(status)] ?? "gray";
}

export function agentApprovalStatusColor(status: string | undefined): Color {
  return APPROVAL_STATUS_COLORS[normalizeStatus(status)] ?? "gray";
}

// The status/trigger vocabularies are owned by the backend (free-form strings),
// so we localize the tokens we know and fall back to the raw value for anything
// unrecognized. This keeps the filter dropdowns and the status chips consistent.
export function agentRunStatusLabel(status: string | undefined): string {
  switch (normalizeStatus(status)) {
    case "queued":
    case "pending":
      return m.agents_run_status_queued();
    case "running":
    case "in_progress":
      return m.agents_run_status_running();
    case "waiting_approval":
    case "pending_approval":
    case "suspended":
      return m.agents_run_status_waiting_approval();
    case "succeeded":
    case "completed":
      return m.agents_run_status_succeeded();
    case "failed":
    case "errored":
    case "error":
      return m.agents_run_status_failed();
    case "cancelled":
    case "canceled":
      return m.agents_run_status_cancelled();
    default:
      return status?.trim() || m.agents_status_unknown();
  }
}

export function agentApprovalStatusLabel(status: string | undefined): string {
  switch (normalizeStatus(status)) {
    case "pending":
      return m.agents_approval_status_pending();
    case "approved":
      return m.agents_approval_status_approved();
    case "denied":
      return m.agents_approval_status_denied();
    case "cancelled":
    case "canceled":
      return m.agents_approval_status_cancelled();
    default:
      return status?.trim() || m.agents_status_unknown();
  }
}

export function agentTriggerLabel(trigger: string | undefined): string {
  switch (normalizeStatus(trigger)) {
    case "manual":
      return m.agents_trigger_manual();
    case "alert":
      return m.agents_trigger_alert();
    default:
      return trigger?.trim() || "—";
  }
}

export function isApprovalPending(status: string | undefined): boolean {
  return normalizeStatus(status) === "pending";
}

/**
 * Timestamps arrive as RFC3339 strings (the proto3 JSON projection of
 * google.protobuf.Timestamp), even though the static type models the message
 * shape, so we accept `unknown` and coerce defensively.
 */
export function formatDateTime(value: unknown): string {
  if (typeof value !== "string" || value === "") return "—";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "—";
  return date.toLocaleString();
}

/**
 * Sort key for a proto timestamp. Like `formatDateTime`, this accepts `unknown`:
 * the value arrives as an RFC3339 string even though the static type models the
 * message shape, and RFC3339 sorts correctly as text. Anything else sorts first.
 */
export function timestampSortKey(value: unknown): string {
  return typeof value === "string" ? value : "";
}

/**
 * Characters that change what the eye reads without changing the text.
 *
 * Two Unicode properties instead of a hand-kept list: `Default_Ignorable_Code_Point`
 * is every code point that renders as nothing (soft hyphen, zero-width space and
 * joiners, variation selectors, the BOM, the astral tag characters that can smuggle
 * a whole hidden ASCII string), and `Bidi_Control` is every one that reorders what
 * follows it (the overrides, embeddings, isolates and directional marks). Between
 * them they are the definition of the problem, and the engine keeps the tables
 * current — a literal range list would silently rot with each Unicode release.
 *
 * Go's json.Marshal leaves all of these literal (they are printable Unicode, not
 * C0 controls), so they travel intact into the signed preimage and the browser
 * then obeys them.
 *
 * Also included: the Unicode spaces (`\p{Zs}`), because a NBSP or a thin space is
 * drawn as a plain one. The ASCII space is exempt — it is the one the others are
 * confused WITH, so it never needs distinguishing from itself.
 *
 * Characters that merely LOOK alike through composition are handled separately,
 * by escapeUnstableUnderNFC, because a code point class cannot express them.
 */
const INVISIBLE_OR_BIDI =
  /[\p{Default_Ignorable_Code_Point}\p{Bidi_Control}\p{Zs}\p{Cc}]/gu;

/**
 * The same, plus the C0/C1 controls and the line/paragraph separators, for text
 * rendered on ONE line. There a real newline does not show as a newline: the
 * element collapses it to a space, so a key signed as `recipient\nadmin` reads
 * as two words and hides which one it really is.
 */
const INVISIBLE_OR_BIDI_OR_CONTROL =
  /[\p{Default_Ignorable_Code_Point}\p{Bidi_Control}\p{Zs}\p{Cc}\u2028\u2029]/gu;

/**
 * Joiners kept as they are in the READING AID only: the emoji ZWJ sequences, the
 * variation selector that turns a character into an emoji, and the ZWNJ that
 * Persian and other scripts use in normal spelling.
 *
 * They are not harmless — `admin` and `ad<ZWNJ>min` are different values that
 * look identical in Latin text — which is exactly why they stay escaped in the
 * signed text and in every one-line identifier. The split is the point: the
 * signed rendering is exact and shows them, the aid beside it is natural and
 * does not, and the reader can see both.
 */
const KEPT_IN_READING_AID_CHARS = new Set(["\u200C", "\u200D", "\uFE0F"]);
const UNICODE_SPACE = /\p{Zs}/u;
const COMBINING_MARK = /[\p{Mn}\p{Me}\p{Mc}]/u;

/** True for the code points the reading aid renders naturally rather than escaped. */
function keptInReadingAid(c: string): boolean {
  return (
    KEPT_IN_READING_AID_CHARS.has(c) || UNICODE_SPACE.test(c) || isLineFeed(c)
  );
}

/**
 * The line feed is the one control that draws as itself, so it is the one that
 * stays. Fields rendered on a single line are the exception: there it collapses
 * to a space and has to be escaped like any other control.
 */
function isLineFeed(c: string): boolean {
  return c === "\n";
}

/**
 * Escapes the code points that make two different texts draw the same glyphs
 * through Unicode composition — which no character class can capture, because
 * every character involved is ordinary on its own.
 *
 * The test is NFC stability. `Jose` + U+0301 composes to `José`, the Hangul jamo
 * U+1100 U+1161 compose to `가`, and the Kelvin sign U+212A folds to `K`: in each
 * case the text is not what NFC would produce, and the browser draws it exactly
 * like its composed twin. A precomposed `é`, a precomposed `가`, an ASCII `K` —
 * what a keyboard actually types — are NFC-stable and stay untouched.
 *
 * It is also gentler than escaping combining marks wholesale: Hebrew niqqud and
 * Arabic harakat are NFC-stable, so ordinary spelling in those scripts stays
 * readable while the decomposed impostor does not.
 *
 * A character is escaped when it is unstable on its own, or when it composes
 * with the character before it; the leading character of a composing pair is
 * left alone so the escape marks what was added, not the word it attached to.
 */
function escapeUnstableUnderNFC(text: string): string {
  // One check up front: NFC-stable text — nearly all of it — needs no work.
  if (text.normalize("NFC") === text) return text;
  // Past that point the text is a decomposed form of something, and pairwise
  // adjacency is not enough to find where: in `a` + U+0327 + U+0301 every
  // adjacent pair is stable on its own, yet the whole composes to the same
  // thing as U+00E1 + U+0327. So once instability is established, every mark in
  // the text is escaped along with every individually unstable code point. The
  // base letters stay, so the word remains readable and the marks are visible.
  const chars = [...text];
  return chars
    .map((c, i) => {
      const prev = i > 0 ? chars[i - 1] : "";
      // Three ways a code point can be part of the disguise: it is a mark (the
      // across-a-mark case, where no adjacent pair is unstable), it folds on its
      // own (the Kelvin sign), or it composes with what precedes it (the Hangul
      // jamo, which are ordinary letters and no mark class would catch).
      const unstableAlone = c.normalize("NFC") !== c;
      const composesWithPrev =
        prev !== "" && (prev + c).normalize("NFC") !== prev + c;
      return COMBINING_MARK.test(c) || unstableAlone || composesWithPrev
        ? escapeCodePoint(c)
        : c;
    })
    .join("");
}

/**
 * One code point as a visible escape. Astral code points use the `\u{...}` form:
 * a bare `\uE0041` could be read as U+E004 followed by "1", and an encoding whose
 * output is ambiguous is exactly the problem this is meant to solve.
 */
function escapeCodePoint(c: string): string {
  const cp = c.codePointAt(0) ?? 0;
  const hex = cp.toString(16).toUpperCase();
  return cp > 0xffff ? `\\u{${hex}}` : "\\u" + hex.padStart(4, "0");
}

function escapeWith(
  text: string,
  pattern: RegExp,
  keep?: (c: string) => boolean,
): string {
  return text.replace(pattern, (c) =>
    // The ASCII space is the one every other space is confused WITH, so it is
    // the one that never needs distinguishing from itself. Escaping it would
    // turn every sentence into \u0020-separated noise for no gain.
    c === " " || keep?.(c) ? c : escapeCodePoint(c),
  );
}

/**
 * The same, on PLAIN text, where a backslash is just a character.
 *
 * Escaping the invisibles is only useful if the result is unambiguous, and it is
 * not unless the escape marker itself is escaped: otherwise the literal ASCII
 * text `ad\u200Cmin` and the same word carrying a real ZWNJ both render as
 * `ad\u200Cmin` — two different signed values, one appearance, which is the very
 * failure being fixed. Doubling the backslash separates them.
 *
 * JSON text does not need this and must not get it: there a real backslash is
 * already doubled by the encoder, so the escapes are distinguishable as they are
 * and doubling again would just misrepresent the signed bytes.
 */
function escapePlainWith(
  text: string,
  pattern: RegExp,
  keep?: (c: string) => boolean,
  composed = true,
): string {
  const doubled = text.replace(/\\/g, "\\\\");
  return escapeWith(
    composed ? escapeUnstableUnderNFC(doubled) : doubled,
    pattern,
    keep,
  );
}

/**
 * Renders invisible and direction-changing characters as visible escapes, for
 * text that is JSON.
 *
 * This is the one place where display deliberately departs from the signed bytes
 * character for character, and it departs toward MORE fidelity, not less: a raw
 * U+202E is not "the signed text shown as it is", it is the signed text shown
 * reversed, chosen by whoever wrote the argument: the stored characters
 * `informe`, U+202E, `fdp.exe` reach the eye as `informeexe.pdf`. An approver
 * reading a filename, a recipient or an amount must see the characters that are
 * actually there.
 *
 * Nothing is exempt here. Every ignorable code point can make two different
 * signed values look like one.
 */
export function escapeInvisible(text: string): string {
  return escapeWith(
    escapeUnstableUnderNFC(text),
    INVISIBLE_OR_BIDI,
    isLineFeed,
  );
}

/**
 * Every control except the line feed. In a `white-space: pre-wrap` block a LF is
 * the one control that renders as itself, so it can stay; CR renders as the SAME
 * line break, and a tab as indistinguishable whitespace, which would make `a\rb`
 * and `a\nb` — different signed texts — look identical.
 */
const CONTROLS_EXCEPT_LF = /[\p{Cc}\u2028\u2029]/gu;

/**
 * `escapeInvisible` for PLAIN text that is not JSON: the non-object preimage the
 * simulated path signs, which is free-form model text. Keeps real line feeds so a
 * multi-line body still reads as one, and escapes every other control so no two
 * different signed texts share an appearance.
 */
export function escapePlain(text: string): string {
  return escapePlainWith(text, INVISIBLE_OR_BIDI, isLineFeed).replace(
    CONTROLS_EXCEPT_LF,
    (c) => (c === "\n" ? c : escapeCodePoint(c)),
  );
}

/**
 * `escapePlain` for a field rendered on one line — an argument key, a tool name,
 * a connector — where a control character silently becomes whitespace.
 */
export function escapeInline(text: string): string {
  // A leading or trailing ASCII space leaves no ink, so `admin` and `admin ` are
  // one label on screen and two different keys to the target system. Inside the
  // text a space is visible between characters, so only the edges need it.
  return escapePlainWith(text, INVISIBLE_OR_BIDI_OR_CONTROL)
    .replace(/^ +/, (m) => "\\u0020".repeat(m.length))
    .replace(/ +$/, (m) => "\\u0020".repeat(m.length));
}

/**
 * Escapes EVERYTHING outside printable ASCII.
 *
 * The surgical escapes above target the tricks that are known and worth naming,
 * and each round of review has found another one, which is the nature of the
 * problem: two different strings drawing the same glyphs is a property of fonts
 * and shaping engines, not of a list this code can finish. Cyrillic `а` beside
 * Latin `a` needs no composition or invisibility at all.
 *
 * So there is also an exact mode. It is unreadable for ordinary text and that is
 * the point: within printable ASCII no two different strings can share an
 * appearance, so a reader who wants certainty rather than legibility can have
 * it, on demand, without imposing `revisi\u00F3n` on everyone who just wants to
 * read an argument.
 */
export function escapeToAscii(text: string): string {
  return [...text]
    .map((c) => {
      const cp = c.codePointAt(0) ?? 0;
      if (c === "\\") return "\\\\";
      // The space goes too, unlike everywhere else. A trailing one leaves no
      // ink, so `admin` and `admin ` would still share an appearance — and this
      // is the mode whose whole promise is that nothing does.
      return cp > 0x20 && cp <= 0x7e ? c : escapeCodePoint(c);
    })
    .join("");
}

/**
 * `escapePlain` for the reading aid shown BESIDE the signed text, where the
 * joiners that build emoji and ordinary Persian spelling are left intact. Safe
 * only because the exact rendering sits right next to it with them escaped.
 */
export function escapeReadable(text: string): string {
  // The aid keeps the natural composition; the exact rendering beside it shows it.
  return escapePlainWith(text, INVISIBLE_OR_BIDI, keptInReadingAid, false);
}

/** One argument of a proposed action, sliced out of the signed canonical JSON. */
export interface ApprovalArgEntry {
  key: string;
  /**
   * The value as it reads in normal mode: the SIGNED text verbatim — a 19-digit
   * account id reads exactly as hashed, and a string keeps its quotes — except
   * that a string's JSON escapes are decoded when decoding them hides nothing
   * (see below). The quotes stay in either case, so the string `"false"` never
   * reads as the boolean `false`.
   *
   * Decoding is where a rendering starts being able to differ from what the
   * signature covers, so it happens only when the decoded text contains nothing
   * `escapeInvisible` would mark: no invisible or direction-changing character,
   * no NFC-unstable composition. When it does contain one, the escape IS the
   * information — the point of showing `‮` is that the reader cannot see
   * it otherwise — and the signed rendering stays, with `readable` beside it.
   *
   * `valueSource` is always the untouched signed token, and exact mode renders
   * from it, so the signed text is one click away in every case.
   */
  value: string;
  /**
   * A human-readable rendering of a string value, present ONLY when the signed
   * rendering above is revealing something the reader would otherwise miss. It
   * is a reading aid beside the signed text, never a replacement for it: a
   * value whose escapes hide nothing gets no second line, it is simply read in
   * `value`.
   */
  readable?: string;
  /**
   * The key and value BEFORE any escaping. Exact mode renders from these rather
   * than from the escaped forms: re-escaping an escape would double its
   * backslashes and turn a faithful rendering into a confusing one.
   */
  keySource: string;
  valueSource: string;
}

export interface ParsedApprovalArgs {
  /** Key/value entries when the preimage is a JSON object; empty otherwise. */
  entries: ApprovalArgEntry[];
  /** True when the preimage decomposed as a JSON object (a gateway action). */
  structured: boolean;
  /**
   * The verbatim preimage: byte-for-byte the material `args_hash` was computed
   * over. Anything presented as "what gets signed" must come from here
   * untouched.
   */
  raw: string;
  /** True when the preimage is an argument-less object (`{}`). */
  isEmpty: boolean;
}

/**
 * Decodes an approval's `canonical_args`: the exact preimage of its `args_hash`,
 * server-verified against the hash before it is served. A gateway action's
 * preimage is a canonical JSON object and decomposes into key/value entries; the
 * legacy simulated path signs the proposal text itself, which stays raw-only.
 * Approvals recorded before arguments were persisted arrive empty and the card
 * falls back to the one-line `proposal`.
 *
 * The entries are NOT produced with JSON.parse: parsing would rebuild the
 * values through JavaScript numbers (rounding a 19-digit id to a different id
 * than the one signed) and reorder integer-like keys, so the table could show
 * something other than the signed text. Instead the preimage is tokenized and
 * each value is a verbatim slice of it, in its textual order; only strings are
 * decoded (lossless). Any input the tokenizer cannot fully account for comes
 * back unstructured, so the card shows the raw preimage rather than a table
 * that might lie.
 */
export function parseCanonicalArgs(
  canonicalArgs: string | undefined,
): ParsedApprovalArgs {
  const raw = canonicalArgs ?? "";
  if (raw === "") {
    return { entries: [], structured: false, raw, isEmpty: false };
  }
  // Validate the grammar with the language's own parser and throw its result
  // away. The tokenizer below is good at slicing but is not a JSON validator:
  // on its own it would accept `{"a":truth}` or `{"a":[1}` and render a table
  // for text that is not JSON at all. Letting JSON.parse decide what is valid,
  // and the tokenizer decide what the bytes are, keeps both jobs with whoever
  // does them correctly. The parsed value itself is deliberately unused: it is
  // exactly the lossy rebuild this whole path exists to avoid.
  try {
    JSON.parse(raw);
  } catch {
    return { entries: [], structured: false, raw, isEmpty: false };
  }
  const entries = sliceTopLevelObject(raw);
  if (entries === null) {
    return { entries: [], structured: false, raw, isEmpty: false };
  }
  return { entries, structured: true, raw, isEmpty: entries.length === 0 };
}

/**
 * Slices the top-level key/value pairs out of a JSON object, keeping each
 * value's exact source text. Returns null when the input is not a single JSON
 * object (or the scan cannot fully account for it), which callers must treat
 * as "show the raw preimage instead".
 */
function sliceTopLevelObject(raw: string): ApprovalArgEntry[] | null {
  let i = 0;

  function skipWhitespace() {
    while (
      raw[i] === " " ||
      raw[i] === "\t" ||
      raw[i] === "\n" ||
      raw[i] === "\r"
    ) {
      i++;
    }
  }

  /** Scans one string token from `i` (which must sit on `"`); returns its raw slice. */
  function scanString(): string | null {
    const start = i;
    if (raw[i] !== '"') return null;
    i++;
    while (i < raw.length) {
      if (raw[i] === "\\") {
        i += 2; // an escape never ends the string, whatever it escapes
        continue;
      }
      if (raw[i] === '"') {
        i++;
        return raw.slice(start, i);
      }
      i++;
    }
    return null; // unterminated
  }

  /** Scans one value token from `i`; returns its raw slice. */
  function scanValue(): string | null {
    const start = i;
    if (raw[i] === '"') {
      return scanString();
    }
    if (raw[i] === "{" || raw[i] === "[") {
      // Balanced-delimiter scan; strings are skipped whole so a `}` or `,`
      // inside one never miscounts.
      let depth = 0;
      while (i < raw.length) {
        if (raw[i] === '"') {
          if (scanString() === null) return null;
          continue;
        }
        if (raw[i] === "{" || raw[i] === "[") depth++;
        else if (raw[i] === "}" || raw[i] === "]") {
          depth--;
          if (depth === 0) {
            i++;
            return raw.slice(start, i);
          }
        }
        i++;
      }
      return null; // unbalanced
    }
    // Number or literal: runs until a structural delimiter.
    while (
      i < raw.length &&
      raw[i] !== "," &&
      raw[i] !== "}" &&
      raw[i] !== "]" &&
      raw[i] !== " " &&
      raw[i] !== "\t" &&
      raw[i] !== "\n" &&
      raw[i] !== "\r"
    ) {
      i++;
    }
    return i > start ? raw.slice(start, i) : null;
  }

  /** JSON-decodes a scanned string token; decoding a lone string is lossless. */
  function decodeString(token: string): string | null {
    try {
      const decoded: unknown = JSON.parse(token);
      return typeof decoded === "string" ? decoded : null;
    } catch {
      return null;
    }
  }

  skipWhitespace();
  if (raw[i] !== "{") return null;
  i++;
  skipWhitespace();
  const entries: ApprovalArgEntry[] = [];
  if (raw[i] === "}") {
    i++;
  } else {
    for (;;) {
      skipWhitespace();
      const keyToken = scanString();
      if (keyToken === null) return null;
      const decodedKey = decodeString(keyToken);
      if (decodedKey === null) return null;
      // The key is proposer-controlled too: a map key carrying a bidi override
      // reorders what the reader sees just as a value would.
      const key = escapeInline(decodedKey);
      skipWhitespace();
      if (raw[i] !== ":") return null;
      i++;
      skipWhitespace();
      const valueToken = scanValue();
      if (valueToken === null) return null;
      // The signed token is the value, with invisible and direction-changing
      // characters made visible (see escapeInvisible: that is more faithful to
      // the signed text, not less).
      const entry: ApprovalArgEntry = {
        key,
        value: escapeInvisible(valueToken),
        keySource: decodedKey,
        valueSource: valueToken,
      };
      if (valueToken[0] === '"') {
        const decoded = decodeString(valueToken);
        if (decoded === null) return null;
        // The natural reading is worth showing whenever it differs from what the
        // signed rendering shows: because of JSON escapes (a body full of \n) or
        // because the signed rendering escaped a joiner the aid keeps (an emoji,
        // a Persian ZWNJ). Comparing against the value without its quotes is
        // comparing like with like; when they match, a second line would be noise.
        const natural = escapeReadable(decoded);
        if (natural !== entry.value.slice(1, -1)) {
          // WHERE it goes depends on why they differ, and the test is whether the
          // decoded text has anything to hide. If `escapeInvisible` leaves it
          // untouched there is no invisible character, no bidi override and no
          // deceiving composition in it: the escapes were pure JSON notation, and
          // showing both forms hands the approver two copies of one text to
          // compare — the alert body arrives as a wall of `\n` above the
          // paragraphs it decodes to. So the decoded form simply becomes the
          // value, quotes kept.
          //
          // Otherwise the escapes are the information, and the pair stays: the
          // signed rendering shows what is actually there, the aid shows how it
          // means to look.
          if (escapeInvisible(decoded) === decoded) {
            entry.value = '"' + natural + '"';
          } else {
            entry.readable = natural;
          }
        }
      }
      entries.push(entry);
      skipWhitespace();
      if (raw[i] === ",") {
        i++;
        continue;
      }
      if (raw[i] === "}") {
        i++;
        break;
      }
      return null;
    }
  }
  skipWhitespace();
  // Anything left over means this was not exactly one object: refuse rather
  // than present a table that omits part of the signed material.
  return i === raw.length ? entries : null;
}

/**
 * The policy vocabulary is owned by the backend; only `approval_required` ever
 * reaches an approval (allow and deny never create one), so we localize it and
 * fall back to the raw token for anything unexpected.
 */
export function approvalPolicyLabel(policy: string | undefined): string {
  switch (normalizeStatus(policy)) {
    case "approval_required":
      return m.agents_approval_policy_approval_required();
    default:
      // An unknown token is backend text rendered on one line, so it gets the
      // same treatment as any other identifier rather than being trusted.
      return escapeInline(policy?.trim() ?? "") || "—";
  }
}

export function formatJson(value: unknown): string {
  if (value === undefined || value === null) return "";
  try {
    return JSON.stringify(value, null, 2);
  } catch {
    // Only reached for values JSON.stringify rejects (e.g. BigInt/circular).
    return "";
  }
}

/**
 * Names the person behind a subject (a user id) using the project's members, and
 * falls back to the raw subject when it cannot be resolved: a decider who left the
 * project, a service principal, or a caller without permission to list members. An
 * unresolvable subject is still a valid audit record, so it is shown as-is.
 */
export function subjectLabel(
  subject: string | undefined,
  names: Map<string, string> | undefined,
): string {
  const raw = subject?.trim() ?? "";
  if (raw === "") return "—";
  return names?.get(raw) ?? raw;
}

/**
 * The identity the proposed action will run as. In v1 runs carry either a user
 * subject or a service principal (see the design's identity model, §6/§14).
 */
export function runActorLabel(
  run: AgentRunData | undefined,
  names?: Map<string, string>,
): string {
  if (!run) return "—";
  if (run.actorServicePrincipal) return m.agents_actor_service_principal();
  return subjectLabel(run.actorSubject, names);
}

/**
 * The subject a lifecycle event attributes its decision to. The executor writes it
 * onto `run.resumed` / `run.rejected` (payload key `decided_by`), which is the only
 * place an API reader can learn who unblocked or stopped a run: the action ledger
 * that stores the same subject is not exposed.
 */
export function eventDecidedBy(payload: unknown): string {
  if (typeof payload !== "object" || payload === null) return "";
  const value = (payload as Record<string, unknown>).decided_by;
  return typeof value === "string" ? value : "";
}

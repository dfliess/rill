// @vitest-environment jsdom
//
// Redundant with the jsdom environment now configured in web-admin's vite
// config, and kept anyway: it states the requirement where the file that needs
// it lives, so moving or copying this spec does not silently lose it.
import { render, screen } from "@testing-library/svelte";
import { tick } from "svelte";
import { describe, expect, it, vi } from "vitest";

// The decision buttons pull in the runtime client and query client, which this
// file does not exercise: every case here renders an approval the caller cannot
// decide. Stub them so the card mounts on its own.
vi.mock("./ApproveDenyButtons.svelte", async () => ({
  default: (
    await import(
      "@rilldata/web-common/features/entity-management/__fixtures__/SlotPassthrough.svelte"
    )
  ).default,
}));

import AgentApprovalCard from "./AgentApprovalCard.svelte";

/**
 * These cover the WIRING, not the escaping — utils.spec.ts already pins what the
 * escapes produce. Three separate rounds of review found the same class of bug
 * right here: a helper written, unit-tested as a function, and then never
 * actually referenced by the template. A test over a pure function cannot see
 * that, which is the whole reason this file exists.
 */
function renderCard(canonicalArgs: string, canonicalArgsIsJson = true) {
  return render(AgentApprovalCard, {
    props: {
      approval: {
        approvalId: "a1",
        runId: "r1",
        status: "pending",
        argsHash: "sha256:abc",
        canonicalArgs,
        canonicalArgsIsJson,
        toolName: "mcp.crm.create_task",
        canDecide: false,
      },
    },
  });
}

describe("AgentApprovalCard rendering", () => {
  it("applies exact mode to the key", async () => {
    // The Cyrillic "а": no composition, no invisibility, nothing a surgical rule
    // catches. Only exact mode separates it from the Latin one, and it has to
    // reach the table for that to be worth anything.
    renderCard('{"\u0430dmin":1}');
    expect(screen.getByText("\u0430dmin")).toBeTruthy();

    screen.getByText("Show exact").click();
    await tick();

    expect(screen.queryByText("\u0430dmin")).toBeNull();
    expect(screen.getByText("\\u0430dmin")).toBeTruthy();
  });

  it("applies exact mode to the value too, not only to the key", async () => {
    // Separate from the key case on purpose: the first version of this file
    // asserted "keys and values" while only exercising the key, so swapping the
    // value back to its unrendered form did not fail anything.
    renderCard('{"user":"\u0430dmin"}');
    expect(screen.getByText('"\u0430dmin"')).toBeTruthy();

    screen.getByText("Show exact").click();
    await tick();

    expect(screen.getByText('"\\u0430dmin"')).toBeTruthy();
  });

  it("shows a real zero-width character escaped in the readable table", () => {
    // A REAL U+200C in the signed bytes, not the ASCII escape: with the escape,
    // source and rendered form coincide and the assertion proves nothing.
    renderCard('{"user":"ad\u200Cmin"}');
    expect(screen.getByText('"ad\\u200Cmin"')).toBeTruthy();
    // And the aid stays under it. This is the pair the rule below must not
    // collapse: here the escape is the only reason the character is visible.
    expect(screen.getByText("Readable form of the signed text")).toBeTruthy();
  });

  it("reads a multi-line body once, decoded, instead of twice", () => {
    // The signed escapes are notation here, so the value IS its decoding and no
    // aid is drawn under it. Asserted through the card, not the parser: the aid
    // is a template branch, and a parser test cannot see one left behind.
    // Read off textContent rather than getByText, whose default normalizer
    // collapses the very line breaks this is about.
    const { container } = renderCard('{"details":"Uno.\\n\\nDos."}');
    expect(container.querySelector("dl dd")?.textContent?.trim()).toBe(
      '"Uno.\n\nDos."',
    );
    expect(screen.queryByText("Readable form of the signed text")).toBeNull();
  });

  it("names an empty key instead of drawing an empty label", () => {
    renderCard('{"":1}');
    expect(screen.getByText("(empty key)")).toBeTruthy();
  });

  it("does not decompose a preimage the server did not certify as canonical", () => {
    // Free-form text that merely parses: the spaces are exactly what tells it
    // apart from its canonical twin, and a table would drop them.
    const { container } = renderCard('{ "a": 1 }', false);
    expect(container.querySelector("dl")).toBeNull();
    // Twice on purpose: as the block itself and inside the signed-JSON details.
    expect(screen.getAllByText('{ "a": 1 }').length).toBeGreaterThan(0);
  });
});

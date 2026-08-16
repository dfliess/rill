<script lang="ts">
  import MetadataLabel from "@rilldata/web-admin/features/scheduled-reports/metadata/MetadataLabel.svelte";
  import MetadataValue from "@rilldata/web-admin/features/scheduled-reports/metadata/MetadataValue.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import type { AgentApprovalData } from "../types";
  import type { ApprovalArgEntry } from "../utils";
  import {
    approvalPolicyLabel,
    escapeInline,
    escapeInvisible,
    escapePlain,
    escapeToAscii,
    formatDateTime,
    isApprovalPending,
    parseCanonicalArgs,
    subjectLabel,
  } from "../utils";
  import AgentApprovalStatusChip from "./AgentApprovalStatusChip.svelte";
  import ApproveDenyButtons from "./ApproveDenyButtons.svelte";

  let {
    approval,
    names,
  }: {
    approval: AgentApprovalData;
    /** Subject -> person, to name the decider instead of showing a raw user id. */
    names?: Map<string, string>;
  } = $props();

  let pending = $derived(isApprovalPending(approval.status));
  // The server resolves can_decide per approval, with its concrete action bound
  // (approval authority can discriminate by tool). When false the card still
  // shows the pending state — it is the run's record — but offers no buttons.
  let canDecide = $derived(!!approval.canDecide);
  // canonical_args is the exact preimage of args_hash, verified against it by
  // the server before it is served: what this card renders as "the arguments"
  // is byte-for-byte the material the approve call's args_hash binds to.
  let args = $derived(parseCanonicalArgs(approval.canonicalArgs));
  // Whether the preimage is JSON comes from the SERVER, which re-canonicalizes to
  // find out, not from whether it happens to parse here. Only the first is proof:
  // free-form text can be valid JSON too, and reading the exemption off its syntax
  // would let a literal "\u200C" and a real zero-width non-joiner render alike.
  // In real JSON the encoder already doubled the backslashes, so doubling again
  // would misrepresent the signed bytes.
  // The table is a decomposition of canonical JSON, so it is offered only when
  // the SERVER proved the preimage is exactly that. Free-form text that merely
  // parses — `{ "a": 1 }`, with spaces the canonical form would not have — is
  // shown as the block it is: decomposing it would drop the very characters that
  // tell it apart from its canonical twin.
  let structured = $derived(!!approval.canonicalArgsIsJson && args.structured);
  let shownRaw = $derived(
    exact
      ? escapeToAscii(args.raw)
      : approval.canonicalArgsIsJson
        ? escapeInvisible(args.raw)
        : escapePlain(args.raw),
  );
  let showPosition = $derived(
    typeof approval.total === "number" && approval.total > 1,
  );

  // A body-sized payload starts collapsed so a big action never buries the
  // decision; the threshold is on the signed bytes' length because measuring
  // rendered height is not worth the complexity here. This one applies only to
  // the single-block (non-object) preimage: a structured one collapses per
  // value instead, so no argument can be pushed out of sight by a fat sibling.
  const COLLAPSE_THRESHOLD = 1200;
  const VALUE_COLLAPSE_THRESHOLD = 400;
  let collapsible = $derived(args.raw.length > COLLAPSE_THRESHOLD);
  let expanded = $state(false);

  // Exact mode renders every non-ASCII code point as an escape. Unreadable on
  // purpose: the surgical escapes cover the tricks worth naming, but two strings
  // drawing the same glyphs is a property of fonts, not a list this code can
  // finish (Cyrillic "а" beside Latin "a" needs no trick at all). This gives
  // certainty on demand without imposing it on everyone who just wants to read.
  let exact = $state(false);
  // Rendered from the SOURCE text, not from the already-escaped form: escaping
  // an escape doubles its backslashes and turns a faithful rendering into a
  // puzzle.
  let showKey = $derived((e: ApprovalArgEntry) =>
    exact ? escapeToAscii(e.keySource) : e.key,
  );
  let showValue = $derived((e: ApprovalArgEntry) =>
    exact ? escapeToAscii(e.valueSource) : e.value,
  );

  let expandedValues = $state<Record<number, boolean>>({});
  function toggleValue(idx: number) {
    expandedValues[idx] = !expandedValues[idx];
  }

  // Anything still clipped is something the approver has not been shown. Two
  // values that differ only past the cut render identically, so signing while a
  // fragment is folded is signing on a prefix. Approve stays disabled until
  // every clipped fragment has been opened; Deny never does, because refusing
  // what you could not read is always a sound answer.
  let hasHiddenText = $derived(
    (structured &&
      args.entries.some(
        (e, idx) =>
          showValue(e).length > VALUE_COLLAPSE_THRESHOLD &&
          !expandedValues[idx],
      )) ||
      (!structured && collapsible && !expanded),
  );
</script>

<!-- A pending approval is a call to action (amber, with the decision buttons); a
     resolved one is the run's audit record (neutral, no buttons) and stays visible
     so the run keeps showing what was proposed, who decided it and when.

     The card reads top-down as the decision itself: WHAT will run (tool on
     connector), on whose request and under which policy, and with EXACTLY which
     arguments — the server-verified preimage of the hash the approval signs. -->
<div
  class="flex flex-col gap-y-4 rounded-lg p-4 border {pending
    ? 'border-amber-300 dark:border-amber-400/30 bg-amber-50 dark:bg-amber-500/10'
    : 'bg-surface-subtle'}"
>
  <div class="flex gap-x-2 items-center flex-wrap">
    <h2 class="text-fg-primary text-base font-semibold">
      {pending ? m.agents_run_approval_required() : m.agents_approval_title()}
    </h2>
    <AgentApprovalStatusChip status={approval.status} />
    {#if showPosition}
      <span class="text-xs text-fg-muted">
        {m.agents_approval_position({
          position: approval.position ?? 0,
          total: approval.total ?? 0,
        })}
      </span>
    {/if}
    <div class="grow"></div>
  </div>

  <!-- The action in one glance: the tool that will run, on which connector. -->
  <div class="flex items-center gap-x-2 flex-wrap min-w-0">
    <span
      class="font-mono text-sm font-semibold text-fg-primary bg-surface-card border rounded px-1.5 py-0.5 break-all whitespace-pre-wrap"
      title={m.agents_approval_proposed_action()}
    >
      {escapeInline(approval.toolName ?? "") || "—"}
    </span>
    {#if approval.connector}
      <span class="text-fg-muted text-sm" aria-hidden="true">·</span>
      <span
        class="text-fg-secondary text-sm whitespace-pre-wrap"
        title={m.agents_field_connector()}
      >
        {escapeInline(approval.connector)}
      </span>
    {/if}
  </div>

  <div class="flex flex-wrap gap-x-16 gap-y-4">
    <div class="flex flex-col gap-y-2">
      <MetadataLabel>{m.agents_approval_requested_by()}</MetadataLabel>
      <MetadataValue>{subjectLabel(approval.requestedBy, names)}</MetadataValue>
    </div>
    <div class="flex flex-col gap-y-2">
      <MetadataLabel>{m.agents_approval_policy()}</MetadataLabel>
      <MetadataValue>{approvalPolicyLabel(approval.policy)}</MetadataValue>
    </div>
    {#if !pending}
      <!-- The audit pair: who decided and when. `decidedBy` is empty for a decision
           no human made (a cancelled approval), so it renders as "—". -->
      <div class="flex flex-col gap-y-2">
        <MetadataLabel>{m.agents_approval_decided_by()}</MetadataLabel>
        <MetadataValue>{subjectLabel(approval.decidedBy, names)}</MetadataValue>
      </div>
      <div class="flex flex-col gap-y-2">
        <MetadataLabel>{m.agents_approval_decided_on()}</MetadataLabel>
        <MetadataValue>{formatDateTime(approval.decidedOn)}</MetadataValue>
      </div>
    {/if}
  </div>

  <div class="flex flex-col gap-y-2 min-w-0">
    <div class="flex items-baseline gap-x-2 flex-wrap">
      <MetadataLabel>{m.agents_approval_args_exact()}</MetadataLabel>
      {#if args.raw !== ""}
        <span class="text-xs text-fg-muted">
          {m.agents_approval_args_signed_note()}
        </span>
      {/if}
      {#if args.raw !== ""}
        <button
          type="button"
          class="text-xs text-primary-600"
          onclick={() => (exact = !exact)}
        >
          {exact
            ? m.agents_approval_args_exact_off()
            : m.agents_approval_args_exact_on()}
        </button>
      {/if}
      {#if hasHiddenText}
        <span class="text-xs text-amber-700 dark:text-amber-400">
          {m.agents_approval_args_expand_to_approve()}
        </span>
      {/if}
    </div>

    {#if args.raw === ""}
      <!-- Recorded before arguments were persisted: the preimage cannot be
           recovered, so say so and fall back to the one-line summary. -->
      <p class="text-xs text-fg-muted">{m.agents_approval_args_missing()}</p>
      {#if approval.proposal}
        <div class="flex flex-col gap-y-1">
          <span class="text-xs text-fg-muted"
            >{m.agents_approval_summary()}</span
          >
          <pre
            class="font-mono text-xs text-fg-primary whitespace-pre-wrap break-words bg-surface-card rounded border p-3 overflow-x-auto">{escapePlain(
              approval.proposal,
            )}</pre>
        </div>
      {/if}
    {:else}
      {#if structured && args.isEmpty}
        <p class="text-xs text-fg-secondary">
          {m.agents_approval_args_empty()}
        </p>
      {:else if structured}
        <!-- One row per argument, in the signed order. Every value is the signed
             text itself, sliced out of the preimage: numbers keep their exact
             digits, and a string keeps its quotes, so the string "false" never
             reads as the boolean false. When escapes make the signed text hard
             to read, the decoded form sits under it, labelled, clearly
             secondary.

             Long values collapse ONE AT A TIME, never the list as a whole. A
             block-level collapse would let whoever writes the arguments hide an
             entire argument below the fold — canonical JSON sorts the keys, so
             a fat `body` pushes `recipient` out of sight — and the decision
             buttons do not require scrolling past it. Collapsing per value
             keeps every key, and the fact that there is a value, on screen. -->
        <dl class="rounded border bg-surface-card divide-y">
          {#each args.entries as entry, idx (idx)}
            {@const long = showValue(entry).length > VALUE_COLLAPSE_THRESHOLD}
            {@const open = !!expandedValues[idx]}
            <div class="flex flex-col gap-y-1 px-3 py-2 min-w-0">
              <dt
                class="font-mono text-xs text-fg-muted whitespace-pre-wrap break-words"
              >
                {#if showKey(entry) === ""}
                  <span class="italic"
                    >{m.agents_approval_args_empty_key()}</span
                  >
                {:else}
                  {showKey(entry)}
                {/if}
              </dt>
              <!-- max-h and overflow-hidden must sit on the SAME element: a
                   capped height without the clip lets a long value paint over
                   the rows below it, which is worse than not collapsing. -->
              <div
                class="relative overflow-hidden"
                class:max-h-32={long && !open}
              >
                <dd
                  class="font-mono text-xs text-fg-primary whitespace-pre-wrap break-words"
                >
                  {showValue(entry)}
                </dd>
                {#if entry.readable !== undefined}
                  <div class="mt-1 flex flex-col gap-y-0.5">
                    <span class="text-[11px] text-fg-muted uppercase">
                      {m.agents_approval_args_readable()}
                    </span>
                    <dd
                      class="text-xs text-fg-secondary whitespace-pre-wrap break-words border-l-2 pl-2"
                    >
                      {entry.readable}
                    </dd>
                  </div>
                {/if}
                {#if long && !open}
                  <div
                    class="pointer-events-none absolute inset-x-0 bottom-0 h-8 bg-gradient-to-t from-surface-card to-transparent"
                  ></div>
                {/if}
              </div>
              {#if long}
                <button
                  type="button"
                  class="text-xs text-primary-600 self-start"
                  onclick={() => toggleValue(idx)}
                >
                  {open
                    ? m.agents_approval_collapse()
                    : m.agents_approval_show_all()}
                </button>
              {/if}
            </div>
          {/each}
        </dl>
      {:else}
        <!-- A preimage that is not a JSON object (the legacy simulated path
             signs the proposal text itself, and anything the tokenizer cannot
             fully account for lands here too): the block below IS the signed
             material. It is one opaque value, so collapsing it hides no
             argument's existence. -->
        <div
          class="relative rounded border bg-surface-card overflow-hidden"
          class:max-h-56={collapsible && !expanded}
        >
          <pre
            class="font-mono text-xs text-fg-primary whitespace-pre-wrap break-words p-3 overflow-x-auto">{shownRaw}</pre>
          {#if collapsible && !expanded}
            <div
              class="pointer-events-none absolute inset-x-0 bottom-0 h-10 bg-gradient-to-t from-surface-card to-transparent"
            ></div>
          {/if}
        </div>
        {#if collapsible}
          <button
            type="button"
            class="text-xs text-primary-600 self-start"
            onclick={() => (expanded = !expanded)}
          >
            {expanded
              ? m.agents_approval_collapse()
              : m.agents_approval_show_all()}
          </button>
        {/if}
      {/if}

      <!-- The raw signed bytes and their hash, always reachable whenever a
           preimage exists — including for `{}`, whose "no arguments" reading is
           itself a claim about signed material and must be checkable. The
           table above is a decomposition of these bytes, never a substitute. -->
      <details class="text-xs">
        <summary
          class="cursor-pointer text-fg-muted hover:text-fg-secondary select-none w-fit"
        >
          {m.agents_approval_args_raw()}
        </summary>
        <div class="mt-2 flex flex-col gap-y-2">
          <pre
            class="font-mono text-xs text-fg-primary whitespace-pre-wrap break-words bg-surface-card rounded border p-3 overflow-x-auto">{shownRaw}</pre>
          <div class="flex items-baseline gap-x-2 flex-wrap">
            <span class="text-fg-muted">{m.agents_approval_args_hash()}</span>
            <code class="font-mono text-fg-secondary break-all"
              >{approval.argsHash}</code
            >
          </div>
        </div>
      </details>
    {/if}
  </div>

  <!-- The decision sits AFTER the arguments, not in the header. Buttons placed
       above the material they authorize can be reached without scrolling past
       it — with enough arguments, or long enough ones, "Approve" stays on
       screen while the recipient is pages below. Putting them last makes the
       reading order and the decision order the same one. -->
  <!-- Without a verified preimage there is nothing to show, so the server refuses
       to approve (only to deny). Offering the button anyway would be a dead end:
       hide it and leave the way out. -->
  {#if pending && canDecide}
    <div class="flex justify-end border-t pt-3">
      <ApproveDenyButtons
        denyOnly={args.raw === "" || hasHiddenText}
        approvalId={approval.approvalId ?? ""}
        argsHash={approval.argsHash ?? ""}
        runId={approval.runId ?? ""}
      />
    </div>
  {/if}
</div>

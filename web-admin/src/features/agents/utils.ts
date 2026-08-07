import type { Color } from "@rilldata/web-common/components/tag/Tag.svelte";
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
 * An approval is resolved once a human decided it (or the run was cancelled).
 * Resolved approvals stay on the run: they are its audit trail, so the detail
 * view keeps showing what was proposed, who decided it and when.
 */
export function isApprovalResolved(status: string | undefined): boolean {
  const s = normalizeStatus(status);
  return s !== "" && s !== "pending";
}

/**
 * A run is finished once the backend stamps `finished_on`. This is the terminal
 * marker we key polling and the Cancel affordance off, independent of the status
 * vocabulary.
 */
export function isRunFinished(run: AgentRunData | undefined): boolean {
  return !!run?.finishedOn;
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

export interface ParsedProposal {
  isJson: boolean;
  /** Pretty-printed JSON when parseable, otherwise the raw string. */
  text: string;
  value: unknown;
}

/**
 * The approval `proposal` is the exact action and its normalized arguments,
 * carried as a string. It is usually JSON; we pretty-print it when it parses and
 * fall back to the raw text otherwise so the approver always sees something.
 */
export function parseProposal(proposal: string | undefined): ParsedProposal {
  const raw = proposal ?? "";
  if (raw.trim() === "") {
    return { isJson: false, text: "", value: undefined };
  }
  try {
    const value: unknown = JSON.parse(raw);
    return { isJson: true, text: JSON.stringify(value, null, 2), value };
  } catch {
    return { isJson: false, text: raw, value: raw };
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

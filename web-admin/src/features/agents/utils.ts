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
  expired: "gray",
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
    case "expired":
      return m.agents_approval_status_expired();
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
 * The identity the proposed action will run as. In v1 runs carry either a user
 * subject or a service principal (see the design's identity model, §6/§14).
 */
export function runActorLabel(run: AgentRunData | undefined): string {
  if (!run) return "—";
  if (run.actorServicePrincipal) return m.agents_actor_service_principal();
  return run.actorSubject || "—";
}

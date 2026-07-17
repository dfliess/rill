// A run's internal id is composed by the backend's `ComposeRunID` as
//   `<len>:<instanceId>/<len>:<agentName>/<idempotencyKey>`
// (see runtime). The `<len>` prefixes are length-of the following segment in
// bytes, which lets the key contain arbitrary characters (a trigger-fired run's
// key may itself contain `/`). We rebuild and parse that composite id on the
// client so the run-detail URL can carry only the human-meaningful parts
// (agentName + idempotency key) while the API still receives the full composite.
//
// The ids we handle are ASCII (instanceId is hex, agentName is a resource name),
// so a JS string's `.length` (UTF-16 code units) equals the byte length the
// backend used, and slicing by code units matches its byte offsets.

export interface ParsedRunId {
  instanceId: string;
  agentName: string;
  key: string;
}

/**
 * Rebuild the backend's composite run id from its parts. Mirrors `ComposeRunID`.
 */
export function composeRunId(
  instanceId: string,
  agentName: string,
  key: string,
): string {
  return `${instanceId.length}:${instanceId}/${agentName.length}:${agentName}/${key}`;
}

/**
 * Parse a composite run id by its length prefixes (not by splitting on `/`,
 * since the key may contain `/`). Throws if the string is not well-formed.
 */
export function parseRunId(composite: string): ParsedRunId {
  let pos = 0;

  // Read a `<len>:<value>` segment starting at `pos` and advance past it.
  function readLengthPrefixed(field: string): string {
    const colon = composite.indexOf(":", pos);
    if (colon === -1) {
      throw new Error(
        `invalid run id: missing ':' for ${field} in "${composite}"`,
      );
    }
    const lenStr = composite.slice(pos, colon);
    if (!/^\d+$/.test(lenStr)) {
      throw new Error(
        `invalid run id: bad length prefix "${lenStr}" for ${field} in "${composite}"`,
      );
    }
    const start = colon + 1;
    const end = start + Number(lenStr);
    if (end > composite.length) {
      throw new Error(
        `invalid run id: length ${lenStr} overruns "${composite}"`,
      );
    }
    pos = end;
    return composite.slice(start, end);
  }

  function expectSlash(field: string): void {
    if (composite[pos] !== "/") {
      throw new Error(
        `invalid run id: expected '/' after ${field} in "${composite}"`,
      );
    }
    pos += 1;
  }

  const instanceId = readLengthPrefixed("instanceId");
  expectSlash("instanceId");
  const agentName = readLengthPrefixed("agentName");
  expectSlash("agentName");
  const key = composite.slice(pos);
  return { instanceId, agentName, key };
}

/**
 * The URL of a run's detail page: `/-/agents/<agent>/runs/<idempotencyKey>`.
 * The instanceId and agentName are implicit in the path, so the composite id is
 * not carried in the URL. Accepts a run that exposes `agentName`+`idempotencyKey`
 * directly, or falls back to deriving them from a composite `runId`. Segments are
 * `encodeURIComponent`d because a trigger-fired key may contain `/` and other
 * reserved characters.
 */
export function runDetailPath(
  organization: string,
  project: string,
  run: { agentName?: string; idempotencyKey?: string; runId?: string },
): string {
  let agentName = run.agentName;
  let key = run.idempotencyKey;
  if ((!agentName || !key) && run.runId) {
    const parsed = parseRunId(run.runId);
    agentName ||= parsed.agentName;
    key ||= parsed.key;
  }
  return `/${organization}/${project}/-/agents/${encodeURIComponent(agentName ?? "")}/runs/${encodeURIComponent(key ?? "")}`;
}

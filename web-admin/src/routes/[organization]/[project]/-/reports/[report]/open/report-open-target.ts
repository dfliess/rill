/**
 * Where a report "open" link should send the recipient.
 *
 * AI reports append `?session_id=<id>` to the open link (see the runtime's
 * `buildAISessionURL`). Everything else opens the mapped explore dashboard.
 */
export type ReportOpenTarget =
  | { type: "ai"; sessionId: string }
  | { type: "explore" };

/**
 * Resolves the open target purely from the link's query parameters.
 *
 * We deep-link to the conversation only when the recipient can actually read it.
 * The `/-/ai/<id>` route resolves the conversation with the viewer's own session,
 * not with a report magic-link token. Creator-mode report links carry such a token
 * (it grants the shared explore snapshot but not conversation access), so we keep
 * the explore mapping there. Recipient-mode links have no usable token and their
 * recipients are project members with their own claims, so those land on the chat.
 */
export function resolveReportOpenTarget(
  searchParams: URLSearchParams,
): ReportOpenTarget {
  const sessionId = searchParams.get("session_id");
  const token = searchParams.get("token");

  if (sessionId && !token) {
    return { type: "ai", sessionId };
  }

  return { type: "explore" };
}

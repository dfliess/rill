import { parseRunId } from "@rilldata/web-admin/features/agents/run-id";
import { redirect } from "@sveltejs/kit";

// The run detail moved to `/-/agents/<agent>/runs/<idempotencyKey>`. Redirect the
// old composite-id URL so existing bookmarks keep working.
export const load = ({ params }) => {
  const { agentName, key } = parseRunId(params.run);
  throw redirect(
    307,
    `/${params.organization}/${params.project}/-/agents/${encodeURIComponent(agentName)}/runs/${encodeURIComponent(key)}`,
  );
};

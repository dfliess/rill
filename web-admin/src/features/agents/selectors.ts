import { createAdminServiceListProjectMemberUsers } from "@rilldata/web-admin/client";
import {
  createAgentServiceApproveAgentApprovalMutation,
  createAgentServiceCancelAgentRunMutation,
  createAgentServiceDenyAgentApprovalMutation,
  createAgentServiceGetAgentRun,
  createAgentServiceListAgentApprovals,
  createAgentServiceListAgentRuns,
  createAgentServiceListAgents,
  createAgentServiceStartAgentRunMutation,
  getAgentServiceListAgentRunsQueryOptions,
} from "@rilldata/web-common/runtime-client/v2/gen/agent-service";
import type { RuntimeClient } from "@rilldata/web-common/runtime-client/v2";
import { createQuery } from "@tanstack/svelte-query";
import { derived, type Readable } from "svelte/store";

// Runs and approvals can change out-of-band (a run advances, a new approval lands
// in the inbox), so the list views poll at a light cadence. Definitions change
// only on redeploy, so the catalog is not polled.
const LIST_REFETCH_INTERVAL = 5_000;

export function useAgents(client: RuntimeClient, enabled = true) {
  return createAgentServiceListAgents(
    client,
    {},
    {
      query: {
        enabled: enabled && !!client.instanceId,
        refetchOnMount: true,
      },
    },
  );
}

export function useAgentRuns(
  client: RuntimeClient,
  filters: { agentName?: string; status?: string } = {},
) {
  return createAgentServiceListAgentRuns(
    client,
    { agentName: filters.agentName, status: filters.status },
    {
      query: {
        enabled: !!client.instanceId,
        refetchOnMount: true,
        refetchInterval: LIST_REFETCH_INTERVAL,
      },
    },
  );
}

// Reactive variant of `useAgentRuns`: the request is a store, so the query
// re-fetches whenever the server-side filters (agent name, status) change. The
// list page bridges its URL-backed filter state into `request`.
export function useAgentRunsReactive(
  client: RuntimeClient,
  request: Readable<{ agentName?: string; status?: string }>,
) {
  return createQuery(
    derived(request, ($r) =>
      getAgentServiceListAgentRunsQueryOptions(
        client,
        { agentName: $r.agentName, status: $r.status },
        {
          query: {
            enabled: !!client.instanceId,
            refetchOnMount: true,
            refetchInterval: LIST_REFETCH_INTERVAL,
          },
        },
      ),
    ),
  );
}

export function useAgentRun(client: RuntimeClient, runId: string) {
  return createAgentServiceGetAgentRun(
    client,
    { runId },
    {
      query: {
        enabled: !!client.instanceId && !!runId,
        // Poll while the run is in flight; `finished_on` is the terminal marker.
        refetchInterval: (query) =>
          query.state.data?.run?.finishedOn ? false : 2_000,
      },
    },
  );
}

export function useAgentApprovals(
  client: RuntimeClient,
  filters: { status?: string; runId?: string } = {},
) {
  return createAgentServiceListAgentApprovals(
    client,
    { status: filters.status, runId: filters.runId },
    {
      query: {
        enabled: !!client.instanceId,
        refetchOnMount: true,
        refetchInterval: LIST_REFETCH_INTERVAL,
      },
    },
  );
}

// A run's actor and an approval's decider are stored as opaque subjects (the user's
// id), which says nothing to a human reading an audit trail. Resolve them against the
// project's members so the UI can name the person. The query is best-effort: a caller
// without permission to list members, or a decider who is no longer a member, falls
// back to the raw subject rather than failing the view (see `subjectLabel`).
export function useSubjectNames(organization: string, project: string) {
  return createAdminServiceListProjectMemberUsers(
    organization,
    project,
    undefined,
    {
      query: {
        enabled: !!organization && !!project,
        retry: false,
        select: (data) => {
          const names = new Map<string, string>();
          for (const member of data.members ?? []) {
            if (!member.userId) continue;
            names.set(
              member.userId,
              member.userName || member.userEmail || member.userId,
            );
          }
          return names;
        },
      },
    },
  );
}

// Mutations
export function useStartAgentRun(client: RuntimeClient) {
  return createAgentServiceStartAgentRunMutation(client);
}

export function useCancelAgentRun(client: RuntimeClient) {
  return createAgentServiceCancelAgentRunMutation(client);
}

export function useApproveAgentApproval(client: RuntimeClient) {
  return createAgentServiceApproveAgentApprovalMutation(client);
}

export function useDenyAgentApproval(client: RuntimeClient) {
  return createAgentServiceDenyAgentApprovalMutation(client);
}

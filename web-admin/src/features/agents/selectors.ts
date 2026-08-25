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
import { countPendingRuns } from "@rilldata/web-admin/features/notifications/app-badge";
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
  // The Act tab polls fast (live operations); passive surfaces like Home
  // pass a gentler interval to keep background request volume down.
  refetchInterval: number = LIST_REFETCH_INTERVAL,
) {
  return createAgentServiceListAgentRuns(
    client,
    { agentName: filters.agentName, status: filters.status },
    {
      query: {
        enabled: !!client.instanceId,
        refetchOnMount: true,
        refetchInterval,
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
  // See useAgentRuns: passive surfaces pass a gentler polling interval.
  refetchInterval: number = LIST_REFETCH_INTERVAL,
) {
  return createAgentServiceListAgentApprovals(
    client,
    { status: filters.status, runId: filters.runId },
    {
      query: {
        enabled: !!client.instanceId,
        refetchOnMount: true,
        refetchInterval,
      },
    },
  );
}

// The app badge follows the viewer around the whole project, not just the Act tab, so its query polls far more
// gently than the lists do. It refetches on focus as well: coming back to the app is exactly when a decision taken
// on another device has to show up here (kairos-cloud#147).
const BADGE_REFETCH_INTERVAL = 60_000;

// usePendingApprovalRunCount counts the runs of this project that are waiting on a decision the viewer may make.
// The server resolves can_decide per approval, so the count is already theirs alone; see countPendingRuns for why
// runs rather than approvals.
export function usePendingApprovalRunCount(client: RuntimeClient) {
  return createAgentServiceListAgentApprovals(
    client,
    { status: "pending" },
    {
      query: {
        enabled: !!client.instanceId,
        refetchOnMount: true,
        refetchOnWindowFocus: true,
        refetchInterval: BADGE_REFETCH_INTERVAL,
        select: (data) => countPendingRuns(data.approvals),
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

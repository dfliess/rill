/**
 * The URL of a run's detail page: `/-/agents/runs/<runId>`. The run id is an
 * opaque UUID (ADR-0034); it is encoded because a run recorded before that has
 * slashes in its id.
 */
export function runDetailPath(
  organization: string,
  project: string,
  run: { runId?: string },
): string {
  return `/${organization}/${project}/-/agents/runs/${encodeURIComponent(run.runId ?? "")}`;
}

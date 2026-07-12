import { getExploreName } from "@rilldata/web-common/features/explore-mappers/utils";
import { redirect } from "@sveltejs/kit";
import { resolveReportOpenTarget } from "./report-open-target";

export async function load({ parent, url, params }) {
  const { report } = await parent();
  const organization = params.organization;
  const project = params.project;
  const reportId = params.report;

  // AI reports deep-link to their conversation via `?session_id=<id>`.
  const openTarget = resolveReportOpenTarget(url.searchParams);
  if (openTarget.type === "ai") {
    throw redirect(
      307,
      `/${organization}/${project}/-/ai/${openTarget.sessionId}`,
    );
  }

  const executionTime = url.searchParams.get("execution_time");
  const token = url.searchParams.get("token");
  const exploreName =
    report.report.spec.annotations["explore"] ??
    getExploreName(report.report.spec.annotations?.web_open_path); // backwards compatibility

  return {
    organization,
    project,
    reportId,
    report,
    executionTime,
    token,
    exploreName,
  };
}

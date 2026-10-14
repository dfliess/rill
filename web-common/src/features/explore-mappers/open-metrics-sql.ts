import type { PartialMessage, Timestamp } from "@bufbuild/protobuf";
import {
  canMapHaving,
  mapMetricsResolverQueryToDashboard,
} from "@rilldata/web-common/features/explore-mappers/map-metrics-resolver-query-to-dashboard.ts";
import {
  generateExploreLink,
  getExploreSpecs,
} from "@rilldata/web-common/features/explore-mappers/open-query.ts";
import type { Schema as MetricsResolverQuery } from "@rilldata/web-common/runtime-client/gen/resolvers/metrics/schema.ts";
import type { RuntimeClient } from "@rilldata/web-common/runtime-client/v2";
import { queryServiceConvertMetricsSQLToQuery } from "@rilldata/web-common/runtime-client/v2/gen/query-service";

export type OpenMetricsSQLArgs = {
  client: RuntimeClient;
  // The metrics SQL of the alert or report.
  sql: string;
  // The explore the alert or report points at.
  exploreName: string;
  // The run the link was sent for: now() and rilltime expressions resolve against it.
  executionTime?: string | null;
  // The magic token of an emailed link. Its row filter already applies the query's WHERE.
  token?: string | null;
  organization: string;
  project: string;
};

/**
 * Builds the explore URL that shows the rows an alert or report written in metrics SQL saw.
 *
 * It throws when the explore can't show the query faithfully, rather than opening it
 * with fewer filters than the alert or report applied: callers fall back to the
 * alert or report page.
 */
export async function getExploreUrlForMetricsSQL({
  client,
  sql,
  exploreName,
  executionTime,
  token,
  organization,
  project,
}: OpenMetricsSQLArgs): Promise<string> {
  const { query: compiled } = await queryServiceConvertMetricsSQLToQuery(
    client,
    {
      sql,
      // The request is decoded from JSON, where a Timestamp is an RFC 3339 string.
      executionTime: (executionTime ?? undefined) as unknown as
        | PartialMessage<Timestamp>
        | undefined,
    },
  );
  const query = compiled as unknown as MetricsResolverQuery | undefined;
  if (!query?.metrics_view) {
    throw new Error("The metrics SQL did not compile to a metrics view query");
  }

  const { metricsViewSpec, exploreSpec } = await getExploreSpecs(
    client,
    exploreName,
  );
  if (exploreSpec.metricsView !== query.metrics_view) {
    throw new Error(
      `Explore "${exploreName}" does not show metrics view "${query.metrics_view}"`,
    );
  }
  if (!canMapHaving(metricsViewSpec, exploreSpec, query)) {
    throw new Error("The explore can't show the query's HAVING");
  }

  const exploreState = mapMetricsResolverQueryToDashboard(
    metricsViewSpec,
    exploreSpec,
    { query, ignoreWhere: !!token },
  );

  const url = new URL(
    await generateExploreLink(
      client,
      exploreState,
      metricsViewSpec,
      exploreSpec,
      exploreName,
      organization,
      project,
    ),
  );
  if (token) {
    url.pathname = `/${organization}/${project}/-/share/${token}/explore/${encodeURIComponent(exploreName)}`;
  }
  return url.toString();
}

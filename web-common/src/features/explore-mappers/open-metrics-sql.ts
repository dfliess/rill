import type { PartialMessage, Timestamp } from "@bufbuild/protobuf";
import {
  canMapHaving,
  mapMetricsResolverQueryToDashboard,
} from "@rilldata/web-common/features/explore-mappers/map-metrics-resolver-query-to-dashboard.ts";
import {
  generateExploreLink,
  getExploreSpecs,
} from "@rilldata/web-common/features/explore-mappers/open-query.ts";
import { ResourceKind } from "@rilldata/web-common/features/entity-management/resource-selectors.js";
import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";
import {
  getRuntimeServiceListResourcesQueryKey,
  runtimeServiceListResources,
} from "@rilldata/web-common/runtime-client";
import type { Schema as MetricsResolverQuery } from "@rilldata/web-common/runtime-client/gen/resolvers/metrics/schema.ts";
import type { RuntimeClient } from "@rilldata/web-common/runtime-client/v2";
import { queryServiceConvertMetricsSQLToQuery } from "@rilldata/web-common/runtime-client/v2/gen/query-service";

export type OpenMetricsSQLArgs = {
  client: RuntimeClient;
  // The metrics SQL of the alert or report.
  sql: string;
  // The explore the alert or report names, if any. Without one, the only explore on the
  // query's metrics view: the runtime grants an emailed link's token that same explore.
  exploreName?: string;
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

  exploreName ||= await findOnlyExploreForMetricsView(
    client,
    query.metrics_view,
  );
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

async function findOnlyExploreForMetricsView(
  client: RuntimeClient,
  metricsViewName: string,
): Promise<string> {
  const request = { kind: ResourceKind.Explore };
  const { resources } = await queryClient.fetchQuery({
    queryKey: getRuntimeServiceListResourcesQueryKey(
      client.instanceId,
      request,
    ),
    queryFn: ({ signal }) =>
      runtimeServiceListResources(client, request, { signal }),
  });
  const explores = (resources ?? []).filter(
    (r) => r.explore?.state?.validSpec?.metricsView === metricsViewName,
  );
  if (explores.length !== 1 || !explores[0].meta?.name?.name) {
    throw new Error(
      `Expected exactly one explore on metrics view "${metricsViewName}", found ${explores.length}`,
    );
  }
  return explores[0].meta.name.name;
}

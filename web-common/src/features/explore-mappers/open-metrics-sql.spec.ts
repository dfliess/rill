import {
  AD_BIDS_EXPLORE_WITH_3_MEASURES_DIMENSIONS_TIME_DIMENSION,
  AD_BIDS_IMPRESSIONS_MEASURE,
  AD_BIDS_METRICS_3_MEASURES_DIMENSIONS_WITH_TIME,
  AD_BIDS_PUBLISHER_DIMENSION,
} from "@rilldata/web-common/features/dashboards/stores/test-data/data.ts";
import type { RuntimeClient } from "@rilldata/web-common/runtime-client/v2";
import { beforeEach, describe, expect, it, vi } from "vitest";

const { convert, generateExploreLink, listResources } = vi.hoisted(() => ({
  convert: vi.fn(),
  generateExploreLink: vi.fn(),
  listResources: vi.fn(),
}));

vi.mock("@rilldata/web-common/runtime-client", async (importOriginal) => ({
  ...(await importOriginal<object>()),
  runtimeServiceListResources: listResources,
}));

vi.mock("@rilldata/web-common/runtime-client/v2/gen/query-service", () => ({
  queryServiceConvertMetricsSQLToQuery: convert,
}));
vi.mock("@rilldata/web-common/features/explore-mappers/open-query.ts", () => ({
  getExploreSpecs: () =>
    Promise.resolve({
      metricsViewSpec: AD_BIDS_METRICS_3_MEASURES_DIMENSIONS_WITH_TIME,
      exploreSpec: AD_BIDS_EXPLORE_WITH_3_MEASURES_DIMENSIONS_TIME_DIMENSION,
    }),
  generateExploreLink,
}));

import { queryClient } from "@rilldata/web-common/lib/svelte-query/globalQueryClient";
import { getExploreUrlForMetricsSQL } from "./open-metrics-sql";

const metricsView =
  AD_BIDS_EXPLORE_WITH_3_MEASURES_DIMENSIONS_TIME_DIMENSION.metricsView;
const having = {
  cond: {
    op: "gt",
    exprs: [{ name: AD_BIDS_IMPRESSIONS_MEASURE }, { val: 10 }],
  },
};
const args = {
  client: {} as RuntimeClient,
  sql: "SELECT publisher, impressions FROM AdBids HAVING impressions > 10",
  exploreName: "AdBids_explore",
  executionTime: "2022-02-01T07:00:00Z",
  organization: "org",
  project: "proj",
};

beforeEach(() => {
  queryClient.clear();
  listResources.mockReset();
  convert.mockReset();
  generateExploreLink.mockReset();
  generateExploreLink.mockResolvedValue(
    "http://localhost/org/proj/explore/AdBids_explore?f=x",
  );
});

describe("getExploreUrlForMetricsSQL", () => {
  it("compiles the SQL at the execution time and opens the explore", async () => {
    convert.mockResolvedValue({
      query: {
        metrics_view: metricsView,
        dimensions: [{ name: AD_BIDS_PUBLISHER_DIMENSION }],
        measures: [{ name: AD_BIDS_IMPRESSIONS_MEASURE }],
        having,
      },
    });
    const url = await getExploreUrlForMetricsSQL(args);
    expect(convert.mock.calls[0][1]).toEqual({
      sql: args.sql,
      executionTime: "2022-02-01T07:00:00Z",
    });
    expect(url).toBe("http://localhost/org/proj/explore/AdBids_explore?f=x");
    // The HAVING reaches the explore as a measure filter.
    const exploreState = generateExploreLink.mock.calls[0][1];
    expect(JSON.stringify(exploreState.whereFilter)).toContain("subquery");
  });

  it("opens the shared explore with an emailed link's token", async () => {
    convert.mockResolvedValue({
      query: {
        metrics_view: metricsView,
        dimensions: [{ name: AD_BIDS_PUBLISHER_DIMENSION }],
        measures: [{ name: AD_BIDS_IMPRESSIONS_MEASURE }],
        where: {
          cond: {
            op: "eq",
            exprs: [{ name: AD_BIDS_PUBLISHER_DIMENSION }, { val: "Google" }],
          },
        },
        having,
      },
    });
    const url = await getExploreUrlForMetricsSQL({ ...args, token: "tkn" });
    expect(new URL(url).pathname).toBe(
      "/org/proj/-/share/tkn/explore/AdBids_explore",
    );
    // The token already filters by the WHERE; the HAVING it doesn't carry stays.
    const exploreState = generateExploreLink.mock.calls[0][1];
    const filter = JSON.stringify(exploreState.whereFilter);
    expect(filter).not.toContain("Google");
    expect(filter).toContain("subquery");
  });

  it("refuses an explore on another metrics view", async () => {
    convert.mockResolvedValue({
      query: { metrics_view: "other", measures: [] },
    });
    await expect(getExploreUrlForMetricsSQL(args)).rejects.toThrow(
      "does not show metrics view",
    );
  });

  it("refuses a HAVING the explore can't show, rather than drop it", async () => {
    convert.mockResolvedValue({
      query: {
        metrics_view: metricsView,
        measures: [{ name: AD_BIDS_IMPRESSIONS_MEASURE }],
        having,
      },
    });
    await expect(getExploreUrlForMetricsSQL(args)).rejects.toThrow("HAVING");
  });

  describe("without an explore named", () => {
    const explore = (name: string, mv = metricsView) => ({
      meta: { name: { name } },
      explore: { state: { validSpec: { metricsView: mv } } },
    });
    beforeEach(() => {
      convert.mockResolvedValue({
        query: {
          metrics_view: metricsView,
          dimensions: [{ name: AD_BIDS_PUBLISHER_DIMENSION }],
          measures: [{ name: AD_BIDS_IMPRESSIONS_MEASURE }],
        },
      });
    });

    it("opens the only explore on the metrics view", async () => {
      listResources.mockResolvedValue({
        resources: [explore("other_explore", "other"), explore("the_one")],
      });
      await getExploreUrlForMetricsSQL({ ...args, exploreName: "" });
      expect(generateExploreLink.mock.calls[0][4]).toBe("the_one");
    });

    it("doesn't guess among several", async () => {
      listResources.mockResolvedValue({
        resources: [explore("one"), explore("two")],
      });
      await expect(
        getExploreUrlForMetricsSQL({ ...args, exploreName: "" }),
      ).rejects.toThrow("found 2");
    });
  });
});

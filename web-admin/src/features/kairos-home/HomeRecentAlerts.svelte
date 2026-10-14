<script lang="ts">
  import AlertIcon from "@rilldata/web-common/components/icons/AlertIcon.svelte";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { createMeasureValueFormatter } from "@rilldata/web-common/lib/number-formatting/format-measure-value";
  import { timeAgo } from "@rilldata/web-common/lib/time/relative-time";
  import type { MetricsViewSpecMeasure } from "@rilldata/web-common/runtime-client";
  import { V1AssertionStatus } from "@rilldata/web-common/runtime-client/gen/index.schemas";
  import { createRuntimeServiceListResources } from "@rilldata/web-common/runtime-client";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import { formatRunDate } from "@rilldata/web-admin/features/scheduled-reports/tableUtils";
  import HomeBox from "./HomeBox.svelte";

  let { organization, project }: { organization: string; project: string } =
    $props();

  const runtimeClient = useRuntimeClient();

  // Unfiltered on purpose: shares the cache entry with the Home menu, and the
  // metrics views ride along for formatting the fail row (no extra request).
  // Polled so a fresh trigger surfaces while Home is open; alerts move at
  // cron/refresh pace, so this stays far gentler than the agents lists (5s).
  const resourcesQuery = createRuntimeServiceListResources(
    runtimeClient,
    {},
    {
      query: {
        enabled: !!runtimeClient.instanceId,
        refetchOnMount: true,
        refetchInterval: 30_000,
      },
    },
  );

  // Measure specs across all metrics views, keyed by measure name, so fail row
  // columns that originate from a measure reuse its format (percentage,
  // currency, ...). Same-named measures across views are assumed equivalent.
  let measureByName = $derived(
    new Map<string, MetricsViewSpecMeasure>(
      ($resourcesQuery.data?.resources ?? [])
        .flatMap((res) => res.metricsView?.state?.validSpec?.measures ?? [])
        .filter(
          (measure): measure is MetricsViewSpecMeasure & { name: string } =>
            !!measure.name,
        )
        .map((measure) => [measure.name, measure]),
    ),
  );

  const compactNumber = new Intl.NumberFormat("en-US", {
    notation: "compact",
    maximumFractionDigits: 1,
  });

  // The fail row (a protobuf Struct) loses the alert query's column order:
  // keys arrive alphabetized, which under the 3-field cap can crowd out the
  // offending measure the query put right after its dimensions. Recover the
  // author's order from the SELECT list when the spec carries the SQL.
  function selectOrder(sql: unknown): string[] {
    if (typeof sql !== "string") return [];
    const list = /select\s+([\s\S]*?)\s+from\s/i.exec(sql)?.[1];
    if (!list) return [];
    return list
      .split(",")
      .map(
        (col) =>
          col
            .trim()
            .split(/\s+as\s+/i)
            .pop()
            ?.trim() ?? "",
      )
      .filter(Boolean);
  }

  // The offending values from the check that fired, in the alert query's
  // column order (dimension values then measures, formatted per their metrics
  // view spec). Capped so a verbose alert query cannot flood the row.
  function failRowSummary(
    failRow: Record<string, unknown> | undefined,
    sql: unknown,
  ) {
    if (!failRow) return "";
    const order = selectOrder(sql);
    const position = (name: string, value: unknown) => {
      const idx = order.indexOf(name);
      if (idx >= 0) return idx;
      // Unknown columns go last, dimension values before measures.
      return order.length + (typeof value === "string" ? 0 : 1);
    };
    return Object.entries(failRow)
      .filter(
        ([, value]) => typeof value === "string" || typeof value === "number",
      )
      .sort(([na, va], [nb, vb]) => position(na, va) - position(nb, vb))
      .map(([name, value]) => {
        if (typeof value === "string") return value;
        const measure = measureByName.get(name);
        return measure
          ? createMeasureValueFormatter(measure)(value as number)
          : compactNumber.format(value as number);
      })
      .slice(0, 3)
      .join(" · ");
  }

  type TriggeredAlert = {
    name: string;
    title: string;
    summary: string;
    firedOn: string;
    timeZone: string;
  };

  // The executive feed: alerts that actually fired, newest first. One entry per
  // alert (its latest triggered check) so a flapping alert doesn't flood the list.
  let allTriggered = $derived(
    ($resourcesQuery.data?.resources ?? [])
      .map((res): TriggeredAlert | undefined => {
        const fired = (res.alert?.state?.executionHistory ?? []).find(
          (e) => e.result?.status === V1AssertionStatus.ASSERTION_STATUS_FAIL,
        );
        const name = res.meta?.name?.name;
        const firedOn = fired?.finishedOn ?? fired?.executionTime;
        if (!fired || !name || !firedOn) return undefined;
        return {
          name,
          title: res.alert?.spec?.displayName || name,
          summary: failRowSummary(
            fired.result?.failRow,
            res.alert?.spec?.resolverProperties?.sql,
          ),
          firedOn,
          timeZone: res.alert?.spec?.refreshSchedule?.timeZone ?? "",
        };
      })
      .filter((a): a is TriggeredAlert => !!a)
      .sort(
        (a, b) => new Date(b.firedOn).getTime() - new Date(a.firedOn).getTime(),
      ),
  );
  // A glance, not the list: the header counts them all, "View all" has the rest.
  let triggered = $derived(allTriggered.slice(0, 4));

  // "Nothing fired" is only a fact once the list has answered.
  let loaded = $derived($resourcesQuery.isSuccess);

  let base = $derived(`/${organization}/${project}/-/alerts`);
</script>

<HomeBox
  title={m.alerts_home_recent_heading()}
  count={allTriggered.length}
  href={base}
  viewAllLabel={m.alerts_home_view_all()}
  emptyText={m.alerts_home_recent_empty()}
  {loaded}
>
  {#each triggered as alert (alert.name)}
    <!-- Every row here fired, so a status tag would say nothing: the row is
         what fired, the offending values, and when. -->
    <li>
      <a
        href={`${base}/${alert.name}`}
        class="flex items-center gap-x-3 px-4 py-2.5 hover:bg-surface-hover"
      >
        <span class="shrink-0 text-amber-500"><AlertIcon size="15px" /></span>
        <div class="flex min-w-0 grow flex-col">
          <span class="truncate text-sm font-medium text-fg-primary">
            {alert.title}
          </span>
          {#if alert.summary}
            <span class="truncate text-xs text-fg-secondary">
              {alert.summary}
            </span>
          {/if}
        </div>
        <span
          class="shrink-0 whitespace-nowrap text-xs text-fg-muted"
          title={formatRunDate(alert.firedOn, alert.timeZone)}
        >
          {timeAgo(new Date(alert.firedOn))}
        </span>
      </a>
    </li>
  {/each}
</HomeBox>

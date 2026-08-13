<script lang="ts">
  import AlertIcon from "@rilldata/web-common/components/icons/AlertIcon.svelte";
  import { Tag } from "@rilldata/web-common/components/tag";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import { createMeasureValueFormatter } from "@rilldata/web-common/lib/number-formatting/format-measure-value";
  import { timeAgo } from "@rilldata/web-common/lib/time/relative-time";
  import type { MetricsViewSpecMeasure } from "@rilldata/web-common/runtime-client";
  import { V1AssertionStatus } from "@rilldata/web-common/runtime-client/gen/index.schemas";
  import { createRuntimeServiceListResources } from "@rilldata/web-common/runtime-client";
  import { useRuntimeClient } from "@rilldata/web-common/runtime-client/v2";
  import { formatRunDate } from "@rilldata/web-admin/features/scheduled-reports/tableUtils";

  let { organization, project }: { organization: string; project: string } =
    $props();

  const runtimeClient = useRuntimeClient();

  // Unfiltered on purpose: shares the cache entry with KairosHomeCanvas, and the
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
  let triggered = $derived(
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
      )
      .slice(0, 5),
  );

  let base = $derived(`/${organization}/${project}/-/alerts`);
</script>

<!-- Hidden entirely when no alert has fired, so Home stays quiet. -->
{#if triggered.length > 0}
  <div class="flex flex-col gap-y-4">
    <h2
      class="flex items-center justify-between text-xl font-semibold text-fg-secondary"
    >
      {m.alerts_home_recent_heading()}
      <a class="text-sm font-normal text-primary-600" href={base}>
        {m.alerts_home_view_all()}
      </a>
    </h2>
    <ul class="flex flex-col rounded-lg border divide-y overflow-hidden">
      {#each triggered as alert (alert.name)}
        <li>
          <!-- Two lines, matching the pending-approvals inbox and the dashboards
               listing: title on its own line, status and context below it. A
               single line made the four fields fight for the width, and on a
               phone the title was the one that lost. -->
          <a
            href={`${base}/${alert.name}`}
            class="flex items-center gap-x-2.5 px-4 py-3 group hover:bg-surface-hover"
          >
            <span class="shrink-0 text-fg-secondary">
              <AlertIcon size="15px" />
            </span>
            <div class="flex flex-col gap-y-1 min-w-0 grow">
              <span
                class="text-fg-primary text-sm font-semibold truncate min-w-0 group-hover:text-accent-primary-action"
              >
                {alert.title}
              </span>
              <!-- The status line does not wrap: the fail row summary is the
                   only field that can run long, so it absorbs the squeeze by
                   truncating and the row stays two lines on a phone. Wrapping
                   instead left the tag stranded on a line of its own. -->
              <div
                class="flex items-center gap-x-2 text-fg-secondary text-xs min-w-0"
              >
                <Tag color="blue">{m.alert_status_triggered()}</Tag>
                {#if alert.summary}
                  <span class="truncate min-w-0">{alert.summary}</span>
                  <span class="shrink-0" aria-hidden="true">·</span>
                {/if}
                <span
                  class="shrink-0 whitespace-nowrap"
                  title={formatRunDate(alert.firedOn, alert.timeZone)}
                >
                  {timeAgo(new Date(alert.firedOn))}
                </span>
              </div>
            </div>
          </a>
        </li>
      {/each}
    </ul>
  </div>
{/if}

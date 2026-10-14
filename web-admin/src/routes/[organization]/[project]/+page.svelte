<script lang="ts">
  // Kairos: the project home is fully owned by KairosProjectHome (welcome
  // heading, executive strip, chat, approvals, alerts, dashboards). Keeping
  // this route a shell confines upstream-merge conflicts to a single mount.
  // With the new Home on (the default), this URL opens it instead; the
  // previous Home lives at /-/home-classic (kairos-cloud#231).
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import KairosProjectHome from "@rilldata/web-admin/features/kairos-home/KairosProjectHome.svelte";
  import { featureFlags } from "@rilldata/web-common/features/feature-flags";
  import { onMount } from "svelte";

  const { kairosNewHome } = featureFlags;

  // Wait for the project's flags: the client-side default is a guess, and a
  // project that turned the new Home off must not be bounced away from here.
  let flagsReady = $state(false);
  onMount(() => {
    void featureFlags.ready.then(() => (flagsReady = true));
  });

  $effect(() => {
    if (!flagsReady || !$kairosNewHome) return;
    const base = $page.url.pathname.replace(/\/$/, "");
    void goto(`${base}/-/home-v2${$page.url.search}`, { replaceState: true });
  });
</script>

{#if flagsReady && !$kairosNewHome}
  <KairosProjectHome />
{/if}

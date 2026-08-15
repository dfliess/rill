<script lang="ts">
  // Kairos wordmark, mirroring the prop shape of Rill.svelte so it can stand in
  // wherever the product used to draw Rill's own logo.
  //
  // Brand rule (same as HeaderLogo.svelte): the primary lockup (anthracite +
  // blue block) on light backgrounds, the inverted one (white + blue block) on
  // dark. The files are served from the marketing site rather than bundled,
  // which is what the header already does; the CSP allows https: images.
  let {
    width = "64",
    height = "40",
    mode = "adapt",
  }: {
    width?: string;
    height?: string;
    mode?: "light" | "dark" | "adapt";
  } = $props();

  const lockupLight =
    "https://kairosagentica.com/img/logo/kairos-lockup-horizontal.svg";
  const lockupDark =
    "https://kairosagentica.com/img/logo/kairos-lockup-horizontal-blanco.svg";
</script>

{#if mode === "adapt"}
  <img
    src={lockupLight}
    alt="Kairos"
    {width}
    {height}
    class="object-contain dark:hidden"
  />
  <img
    src={lockupDark}
    alt="Kairos"
    {width}
    {height}
    class="object-contain hidden dark:block"
  />
{:else}
  <img
    src={mode === "dark" ? lockupDark : lockupLight}
    alt="Kairos"
    {width}
    {height}
    class="object-contain"
  />
{/if}

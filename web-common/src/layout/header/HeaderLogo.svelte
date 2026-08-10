<script lang="ts">
  export let href: string = "/";
  export let logoUrl: string | undefined = undefined;

  // Brand rule: the primary lockup (anthracite + blue block) on light,
  // the inverted one (white + blue block) on dark.
  const kairosLogoLight =
    "https://kairosagentica.com/img/logo/kairos-lockup-horizontal.svg";
  const kairosLogoDark =
    "https://kairosagentica.com/img/logo/kairos-lockup-horizontal-blanco.svg";
</script>

<a {href} class="grid place-content-center rounded pl-2 pr-2 shrink-0">
  {#if logoUrl}
    <img src={logoUrl} alt="logo" class="h-7" />
  {:else}
    <img src={kairosLogoLight} alt="logo" class="h-7 dark:hidden" />
    <img src={kairosLogoDark} alt="logo" class="h-7 hidden dark:block" />
  {/if}
</a>

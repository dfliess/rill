<script lang="ts">
  import KairosSymbol from "@rilldata/web-common/components/icons/RillLogoSquareNegative.svelte";

  export let href: string = "/";
  export let logoUrl: string | undefined = undefined;

  // Brand rule: the primary lockup (anthracite + blue block) on light,
  // the inverted one (white + blue block) on dark. Below `sm` only the
  // symbol is drawn, so the header fits one row on a phone; 20px is the
  // symbol's height inside the 28px lockup.
  const kairosLogoLight =
    "https://kairosagentica.com/img/logo/kairos-lockup-horizontal.svg";
  const kairosLogoDark =
    "https://kairosagentica.com/img/logo/kairos-lockup-horizontal-blanco.svg";
</script>

<a {href} class="grid place-content-center rounded pl-2 pr-2 shrink-0">
  {#if logoUrl}
    <img src={logoUrl} alt="logo" class="h-7" />
  {:else}
    <span class="h-7 grid place-content-center sm:hidden">
      <KairosSymbol size="20px" />
    </span>
    <img
      src={kairosLogoLight}
      alt="logo"
      class="h-7 hidden sm:block sm:dark:hidden"
    />
    <img src={kairosLogoDark} alt="logo" class="h-7 hidden sm:dark:block" />
  {/if}
</a>

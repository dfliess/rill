<script lang="ts">
  export let href: string = "/";
  export let logoUrl: string | undefined = undefined;

  const kairosLogo =
    "https://kairosagentica.com/img/logo/kairos-lockup-horizontal-mono-azul.svg";
</script>

<a
  {href}
  class="grid place-content-center rounded pl-2 pr-2 shrink-0"
>
  <img src={logoUrl || kairosLogo} alt="logo" class="h-7" />
</a>

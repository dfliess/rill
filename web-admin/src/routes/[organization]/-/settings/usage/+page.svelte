<script lang="ts">
  import { createAdminServiceGetBillingProjectCredentials } from "@rilldata/web-admin/client";
  import { ResourceKind } from "@rilldata/web-common/features/entity-management/resource-selectors";
  import Spinner from "@rilldata/web-common/features/entity-management/Spinner.svelte";
  import { EntityStatus } from "@rilldata/web-common/features/entity-management/types";
  import { themeControl } from "@rilldata/web-common/features/themes/theme-control";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  let embedLoaded = $state(false);

  // El hueco del iframe se tapa hasta que el panel se ha pintado de verdad. Dar
  // fondo al elemento no vale: el documento de dentro lo cubre con su blanco
  // hasta que la aplicación monta y aplica el tema. Y `load` tampoco basta,
  // porque salta cuando el documento está listo pero aún vacío.
  const REVEAL_TIMEOUT_MS = 15000;

  function revealWhenPainted(event: Event) {
    const frame = event.currentTarget as HTMLIFrameElement;
    const deadline = performance.now() + REVEAL_TIMEOUT_MS;
    // El documento del embed nace blanco y se oscurece cuando la aplicación
    // monta y resuelve el tema, así que la señal es que el iframe lleve ya la
    // misma marca de tema que la página que lo contiene. En claro no hay nada
    // que esperar, porque no hay salto que tapar.
    const hostIsDark = document.documentElement.classList.contains("dark");
    const check = () => {
      const painted =
        !hostIsDark ||
        !!frame.contentDocument?.documentElement.classList.contains("dark");
      // Se descubre igualmente al agotar el plazo: más vale enseñar un panel a
      // medias que un hueco eterno si el embed falla.
      if (painted || performance.now() > deadline) embedLoaded = true;
      else requestAnimationFrame(check);
    };
    check();
  }

  // Kairos: donde upstream embebe el portal de Orb, nosotros embebemos el
  // consumo que mide la propia plataforma (ADR-0017).
  //
  // El admin acuña un token de vida corta contra el proyecto de métricas con el
  // atributo organization_id de ESTA organización, y la política de la metrics
  // view filtra por él. O sea que el recorte por Tenant no lo hace este
  // componente: lo hace el runtime, y no hay forma de pedir más de lo tuyo
  // manipulando la URL. La RPC además exige ManageOrg, así que esto solo lo ven
  // los administradores de la organización.
  let creds = $derived(
    createAdminServiceGetBillingProjectCredentials({ org: data.organization }),
  );

  // El embed vive en el mismo origen, así que no hace falta abrir la CSP a
  // terceros (ADR-0006).
  let embedUrl = $derived.by(() => {
    const c = $creds.data;
    if (!c?.runtimeHost || !c?.instanceId || !c?.accessToken) return "";
    const params = new URLSearchParams({
      instance_id: c.instanceId,
      runtime_host: c.runtimeHost,
      access_token: c.accessToken,
      resource: "consumo_tenant",
      // El embed compara contra ResourceKind, no contra la cadena "canvas":
      // pasar "canvas" a secas cae al explorador y da un 404.
      type: ResourceKind.Canvas,
      hide_navigation_bar: "true",
      // Sin este parámetro el embed fuerza claro (init-embed-public-api.ts),
      // y se vería un panel blanco dentro de una app oscura. Pasamos el modo ya
      // resuelto, no la preferencia: si el usuario tiene "sistema", el iframe
      // hereda lo que esté viendo de verdad.
      theme_mode: $themeControl,
    });
    return `/-/embed?${params.toString()}`;
  });
</script>

<section class="usage-page">
  <h1 class="text-xl font-semibold text-fg-primary mb-2">
    {m.billing_usage_title()}
  </h1>

  {#if $creds.isLoading}
    <Spinner status={EntityStatus.Running} size="16px" />
  {:else if embedUrl}
    <div class="embed-frame">
      {#if !embedLoaded}
        <div class="embed-placeholder">
          <Spinner status={EntityStatus.Running} size="16px" />
        </div>
      {/if}
      <iframe
        src={embedUrl}
        title={m.billing_usage_title()}
        class="embed"
        class:opacity-0={!embedLoaded}
        onload={revealWhenPainted}
      ></iframe>
    </div>
  {/if}
</section>

<style lang="postcss">
  .usage-page {
    @apply w-full;
  }

  .embed-frame {
    @apply relative w-full h-[1440px];
  }

  .embed-placeholder {
    @apply absolute inset-0 grid place-content-center bg-surface-base;
  }

  .embed {
    @apply w-full h-full border-0 bg-surface-base transition-opacity;
  }
</style>

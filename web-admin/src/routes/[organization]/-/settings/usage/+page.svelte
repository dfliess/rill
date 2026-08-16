<script lang="ts">
  import { createAdminServiceGetBillingProjectCredentials } from "@rilldata/web-admin/client";
  import { ResourceKind } from "@rilldata/web-common/features/entity-management/resource-selectors";
  import Spinner from "@rilldata/web-common/features/entity-management/Spinner.svelte";
  import { EntityStatus } from "@rilldata/web-common/features/entity-management/types";
  import { themeControl } from "@rilldata/web-common/features/themes/theme-control";
  import { m } from "@rilldata/web-common/lib/i18n/gen/messages";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

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
    <iframe src={embedUrl} title={m.billing_usage_title()} class="embed"
    ></iframe>
  {/if}
</section>

<style lang="postcss">
  .usage-page {
    @apply w-full;
  }

  .embed {
    @apply w-full h-[1000px] border-0;
    /* Un iframe sin fondo lo pinta el navegador de blanco hasta que el
       documento de dentro dibuja el suyo, así que al cargar la página daba un
       fogonazo blanco dentro de una app oscura. Con el color de la superficie
       el hueco es del mismo color que lo que va a aparecer, y el cambio deja
       de verse. Sigue al tema, igual que el propio embed. */
    @apply bg-surface-base;
  }
</style>

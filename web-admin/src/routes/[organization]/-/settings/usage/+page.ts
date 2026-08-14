import type { PageLoad } from "./$types";

// Kairos: upstream redirige fuera de esta página si no hay portal de Orb y el
// plan no es de pago. Nosotros no facturamos con Orb: el consumo lo mide la
// propia plataforma y se lee del proyecto de métricas (ADR-0017), así que la
// página vale para cualquier organización. Las credenciales del embed se piden
// en el cliente, acotadas a la organización por el propio admin.
export const load: PageLoad = ({ params: { organization } }) => {
  return { organization };
};

import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

// Kairos: la pestaña de Facturación está retirada del menú (ver el layout de
// settings). Esta redirección cierra la puerta de atrás: sin biller conectado
// la página no encuentra suscripción y cae en su plan por defecto, que muestra
// la tarifa comercial de Rill, de modo que llegar por URL directa enseñaría
// unos precios que no son los nuestros.
export const load: PageLoad = ({ params: { organization } }) => {
  throw redirect(307, `/${organization}/-/settings`);
};

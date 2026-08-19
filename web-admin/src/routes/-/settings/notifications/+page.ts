import { redirectToLogin } from "@rilldata/web-admin/client/redirect-utils";
import type { PageLoad } from "./$types";

// Personal settings, so they hang off `/-/` rather than off an organization.
// Without a session there is nothing to show and every query on the page would
// 401 on its own, so send the visitor to log in first.
export const load: PageLoad = async ({ parent }) => {
  const { user } = await parent();
  if (!user) redirectToLogin();
};

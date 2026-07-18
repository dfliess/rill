<script lang="ts">
  import { goto } from "$app/navigation";
  import {
    adminServiceGetCurrentUser,
    adminServiceGetOrganizationNameForDomain,
    adminServiceListOrganizations,
    adminServiceListProjectsForOrganization,
  } from "@rilldata/web-admin/client";
  import { onMount } from "svelte";
  import { getActiveOrgLocalStorageKey } from "./local-storage";
  import { ADMIN_URL, CANONICAL_ADMIN_URL } from "../../../client/http-client";

  let showWelcomeMessage = false;

  // Land on the org overview, unless the user can see exactly one project in it,
  // in which case go straight to that project. This only runs on the post-login
  // "/" landing; navigating to "/{org}" on purpose still shows the project list.
  async function gotoOrgOrSoleProject(org: string | undefined) {
    if (org) {
      try {
        // pageSize 2 is enough to distinguish "exactly one" from "more than one".
        const { projects } = await adminServiceListProjectsForOrganization(
          org,
          {
            pageSize: 2,
          },
        );
        if (projects?.length === 1) {
          await goto(`/${org}/${projects[0].name}`);
          return;
        }
      } catch (e) {
        console.error("Failed to list projects for sole-project redirect", e);
        // Fall back to the org overview.
      }
    }
    await goto(`/${org}`);
  }

  onMount(async () => {
    // Scenario 1: If running on a custom domain, redirect to the org for the custom domain.
    if (ADMIN_URL !== CANONICAL_ADMIN_URL) {
      try {
        const res = await adminServiceGetOrganizationNameForDomain(
          window.location.hostname,
        );
        await gotoOrgOrSoleProject(res.name);
        return;
      } catch (e) {
        console.error("Failed to get organization for custom domain", e);
        // Fall back to the default behavior
      }
    }

    // Get the activeOrg local storage key for the current user
    const userId = (await adminServiceGetCurrentUser())?.user?.id;
    const activeOrgLocalStorageKey = getActiveOrgLocalStorageKey(userId);

    // Scenario 2: User has an activeOrg in localStorage
    const activeOrg = localStorage.getItem(activeOrgLocalStorageKey);
    if (activeOrg) {
      await gotoOrgOrSoleProject(activeOrg);
      return;
    }

    const orgs = (await adminServiceListOrganizations()).organizations;

    // Scenario 3: User has no activeOrg in localStorage, but does belong to an org
    if (orgs.length > 0) {
      await gotoOrgOrSoleProject(orgs[0].name);
      return;
    }

    // Scenario 4: User does not belong to an org
    showWelcomeMessage = true;
  });
</script>

{#if showWelcomeMessage}
  <slot />
{/if}

<!-- PERSONAL SETTINGS LAYOUT -->

<script lang="ts">
  import Header from "@rilldata/web-common/layout/header/Header.svelte";
  import HeaderLogo from "@rilldata/web-common/layout/header/HeaderLogo.svelte";
  import AvatarButton from "@rilldata/web-admin/features/authentication/AvatarButton.svelte";

  // These pages hang off `/-/` because they belong to the user, not to an
  // organization, so the root layout's OrgHeader (which needs an organization
  // in the URL) never renders for them. Without a header of their own they
  // would be a dead end: nothing to click to get back into the app.
</script>

<Header borderBottom>
  <HeaderLogo href="/" />
  <div class="grow"></div>
  <AvatarButton />
</Header>

<slot />

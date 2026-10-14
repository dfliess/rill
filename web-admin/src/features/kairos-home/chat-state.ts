import { createChatActions } from "@rilldata/web-common/features/chat/layouts/sidebar/sidebar-store";
import { localStorageStore } from "@rilldata/web-common/lib/store-utils/local-storage";
import { writable } from "svelte/store";

/**
 * Open state of the Home AI panel. Shared by the page, which renders the
 * panel, and the project header, which holds its toggle. Its own store:
 * opening it here must not open the side panel on every dashboard.
 */
export const homeChatOpen = writable(false);
export const homeChatActions = createChatActions(homeChatOpen);

/** Side by side with the canvas from this width; over it below. */
export const HOME_SIDE_BY_SIDE_QUERY = "(min-width: 1280px)";

/** Width of the dashboards menu, resizable like the AI panel and remembered. */
export const HOME_MENU_WIDTH = { BASIS: 224, MIN: 160, MAX: 420 };
export const homeMenuWidth = localStorageStore<number>(
  "home-menu-width",
  HOME_MENU_WIDTH.BASIS,
);

import { writable } from "svelte/store";

export type SteamPageTab = "switcher" | "games";

/** Which half of the Steam platform page is showing. Read by the page toolbar and the page. */
export const steamPageTab = writable<SteamPageTab>("switcher");

export function resetSteamPageTab(): void {
  steamPageTab.set("switcher");
}

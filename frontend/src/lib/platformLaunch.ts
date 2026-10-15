import { get, writable } from "svelte/store";
import * as PlatformService from "../../bindings/TcNo-Acc-Switcher/internal/platform/platformservice.js";
import { route } from "../stores/nav";
import { t } from "../stores/i18n";
import { openFolderPicker } from "../stores/modal";
import { pushToast } from "../stores/toast";

export const platformNavigationBusy = writable(false);

export async function openPlatformPage(name: string): Promise<boolean> {
  if (get(platformNavigationBusy)) return false;
  platformNavigationBusy.set(true);
  try {
    return await resolveAndOpenPlatform(name);
  } finally {
    platformNavigationBusy.set(false);
  }
}

async function resolveAndOpenPlatform(name: string): Promise<boolean> {
  const tr = get(t);
  const result = await PlatformService.ResolvePlatformLaunch(name);
  if (result.ok) {
    if (result.foundViaShortcut) {
      pushToast({ type: "success", title: "", message: tr("Toast_FoundExeViaShortcut"), duration: 30000 });
    }
    route.set({ page: "platform", platformName: name });
    return true;
  }
  if (result.needsManualLocate) {
    const picked = await openFolderPicker({
      title: tr("Modal_Title_LocatePlatform", { platform: name }),
      body: `<p>${tr("Modal_LocatePlatform", { platformExe: result.soughtExeName })}</p>`,
      initialPath: result.initialPath ?? "",
      dirsOnly: false,
      soughtFilename: result.soughtExeName,
      positiveLabel: tr("Modal_Button_Select"),
    });
    if (picked) {
      await PlatformService.ConfirmPlatformExePath(name, picked);
      route.set({ page: "platform", platformName: name });
      return true;
    }
  }
  return false;
}

import { get } from "svelte/store";
import { beforeEach, describe, expect, it, vi } from "vitest";

const native = vi.hoisted(() => ({
  ResolvePlatformLaunch: vi.fn(),
  ConfirmPlatformExePath: vi.fn(),
  openFolderPicker: vi.fn(),
  navigate: vi.fn(),
  pushToast: vi.fn(),
}));

vi.mock("../../bindings/TcNo-Acc-Switcher/internal/platform/platformservice.js", () => native);
vi.mock("../stores/modal", () => ({ openFolderPicker: native.openFolderPicker }));
vi.mock("../stores/nav", () => ({ route: { set: native.navigate } }));
vi.mock("../stores/toast", () => ({ pushToast: native.pushToast }));
vi.mock("../stores/i18n", () => ({
  t: { subscribe: (run: (translate: (key: string) => string) => void) => {
    run((key) => key);
    return () => {};
  } },
}));

import { openPlatformPage, platformNavigationBusy } from "./platformLaunch";

beforeEach(() => {
  vi.resetAllMocks();
  platformNavigationBusy.set(false);
});

describe("opening a platform from either navigation surface", () => {
  it("opens the account page after the executable resolves", async () => {
    native.ResolvePlatformLaunch.mockResolvedValue({ ok: true });
    expect(await openPlatformPage("Steam")).toBe(true);
    expect(native.navigate).toHaveBeenCalledWith({ page: "platform", platformName: "Steam" });
    expect(native.openFolderPicker).not.toHaveBeenCalled();
  });

  it("stays on the current page when manual location is cancelled", async () => {
    native.ResolvePlatformLaunch.mockResolvedValue({ needsManualLocate: true });
    native.openFolderPicker.mockResolvedValue(null);
    expect(await openPlatformPage("Steam")).toBe(false);
    expect(native.navigate).not.toHaveBeenCalled();
    expect(native.ConfirmPlatformExePath).not.toHaveBeenCalled();
  });

  it("confirms a selected executable before opening the account page", async () => {
    native.ResolvePlatformLaunch.mockResolvedValue({ needsManualLocate: true, soughtExeName: "steam.exe" });
    native.openFolderPicker.mockResolvedValue("C:\\Steam\\steam.exe");
    native.ConfirmPlatformExePath.mockResolvedValue(undefined);
    expect(await openPlatformPage("Steam")).toBe(true);
    expect(native.ConfirmPlatformExePath).toHaveBeenCalledWith("Steam", "C:\\Steam\\steam.exe");
    expect(native.navigate).toHaveBeenCalledOnce();
  });

  it("does not navigate when confirming the path fails, and releases the busy state", async () => {
    native.ResolvePlatformLaunch.mockResolvedValue({ needsManualLocate: true });
    native.openFolderPicker.mockResolvedValue("C:\\Steam\\steam.exe");
    native.ConfirmPlatformExePath.mockRejectedValue(new Error("Path unavailable"));
    await expect(openPlatformPage("Steam")).rejects.toThrow("Path unavailable");
    expect(native.navigate).not.toHaveBeenCalled();
    expect(get(platformNavigationBusy)).toBe(false);
  });

  it("keeps shortcut discovery feedback", async () => {
    native.ResolvePlatformLaunch.mockResolvedValue({ ok: true, foundViaShortcut: true });
    await openPlatformPage("Steam");
    expect(native.pushToast).toHaveBeenCalledWith(expect.objectContaining({ message: "Toast_FoundExeViaShortcut" }));
  });

  it("prevents a second request while native resolution is pending", async () => {
    let finish!: (value: { ok: boolean }) => void;
    native.ResolvePlatformLaunch.mockReturnValue(new Promise((resolve) => { finish = resolve; }));
    const first = openPlatformPage("Steam");
    expect(await openPlatformPage("Discord")).toBe(false);
    finish({ ok: true });
    await first;
    expect(native.ResolvePlatformLaunch).toHaveBeenCalledOnce();
    expect(get(platformNavigationBusy)).toBe(false);
  });
});

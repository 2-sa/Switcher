import { describe, expect, it } from "vitest";
import { applyViewportDropdownLayout, computeViewportDropdownLayout, computeAboveDropdownLayout } from "./viewportDropdown";

const VIEWPORT = { viewportHeight: 600, viewportWidth: 1000 };

describe("shortcut dropdown bounds", () => {
  const bounds = { left: 0, right: 718, top: 32 };
  it("anchors Arabic menus at the trigger's right edge", () => {
    expect(computeAboveDropdownLayout({ left: 630, right: 666, top: 540 }, bounds, 200, true))
      .toEqual({ left: -164, maxWidth: 702, maxHeight: 260 });
  });
  it("keeps menus inside both content edges after footer wrapping", () => {
    expect(computeAboveDropdownLayout({ left: 20, right: 56, top: 540 }, bounds, 200, true).left).toBe(-12);
    expect(computeAboveDropdownLayout({ left: 650, right: 686, top: 540 }, bounds, 200, false).left).toBe(-140);
  });
  it("limits width and height in a compact window", () => {
    expect(computeAboveDropdownLayout({ left: 18, right: 54, top: 180 }, { left: 0, right: 160, top: 32 }, 200, false))
      .toEqual({ left: -10, maxWidth: 144, maxHeight: 132 });
  });
});

describe("viewport dropdown layout", () => {
  it("keeps a menu below when it fits there, and asks for no scroll", () => {
    expect(
      computeViewportDropdownLayout(
        { top: 80, bottom: 118, left: 20 },
        { ...VIEWPORT, menuHeight: 180, menuWidth: 200 },
      ),
    ).toEqual({ maxHeight: 180, shift: 0, scrollBy: 0 });
  });

  it("stays below a trigger near the bottom and asks to scroll it into room", () => {
    const layout = computeViewportDropdownLayout(
      { top: 540, bottom: 578, left: 20 },
      { ...VIEWPORT, menuHeight: 360, menuWidth: 200 },
    );

    expect(layout.scrollBy).toBe(346);
    // 14px of real space is unusable, so it takes the floor and the scroll fixes it.
    expect(layout.maxHeight).toBe(144);
  });

  it("caps a tall menu to a share of the viewport rather than the full list", () => {
    const layout = computeViewportDropdownLayout(
      { top: 80, bottom: 118, left: 20 },
      { ...VIEWPORT, menuHeight: 900, menuWidth: 200 },
    );

    expect(layout.maxHeight).toBe(360);
  });

  it("pulls a menu back from the right edge", () => {
    const layout = computeViewportDropdownLayout(
      { top: 80, bottom: 118, left: 900 },
      { ...VIEWPORT, menuHeight: 180, menuWidth: 220 },
    );

    expect(layout.shift).toBe(-128);
  });

  it("never pushes a menu off the left edge to save the right", () => {
    const layout = computeViewportDropdownLayout(
      { top: 80, bottom: 118, left: 12 },
      { ...VIEWPORT, menuHeight: 180, menuWidth: 1200 },
    );

    expect(layout.shift).toBe(-4);
  });

  it("keeps runtime placement authoritative over theme styles", () => {
    const declarations: Array<[string, string, string]> = [];
    const style = {
      setProperty(name: string, value: string, priority = "") {
        declarations.push([name, value, priority]);
      },
    };

    applyViewportDropdownLayout(style, { maxHeight: 272, shift: -40, scrollBy: 0 });

    expect(declarations).toEqual([
      ["top", "100%", "important"],
      ["bottom", "auto", "important"],
      ["left", "-40px", "important"],
      ["max-height", "272px", ""],
    ]);
  });
});

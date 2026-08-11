import { get } from "svelte/store";
import { beforeEach, describe, expect, it } from "vitest";
import { themeControl } from "./theme-control";

describe("themeControl.forceLight", () => {
  beforeEach(() => {
    themeControl.set.dark();
  });

  it("switches to light and restores dark", () => {
    const restore = themeControl.forceLight();

    expect(get(themeControl.current)).toBe("light");
    expect(document.documentElement.classList.contains("dark")).toBe(false);

    restore();

    expect(get(themeControl.current)).toBe("dark");
    expect(document.documentElement.classList.contains("dark")).toBe(true);
  });

  it("leaves the persisted preference alone", () => {
    const before = get(themeControl.preference);

    const restore = themeControl.forceLight();
    expect(get(themeControl.preference)).toBe(before);

    restore();
    expect(get(themeControl.preference)).toBe(before);
  });

  it("is a no-op when already light", () => {
    themeControl.set.light();

    const restore = themeControl.forceLight();
    expect(get(themeControl.current)).toBe("light");

    restore();
    expect(get(themeControl.current)).toBe("light");
    expect(document.documentElement.classList.contains("dark")).toBe(false);
  });
});

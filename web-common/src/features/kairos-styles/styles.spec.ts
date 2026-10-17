import { get } from "svelte/store";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { themeControl } from "../themes/theme-control";
import { styleControl } from "./style-control";
import { DEFAULT_STYLE, STYLE_ATTRIBUTE, resolveStyle } from "./styles";

describe("resolveStyle", () => {
  it.each([
    ["app: user wins over Org", { user: "aurora", org: "kairos-v1" }, "aurora"],
    [
      "app: Org when the user follows it",
      { user: null, org: "aurora" },
      "aurora",
    ],
    ["app: default with neither", {}, DEFAULT_STYLE],
    [
      "app: unknown user value falls to Org",
      { user: "retired", org: "aurora" },
      "aurora",
    ],
    ["app: unknown Org value falls to default", { org: "x" }, DEFAULT_STYLE],
    [
      "app: embed param ignored outside embeds",
      { embedParam: "aurora" },
      DEFAULT_STYLE,
    ],
    [
      "embed: param wins over Org",
      { isEmbed: true, embedParam: "aurora", org: "kairos-v1" },
      "aurora",
    ],
    [
      "embed: viewer's own choice is ignored",
      { isEmbed: true, user: "aurora", org: "kairos-v1" },
      "kairos-v1",
    ],
    ["embed: default with neither", { isEmbed: true }, DEFAULT_STYLE],
  ])("%s", (_, sources, expected) => {
    expect(resolveStyle(sources)).toBe(expected);
  });
});

describe("styleControl", () => {
  beforeEach(() => {
    localStorage.clear();
    styleControl.setOrg(null);
    styleControl.setPreference(null);
  });

  it("paints the resolved Style on <html>", () => {
    styleControl.setPreference("aurora");
    expect(document.documentElement.getAttribute(STYLE_ATTRIBUTE)).toBe(
      "aurora",
    );
  });

  it("persists the viewer's choice and clears it to follow the Org", () => {
    styleControl.setOrg("acme", "kairos-v1");
    styleControl.setPreference("aurora");
    expect(localStorage.getItem("kairos:style")).toBe("aurora");

    styleControl.setPreference(null);
    expect(localStorage.getItem("kairos:style")).toBeNull();
    expect(get(styleControl)).toBe("kairos-v1");
  });

  it("uses the cached Org default while the Org loads", () => {
    styleControl.setOrg("acme", "aurora");
    styleControl.setOrg(null);
    expect(get(styleControl)).toBe(DEFAULT_STYLE);

    styleControl.setOrg("acme");
    expect(get(styleControl)).toBe("aurora");
  });

  it("asks charts to re-render when the Style changes, not when it stays", () => {
    const refresh = vi.spyOn(themeControl, "refresh");
    styleControl.setPreference("aurora");
    expect(refresh).toHaveBeenCalledTimes(1);

    styleControl.setPreference("aurora");
    expect(refresh).toHaveBeenCalledTimes(1);
    refresh.mockRestore();
  });

  it("an embed paints its URL parameter, never the viewer's choice", () => {
    styleControl.setPreference("aurora");
    styleControl.setEmbed("kairos-v1");
    expect(get(styleControl)).toBe("kairos-v1");

    styleControl.setEmbed(null);
    expect(get(styleControl)).toBe(DEFAULT_STYLE);
  });
});

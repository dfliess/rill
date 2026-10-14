import type { V1Resource } from "@rilldata/web-common/runtime-client";
import { describe, expect, it } from "vitest";
import { homeV2MenuItems } from "./menu";

const canvas = (
  name: string,
  displayName: string,
  metadata?: Record<string, string>,
): V1Resource => ({
  meta: { name: { name }, metadata },
  canvas: { spec: { displayName }, state: { validSpec: {} } },
});

const titles = (resources: V1Resource[]) =>
  homeV2MenuItems(resources, "/o/p/-/home-v2").map((i) => i.title);

describe("homeV2MenuItems", () => {
  it("puts ordered canvases first, then the rest alphabetically", () => {
    expect(
      titles([
        canvas("resultados", "Resultados"),
        canvas("caja", "Caja y Liquidez", { kairos_order: "2" }),
        canvas("vision", "Visión general", { kairos_order: "1" }),
        canvas("anual", "Anual"),
      ]),
    ).toEqual(["Visión general", "Caja y Liquidez", "Anual", "Resultados"]);
  });

  it("leaves hidden canvases out", () => {
    expect(
      titles([
        canvas("caja", "Caja y Liquidez"),
        canvas("borrador", "Borrador", { kairos_hidden: "true" }),
      ]),
    ).toEqual(["Caja y Liquidez"]);
  });

  it("serves the first canvas at the bare URL and the rest by name", () => {
    const items = homeV2MenuItems(
      [
        canvas("caja_y_liquidez", "Caja y Liquidez"),
        canvas("home", "Inicio", { kairos_order: "1" }),
      ],
      "/o/p/-/home-v2",
    );
    expect(items.map((i) => [i.name, i.href])).toEqual([
      ["home", "/o/p/-/home-v2"],
      ["caja_y_liquidez", "/o/p/-/home-v2/caja_y_liquidez"],
    ]);
  });

  it("falls back to the file name and skips invalid canvases", () => {
    expect(
      titles([
        canvas("sin_titulo", ""),
        { meta: { name: { name: "roto" } }, canvas: { state: {} } },
        { meta: { name: { name: "ventas" } }, explore: {} },
      ]),
    ).toEqual(["sin_titulo"]);
  });
});

import type { V1Resource } from "@rilldata/web-common/runtime-client";
import { describe, expect, it } from "vitest";
import { homeMenuItems, homePersonalMenuItems } from "./menu";

const canvas = (
  name: string,
  displayName: string,
  metadata?: Record<string, string>,
  annotations?: Record<string, string>,
): V1Resource => ({
  meta: { name: { name }, metadata },
  canvas: { spec: { displayName }, state: { validSpec: { annotations } } },
});

const titles = (resources: V1Resource[]) =>
  homeMenuItems(resources, "/o/p").map((i) => i.title);

describe("homeMenuItems", () => {
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
    const items = homeMenuItems(
      [
        canvas("caja_y_liquidez", "Caja y Liquidez"),
        canvas("home", "Inicio", { kairos_order: "1" }),
      ],
      "/o/p",
    );
    expect(items.map((i) => [i.name, i.href])).toEqual([
      ["home", "/o/p"],
      ["caja_y_liquidez", "/o/p/-/home/caja_y_liquidez"],
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

  it("carries the canvas's icon from its metadata", () => {
    const [caja, ventas] = homeMenuItems(
      [
        canvas("caja", "Caja", { kairos_order: "1", kairos_icon: "wallet" }),
        canvas("ventas", "Ventas", { kairos_order: "2" }),
      ],
      "/o/p",
    );
    expect(caja.icon).toBe("wallet");
    expect(ventas.icon).toBeUndefined();
  });

  it("lists a personal canvas only once its owner shares it, like the Dashboards tab", () => {
    expect(
      titles([
        canvas("caja", "Caja"),
        canvas("mio", "Mío", undefined, { admin_managed: "true" }),
        canvas("compartido", "Compartido", undefined, {
          admin_managed: "true",
          admin_shared: "true",
        }),
      ]),
    ).toEqual(["Caja", "Compartido"]);
  });

  it("leaves the viewer's own personal canvases to their group", () => {
    const shared = canvas("compartido", "Compartido", undefined, {
      admin_managed: "true",
      admin_shared: "true",
    });
    expect(
      homeMenuItems(
        [canvas("caja", "Caja"), shared],
        "/o/p",
        new Set(["compartido"]),
      ).map((i) => i.name),
    ).toEqual(["caja"]);
  });
});

describe("homePersonalMenuItems", () => {
  it("lists them alphabetically, opening the personal canvas page", () => {
    const items = homePersonalMenuItems(
      [canvas("zeta", "Zeta"), canvas("mi_vista", "Alfa")],
      "/o/p",
    );
    expect(items.map((i) => [i.title, i.href])).toEqual([
      ["Alfa", "/o/p/-/personal/mi_vista"],
      ["Zeta", "/o/p/-/personal/zeta"],
    ]);
  });
});

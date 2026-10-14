import type { RuntimeClient } from "@rilldata/web-common/runtime-client/v2";
import { describe, expect, it, vi } from "vitest";
import { isInitialBuild, useDashboards } from "./selectors";

// Hands back the query options, so the listing's filter can be run on its own.
vi.mock("@rilldata/web-common/runtime-client", async (importOriginal) => ({
  ...(await importOriginal<object>()),
  createRuntimeServiceListResources: (
    _client: unknown,
    _req: unknown,
    options: unknown,
  ) => options,
}));

const dashboard = { explore: {} };
const model = { model: {} };

describe("isInitialBuild", () => {
  it("is false once a dashboard exists, even while the runtime is still initializing", () => {
    expect(
      isInitialBuild({ resources: [dashboard, model], initializing: true }),
    ).toBe(false);
  });

  it("is true when no dashboards exist yet and the runtime is still initializing", () => {
    expect(isInitialBuild({ resources: [model], initializing: true })).toBe(
      true,
    );
  });

  it("is false when security policies denied every resource", () => {
    // ListResources returns 200 with an empty list. Reading that as "still building"
    // left embed users on a deny-by-default project staring at a permanent spinner.
    expect(isInitialBuild({ resources: [], initializing: false })).toBe(false);
  });

  it("is false when the project genuinely has no dashboards", () => {
    expect(isInitialBuild({ resources: [model], initializing: false })).toBe(
      false,
    );
  });
});

describe("useDashboards", () => {
  const listed = (resources: object[]) => {
    const options = useDashboards({
      instanceId: "i",
    } as RuntimeClient) as unknown as {
      query: { select: (d: object) => { meta: { name: { name: string } } }[] };
    };
    return options.query.select({ resources }).map((r) => r.meta.name.name);
  };
  const canvas = (name: string, annotations?: Record<string, string>) => ({
    meta: { name: { name } },
    canvas: { state: { validSpec: { annotations } } },
  });

  // Kairos #231: the Home canvases used to be hidden here (`kairos_home`);
  // they are ordinary canvases now.
  it("lists canvases annotated kairos_home like any other", () => {
    expect(listed([canvas("home", { kairos_home: "true" })])).toEqual(["home"]);
  });

  it("leaves unshared personal canvases out", () => {
    expect(
      listed([
        canvas("mio", { admin_managed: "true" }),
        canvas("compartido", { admin_managed: "true", admin_shared: "true" }),
      ]),
    ).toEqual(["compartido"]);
  });
});

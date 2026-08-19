import { describe, expect, it } from "vitest";
import {
  derivePushSectionState,
  describeUserAgent,
  listPushDevices,
  readNotificationPreferences,
  subscriptionCredentials,
  urlBase64ToUint8Array,
} from "./push-utils";

describe("urlBase64ToUint8Array", () => {
  it("decodes unpadded base64", () => {
    expect(Array.from(urlBase64ToUint8Array("AQID"))).toEqual([1, 2, 3]);
  });

  it("restores the padding the base64url alphabet drops", () => {
    // "SGVsbG8" is "Hello" with its "=" stripped.
    expect(Array.from(urlBase64ToUint8Array("SGVsbG8"))).toEqual([
      72, 101, 108, 108, 111,
    ]);
  });

  it("maps the url-safe alphabet back to standard base64", () => {
    // "_w" is base64url for 0xff ("/w==" in standard base64), "-_8" for 0xfb 0xff.
    expect(Array.from(urlBase64ToUint8Array("_w"))).toEqual([255]);
    expect(Array.from(urlBase64ToUint8Array("-_8"))).toEqual([251, 255]);
  });

  it("round-trips a VAPID-sized key", () => {
    // Uncompressed P-256 points are 65 bytes, the size of a real VAPID public key.
    const bytes = Array.from({ length: 65 }, (_, i) => (i * 7) % 256);
    const encoded = Buffer.from(bytes).toString("base64url");
    expect(Array.from(urlBase64ToUint8Array(encoded))).toEqual(bytes);
  });
});

describe("derivePushSectionState", () => {
  const ready = {
    configLoaded: true,
    vapidPublicKey: "key",
    supported: true,
    permission: "granted" as NotificationPermission,
  };

  it("is loading until the server config arrives", () => {
    expect(derivePushSectionState({ ...ready, configLoaded: false })).toBe(
      "loading",
    );
  });

  it("is disabled when the deployment has no VAPID key, even without browser support", () => {
    expect(
      derivePushSectionState({
        ...ready,
        vapidPublicKey: "",
        supported: false,
      }),
    ).toBe("disabled");
  });

  it("is unsupported when the browser lacks the Push API", () => {
    expect(derivePushSectionState({ ...ready, supported: false })).toBe(
      "unsupported",
    );
  });

  it("is denied when the user blocked notifications", () => {
    expect(derivePushSectionState({ ...ready, permission: "denied" })).toBe(
      "denied",
    );
  });

  it("is ready otherwise, whether or not permission was granted yet", () => {
    expect(derivePushSectionState(ready)).toBe("ready");
    expect(derivePushSectionState({ ...ready, permission: "default" })).toBe(
      "ready",
    );
  });
});

describe("describeUserAgent", () => {
  it("recognizes Chrome on macOS", () => {
    expect(
      describeUserAgent(
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
      ),
    ).toBe("Chrome · macOS");
  });

  it("recognizes Edge despite its UA also containing Chrome and Safari", () => {
    expect(
      describeUserAgent(
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0",
      ),
    ).toBe("Edge · Windows");
  });

  it("recognizes Firefox on Linux", () => {
    expect(
      describeUserAgent(
        "Mozilla/5.0 (X11; Linux x86_64; rv:127.0) Gecko/20100101 Firefox/127.0",
      ),
    ).toBe("Firefox · Linux");
  });

  it("reports Android rather than Linux for Android Chrome", () => {
    expect(
      describeUserAgent(
        "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36",
      ),
    ).toBe("Chrome · Android");
  });

  it("recognizes Safari on iOS", () => {
    expect(
      describeUserAgent(
        "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
      ),
    ).toBe("Safari · iOS");
  });

  it("falls back to a truncated raw UA when nothing matches", () => {
    const exotic = "SomeBot/1.0 (unknown platform)".repeat(10);
    expect(describeUserAgent(exotic)).toBe(exotic.slice(0, 80));
  });
});

describe("subscriptionCredentials", () => {
  const full = {
    endpoint: "https://push.example.com/sub/abc",
    keys: { p256dh: "p256dh-key", auth: "auth-secret" },
  };

  it("extracts endpoint and keys", () => {
    expect(subscriptionCredentials({ toJSON: () => full })).toEqual({
      endpoint: "https://push.example.com/sub/abc",
      p256dh: "p256dh-key",
      auth: "auth-secret",
    });
  });

  it("returns null when the endpoint or either key is missing", () => {
    expect(
      subscriptionCredentials({ toJSON: () => ({ ...full, endpoint: "" }) }),
    ).toBeNull();
    expect(
      subscriptionCredentials({
        toJSON: () => ({ ...full, keys: { p256dh: "p256dh-key" } }),
      }),
    ).toBeNull();
    expect(subscriptionCredentials({ toJSON: () => ({}) })).toBeNull();
  });
});

describe("readNotificationPreferences", () => {
  it("reads a flag the response omits as off, since proto3 JSON drops false", () => {
    expect(readNotificationPreferences({ pushAlerts: true })).toEqual({
      pushAlerts: true,
      pushReports: false,
      pushActApprovals: false,
    });
  });

  it("reads a missing preferences object as all off", () => {
    expect(readNotificationPreferences(undefined)).toEqual({
      pushAlerts: false,
      pushReports: false,
      pushActApprovals: false,
    });
  });
});

describe("listPushDevices", () => {
  const subscriptions = [
    {
      id: "old",
      endpoint: "https://push.example.com/a",
      userAgent: "Firefox · Linux",
      createdOn: "2026-08-01T10:00:00Z",
    },
    {
      id: "new",
      endpoint: "https://push.example.com/b",
      userAgent: "Chrome · macOS",
      createdOn: "2026-08-17T10:00:00Z",
    },
    {
      id: "here",
      endpoint: "https://push.example.com/c",
      userAgent: "Safari · iOS",
      createdOn: "2026-07-01T10:00:00Z",
    },
  ];

  it("puts this device first and flags it, whatever its age", () => {
    const devices = listPushDevices(
      subscriptions,
      "https://push.example.com/c",
    );
    expect(devices.map((d) => d.id)).toEqual(["here", "new", "old"]);
    expect(devices.map((d) => d.isCurrent)).toEqual([true, false, false]);
  });

  it("sorts by recency and flags nothing when the browser has no subscription", () => {
    const devices = listPushDevices(subscriptions, undefined);
    expect(devices.map((d) => d.id)).toEqual(["new", "old", "here"]);
    expect(devices.some((d) => d.isCurrent)).toBe(false);
  });

  it("handles an empty or missing list", () => {
    expect(listPushDevices(undefined, "https://push.example.com/c")).toEqual(
      [],
    );
    expect(listPushDevices([], undefined)).toEqual([]);
  });

  it("drops subscriptions the server did not identify and defaults a missing name to empty", () => {
    const devices = listPushDevices(
      [{ endpoint: "https://push.example.com/d" }, { id: "x", createdOn: "z" }],
      undefined,
    );
    expect(devices).toEqual([
      { id: "x", label: "", createdOn: "z", isCurrent: false },
    ]);
  });
});

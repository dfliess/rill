import { describe, expect, it, vi } from "vitest";
import {
  ensurePushSubscription,
  usesApplicationServerKey,
  type PushManagerLike,
  type PushSubscriptionLike,
} from "./push-subscription";

// "AQID" is base64url for the bytes 1, 2, 3.
const VAPID_KEY = "AQID";
const VAPID_BYTES = new Uint8Array([1, 2, 3]);

function subscriptionStub(
  endpoint: string,
  applicationServerKey: ArrayBuffer | null,
): PushSubscriptionLike & { unsubscribe: ReturnType<typeof vi.fn> } {
  return {
    endpoint,
    options: { applicationServerKey },
    toJSON: () => ({ endpoint }),
    unsubscribe: vi.fn().mockResolvedValue(true),
  };
}

function managerStub(existing: PushSubscriptionLike | null): PushManagerLike & {
  subscribe: ReturnType<typeof vi.fn>;
} {
  return {
    getSubscription: () => Promise.resolve(existing),
    subscribe: vi
      .fn()
      .mockResolvedValue(
        subscriptionStub("https://push.example.com/new", null),
      ),
  };
}

describe("ensurePushSubscription", () => {
  it("subscribes with the deployment's key when the browser has none", async () => {
    const manager = managerStub(null);

    const subscription = await ensurePushSubscription(manager, VAPID_KEY);

    expect(subscription.endpoint).toBe("https://push.example.com/new");
    expect(manager.subscribe).toHaveBeenCalledWith({
      userVisibleOnly: true,
      applicationServerKey: VAPID_BYTES,
    });
  });

  it("reuses a subscription created with the same key", async () => {
    const existing = subscriptionStub(
      "https://push.example.com/old",
      VAPID_BYTES.buffer,
    );
    const manager = managerStub(existing);

    const subscription = await ensurePushSubscription(manager, VAPID_KEY);

    expect(subscription).toBe(existing);
    expect(manager.subscribe).not.toHaveBeenCalled();
    expect(existing.unsubscribe).not.toHaveBeenCalled();
  });

  it("drops a subscription from a rotated key before subscribing again", async () => {
    const existing = subscriptionStub(
      "https://push.example.com/stale",
      new Uint8Array([9, 9, 9]).buffer,
    );
    const manager = managerStub(existing);

    const subscription = await ensurePushSubscription(manager, VAPID_KEY);

    expect(existing.unsubscribe).toHaveBeenCalled();
    expect(manager.subscribe).toHaveBeenCalledWith({
      userVisibleOnly: true,
      applicationServerKey: VAPID_BYTES,
    });
    expect(subscription.endpoint).toBe("https://push.example.com/new");
  });
});

describe("usesApplicationServerKey", () => {
  it("compares the bytes, not the buffer identity", () => {
    const subscription = subscriptionStub(
      "https://push.example.com/a",
      new Uint8Array([1, 2, 3]).buffer,
    );
    expect(usesApplicationServerKey(subscription, VAPID_BYTES)).toBe(true);
  });

  it("rejects a key of the same length with different bytes", () => {
    const subscription = subscriptionStub(
      "https://push.example.com/a",
      new Uint8Array([1, 2, 4]).buffer,
    );
    expect(usesApplicationServerKey(subscription, VAPID_BYTES)).toBe(false);
  });

  it("rejects a prefix of the key", () => {
    const subscription = subscriptionStub(
      "https://push.example.com/a",
      new Uint8Array([1, 2]).buffer,
    );
    expect(usesApplicationServerKey(subscription, VAPID_BYTES)).toBe(false);
  });

  it("accepts a subscription from a browser that hides its key", () => {
    const subscription = subscriptionStub("https://push.example.com/a", null);
    expect(usesApplicationServerKey(subscription, VAPID_BYTES)).toBe(true);
    expect(
      usesApplicationServerKey(
        { ...subscription, options: undefined },
        VAPID_BYTES,
      ),
    ).toBe(true);
  });
});

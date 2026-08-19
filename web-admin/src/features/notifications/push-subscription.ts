// The half of the push settings that talks to the browser's Push API.
// `ensurePushSubscription` takes the push manager as an argument rather than
// reaching for `navigator` itself, so the branch that actually bites — a
// subscription left over from a previous VAPID key — is unit-testable.

import { urlBase64ToUint8Array } from "./push-utils";

// Structural subsets of the DOM types: a real PushSubscription satisfies them,
// and so can a stub in a test.
export type PushSubscriptionLike = {
  readonly endpoint: string;
  readonly options?: { applicationServerKey?: ArrayBuffer | null } | null;
  toJSON(): PushSubscriptionJSON;
  unsubscribe(): Promise<boolean>;
};

export type PushManagerLike = {
  getSubscription(): Promise<PushSubscriptionLike | null>;
  subscribe(options: {
    userVisibleOnly: boolean;
    applicationServerKey: Uint8Array<ArrayBuffer>;
  }): Promise<PushSubscriptionLike>;
};

// `navigator.serviceWorker.ready` never rejects: it waits forever when no
// worker ever activates, which happens on a build without the worker or after
// a failed registration. The timeout turns that into an error the UI can
// report instead of a switch that spins for good.
const WORKER_READY_TIMEOUT_MS = 10_000;

export async function pushManagerOfActiveWorker(): Promise<PushManagerLike> {
  const ready = navigator.serviceWorker.ready;
  const timeout = new Promise<never>((_, reject) =>
    setTimeout(
      () => reject(new Error("service worker did not become ready")),
      WORKER_READY_TIMEOUT_MS,
    ),
  );
  const registration = await Promise.race([ready, timeout]);
  return registration.pushManager;
}

// Returns the subscription this browser holds for the deployment's current
// VAPID key, subscribing if there is none. A subscription left over from a
// previous key is dropped first: `subscribe` refuses to replace it and throws
// InvalidStateError instead, which would leave the device stuck on a key the
// server no longer signs with.
export async function ensurePushSubscription(
  manager: PushManagerLike,
  vapidPublicKey: string,
): Promise<PushSubscriptionLike> {
  const applicationServerKey = urlBase64ToUint8Array(vapidPublicKey);

  const existing = await manager.getSubscription();
  if (existing) {
    if (usesApplicationServerKey(existing, applicationServerKey)) {
      return existing;
    }
    await existing.unsubscribe();
  }

  return manager.subscribe({ userVisibleOnly: true, applicationServerKey });
}

export function usesApplicationServerKey(
  subscription: PushSubscriptionLike,
  key: Uint8Array,
): boolean {
  const current = subscription.options?.applicationServerKey;
  // Browsers that don't expose the key can't be checked, so the existing
  // subscription is taken at face value rather than churned on every visit.
  if (!current) return true;
  const bytes = new Uint8Array(current);
  return bytes.length === key.length && bytes.every((b, i) => b === key[i]);
}

// Reads the endpoint this browser is already subscribed with, if any. Uses
// `getRegistration` rather than `ready` so a page without an active worker
// answers "none" instead of waiting.
export async function currentPushEndpoint(): Promise<string | undefined> {
  const registration = await navigator.serviceWorker.getRegistration();
  const subscription = await registration?.pushManager.getSubscription();
  return subscription?.endpoint;
}

export async function unsubscribeThisDevice(): Promise<void> {
  const registration = await navigator.serviceWorker.getRegistration();
  const subscription = await registration?.pushManager.getSubscription();
  await subscription?.unsubscribe();
}

// The badge on the installed app's icon: how many things are waiting for the person looking at it
// (kairos-cloud#147). Two very different contexts write it, so everything here has to work in both: the service
// worker, when a notification arrives with the app closed, and the page, when it is open and can ask the server.
//
// The badge is one number for the whole origin, while work is pending per project. So neither writer knows the
// total on its own: each records what it learned about ONE project in a small ledger, and the badge is the sum.
// A push updates the project it came from; the open app updates the project it is on. What neither has heard from
// recently keeps its last known value, which is the honest answer: not zero, and not a number we made up.
//
// Nothing here throws. A browser without the Badging API (Firefox; Safari outside an installed app) or without a
// usable IndexedDB (private windows, storage blocked) is a browser where this feature does not exist, not one
// where the app should report errors.

const DB_NAME = "kairos-app-badge";
const DB_VERSION = 1;
// One record per project: the key is its scope, the value the count last known for it.
const STORE = "pending";

// PendingUpdate is what a push payload carries about the badge: one project's count, and which project it is.
export type PendingUpdate = {
  // "organization/project", lowercased. Composed by the admin service for a push; by scopeOf for the open app.
  scope: string;
  count: number;
};

// scopeOf keys the ledger the same way the admin service does, so a push and the open app write the same entry
// for the same project instead of each keeping its own.
export function scopeOf(organization: string, project: string): string {
  return `${organization}/${project}`.toLowerCase();
}

// readPendingUpdate extracts the badge instruction from a push payload, or nothing if it carries none.
// A message about something there is nothing to resolve (a report, an alert) leaves the badge alone, and so does a
// malformed one: an app icon showing a number nobody can explain is worse than one showing yesterday's.
export function readPendingUpdate(raw: unknown): PendingUpdate | null {
  if (typeof raw !== "object" || raw === null) return null;
  const { scope, count } = raw as Record<string, unknown>;
  if (typeof scope !== "string" || !scope) return null;
  // Proto3 JSON omits a zero, so a payload with a scope and no count means zero, not "unset".
  const parsed = count === undefined ? 0 : count;
  if (typeof parsed !== "number" || !Number.isInteger(parsed) || parsed < 0) {
    return null;
  }
  return { scope, count: parsed };
}

// countPendingRuns is the number this app puts in the ledger for the project it is looking at: how many runs are
// stopped waiting for a decision this viewer may make.
//
// Runs, not approvals, and deliberately: a run that proposed three actions in one turn is one pause, one
// notification and one row in the inbox, so counting its approvals would put a 3 on the icon for one thing to look
// at. It mirrors pendingRunsByRecipient in runtime/server/act_approval_push.go, which is what a push carries. The
// two have to agree, or the badge would jump every time the app opens.
export function countPendingRuns(
  approvals: { runId?: string; canDecide?: boolean }[] | undefined,
): number {
  const runs = new Set<string>();
  for (const approval of approvals ?? []) {
    if (approval.canDecide && approval.runId) runs.add(approval.runId);
  }
  return runs.size;
}

// ledgerTotal is what the icon shows once one project has reported in: its new count, plus the last thing every
// other project said about itself.
//
// A project that has not been heard from keeps its number rather than dropping out. That is the whole point of
// keeping a ledger instead of a single figure: a push about one project says nothing about the others, and neither
// does the app being open on one of them. Zero is a report like any other, so the badge clears only once every
// project has said it has nothing left.
export function ledgerTotal(
  stored: Record<string, number>,
  update: PendingUpdate,
): number {
  let total = 0;
  for (const [scope, count] of Object.entries({
    ...stored,
    [update.scope]: update.count,
  })) {
    // A value that is not a count is a ledger entry from a version that wrote something else, or storage that came
    // back corrupted. Skipping it beats adding NaN to the icon.
    if (scope && Number.isInteger(count) && count > 0) total += count;
  }
  return total;
}

// recordPending stores one project's count and returns the new total across all of them.
// The read and the write share one transaction, so the page and the service worker reporting on different projects
// at the same moment cannot each compute a total from a ledger the other has already moved on from.
async function recordPending(update: PendingUpdate): Promise<number | null> {
  const db = await openLedger();
  if (!db) return null;

  try {
    return await new Promise<number | null>((resolve) => {
      const tx = db.transaction(STORE, "readwrite");
      const store = tx.objectStore(STORE);
      // Requests complete in the order they were issued, so the keys are in hand by the time the values arrive.
      const keys = store.getAllKeys();
      const values = store.getAll();

      let total: number | null = null;
      values.onsuccess = () => {
        const stored: Record<string, number> = {};
        (keys.result ?? []).forEach((key, i) => {
          // Scopes are strings. Anything else is a key this never wrote, and there is no project it names.
          if (typeof key !== "string") return;
          stored[key] = (values.result ?? [])[i] as number;
        });
        total = ledgerTotal(stored, update);
        store.put(update.count, update.scope);
      };

      tx.oncomplete = () => resolve(total);
      tx.onerror = () => resolve(null);
      tx.onabort = () => resolve(null);
    });
  } finally {
    db.close();
  }
}

function openLedger(): Promise<IDBDatabase | null> {
  if (typeof indexedDB === "undefined") return Promise.resolve(null);

  return new Promise((resolve) => {
    let request: IDBOpenDBRequest;
    try {
      request = indexedDB.open(DB_NAME, DB_VERSION);
    } catch {
      // Some browsers throw here rather than failing the request (storage blocked outright).
      resolve(null);
      return;
    }
    request.onupgradeneeded = () => {
      if (!request.result.objectStoreNames.contains(STORE)) {
        request.result.createObjectStore(STORE);
      }
    };
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => resolve(null);
    request.onblocked = () => resolve(null);
  });
}

// Badging is the part of the Badging API this uses. Both methods are optional because in most browsers they are
// simply not there: `navigator` and a service worker's `WorkerNavigator` both satisfy this.
export type Badging = {
  setAppBadge?: (contents?: number) => Promise<void>;
  clearAppBadge?: () => Promise<void>;
};

// applyBadge paints a total on the app icon, or takes the badge off when nothing is waiting.
// Absent support is the common case, not a failure: Firefox has no Badging API at all, and Safari exposes it only
// to an app installed on the Home Screen or in the Dock. Even where it exists the call can be refused, so the
// missing method and the rejected promise are both swallowed: this must never reach the console.
export async function applyBadge(
  badging: Badging,
  total: number,
): Promise<void> {
  try {
    if (total > 0) {
      await badging.setAppBadge?.(total);
    } else {
      await badging.clearAppBadge?.();
    }
  } catch {
    // Nothing to recover: the icon simply keeps whatever it was showing.
  }
}

// syncBadge records what one project has waiting and repaints the icon with the new total.
// It is the whole feature in one call, and the only one either writer needs.
export async function syncBadge(update: PendingUpdate): Promise<void> {
  const total = await recordPending(update);
  if (total === null) return; // No ledger, so no total worth painting: leave the icon alone.
  await applyBadge(navigator, total);
}

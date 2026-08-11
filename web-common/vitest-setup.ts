import "@testing-library/jest-dom";
import { vi } from "vitest";
import { Settings } from "luxon";

// Neither jsdom nor Node exposes Web Storage here, so a module that reads it at
// import time (theme-control, and anything importing it) throws during
// collection. Defined on globalThis so the bare and window-qualified references
// resolve to the same store.
function inMemoryStorage(): Storage {
  const entries = new Map<string, string>();
  return {
    get length() {
      return entries.size;
    },
    key: (index: number) => [...entries.keys()][index] ?? null,
    getItem: (key: string) => entries.get(key) ?? null,
    setItem: (key: string, value: string) =>
      void entries.set(key, String(value)),
    removeItem: (key: string) => void entries.delete(key),
    clear: () => entries.clear(),
  } as Storage;
}

for (const name of ["localStorage", "sessionStorage"] as const) {
  Object.defineProperty(globalThis, name, {
    writable: true,
    configurable: true,
    value: inMemoryStorage(),
  });
}

// required for svelte5 + jsdom as jsdom does not support matchMedia
Object.defineProperty(window, "matchMedia", {
  writable: true,
  enumerable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

Object.defineProperty(window, "scrollTo", {
  writable: true,
  enumerable: true,
  value: vi.fn(),
});

Settings.defaultWeekSettings = {
  minimalDays: 4,
  firstDay: 1,
  weekend: [6, 7],
};

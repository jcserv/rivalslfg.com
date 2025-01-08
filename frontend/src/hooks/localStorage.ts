/* eslint-disable no-console */
import { useEffect, useState } from "react";

import { Profile } from "@/types";

// Type-safe storage keys and their corresponding value types
export type StorageKeyValueMap = {
  profile: Profile;
};

export type StorageKey = keyof StorageKeyValueMap;

type StorageEntry<T> = {
  value: T;
  expiration?: Date;
};

// Generic helper to get the value type for a given key
type StorageValue<K extends StorageKey> = StorageKeyValueMap[K];

function isEntryExpired(entry: StorageEntry<unknown> | null): boolean {
  if (!entry?.expiration) return false;
  return new Date(entry.expiration) < new Date();
}

export function getStorageValue<K extends StorageKey>(
  key: K,
  initialValue: StorageValue<K>,
): StorageEntry<StorageValue<K>> {
  try {
    const savedValue = localStorage.getItem(key);
    if (!savedValue) return { value: initialValue };

    const entry = JSON.parse(savedValue) as StorageEntry<StorageValue<K>>;
    if (!entry?.value || isEntryExpired(entry)) {
      localStorage.removeItem(key);
      return { value: initialValue };
    }
    return entry;
  } catch {
    return { value: initialValue };
  }
}

export function setStorageValue<K extends StorageKey>(
  key: K,
  value: StorageValue<K>,
  expiration?: Date,
): void {
  try {
    const entry: StorageEntry<StorageValue<K>> = {
      value: value || ({} as StorageValue<K>),
      expiration,
    };
    localStorage.setItem(key, JSON.stringify(entry));
  } catch (error) {
    console.error(`Error setting localStorage key "${key}":`, error);
  }
}

export function useLocalStorage<K extends StorageKey>(
  key: K,
  initialValue: StorageValue<K>,
  expiration?: Date,
) {
  const [entry, setEntry] = useState<StorageEntry<StorageValue<K>>>(() =>
    getStorageValue(key, initialValue),
  );

  useEffect(() => {
    try {
      const newEntry: StorageEntry<StorageValue<K>> = {
        value: entry.value || ({} as StorageValue<K>),
        expiration: entry.expiration || expiration,
      };
      localStorage.setItem(key, JSON.stringify(newEntry));
    } catch (error) {
      console.error(`Error setting localStorage key "${key}":`, error);
    }
  }, [key, entry, expiration]);

  const setValue = (newValue: StorageValue<K>) => {
    setEntry({ value: newValue || ({} as StorageValue<K>), expiration });
  };

  // Ensure we never return undefined/null
  return [entry.value || initialValue, setValue] as const;
}

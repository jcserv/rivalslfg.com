/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useState } from "react";

type StorageValue = {
  value: any;
  expiration: Date;
};

function isEntryExpired(s: StorageValue | null): boolean {
  if (!s) return false;
  return s.expiration && new Date(s.expiration) < new Date();
}

function getStorageValue(key: string, initialValue: any) {
  const savedValue = localStorage.getItem(key);
  if (!savedValue) return initialValue;

  const entry = JSON.parse(savedValue) as StorageValue;
  if (isEntryExpired(entry)) {
    localStorage.removeItem(key);
    return initialValue;
  }
  return entry;
}

export const useLocalStorage = (
  key: string,
  initialValue: any,
  expiration?: Date,
) => {
  const [entry, setEntry] = useState(() => getStorageValue(key, initialValue));
  useEffect(() => {
    if (entry === null) {
      localStorage.setItem(
        key,
        JSON.stringify({ value: initialValue, expiration }),
      );
    } else if (entry.expiration) {
      localStorage.setItem(
        key,
        JSON.stringify({ value: entry.value, expiration }),
      );
    } else {
      localStorage.setItem(key, JSON.stringify(entry));
    }
  }, [key, entry]);

  return [entry, setEntry] as const;
};

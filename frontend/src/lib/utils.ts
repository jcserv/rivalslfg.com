import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function toTitleCase(str: string) {
  return str
    .toLowerCase()
    .split(" ")
    .map((word: string) => {
      return word.charAt(0).toUpperCase() + word.slice(1);
    })
    .join(" ");
}

export function strArrayToTitleCase(strArray: string[]) {
  return strArray.map((str) => toTitleCase(str)).join(", ");
}

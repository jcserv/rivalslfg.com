import { Profanity } from "@2toad/profanity";

const allowList = (import.meta.env.VITE_PROFANITY_ALLOWLIST ?? "")
  .split(",")
  .map((s: string) => s.trim());

const profanity = new Profanity({
  languages: ["ar", "de", "en", "es", "fr", "hi", "ja", "ko", "pt", "ru", "zh"],
});

profanity.whitelist.addWords(allowList);

export const containsProfanity = (text: string): boolean => {
  return profanity.exists(text);
};

export const filterProfanity = (text: string): string => {
  return profanity.censor(text);
};

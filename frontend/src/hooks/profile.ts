import { useLocalStorage } from "@/hooks";
import { FOURTEEN_DAYS_FROM_TODAY, Profile } from "@/types";

export const useProfile = (): readonly [
  profile: Profile,
  setProfile: (p: Profile) => void,
] => {
  return useLocalStorage("profile", {}, FOURTEEN_DAYS_FROM_TODAY);
};

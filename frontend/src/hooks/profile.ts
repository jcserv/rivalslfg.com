import { useLocalStorage } from "@/hooks";
import { FOURTEEN_DAYS_FROM_TODAY, Profile } from "@/types";

export const useProfile = (): readonly [
  profile: Profile,
  setProfile: (p: Profile) => void,
  isProfileConfigured: boolean,
] => {
  const [profile, setProfile] = useLocalStorage(
    "profile",
    {},
    FOURTEEN_DAYS_FROM_TODAY,
  );
  const isProfileConfigured = Object.keys(profile).length > 0;

  return [profile, setProfile, isProfileConfigured];
};

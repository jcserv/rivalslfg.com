import { useLocalStorage } from "@/hooks";
import { FOURTEEN_DAYS_FROM_TODAY, Profile } from "@/types";

export const useProfile = (): readonly [
  profile: Profile,
  setProfile: (p: Profile) => void,
  isProfileConfigured: boolean,
  setProfileId: (id: number) => void,
] => {
  const [profile, setProfile] = useLocalStorage(
    "profile",
    {} as Profile,
    FOURTEEN_DAYS_FROM_TODAY,
  );
  const isProfileConfigured = Object.keys(profile || {}).length > 0;

  const setProfileId = (id: number) => {
    setProfile({
      ...profile,
      id,
    });
  };

  return [profile, setProfile, isProfileConfigured, setProfileId];
};

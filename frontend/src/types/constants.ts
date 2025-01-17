export const TEAM_SIZE = 6;

export const FOURTEEN_DAYS_FROM_TODAY = new Date(
  new Date().setDate(new Date().getDate() + 14),
);

export const ONE_MINUTE_IN_MILLISECONDS = 60000;

export type ProfileFormType = "find" | "create" | "profile";

export function getSubmitButtonLabel(p: ProfileFormType): string {
  switch (p) {
    case "find":
      return "Find Group";
    case "create":
      return "Create Group";
    default:
      return "Save";
  }
}

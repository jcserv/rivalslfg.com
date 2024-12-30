import { Region } from "@/types";

export const regions = Object.entries(Region).map(([key, value]) => ({
  value: key,
  label: value,
}));

export const gamemodes = [
  {
    value: "competitive",
    label: "Competitive",
    icon: "👑",
  },
  {
    value: "quickplay",
    label: "Quickplay",
    icon: "⚡",
  },
];

export const areRequirementsMet = [
  {
    value: "true",
    label: "Yes",
    icon: "✅",
  },
  {
    value: "false",
    label: "No",
    icon: "❌",
  },
];

import { Platform, Region } from "@/types";

export const regions = Object.entries(Region).map(([key, value]) => ({
  value: key,
  label: value,
}));

export const open = [
  {
    value: "true",
    label: "Public",
    icon: "👀",
  },
  {
    value: "false",
    label: "Private",
    icon: "🔒",
  },
];

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

export const platforms = [
  {
    value: Platform.PC,
    label: "PC",
    icon: "🖥️",
  },
  {
    value: Platform.Console,
    label: "Console",
    icon: "🎮",
  },
];

export const areRequirementsMet = [
  {
    value: "true",
    label: "Yes",
    icon: "✅",
  },
];

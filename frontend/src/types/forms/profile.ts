import { z } from "zod";

import { Gamemode, Platform, Rank, Region, Roles } from "@/types/types";

export const formSchema = z.object({
  name: z
    .string()
    .min(3, "Username must be at least 3 characters")
    .max(14, "Username cannot exceed 14 characters")
    .regex(/^[a-zA-Z0-9.\-_'<>]+$/, "Username contains invalid characters."),
  region: z.nativeEnum(Region).or(z.string()),
  platform: z.nativeEnum(Platform).or(z.string()),
  gamemode: z.nativeEnum(Gamemode).or(z.string()),
  role: z.enum(Roles).or(z.string()),
  rank: z.nativeEnum(Rank).or(z.string()),
  voiceChat: z.boolean(),
  mic: z.boolean(),
  characters: z.array(z.string()),
  roleQueue: z
    .object({
      vanguards: z
        .number()
        .min(0, "Please select a minimum of 0 vanguards")
        .max(6, "Please select a maximum of 6 vanguards"),
      duelists: z
        .number()
        .min(0, "Please select a minimum of 0 duelists")
        .max(6, "Please select a maximum of 6 duelists"),
      strategists: z
        .number()
        .min(0, "Please select a minimum of 0 strategists")
        .max(6, "Please select a maximum of 6 strategists"),
      sum: z.any().optional(), // Used to render the error message
    })
    .optional(),
  groupSettings: z
    .object({
      platform: z.nativeEnum(Platform),
      voiceChat: z.boolean(),
      mic: z.boolean(),
    })
    .optional(),
});

export const emptyState = {
  name: "",
  region: "",
  platform: "",
  gamemode: "",
  role: "",
  rank: "",
  characters: [] as string[],
  voiceChat: false,
  mic: false,
  roleQueue: {
    vanguards: 0,
    duelists: 0,
    strategists: 0,
  },
  groupSettings: {
    platform: "pc" as Platform,
    voiceChat: false,
    mic: false,
  },
};

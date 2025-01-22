import { z } from "zod";

import { containsProfanity } from "@/lib";

export const chatSchema = z.object({
  message: z
    .string()
    .min(1, "Please enter a message.")
    .max(1000, "Message must be less than 1000 characters.")
    .refine((val) => !val.includes("http"), {
      message: "Message cannot contain links.",
    })
    .refine((val) => !containsProfanity(val), {
      message: "Message contains restricted words.",
    }),
});

export const emptyState = {
  message: "",
};

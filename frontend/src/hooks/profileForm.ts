import { useMemo } from "react";
import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "@tanstack/react-router";
import { z } from "zod";

import { HTTPError } from "@/api";
import { useCreateGroup, useToast } from "@/hooks";
import { emptyState, formSchema, Profile, ProfileFormType } from "@/types";

interface UseProfileFormProps {
  profileFormType: ProfileFormType;
  profile?: Profile;
  setProfile: (profile: Profile) => void;
}

const formMessages: Record<
  ProfileFormType,
  { success: string; error: string }
> = {
  find: { success: "Preferences saved", error: "Failed to save preferences." },
  create: { success: "Group created!", error: "Failed to create group." },
  profile: {
    success: "Preferences saved",
    error: "Failed to save preferences.",
  },
};

export function useProfileForm({
  profileFormType,
  profile,
  setProfile,
}: UseProfileFormProps) {
  const defaultValues = useMemo(
    () => ({
      ...emptyState,
      ...profile,
    }),
    [profile],
  );

  const router = useRouter();
  const { toast } = useToast();
  const createGroup = useCreateGroup();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      ...defaultValues,
      ...profile,
    },
  });

  const isGroup = profileFormType === "create";

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      let groupId = "",
        playerId = 0;
      if (isGroup) {
        const data = await createGroup({
          profile: values as Profile,
        });
        groupId = data.groupId;
        playerId = data.playerId;
      }

      setProfile({
        id: playerId,
        ...values,
      } as Profile);

      toast({
        title: formMessages[profileFormType].success,
        variant: "success",
      });

      if (isGroup) {
        router.navigate({ to: `/groups/${groupId}` });
      } else if (profileFormType === "find") {
        router.navigate({ to: "/browse", search: { queue: true } });
      }
    } catch (error) {
      if (!(error instanceof HTTPError)) {
        toast({
          title: formMessages[profileFormType].error,
          description: "Please try again.",
          variant: "destructive",
        });
        return;
      }
      toast({
        title: error.statusText,
        description: error.message,
        variant: "destructive",
      });
    }
  };

  function onClear() {
    form.reset(emptyState);
  }

  const onReset = () => {
    form.reset(defaultValues);
  };

  return {
    form,
    onSubmit,
    onClear,
    onReset,
    isSubmitting: form.formState.isSubmitting,
    isDirty: form.formState.isDirty,
  };
}

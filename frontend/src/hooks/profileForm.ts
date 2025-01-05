import { useMemo } from "react";
import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "@tanstack/react-router";
import { z } from "zod";

import { useToast, useUpsertGroup } from "@/hooks";
import { emptyState, formSchema, Profile, ProfileFormType } from "@/types";

interface UseProfileFormProps {
  profileFormType: ProfileFormType;
  profile?: Profile;
  setProfile: (profile: Profile) => void;
}

const formMessages: Record<ProfileFormType, { success: string; error: string }> = {
  find: { success: "Preferences saved", error: "Failed to save preferences." },
  create: { success: "Group created!", error: "Failed to create group." },
  profile: { success: "Preferences saved", error: "Failed to save preferences." },
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
    [profile]
  );

  const router = useRouter();
  const { toast } = useToast();
  const upsertGroup = useUpsertGroup();

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
      let groupId = "";
      if (isGroup) {
        groupId = await upsertGroup({
          profile: values as Profile,
          id: "",
        });
      }

      setProfile({
        id: profile?.id ?? 0,
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
    } catch {
      toast({
        title: formMessages[profileFormType].error,
        description: "Please try again.",
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

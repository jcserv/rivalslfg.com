import { useMemo } from "react";
import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "@tanstack/react-router";
import { z } from "zod";

import { useToast, useUpsertGroup } from "@/hooks";
import { emptyState, formSchema, Profile } from "@/types";

interface UseProfileFormProps {
  isGroup?: boolean;
  profile?: Profile;
  setProfile: (profile: Profile) => void;
}

export function useProfileForm({
  isGroup = false,
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
  const upsertGroup = useUpsertGroup();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      ...defaultValues,
      ...profile,
    },
  });

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
        id: 1, // TODO: This should be generated server-side
        ...values,
      } as Profile);

      if (!isGroup) {
        toast({
          title: "Preferences saved",
          variant: "success",
        });
        return;
      }

      toast({
        title: "Group created",
        variant: "success",
      });
      router.navigate({ to: `/groups/${groupId}` });
    } catch {
      toast({
        title: "Failed to save preferences",
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

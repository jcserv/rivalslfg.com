import { useState } from "react";
import { UseFormReturn } from "react-hook-form";

import { z } from "zod";

import {
  CharactersField,
  GamemodeField,
  MicField,
  PlatformField,
  RankField,
  RegionField,
  RoleQueueEnabledField,
  RoleQueueFields,
  RolesField,
  UsernameField,
  VoiceChatField,
} from "@/components/ProfileForm.Fields";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Form } from "@/components/ui/form";
import { useProfileForm } from "@/hooks";
import {
  formSchema,
  getSubmitButtonLabel,
  Profile,
  ProfileFormType,
} from "@/types";

interface ProfileFormProps {
  profileFormType: ProfileFormType;
  profile?: Profile;
  setProfile: (profile: Profile) => void;
}

export function ProfileForm({
  profileFormType,
  profile,
  setProfile,
}: ProfileFormProps) {
  const [roleQueueEnabled, setRoleQueueEnabled] = useState(
    profile?.roleQueue ? true : false,
  );

  const { form, onSubmit, onClear, onReset } = useProfileForm({
    profileFormType,
    profile,
    setProfile,
  });

  const isGroup = profileFormType === "create";

  const personalInfo = (form: UseFormReturn<z.infer<typeof formSchema>>) => (
    <div>
      <div className="grid grid-cols-12 gap-4">
        <div className="col-span-12">
          <UsernameField form={form} />
        </div>
      </div>
      <div className="grid grid-cols-12 gap-4">
        <div className="col-span-6">
          <RegionField form={form} />
        </div>
        <div className="col-span-6">
          <PlatformField form={form} />
        </div>
      </div>
      <div className="grid grid-cols-12 gap-4">
        <div className="col-span-6">
          <GamemodeField form={form} />
        </div>
        <div className="col-span-6">
          <RankField form={form} />
        </div>
      </div>
      <div className="grid grid-cols-12 gap-4">
        <div className="col-span-6">
          <RolesField form={form} />
        </div>
        <div className="col-span-6">
          <CharactersField form={form} />
        </div>
      </div>
      <div className="grid grid-cols-12 gap-4">
        <div className="col-span-6">
          <VoiceChatField form={form} isGroup={false} />
        </div>
        <div className="col-span-6">
          <MicField form={form} isGroup={false} />
        </div>
      </div>
    </div>
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>Preferences</CardTitle>
        <CardDescription>
          Make changes to your matchmaking preferences here. Click submit when
          you&apos;re done.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="space-y-8 max-w-3xl mx-auto"
          >
            {profile?.name ? (
              <Accordion type="single" collapsible className="w-full">
                <AccordionItem value="item-1">
                  <AccordionTrigger>User Info</AccordionTrigger>
                  <AccordionContent className="flex flex-col gap-2 m-auto">
                    {personalInfo(form)}
                  </AccordionContent>
                </AccordionItem>
              </Accordion>
            ) : (
              personalInfo(form)
            )}
            <Accordion type="single" collapsible className="w-full">
              <AccordionItem value="item-1">
                <AccordionTrigger>
                  {isGroup ? "Group" : "Advanced"}
                </AccordionTrigger>
                <AccordionContent className="flex flex-col gap-2 m-auto">
                  <RoleQueueEnabledField
                    form={form}
                    roleQueueEnabled={roleQueueEnabled}
                    setRoleQueueEnabled={setRoleQueueEnabled}
                  />
                  <div className="grid grid-cols-12 gap-4">
                    <div className="col-span-6">
                      {roleQueueEnabled && <RoleQueueFields form={form} />}
                    </div>
                  </div>
                  {isGroup && (
                    <div className="grid grid-cols-12 gap-4">
                      <div className="col-span-6">
                        <VoiceChatField form={form} isGroup={true} />
                      </div>
                      <div className="col-span-6">
                        <MicField form={form} isGroup={true} />
                      </div>
                    </div>
                  )}
                  {form.formState.errors?.roleQueue?.sum?.message && (
                    <p className="text-sm font-medium text-destructive px-2">
                      {String(form.formState.errors?.roleQueue?.sum?.message)}
                    </p>
                  )}
                </AccordionContent>
              </AccordionItem>
            </Accordion>
            <div className="flex justify-end space-x-2">
              <Button type="button" variant="destructive" onClick={onClear}>
                Clear
              </Button>
              <Button type="button" variant="secondary" onClick={onReset}>
                Reset
              </Button>
              <Button type="submit" variant="success">
                {getSubmitButtonLabel(profileFormType)}
              </Button>
            </div>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}

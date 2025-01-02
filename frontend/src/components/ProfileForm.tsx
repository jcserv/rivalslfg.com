import { useMemo, useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "@tanstack/react-router";
import { Check, ChevronsUpDown } from "lucide-react";
import { useForm, UseFormReturn } from "react-hook-form";
import { z } from "zod";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
  Input,
  Label,
  Switch,
} from "@/components/ui";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  MultiSelector,
  MultiSelectorContent,
  MultiSelectorInput,
  MultiSelectorItem,
  MultiSelectorList,
  MultiSelectorTrigger,
} from "@/components/ui/multi-select";
import { useToast, useUpsertGroup } from "@/hooks";
import { cn } from "@/lib/utils";
import {
  Gamemode,
  Platform,
  Profile,
  Rank,
  Region,
  Roles,
  TEAM_SIZE,
} from "@/types";

import characters from "@/assets/characters.json";
import gamemodes from "@/assets/gamemodes.json";
import platforms from "@/assets/platforms.json";
import ranks from "@/assets/ranks.json";
import regions from "@/assets/regions.json";
import roles from "@/assets/roles.json";

const formSchema = z.object({
  name: z
    .string()
    .min(3, "Username must be at least 3 characters")
    .max(14, "Username cannot exceed 14 characters")
    .regex(/^[a-zA-Z0-9.\-_'<>]+$/, "Username contains invalid characters."),
  region: z.nativeEnum(Region).or(z.string()),
  platform: z.nativeEnum(Platform).or(z.string()),
  gamemode: z.nativeEnum(Gamemode).or(z.string()),
  roles: z
    .array(z.enum(Roles).or(z.string()))
    .min(1, "Please select at least one role"),
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
    .optional()
    .refine(
      (data) =>
        data &&
        data?.vanguards + data?.duelists + data?.strategists === TEAM_SIZE,
      {
        message:
          "Number of desired vanguards, duelists, and strategists must add up to 6",
        path: ["sum"],
      },
    ),
  groupSettings: z
    .object({
      platforms: z.array(z.nativeEnum(Platform)),
      voiceChat: z.boolean(),
      mic: z.boolean(),
    })
    .optional(),
});

interface ProfileFormProps {
  isGroup?: boolean;
  profile?: Profile;
  setProfile: (profile: Profile) => void;
}

export function ProfileForm({
  isGroup = false,
  profile,
  setProfile,
}: ProfileFormProps) {
  const [roleQueueEnabled, setRoleQueueEnabled] = useState(
    profile?.roleQueue ? true : false,
  );
  const upsertGroup = useUpsertGroup();
  const { toast } = useToast();
  const router = useRouter();

  const defaultValues = useMemo(
    () => ({
      region: "",
      platform: "",
      gamemode: "",
      roles: [] as string[],
      rank: "",
      characters: [] as string[],
      voiceChat: false,
      mic: false,
      roleQueue: {
        vanguards: 2,
        duelists: 2,
        strategists: 2,
      },
      groupSettings: {
        platforms: [],
        voiceChat: false,
        mic: false,
      },
      ...profile,
    }),
    [profile],
  );

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues,
  });

  function onReset() {
    form.reset(defaultValues);
  }

  async function onSubmit(values: z.infer<typeof formSchema>) {
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
    } catch (error) {
      console.error("Form submission error", error);
      toast({
        title: "Failed to save preferences",
        description: "Please try again.",
        variant: "destructive",
      });
    }
  }

  const personalInfo = (form: UseFormReturn<z.infer<typeof formSchema>>) => (
    <div>
      <div className="grid grid-cols-12 gap-4">
        <div className="col-span-12">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem className="m-2">
                <FormLabel>Username</FormLabel>
                <FormDescription>
                  This should match your in-game name in Marvel Rivals.
                </FormDescription>
                <Input id="name" {...field} />
              </FormItem>
            )}
          />
        </div>
      </div>
      <div className="grid grid-cols-12 gap-4">
        <div className="col-span-6">
          <FormField
            control={form.control}
            name="region"
            render={({ field }) => (
              <FormItem className="m-2">
                <FormLabel>Region</FormLabel>
                <Select onValueChange={field.onChange} value={field.value}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select your region" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {regions.map((region) => (
                      <SelectItem key={region.value} value={region.value}>
                        {region.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
        <div className="col-span-6">
          <FormField
            control={form.control}
            name="platform"
            render={({ field }) => (
              <FormItem className="m-2">
                <FormLabel>Platform</FormLabel>
                <Select onValueChange={field.onChange} value={field.value}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select your platform" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {platforms.map((platform) => (
                      <SelectItem key={platform.value} value={platform.value}>
                        {platform.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>

                <FormMessage />
              </FormItem>
            )}
          />
        </div>
      </div>
      <div className="grid grid-cols-12 gap-4">
        <div className="col-span-6">
          <FormField
            control={form.control}
            name="gamemode"
            render={({ field }) => (
              <FormItem className="m-2">
                <FormLabel>Gamemode</FormLabel>
                <Select onValueChange={field.onChange} value={field.value}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select the gamemode you want to play" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {gamemodes.map((gamemode) => (
                      <SelectItem key={gamemode.value} value={gamemode.value}>
                        {gamemode.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>

                <FormMessage />
              </FormItem>
            )}
          />
        </div>
        <div className="col-span-6">
          <FormField
            control={form.control}
            name="rank"
            render={({ field }) => (
              <FormItem className="m-2">
                <FormLabel>Rank</FormLabel>
                <Popover>
                  <PopoverTrigger asChild>
                    <FormControl>
                      <Button
                        variant="outline"
                        role="combobox"
                        className={cn(
                          "w-full justify-between",
                          !field.value && "text-muted-foreground",
                        )}
                      >
                        {field.value
                          ? ranks.find((rank) => rank.value === field.value)
                              ?.label
                          : "Select your rank"}
                        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                      </Button>
                    </FormControl>
                  </PopoverTrigger>
                  <PopoverContent className="w-[200px] p-0">
                    <Command>
                      <CommandInput placeholder="Search ranks..." />
                      <CommandList>
                        <CommandEmpty>No results found.</CommandEmpty>
                        <CommandGroup>
                          {ranks.map((rank) => (
                            <CommandItem
                              value={rank.label}
                              key={rank.value}
                              onSelect={() => {
                                form.setValue("rank", rank.value);
                              }}
                            >
                              <Check
                                className={cn(
                                  "mr-2 h-4 w-4",
                                  rank.value === field.value
                                    ? "opacity-100"
                                    : "opacity-0",
                                )}
                              />
                              {rank.label}
                            </CommandItem>
                          ))}
                        </CommandGroup>
                      </CommandList>
                    </Command>
                  </PopoverContent>
                </Popover>
                <FormDescription>
                  You&apos;ll be matched with players within adjacent ranks if
                  your selected gamemode is competitive.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
      </div>
      <div className="grid grid-cols-12 gap-4">
        <div className="col-span-6">
          <FormField
            control={form.control}
            name="roles"
            render={({ field }) => (
              <FormItem className="m-2">
                <FormLabel>Roles</FormLabel>
                <FormControl>
                  <MultiSelector
                    values={field.value}
                    onValuesChange={field.onChange}
                    loop
                    className="max-w-xs"
                  >
                    <MultiSelectorTrigger>
                      <MultiSelectorInput placeholder="Select your preferred role(s)" />
                    </MultiSelectorTrigger>
                    <MultiSelectorContent>
                      <MultiSelectorList className="z-50">
                        {roles.map((role) => (
                          <MultiSelectorItem key={role} value={role}>
                            {role}
                          </MultiSelectorItem>
                        ))}
                      </MultiSelectorList>
                    </MultiSelectorContent>
                  </MultiSelector>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
        <div className="col-span-6">
          <FormField
            control={form.control}
            name="characters"
            render={({ field }) => (
              <FormItem className="m-2">
                <FormLabel>Characters</FormLabel>
                <FormControl>
                  <MultiSelector
                    values={field.value}
                    onValuesChange={field.onChange}
                    loop
                    className="max-w-xs"
                  >
                    <MultiSelectorTrigger>
                      <MultiSelectorInput placeholder="Select your preferred character(s)" />
                    </MultiSelectorTrigger>
                    <MultiSelectorContent>
                      <MultiSelectorList className="h-28">
                        {characters.map((character) => (
                          <MultiSelectorItem
                            key={character.name}
                            value={character.name}
                          >
                            {character.name} - {character.role}
                          </MultiSelectorItem>
                        ))}
                      </MultiSelectorList>
                    </MultiSelectorContent>
                  </MultiSelector>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
      </div>
      <div className="grid grid-cols-12 gap-4">
        <div className="col-span-6">
          <FormField
            control={form.control}
            name="voiceChat"
            render={({ field }) => (
              <FormItem className="flex flex-row items-start gap-2 m-2">
                <div className="space-y-0.5">
                  <FormLabel className="self-center leading-none mt-1">
                    Voice Chat
                  </FormLabel>
                  <FormDescription>
                    This indicates whether you are able to listen via voice
                    chat, either in-game or through Discord.
                  </FormDescription>
                </div>
                <FormControl>
                  <Switch
                    id="voiceChat"
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
        <div className="col-span-6">
          <FormField
            control={form.control}
            name="mic"
            render={({ field }) => (
              <FormItem className="flex flex-row items-start gap-2 m-2">
                <div className="space-y-0.5">
                  <FormLabel className="self-center leading-none mt-1">
                    Mic
                  </FormLabel>
                  <FormDescription>
                    This indicates whether you are able to speak via voice chat,
                    either in-game or through Discord.
                  </FormDescription>
                </div>
                <FormControl>
                  <Switch
                    id="mic"
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
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
                  <div className="flex items-center space-x-2 p-2">
                    <Checkbox
                      id="roleQueue"
                      checked={roleQueueEnabled}
                      onClick={() => setRoleQueueEnabled(!roleQueueEnabled)}
                    />
                    <Label
                      htmlFor="roleQueue"
                      className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                    >
                      Enable Role Queue
                    </Label>
                  </div>
                  <div>
                    <FormDescription>
                      If this setting is enabled, you&apos;ll be matched to
                      players according to your desired role counts.
                    </FormDescription>
                  </div>
                  {roleQueueEnabled && (
                    <div className="flex flex-row gap-4">
                      <div className="flex flex-col space-y-2">
                        <FormField
                          control={form.control}
                          name="roleQueue.vanguards"
                          render={({ field }) => (
                            <FormItem className="mx-2">
                              <FormLabel>Vanguards</FormLabel>
                              <FormControl>
                                <Input
                                  type="number"
                                  id="roleQueue.vanguards"
                                  value={field.value}
                                  onChange={(e) =>
                                    field.onChange(+e.target.value)
                                  }
                                  min={0}
                                  max={6}
                                  className="w-[75px]"
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>

                      <div className="flex flex-col space-y-2">
                        <FormField
                          control={form.control}
                          name="roleQueue.duelists"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Duelists</FormLabel>
                              <FormControl>
                                <Input
                                  type="number"
                                  id="roleQueue.duelists"
                                  value={field.value}
                                  onChange={(e) =>
                                    field.onChange(+e.target.value)
                                  }
                                  min={0}
                                  max={6}
                                  className="w-[75px]"
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>

                      <div className="flex flex-col space-y-2">
                        <FormField
                          control={form.control}
                          name="roleQueue.strategists"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Strategists</FormLabel>
                              <FormControl>
                                <Input
                                  type="number"
                                  id="roleQueue.strategists"
                                  value={field.value}
                                  onChange={(e) =>
                                    field.onChange(+e.target.value)
                                  }
                                  min={0}
                                  max={6}
                                  className="w-[75px]"
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>
                    </div>
                  )}
                  {isGroup && (
                    <div className="grid grid-cols-12 gap-4">
                      <div className="col-span-6">
                        <FormField
                          control={form.control}
                          name="groupSettings.voiceChat"
                          render={({ field }) => (
                            <FormItem className="flex flex-row items-start gap-2 m-2">
                              <div className="space-y-0.5">
                                <FormLabel className="self-center leading-none mt-1">
                                  Voice Chat
                                </FormLabel>
                                <FormDescription>
                                  This indicates whether you want all players to
                                  be able to listen via voice chat, either
                                  in-game or through Discord.
                                </FormDescription>
                              </div>
                              <FormControl>
                                <Switch
                                  id="groupSettings.voiceChat"
                                  checked={field.value}
                                  onCheckedChange={field.onChange}
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>
                      <div className="col-span-6">
                        <FormField
                          control={form.control}
                          name="groupSettings.mic"
                          render={({ field }) => (
                            <FormItem className="flex flex-row items-start gap-2 m-2">
                              <div className="space-y-0.5">
                                <FormLabel className="self-center leading-none mt-1">
                                  Mic
                                </FormLabel>
                                <FormDescription>
                                  This indicates whether you want all players to
                                  be able to speak via voice chat, either
                                  in-game or through Discord.
                                </FormDescription>
                              </div>
                              <FormControl>
                                <Switch
                                  id="groupSettings.mic"
                                  checked={field.value}
                                  onCheckedChange={field.onChange}
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
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
            <div className="flex space-x-2">
              <Button type="button" variant="destructive" onClick={onReset}>
                Clear
              </Button>
              <Button type="submit">Submit</Button>
            </div>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}

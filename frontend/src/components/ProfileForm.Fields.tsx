import { UseFormReturn } from "react-hook-form";

import { Check, ChevronsUpDown } from "lucide-react";
import { z } from "zod";

import {
  Button,
  Checkbox,
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  Input,
  Label,
  MultiSelect,
  Popover,
  PopoverContent,
  PopoverTrigger,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Switch,
} from "@/components/ui";
import { cn, toTitleCase } from "@/lib/utils";
import { formSchema } from "@/types";

import characters from "@/assets/characters.json";
import gamemodes from "@/assets/gamemodes.json";
import platforms from "@/assets/platforms.json";
import ranks from "@/assets/ranks.json";
import regions from "@/assets/regions.json";
import roles from "@/assets/roles.json";

interface FormFieldProps {
  form: UseFormReturn<z.infer<typeof formSchema>>;
}

export function UsernameField({ form }: FormFieldProps) {
  return (
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
          <FormMessage />
        </FormItem>
      )}
    />
  );
}

export function RegionField({ form }: FormFieldProps) {
  return (
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
  );
}

export function PlatformField({ form }: FormFieldProps) {
  return (
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
  );
}

export function GamemodeField({ form }: FormFieldProps) {
  return (
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
  );
}

export function RankField({ form }: FormFieldProps) {
  return (
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
                    ? ranks.find((rank) => rank.value === field.value)?.label
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
            You&apos;ll be matched with players within adjacent ranks if your
            selected gamemode is competitive.
          </FormDescription>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}

export function RolesField({ form }: FormFieldProps) {
  return (
    <FormField
      control={form.control}
      name="role"
      render={({ field }) => (
        <FormItem className="m-2">
          <FormLabel>Roles</FormLabel>
          <Select onValueChange={field.onChange} value={field.value}>
            <FormControl>
              <SelectTrigger>
                <SelectValue placeholder="Select your preferred role" />
              </SelectTrigger>
            </FormControl>
            <SelectContent>
              {roles.map((role) => (
                <SelectItem key={role} value={role}>
                  {toTitleCase(role)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}

export function CharactersField({ form }: FormFieldProps) {
  return (
    <FormField
      control={form.control}
      name="characters"
      render={({ field }) => (
        <FormItem className="m-2">
          <FormLabel>Characters</FormLabel>
          <FormControl>
            <MultiSelect
              value={field.value}
              defaultValue={field.value}
              options={characters}
              onValueChange={field.onChange}
              placeholder="Select your preferred character(s)"
              variant="inverted"
              animation={2}
              maxCount={3}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}

interface RoleQueueEnabledFieldProps extends FormFieldProps {
  roleQueueEnabled: boolean;
  setRoleQueueEnabled: (enabled: boolean) => void;
}

export function RoleQueueEnabledField({
  roleQueueEnabled,
  setRoleQueueEnabled,
}: RoleQueueEnabledFieldProps) {
  return (
    <>
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
          If this setting is enabled, you&apos;ll be matched to players
          according to your desired role counts.
        </FormDescription>
      </div>
    </>
  );
}

interface RoleQueueFieldProps extends FormFieldProps {
  role: "vanguards" | "duelists" | "strategists";
  label: string;
}

export function RoleQueueField({ form, role, label }: RoleQueueFieldProps) {
  return (
    <FormField
      control={form.control}
      name={`roleQueue.${role}`}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{label}</FormLabel>
          <Input
            type="number"
            value={field.value}
            onChange={(e) => field.onChange(+e.target.value)}
            min={0}
            max={6}
            className="w-[75px]"
          />
          <FormMessage />
        </FormItem>
      )}
    />
  );
}

export function RoleQueueFields({ form }: FormFieldProps) {
  return (
    <div className="flex flex-row gap-4">
      <RoleQueueField form={form} role="vanguards" label="Vanguards" />
      <RoleQueueField form={form} role="duelists" label="Duelists" />
      <RoleQueueField form={form} role="strategists" label="Strategists" />
    </div>
  );
}

interface SwitchFieldProps extends FormFieldProps {
  isGroup?: boolean;
}

export function VoiceChatField({ form, isGroup = false }: SwitchFieldProps) {
  const name = isGroup ? "groupSettings.voiceChat" : "voiceChat";
  const description = `This indicates whether ${
    isGroup ? "you want all players" : "you are able"
  } to listen via voice chat, either in-game or through Discord.`;

  return (
    <FormField
      control={form.control}
      name={name}
      render={({ field }) => (
        <FormItem className="flex flex-row items-start gap-2 m-2">
          <div className="space-y-0.5">
            <FormLabel className="self-center leading-none mt-1">
              Voice Chat
            </FormLabel>
            <FormDescription>{description}</FormDescription>
          </div>
          <FormControl>
            <Switch
              id={name}
              checked={field.value}
              onCheckedChange={field.onChange}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}

export function MicField({ form, isGroup = false }: SwitchFieldProps) {
  const name = isGroup ? "groupSettings.mic" : "mic";
  const description = `This indicates whether ${
    isGroup ? "you want all players" : "you are able"
  } to speak via voice chat, either in-game or through Discord.`;

  return (
    <FormField
      control={form.control}
      name={name}
      render={({ field }) => (
        <FormItem className="flex flex-row items-start gap-2 m-2">
          <div className="space-y-0.5">
            <FormLabel className="self-center leading-none mt-1">Mic</FormLabel>
            <FormDescription>{description}</FormDescription>
          </div>
          <FormControl>
            <Switch
              id={name}
              checked={field.value}
              onCheckedChange={field.onChange}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}

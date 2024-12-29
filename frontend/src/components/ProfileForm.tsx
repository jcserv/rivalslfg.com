import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { cn } from "@/lib/utils";
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
import { Check, ChevronsUpDown } from "lucide-react";
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

import characters from "@/assets/characters.json";
import gamemodes from "@/assets/gamemodes.json";
import platforms from "@/assets/platforms.json";
import ranks from "@/assets/ranks.json";
import regions from "@/assets/regions.json";
import roles from "@/assets/roles.json";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
  Input,
  Label,
} from "./ui";
import { Checkbox } from "./ui/checkbox";
import { useState } from "react";

const TEAM_SIZE = 6;

const formSchema = z.object({
  name: z.string().min(1, "Please enter your in-game name"),
  region: z.string().min(1, "Please select a region"),
  platform: z.string().min(1, "Please select a platform"),
  gamemode: z.string().min(1, "Please select a gamemode"),
  roles: z.array(z.string()).min(1, "Please select at least one role"),
  rank: z.string().min(1, "Please select a rank"),
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
});
export function ProfileForm() {
  const [roleQueueEnabled, setRoleQueueEnabled] = useState(false);

  const defaultValues = {
    region: "",
    platform: "",
    gamemode: "",
    roles: [] as string[],
    rank: "",
    characters: [] as string[],
    roleQueue: {
      vanguards: 2,
      duelists: 2,
      strategists: 2,
    },
  };

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues,
  });

  function onReset() {
    form.reset(defaultValues);
  }

  function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      console.log(values);
      toast(
        <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
          <code className="text-white">{JSON.stringify(values, null, 2)}</code>
        </pre>,
      );
    } catch (error) {
      console.error("Form submission error", error);
      toast.error("Failed to submit the form. Please try again.");
    }
  }

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
            <div className="grid grid-cols-12 gap-4">
              <div className="col-span-12">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
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
                    <FormItem>
                      <FormLabel>Region</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        value={field.value}
                      >
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
                    <FormItem>
                      <FormLabel>Platform</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        value={field.value}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Select your platform" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {platforms.map((platform) => (
                            <SelectItem
                              key={platform.value}
                              value={platform.value}
                            >
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
                    <FormItem>
                      <FormLabel>Gamemode</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        value={field.value}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Select the gamemode you want to play" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {gamemodes.map((gamemode) => (
                            <SelectItem
                              key={gamemode.value}
                              value={gamemode.value}
                            >
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
                    <FormItem>
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
                                ? ranks.find(
                                    (rank) => rank.value === field.value,
                                  )?.label
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
                        You&apos;ll be matched with players within adjacent
                        ranks if your selected gamemode is competitive.
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
                    <FormItem>
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
                            <MultiSelectorList>
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
                    <FormItem>
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
                            <MultiSelectorList>
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
            {/* <Accordion type="single" collapsible className="w-full">
              <AccordionItem value="item-1">
                <AccordionTrigger>Social</AccordionTrigger>
                <AccordionContent className="flex flex-col gap-2 m-auto">
                  Discord
                </AccordionContent>
              </AccordionItem>
            </Accordion> */}
            <Accordion type="single" collapsible className="w-full">
              <AccordionItem value="item-1">
                <AccordionTrigger>Advanced</AccordionTrigger>
                <AccordionContent className="flex flex-col gap-2 m-auto">
                  <div className="flex items-center space-x-2 px-2">
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
                      groups according to your desired role counts.
                    </FormDescription>
                  </div>
                  {roleQueueEnabled && (
                    <div className="flex flex-row gap-4">
                      <div className="flex flex-col space-y-2">
                        <FormField
                          control={form.control}
                          name="roleQueue.vanguards"
                          render={({ field }) => (
                            <FormItem>
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

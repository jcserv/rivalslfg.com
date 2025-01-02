import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  Skeleton,
  Switch,
} from "@/components/ui";

const formSchema = z.object({
  open: z.boolean(),
});

interface GroupControlsProps {
  isGroupOpen: boolean;
  canUserAccessGroup: boolean | null;
}

export function GroupControls({
  isGroupOpen,
  canUserAccessGroup,
}: GroupControlsProps) {
  const passcode = "abcd";

  const defaultValues = {
    open: isGroupOpen,
  };

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues,
  });

  const isClosed = form.watch("open") === false;

  return (
    <Card>
      <CardHeader>
        <div className="flex flex-row items-center justify-between">
          <CardTitle>Settings</CardTitle>
          {isClosed &&
            (canUserAccessGroup ? (
              <p>
                <strong>Passcode:</strong> {passcode}
              </p>
            ) : (
              <Skeleton className="h-12 w-1/2 rounded-xl" />
            ))}
        </div>
      </CardHeader>
      <CardContent className="h-full">
        {!canUserAccessGroup ? (
          <Skeleton className="h-[125px] w-full rounded-xl" />
        ) : (
          <Form {...form}>
            <form>
              <div className="grid grid-cols-12 gap-4">
                <div className="col-span-12">
                  <FormField
                    control={form.control}
                    name="open"
                    render={({ field }) => (
                      <FormItem className="flex flex-row items-start gap-2">
                        <div className="space-y-0.5">
                          <FormLabel
                            htmlFor="open"
                            className="self-center leading-none mt-1"
                          >
                            Open
                          </FormLabel>
                          <FormDescription>
                            When this is on, players can join the group.
                          </FormDescription>
                        </div>
                        <FormControl>
                          <Switch
                            id="open"
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
            </form>
          </Form>
        )}
      </CardContent>
    </Card>
  );
}

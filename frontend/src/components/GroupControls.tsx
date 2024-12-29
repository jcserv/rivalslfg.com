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
  Switch,
} from "@/components/ui";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
  open: z.boolean(),
});

interface GroupControlsProps {
  passcode: string;
}

export function GroupControls({ passcode }: GroupControlsProps) {
  const defaultValues = {
    open: true,
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
          {isClosed && (
            <p>
              <strong>Passcode:</strong> {passcode}
            </p>
          )}
        </div>
      </CardHeader>
      <CardContent>
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
      </CardContent>
    </Card>
  );
}

import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useProfile } from "@/hooks";
import { Profile } from "@/types";

import { BackButton } from "./BackButton";

const formSchema = z.object({
  passcode: z
    .string()
    .min(1, "Please enter a passcode")
    .max(4, "Passcode must be 4 characters"),
});

interface AccessGroupDialogProps {
  open: boolean;
  onJoin: (p: Profile, passcode: string) => void;
}

export function AccessGroupDialog({ open, onJoin }: AccessGroupDialogProps) {
  const [profile] = useProfile();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      passcode: "",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    onJoin(profile, values.passcode);
  }

  return (
    <Dialog open={open}>
      <DialogContent className="sm:max-w-[425px]" showClose={false}>
        <DialogHeader>
          <DialogTitle>Enter group passcode</DialogTitle>
          <DialogDescription>
            This is a private group. Please enter the passcode to join.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="grid gap-4">
              <div className="grid items-center gap-4">
                <FormField
                  control={form.control}
                  name="passcode"
                  render={({ field }) => (
                    <FormItem className="mb-2">
                      <FormLabel htmlFor="passcode" className="text-right">
                        Passcode
                      </FormLabel>
                      <Input id="passcode" {...field} className="col-span-3" />
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </div>
            <DialogFooter className="flex flex-row justify-between sm:justify-between">
              <BackButton text="Back" />
              <Button type="submit">Submit</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

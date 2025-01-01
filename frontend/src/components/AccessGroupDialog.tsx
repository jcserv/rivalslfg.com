import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
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
import { useToast } from "@/hooks/use-toast";

const formSchema = z.object({
  passcode: z
    .string()
    .min(1, "Please enter a passcode")
    .max(4, "Passcode must be 4 characters"),
});

interface AccessGroupDialogProps {
  open: boolean;
  onClose: () => void;
}

export function AccessGroupDialog({ open, onClose }: AccessGroupDialogProps) {
  const { toast } = useToast();
  const router = useRouter();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      passcode: "",
    },
  });

  function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      console.log("values", values);
      toast({
        title: "Access granted",
        variant: "success",
      });
      onClose();
    } catch (error) {
      console.error("Form submission error", error);
      toast({
        title: "Access denied",
        description: "Please try again.",
        variant: "destructive",
      });
    }
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
              <Button variant="outline" onClick={() => router.history.back()}>
                Back
              </Button>
              <Button type="submit">Submit</Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

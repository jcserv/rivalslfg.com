import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { chatSchema, emptyState } from "@/types/forms/chat";

export function useChatForm({
  onSubmit,
}: {
  onSubmit: (message: string) => void;
}) {
  const form = useForm<z.infer<typeof chatSchema>>({
    resolver: zodResolver(chatSchema),
    defaultValues: emptyState,
  });

  const handleSubmit = form.handleSubmit((values) => {
    onSubmit(values.message);
    form.reset();
  });

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e as unknown as React.BaseSyntheticEvent);
    }
  };

  return {
    form,
    handleSubmit,
    handleKeyDown,
  };
}

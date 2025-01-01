import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui";

interface TooltipItemProps {
  children: React.ReactNode;
  content: React.ReactNode;
  enabled?: boolean;
}

export function TooltipItem({
  children,
  content,
  enabled = true,
}: TooltipItemProps) {
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger>{children}</TooltipTrigger>
        {enabled && <TooltipContent>{content}</TooltipContent>}
      </Tooltip>
    </TooltipProvider>
  );
}

import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui";

interface TooltipItemProps {
  children: React.ReactNode;
  content: React.ReactNode;
  disabled?: boolean;
}

export function TooltipItem({ children, content, disabled }: TooltipItemProps) {
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger>{children}</TooltipTrigger>
        {disabled && <TooltipContent>{content}</TooltipContent>}
      </Tooltip>
    </TooltipProvider>
  );
}

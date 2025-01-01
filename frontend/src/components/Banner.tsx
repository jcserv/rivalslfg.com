import {
  Accordion,
  AccordionItem,
  AccordionTrigger,
  AccordionContent,
} from "@/components/ui";
import { Info } from "lucide-react";

import { Alert, AlertTitle, AlertDescription } from "@/components/ui";
import { cn } from "@/lib/utils";

export const InfoBanner = ({ children }: { children: JSX.Element }) => (
  <Alert className="w-1/2 max-w-4xl bg-blue-400 dark:bg-blue-600 my-4 md:my-8">
    <Info className="h-4 w-4" />
    <AlertTitle>Info</AlertTitle>
    <AlertDescription className="space-y-1">{children}</AlertDescription>
  </Alert>
);

export const ErrorBanner = ({
  message,
  error,
  className,
  children,
}: {
  message: string;
  error?: string;
  className?: string;
  children?: React.ReactNode;
}) => (
  <Alert
    className={cn(
      "w-1/2 max-w-4xl bg-red-400 dark:bg-red-600 my-4 md:my-8 overflow-x-auto",
      className,
    )}
  >
    <Info className="h-4 w-4" />
    <AlertTitle>Error</AlertTitle>
    <AlertDescription className="space-y-1">
      <p>{message}</p>
      {children}
      {error && (
        <Accordion type="single" collapsible>
          <AccordionItem value="item-1">
            <AccordionTrigger>Click to View Details</AccordionTrigger>
            <AccordionContent className="overflow-y-scroll">
              <pre>{error}</pre>
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      )}
    </AlertDescription>
  </Alert>
);

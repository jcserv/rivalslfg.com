import { useRouter } from "@tanstack/react-router";

import { Button } from "./ui";

interface BackButtonProps {
  text?: string;
  link?: boolean;
  className?: string;
}

export function BackButton({
  text = "Go Back",
  link = false,
  className,
}: BackButtonProps) {
  const router = useRouter();
  const handleBack = () => {
    router.history.back();
  };
  return (
    <Button
      variant={link ? "link" : "outline"}
      onClick={handleBack}
      className={className}
    >
      {text}
    </Button>
  );
}

import { useEffect, useMemo, useState } from "react";

import { useRouter } from "@tanstack/react-router";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Progress } from "@/components/ui/progress";

interface FindGroupDialogProps {
  open: boolean;
  onClose: () => void;
}

const ACCEPT_TIME_LIMIT_SECONDS = 7;

export function FindGroupDialog({ open, onClose }: FindGroupDialogProps) {
  const router = useRouter();
  const groupId: string = "AAAA";

  const [queueTime, setQueueTime] = useState(0);
  const [groupFound, setGroupFound] = useState(false);
  const [acceptProgress, setAcceptProgress] = useState(0);

  const formattedTime = useMemo(() => {
    const minutes = Math.floor(queueTime / 60);
    const seconds = queueTime % 60;
    return `${String(minutes).padStart(2, "0")}:${String(seconds).padStart(2, "0")}`;
  }, [queueTime]);

  const titleDots = useMemo(() => {
    if (groupFound) return "";
    return ".".repeat(queueTime % 4);
  }, [queueTime, groupFound]);

  useEffect(() => {
    if (!open || groupFound) return;

    const timer = setInterval(() => {
      setQueueTime((prev) => prev + 1);
    }, 1000);

    return () => clearInterval(timer);
  }, [open, groupFound]);

  useEffect(() => {
    if (!open || groupFound) return;

    const timer = setTimeout(() => {
      setGroupFound(true);
    }, 5000);

    return () => clearTimeout(timer);
  }, [open]);

  useEffect(() => {
    if (!groupFound) return;

    const startTime = Date.now();
    const timer = setInterval(() => {
      const elapsedTime = Date.now() - startTime;
      const progress = (elapsedTime / (ACCEPT_TIME_LIMIT_SECONDS * 1000)) * 100;

      if (progress >= 100) {
        clearInterval(timer);
        onLeave();
        setAcceptProgress(0);
        return;
      }

      setAcceptProgress(progress);
    }, 100);

    return () => clearInterval(timer);
  }, [groupFound]);

  const onJoin = () => {
    router.navigate({ to: `/groups/${groupId}`, search: { join: true } });
  };

  const onLeave = () => {
    onClose();
    setQueueTime(0);
    setGroupFound(false);
  };

  return (
    <Dialog open={open}>
      <DialogContent className="sm:max-w-[425px]" showClose={false}>
        <DialogHeader>
          <DialogTitle>
            {groupFound ? "Found group!" : `Finding group${titleDots}`}
          </DialogTitle>
          <DialogDescription>Time in queue: {formattedTime}</DialogDescription>
        </DialogHeader>
        {groupFound && (
          <div className="py-4">
            <Progress value={acceptProgress} className="w-full" />
          </div>
        )}
        <DialogFooter className="flex flex-row justify-between sm:justify-between">
          <Button variant="outline" onClick={onLeave}>
            {groupFound ? "Decline" : "Leave Queue"}
          </Button>
          <Button variant="success" onClick={onJoin} disabled={!groupFound}>
            Accept
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

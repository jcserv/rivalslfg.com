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

interface FindGroupDialogProps {
  open: boolean;
  onClose: () => void;
}

export function FindGroupDialog({ open, onClose }: FindGroupDialogProps) {
  const router = useRouter();
  const groupId: string = "AAAA";

  const onJoin = () => {
    router.navigate({ to: `/groups/${groupId}` });
  };

  return (
    <Dialog open={open}>
      <DialogContent className="sm:max-w-[425px]" showClose={false}>
        <DialogHeader>
          <DialogTitle>Finding group</DialogTitle>
          <DialogDescription>Time in queue: 1 minute(s)</DialogDescription>
        </DialogHeader>
        <DialogFooter className="flex flex-row justify-between sm:justify-between">
          <Button variant="outline" onClick={onClose}>
            Leave Queue
          </Button>
          <Button variant="success" onClick={onJoin} disabled={true}>
            Join
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

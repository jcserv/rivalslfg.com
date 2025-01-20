import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { Button } from "./ui";

interface DeleteGroupDialogProps {
  open: boolean;
  onClose: () => void;
  onDelete: () => void;
}

export function DeleteGroupDialog({
  open,
  onClose,
  onDelete,
}: DeleteGroupDialogProps) {
  return (
    <Dialog open={open}>
      <DialogContent className="sm:max-w-[425px]" showClose={false}>
        <DialogHeader>
          <DialogTitle>Delete group</DialogTitle>
          <DialogDescription>
            This will delete the group and all players in the group will be
            removed.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="flex flex-row justify-between sm:justify-between">
          <Button variant="outline" onClick={onClose}>Cancel</Button>
          <Button variant="destructive" onClick={onDelete}>Confirm</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

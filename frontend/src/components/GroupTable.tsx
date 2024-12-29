import { DataTable } from "./ui/data-table";
import { columns } from "./GroupTable.Columns";

import { useLocalStorage } from "@/hooks/localStorage";
import { Group } from "@/types";

interface GroupTableProps {
  groups: Group[];
}

export function GroupTable({ groups }: GroupTableProps) {
  const fourteenDaysFromToday = new Date(
    new Date().setDate(new Date().getDate() + 14),
  );
  const [profile] = useLocalStorage("profile", {}, fourteenDaysFromToday);
  return (
    <>
      <div className="hidden h-full flex-1 flex-col space-y-8 p-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Groups</h2>
          </div>
        </div>
        <DataTable data={groups} columns={columns(profile)} />
      </div>
    </>
  );
}

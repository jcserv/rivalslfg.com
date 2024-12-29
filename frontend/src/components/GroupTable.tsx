import { DataTable } from "./ui/data-table";
import { columns } from "./GroupTable.Columns";

import { Group, Profile } from "@/types";

interface GroupTableProps {
  groups: Group[];
  profile: Profile | null;
  isProfileEmpty: boolean;
}

export function GroupTable({
  groups,
  profile,
  isProfileEmpty,
}: GroupTableProps) {
  return (
    <>
      <div className="hidden h-full flex-1 flex-col space-y-8 p-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Groups</h2>
          </div>
        </div>
        <DataTable data={groups} columns={columns(profile, isProfileEmpty)} />
      </div>
    </>
  );
}

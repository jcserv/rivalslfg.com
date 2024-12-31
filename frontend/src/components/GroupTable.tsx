import { DataTable } from "./ui/data-table";
import { columns } from "./GroupTable.Columns";

import { areRequirementsMet, getRequirements, Group, Profile } from "@/types";
import { useMemo } from "react";

interface GroupTableProps {
  groups: Group[];
  profile: Profile | undefined;
  isProfileEmpty: boolean;
}

export function GroupTable({
  groups,
  profile,
  isProfileEmpty,
}: GroupTableProps) {
  const groupTableData = useMemo(() => {
    return groups.map((group) => {
      const requirements = getRequirements(group);
      const areReqsMet = areRequirementsMet(group, requirements, profile);
      return {
        ...group,
        requirements,
        areRequirementsMet: areReqsMet,
      };
    });
  }, [groups, profile]);

  return (
    <>
      <div className="hidden h-full flex-1 flex-col space-y-8 p-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Groups</h2>
          </div>
        </div>
        <DataTable data={groupTableData} columns={columns(isProfileEmpty)} />
      </div>
    </>
  );
}

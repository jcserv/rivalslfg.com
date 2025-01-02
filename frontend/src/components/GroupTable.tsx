import { useMemo } from "react";

import { columns } from "@/components/GroupTable.Columns";
import { DataTable } from "@/components/ui";
import { useGroups } from "@/hooks";
import { areRequirementsMet, getRequirements, Group, Profile } from "@/types";

interface GroupTableProps {
  profile: Profile | undefined;
  isProfileEmpty: boolean;
}

export function GroupTable({ profile, isProfileEmpty }: GroupTableProps) {
  const { data, pagination, isLoading } = useGroups();
  const groups: Group[] = !isLoading ? data : [];

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
    <DataTable
      data={groupTableData}
      columns={columns(isProfileEmpty)}
      pagination={pagination}
      isLoading={isLoading}
    />
  );
}

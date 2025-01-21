import { useMemo } from "react";

import { ColumnFiltersState, OnChangeFn } from "@tanstack/react-table";

import { columns } from "@/components/GroupTable.Columns";
import { DataTable } from "@/components/ui";
import { useGroups } from "@/hooks";
import { areRequirementsMet, getRequirements, Group, Profile } from "@/types";

interface GroupTableProps {
  profile: Profile | undefined;
  isProfileEmpty: boolean;
  filters: ColumnFiltersState;
  handleColumnFiltersChange: OnChangeFn<ColumnFiltersState>;
  handleClearFilters: () => void;
}

export function GroupTable({
  profile,
  isProfileEmpty,
  filters,
  handleColumnFiltersChange,
  handleClearFilters,
}: GroupTableProps) {
  const { data, pagination, isLoading } = useGroups();
  const groups: Group[] = !isLoading ? data : [];

  const groupTableData = useMemo(() => {
    return groups.map((group) => {
      const requirements = getRequirements(group);
      const areReqsMet = areRequirementsMet(group, requirements, profile);
      return {
        ...group,
        platform: group.groupSettings.platform,
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
      filters={filters}
      handleColumnFiltersChange={handleColumnFiltersChange}
      handleClearFilters={handleClearFilters}
    />
  );
}

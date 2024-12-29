import tasks from "@/assets/tasks.json";
import { DataTable } from "./ui/data-table";
import { columns } from "./GroupTable.Columns";

// import { Group } from "@/types";

// interface GroupTableProps {
//   groups: Group[];
// }

export function GroupTable() {
  return (
    <>
      <div className="hidden h-full flex-1 flex-col space-y-8 p-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Welcome back!</h2>
            <p className="text-muted-foreground">
              Here&apos;s a list of your tasks for this month!
            </p>
          </div>
        </div>
        <DataTable data={tasks} columns={columns} />
      </div>
    </>
  );
}

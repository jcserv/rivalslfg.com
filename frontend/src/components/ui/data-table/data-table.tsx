import * as React from "react";
import {
  ColumnDef,
  ColumnFiltersState,
  SortingState,
  VisibilityState,
  flexRender,
  getCoreRowModel,
  getFacetedRowModel,
  getFacetedUniqueValues,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui";

import { DataTablePagination } from "./data-table-pagination";
import { DataTableToolbar } from "./data-table-toolbar";
import { PaginationState } from "@/types";

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  pagination?: PaginationState;
  isLoading?: boolean;
}

export function DataTable<TData, TValue>({
  columns,
  data,
  pagination,
  isLoading,
}: DataTableProps<TData, TValue>) {
  const [rowSelection, setRowSelection] = React.useState({});
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    [],
  );

  React.useEffect(() => {
    if (pagination?.setFilters) {
      pagination.setFilters(columnFilters);
    }
  }, [columnFilters]);

  const [sorting, setSorting] = React.useState<SortingState>([]);
  const table = useReactTable({
    data,
    columns,
    state: {
      sorting,
      columnVisibility: {
        open: false,
      },
      rowSelection,
      columnFilters,
      ...(pagination
        ? {
            pageSize: pagination.pageSize,
            pageIndex: pagination.pageIndex,
          }
        : {}),
    },
    enableRowSelection: true,
    onRowSelectionChange: setRowSelection,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFacetedRowModel: getFacetedRowModel(),
    getFacetedUniqueValues: getFacetedUniqueValues(),
    ...(pagination
      ? {
          rowCount: pagination.totalCount ?? 0,
          getPaginationRowModel: getPaginationRowModel(),
          manualPagination: true,
          manualFiltering: true,
        }
      : {
          getPaginationRowModel: getPaginationRowModel(),
          manualPagination: false,
          manualFiltering: false,
        }),
  });

  // Trigger refetch when filters change
  React.useEffect(() => {
    if (pagination?.refetch) {
      pagination.refetch();
    }
  }, [columnFilters]);

  return (
    <div className="space-y-4">
      <DataTableToolbar
        table={table}
        rightAdornment={
          <DataTablePagination
            table={table}
            pagination={pagination}
            showCompactMode
          />
        }
      />
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id} colSpan={header.colSpan}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext(),
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id} isLoading={isLoading}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext(),
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : isLoading ? (
              // Render dummy rows when data is loading
              Array.from({ length: pagination?.pageSize || 10 }, (_, index) => (
                <TableRow key={`loading-${index}`} isLoading>
                  {columns.map((column) => (
                    <TableCell key={`col-${column.id}`}></TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <DataTablePagination table={table} pagination={pagination} />
    </div>
  );
}

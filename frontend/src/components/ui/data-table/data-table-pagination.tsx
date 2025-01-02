import { Table } from "@tanstack/react-table";
import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { PaginationState } from "@/types";

interface DataTablePaginationProps<TData> {
  table: Table<TData>;
  pagination?: PaginationState;
}

export function DataTablePagination<TData>({
  table,
  pagination,
}: DataTablePaginationProps<TData>) {
  const totalRows = pagination
    ? pagination.totalCount
    : table.getRowModel().rows.length;
  const pageSize = pagination
    ? pagination.pageSize
    : table.getState().pagination.pageSize;
  const pageIndex = pagination
    ? pagination.pageIndex
    : table.getState().pagination.pageIndex;
  const totalPages = pagination ? pagination.pageCount : table.getPageCount();

  const onPageSizeChange = (value: string) => {
    if (pagination) {
      pagination.setPageSize(Number(value));
      return;
    }
    table.setPageSize(Number(value));
  };

  const canPreviousPage = pagination
    ? pagination.canPreviousPage
    : table.getCanPreviousPage();
  const canNextPage = pagination
    ? pagination.canNextPage
    : table.getCanNextPage();

  const firstPage = () => {
    if (pagination) {
      pagination.firstPage();
      return;
    }
    table.setPageIndex(0);
  };

  const previousPage = () => {
    if (pagination) {
      pagination.previousPage();
      return;
    }
    table.previousPage();
  };

  const nextPage = () => {
    if (pagination) {
      pagination.nextPage();
      return;
    }
    table.nextPage();
  };

  const lastPage = () => {
    if (pagination) {
      pagination.lastPage();
      return;
    }
    table.setPageIndex(table.getPageCount() - 1);
  };

  return (
    <div className="flex items-center justify-between px-2">
      <div className="flex-1 text-sm text-muted-foreground">
        {totalRows} row(s) total
      </div>
      <div className="flex items-center space-x-6 lg:space-x-8">
        <div className="flex items-center space-x-2">
          <p className="text-sm font-medium">Rows per page</p>
          <Select value={`${pageSize}`} onValueChange={onPageSizeChange}>
            <SelectTrigger className="h-8 w-[70px]">
              <SelectValue placeholder={pageSize} />
            </SelectTrigger>
            <SelectContent side="top">
              {[10, 20, 30, 40, 50].map((pageSize) => (
                <SelectItem key={pageSize} value={`${pageSize}`}>
                  {pageSize}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex w-[100px] items-center justify-center text-sm font-medium">
          Page {pageIndex + 1} of {totalPages}
        </div>
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            className="hidden h-8 w-8 p-0 lg:flex"
            onClick={firstPage}
            disabled={!canPreviousPage}
          >
            <span className="sr-only">Go to first page</span>
            <ChevronsLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={previousPage}
            disabled={!canPreviousPage}
          >
            <span className="sr-only">Go to previous page</span>
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={nextPage}
            disabled={!canNextPage}
          >
            <span className="sr-only">Go to next page</span>
            <ChevronRight className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="hidden h-8 w-8 p-0 lg:flex"
            onClick={lastPage}
            disabled={!canNextPage}
          >
            <span className="sr-only">Go to last page</span>
            <ChevronsRight className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}

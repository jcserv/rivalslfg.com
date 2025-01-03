import { ColumnFilter } from "@tanstack/react-table";

import { QueryParams } from "./query";

export type OffsetPagination = {
  limit: number;
  offset: number;
  count?: boolean;
};

export type PaginatedQueryFn<TData> = (
  params: QueryParams,
) => PaginatedQueryFnResponse<TData>;

export type PaginatedQueryFnResponse<TData> = Promise<{
  data: TData[];
  totalCount: number;
}>;

export interface PaginationState {
  pageSize: number;
  pageIndex: number;
  pageCount: number;
  totalCount: number | null;
  setPageSize: (size: number) => void;
  setPageIndex: (index: number) => void;
  setFilters: (filters: ColumnFilter[]) => void;
  canPreviousPage: boolean;
  canNextPage: boolean;
  previousPage: () => void;
  nextPage: () => void;
  firstPage: () => void;
  lastPage: () => void;
  refetch: () => void;
}

export type PaginatedQueryFn<TData> = (params: {
  limit: number;
  offset: number;
  count: boolean;
}) => PaginatedQueryFnResponse<TData>;

export type PaginatedQueryFnResponse<TData> = Promise<{
  data: TData[];
  totalCount: number;
}>;

export interface PaginationState {
  pageSize: number;
  pageIndex: number;
  pageCount: number;
  totalCount: number;
  setPageSize: (size: number) => void;
  setPageIndex: (index: number) => void;
  canPreviousPage: boolean;
  canNextPage: boolean;
  previousPage: () => void;
  nextPage: () => void;
  firstPage: () => void;
  lastPage: () => void;
}

import { useCallback, useEffect, useState } from "react";

import { useQuery } from "@tanstack/react-query";
import { ColumnFiltersState } from "@tanstack/react-table";

import { queryClient } from "@/routes/__root";
import { getFilterBy } from "@/types";
import { PaginatedQueryFn } from "@/types/paginate";

interface PaginationState {
  pageSize: number;
  pageIndex: number;
  filters?: ColumnFiltersState;
}

interface PaginationOptions<TData> {
  queryKey: readonly string[];
  queryFn: PaginatedQueryFn<TData>;
  initialState?: Partial<PaginationState>;
}

const getQueryKey = (
  base: readonly string[],
  { pageSize, pageIndex, filters }: PaginationState,
) => [
  ...base,
  pageSize,
  pageIndex,
  ...(filters?.map((f) => [f.id, f.value]) ?? []),
];

export function usePagination<TData>({
  queryKey: baseQueryKey,
  queryFn,
  initialState,
}: PaginationOptions<TData>) {
  const [pageSize, setPageSize] = useState(initialState?.pageSize || 10);
  const [pageIndex, setPageIndex] = useState(initialState?.pageIndex || 0);
  const [filters, setFilters] = useState<ColumnFiltersState>([]);
  const [totalCount, setTotalCount] = useState<number | null>(null);

  // Current page query
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: getQueryKey(baseQueryKey, { pageSize, pageIndex, filters }),
    queryFn: async () => {
      const response = await queryFn({
        paginateBy: {
          limit: pageSize,
          offset: pageIndex * pageSize,
          count: totalCount === null, // Only request total count on first load
        },
        filterBy: getFilterBy(filters),
      });

      if (totalCount === null && response.totalCount) {
        setTotalCount(response.totalCount);
      }

      return response;
    },
  });

  // Prefetch page based on page offset
  const prefetchPage = useCallback(
    async (pageOffset: number) => {
      const targetPageIndex = pageIndex + pageOffset;
      const nextPageKey = getQueryKey(baseQueryKey, {
        pageSize,
        pageIndex: targetPageIndex,
      });

      await queryClient.prefetchQuery({
        queryKey: nextPageKey,
        queryFn: async () =>
          queryFn({
            paginateBy: {
              limit: pageSize,
              offset: targetPageIndex * pageSize,
              count: false,
            },
          }),
        staleTime: 60000,
      });
    },
    [pageSize, pageIndex, queryClient, baseQueryKey, queryFn],
  );

  // Automatically prefetch prev and next page whenever current page changes
  useEffect(() => {
    const pageCount = Math.ceil((totalCount || 0) / pageSize);
    if (pageIndex > 0) {
      prefetchPage(-1);
    }
    if (pageIndex < pageCount - 1) {
      prefetchPage(1);
    }
  }, [pageIndex, pageSize, totalCount, prefetchPage]);

  const pageCount = Math.ceil((totalCount || 0) / pageSize);

  // Prefetch last page when total count is first set
  useEffect(() => {
    if (totalCount !== null && pageCount > 1) {
      const lastPageOffset = pageCount - 1 - pageIndex;
      prefetchPage(lastPageOffset);
    }
  }, [totalCount, pageCount, pageIndex, prefetchPage]);

  return {
    data: data?.data || [],
    pagination: {
      pageSize,
      pageIndex,
      pageCount,
      totalCount,
      setPageSize: (size: number) => {
        setPageSize(size);
        setPageIndex(0);
      },
      setPageIndex,
      setFilters,
      canPreviousPage: pageIndex > 0,
      canNextPage: pageIndex < pageCount - 1,
      previousPage: () => setPageIndex((old) => Math.max(0, old - 1)),
      nextPage: () => setPageIndex((old) => Math.min(old + 1, pageCount - 1)),
      firstPage: () => setPageIndex(0),
      lastPage: () => setPageIndex(Math.max(0, pageCount - 1)),
      prefetchPage,
      refetch,
    },
    isLoading,
    error,
    refetch,
  };
}

import { useState } from "react";

import { useQuery } from "@tanstack/react-query";

import { PaginatedQueryFn } from "@/types/paginate";

interface PaginationState {
  pageSize: number;
  pageIndex: number;
}

interface PaginationOptions<TData> {
  queryKey: readonly string[];
  queryFn: PaginatedQueryFn<TData>;
  initialState?: Partial<PaginationState>;
}

export function usePagination<TData>({
  queryKey,
  queryFn,
  initialState,
}: PaginationOptions<TData>) {
  const [pageSize, setPageSize] = useState(initialState?.pageSize || 10);
  const [pageIndex, setPageIndex] = useState(initialState?.pageIndex || 0);
  const [totalCount, setTotalCount] = useState<number | null>(null);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: [...queryKey, pageSize, pageIndex],
    queryFn: async () => {
      const response = await queryFn({
        limit: pageSize,
        offset: pageIndex * pageSize,
        count: totalCount === null, // Only request count on first load
      });

      if (totalCount === null && response.totalCount) {
        setTotalCount(response.totalCount);
      }

      return response;
    },
  });

  const pageCount = Math.ceil((totalCount || 0) / pageSize);

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
      canPreviousPage: pageIndex > 0,
      canNextPage: pageIndex < pageCount - 1,
      previousPage: () => setPageIndex((old) => Math.max(0, old - 1)),
      nextPage: () => setPageIndex((old) => Math.min(old + 1, pageCount - 1)),
      firstPage: () => setPageIndex(0),
      lastPage: () => setPageIndex(Math.max(0, pageCount - 1)),
    },
    isLoading,
    error,
    refetch,
  };
}

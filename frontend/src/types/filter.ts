type FilterOp = "eq";

export interface Filter {
  field: string;
  op: FilterOp;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  value: any;
}

export function buildFilterQuery(filters: Filter[]): string {
  if (!filters || filters.length === 0) {
    return "";
  }

  const filterString = filters
    .map((filter) => {
      const value =
        typeof filter.value === "string" ? `"${filter.value}"` : filter.value;
      return `${filter.field} ${filter.op} ${value}`;
    })
    .join(" and ");

  return filterString;
}

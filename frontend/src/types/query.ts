import { buildFilterQuery, Filter } from "./filter";
import { OffsetPagination } from "./paginate";
import { Profile } from "./types";

export interface QueryParams {
  filterBy?: Filter[];
  paginateBy?: OffsetPagination;
  playerRequirements?: Profile;
}

export function toURLSearchParams(q: QueryParams): URLSearchParams {
  const params = new URLSearchParams();
  if (q.paginateBy) {
    params.set("limit", q.paginateBy.limit.toString());
    params.set("offset", q.paginateBy.offset.toString());
    params.set("count", `${q.paginateBy.count ?? false}`);
  }
  if (q.filterBy) {
    const filterStr = buildFilterQuery(q.filterBy);
    if (filterStr) {
      params.set("filter", filterStr);
    }
  }
  return params;
}

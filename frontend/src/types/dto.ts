import { Group } from "./types";

export interface PaginatedGroupsResponse {
  groups: Group[];
  pageCount: number;
  totalCount: number;
}

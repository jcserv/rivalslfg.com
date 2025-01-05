import { HTTPClient, StatusCode, StatusCodes } from "@/api";
import {
  CreateGroupResponse,
  getCreateGroupFromProfile,
  Group,
  PaginatedGroupsResponse,
  Profile,
  QueryParams,
  toURLSearchParams,
} from "@/types";

export class RivalsLFGClient extends HTTPClient {
  private readonly baseURL: string;

  constructor(baseURL: string) {
    super();
    this.baseURL = baseURL;
  }

  async createGroup(owner: Profile): Promise<CreateGroupResponse> {
    const body = getCreateGroupFromProfile(owner);
    const response = await this.fetchWithAuth(`${this.baseURL}/api/v1/groups`, {
      method: "POST",
      body: JSON.stringify(body),
    });
    if (!response.ok) {
      throw new Error(`HTTP error, status: ${response.status}`);
    }

    const data = await response.json();
    return data ?? { groupId: "", playerId: 0 };
  }

  async getGroups(query?: QueryParams): Promise<PaginatedGroupsResponse> {
    try {
      const params = query ? toURLSearchParams(query) : new URLSearchParams();
      const response = await this.fetchWithRetry(
        `${this.baseURL}/api/v1/groups?${params.toString()}`,
      );
      const totalCount = parseInt(response.headers.get("X-Total-Count") ?? "0");
      const data = await response.json();
      return {
        groups: data ?? [],
        pageCount: Math.ceil(totalCount / (query?.paginateBy?.limit || 10)),
        totalCount,
      };
    } catch {
      return {
        groups: [],
        pageCount: 0,
        totalCount: 0,
      };
    }
  }

  async getGroup(id: string): Promise<Group | undefined> {
    try {
      const response = await this.fetchWithRetry(
        `${this.baseURL}/api/v1/groups/${id}`,
      );
      const data = await response.json();
      return data;
    } catch {
      return undefined;
    }
  }

  async joinGroup(
    groupId: string,
    player: Profile,
    passcode: string,
  ): Promise<StatusCode> {
    try {
      const response = await this.fetchWithAuth(
        `${this.baseURL}/api/v1/groups/${groupId}/players/${player.id}`,
        {
          method: "POST",
          body: JSON.stringify({
            player,
            passcode,
          }),
        },
      );
      return response.status as StatusCode;
    } catch {
      return StatusCodes.InternalServerError as StatusCode;
    }
  }

  async removePlayer(
    groupId: string,
    playerToRemoveId: number,
  ): Promise<StatusCode> {
    try {
      const response = await this.fetchWithAuth(
        `${this.baseURL}/api/v1/groups/${groupId}/players/${playerToRemoveId}`,
        {
          method: "DELETE",
        },
      );
      return response.status as StatusCode;
    } catch {
      return StatusCodes.InternalServerError as StatusCode;
    }
  }
}

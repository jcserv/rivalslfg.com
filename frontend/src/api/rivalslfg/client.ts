import { HTTPClient, StatusCode, StatusCodes } from "@/api";
import {
  getGroupFromProfile,
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

  async upsertGroup(owner: Profile, id: string = ""): Promise<string> {
    const newGroup = getGroupFromProfile(owner, id);
    try {
      const response = await fetch(`${this.baseURL}/api/v1/groups`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newGroup),
      });

      if (!response.ok) {
        throw new Error(`HTTP error, status: ${response.status}`);
      }

      const data = await response.json();
      return data.id || "";
    } catch {
      return "";
    }
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
      const response = await fetch(
        `${this.baseURL}/api/v1/groups/${groupId}/players/${player.id}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
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
    playerId: number,
    requesterName: string,
    playerName: string,
  ): Promise<StatusCode> {
    try {
      const response = await fetch(
        `${this.baseURL}/api/v1/groups/${groupId}/players/${playerId}`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            requesterName,
            playerName,
          }),
        },
      );
      return response.status as StatusCode;
    } catch {
      return StatusCodes.InternalServerError as StatusCode;
    }
  }
}

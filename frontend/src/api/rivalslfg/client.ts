import { HTTPClient, StatusCode } from "@/api";
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
    const response = await this.fetchWithRetry(
      `${this.baseURL}/api/v1/groups/${id}`,
    );
    const data = await response.json();
    return data;
  }

  async joinGroup(
    groupId: string,
    player: Profile,
    passcode: string,
  ): Promise<StatusCode> {
    const response = await this.fetchWithAuth(
      `${this.baseURL}/api/v1/groups/${groupId}/players/${player.id ?? 0}`,
      {
        method: "POST",
        body: JSON.stringify({
          name: player.name,
          passcode,
          gamemode: player.gamemode,
          region: player.region,
          platform: player.platform,
          role: player.role,
          rankId: player.rank,
          characters: player.characters,
          voiceChat: player.voiceChat,
          mic: player.mic,
          vanguards: player.roleQueue?.vanguards ?? 0,
          duelists: player.roleQueue?.duelists ?? 0,
          strategists: player.roleQueue?.strategists ?? 0,
        }),
      },
    );
    return response.status as StatusCode;
  }

  async removePlayer(groupId: string, playerId: number): Promise<StatusCode> {
    const response = await this.fetchWithAuth(
      `${this.baseURL}/api/v1/groups/${groupId}/players/${playerId}`,
      {
        method: "DELETE",
      },
    );
    return response.status as StatusCode;
  }
}

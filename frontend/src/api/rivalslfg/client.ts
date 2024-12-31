import { HTTPClient } from "@/api/base/client";
import { Group, Profile } from "@/types";

export class RivalsLFGClient extends HTTPClient {
  private readonly baseURL: string;

  constructor(baseURL: string) {
    super();
    this.baseURL = baseURL;
  }

  async getGroups(): Promise<Group[]> {
    const response = await this.fetchWithRetry(`${this.baseURL}/api/v1/groups`);
    const data = await response.json();
    return data;
  }

  async getGroup(id: string): Promise<Group> {
    const response = await this.fetchWithRetry(
      `${this.baseURL}/api/v1/groups/${id}`,
    );
    const data = await response.json();
    return data;
  }

  async createPlayer(
    profile: Profile,
    id: number | undefined,
  ): Promise<string> {
    const response = await fetch(
      `${this.baseURL}/api/v1/players${id ? `/${id}` : ""}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          displayName: profile.name,
          region: profile.region,
          platform: profile.platform,
          gamemode: profile.gamemode,
          roles: profile.roles,
          rank: profile.rank,
          characters: profile.characters,
          voiceChat: profile.voiceChat,
          mic: profile.mic,
          roleQueue: {
            vanguards: profile.roleQueue?.vanguards ?? 0,
            duelists: profile.roleQueue?.duelists ?? 0,
            strategists: profile.roleQueue?.strategists ?? 0,
          },
          groupSettings: {
            platforms: profile.groupSettings?.platforms ?? [],
            voiceChat: profile.groupSettings?.voiceChat ?? false,
            mic: profile.groupSettings?.mic ?? false,
          },
        }),
      },
    );
    return response.json().then((data) => data.id);
  }

  async getPlayer(id: string): Promise<Profile> {
    const response = await this.fetchWithRetry(
      `${this.baseURL}/api/v1/players/${id}`,
    );
    const data = await response.json();
    return data;
  }
}

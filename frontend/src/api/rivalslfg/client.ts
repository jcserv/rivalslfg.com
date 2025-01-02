import { HTTPClient } from "@/api/base/client";
import { getGroupFromProfile, Group, Profile } from "@/types";
import { StatusCode, StatusCodes } from "@/types/http";

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
    } catch (error) {
      console.error("Error creating/updating group:", error);
      return "";
    }
  }

  async getGroups(): Promise<Group[]> {
    try {
      const response = await this.fetchWithRetry(
        `${this.baseURL}/api/v1/groups`,
      );
      const data = await response.json();
      return data;
    } catch (error) {
      console.error("Error fetching groups", error);
      return [];
    }
  }

  async getGroup(id: string): Promise<Group | undefined> {
    try {
      const response = await this.fetchWithRetry(
        `${this.baseURL}/api/v1/groups/${id}`,
      );
      const data = await response.json();
      return data;
    } catch (error) {
      console.error("Error fetching group with id", error);
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
    } catch (error) {
      console.error("Error joining group", error);
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
    } catch (error) {
      console.error("Error removing player from group", error);
      return StatusCodes.InternalServerError as StatusCode;
    }
  }
}

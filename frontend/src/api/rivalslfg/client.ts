import { HTTPClient } from "@/api/base/client";
import { Group } from "@/types";

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
}

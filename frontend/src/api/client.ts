export class HTTPClient {
  private lastRequestTime: number = 0;
  private readonly minRequestInterval: number = 75; // 75ms between requests

  private async delay(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  private async waitForRateLimit(): Promise<void> {
    const now = Date.now();
    const timeSinceLastRequest = now - this.lastRequestTime;

    if (timeSinceLastRequest < this.minRequestInterval) {
      await this.delay(this.minRequestInterval - timeSinceLastRequest);
    }

    this.lastRequestTime = Date.now();
  }

  public async fetchWithRetry(
    url: string,
    retries: number = 3,
  ): Promise<Response> {
    await this.waitForRateLimit();

    const headers = {
      Accept: "application/json;q=0.9,*/*;q=0.8",
    };

    for (let attempt = 0; attempt < retries; attempt++) {
      try {
        const response = await fetch(url, { headers });

        if (response.status === 429) {
          const retryAfter = response.headers.get("Retry-After");
          const delayMs = retryAfter ? parseInt(retryAfter) * 1000 : 1000;
          await this.delay(delayMs);
          continue;
        }

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        return response;
      } catch (error) {
        if (attempt === retries - 1) throw error;
        await this.delay(1000 * (attempt + 1)); // Exponential backoff
      }
    }

    throw new Error("Max retries exceeded");
  }
}

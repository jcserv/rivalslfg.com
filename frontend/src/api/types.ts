/* eslint-disable @typescript-eslint/no-explicit-any */
export type StatusCode = number;

export const StatusCodes: Record<string, StatusCode> = {
  OK: 200,
  Created: 201,
  Accepted: 202,
  NoContent: 204,
  BadRequest: 400,
  Unauthorized: 401,
  Forbidden: 403,
  NotFound: 404,
  InternalServerError: 500,
};

export type PaginationParams = {
  limit: number;
  offset: number;
};

export function toString(p: PaginationParams): string {
  return `${p.limit}-${p.offset}`;
}

export class HTTPError extends Error {
  public readonly statusCode: StatusCode;
  public readonly statusText: string;
  public readonly details: Map<string, any>;

  constructor(statusCode: StatusCode, statusText: string, respBody: any) {
    super(respBody.message);
    this.statusCode = statusCode;
    this.statusText = statusText;
    this.details = respBody.details;
  }
}

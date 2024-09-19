let API_URL: string | undefined = undefined;
const apiPath = "/api/v1";

export function GetApiUrl(): string {
  if (API_URL === undefined) {
    const BACKEND_HOST = import.meta.env.VITE_BACKEND_HOST;
    const BACKEND_USE_TLS = import.meta.env.VITE_BACKEND_USE_TLS;

    if (BACKEND_HOST !== undefined) {
      if (BACKEND_USE_TLS !== undefined && BACKEND_USE_TLS === "true") {
        API_URL = `https://${BACKEND_HOST}${apiPath}`;
      } else {
        API_URL = `http://${BACKEND_HOST}${apiPath}`;
      }
    } else {
      if (window.location.protocol === "https:") {
        API_URL = `https://${window.location.host}${apiPath}`;
      } else {
        API_URL = `http://${window.location.host}${apiPath}`;
      }
    }
  }

  return API_URL;
}

export class ErrorResponseError extends Error {
  public errorResponse: ErrorResponse;

  constructor(errorResponse: ErrorResponse) {
    super(`Request failed with status code ${errorResponse.status}`);
    this.errorResponse = errorResponse;
  }
}

export interface ErrorResponse {
  status: number;
  message: string;
  request_id: string;
}

export async function SendWithoutResponse(
  path: string,
  requestInit?: RequestInit,
): Promise<void> {
  try {
    const response = await fetch(`${GetApiUrl()}${path}`, requestInit);
    if (!response.ok) {
      const errorResponse = (await response.json()) as ErrorResponse;
      if (isValidErrorResponse(errorResponse)) {
        throw new ErrorResponseError(errorResponse);
      }
      throw new Error(
        `Response failed with status code ${response.status} but no ErrorResponse returned`,
      );
    }
  } catch (e) {
    console.error(`request to path "${path}" failed with error ${e}`);
    throw e;
  }
}

export async function SendBase(
  path: string,
  requestInit?: RequestInit,
): Promise<Response> {
  const response = await fetch(`${GetApiUrl()}${path}`, requestInit);
  if (!response.ok) {
    const errorResponse = (await response.json()) as ErrorResponse;
    if (isValidErrorResponse(errorResponse)) {
      throw new ErrorResponseError(errorResponse);
    }
    throw new Error(
      `Response failed with status code ${response.status} but no ErrorResponse returned`,
    );
  }

  return response;
}

export async function Send<T>(
  path: string,
  requestInit?: RequestInit,
): Promise<T> {
  try {
    const response = await SendBase(path, requestInit);
    const respJson = (await response.json()) as T;
    if (respJson === undefined) {
      throw new Error("response json in undefined");
    }

    return respJson;
  } catch (e) {
    console.error(`request to path "${path}" failed with error ${e}`);
    throw e;
  }
}

export async function Get<T>(path: string): Promise<T> {
  return await Send<T>(path);
}

export async function PostWithoutResponse(path: string): Promise<void> {
  return await SendWithoutResponse(path, { method: "POST" });
}

export async function Post<T>(path: string): Promise<T> {
  return await Send<T>(path, { method: "POST" });
}

export async function PostJson<T>(path: string, body: any): Promise<T> {
  return await Send<T>(path, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export async function PostJsonWithoutResponse(
  path: string,
  body: any,
): Promise<void> {
  return await SendWithoutResponse(path, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export async function DeleteWithoutResponse(path: string): Promise<void> {
  return await SendWithoutResponse(path, { method: "DELETE" });
}

function isValidErrorResponse(errorResponse: any): boolean {
  if (typeof errorResponse !== "object") {
    return false;
  }
  if (errorResponse === null) {
    return false;
  }

  const status = errorResponse.status;
  if (typeof status !== "number" || status === 0) {
    return false;
  }

  const message = errorResponse.message;
  if (typeof message !== "string" || message === "") {
    return false;
  }

  const requestId = errorResponse.request_id;
  if (typeof requestId !== "string" || requestId === "") {
    return false;
  }

  return true;
}

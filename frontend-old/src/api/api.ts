import { API_URL, WEBSOCKET_URL } from '@/main'

export * from '@/api/logs'
export * from '@/api/start'
export * from '@/api/status'

export class ErrorResponse {
    status: string | undefined
    message: string | undefined
    request_id: string | undefined
}

export async function SendWithoutResponse(path: string, requestInit?: RequestInit): Promise<undefined | ErrorResponse> {
    const response = await fetch(`${API_URL}${path}`, requestInit)
    if (!response.ok) {
        const errorResponse = await response.json() as ErrorResponse
        checkIfAllValuesAreDefined(Object.keys(new ErrorResponse()) as (keyof ErrorResponse)[], errorResponse)
        return errorResponse
    }

    return undefined
}

export async function GetWithoutResponse(path: string): Promise<undefined | ErrorResponse> {
    return await SendWithoutResponse(path)
}

export async function PostWithoutResponse(path: string): Promise<undefined | ErrorResponse> {
    return await SendWithoutResponse(path, { method: "POST" })
}

export async function Send<T>(path: string, requestInit?: RequestInit): Promise<T | ErrorResponse> {
    const response = await fetch(`${API_URL}${path}`, requestInit)
    if (!response.ok) {
        const errorResponse = await response.json() as ErrorResponse
        checkIfAllValuesAreDefined(Object.keys(new ErrorResponse()) as (keyof ErrorResponse)[], errorResponse)
        return errorResponse
    }

    return await response.json() as T
}

export async function Get<T>(path: string): Promise<T | ErrorResponse> {
    return await Send<T>(path)
}

export async function Post<T>(path: string): Promise<T | ErrorResponse> {
    return await Send<T>(path, { method: "POST" })
}

export function checkIfAllValuesAreDefined<T>(keys: (keyof T)[], obj: T) {
    for (const key of keys) {
        if (obj[key] === undefined) {
            throw new Error(`Property ${String(key)} is undefined`);
        }
    }
}

///

export class WebSocketMessage {
    type: string | undefined
    message: any
}

export function ConnectToWebSocket() {
    return new WebSocket(WEBSOCKET_URL)
}
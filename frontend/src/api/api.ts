import {Show} from "@/errorDisplay";

let API_URL: string
let WEBSOCKET_URL: string

function isValidHostname(hostname: string): boolean {
    const hostnameRegex = new RegExp("^[a-zA-Z0-9-.]+$")
    const hostnamePortRegex = new RegExp("^[a-zA-Z0-9-.]+:\\d{1,5}$")
    return hostnameRegex.test(hostname) || hostnamePortRegex.test(hostname);
}

export function Setup(backendAddress: string | undefined, useTls: boolean | undefined) {
    const apiPath = "/api/v1"

    if (backendAddress == null || backendAddress == "" || !isValidHostname(backendAddress)) {
        if (window.location.protocol === "https") {
            API_URL = `https://${window.location.host}${apiPath}`
            WEBSOCKET_URL = `wss://${window.location.host}/ws`
        } else {
            API_URL = `http://${window.location.host}${apiPath}`
            WEBSOCKET_URL = `ws://${window.location.host}/ws`
        }
    } else {
        if (useTls == true) {
            API_URL = `https://${backendAddress}${apiPath}`
            WEBSOCKET_URL = `wss://${backendAddress}${apiPath}/ws`
        } else {
            API_URL = `http://${backendAddress}${apiPath}`
            WEBSOCKET_URL = `ws://${backendAddress}${apiPath}/ws`
        }
    }

    console.log("API_URL: ", API_URL)
    console.log("WEBSOCKET_URL: ", WEBSOCKET_URL)
}

export class ErrorResponse {
    status: string | undefined
    message: string | undefined
    request_id: string | undefined
}

export async function SendWithoutResponse(path: string, requestInit?: RequestInit): Promise<undefined | ErrorResponse> {
    const response = await fetch(`${API_URL}${path}`, requestInit)
    if (!response.ok) {
        const errorResponse = await response.json() as ErrorResponse
        throwIfErrorResponseIsNotValid(errorResponse)
        return errorResponse
    }

    return undefined
}

export async function PostWithoutResponse(path: string): Promise<undefined | ErrorResponse> {
    return await SendWithoutResponse(path, {method: "POST"})
}

export async function Send<T>(path: string, requestInit?: RequestInit): Promise<T | ErrorResponse> {
    try {
        const response = await fetch(`${API_URL}${path}`, requestInit)
        if (!response.ok) {
            const errorResponse = await response.json() as ErrorResponse
            throwIfErrorResponseIsNotValid(errorResponse)
            return errorResponse
        }

        const respJson = await response.json() as T
        if (respJson === undefined) {
            throw new Error("response json in undefined")
        }

        return respJson
    } catch (e) {
        console.error(`request to path "${path}" failed with error ${e}`);
        throw e
    }
}

export async function Get<T>(path: string): Promise<T | ErrorResponse> {
    return await Send<T>(path)
}

export async function Post<T>(path: string): Promise<T | ErrorResponse> {
    return await Send<T>(path, {method: "POST"})
}

export async function PostJson<T>(path: string, body: any): Promise<T | ErrorResponse> {
    return await Send<T>(path, {
        method: "POST",
        body: JSON.stringify(body),
    })
}

function throwIfErrorResponseIsNotValid(errorResponse: ErrorResponse) {
    const keys = Object.keys(new ErrorResponse()) as (keyof ErrorResponse)[]
    for (const key of keys) {
        if (errorResponse[key] == null || errorResponse[key] === "") {
            throw new Error(`Property ${String(key)} is undefined`);
        }
    }
}

function isValidErrorResponse(errorResponse: any): boolean {
    if (errorResponse == null) {
        return false
    }

    const keys = Object.keys(new ErrorResponse()) as (keyof ErrorResponse)[]
    for (const key of keys) {
        if (errorResponse[key] == null || errorResponse[key] === "") {
            return false
        }
    }

    return true
}

export function handleErrorResponse(title: string, response: any) {
    if (isValidErrorResponse(response)) {
        console.log(`is valid: `, response)
        ShowErrorResponse(title, response)
        throw new Error(`${title}. Response: ${response}`)
    }
}

export function handleErrorResponseWithMessage(title: string, message: string, response: any) {
    if (isValidErrorResponse(response)) {
        ShowErrorResponseWithMessage(title, message, response)
        throw new Error(`${title}. Message: ${message} Response: ${response}`)
    }
}

export function ShowErrorResponse(title: string, errorResponse: ErrorResponse) {
    Show(title, `${errorResponse.message}`, `RequestId: ${errorResponse.request_id}`)
}

export function ShowErrorResponseWithMessage(title: string, message: string, errorResponse: ErrorResponse) {
    Show(title, message, `${errorResponse.message}`, `RequestId: ${errorResponse.request_id}`)
}

///

export class WebSocketMessage {
    type: string | undefined
    message: any
}

export function ConnectToWebSocket() {
    return new WebSocket(WEBSOCKET_URL)
}
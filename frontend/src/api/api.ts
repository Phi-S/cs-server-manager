import {Show} from "@/errorDisplay";

let API_URL: string

export function SetAPIUrl(ApiUrl: string) {
    API_URL = ApiUrl
}

export interface ErrorResponse {
    status: string
    message: string
    request_id: string
}

export async function SendWithoutResponse(path: string, requestInit?: RequestInit): Promise<undefined | ErrorResponse> {
    try {
        const response = await fetch(`${API_URL}${path}`, requestInit)
        if (!response.ok) {
            const errorResponse = await response.json() as ErrorResponse
            if (!isValidErrorResponse(errorResponse)) {
                throw new Error(`Response failed with status code ${response.status} but no ErrorResponse returned`)
            }
            return errorResponse
        }

        return undefined
    } catch (e) {
        console.error(`request to path "${path}" failed with error ${e}`);
        throw e
    }
}

export async function Send<T>(path: string, requestInit?: RequestInit): Promise<T | ErrorResponse> {
    try {
        const response = await fetch(`${API_URL}${path}`, requestInit)
        if (!response.ok) {
            const errorResponse = await response.json() as ErrorResponse
            if (!isValidErrorResponse(errorResponse)) {
                throw new Error(`Response failed with status code ${response.status} but no ErrorResponse returned`)
            }
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

export async function PostWithoutResponse(path: string): Promise<undefined | ErrorResponse> {
    return await SendWithoutResponse(path, {method: "POST"})
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
    if (errorResponse.status == null || errorResponse.status === "") {
        throw new Error(`Property ${String("status")} is undefined`);
    }

}

function isValidErrorResponse(errorResponse: any): boolean {
    if (typeof errorResponse !== "object") {
        return false
    }
    if (errorResponse === null) {
        return false
    }

    const status = errorResponse.status
    if (typeof status !== "string" || status === "") {
        return false
    }

    const message = errorResponse.message
    if (typeof message !== "string" || message === "") {
        return false
    }

    const requestId = errorResponse.request_id
    if (typeof requestId !== "string" || requestId === "") {
        return false
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
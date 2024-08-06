import { ErrorResponse, Start } from '@/api/api';
import { SendCommand } from './api/send-command';
import { Stop } from './api/stop';
import { CancelUpdate, Update } from './api/update';
import { Show } from './popup';

export function ShowErrorResponse(title: string, errorResponse: ErrorResponse) {
    Show(title, `Error: ${errorResponse.message}`, `RequestId: ${errorResponse.request_id}`)
}

export function ShowErrorResponseWithMessage(title: string, message: string, errorResponse: ErrorResponse) {
    Show(title, message, `Error: ${errorResponse.message}`, `RequestId: ${errorResponse.request_id}`)
}

export function isErrorResponse(response: any): response is ErrorResponse {
    return response !== undefined
        && response !== null
        && (response as ErrorResponse).status !== undefined
        && (response as ErrorResponse).message !== undefined
        && (response as ErrorResponse).request_id !== undefined
}

export async function StartHandler() {
    try {
        const resp = await Start()
        if (isErrorResponse(resp)) {
            ShowErrorResponse("Failed to start server", resp)
        }
    } catch (error) {
        Show("Failed to start server", `${error}`)
    }
}


export async function StopHandler() {
    try {
        const resp = await Stop()
        if (isErrorResponse(resp)) {
            ShowErrorResponse("Failed to stop server", resp)
        }
    } catch (error) {
        Show("Failed to stop server", `${error}`)
    }
}

export async function RestartHandler() {

    try {
        const stopResp = await Stop()
        if (isErrorResponse(stopResp)) {
            ShowErrorResponseWithMessage("Failed to restart server", "Server failed to stop", stopResp)
        }

        const startResp = await Start()
        if (isErrorResponse(startResp)) {
            ShowErrorResponse("Failed to restart server", startResp)
        }
    } catch (error) {
        Show("Failed to restart server", `${error}`)
    }
}


export async function UpdateHandler() {

    try {
        const stopResp = await Stop()
        if (isErrorResponse(stopResp)) {
            ShowErrorResponse("Failed to update server. Server failed to stop", stopResp)
        }
        const updateResp = await Update()
        if (isErrorResponse(updateResp)) {
            ShowErrorResponse("Failed to update server", updateResp)
        }
    } catch (error) {
        Show("Failed to update server", `${error}`)
    }
}


export async function CancelUpdateHandler() {

    try {
        const resp = await CancelUpdate()
        if (isErrorResponse(resp)) {
            ShowErrorResponse("Failed to cancel update server", resp)
        }
    } catch (error) {
        Show("Failed to cancel update server", `${error}`)
    }
}

export async function SendCommandHandler(command: string) {

    try {
        const resp = await SendCommand(command)
        if (isErrorResponse(resp)) {
            ShowErrorResponse("Failed to send command", resp)
        }
    } catch (error) {
        Show("Failed to send command", `${error}`)
    }
}
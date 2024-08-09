import {
    ErrorResponse,
    Get,
    handleErrorResponse,
    handleErrorResponseWithMessage,
    Post,
    PostWithoutResponse,
    ShowErrorResponse
} from "@/api/api";

export class SendCommandResponse {
    output: string[]
}

export enum ServerStatus {
    ServerStatusStarted = "server-status-started",
    ServerStatusStarting = "server-status-starting",
    ServerStatusStopped = "server-status-stopped",
    ServerStatusStopping = "server-status-stopping"
}

export enum SteamcmdStatus {
    SteamcmdStatusStopped = "steamcmd-status-stopped",
    SteamcmdStatusUpdating = "steamcmd-status-updating"
}

export class Status {
    hostname: string
    server: ServerStatus
    steamcmd: SteamcmdStatus
    player_count: number
    max_player_count: number
    map: string
}

export class LogEntry {
    timestamp: string
    log_type: string
    message: string
}

export async function startServer() {
    const resp = await PostWithoutResponse("/start")
    handleErrorResponse("Failed to start server", resp)
}

export async function stopServer() {
    const resp = await PostWithoutResponse("/stop")
    handleErrorResponse("Failed to stop server", resp)
}

export async function restartServer() {
    let resp = await PostWithoutResponse("/stop")
    handleErrorResponseWithMessage("Failed to restart server", "Server failed to stop", resp)
    resp = await PostWithoutResponse("/start")
    handleErrorResponse("Failed to restart server", resp)
}

export async function sendCommand(command: string) {
    const resp = await Post<SendCommandResponse>(`/send-command?command=${command}`)
    handleErrorResponse("Failed to send command", resp)
}

export async function startUpdate() {
    const resp = await PostWithoutResponse("/update")
    handleErrorResponse("Failed to start server update", resp)
}

export async function cancelUpdate() {
    const resp = await PostWithoutResponse("/update/cancel")
    handleErrorResponse("Failed to cancel server update", resp)
}

export async function getStatus(): Promise<Status> {
    const resp = await Get<Status>("/status")
    handleErrorResponse("Failed to get server status", resp)
    return resp as Status
}

export async function getLogs(count: number): Promise<LogEntry[]> {
    const resp = await Get<LogEntry[]>(`/log/${count}`)
    if (resp instanceof ErrorResponse) {
        ShowErrorResponse("Failed to stop server", resp)
        throw new Error(`Failed to stop server. Response: ${resp}`)
    }
}
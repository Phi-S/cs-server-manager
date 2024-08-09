import {Get, handleErrorResponse, handleErrorResponseWithMessage, Post, PostWithoutResponse} from "@/api/api";

export class SendCommandResponse {
    output: string[] | undefined
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
    hostname: string | undefined
    server: ServerStatus | undefined
    steamcmd: SteamcmdStatus | undefined
    player_count: number | undefined
    max_player_count: number | undefined
    map: string | undefined
    ip: string | undefined
    port: string | undefined
    password: string | undefined
}

export class LogEntry {
    timestamp: string | undefined
    log_type: string | undefined
    message: string | undefined
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
    handleErrorResponse("Failed to get server logs", resp)
    return resp as LogEntry[]
}
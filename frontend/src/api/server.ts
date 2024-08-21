import {Get, handleErrorResponse, handleErrorResponseWithMessage, Post, PostWithoutResponse} from "@/api/api";

export class SendCommandResponse {
    output: string[] | undefined
}

export enum State {
    Idle = "idle",
    ServerStarting = "server-starting",
    ServerStarted = "server-started",
    ServerStopping = "server-stopping",
    SteamcmdUpdating = "steamcmd-updating",
    PluginInstalling = "plugin-installing",
    PluginUninstalling = "plugin-uninstalling"
}

export class Status {
    state: State;
    hostname: string;
    player_count: number;
    max_player_count: number;
    map: string;
    ip: string;
    port: string;
    password: string;
}

export class LogEntry {
    timestamp: string | undefined
    log_type: string | undefined
    message: string | undefined
}

export function startServer() {
    PostWithoutResponse("/start").then(value => handleErrorResponse("Failed to start server", value))
}

export function stopServer() {
    PostWithoutResponse("/stop").then(value => handleErrorResponse("Failed to stop server", value))
}

export function restartServer() {
    PostWithoutResponse("/stop").then(value => {
        handleErrorResponseWithMessage("Failed to restart server", "Server failed to stop", value)
        PostWithoutResponse("/start").then(value => handleErrorResponse("Failed to restart server", value))
    })
}

export function sendCommandWithoutResponse(command: string) {
    Post<SendCommandResponse>(`/send-command?command=${command}`).then(value => handleErrorResponse("Failed to send command", value))
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
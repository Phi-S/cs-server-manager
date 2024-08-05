
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
}

const uri = "localhost:8080"
const apiUrl = `http://${uri}/v1`

export async function GetStatus(): Promise<Status> {
    try {
        const response = await fetch(`${apiUrl}/status`)
        const status = await response.json() as Status;

        checkIfAllValuesAreDefined(Object.keys(new Status()) as (keyof Status)[], status)

        return status
    } catch (err) {
        throw new Error(`failed to get status: ${err}`);
    }
}

export function checkIfAllValuesAreDefined<T>(keys: (keyof T)[], obj: T) {
    for (const key of keys) {
        if (obj[key] === undefined) {
            throw new Error(`Property ${String(key)} is undefined`);
        }
    }
}

export class LogEntry {
    timestamp: string | undefined
    log_type: string | undefined
    message: string | undefined
}

export async function GetLogs(): Promise<LogEntry[]> {
    try {
        const response = await fetch(`${apiUrl}/log/500`)
        const logEntries = await response.json() as LogEntry[];
        return logEntries
    } catch (err) {
        throw new Error(`failed to get log entries: ${err}`);
    }
}

export class WebSocketMessage {
    type: string | undefined
    message: any
}

export function ConnectToWebSocket() {
    return new WebSocket(`ws://${uri}/v1/ws`)
}
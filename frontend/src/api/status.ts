import { checkIfAllValuesAreDefined, ErrorResponse, Get } from "./api"

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

export async function GetStatus(): Promise<Status | ErrorResponse> {
    try {
        const status = await Get<Status>(`/status`)
        if (status instanceof ErrorResponse) {
            return status
        }

        checkIfAllValuesAreDefined(Object.keys(new Status()) as (keyof Status)[], status)
        return status
    } catch (err) {
        throw new Error(`failed to get status: ${err}`);
    }
}
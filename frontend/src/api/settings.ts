import {Get, handleErrorResponse, PostJson} from "@/api/api";

export class Settings {
    hostname: string | undefined
    password: string | undefined
    start_map: string | undefined
    max_players: number | undefined
    steam_login_token: string | undefined
}

export async function getSettings(): Promise<Settings> {
    const resp = await Get<Settings>("/settings")
    handleErrorResponse("Failed to get settings", resp)
    return resp as Settings;
}

export async function updateSettings(settings: Settings): Promise<Settings> {
    const resp = await PostJson<Settings>("/settings", settings)
    handleErrorResponse("Failed to update server settings", resp)
    return resp as Settings;
}
import {Get, handleErrorResponse, PostJson} from "@/api/api";

export class Settings {
    hostname: string
    password: string
    start_map: string
    max_players: number
    steam_login_token: string
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
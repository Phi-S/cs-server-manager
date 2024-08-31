import { Get, PostJson } from "./api";

export class Settings {
  hostname: string | undefined;
  password: string | undefined;
  start_map: string | undefined;
  max_players: number | undefined;
  steam_login_token: string | undefined;
}

export async function getSettings(): Promise<Settings> {
  return await Get<Settings>("/settings");
}

export async function updateSettings(settings: Settings): Promise<Settings> {
  return await PostJson<Settings>("/settings", settings);
}

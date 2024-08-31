import { Get, PostJson, PostWithoutResponse } from "./api";

export interface SendCommandResponse {
  output: string[];
}

export enum State {
  Idle = "idle",
  ServerStarting = "server-starting",
  ServerStarted = "server-started",
  ServerStopping = "server-stopping",
  SteamcmdUpdating = "steamcmd-updating",
  PluginInstalling = "plugin-installing",
  PluginUninstalling = "plugin-uninstalling",
}

export interface Status {
  state: State;
  hostname: string;
  player_count: number;
  max_player_count: number;
  map: string;
  ip: string;
  port: string;
  password: string;
}

export interface LogEntry {
  timestamp: string;
  log_type: string;
  message: string;
}

export async function startServer() {
  return await PostWithoutResponse("/start");
}

export async function stopServer() {
  return await PostWithoutResponse("/stop");
}

export async function restartServer() {
  await stopServer();
  await startServer();
}

export async function sendCommandWithoutResponse(c: string) {
  return PostJson<SendCommandResponse>(`/command`, { command: c });
}

export async function startUpdate() {
  return await PostWithoutResponse("/update");
}

export async function cancelUpdate() {
  return await PostWithoutResponse("/update/cancel");
}

export async function getStatus(): Promise<Status> {
  return await Get<Status>("/status");
}

export async function getLogs(count: number): Promise<LogEntry[]> {
  return await Get<LogEntry[]>(`/logs/${count}`);
}

import { DeleteWithoutResponse, Get, PostJsonWithoutResponse } from "./api";

export interface PluginResp {
  name: string;
  description: string;
  url: string;
  versions: Version[];
}

export interface Version {
  name: string;
  installed: boolean;
  dependencies: Dependency[];
}

export interface Dependency {
  plugin_name: string;
  version_name: string;
}

export async function getPlugins(): Promise<PluginResp[]> {
  return await Get<PluginResp[]>("/plugins");
}

export async function installPlugin(name: string, version: string) {
  return await PostJsonWithoutResponse("/plugins", {
    name: name,
    version: version,
  });
}

export async function uninstallPlugin() {
  return await DeleteWithoutResponse(`/plugins`);
}

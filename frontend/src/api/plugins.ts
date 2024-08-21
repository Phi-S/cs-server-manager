import {Get, handleErrorResponse, PostWithoutResponse} from "@/api/api";

export interface PluginResp {
    name: string
    description: string
    url: string
    versions: Version[]
}

export interface Version {
    name: string
    installed: boolean
    dependencies: Dependency[]
}

export interface Dependency {
    plugin_name: string
    version_name: string
}

export async function getPlugins(): Promise<PluginResp[]> {
    const resp = await Get<PluginResp[]>("/plugins")
    handleErrorResponse("Failed to get plugins", resp)
    return resp as PluginResp[]
}

export async function installPlugin(name: string, version: string) {
    const resp = await PostWithoutResponse(`/plugins/install?name=${name}&version=${version}`)
    handleErrorResponse("Failed to install plugin", resp)
}

export async function uninstallPlugin(name: string) {
    const resp = await PostWithoutResponse(`/plugins/uninstall?name=${name}`)
    handleErrorResponse("Failed to uninstall plugin", resp)
}
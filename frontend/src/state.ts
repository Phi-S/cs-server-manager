import { reactive } from 'vue'

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

function Setup() {
    fetch("http://localhost:8080/v1/status", {
        method: "get",
        headers: {
            'content-type': 'application/json'
        }
    }).then(res => {
        return res.json()
    }).then(json => {
        console.log(json)
        const hostname: string = json["hostname"]
        const serverStatus: string = json["server"]
        const steamcmdStatus: string = json["steamcmd"]
        const playerCount: string = json["player-count"]
        const maxPlayerCount: string = json["max-player-count"]
        const map: string = json["map"]

        state.Hostname = hostname
        state.ServerStatus = serverStatus as ServerStatus
        state.SteamcmdStatus = steamcmdStatus as SteamcmdStatus
        state.PlayerCount = Number(playerCount)
        state.MaxPlayerCount = Number(maxPlayerCount)
        state.Map = map


    }).catch(ex => {
        console.log("failed to get initial status values. " + ex)
    })
}

function SetupWebSocket() {

}

export const state = reactive({
    IsServerBusy() {
        return state.ServerStatus == ServerStatus.ServerStatusStarting ||
            state.ServerStatus == ServerStatus.ServerStatusStopping ||
            state.SteamcmdStatus == SteamcmdStatus.SteamcmdStatusUpdating;
    },
    init() {
        Setup()
        SetupWebSocket()
    },
    //TODO api request
    Hostname: "",
    ServerStatus: ServerStatus.ServerStatusStopped,
    SteamcmdStatus: SteamcmdStatus.SteamcmdStatusStopped,
    PlayerCount: 0,
    MaxPlayerCount: 0,
    Map: ""
})
import { ref } from 'vue';
import { ConnectToWebSocket, GetLogs, GetStatus, LogEntry, ServerStatus, Status, SteamcmdStatus, WebSocketMessage, checkIfAllValuesAreDefined } from './api';

export var status = ref<Status>()
export var logEntires = ref<LogEntry[]>()

export function IsServerBusy() {
    return status.value?.server == ServerStatus.ServerStatusStarting ||
        status.value?.server == ServerStatus.ServerStatusStopping ||
        status.value?.steamcmd == SteamcmdStatus.SteamcmdStatusUpdating;
}

export async function Setup() {
    status.value = await GetStatus()
    logEntires.value = await GetLogs()
    SetupWebSocket()
}

function SetupWebSocket() {
    var socket = ConnectToWebSocket()
    socket.onmessage = (event) => {
        const msg = JSON.parse(event.data) as WebSocketMessage

        if (msg.message === undefined || msg.type == undefined) {
            console.log(`unexpected websocket message received ${event.data}`)
        }

        if (msg.type === "status") {
            const newStatus = msg.message as Status
            checkIfAllValuesAreDefined(Object.keys(new Status()) as (keyof Status)[], newStatus)
            status.value = newStatus
        } else if (msg.type === "log") {
            const newLog = msg.message as LogEntry
            checkIfAllValuesAreDefined(Object.keys(new LogEntry()) as (keyof LogEntry)[], newLog)
            logEntires.value?.unshift(newLog)
        } else {
            console.log(`unexpected websocket message received ${event.data}`)
        }
    }
}
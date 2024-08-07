import { ref } from 'vue';
import { ConnectToWebSocket, ErrorResponse, GetLogs, GetStatus, LogEntry, ServerStatus, Status, SteamcmdStatus, WebSocketMessage, checkIfAllValuesAreDefined } from '@/api/api';

export var connected = ref<boolean>(false)
export var status = ref<Status>()
export var logEntires = ref<LogEntry[]>()

export function IsServerBusy() {
    return status.value?.server == ServerStatus.ServerStatusStarting ||
        status.value?.server == ServerStatus.ServerStatusStopping ||
        status.value?.steamcmd == SteamcmdStatus.SteamcmdStatusUpdating;
}

export async function Setup() {
    connected.value = false
    while (true) {
        try {
            await initStatusAndLogs()
            SetupWebSocket()
            console.log("setup finished")
            connected.value = true
            break
        } catch (error) {
            console.log(error)
            await new Promise(r => setTimeout(r, 2000));
        }
    }
}

async function initStatusAndLogs() {
    const newStatus = await GetStatus()
    if (newStatus instanceof ErrorResponse) {
        throw new Error(`failed to initialize status. ${newStatus}`)
    }
    status.value = newStatus

    const newLogs = await GetLogs()
    if (newLogs instanceof ErrorResponse) {
        throw new Error(`failed to initialize logs. ${newLogs}`)
    }
    logEntires.value = newLogs
}

function SetupWebSocket() {
    var socket = ConnectToWebSocket()

    socket.onerror = async (event) => {
        socket.close()
        console.info(`websocket closed with error ${event}. Trying to reconnect`)
        Setup()
    }

    socket.onclose = (event) => {
        socket.close()
        console.info(`websocket closed ${event}. Trying to reconnect`)
        Setup()
    }

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
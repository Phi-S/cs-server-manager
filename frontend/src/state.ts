import {ref} from 'vue';
import type {Status} from "@/api/server";
import {getLogs, getStatus, LogEntry, ServerStatus, SteamcmdStatus} from "@/api/server";
import {ConnectToWebSocket, WebSocketMessage} from "@/api/api";


export var connected = ref<boolean>(false)
export var status = ref<Status>()
export var logEntries = ref<LogEntry[]>()

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
    status.value = await getStatus()
    logEntries.value = await getLogs(500)
}

function SetupWebSocket() {
    const socket = ConnectToWebSocket();

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
        console.log("ws: " + event.data)
        const msg = JSON.parse(event.data) as WebSocketMessage

        if (msg.message === undefined || msg.type == undefined) {
            console.log(`unexpected websocket message received ${event.data}`)
        }

        if (msg.type === "status") {
            status.value = msg.message as Status
        } else if (msg.type === "log") {
            const newLog = msg.message as LogEntry
            if (logEntries.value === undefined) {
                logEntries.value = [newLog]
            } else {
                logEntries.value.unshift(newLog)
            }
        } else {
            console.log(`unexpected websocket message received ${event.data}`)
        }
    }
}
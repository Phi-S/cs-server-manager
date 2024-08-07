<script setup lang="ts">
import { ServerStatus, SteamcmdStatus } from '@/api/api';
import { SendCommand } from '@/api/send-command';
import { CancelUpdateHandler, isErrorResponse, RestartHandler, ShowErrorResponse, StartHandler, StopHandler, UpdateHandler } from '@/buttonHandler';
import { Show } from '@/popup';
import { IsServerBusy, status } from '@/state';

function disableUpdateButton(): boolean {
    return status.value?.server === ServerStatus.ServerStatusStarting || status.value?.server === ServerStatus.ServerStatusStopping
}

function StartStop() {
    if (status.value?.server == ServerStatus.ServerStatusStopped) {
        StartHandler()
    } else {
        StopHandler()
    }
}

async function UpdateOrCancelUpdate() {
    if (status.value?.steamcmd === SteamcmdStatus.SteamcmdStatusStopped) {
        UpdateHandler()
    } else {
        CancelUpdateHandler()
    }
}

async function changeMap(map: string) {
    try {
        const resp = await SendCommand(`changelevel ${map}`)
        if (isErrorResponse(resp)) {
            ShowErrorResponse("Failed to change map", resp)
        }
    } catch (error) {
        Show("Failed to change map", `${error}`)
    }
}

const maps = [
    "de_anubis",
    "de_inferno",
    "de_dust2",
    "de_ancient",
    "de_mirage",
    "de_vertigo",
    "de_nuke"
]
</script>

<template>
    <div class="container-xxl" style="height: 60px">
        <div class="row input-group flex-nowrap w-100 m-0 h-100">
            <button @click="StartStop" class="col-3 btn btn-outline-info" :disabled="IsServerBusy()">
                <div v-if="status?.server == ServerStatus.ServerStatusStopped">Start</div>
                <div v-else>Stop</div>
            </button>
            <button @click="RestartHandler" class="col-3 btn btn-outline-info" :disabled="IsServerBusy()">
                Restart
            </button>
            <button @click="UpdateOrCancelUpdate" class="col-3 btn btn-outline-info" :disabled="disableUpdateButton()">
                <div v-if="status?.steamcmd == SteamcmdStatus.SteamcmdStatusStopped">Update</div>
                <div v-else>Cancel Update</div>
            </button>
            <button class="col-3 btn btn-outline-info dropdown" data-bs-toggle="dropdown" aria-expanded="false"
                :disabled="status?.server !== ServerStatus.ServerStatusStarted">
                Change Map
            </button>
            <ul class="dropdown-menu col-3 text-center">
                <li>
                    <button v-for="map in maps" @click="changeMap(map)" class="dropdown-item">
                        {{ map }}
                    </button>
                </li>
            </ul>

        </div>
    </div>
</template>
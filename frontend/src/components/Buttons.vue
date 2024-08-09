<script setup lang="ts">
import {IsServerBusy, status} from '@/state';
import {
  cancelUpdate, restartServer,
  sendCommand,
  ServerStatus,
  startServer,
  startUpdate,
  SteamcmdStatus,
  stopServer
} from "@/api/server";

function disableUpdateButton(): boolean {
  return status.value?.server === ServerStatus.ServerStatusStarting || status.value?.server === ServerStatus.ServerStatusStopping
}

async function StartStop() {
  if (status.value?.server == ServerStatus.ServerStatusStopped) {
    await startServer()
  } else {
    await stopServer()
  }
}

async function UpdateOrCancelUpdate() {
  if (status.value?.steamcmd === SteamcmdStatus.SteamcmdStatusStopped) {
    await startUpdate()
  } else {
    await cancelUpdate()
  }
}

async function changeMap(map: string) {
  await sendCommand(`changelevel ${map}`)
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

  <div class="input-group flex-nowrap w-100 m-0 h-100">
    <button @click="StartStop" class="col-3 btn btn-outline-info" :disabled="IsServerBusy()">
      <div v-if="status?.server == ServerStatus.ServerStatusStopped">Start</div>
      <div v-else>Stop</div>
    </button>
    <button @click="restartServer()" class="col-3 btn btn-outline-info" :disabled="IsServerBusy()">
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

</template>
<script setup lang="ts">
import {
  cancelUpdate,
  restartServer,
  sendCommandWithoutResponse,
  startServer,
  startUpdate,
  State,
  stopServer
} from "@/api/server";
import {useStatusStore} from "@/stores/status";

const statusStore = useStatusStore()

async function StartStop() {
  if (statusStore.state == State.Idle) {
    await startServer()
  } else {
    await stopServer()
  }
}

async function UpdateOrCancelUpdate() {
  if (statusStore.state === State.Idle) {
    await startUpdate()
  } else {
    await cancelUpdate()
  }
}

async function changeMap(map: string) {
  await sendCommandWithoutResponse(`changelevel ${map}`)
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
    <button @click="StartStop" class="col-3 btn btn-outline-info">
      <div v-if="statusStore.state == State.Idle">Start</div>
      <div v-else>Stop</div>
    </button>
    <button @click="restartServer()" class="col-3 btn btn-outline-info">
      Restart
    </button>
    <button @click="UpdateOrCancelUpdate" class="col-3 btn btn-outline-info" :disabled="statusStore.isServerBusy">
      <div v-if="statusStore.state !== State.SteamcmdUpdating">Update</div>
      <div v-else>Cancel Update</div>
    </button>
    <button class="col-3 btn btn-outline-info dropdown" data-bs-toggle="dropdown" aria-expanded="false"
            :disabled="statusStore.state !== State.ServerStarted">
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
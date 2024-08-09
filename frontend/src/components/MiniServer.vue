<script setup lang="ts">
import {getConnectionString, getConnectionUrl, IsServerBusy, status} from '@/state';
import {ServerStatus, startServer, stopServer} from "@/api/server";
import {copyToClipboard, navigateTo} from "@/util";
</script>

<template>
  <div
      :class="[(status?.server === ServerStatus.ServerStatusStarted) ? 'bg-success' : 'bg-warning', 'rounded-2 row p-0 m-0 h-100 ps-3 pe-2 me-3']">
    <div @click="navigateTo(getConnectionUrl())" class="col-10 col-sm-10 btn text-start black p-0 m-0">
      <div class="row h-100 align-items-center p-0 m-0">
        <div class="col-12 col-sm-7 text-truncate p-0 m-0 pe-3">
          {{ status?.hostname }}
        </div>
        <div class="col-5 d-none d-sm-block text-nowrap text-end p-0 m-0">
          {{ status?.map }}
          [ {{ status?.player_count }} / {{ status?.max_player_count }} ]
        </div>
      </div>
    </div>
    <div class="col-2 col-sm-2 h-100 p-0 m-0">
      <div class="flex-nowrap row h-100 ps-2 ps-sm-4 ps-md-5">
        <button @click="copyToClipboard(getConnectionString())" class="col-6 btn bi-copy fs-2 black p-0 m-0"></button>
        <div class="col-6 p-0 m-0">
          <button @click="stopServer()" v-if="status?.server == ServerStatus.ServerStatusStarted"
                  class="h-100 btn bi-stop fs-1 p-0 m-0 black"></button>
          <div v-else-if="IsServerBusy()" class="h-100 align-content-center pb-1">
            <button class="spinner-grow black border-0">
              <span class="visually-hidden">Loading...</span>
            </button>
          </div>
          <button @click="startServer()" v-else class="h-100 btn bi-play fs-1 p-0 m-0 black"></button>
        </div>
      </div>
    </div>
  </div>

</template>

<style scoped>
.black {
  color: black
}
</style>
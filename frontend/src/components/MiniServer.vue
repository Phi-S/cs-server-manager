<script setup lang="ts">
import {IsServerBusy, status} from '@/state';
import {ServerStatus, startServer, stopServer} from "@/api/server";
</script>

<template>
  <div
      :class="[(status?.server === ServerStatus.ServerStatusStarted) ? 'bg-success' : 'bg-warning', 'rounded-2 row p-0 m-0 h-100']">
    <div class="col-11 btn text-start black p-0 m-0 ps-2">
      <div class="row text-nowrap h-100 align-items-center">
        <div class="col-12 col-sm-8 text-truncate">
          {{ status?.hostname }}
        </div>
        <div class="col-4 d-none d-sm-block text-end">
          {{ status?.map }}
          [ {{ status?.player_count }} / {{ status?.max_player_count }}]
        </div>
      </div>
    </div>
    <div class="col-1 text-end p-0 m-0 pe-2">
      <button @click="stopServer()" v-if="status?.server == ServerStatus.ServerStatusStarted"
              class="h-100 btn bi-stop fs-2 p-0 m-0 black"></button>

      <div v-else-if="IsServerBusy()" class="h-100 align-content-center">
        <button class="spinner-grow mb-1 black border-0">
          <span class="visually-hidden">Loading...</span>
        </button>
      </div>
      <button @click="startServer()" v-else class="h-100 btn bi-play fs-2 p-0 m-0 black"></button>
    </div>
  </div>

</template>

<style scoped>
.black {
  color: black
}
</style>
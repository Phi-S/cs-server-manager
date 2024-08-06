<script setup lang="ts">
import { ServerStatus, Start } from '@/api/api';
import { Stop } from '@/api/stop';
import { IsServerBusy, status } from '@/state';

function getBackgroundClass(): string {
    if (status.value?.server === ServerStatus.ServerStatusStarted) {
        return "bg-success"
    } else {
        return ""
    }
}

</script>

<template>
    <div :class="[(status?.server === ServerStatus.ServerStatusStarted) ? 'bg-success' : 'bg-warning', 'container-fluid rounded-2']"
        style="height: 50px; max-width: 600px;">
        <div class="row h-100 align-items-center flex-nowrap m-0 p-0">
            <div class="col-11 btn text-start black p-0 m-0">
                <div class="row text-nowrap flex-nowrap">
                    <div class="col-12 col-sm-9 text-truncate">
                        {{ status?.hostname }}
                    </div>
                    <div class="col-3 d-none d-sm-block text-end">
                        {{ status?.map }}
                        [ {{ status?.player_count }} / {{ status?.max_player_count }}]
                    </div>
                </div>
            </div>
            <div class="col-1 text-end p-0 m-0">
                <button @click="Stop" v-if="status?.server == ServerStatus.ServerStatusStarted"
                    class="btn bi-stop fs-2 p-0 m-0 black"></button>
                <button v-else-if="IsServerBusy()" class="spinner-grow fs-2 p-0 m-0 black"> </button>
                <button @click="Start()" v-else class="btn bi-play fs-2 p-0 m-0 black"></button>
            </div>
        </div>
    </div>

</template>

<style scoped>
.black {
    color: black
}
</style>
<script setup lang="ts">
import { ServerStatus, state } from "@/state";
</script>

<template>
    <div class="container bg-success rounded-2" style="height: 3rem; max-width: 40rem">
        <div class="row h-100 align-items-center">
            <div class="col-11 btn text-start">
                <div class="row text-nowrap flex-nowrap">
                    <div class="col-12 col-sm-9 text-truncate">
                        {{ state.Hostname }}
                    </div>
                    <div class="col-3 d-none d-sm-block text-end">
                        {{ state.Map }}
                        [ {{ state.PlayerCount }} / {{ state.MaxPlayerCount }}]
                    </div>
                </div>
            </div>
            <div class="col-1 text-end">
                <button v-if="state.ServerStatus == ServerStatus.ServerStatusStarted"
                    class="btn bi-stop fs-2 p-0 m-0"></button>
                <button v-else-if="state.IsServerBusy()" class="spinner-grow fs-2 p-0 m-0"> </button>
                <button v-else class="btn bi-play fs-2 p-0 m-0"></button>
            </div>
        </div>
    </div>

</template>
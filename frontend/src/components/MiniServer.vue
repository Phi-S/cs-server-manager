<script setup lang="ts">
import { ServerStatus } from '@/api';
import { IsServerBusy, status } from '@/state';

</script>

<template>
    <div class="container bg-success rounded-2" style="height: 3rem; max-width: 40rem">
        <div class="row h-100 align-items-center">
            <div class="col-11 btn text-start">
                <div class="row text-nowrap flex-nowrap">
                    <div class="col-12 col-sm-9 text-truncate">
                        {{ status?.hostname }}
                    </div>
                    <div class="col-3 d-none d-sm-block text-end">
                        {{ status?.map }}
                        [ {{ status?.player_countasdfasdf }} / {{ status?.max_player_count }}]
                    </div>
                </div>
            </div>
            <div class="col-1 text-end">
                <button v-if="status?.server == ServerStatus.ServerStatusStarted"
                    class="btn bi-stop fs-2 p-0 m-0"></button>
                <button v-else-if="IsServerBusy()" class="spinner-grow fs-2 p-0 m-0"> </button>
                <button v-else class="btn bi-play fs-2 p-0 m-0"></button>
            </div>
        </div>
    </div>

</template>
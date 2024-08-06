<script setup lang="ts">
import { ServerStatus, Start, SteamcmdStatus } from '@/api/api';
import { Stop } from '@/api/stop';
import { CancelUpdate, Update } from '@/api/update';
import { IsServerBusy, status } from '@/state';

function disableUpdateButton(): boolean {
    return status.value?.server === ServerStatus.ServerStatusStarting || status.value?.server === ServerStatus.ServerStatusStopping
}

</script>

<template>
    <div class="container-xxl" style="height: 60px">
        <div class="row input-group flex-nowrap w-100 m-0 h-100">
            <button class="col-3 btn btn-outline-info" :disabled="IsServerBusy()">
                <div @click="Start()" v-if="status?.server == ServerStatus.ServerStatusStopped">Start</div>
                <div @click="Stop" v-else>Stop</div>
            </button>
            <button @click="async () => {
                await Stop()
                await Start()
            }" class="col-3 btn btn-outline-info" :disabled="IsServerBusy()">
                Restart
            </button>
            <button class="col-3 btn btn-outline-info" :disabled="disableUpdateButton()">
                <div @click="async () => {
                    await Stop()
                    await Update()
                }" v-if="status?.steamcmd == SteamcmdStatus.SteamcmdStatusStopped">Update</div>
                <div @click="CancelUpdate" v-else>Cancel Update</div>
            </button>
            <button class="col-3 btn btn-outline-info" :disabled="status?.server !== ServerStatus.ServerStatusStarted">
                Change Map
            </button>
        </div>
    </div>
</template>
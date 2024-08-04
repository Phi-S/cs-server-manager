<script setup lang="ts">
import { ServerStatus, state, SteamcmdStatus } from '@/state';
</script>

<template>
    <div class="container" style="height: 60px">
        <div class="row input-group flex-nowrap h-100 w-100 m-0">
            <button class="col-3 btn btn-outline-info" :disabled="state.IsServerBusy()">
                <div v-if="state.ServerStatus == ServerStatus.ServerStatusStopped">Start</div>
                <div v-else>Stop</div>
            </button>
            <button class="col-3 btn btn-outline-info" :disabled="state.IsServerBusy()">
                Restart
            </button>
            <button class="col-3 btn btn-outline-info" :disabled="state.IsServerBusy()">
                <div v-if="state.SteamcmdStatus == SteamcmdStatus.SteamcmdStatusUpdating">Cancel Update</div>
                <div v-else>Update</div>
            </button>
            <button class="col-3 btn btn-outline-info" :disabled="state.IsServerBusy()">
                Change Map
            </button>
        </div>
    </div>
</template>
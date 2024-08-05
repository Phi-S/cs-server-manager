<script setup lang="ts">
import { ServerStatus, SteamcmdStatus } from '@/api';
import { IsServerBusy, status } from '@/state';
</script>

<template>
    <div class="container" style="height: 60px">
        <div class="row input-group flex-nowrap h-100 w-100 m-0">
            <button class="col-3 btn btn-outline-info" :disabled="IsServerBusy()">
                <div v-if="status?.server == ServerStatus.ServerStatusStopped">Start</div>
                <div v-else>Stop</div>
            </button>
            <button class="col-3 btn btn-outline-info" :disabled="IsServerBusy()">
                Restart
            </button>
            <button class="col-3 btn btn-outline-info" :disabled="IsServerBusy()">
                <div v-if="status?.steamcmd == SteamcmdStatus.SteamcmdStatusUpdating">Cancel Update</div>
                <div v-else>Update</div>
            </button>
            <button class="col-3 btn btn-outline-info" :disabled="IsServerBusy()">
                Change Map
            </button>
        </div>
    </div>
</template>
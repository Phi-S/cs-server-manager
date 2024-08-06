<script setup lang="ts">
import type { LogEntry } from '@/api/logs';
import { ServerStatus } from '@/api/status';
import { SendCommandHandler } from '@/buttonHandler';
import { logEntires, status } from '@/state';
import moment from 'moment';
import { onMounted } from 'vue';

function TimestampString(timestampUtc: string): string {
    const offset = new Date().getTimezoneOffset();
    return moment(timestampUtc).add(offset).format("HH:mm:ss")
}

function GetLogBackground(log: LogEntry): string {
    if (log.log_type === "system-info") {
        return "bg-success bg-opacity-75"
    } else if (log.log_type === "system-error") {
        return "bg-danger bg-opacity-75"
    } else {
        return ""
    }
}

var command = ""

async function SendCommandIfServerIsRunning() {
    if (status.value?.server !== ServerStatus.ServerStatusStarted) {
        return
    }

    if (command.trim().length == 0) {
        return
    }

    await SendCommandHandler(command)
}

onMounted(() => {
    document.getElementById("command-input")?.focus()
})
</script>

<template>
    <div class="container-xxl">
        <div class="text-nowrap align-items-center d-flex justify-content-center pb-2 pt-1" style="height: 40px;">
            <div class="input-group justify-content-center d-flex w-100">
                <input id="command-input" v-on:keyup.enter="SendCommandIfServerIsRunning" v-model="command"
                    class="input-group-text" style="width: 70%" placeholder="Server command" />
                <button @click="SendCommandIfServerIsRunning"
                    :disabled="status?.server !== ServerStatus.ServerStatusStarted" class="btn btn-outline-info"
                    style="width: 30%">Send</button>
            </div>
        </div>

        <div class="overflow-x-scroll rounded-3 border border-2" style="height: calc(100vh - 180px);">
            <table class="table table-sm table-striped">
                <tr v-for="log in logEntires" :class="[GetLogBackground(log), 'border-bottom']">
                    <td class="ps-2 pe-2 pt-1 text-nowrap border-end">
                        {{ TimestampString(log.timestamp as string) }}
                    </td>
                    <td class="ps-2" style="word-wrap: anywhere">
                        {{ log.message }}
                    </td>
                </tr>
            </table>
        </div>
    </div>

</template>
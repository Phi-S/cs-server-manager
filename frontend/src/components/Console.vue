<script setup lang="ts">
import {logEntries, status} from '@/state';
import moment from 'moment';
import {onMounted, ref} from 'vue';
import {type LogEntry, sendCommand, ServerStatus} from "@/api/server";

const command = ref("");

function timestampString(timestampUtc: string): string {
  const offset = new Date().getTimezoneOffset();
  return moment(timestampUtc).add(offset).format("HH:mm:ss")
}

function getLogBackgroundColor(log: LogEntry): string {
  if (log.log_type === "system_info") {
    return "bg-success bg-opacity-75"
  } else if (log.log_type === "system_error") {
    return "bg-danger bg-opacity-75"
  } else {
    return ""
  }
}

async function sendCommandIfServerIsRunning() {
  if (status.value?.server !== ServerStatus.ServerStatusStarted) {
    return
  }

  if (command.value.trim().length == 0) {
    return
  }

  await sendCommand(command.value)
}

onMounted(() => {
  document.getElementById("command-input")?.focus()
})
</script>

<template>
  <div class="text-nowrap align-items-center d-flex justify-content-center pb-2 pt-1" style="height: 40px;">
    <div class="input-group justify-content-center d-flex w-100">
      <input id="command-input" v-on:keyup.enter="sendCommandIfServerIsRunning" v-model="command"
             class="input-group-text" style="width: 70%" placeholder="Server command"/>
      <button @click="sendCommandIfServerIsRunning"
              :disabled="status?.server !== ServerStatus.ServerStatusStarted" class="btn btn-outline-info"
              style="width: 30%">Send
      </button>
    </div>
  </div>

  <div class="overflow-x-scroll rounded-3 border border-2" style="height: calc(100% - 40px);">
    <table class="table table-sm table-striped">
      <tr v-for="log in logEntries" :class="[getLogBackgroundColor(log), 'border-bottom']">
        <td class="ps-2 pe-2 pt-1 text-nowrap border-end">
          {{ timestampString(log.timestamp as string) }}
        </td>
        <td class="ps-2" style="word-wrap: anywhere">
          {{ log.message }}
        </td>
      </tr>
    </table>
  </div>
</template>
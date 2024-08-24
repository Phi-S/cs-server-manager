import {defineStore} from "pinia";
import type {LogEntry} from "@/api/server";
import {getLogs} from "@/api/server";

export const useLogsStore = defineStore("logs", {
    state: () => {
        return {
            logs: [] as LogEntry[]
        }
    },
    getters: {
        sortedByDate(state) {
            return state.logs
        }
    },
    actions: {
        async init() {
            this.logs = await getLogs(500)
        },
        unshift(log: LogEntry) {
            this.logs.unshift(log)
        }
    }
})
import {defineStore} from "pinia";
import type {Status} from "@/api/server"
import {State} from "@/api/server";

export const useStatusStore = defineStore("status", {
    state: (): Status => {
        return {
            state: State.Idle,
            port: "",
            password: "",
            hostname: "",
            map: "",
            ip: "",
            max_player_count: 0,
            player_count: 0
        }
    },
    getters: {
        isServerBusy: (state: Status) => state.state !== State.Idle,
        getConnectionUrl(state: Status) {
            let connectUrl = `steam://connect/${state.ip}:${state.port}`;
            if (state.password !== "") {
                connectUrl = `${connectUrl}/${state.password}`;
            }
            return connectUrl
        },
        getConnectionString(state: Status) {
            let connectString = `connect ${state.ip}:${state.port}`;
            if (state.password !== undefined && state.password !== "") {
                connectString = `${connectString}; password ${state.password}`;
            }
            return connectString
        }
    },
    actions: {
        async init() {
            //const newStatus = await getStatus()
            //this.update(newStatus)
        },
        update(status: Status) {
            this.state = status.state;
            this.hostname = status.hostname;
            this.player_count = status.player_count;
            this.max_player_count = status.max_player_count;
            this.map = status.map;
            this.ip = status.ip;
            this.port = status.port;
            this.password = status.password;
        }
    }
})
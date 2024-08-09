<script setup lang="ts">
import {getSettings, Settings, updateSettings} from "@/api/settings";
import {ref} from "vue";
import {restartServer} from "@/api/server";

let settings = ref<Settings>(new Settings());

getSettings().then(respSettings => {
  settings.value = respSettings
})

async function saveAndRestart() {
  await updateSettings(settings.value)
  await restartServer()
}
</script>

<template>
  <div class="w-100 d-flex justify-content-center">
    <form class="w-100" style="max-width: 600px">
      <div class="form-group mb-3">
        <label for="hostname">Hostname</label>
        <input v-model="settings.hostname" type="text" class="form-control" id="hostname">
      </div>
      <div class="form-group mb-3">
        <label for="password">Server password</label>
        <input v-model="settings.password" type="text" class="form-control" id="password">
      </div>
      <div class="form-group mb-3">
        <label for="startMap">Start map</label>
        <input v-model="settings.start_map" type="text" class="form-control" id="startMap"
               aria-describedby="statMapHelp">
        <small id="StartMapHelp" class="form-text text-muted">The map which is loaded when server starts</small>
      </div>
      <div class="form-group mb-3">
        <label for="maxPlayers">May Players</label>
        <input v-model.number="settings.max_players" type="number" class="form-control" id="maxPlayers">
      </div>
      <div class="form-group mb-3">
        <label for="steamLoginToken">Steam login token</label>
        <input v-model="settings.steam_login_token" type="text" class="form-control" id="steamLoginToken"
               aria-describedby="steamLoginTokenHelp">
        <small id="steamLoginTokenHelp" class="form-text text-muted">
          You can generate a token <a href="https://steamcommunity.com/dev/managegameservers">Here</a>
        </small>
      </div>
      <button @click.prevent="saveAndRestart" type="submit" class="btn btn-outline-info me-2">
        Save and Restart
      </button>
      <button @click.prevent="updateSettings(settings)" type="submit" class="btn btn-outline-info">
        Save
      </button>
      <br>
      <small class="form-text text-muted">A server restart is required for those settings to take effect</small>
    </form>
  </div>
</template>

<style>
</style>

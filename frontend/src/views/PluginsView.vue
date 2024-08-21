<script setup lang="ts">
import {getPlugins, installPlugin, type PluginResp, uninstallPlugin} from "@/api/plugins";
import {onMounted, ref} from "vue";
import OverlaySpinner from "@/components/OverlaySpinner.vue";
import {useStatusStore} from "@/stores/status";
import {State} from "@/api/server";

const statusStore = useStatusStore()

const selectedVersions = ref<Map<string, string>>(new Map<string, string>())

let plugins = ref<PluginResp[]>();

function isPluginInstalled(name: string, version: string): boolean {
  for (const plugin of plugins.value) {
    if (plugin.name === name) {
      for (const pluginVersion of plugin.versions) {
        if (pluginVersion.name === version) {
          return pluginVersion.installed
        }
      }
    }
  }

  return false
}

function updatePlugins() {
  getPlugins().then(value => {
    plugins.value = value
    for (const plugin of plugins.value) {
      if (plugin.versions.length <= 0) {
        continue
      }

      selectedVersions.value[plugin.name] = plugin.versions[0].name
    }
  })
}

onMounted(() => {
  updatePlugins()
})
</script>

<template>
  <overlay-spinner v-if="statusStore.state === State.PluginInstalling" message="installing plugin"/>
  <overlay-spinner v-else-if="statusStore.state === State.PluginUninstalling" message="uninstalling plugin"/>
  <table>
    <tr class="row pb-4" v-for="plugin in plugins">
      <td class="col-2 text-center">
        <a :href="plugin.url">{{ plugin.name }}</a>
      </td>
      <td class="col-6">{{ plugin.description }}</td>
      <td class="col-2">
        <select class="w-100 form-select" v-model="selectedVersions[plugin.name]">
          <option v-for="version in plugin.versions" :value="version.name">
            {{ version.name }}
          </option>
        </select>
      </td>
      <td class="col-2">
        <button class="btn btn-outline-info w-100" v-if="isPluginInstalled(plugin.name, selectedVersions[plugin.name])"
                @click="uninstallPlugin(plugin.name).then(_ => updatePlugins())">
          Uninstall
        </button>
        <button class="btn btn-outline-info w-100" v-else
                @click="installPlugin(plugin.name, selectedVersions[plugin.name]).then(_ => updatePlugins())">
          Install
        </button>
      </td>
    </tr>
  </table>
</template>

<style></style>

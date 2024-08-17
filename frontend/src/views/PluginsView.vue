<script setup lang="ts">
import {getPlugins, type PluginResp} from "@/api/plugins";
import {onMounted, ref} from "vue";

const selectedVersions = ref<Map<string, string>>(new Map<string, string>())

let plugins = ref<PluginResp[]>();
onMounted(() => {
  getPlugins().then(value => {
    plugins.value = value
    for (const plugin of plugins.value) {
      if (plugin.versions.length <= 0) {
        continue
      }

      selectedVersions.value[plugin.name] = plugin.versions[0].name
    }
  })
})
</script>

<template>
  <table>
    <tr v-for="plugin in plugins">
      <td>{{ plugin.name }}</td>
      <td>{{ plugin.url }}</td>
      <td>{{ plugin.description }}</td>
      <td>
        <select v-model="selectedVersions[plugin.name]">
          <option v-for="version in plugin.versions" :value="version.name">
            {{ version.name }}
          </option>
        </select>
      </td>
      <td>
        <button>Install</button>
      </td>
    </tr>
  </table>
</template>

<style></style>

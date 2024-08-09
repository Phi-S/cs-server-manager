<script setup lang="ts">
import Buttons from "@/components/Buttons.vue";
import Console from "@/components/Console.vue";
import {connected} from "@/state";
</script>

<template>
  <div v-if="connected === false" class="position-absolute bg-dark w-100 h-100 text-center fs-1 bg-opacity-75 z-3 pt-3">
    <div class="bi spinner-border"></div>
    connecting...
  </div>

  <div class="mb-2"></div>
  <div style="height: 40px;">
    <Buttons></Buttons>
  </div>
  <div class="mb-2"></div>

  <div style="height: calc(100vh - 130px);">
    <Console></Console>
  </div>
</template>
<script setup lang="ts">
import Buttons from "@/components/Buttons.vue";
import Console from "@/components/Console.vue";
</script>

<template>
  <div class="mb-2"></div>
  <div style="height: 40px;">
    <Buttons></Buttons>
  </div>
  <div class="mb-2"></div>

  <div style="height: calc(100vh - 130px);">
    <Console></Console>
  </div>
</template>
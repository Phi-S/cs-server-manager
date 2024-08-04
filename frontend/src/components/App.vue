<script setup lang="ts">
import Buttons from "@/components/Buttons.vue";
import MiniServer from "@/components/MiniServer.vue";

const loc = window.location

</script>

<template>

    <h1>{{ loc }}</h1>
    <MiniServer></MiniServer>
    <div class="mb-2"></div>
    <Buttons></Buttons>
</template>

<style scoped></style>

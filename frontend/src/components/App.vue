<script setup lang="ts">
import Buttons from "@/components/Buttons.vue";
import Console from "@/components/Console.vue";
import MiniServer from "@/components/MiniServer.vue";
import Modal from "@/components/Modal.vue";
import { connected } from "@/state";
</script>

<template>
    <Modal></Modal>
    <div v-if="connected === false"
        class="position-absolute bg-dark w-100 h-100 text-center fs-1 bg-opacity-75 z-3 pt-3">
        <div class="bi spinner-border"> </div>
        connecting...
    </div>

    <div class="mb-2"></div>
    <MiniServer></MiniServer>
    <div class="mb-2"></div>
    <Buttons></Buttons>
    <div class="mb-2"></div>
    <Console></Console>
</template>
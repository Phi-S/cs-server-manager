<script setup lang="ts">
import { Hide, messages, shouldShow, title } from '@/popup';
</script>

<template>
    <div v-if="shouldShow" class="modal modal-dialog-centered justify-content-center" tabindex="-1">
        <div class="modal-dialog" style="width: 400px;">
            <div class="modal-content border border-2 border-danger">
                <div class="modal-header">
                    <h5 class="modal-title w-100 text-center">{{ title }} </h5>
                    <button @click="Hide" class="btn-close"></button>
                </div>
                <div class="modal-body text-center">
                    <p v-for="msg in messages">
                        {{ msg }}</p>
                </div>
                <div class="modal-footer">
                    <button @click="Hide" class="btn btn-secondary">Close</button>
                </div>
            </div>
        </div>
    </div>
</template>
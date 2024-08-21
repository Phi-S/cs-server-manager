<script setup lang="ts">
defineProps({
  message: String
})
</script>


<template>
  <div class="position-absolute bg-dark w-100 h-100 text-center fs-1 bg-opacity-75 z-3 pt-3">
    <div class="bi spinner-border"></div>
    {{ message }}
  </div>
</template>

<style scoped>

</style>
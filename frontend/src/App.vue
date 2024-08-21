<script setup lang="ts">
import {RouterView} from 'vue-router';
import MiniServer from './components/MiniServer.vue';
import ErrorDisplay from "@/components/ErrorDisplay.vue";
import router from "@/router";
import {connected} from "@/webSocket";
import OverlaySpinner from "@/components/OverlaySpinner.vue";

function navigateToHome() {
  router.push("/")
}
</script>

<template>
  <OverlaySpinner message="connecting..." v-if="connected === false"/>

  <div class="container-lg">

    <ErrorDisplay></ErrorDisplay>

    <header class="mt-2 container-lg mb-4 pe-5">
      <div class="row w-100 flex-nowrap" style="height: 50px;">
        <button tabindex="1" class="col-1 btn btn-outline-info bi bi-list-stars border-0 fs-1 me-2" style="width: 45px;"
                data-bs-toggle="offcanvas"
                data-bs-target="#navbar" aria-controls="navbar">
        </button>

        <img @click="navigateToHome" class="mt-1 me-2 border-start" src="@/assets/logo.png" alt="logo" href="/"
             style="height: 40px; width: 60px; cursor: pointer"/>

        <div class="col-10 mx-auto" style="max-width: 600px;">
          <MiniServer></MiniServer>
        </div>
        <div class="col-1" style="width: 45px;">

        </div>
      </div>
    </header>

    <hr>

    <div class="offcanvas offcanvas-start" style="width: 20rem;" tabindex="-1" id="navbar"
         aria-labelledby="navbarLabel">
      <div class="offcanvas-header">
        <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close"></button>
      </div>
      <div class="offcanvas-body">
        <div class="fs-3 mx-2">
          <RouterLink to="/" active-class="active" class="btn btn-outline-info nav-link d-flex flex-nowrap">
            <div class="bi bi-house col-2">
            </div>
            <div class="col-10">
              Home
            </div>
          </RouterLink>
          <RouterLink to="/settings" active-class="active" class="btn btn-outline-info nav-link d-flex flex-nowrap">
            <div class="bi bi-gear col-2">
            </div>
            <div class="col-10">
              Settings
            </div>
          </RouterLink>
          <RouterLink to="/plugins" active-class="active" class="btn btn-outline-info nav-link d-flex flex-nowrap">
            <div class="bi bi-plugin col-2">
            </div>
            <div class="col-10">
              Plugins
            </div>
          </RouterLink>
          <RouterLink to="/logs" active-class="active" class="btn btn-outline-info nav-link d-flex flex-nowrap">
            <div class="bi bi-card-text col-2">
            </div>
            <div class="col-10">
              Logs
            </div>
          </RouterLink>
          <RouterLink to="/about" active-class="active" class="btn btn-outline-info nav-link d-flex flex-nowrap">
            <div class="bi bi-info-square col-2">
            </div>
            <div class="col-10">
              About
            </div>
          </RouterLink>
        </div>
      </div>
    </div>

    <RouterView/>

  </div>
</template>

<style scoped>

</style>

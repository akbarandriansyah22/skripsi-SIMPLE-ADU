<template>
  <div class="min-h-screen bg-slate-50 text-slate-900">
    <router-view />
    <ToastContainer />
  </div>
</template>

<script setup>
import ToastContainer from "./components/ToastContainer.vue";
</script>

<template>
  <div class="fixed right-4 top-6 z-50 flex flex-col gap-3">
    <transition-group name="toast" tag="div">
      <div
        v-for="t in toasts"
        :key="t.id"
        class="w-80 max-w-full rounded-lg border bg-white p-3 shadow-lg ring-1 ring-slate-200 flex items-start gap-3"
      >
        <div v-if="t.type === 'success'" class="text-emerald-500">✓</div>
        <div v-else-if="t.type === 'danger'" class="text-red-500">⚠</div>
        <div v-else class="text-slate-500">i</div>
        <div class="flex-1 text-sm text-slate-800">{{ t.message }}</div>
        <button
          @click="remove(t.id)"
          class="text-slate-400 hover:text-slate-600"
        >
          ✕
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { useToastStore } from "../stores/toast";

const store = useToastStore();
const toasts = computed(() => store.toasts);
function remove(id) {
  store.remove(id);
}
</script>

<style scoped>
.toast-enter-from {
  opacity: 0;
  transform: translateY(-6px);
}
.toast-enter-active {
  transition: all 200ms ease;
}
.toast-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
.toast-leave-active {
  transition: all 200ms ease;
}
</style>

<template>
  <transition name="modal-fade">
    <div
      v-if="visible"
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 p-4"
    >
      <div
        class="w-full max-w-xl rounded-3xl bg-white shadow-elevated ring-1 ring-slate-200 overflow-hidden"
      >
        <div
          class="flex items-center justify-between border-b border-slate-200 px-6 py-4"
        >
          <div>
            <p class="text-sm font-semibold text-slate-900">{{ title }}</p>
            <p v-if="subtitle" class="mt-1 text-sm text-slate-500">
              {{ subtitle }}
            </p>
          </div>
          <button
            @click="$emit('close')"
            class="text-slate-400 hover:text-slate-700"
          >
            ✕
          </button>
        </div>
        <div class="px-6 py-5">
          <slot />
        </div>
        <div v-if="$slots.footer" class="border-t border-slate-200 px-6 py-4">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
const props = defineProps({
  visible: Boolean,
  title: {
    type: String,
    default: "",
  },
  subtitle: {
    type: String,
    default: "",
  },
});
</script>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition:
    opacity 200ms ease,
    transform 200ms ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
  transform: translateY(-12px);
}
</style>

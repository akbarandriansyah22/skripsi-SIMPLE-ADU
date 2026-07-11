import { defineStore } from "pinia";
import { ref } from "vue";

export const useToastStore = defineStore("toast", () => {
  const toasts = ref([]);

  function add({ type = "info", message = "", timeout = 4000 } = {}) {
    const id = Date.now().toString() + Math.floor(Math.random() * 1000);
    toasts.value.push({ id, type, message });
    // auto remove
    setTimeout(() => remove(id), timeout);
    return id;
  }

  function remove(id) {
    toasts.value = toasts.value.filter((t) => t.id !== id);
  }

  function clear() {
    toasts.value = [];
  }

  return { toasts, add, remove, clear };
});

export default useToastStore;

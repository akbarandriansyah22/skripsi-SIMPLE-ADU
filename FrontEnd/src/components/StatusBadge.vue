<template>
  <span :class="['inline-flex rounded-full px-3 py-1 text-xs font-semibold', badgeClass]">
    {{ label || 'Belum tersedia' }}
  </span>
</template>

<script setup>
import { computed } from 'vue'
import { statusClass, urgencyClass, sentimentClass } from '../utils/api'

const props = defineProps({
  label: { type: [String, Number], default: '' },
  kind: { type: String, default: 'status' },
})

const badgeClass = computed(() => {
  if (props.kind === 'urgency') return urgencyClass(props.label)
  if (props.kind === 'sentiment') return sentimentClass(props.label)
  return statusClass(props.label)
})
</script>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  status: string
}

const props = defineProps<Props>()

const statusConfig = computed(() => {
  const status = props.status.toLowerCase()
  switch (status) {
    case 'active':
      return { bg: 'bg-success/10', text: 'text-success', dot: 'bg-success' }
    case 'inactive':
      return { bg: 'bg-text-muted/10', text: 'text-text-muted', dot: 'bg-text-muted' }
    case 'recovery':
      return { bg: 'bg-warning/10', text: 'text-warning', dot: 'bg-warning' }
    case 'banned':
      return { bg: 'bg-danger/10', text: 'text-danger', dot: 'bg-danger' }
    default:
      return { bg: 'bg-text-muted/10', text: 'text-text-muted', dot: 'bg-text-muted' }
  }
})

const displayText = computed(() => {
  return props.status.charAt(0).toUpperCase() + props.status.slice(1).toLowerCase()
})
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium"
    :class="[statusConfig.bg, statusConfig.text]"
  >
    <span class="w-1.5 h-1.5 rounded-full" :class="statusConfig.dot"></span>
    {{ displayText }}
  </span>
</template>





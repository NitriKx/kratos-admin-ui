<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  name?: string
  email?: string
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md'
})

const initials = computed(() => {
  if (props.name) {
    const parts = props.name.split(' ')
    if (parts.length >= 2) {
      const first = parts[0] ?? ''
      const second = parts[1] ?? ''
      if (first[0] && second[0]) {
        return (first[0] + second[0]).toUpperCase()
      }
    }
    return props.name.substring(0, 2).toUpperCase()
  }
  if (props.email) {
    return props.email.substring(0, 2).toUpperCase()
  }
  return '??'
})

const sizeClasses = computed(() => {
  switch (props.size) {
    case 'sm': return 'w-8 h-8 text-xs'
    case 'lg': return 'w-12 h-12 text-base'
    default: return 'w-10 h-10 text-sm'
  }
})

// Generate consistent color based on name/email
const bgColor = computed(() => {
  const colors = [
    'bg-blue-600',
    'bg-green-600',
    'bg-purple-600',
    'bg-pink-600',
    'bg-orange-600',
    'bg-teal-600',
    'bg-indigo-600',
    'bg-cyan-600'
  ]
  const str = props.name || props.email || ''
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
})
</script>

<template>
  <div
    class="rounded-full flex items-center justify-center font-medium text-white shrink-0"
    :class="[sizeClasses, bgColor]"
  >
    {{ initials }}
  </div>
</template>


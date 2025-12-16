<script setup lang="ts">
import { X, CheckCircle, XCircle, AlertTriangle, Info } from 'lucide-vue-next'
import { useToast } from '@/composables/useToast'

const { toasts, removeToast } = useToast()

const icons = {
  success: CheckCircle,
  error: XCircle,
  warning: AlertTriangle,
  info: Info
}

const colors = {
  success: {
    bg: 'bg-success/10',
    border: 'border-success/20',
    icon: 'text-success',
    title: 'text-success'
  },
  error: {
    bg: 'bg-danger/10',
    border: 'border-danger/20',
    icon: 'text-danger',
    title: 'text-danger'
  },
  warning: {
    bg: 'bg-warning/10',
    border: 'border-warning/20',
    icon: 'text-warning',
    title: 'text-warning'
  },
  info: {
    bg: 'bg-primary/10',
    border: 'border-primary/20',
    icon: 'text-primary',
    title: 'text-primary'
  }
}
</script>

<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-[100] flex flex-col gap-3 max-w-md w-full pointer-events-none">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="[
            'pointer-events-auto flex items-start gap-3 p-4 rounded-lg border shadow-lg backdrop-blur-sm',
            colors[toast.type].bg,
            colors[toast.type].border
          ]"
        >
          <component
            :is="icons[toast.type]"
            :class="['w-5 h-5 flex-shrink-0 mt-0.5', colors[toast.type].icon]"
          />
          <div class="flex-1 min-w-0">
            <p :class="['font-medium text-sm', colors[toast.type].title]">
              {{ toast.title }}
            </p>
            <p v-if="toast.message" class="text-sm text-text-secondary mt-1 break-words">
              {{ toast.message }}
            </p>
          </div>
          <button
            @click="removeToast(toast.id)"
            class="flex-shrink-0 p-1 text-text-muted hover:text-text-primary transition-colors rounded"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active {
  transition: all 0.3s ease-out;
}

.toast-leave-active {
  transition: all 0.2s ease-in;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.toast-move {
  transition: transform 0.3s ease;
}
</style>




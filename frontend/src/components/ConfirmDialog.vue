<script setup lang="ts">
import { computed } from 'vue'
import { AlertTriangle, Trash2, X, AlertCircle } from 'lucide-vue-next'

interface Props {
  isOpen: boolean
  title?: string
  message: string
  confirmText?: string
  cancelText?: string
  variant?: 'danger' | 'warning' | 'info'
  loading?: boolean
}

interface Emits {
  (e: 'confirm'): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Confirm Action',
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  variant: 'danger',
  loading: false
})

const emit = defineEmits<Emits>()

const variantClasses = computed(() => {
  switch (props.variant) {
    case 'danger':
      return {
        iconBg: 'bg-danger/10',
        iconColor: 'text-danger',
        buttonBg: 'bg-danger hover:bg-danger/90',
        icon: Trash2
      }
    case 'warning':
      return {
        iconBg: 'bg-warning/10',
        iconColor: 'text-warning',
        buttonBg: 'bg-warning hover:bg-warning/90',
        icon: AlertTriangle
      }
    case 'info':
    default:
      return {
        iconBg: 'bg-primary/10',
        iconColor: 'text-primary',
        buttonBg: 'bg-primary hover:bg-primary-hover',
        icon: AlertCircle
      }
  }
})

const handleConfirm = () => {
  if (!props.loading) {
    emit('confirm')
  }
}

const handleCancel = () => {
  if (!props.loading) {
    emit('cancel')
  }
}

// Close on escape key
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && !props.loading) {
    emit('cancel')
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="isOpen"
        class="fixed inset-0 z-[60] flex items-center justify-center p-4"
        @keydown="handleKeydown"
      >
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-black/50 backdrop-blur-sm"
          @click="handleCancel"
        />
        
        <!-- Dialog -->
        <div 
          class="relative bg-card border border-border rounded-xl shadow-xl w-full max-w-md overflow-hidden"
          role="alertdialog"
          aria-modal="true"
          :aria-labelledby="title ? 'confirm-dialog-title' : undefined"
          :aria-describedby="'confirm-dialog-message'"
        >
          <!-- Header -->
          <div class="flex items-start gap-4 p-6 pb-4">
            <!-- Icon -->
            <div :class="['p-3 rounded-full flex-shrink-0', variantClasses.iconBg]">
              <component :is="variantClasses.icon" :class="['w-6 h-6', variantClasses.iconColor]" />
            </div>
            
            <!-- Content -->
            <div class="flex-1 min-w-0">
              <h3 
                v-if="title"
                id="confirm-dialog-title"
                class="text-lg font-semibold text-text-primary"
              >
                {{ title }}
              </h3>
              <p 
                id="confirm-dialog-message"
                class="mt-2 text-sm text-text-muted leading-relaxed"
              >
                {{ message }}
              </p>
            </div>
            
            <!-- Close button -->
            <button
              @click="handleCancel"
              :disabled="loading"
              class="p-1 text-text-muted hover:text-text-primary transition-colors rounded-lg hover:bg-background disabled:opacity-50 disabled:cursor-not-allowed flex-shrink-0"
            >
              <X class="w-5 h-5" />
            </button>
          </div>
          
          <!-- Actions -->
          <div class="flex items-center justify-end gap-3 px-6 py-4 bg-background/50 border-t border-border">
            <button
              @click="handleCancel"
              :disabled="loading"
              class="px-4 py-2 text-sm font-medium text-text-primary hover:bg-background rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ cancelText }}
            </button>
            <button
              @click="handleConfirm"
              :disabled="loading"
              :class="[
                'px-4 py-2 text-sm font-medium text-white rounded-lg transition-colors flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed',
                variantClasses.buttonBg
              ]"
            >
              <div v-if="loading" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
              {{ loading ? 'Please wait...' : confirmText }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .relative,
.modal-leave-to .relative {
  transform: scale(0.95);
}
</style>

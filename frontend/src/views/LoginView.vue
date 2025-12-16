<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Key, AlertCircle } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const password = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  if (!password.value) {
    error.value = 'Password is required'
    return
  }

  loading.value = true
  error.value = ''

  try {
    await authStore.login(password.value)
    router.push({ name: 'dashboard' })
  } catch (e) {
    error.value = 'Invalid password. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-background flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <!-- Logo -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-primary rounded-2xl mb-4">
          <Key class="w-8 h-8 text-white" />
        </div>
        <h1 class="text-2xl font-bold text-text-primary">Ory Kratos</h1>
        <p class="text-text-muted mt-1">Admin Console</p>
      </div>

      <!-- Login form -->
      <div class="bg-card border border-border rounded-2xl p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-1">Welcome back</h2>
        <p class="text-sm text-text-muted mb-6">Enter your admin password to continue</p>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <!-- Error message -->
          <div v-if="error" class="flex items-center gap-2 p-3 bg-danger/10 border border-danger/20 rounded-lg text-danger text-sm">
            <AlertCircle class="w-4 h-4 shrink-0" />
            {{ error }}
          </div>

          <!-- Password input -->
          <div>
            <label for="password" class="block text-sm font-medium text-text-secondary mb-2">
              Admin Password
            </label>
            <input
              id="password"
              v-model="password"
              type="password"
              placeholder="Enter your password"
              class="w-full px-4 py-3 bg-background border border-border rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:border-primary transition-colors"
              :disabled="loading"
            />
          </div>

          <!-- Submit button -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 bg-primary hover:bg-primary-hover text-white font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <div v-if="loading" class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            <span>{{ loading ? 'Signing in...' : 'Sign In' }}</span>
          </button>
        </form>
      </div>

      <!-- Footer -->
      <p class="text-center text-sm text-text-muted mt-6">
        Powered by <a href="https://www.ory.sh" target="_blank" class="text-primary hover:text-primary-light">Ory</a>
      </p>
    </div>
  </div>
</template>





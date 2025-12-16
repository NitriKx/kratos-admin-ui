import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('auth_token'))
  const expiresAt = ref<number | null>(
    localStorage.getItem('auth_expires_at') 
      ? parseInt(localStorage.getItem('auth_expires_at')!) 
      : null
  )

  const isAuthenticated = computed(() => {
    if (!token.value || !expiresAt.value) return false
    return Date.now() < expiresAt.value * 1000
  })

  async function login(password: string) {
    const response = await api.login(password)
    token.value = response.token
    expiresAt.value = response.expires_at
    localStorage.setItem('auth_token', response.token)
    localStorage.setItem('auth_expires_at', response.expires_at.toString())
  }

  function logout() {
    token.value = null
    expiresAt.value = null
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_expires_at')
  }

  return {
    token,
    expiresAt,
    isAuthenticated,
    login,
    logout
  }
})





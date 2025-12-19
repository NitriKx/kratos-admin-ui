import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/api/client'
import type { Identity, Session } from '@/types'

export const useIdentitiesStore = defineStore('identities', () => {
  const identities = ref<Identity[]>([])
  const currentIdentity = ref<Identity | null>(null)
  const currentIdentitySessions = ref<Session[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const page = ref(1)
  const perPage = ref(20)
  const total = ref(0)

  async function fetchIdentities(pageNum = 1, itemsPerPage = 20) {
    loading.value = true
    error.value = null
    try {
      const response = await api.getIdentities(pageNum, itemsPerPage)
      identities.value = response.data
      page.value = response.page
      perPage.value = response.per_page
      total.value = response.total || response.data.length
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch identities'
    } finally {
      loading.value = false
    }
  }

  async function fetchIdentity(id: string) {
    loading.value = true
    error.value = null
    try {
      currentIdentity.value = await api.getIdentity(id)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch identity'
    } finally {
      loading.value = false
    }
  }

  async function fetchIdentitySessions(id: string) {
    loading.value = true
    error.value = null
    try {
      const response = await api.getIdentitySessions(id)
      currentIdentitySessions.value = response.data
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch sessions'
    } finally {
      loading.value = false
    }
  }

  async function createIdentity(data: { schema_id: string; traits: Record<string, unknown>; state?: string }) {
    loading.value = true
    error.value = null
    try {
      const identity = await api.createIdentity(data)
      identities.value.unshift(identity)
      return identity
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create identity'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateIdentity(id: string, data: { schema_id: string; traits: Record<string, unknown>; state?: string }) {
    loading.value = true
    error.value = null
    try {
      const identity = await api.updateIdentity(id, data)
      const index = identities.value.findIndex(i => i.id === id)
      if (index !== -1) {
        identities.value[index] = identity
      }
      if (currentIdentity.value?.id === id) {
        currentIdentity.value = identity
      }
      return identity
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to update identity'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteIdentity(id: string) {
    loading.value = true
    error.value = null
    try {
      await api.deleteIdentity(id)
      identities.value = identities.value.filter(i => i.id !== id)
      if (currentIdentity.value?.id === id) {
        currentIdentity.value = null
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete identity'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchIdentityWithCredentials(id: string) {
    loading.value = true
    error.value = null
    try {
      currentIdentity.value = await api.getIdentityWithCredentials(id)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch identity'
    } finally {
      loading.value = false
    }
  }

  async function resetPassword(id: string, password: string) {
    loading.value = true
    error.value = null
    try {
      await api.resetPassword(id, password)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to reset password'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteCredential(id: string, credentialType: string) {
    loading.value = true
    error.value = null
    try {
      await api.deleteCredential(id, credentialType)
      // Refresh the identity to get updated credentials
      await fetchIdentityWithCredentials(id)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete credential'
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    identities,
    currentIdentity,
    currentIdentitySessions,
    loading,
    error,
    page,
    perPage,
    total,
    fetchIdentities,
    fetchIdentity,
    fetchIdentitySessions,
    fetchIdentityWithCredentials,
    createIdentity,
    updateIdentity,
    deleteIdentity,
    resetPassword,
    deleteCredential
  }
})





<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Key, Search, Filter, Trash2, Eye } from 'lucide-vue-next'
import { api } from '@/api/client'
import { useToast } from '@/composables/useToast'
import DataTable from '@/components/DataTable.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import type { Session } from '@/types'

const router = useRouter()
const toast = useToast()

const sessions = ref<Session[]>([])
const loading = ref(true)
const page = ref(1)
const perPage = ref(20)
const total = ref(0)
const searchQuery = ref('')

// Revoke confirmation dialog state
const isRevokeDialogOpen = ref(false)
const sessionToRevoke = ref<string | null>(null)
const isRevoking = ref(false)

const columns = [
  { key: 'id', label: 'Session ID', width: '200px' },
  { key: 'identity', label: 'Identity' },
  { key: 'status', label: 'Status', width: '100px' },
  { key: 'authenticated_at', label: 'Authenticated', width: '180px' },
  { key: 'expires_at', label: 'Expires', width: '180px' },
  { key: 'actions', label: 'Actions', width: '100px' }
]

const truncateId = (id: string) => {
  if (id.length <= 16) return id
  return `${id.substring(0, 8)}...${id.substring(id.length - 4)}`
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return {
    date: date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }),
    time: date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
  }
}

const getEmail = (session: Session) => {
  const traits = session.identity?.traits as Record<string, unknown>
  return traits?.email as string || traits?.username as string || 'Unknown'
}

const fetchSessions = async (pageNum = 1) => {
  loading.value = true
  try {
    const response = await api.getSessions(pageNum, perPage.value)
    sessions.value = response.data
    page.value = response.page
    total.value = response.total || response.data.length
  } catch (e) {
    console.error('Failed to fetch sessions', e)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (newPage: number) => {
  fetchSessions(newPage)
}

const viewIdentity = (identityId: string) => {
  router.push({ name: 'identity-detail', params: { id: identityId } })
}

const openRevokeDialog = (sessionId: string) => {
  sessionToRevoke.value = sessionId
  isRevokeDialogOpen.value = true
}

const closeRevokeDialog = () => {
  isRevokeDialogOpen.value = false
  sessionToRevoke.value = null
}

const confirmRevoke = async () => {
  if (!sessionToRevoke.value) return
  
  isRevoking.value = true
  try {
    await api.revokeSession(sessionToRevoke.value)
    toast.success('Session revoked', 'The session has been revoked successfully.')
    await fetchSessions(page.value)
    closeRevokeDialog()
  } catch (e) {
    console.error('Failed to revoke session', e)
    toast.error('Failed to revoke session', 'An error occurred while revoking the session.')
  } finally {
    isRevoking.value = false
  }
}

onMounted(() => {
  fetchSessions()
})
</script>

<template>
  <div>
    <!-- Breadcrumb -->
    <div class="mb-2">
      <span class="text-sm text-text-muted">User Management</span>
      <span class="text-sm text-text-muted mx-2">/</span>
      <span class="text-sm text-text-primary font-medium">Sessions</span>
    </div>

    <!-- Page header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-text-primary">Sessions</h1>
        <p class="text-text-muted mt-1">View and manage active user sessions.</p>
      </div>
    </div>

    <!-- Search and filters -->
    <div class="flex items-center gap-4 mb-6">
      <div class="flex-1 relative">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-text-muted" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search sessions..."
          class="w-full pl-10 pr-4 py-2.5 bg-card border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-primary transition-colors"
        />
      </div>
      <button class="flex items-center gap-2 px-4 py-2.5 bg-card border border-border rounded-lg text-sm text-text-secondary hover:text-text-primary hover:bg-background transition-colors">
        <Filter class="w-4 h-4" />
        Filter
      </button>
    </div>

    <!-- Data table -->
    <DataTable
      :columns="columns"
      :data="sessions"
      :loading="loading"
      :page="page"
      :per-page="perPage"
      :total="total"
      @page-change="handlePageChange"
    >
      <template #cell-id="{ row }">
        <span class="font-mono text-sm text-text-muted">{{ truncateId(row.id) }}</span>
      </template>

      <template #cell-identity="{ row }">
        <div class="flex items-center gap-3">
          <UserAvatar :email="getEmail(row)" size="sm" />
          <div>
            <p class="text-sm font-medium text-text-primary">{{ getEmail(row) }}</p>
            <p class="text-xs text-text-muted font-mono">{{ truncateId(row.identity?.id || '') }}</p>
          </div>
        </div>
      </template>

      <template #cell-status="{ row }">
        <StatusBadge :status="row.active ? 'active' : 'inactive'" />
      </template>

      <template #cell-authenticated_at="{ row }">
        <div>
          <p class="text-sm text-text-primary">{{ formatDate(row.authenticated_at).date }}</p>
          <p class="text-xs text-text-muted">{{ formatDate(row.authenticated_at).time }}</p>
        </div>
      </template>

      <template #cell-expires_at="{ row }">
        <div>
          <p class="text-sm text-text-primary">{{ formatDate(row.expires_at).date }}</p>
          <p class="text-xs text-text-muted">{{ formatDate(row.expires_at).time }}</p>
        </div>
      </template>

      <template #cell-actions="{ row }">
        <div class="flex items-center gap-1">
          <button
            @click="viewIdentity(row.identity?.id)"
            class="p-1.5 text-text-muted hover:text-text-primary transition-colors"
            title="View identity"
          >
            <Eye class="w-4 h-4" />
          </button>
          <button
            @click="openRevokeDialog(row.id)"
            class="p-1.5 text-text-muted hover:text-danger transition-colors"
            title="Revoke session"
          >
            <Trash2 class="w-4 h-4" />
          </button>
        </div>
      </template>
    </DataTable>

    <!-- Revoke Confirmation Dialog -->
    <ConfirmDialog
      :is-open="isRevokeDialogOpen"
      title="Revoke Session"
      message="Are you sure you want to revoke this session? The user will be logged out immediately."
      confirm-text="Revoke"
      cancel-text="Cancel"
      variant="warning"
      :loading="isRevoking"
      @confirm="confirmRevoke"
      @cancel="closeRevokeDialog"
    />
  </div>
</template>




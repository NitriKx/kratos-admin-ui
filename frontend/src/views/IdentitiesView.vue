<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { AxiosError } from 'axios'
import { Plus, Search, Pencil, Eye, Trash2, UserCheck, UserX } from 'lucide-vue-next'
import { useIdentitiesStore } from '@/stores/identities'
import { useToast } from '@/composables/useToast'
import DataTable from '@/components/DataTable.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import IdentityModal from '@/components/IdentityModal.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import type { Identity } from '@/types'

const router = useRouter()
const identitiesStore = useIdentitiesStore()
const toast = useToast()

// Helper to extract error message from API errors
function getErrorMessage(error: unknown): string {
  if (error instanceof AxiosError) {
    // Try to get the error message from the API response
    const data = error.response?.data
    if (data?.details) return data.details
    if (data?.error) return data.error
    if (data?.message) return data.message
    if (error.message) return error.message
  }
  if (error instanceof Error) {
    return error.message
  }
  return 'An unexpected error occurred'
}

const searchQuery = ref('')
const statusFilter = ref('all')

// Filter identities based on search query and status
const filteredIdentities = computed(() => {
  let result = identitiesStore.identities

  // Filter by status
  if (statusFilter.value !== 'all') {
    result = result.filter(identity => identity.state === statusFilter.value)
  }

  // Filter by search query
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase().trim()
    result = result.filter(identity => {
      // Search in ID
      if (identity.id.toLowerCase().includes(query)) return true

      // Search in traits (email, username, name, etc.)
      const traits = identity.traits as Record<string, unknown>
      if (traits) {
        // Check email
        if (typeof traits.email === 'string' && traits.email.toLowerCase().includes(query)) return true
        // Check username
        if (typeof traits.username === 'string' && traits.username.toLowerCase().includes(query)) return true
        // Check name
        if (traits.name && typeof traits.name === 'object') {
          const name = traits.name as Record<string, string>
          const fullName = `${name.first || ''} ${name.last || ''}`.toLowerCase()
          if (fullName.includes(query)) return true
        }
        // Check schema_id
        if (identity.schema_id.toLowerCase().includes(query)) return true
      }

      return false
    })
  }

  return result
})

// Modal state
const isModalOpen = ref(false)
const editingIdentity = ref<Identity | null>(null)
const isSaving = ref(false)

// Delete confirmation dialog state
const isDeleteDialogOpen = ref(false)
const identityToDelete = ref<string | null>(null)
const isDeleting = ref(false)

// Toggle state tracking
const togglingIdentityId = ref<string | null>(null)

const columns = [
  { key: 'id', label: 'ID / UUID', width: '180px' },
  { key: 'traits', label: 'Traits' },
  { key: 'state', label: 'State', width: '120px' },
  { key: 'created_at', label: 'Created At', width: '150px' },
  { key: 'actions', label: 'Actions', width: '150px' }
]

const truncateId = (id: string) => {
  if (id.length <= 12) return id
  return `${id.substring(0, 6)}...${id.substring(id.length - 4)}`
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return {
    date: date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }),
    time: date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
  }
}

const getEmail = (identity: Identity) => {
  const traits = identity.traits as Record<string, unknown>
  return traits?.email as string || traits?.username as string || 'Unknown'
}

const getName = (identity: Identity): string | undefined => {
  const traits = identity.traits as Record<string, unknown>
  if (traits?.name) {
    const name = traits.name as Record<string, string>
    return `${name.first || ''} ${name.last || ''}`.trim() || undefined
  }
  return undefined
}

const handlePageChange = (page: number) => {
  identitiesStore.fetchIdentities(page, identitiesStore.perPage)
}

const viewIdentity = (id: string) => {
  router.push({ name: 'identity-detail', params: { id } })
}

const openCreateModal = () => {
  editingIdentity.value = null
  isModalOpen.value = true
}

const openEditModal = (identity: Identity) => {
  editingIdentity.value = identity
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
  editingIdentity.value = null
}

const handleSaveIdentity = async (data: { schema_id: string; traits: Record<string, unknown>; state: string }) => {
  isSaving.value = true
  try {
    if (editingIdentity.value) {
      await identitiesStore.updateIdentity(editingIdentity.value.id, data)
      toast.success('Identity updated', 'The identity has been updated successfully.')
    } else {
      await identitiesStore.createIdentity(data)
      toast.success('Identity created', 'The new identity has been created successfully.')
    }
    closeModal()
  } catch (e) {
    console.error('Failed to save identity', e)
    const action = editingIdentity.value ? 'update' : 'create'
    toast.error(`Failed to ${action} identity`, getErrorMessage(e))
  } finally {
    isSaving.value = false
  }
}

const openDeleteDialog = (id: string) => {
  identityToDelete.value = id
  isDeleteDialogOpen.value = true
}

const closeDeleteDialog = () => {
  isDeleteDialogOpen.value = false
  identityToDelete.value = null
}

const confirmDelete = async () => {
  if (!identityToDelete.value) return
  
  isDeleting.value = true
  try {
    await identitiesStore.deleteIdentity(identityToDelete.value)
    toast.success('Identity deleted', 'The identity has been deleted successfully.')
    closeDeleteDialog()
  } catch (e) {
    console.error('Failed to delete identity', e)
    toast.error('Failed to delete identity', getErrorMessage(e))
  } finally {
    isDeleting.value = false
  }
}

const toggleIdentityState = async (identity: Identity) => {
  togglingIdentityId.value = identity.id
  const newState = identity.state === 'active' ? 'inactive' : 'active'
  try {
    await identitiesStore.updateIdentity(identity.id, {
      schema_id: identity.schema_id,
      traits: identity.traits as Record<string, unknown>,
      state: newState
    })
    const action = newState === 'active' ? 'enabled' : 'disabled'
    toast.success(`Identity ${action}`, `The identity has been ${action} successfully.`)
  } catch (e) {
    console.error('Failed to toggle identity state', e)
    toast.error('Failed to update identity', getErrorMessage(e))
  } finally {
    togglingIdentityId.value = null
  }
}

onMounted(async () => {
  await identitiesStore.fetchIdentities()
})
</script>

<template>
  <div>
    <!-- Breadcrumb -->
    <div class="mb-2">
      <span class="text-sm text-text-muted">User Management</span>
      <span class="text-sm text-text-muted mx-2">/</span>
      <span class="text-sm text-text-primary font-medium">Identities</span>
    </div>

    <!-- Page header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-text-primary">User Management</h1>
        <p class="text-text-muted mt-1">Manage user identities, credentials, and sessions.</p>
      </div>
      <button
        @click="openCreateModal"
        class="flex items-center gap-2 px-4 py-2 bg-primary hover:bg-primary-hover text-white font-medium rounded-lg transition-colors"
      >
        <Plus class="w-4 h-4" />
        Create Identity
      </button>
    </div>

    <!-- Search and filters -->
    <div class="flex items-center gap-4 mb-6">
      <div class="flex-1 relative">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-text-muted" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Filter by ID, email, or traits..."
          class="w-full pl-10 pr-4 py-2.5 bg-card border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-primary transition-colors"
        />
      </div>
      <select
        v-model="statusFilter"
        class="px-4 py-2.5 bg-card border border-border rounded-lg text-sm text-text-primary focus:outline-none focus:border-primary transition-colors"
      >
        <option value="all">All States</option>
        <option value="active">Active</option>
        <option value="inactive">Inactive</option>
      </select>
    </div>

    <!-- Data table -->
    <DataTable
      :columns="columns"
      :data="filteredIdentities"
      :loading="identitiesStore.loading"
      :page="identitiesStore.page"
      :per-page="identitiesStore.perPage"
      :total="filteredIdentities.length"
      @page-change="handlePageChange"
    >
      <template #cell-id="{ row }">
        <span class="font-mono text-sm text-text-muted">{{ truncateId(row.id) }}</span>
      </template>

      <template #cell-traits="{ row }">
        <div class="flex items-center gap-3">
          <UserAvatar :email="getEmail(row)" :name="getName(row)" size="sm" />
          <div>
            <p class="text-sm font-medium text-text-primary">{{ getEmail(row) }}</p>
            <p class="text-xs text-text-muted">Schema: {{ row.schema_id }}</p>
          </div>
        </div>
      </template>

      <template #cell-state="{ row }">
        <StatusBadge :status="row.state" />
      </template>

      <template #cell-created_at="{ row }">
        <div>
          <p class="text-sm text-text-primary">{{ formatDate(row.created_at).date }}</p>
          <p class="text-xs text-text-muted">{{ formatDate(row.created_at).time }}</p>
        </div>
      </template>

      <template #cell-actions="{ row }">
        <div class="flex items-center gap-1">
          <button
            @click="openEditModal(row)"
            class="p-1.5 text-text-muted hover:text-text-primary transition-colors"
            title="Edit"
          >
            <Pencil class="w-4 h-4" />
          </button>
          <button
            @click="viewIdentity(row.id)"
            class="p-1.5 text-text-muted hover:text-text-primary transition-colors"
            title="View details"
          >
            <Eye class="w-4 h-4" />
          </button>
          <button
            @click="toggleIdentityState(row)"
            class="p-1.5 transition-colors"
            :class="row.state === 'active' ? 'text-text-muted hover:text-warning' : 'text-text-muted hover:text-success'"
            :title="row.state === 'active' ? 'Disable user' : 'Enable user'"
            :disabled="togglingIdentityId === row.id"
          >
            <div v-if="togglingIdentityId === row.id" class="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin"></div>
            <UserX v-else-if="row.state === 'active'" class="w-4 h-4" />
            <UserCheck v-else class="w-4 h-4" />
          </button>
          <button
            @click="openDeleteDialog(row.id)"
            class="p-1.5 text-text-muted hover:text-danger transition-colors"
            title="Delete"
          >
            <Trash2 class="w-4 h-4" />
          </button>
        </div>
      </template>
    </DataTable>

    <!-- Identity Modal -->
    <IdentityModal
      :is-open="isModalOpen"
      :identity="editingIdentity"
      :saving="isSaving"
      @close="closeModal"
      @save="handleSaveIdentity"
    />

    <!-- Delete Confirmation Dialog -->
    <ConfirmDialog
      :is-open="isDeleteDialogOpen"
      title="Delete Identity"
      message="Are you sure you want to delete this identity? This action cannot be undone and will remove all associated data."
      confirm-text="Delete"
      cancel-text="Cancel"
      variant="danger"
      :loading="isDeleting"
      @confirm="confirmDelete"
      @cancel="closeDeleteDialog"
    />
  </div>
</template>


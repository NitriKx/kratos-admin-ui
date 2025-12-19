<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Pencil, Trash2, Key, Save, X, KeyRound, ShieldOff, Shield, Eye, EyeOff } from 'lucide-vue-next'
import { useIdentitiesStore } from '@/stores/identities'
import { useToast } from '@/composables/useToast'
import { getErrorMessage } from '@/utils/errorMessage'
import StatusBadge from '@/components/StatusBadge.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

const route = useRoute()
const router = useRouter()
const identitiesStore = useIdentitiesStore()
const toast = useToast()

const isEditing = ref(false)
const editedTraits = ref<string>('')

// Delete confirmation dialog state
const isDeleteDialogOpen = ref(false)
const isDeleting = ref(false)

// Password reset state
const isPasswordModalOpen = ref(false)
const newPassword = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const isResettingPassword = ref(false)

// Delete credential (MFA) state
const isDeleteCredentialDialogOpen = ref(false)
const credentialToDelete = ref<string | null>(null)
const isDeletingCredential = ref(false)

const identity = computed(() => identitiesStore.currentIdentity)
const sessions = computed(() => identitiesStore.currentIdentitySessions)

// Check if identity has specific credential types
const hasTotp = computed(() => identity.value?.credentials?.totp !== undefined)
const hasWebauthn = computed(() => identity.value?.credentials?.webauthn !== undefined)
const hasLookupSecret = computed(() => identity.value?.credentials?.lookup_secret !== undefined)
const hasMfa = computed(() => hasTotp.value || hasWebauthn.value || hasLookupSecret.value)

const getEmail = () => {
  if (!identity.value) return 'Unknown'
  const traits = identity.value.traits as Record<string, unknown>
  return traits?.email as string || traits?.username as string || 'Unknown'
}

const getName = (): string | undefined => {
  if (!identity.value) return undefined
  const traits = identity.value.traits as Record<string, unknown>
  if (traits?.name) {
    const name = traits.name as Record<string, string>
    return `${name.first || ''} ${name.last || ''}`.trim() || undefined
  }
  return undefined
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const startEditing = () => {
  if (identity.value) {
    editedTraits.value = JSON.stringify(identity.value.traits, null, 2)
    isEditing.value = true
  }
}

const cancelEditing = () => {
  isEditing.value = false
  editedTraits.value = ''
}

const saveChanges = async () => {
  if (!identity.value) return
  try {
    const traits = JSON.parse(editedTraits.value)
    await identitiesStore.updateIdentity(identity.value.id, {
      schema_id: identity.value.schema_id,
      traits,
      state: identity.value.state
    })
    isEditing.value = false
  } catch (e) {
    alert('Invalid JSON or failed to save changes')
  }
}

const openDeleteDialog = () => {
  isDeleteDialogOpen.value = true
}

const closeDeleteDialog = () => {
  isDeleteDialogOpen.value = false
}

const confirmDelete = async () => {
  if (!identity.value) return
  
  isDeleting.value = true
  try {
    await identitiesStore.deleteIdentity(identity.value.id)
    toast.success('Identity deleted', 'The identity has been deleted successfully.')
    router.push({ name: 'identities' })
  } catch (e) {
    console.error('Failed to delete identity', e)
    toast.error('Failed to delete identity', 'An error occurred while deleting the identity.')
  } finally {
    isDeleting.value = false
    closeDeleteDialog()
  }
}

// Password reset functions
const openPasswordModal = () => {
  newPassword.value = ''
  confirmPassword.value = ''
  showPassword.value = false
  isPasswordModalOpen.value = true
}

const closePasswordModal = () => {
  isPasswordModalOpen.value = false
  newPassword.value = ''
  confirmPassword.value = ''
}

const passwordsMatch = computed(() => newPassword.value === confirmPassword.value)
const passwordValid = computed(() => newPassword.value.length >= 8 && passwordsMatch.value)

const resetPassword = async () => {
  if (!identity.value || !passwordValid.value) return
  
  isResettingPassword.value = true
  try {
    await identitiesStore.resetPassword(identity.value.id, newPassword.value)
    toast.success('Password reset', 'The password has been reset successfully.')
    closePasswordModal()
  } catch (e) {
    console.error('Failed to reset password', e)
    toast.error('Failed to reset password', getErrorMessage(e))
  } finally {
    isResettingPassword.value = false
  }
}

// MFA/Credential deletion functions
const openDeleteCredentialDialog = (credType: string) => {
  credentialToDelete.value = credType
  isDeleteCredentialDialogOpen.value = true
}

const closeDeleteCredentialDialog = () => {
  isDeleteCredentialDialogOpen.value = false
  credentialToDelete.value = null
}

const getCredentialDisplayName = (type: string) => {
  const names: Record<string, string> = {
    totp: 'TOTP (Authenticator App)',
    webauthn: 'WebAuthn (Security Key)',
    lookup_secret: 'Backup Codes'
  }
  return names[type] || type
}

const confirmDeleteCredential = async () => {
  if (!identity.value || !credentialToDelete.value) return
  
  isDeletingCredential.value = true
  try {
    await identitiesStore.deleteCredential(identity.value.id, credentialToDelete.value)
    toast.success('MFA disabled', `${getCredentialDisplayName(credentialToDelete.value)} has been removed.`)
    closeDeleteCredentialDialog()
  } catch (e) {
    console.error('Failed to delete credential', e)
    toast.error('Failed to disable MFA', getErrorMessage(e))
  } finally {
    isDeletingCredential.value = false
  }
}

onMounted(async () => {
  const id = route.params.id as string
  await identitiesStore.fetchIdentityWithCredentials(id)
  await identitiesStore.fetchIdentitySessions(id)
})
</script>

<template>
  <div>
    <!-- Breadcrumb -->
    <div class="mb-2">
      <span class="text-sm text-text-muted">User Management</span>
      <span class="text-sm text-text-muted mx-2">/</span>
      <router-link :to="{ name: 'identities' }" class="text-sm text-text-muted hover:text-primary">Identities</router-link>
      <span class="text-sm text-text-muted mx-2">/</span>
      <span class="text-sm text-text-primary font-medium">{{ identity?.id?.substring(0, 8) }}...</span>
    </div>

    <!-- Back button -->
    <button
      @click="router.push({ name: 'identities' })"
      class="flex items-center gap-2 text-text-secondary hover:text-text-primary mb-4 transition-colors"
    >
      <ArrowLeft class="w-4 h-4" />
      Back to Identities
    </button>

    <!-- Loading state -->
    <div v-if="identitiesStore.loading && !identity" class="bg-card border border-border rounded-xl p-8 text-center">
      <div class="flex items-center justify-center gap-2 text-text-muted">
        <div class="w-5 h-5 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
        <span>Loading identity...</span>
      </div>
    </div>

    <!-- Error state -->
    <div v-else-if="identitiesStore.error" class="bg-danger/10 border border-danger/20 rounded-xl p-6 text-center">
      <p class="text-danger">{{ identitiesStore.error }}</p>
    </div>

    <!-- Identity details -->
    <div v-else-if="identity" class="space-y-6">
      <!-- Header card -->
      <div class="bg-card border border-border rounded-xl p-6">
        <div class="flex items-start justify-between">
          <div class="flex items-center gap-4">
            <UserAvatar :email="getEmail()" :name="getName()" size="lg" />
            <div>
              <h1 class="text-xl font-bold text-text-primary">{{ getEmail() }}</h1>
              <p class="text-sm text-text-muted font-mono">{{ identity.id }}</p>
              <div class="flex items-center gap-2 mt-2">
                <StatusBadge :status="identity.state" />
                <span class="text-xs text-text-muted">Schema: {{ identity.schema_id }}</span>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <button
              @click="openPasswordModal"
              class="flex items-center gap-2 px-3 py-2 bg-card border border-border rounded-lg text-sm text-text-secondary hover:text-text-primary hover:bg-background transition-colors"
              title="Reset Password"
            >
              <KeyRound class="w-4 h-4" />
              Reset Password
            </button>
            <button
              v-if="!isEditing"
              @click="startEditing"
              class="flex items-center gap-2 px-3 py-2 bg-card border border-border rounded-lg text-sm text-text-secondary hover:text-text-primary hover:bg-background transition-colors"
            >
              <Pencil class="w-4 h-4" />
              Edit
            </button>
            <button
              @click="openDeleteDialog"
              class="flex items-center gap-2 px-3 py-2 bg-danger/10 border border-danger/20 rounded-lg text-sm text-danger hover:bg-danger/20 transition-colors"
            >
              <Trash2 class="w-4 h-4" />
              Delete
            </button>
          </div>
        </div>
      </div>

      <!-- Traits -->
      <div class="bg-card border border-border rounded-xl p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-text-primary">Traits</h2>
          <div v-if="isEditing" class="flex items-center gap-2">
            <button
              @click="cancelEditing"
              class="flex items-center gap-1 px-3 py-1.5 text-sm text-text-secondary hover:text-text-primary transition-colors"
            >
              <X class="w-4 h-4" />
              Cancel
            </button>
            <button
              @click="saveChanges"
              class="flex items-center gap-1 px-3 py-1.5 bg-primary hover:bg-primary-hover text-white text-sm rounded-lg transition-colors"
            >
              <Save class="w-4 h-4" />
              Save
            </button>
          </div>
        </div>
        <div v-if="isEditing">
          <textarea
            v-model="editedTraits"
            class="w-full h-64 p-4 bg-background border border-border rounded-lg font-mono text-sm text-text-primary focus:outline-none focus:border-primary transition-colors"
          ></textarea>
        </div>
        <pre v-else class="p-4 bg-background rounded-lg overflow-auto text-sm text-text-primary font-mono">{{ JSON.stringify(identity.traits, null, 2) }}</pre>
      </div>

      <!-- Security / MFA -->
      <div class="bg-card border border-border rounded-xl p-6">
        <div class="flex items-center gap-2 mb-4">
          <Shield class="w-5 h-5 text-text-muted" />
          <h2 class="text-lg font-semibold text-text-primary">Security & MFA</h2>
        </div>
        
        <div class="space-y-3">
          <!-- Password -->
          <div class="flex items-center justify-between p-4 bg-background rounded-lg">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
                <KeyRound class="w-5 h-5 text-primary" />
              </div>
              <div>
                <p class="text-sm font-medium text-text-primary">Password</p>
                <p class="text-xs text-text-muted">Primary authentication method</p>
              </div>
            </div>
            <button
              @click="openPasswordModal"
              class="px-3 py-1.5 text-sm text-primary hover:bg-primary/10 rounded-lg transition-colors"
            >
              Reset
            </button>
          </div>

          <!-- TOTP -->
          <div class="flex items-center justify-between p-4 bg-background rounded-lg">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center" :class="hasTotp ? 'bg-success/10' : 'bg-background border border-border'">
                <Shield class="w-5 h-5" :class="hasTotp ? 'text-success' : 'text-text-muted'" />
              </div>
              <div>
                <p class="text-sm font-medium text-text-primary">TOTP (Authenticator App)</p>
                <p class="text-xs" :class="hasTotp ? 'text-success' : 'text-text-muted'">
                  {{ hasTotp ? 'Enabled' : 'Not configured' }}
                </p>
              </div>
            </div>
            <button
              v-if="hasTotp"
              @click="openDeleteCredentialDialog('totp')"
              class="flex items-center gap-1 px-3 py-1.5 text-sm text-danger hover:bg-danger/10 rounded-lg transition-colors"
            >
              <ShieldOff class="w-4 h-4" />
              Disable
            </button>
            <span v-else class="text-xs text-text-muted">—</span>
          </div>

          <!-- WebAuthn -->
          <div class="flex items-center justify-between p-4 bg-background rounded-lg">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center" :class="hasWebauthn ? 'bg-success/10' : 'bg-background border border-border'">
                <Key class="w-5 h-5" :class="hasWebauthn ? 'text-success' : 'text-text-muted'" />
              </div>
              <div>
                <p class="text-sm font-medium text-text-primary">WebAuthn (Security Key)</p>
                <p class="text-xs" :class="hasWebauthn ? 'text-success' : 'text-text-muted'">
                  {{ hasWebauthn ? 'Enabled' : 'Not configured' }}
                </p>
              </div>
            </div>
            <button
              v-if="hasWebauthn"
              @click="openDeleteCredentialDialog('webauthn')"
              class="flex items-center gap-1 px-3 py-1.5 text-sm text-danger hover:bg-danger/10 rounded-lg transition-colors"
            >
              <ShieldOff class="w-4 h-4" />
              Disable
            </button>
            <span v-else class="text-xs text-text-muted">—</span>
          </div>

          <!-- Lookup Secret (Backup Codes) -->
          <div class="flex items-center justify-between p-4 bg-background rounded-lg">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center" :class="hasLookupSecret ? 'bg-success/10' : 'bg-background border border-border'">
                <Key class="w-5 h-5" :class="hasLookupSecret ? 'text-success' : 'text-text-muted'" />
              </div>
              <div>
                <p class="text-sm font-medium text-text-primary">Backup Codes</p>
                <p class="text-xs" :class="hasLookupSecret ? 'text-success' : 'text-text-muted'">
                  {{ hasLookupSecret ? 'Configured' : 'Not configured' }}
                </p>
              </div>
            </div>
            <button
              v-if="hasLookupSecret"
              @click="openDeleteCredentialDialog('lookup_secret')"
              class="flex items-center gap-1 px-3 py-1.5 text-sm text-danger hover:bg-danger/10 rounded-lg transition-colors"
            >
              <ShieldOff class="w-4 h-4" />
              Remove
            </button>
            <span v-else class="text-xs text-text-muted">—</span>
          </div>
        </div>
      </div>

      <!-- Sessions -->
      <div class="bg-card border border-border rounded-xl p-6">
        <div class="flex items-center gap-2 mb-4">
          <Key class="w-5 h-5 text-text-muted" />
          <h2 class="text-lg font-semibold text-text-primary">Active Sessions</h2>
          <span class="px-2 py-0.5 bg-primary/10 text-primary text-xs font-medium rounded-full">
            {{ sessions.length }}
          </span>
        </div>
        <div v-if="sessions.length === 0" class="text-center py-8 text-text-muted">
          No active sessions
        </div>
        <div v-else class="space-y-3">
          <div
            v-for="session in sessions"
            :key="session.id"
            class="p-4 bg-background rounded-lg"
          >
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-mono text-text-muted">{{ session.id.substring(0, 16) }}...</p>
                <p class="text-xs text-text-muted mt-1">
                  Authenticated: {{ formatDate(session.authenticated_at) }}
                </p>
              </div>
              <StatusBadge :status="session.active ? 'active' : 'inactive'" />
            </div>
          </div>
        </div>
      </div>

      <!-- Metadata -->
      <div class="grid grid-cols-2 gap-6">
        <div class="bg-card border border-border rounded-xl p-6">
          <h2 class="text-lg font-semibold text-text-primary mb-4">Created</h2>
          <p class="text-text-primary">{{ formatDate(identity.created_at) }}</p>
        </div>
        <div class="bg-card border border-border rounded-xl p-6">
          <h2 class="text-lg font-semibold text-text-primary mb-4">Last Updated</h2>
          <p class="text-text-primary">{{ formatDate(identity.updated_at) }}</p>
        </div>
      </div>
    </div>

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

    <!-- Delete Credential Confirmation Dialog -->
    <ConfirmDialog
      :is-open="isDeleteCredentialDialogOpen"
      title="Disable MFA"
      :message="`Are you sure you want to disable ${credentialToDelete ? getCredentialDisplayName(credentialToDelete) : 'this MFA method'}? The user will need to set it up again.`"
      confirm-text="Disable"
      cancel-text="Cancel"
      variant="danger"
      :loading="isDeletingCredential"
      @confirm="confirmDeleteCredential"
      @cancel="closeDeleteCredentialDialog"
    />

    <!-- Password Reset Modal -->
    <Teleport to="body">
      <div
        v-if="isPasswordModalOpen"
        class="fixed inset-0 z-50 flex items-center justify-center"
      >
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-black/50"
          @click="closePasswordModal"
        ></div>

        <!-- Modal -->
        <div class="relative bg-card border border-border rounded-xl shadow-xl w-full max-w-md mx-4 p-6">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
              <KeyRound class="w-5 h-5 text-primary" />
            </div>
            <div>
              <h2 class="text-lg font-semibold text-text-primary">Reset Password</h2>
              <p class="text-sm text-text-muted">Set a new password for this user</p>
            </div>
          </div>

          <form @submit.prevent="resetPassword" class="space-y-4">
            <!-- New Password -->
            <div>
              <label class="block text-sm font-medium text-text-secondary mb-2">
                New Password
              </label>
              <div class="relative">
                <input
                  v-model="newPassword"
                  :type="showPassword ? 'text' : 'password'"
                  placeholder="Enter new password"
                  class="w-full px-4 py-3 pr-12 bg-background border border-border rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:border-primary transition-colors"
                  :class="{ 'border-danger': newPassword.length > 0 && newPassword.length < 8 }"
                />
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-primary transition-colors"
                >
                  <Eye v-if="!showPassword" class="w-5 h-5" />
                  <EyeOff v-else class="w-5 h-5" />
                </button>
              </div>
              <p v-if="newPassword.length > 0 && newPassword.length < 8" class="text-xs text-danger mt-1">
                Password must be at least 8 characters
              </p>
            </div>

            <!-- Confirm Password -->
            <div>
              <label class="block text-sm font-medium text-text-secondary mb-2">
                Confirm Password
              </label>
              <input
                v-model="confirmPassword"
                :type="showPassword ? 'text' : 'password'"
                placeholder="Confirm new password"
                class="w-full px-4 py-3 bg-background border border-border rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:border-primary transition-colors"
                :class="{ 'border-danger': confirmPassword.length > 0 && !passwordsMatch }"
              />
              <p v-if="confirmPassword.length > 0 && !passwordsMatch" class="text-xs text-danger mt-1">
                Passwords do not match
              </p>
            </div>

            <!-- Actions -->
            <div class="flex justify-end gap-3 pt-2">
              <button
                type="button"
                @click="closePasswordModal"
                class="px-4 py-2 text-sm text-text-secondary hover:text-text-primary transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="!passwordValid || isResettingPassword"
                class="flex items-center gap-2 px-4 py-2 bg-primary hover:bg-primary-hover text-white text-sm font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <div v-if="isResettingPassword" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                <span>{{ isResettingPassword ? 'Resetting...' : 'Reset Password' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>


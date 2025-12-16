<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Users, Key } from 'lucide-vue-next'
import { api } from '@/api/client'
import StatsCard from '@/components/StatsCard.vue'
import type { Stats } from '@/types'

const stats = ref<Stats | null>(null)
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    stats.value = await api.getStats()
  } catch (e) {
    error.value = 'Failed to load statistics'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <!-- Breadcrumb -->
    <div class="mb-2">
      <span class="text-sm text-text-primary font-medium">Dashboard</span>
    </div>

    <!-- Page header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-text-primary">Dashboard</h1>
        <p class="text-text-muted mt-1">Overview of your Kratos deployment</p>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
      <div v-for="i in 2" :key="i" class="bg-card border border-border rounded-xl p-5 animate-pulse">
        <div class="h-4 bg-border rounded w-24 mb-3"></div>
        <div class="h-8 bg-border rounded w-32"></div>
      </div>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="bg-danger/10 border border-danger/20 rounded-xl p-6 text-center">
      <p class="text-danger">{{ error }}</p>
    </div>

    <!-- Stats cards -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
      <StatsCard
        title="Active Users"
        :value="stats?.active_identities?.toLocaleString() ?? '0'"
        :icon="Users"
      />
      <StatsCard
        title="Active Sessions"
        :value="stats?.active_sessions?.toLocaleString() ?? '0'"
        :icon="Key"
      />
    </div>

    <!-- Quick actions -->
    <div class="bg-card border border-border rounded-xl p-6">
      <h2 class="text-lg font-semibold text-text-primary mb-4">Quick Actions</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <router-link
          :to="{ name: 'identities' }"
          class="p-4 bg-background border border-border rounded-lg hover:border-primary transition-colors group"
        >
          <Users class="w-8 h-8 text-primary mb-2" />
          <h3 class="font-medium text-text-primary group-hover:text-primary transition-colors">Manage Identities</h3>
          <p class="text-sm text-text-muted mt-1">View and manage user accounts</p>
        </router-link>
        <router-link
          :to="{ name: 'sessions' }"
          class="p-4 bg-background border border-border rounded-lg hover:border-primary transition-colors group"
        >
          <Key class="w-8 h-8 text-primary mb-2" />
          <h3 class="font-medium text-text-primary group-hover:text-primary transition-colors">View Sessions</h3>
          <p class="text-sm text-text-muted mt-1">Monitor and revoke active sessions</p>
        </router-link>
        <router-link
          :to="{ name: 'schemas' }"
          class="p-4 bg-background border border-border rounded-lg hover:border-primary transition-colors group"
        >
          <UserPlus class="w-8 h-8 text-primary mb-2" />
          <h3 class="font-medium text-text-primary group-hover:text-primary transition-colors">Identity Schemas</h3>
          <p class="text-sm text-text-muted mt-1">View configured identity schemas</p>
        </router-link>
      </div>
    </div>
  </div>
</template>




<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  LayoutDashboard,
  Users,
  Key,
  Clock,
  Settings,
  FileCode,
  LogOut
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const navigation = [
  {
    title: 'USER MANAGEMENT',
    items: [
      { name: 'Identities', icon: Users, route: 'identities' },
      { name: 'Sessions', icon: Key, route: 'sessions' }
    ]
  },
  {
    title: 'SYSTEM',
    items: [
      { name: 'Schemas', icon: FileCode, route: 'schemas' }
    ]
  }
]

const isActive = (routeName: string) => {
  return route.name === routeName
}

const handleLogout = () => {
  authStore.logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <aside class="fixed left-0 top-0 h-screen w-60 bg-sidebar flex flex-col border-r border-border">
    <!-- Logo -->
    <div class="p-4 border-b border-border">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
          <Key class="w-4 h-4 text-white" />
        </div>
        <div>
          <h1 class="text-sm font-semibold text-text-primary">Ory Kratos</h1>
          <p class="text-xs text-text-muted">Admin Console</p>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto p-4">
      <!-- Dashboard -->
      <router-link
        :to="{ name: 'dashboard' }"
        class="flex items-center gap-3 px-3 py-2 rounded-lg mb-4 transition-colors"
        :class="isActive('dashboard') ? 'bg-primary text-white' : 'text-text-secondary hover:bg-card hover:text-text-primary'"
      >
        <LayoutDashboard class="w-5 h-5" />
        <span class="text-sm font-medium">Dashboard</span>
      </router-link>

      <!-- Sections -->
      <div v-for="section in navigation" :key="section.title" class="mb-6">
        <h2 class="px-3 mb-2 text-xs font-semibold text-text-muted uppercase tracking-wider">
          {{ section.title }}
        </h2>
        <div class="space-y-1">
          <router-link
            v-for="item in section.items"
            :key="item.route"
            :to="{ name: item.route }"
            class="flex items-center gap-3 px-3 py-2 rounded-lg transition-colors"
            :class="isActive(item.route) ? 'bg-primary text-white' : 'text-text-secondary hover:bg-card hover:text-text-primary'"
          >
            <component :is="item.icon" class="w-5 h-5" />
            <span class="text-sm font-medium">{{ item.name }}</span>
          </router-link>
        </div>
      </div>
    </nav>

    <!-- Sign Out -->
    <div class="p-4 border-t border-border">
      <button
        @click="handleLogout"
        class="flex items-center gap-2 w-full px-3 py-2 text-sm text-text-secondary hover:text-text-primary hover:bg-card rounded-lg transition-colors"
      >
        <LogOut class="w-4 h-4" />
        <span>Sign Out</span>
      </button>
    </div>
  </aside>
</template>





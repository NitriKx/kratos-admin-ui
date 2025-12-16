<script setup lang="ts">
import { ref } from 'vue'
import { Search, HelpCircle, ExternalLink } from 'lucide-vue-next'

const searchQuery = ref('')
</script>

<template>
  <header class="h-16 bg-sidebar border-b border-border flex items-center justify-between px-6">
    <!-- Search -->
    <div class="flex-1 max-w-xl">
      <div class="relative">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-text-muted" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search for users, sessions, or logs..."
          class="w-full pl-10 pr-4 py-2 bg-card border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-primary transition-colors"
        />
      </div>
    </div>

    <!-- Right side -->
    <div class="flex items-center gap-4">
      <!-- Help -->
      <button class="p-2 text-text-secondary hover:text-text-primary transition-colors">
        <HelpCircle class="w-5 h-5" />
      </button>

      <!-- Documentation link -->
      <a
        href="https://www.ory.sh/docs/kratos"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center gap-1 text-sm text-primary hover:text-primary-light transition-colors"
      >
        Documentation
        <ExternalLink class="w-4 h-4" />
      </a>
    </div>
  </header>
</template>





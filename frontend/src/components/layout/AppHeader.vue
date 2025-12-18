<script setup lang="ts">
import { HelpCircle, ExternalLink } from 'lucide-vue-next'
</script>

<template>
  <header class="h-16 bg-sidebar border-b border-border flex items-center justify-end px-6">
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






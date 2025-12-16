<script setup lang="ts">
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
</script>

<template>
  <div class="min-h-screen bg-background">
    <AppSidebar />
    <div class="ml-60">
      <AppHeader />
      <main class="p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>





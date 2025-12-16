<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { FileCode, ChevronDown, ChevronRight, Info } from 'lucide-vue-next'
import { api } from '@/api/client'
import type { IdentitySchema } from '@/types'

const schemas = ref<IdentitySchema[]>([])
const loading = ref(true)
const error = ref('')
const expandedSchemas = ref<Set<string>>(new Set())

const toggleSchema = (id: string) => {
  if (expandedSchemas.value.has(id)) {
    expandedSchemas.value.delete(id)
  } else {
    expandedSchemas.value.add(id)
  }
}

const isExpanded = (id: string) => expandedSchemas.value.has(id)

onMounted(async () => {
  try {
    const response = await api.getSchemas()
    schemas.value = response.data
    // Expand first schema by default
    const firstSchema = schemas.value[0]
    if (firstSchema) {
      expandedSchemas.value.add(firstSchema.id)
    }
  } catch (e) {
    error.value = 'Failed to load schemas'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <!-- Breadcrumb -->
    <div class="mb-2">
      <span class="text-sm text-text-muted">System</span>
      <span class="text-sm text-text-muted mx-2">/</span>
      <span class="text-sm text-text-primary font-medium">Schemas</span>
    </div>

    <!-- Page header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-text-primary">Identity Schemas</h1>
        <p class="text-text-muted mt-1">View configured identity schemas for your Kratos deployment.</p>
      </div>
    </div>

    <!-- Info banner about read-only -->
    <div class="bg-primary/10 border border-primary/20 rounded-xl p-4 mb-6 flex items-start gap-3">
      <Info class="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
      <div>
        <h4 class="font-medium text-text-primary">Schemas are read-only</h4>
        <p class="text-sm text-text-muted mt-1">
          Identity schemas are managed through Kubernetes configuration (ConfigMaps/Secrets) or the Kratos configuration file. 
          Changes to schemas should be made through your infrastructure deployment process.
        </p>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="bg-card border border-border rounded-xl p-8 text-center">
      <div class="flex items-center justify-center gap-2 text-text-muted">
        <div class="w-5 h-5 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
        <span>Loading schemas...</span>
      </div>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="bg-danger/10 border border-danger/20 rounded-xl p-6 text-center">
      <p class="text-danger">{{ error }}</p>
    </div>

    <!-- Empty state -->
    <div v-else-if="schemas.length === 0" class="bg-card border border-border rounded-xl p-8 text-center">
      <FileCode class="w-12 h-12 text-text-muted mx-auto mb-4" />
      <h3 class="text-lg font-medium text-text-primary mb-2">No Schemas Found</h3>
      <p class="text-text-muted">No identity schemas have been configured yet.</p>
    </div>

    <!-- Schemas list -->
    <div v-else class="space-y-4">
      <div
        v-for="schema in schemas"
        :key="schema.id"
        class="bg-card border border-border rounded-xl overflow-hidden"
      >
        <!-- Schema header -->
        <button
          @click="toggleSchema(schema.id)"
          class="w-full flex items-center justify-between p-4 hover:bg-background/50 transition-colors"
        >
          <div class="flex items-center gap-3">
            <div class="p-2 bg-primary/10 rounded-lg">
              <FileCode class="w-5 h-5 text-primary" />
            </div>
            <div class="text-left">
              <h3 class="font-semibold text-text-primary">{{ schema.id }}</h3>
              <p class="text-sm text-text-muted">Identity Schema</p>
            </div>
          </div>
          <component
            :is="isExpanded(schema.id) ? ChevronDown : ChevronRight"
            class="w-5 h-5 text-text-muted"
          />
        </button>

        <!-- Schema content -->
        <div v-if="isExpanded(schema.id)" class="border-t border-border">
          <div class="p-4 bg-background/30">
            <pre class="p-4 bg-background rounded-lg overflow-auto text-sm text-text-primary font-mono max-h-96">{{ JSON.stringify(schema.schema, null, 2) }}</pre>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


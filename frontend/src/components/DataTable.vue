<script setup lang="ts" generic="T extends Record<string, any>">
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

interface Column {
  key: string
  label: string
  width?: string
}

interface Props {
  columns: Column[]
  data: T[]
  loading?: boolean
  page?: number
  perPage?: number
  total?: number
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  page: 1,
  perPage: 20,
  total: 0
})

const emit = defineEmits<{
  (e: 'page-change', page: number): void
}>()

const startItem = () => {
  return (props.page - 1) * props.perPage + 1
}

const endItem = () => {
  return Math.min(props.page * props.perPage, props.total || props.data.length)
}

const totalPages = () => {
  return Math.ceil((props.total || props.data.length) / props.perPage)
}

const canGoPrevious = () => props.page > 1
const canGoNext = () => props.page < totalPages()

const goToPrevious = () => {
  if (canGoPrevious()) {
    emit('page-change', props.page - 1)
  }
}

const goToNext = () => {
  if (canGoNext()) {
    emit('page-change', props.page + 1)
  }
}
</script>

<template>
  <div class="bg-card border border-border rounded-xl overflow-hidden">
    <!-- Table -->
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-border">
            <th
              v-for="column in columns"
              :key="column.key"
              class="px-4 py-3 text-left text-xs font-semibold text-text-muted uppercase tracking-wider"
              :style="column.width ? { width: column.width } : {}"
            >
              {{ column.label }}
            </th>
          </tr>
        </thead>
        <tbody>
          <!-- Loading state -->
          <tr v-if="loading">
            <td :colspan="columns.length" class="px-4 py-8 text-center">
              <div class="flex items-center justify-center gap-2 text-text-muted">
                <div class="w-5 h-5 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
                <span>Loading...</span>
              </div>
            </td>
          </tr>

          <!-- Empty state -->
          <tr v-else-if="data.length === 0">
            <td :colspan="columns.length" class="px-4 py-8 text-center text-text-muted">
              No data available
            </td>
          </tr>

          <!-- Data rows -->
          <tr
            v-else
            v-for="(row, index) in data"
            :key="index"
            class="border-b border-border last:border-b-0 hover:bg-background/50 transition-colors"
          >
            <td
              v-for="column in columns"
              :key="column.key"
              class="px-4 py-4 text-sm text-text-primary"
            >
              <slot :name="`cell-${column.key}`" :row="row" :value="row[column.key]">
                {{ row[column.key] }}
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="px-4 py-3 border-t border-border flex items-center justify-between">
      <p class="text-sm text-text-muted">
        Showing <span class="text-text-primary font-medium">{{ startItem() }}</span>
        to <span class="text-text-primary font-medium">{{ endItem() }}</span>
        of <span class="text-text-primary font-medium">{{ total || data.length }}</span> results
      </p>
      <div class="flex items-center gap-2">
        <button
          @click="goToPrevious"
          :disabled="!canGoPrevious()"
          class="px-3 py-1.5 text-sm border border-border rounded-lg text-text-secondary hover:bg-card disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center gap-1"
        >
          <ChevronLeft class="w-4 h-4" />
          Previous
        </button>
        <button
          @click="goToNext"
          :disabled="!canGoNext()"
          class="px-3 py-1.5 text-sm border border-border rounded-lg text-text-secondary hover:bg-card disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center gap-1"
        >
          Next
          <ChevronRight class="w-4 h-4" />
        </button>
      </div>
    </div>
  </div>
</template>





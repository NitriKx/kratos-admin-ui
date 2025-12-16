<script setup lang="ts">
import { computed } from 'vue'
import { TrendingUp, TrendingDown, Minus } from 'lucide-vue-next'

interface Props {
  title: string
  value: string | number
  icon?: object
  trend?: number
  trendLabel?: string
}

const props = defineProps<Props>()

const trendType = computed(() => {
  if (!props.trend) return 'neutral'
  return props.trend > 0 ? 'positive' : 'negative'
})

const trendIcon = computed(() => {
  if (!props.trend) return Minus
  return props.trend > 0 ? TrendingUp : TrendingDown
})

const trendClass = computed(() => {
  switch (trendType.value) {
    case 'positive': return 'text-success'
    case 'negative': return 'text-danger'
    default: return 'text-text-muted'
  }
})
</script>

<template>
  <div class="bg-card border border-border rounded-xl p-5">
    <div class="flex items-start justify-between">
      <div>
        <p class="text-sm text-text-muted mb-1">{{ title }}</p>
        <p class="text-3xl font-semibold text-text-primary">{{ value }}</p>
        <div v-if="trend !== undefined || trendLabel" class="flex items-center gap-1 mt-2" :class="trendClass">
          <component :is="trendIcon" class="w-4 h-4" />
          <span class="text-sm">
            <template v-if="trend">{{ trend > 0 ? '+' : '' }}{{ trend }}%</template>
            <template v-else-if="trendLabel">{{ trendLabel }}</template>
          </span>
        </div>
      </div>
      <div v-if="icon" class="p-2 bg-background rounded-lg">
        <component :is="icon" class="w-5 h-5 text-text-muted" />
      </div>
    </div>
  </div>
</template>





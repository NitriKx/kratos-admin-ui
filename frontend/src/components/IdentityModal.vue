<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { X, User, Server, ChevronDown, Plus, Trash2 } from 'lucide-vue-next'
import { api } from '@/api/client'
import type { Identity, IdentitySchema } from '@/types'

interface Props {
  isOpen: boolean
  identity: Identity | null
  saving: boolean
}

interface Emits {
  (e: 'close'): void
  (e: 'save', data: { schema_id: string; traits: Record<string, unknown>; state: string }): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Schema management
const schemas = ref<IdentitySchema[]>([])
const loadingSchemas = ref(true)
const selectedSchemaId = ref('')
const traits = ref<Record<string, unknown>>({})
const state = ref<'active' | 'inactive'>('active')

// Computed
const isEditing = computed(() => !!props.identity)
const modalTitle = computed(() => isEditing.value ? 'Edit Identity' : 'Create Identity')
const saveButtonText = computed(() => props.saving ? 'Saving...' : (isEditing.value ? 'Update' : 'Create'))

const selectedSchema = computed(() => {
  return schemas.value.find(s => s.id === selectedSchemaId.value)
})

const schemaProperties = computed(() => {
  if (!selectedSchema.value?.schema) return []
  const schema = selectedSchema.value.schema as Record<string, unknown>
  const properties = schema.properties as Record<string, unknown>
  const traitsSchema = properties?.traits as Record<string, unknown>
  const traitsProperties = traitsSchema?.properties as Record<string, unknown>
  const required = (traitsSchema?.required as string[]) || []
  
  if (!traitsProperties) return []
  
  return Object.entries(traitsProperties).map(([key, value]) => {
    const prop = value as Record<string, unknown>
    return {
      key,
      title: (prop.title as string) || key,
      type: prop.type as string,
      format: prop.format as string | undefined,
      enum: prop.enum as string[] | undefined,
      default: prop.default,
      minimum: prop.minimum as number | undefined,
      maximum: prop.maximum as number | undefined,
      minLength: prop.minLength as number | undefined,
      maxLength: prop.maxLength as number | undefined,
      pattern: prop.pattern as string | undefined,
      items: prop.items as Record<string, unknown> | undefined,
      properties: prop.properties as Record<string, unknown> | undefined,
      required: required.includes(key)
    }
  })
})

// Get schema icon
const getSchemaIcon = (schemaId: string) => {
  if (schemaId.includes('service') || schemaId.includes('account')) {
    return Server
  }
  return User
}

// Load schemas
const loadSchemas = async () => {
  loadingSchemas.value = true
  try {
    const response = await api.getSchemas()
    schemas.value = response.data
    // Set default schema if not editing
    const firstSchema = schemas.value[0]
    if (!props.identity && firstSchema) {
      selectedSchemaId.value = firstSchema.id
    }
  } catch (e) {
    console.error('Failed to load schemas', e)
  } finally {
    loadingSchemas.value = false
  }
}

// Initialize form when modal opens or identity changes
const initForm = () => {
  if (props.identity) {
    selectedSchemaId.value = props.identity.schema_id
    traits.value = JSON.parse(JSON.stringify(props.identity.traits || {}))
    state.value = props.identity.state
  } else {
    traits.value = {}
    state.value = 'active'
    const firstSchema = schemas.value[0]
    if (firstSchema && !selectedSchemaId.value) {
      selectedSchemaId.value = firstSchema.id
    }
  }
}

// Get trait value with support for nested properties
const getTraitValue = (key: string): unknown => {
  return traits.value[key]
}

// Set trait value
const setTraitValue = (key: string, value: unknown) => {
  traits.value[key] = value
}

// Get nested object value
const getNestedValue = (parentKey: string, childKey: string): unknown => {
  const parent = traits.value[parentKey] as Record<string, unknown> | undefined
  return parent?.[childKey]
}

// Set nested object value
const setNestedValue = (parentKey: string, childKey: string, value: unknown) => {
  if (!traits.value[parentKey]) {
    traits.value[parentKey] = {}
  }
  (traits.value[parentKey] as Record<string, unknown>)[childKey] = value
}

// Array handling
const getArrayValue = (key: string): string[] => {
  const value = traits.value[key]
  return Array.isArray(value) ? value : []
}

const addArrayItem = (key: string) => {
  if (!traits.value[key]) {
    traits.value[key] = []
  }
  (traits.value[key] as string[]).push('')
}

const removeArrayItem = (key: string, index: number) => {
  const arr = traits.value[key] as string[]
  arr.splice(index, 1)
}

const updateArrayItem = (key: string, index: number, value: string) => {
  const arr = traits.value[key] as string[]
  arr[index] = value
}

// Handle form submission
const handleSubmit = () => {
  // Clean up empty values
  const cleanedTraits = cleanTraits(traits.value)
  
  emit('save', {
    schema_id: selectedSchemaId.value,
    traits: cleanedTraits,
    state: state.value
  })
}

// Clean empty values from traits
const cleanTraits = (obj: Record<string, unknown>): Record<string, unknown> => {
  const result: Record<string, unknown> = {}
  
  for (const [key, value] of Object.entries(obj)) {
    if (value === undefined || value === null || value === '') {
      continue
    }
    
    if (Array.isArray(value)) {
      const filtered = value.filter(v => v !== '' && v !== null && v !== undefined)
      if (filtered.length > 0) {
        result[key] = filtered
      }
    } else if (typeof value === 'object') {
      const cleaned = cleanTraits(value as Record<string, unknown>)
      if (Object.keys(cleaned).length > 0) {
        result[key] = cleaned
      }
    } else {
      result[key] = value
    }
  }
  
  return result
}

// Close modal
const close = () => {
  emit('close')
}

// Watch for schema changes to reset traits when switching schemas (only for new identities)
watch(selectedSchemaId, (newId, oldId) => {
  if (newId !== oldId && !props.identity) {
    // Initialize default values from schema
    traits.value = {}
    for (const prop of schemaProperties.value) {
      if (prop.default !== undefined) {
        traits.value[prop.key] = prop.default
      }
    }
  }
})

// Watch for modal open
watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    initForm()
  }
})

// Watch for identity changes
watch(() => props.identity, () => {
  if (props.isOpen) {
    initForm()
  }
})

onMounted(() => {
  loadSchemas()
})
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
      >
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-black/50 backdrop-blur-sm"
          @click="close"
        />
        
        <!-- Modal -->
        <div class="relative bg-card border border-border rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] overflow-hidden flex flex-col">
          <!-- Header -->
          <div class="flex items-center justify-between p-4 border-b border-border">
            <h2 class="text-xl font-semibold text-text-primary">{{ modalTitle }}</h2>
            <button
              @click="close"
              class="p-1.5 text-text-muted hover:text-text-primary transition-colors rounded-lg hover:bg-background"
            >
              <X class="w-5 h-5" />
            </button>
          </div>
          
          <!-- Content -->
          <div class="flex-1 overflow-y-auto p-4 space-y-6">
            <!-- Loading state -->
            <div v-if="loadingSchemas" class="flex items-center justify-center py-8">
              <div class="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
              <span class="ml-2 text-text-muted">Loading schemas...</span>
            </div>
            
            <template v-else>
              <!-- Schema Selection (only for new identities) -->
              <div v-if="!isEditing" class="space-y-2">
                <label class="block text-sm font-medium text-text-primary">
                  Identity Schema <span class="text-danger">*</span>
                </label>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <button
                    v-for="schema in schemas"
                    :key="schema.id"
                    type="button"
                    @click="selectedSchemaId = schema.id"
                    :class="[
                      'flex items-center gap-3 p-4 rounded-lg border-2 transition-all text-left',
                      selectedSchemaId === schema.id
                        ? 'border-primary bg-primary/5'
                        : 'border-border hover:border-primary/50 hover:bg-background/50'
                    ]"
                  >
                    <div :class="[
                      'p-2 rounded-lg',
                      selectedSchemaId === schema.id ? 'bg-primary/20' : 'bg-background'
                    ]">
                      <component
                        :is="getSchemaIcon(schema.id)"
                        :class="[
                          'w-5 h-5',
                          selectedSchemaId === schema.id ? 'text-primary' : 'text-text-muted'
                        ]"
                      />
                    </div>
                    <div>
                      <p :class="[
                        'font-medium',
                        selectedSchemaId === schema.id ? 'text-primary' : 'text-text-primary'
                      ]">
                        {{ (schema.schema as Record<string, unknown>).title || schema.id }}
                      </p>
                      <p class="text-xs text-text-muted">{{ schema.id }}</p>
                    </div>
                  </button>
                </div>
              </div>
              
              <!-- Schema info for editing -->
              <div v-else class="bg-background/50 rounded-lg p-4 flex items-center gap-3">
                <div class="p-2 bg-primary/10 rounded-lg">
                  <component :is="getSchemaIcon(selectedSchemaId)" class="w-5 h-5 text-primary" />
                </div>
                <div>
                  <p class="text-sm text-text-muted">Schema</p>
                  <p class="font-medium text-text-primary">{{ selectedSchemaId }}</p>
                </div>
              </div>
              
              <!-- State Selection -->
              <div class="space-y-2">
                <label class="block text-sm font-medium text-text-primary">State</label>
                <select
                  v-model="state"
                  class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors"
                >
                  <option value="active">Active</option>
                  <option value="inactive">Inactive</option>
                </select>
              </div>
              
              <!-- Dynamic Traits Form -->
              <div v-if="selectedSchema" class="space-y-4">
                <h3 class="text-sm font-semibold text-text-primary uppercase tracking-wide">Traits</h3>
                
                <div
                  v-for="prop in schemaProperties"
                  :key="prop.key"
                  class="space-y-2"
                >
                  <!-- String fields -->
                  <template v-if="prop.type === 'string' && !prop.enum && !prop.properties">
                    <label class="block text-sm font-medium text-text-primary">
                      {{ prop.title }}
                      <span v-if="prop.required" class="text-danger">*</span>
                    </label>
                    
                    <!-- Date input -->
                    <input
                      v-if="prop.format === 'date'"
                      type="date"
                      :value="getTraitValue(prop.key) as string || ''"
                      @input="setTraitValue(prop.key, ($event.target as HTMLInputElement).value)"
                      class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors"
                    />
                    
                    <!-- Email input -->
                    <input
                      v-else-if="prop.format === 'email'"
                      type="email"
                      :value="getTraitValue(prop.key) as string || ''"
                      @input="setTraitValue(prop.key, ($event.target as HTMLInputElement).value)"
                      :placeholder="`Enter ${prop.title.toLowerCase()}`"
                      class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors"
                    />
                    
                    <!-- URL input -->
                    <input
                      v-else-if="prop.format === 'uri'"
                      type="url"
                      :value="getTraitValue(prop.key) as string || ''"
                      @input="setTraitValue(prop.key, ($event.target as HTMLInputElement).value)"
                      :placeholder="`Enter ${prop.title.toLowerCase()}`"
                      class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors"
                    />
                    
                    <!-- Tel input -->
                    <input
                      v-else-if="prop.format === 'tel'"
                      type="tel"
                      :value="getTraitValue(prop.key) as string || ''"
                      @input="setTraitValue(prop.key, ($event.target as HTMLInputElement).value)"
                      :placeholder="`Enter ${prop.title.toLowerCase()}`"
                      class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors"
                    />
                    
                    <!-- Textarea for long text -->
                    <textarea
                      v-else-if="prop.maxLength && prop.maxLength > 100"
                      :value="getTraitValue(prop.key) as string || ''"
                      @input="setTraitValue(prop.key, ($event.target as HTMLTextAreaElement).value)"
                      :placeholder="`Enter ${prop.title.toLowerCase()}`"
                      :maxlength="prop.maxLength"
                      rows="3"
                      class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors resize-none"
                    />
                    
                    <!-- Default text input -->
                    <input
                      v-else
                      type="text"
                      :value="getTraitValue(prop.key) as string || ''"
                      @input="setTraitValue(prop.key, ($event.target as HTMLInputElement).value)"
                      :placeholder="`Enter ${prop.title.toLowerCase()}`"
                      :pattern="prop.pattern"
                      :minlength="prop.minLength"
                      :maxlength="prop.maxLength"
                      class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors"
                    />
                    
                    <p v-if="prop.pattern" class="text-xs text-text-muted">
                      Pattern: {{ prop.pattern }}
                    </p>
                  </template>
                  
                  <!-- Enum/Select fields -->
                  <template v-else-if="prop.type === 'string' && prop.enum">
                    <label class="block text-sm font-medium text-text-primary">
                      {{ prop.title }}
                      <span v-if="prop.required" class="text-danger">*</span>
                    </label>
                    <div class="relative">
                      <select
                        :value="getTraitValue(prop.key) as string || ''"
                        @change="setTraitValue(prop.key, ($event.target as HTMLSelectElement).value)"
                        class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors appearance-none pr-10"
                      >
                        <option value="">Select {{ prop.title.toLowerCase() }}</option>
                        <option v-for="option in prop.enum" :key="option" :value="option">
                          {{ option }}
                        </option>
                      </select>
                      <ChevronDown class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted pointer-events-none" />
                    </div>
                  </template>
                  
                  <!-- Integer fields -->
                  <template v-else-if="prop.type === 'integer'">
                    <label class="block text-sm font-medium text-text-primary">
                      {{ prop.title }}
                      <span v-if="prop.required" class="text-danger">*</span>
                    </label>
                    <input
                      type="number"
                      :value="getTraitValue(prop.key) as number || ''"
                      @input="setTraitValue(prop.key, parseInt(($event.target as HTMLInputElement).value) || undefined)"
                      :min="prop.minimum"
                      :max="prop.maximum"
                      :placeholder="`Enter ${prop.title.toLowerCase()}`"
                      class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors"
                    />
                  </template>
                  
                  <!-- Number fields -->
                  <template v-else-if="prop.type === 'number'">
                    <label class="block text-sm font-medium text-text-primary">
                      {{ prop.title }}
                      <span v-if="prop.required" class="text-danger">*</span>
                    </label>
                    <input
                      type="number"
                      step="any"
                      :value="getTraitValue(prop.key) as number || ''"
                      @input="setTraitValue(prop.key, parseFloat(($event.target as HTMLInputElement).value) || undefined)"
                      :min="prop.minimum"
                      :max="prop.maximum"
                      :placeholder="`Enter ${prop.title.toLowerCase()}`"
                      class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors"
                    />
                  </template>
                  
                  <!-- Boolean fields -->
                  <template v-else-if="prop.type === 'boolean'">
                    <label class="flex items-center gap-3 cursor-pointer">
                      <input
                        type="checkbox"
                        :checked="getTraitValue(prop.key) as boolean || false"
                        @change="setTraitValue(prop.key, ($event.target as HTMLInputElement).checked)"
                        class="w-4 h-4 rounded border-border text-primary focus:ring-primary focus:ring-offset-0 bg-background"
                      />
                      <span class="text-sm font-medium text-text-primary">
                        {{ prop.title }}
                        <span v-if="prop.required" class="text-danger">*</span>
                      </span>
                    </label>
                  </template>
                  
                  <!-- Array fields -->
                  <template v-else-if="prop.type === 'array'">
                    <label class="block text-sm font-medium text-text-primary">
                      {{ prop.title }}
                      <span v-if="prop.required" class="text-danger">*</span>
                    </label>
                    
                    <!-- Array with enum items -->
                    <div v-if="prop.items?.enum" class="space-y-2">
                      <div class="flex flex-wrap gap-2">
                        <label
                          v-for="option in (prop.items.enum as string[])"
                          :key="option"
                          class="flex items-center gap-2 px-3 py-1.5 bg-background border border-border rounded-lg cursor-pointer hover:border-primary/50 transition-colors"
                        >
                          <input
                            type="checkbox"
                            :checked="getArrayValue(prop.key).includes(option)"
                            @change="(e) => {
                              const arr = getArrayValue(prop.key)
                              if ((e.target as HTMLInputElement).checked) {
                                setTraitValue(prop.key, [...arr, option])
                              } else {
                                setTraitValue(prop.key, arr.filter(v => v !== option))
                              }
                            }"
                            class="w-4 h-4 rounded border-border text-primary focus:ring-primary focus:ring-offset-0 bg-background"
                          />
                          <span class="text-sm text-text-primary">{{ option }}</span>
                        </label>
                      </div>
                    </div>
                    
                    <!-- Array with string items -->
                    <div v-else class="space-y-2">
                      <div
                        v-for="(item, index) in getArrayValue(prop.key)"
                        :key="index"
                        class="flex items-center gap-2"
                      >
                        <input
                          type="text"
                          :value="item"
                          @input="updateArrayItem(prop.key, index, ($event.target as HTMLInputElement).value)"
                          :placeholder="`Enter ${prop.title.toLowerCase()}`"
                          class="flex-1 px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors"
                        />
                        <button
                          type="button"
                          @click="removeArrayItem(prop.key, index)"
                          class="p-2 text-text-muted hover:text-danger transition-colors"
                        >
                          <Trash2 class="w-4 h-4" />
                        </button>
                      </div>
                      <button
                        type="button"
                        @click="addArrayItem(prop.key)"
                        class="flex items-center gap-1 text-sm text-primary hover:text-primary-hover transition-colors"
                      >
                        <Plus class="w-4 h-4" />
                        Add {{ prop.title.toLowerCase() }}
                      </button>
                    </div>
                  </template>
                  
                  <!-- Nested Object fields -->
                  <template v-else-if="prop.type === 'object' && prop.properties">
                    <div class="border border-border rounded-lg p-4 space-y-4 bg-background/30">
                      <label class="block text-sm font-semibold text-text-primary">
                        {{ prop.title }}
                        <span v-if="prop.required" class="text-danger">*</span>
                      </label>
                      
                      <div
                        v-for="[childKey, childValue] in Object.entries(prop.properties)"
                        :key="childKey"
                        class="space-y-1"
                      >
                        <label class="block text-sm font-medium text-text-muted">
                          {{ (childValue as Record<string, unknown>).title || childKey }}
                        </label>
                        
                        <!-- Nested enum select -->
                        <template v-if="(childValue as Record<string, unknown>).enum">
                          <div class="relative">
                            <select
                              :value="getNestedValue(prop.key, childKey) as string || ''"
                              @change="setNestedValue(prop.key, childKey, ($event.target as HTMLSelectElement).value)"
                              class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors appearance-none pr-10"
                            >
                              <option value="">Select {{ ((childValue as Record<string, unknown>).title as string || childKey).toLowerCase() }}</option>
                              <option
                                v-for="option in ((childValue as Record<string, unknown>).enum as string[])"
                                :key="option"
                                :value="option"
                              >
                                {{ option }}
                              </option>
                            </select>
                            <ChevronDown class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted pointer-events-none" />
                          </div>
                        </template>
                        
                        <!-- Nested boolean -->
                        <template v-else-if="(childValue as Record<string, unknown>).type === 'boolean'">
                          <label class="flex items-center gap-2 cursor-pointer">
                            <input
                              type="checkbox"
                              :checked="getNestedValue(prop.key, childKey) as boolean || false"
                              @change="setNestedValue(prop.key, childKey, ($event.target as HTMLInputElement).checked)"
                              class="w-4 h-4 rounded border-border text-primary focus:ring-primary focus:ring-offset-0 bg-background"
                            />
                            <span class="text-sm text-text-muted">Enable</span>
                          </label>
                        </template>
                        
                        <!-- Nested text input -->
                        <input
                          v-else
                          type="text"
                          :value="getNestedValue(prop.key, childKey) as string || ''"
                          @input="setNestedValue(prop.key, childKey, ($event.target as HTMLInputElement).value)"
                          :placeholder="`Enter ${((childValue as Record<string, unknown>).title as string || childKey).toLowerCase()}`"
                          class="w-full px-3 py-2 bg-background border border-border rounded-lg text-text-primary focus:outline-none focus:border-primary transition-colors"
                        />
                      </div>
                    </div>
                  </template>
                </div>
              </div>
            </template>
          </div>
          
          <!-- Footer -->
          <div class="flex items-center justify-end gap-3 p-4 border-t border-border bg-background/50">
            <button
              type="button"
              @click="close"
              class="px-4 py-2 text-text-primary hover:bg-background rounded-lg transition-colors"
            >
              Cancel
            </button>
            <button
              type="button"
              @click="handleSubmit"
              :disabled="saving || !selectedSchemaId"
              class="px-4 py-2 bg-primary hover:bg-primary-hover disabled:bg-primary/50 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors"
            >
              {{ saveButtonText }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .relative,
.modal-leave-to .relative {
  transform: scale(0.95);
}
</style>

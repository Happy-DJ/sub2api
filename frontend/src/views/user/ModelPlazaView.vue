<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Header -->
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">
          {{ t('modelPlaza.title') }}
        </h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ t('modelPlaza.description') }}
        </p>
      </div>

      <!-- Platform filter tabs -->
      <div class="flex items-center gap-2 overflow-x-auto pb-1 scrollbar-thin">
        <button
          class="px-3.5 py-1.5 rounded-full text-sm font-medium transition-all duration-200 whitespace-nowrap"
          :class="selectedPlatform === ''
            ? 'bg-primary-600 text-white shadow-sm'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-dark-800 dark:text-gray-300 dark:hover:bg-dark-700'"
          @click="selectedPlatform = ''"
        >
          {{ t('modelPlaza.allPlatforms') }}
        </button>
        <button
          v-for="platform in availablePlatforms"
          :key="platform.platform"
          class="px-3.5 py-1.5 rounded-full text-sm font-medium transition-all duration-200 whitespace-nowrap"
          :class="selectedPlatform === platform.platform
            ? 'bg-primary-600 text-white shadow-sm'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-dark-800 dark:text-gray-300 dark:hover:bg-dark-700'"
          @click="selectedPlatform = platform.platform"
        >
          {{ platform.display_name }}
        </button>
      </div>

      <!-- Search -->
      <div class="relative max-w-sm">
        <svg
          class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          :placeholder="t('modelPlaza.searchPlaceholder')"
          class="w-full pl-10 pr-4 py-2 rounded-xl border border-gray-200 bg-white/70 backdrop-blur-xl text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500 dark:bg-dark-800/60 dark:border-dark-700/70 dark:text-gray-100 dark:placeholder-gray-500 transition-all"
        />
      </div>

      <!-- Loading skeleton -->
      <div v-if="loading" class="grid gap-5 grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
        <div
          v-for="n in 8"
          :key="n"
          class="p-5 rounded-2xl min-h-[220px] bg-white/70 backdrop-blur-xl border border-gray-200/80 dark:bg-dark-800/60 dark:border-dark-700/70 animate-pulse"
        >
          <div class="flex items-start gap-3">
            <div class="w-9 h-9 rounded-xl bg-gray-200 dark:bg-dark-700" />
            <div class="flex-1 space-y-2">
              <div class="h-4 w-24 bg-gray-200 dark:bg-dark-700 rounded" />
              <div class="h-3 w-16 bg-gray-100 dark:bg-dark-700 rounded" />
            </div>
          </div>
          <div class="mt-4 space-y-3">
            <div class="h-3 w-full bg-gray-100 dark:bg-dark-700 rounded" />
            <div class="h-3 w-3/4 bg-gray-100 dark:bg-dark-700 rounded" />
          </div>
        </div>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="text-center py-20">
        <p class="text-gray-500 dark:text-gray-400">{{ error }}</p>
        <button
          class="mt-4 px-4 py-2 rounded-lg bg-primary-600 text-white text-sm font-medium hover:bg-primary-700 transition-colors"
          @click="fetchModels"
        >
          {{ t('common.retry') }}
        </button>
      </div>

      <!-- Empty state -->
      <div v-else-if="filteredGroups.length === 0" class="text-center py-20">
        <p class="text-gray-400 dark:text-gray-500">{{ t('modelPlaza.empty') }}</p>
      </div>

      <!-- Card grid -->
      <div v-else class="grid gap-5 grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
        <template v-for="group in filteredGroups" :key="group.platform">
          <ModelPricingCard
            v-for="model in group.models"
            :key="model.id"
            :model="model"
            :platform="group.platform"
          />
        </template>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import ModelPricingCard from '@/components/models/ModelPricingCard.vue'
import { modelsAPI } from '@/api/models'
import type { PlatformModelGroup } from '@/types/models'

const { t } = useI18n()

const platformGroups = ref<PlatformModelGroup[]>([])
const loading = ref(true)
const error = ref('')
const selectedPlatform = ref('')
const searchQuery = ref('')

// Available platforms derived from data (sorted by sort_order)
const availablePlatforms = computed(() =>
  [...platformGroups.value].sort((a, b) => a.sort_order - b.sort_order)
)

// Filtered groups by platform, then by search
const filteredGroups = computed(() => {
  let groups = platformGroups.value

  // Filter by selected platform
  if (selectedPlatform.value) {
    groups = groups.filter(g => g.platform === selectedPlatform.value)
  }

  // Filter by search query
  const query = searchQuery.value.trim().toLowerCase()
  if (query) {
    groups = groups
      .map(g => ({
        ...g,
        models: g.models.filter(
          m =>
            m.id.toLowerCase().includes(query) ||
            m.display_name.toLowerCase().includes(query)
        ),
      }))
      .filter(g => g.models.length > 0)
  }

  return groups
})

async function fetchModels() {
  loading.value = true
  error.value = ''
  try {
    platformGroups.value = await modelsAPI.getModelPlaza()
  } catch (e: unknown) {
    const msg = e instanceof Error ? e.message : 'Failed to load model plaza'
    error.value = msg
    console.error('[ModelPlaza] fetch failed:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchModels()
})
</script>

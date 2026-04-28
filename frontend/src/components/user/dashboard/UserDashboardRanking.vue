<template>
  <div class="card">
    <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('dashboard.tokenRankingTitle') }}</h2>
      <div class="flex items-center gap-2">
        <!-- Tabs -->
        <div class="inline-flex rounded-lg bg-gray-100 p-1 dark:bg-dark-800">
          <button
            type="button"
            class="rounded-md px-3 py-1 text-xs font-medium transition-colors"
            :class="activeTab === 'today'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            @click="switchTab('today')"
          >
            {{ t('dashboard.tokenRankingToday') }}
          </button>
          <button
            type="button"
            class="rounded-md px-3 py-1 text-xs font-medium transition-colors"
            :class="activeTab === 'monthly'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            @click="switchTab('monthly')"
          >
            {{ t('dashboard.tokenRankingMonthly') }}
          </button>
          <button
            type="button"
            class="rounded-md px-3 py-1 text-xs font-medium transition-colors"
            :class="activeTab === 'history'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            @click="switchTab('history')"
          >
            {{ t('dashboard.tokenRankingHistory') }}
          </button>
        </div>

        <!-- Month Picker -->
        <input
          v-if="activeTab === 'monthly'"
          v-model="monthInput"
          type="month"
          :max="maxMonth"
          class="rounded-md border border-gray-200 bg-white px-2 py-1 text-xs dark:border-gray-700 dark:bg-dark-800 dark:text-white"
          @change="onMonthChange"
        />
      </div>
    </div>

    <div class="p-6">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <LoadingSpinner size="lg" />
      </div>

      <!-- Empty -->
      <div v-else-if="items.length === 0" class="py-8">
        <EmptyState :title="t('dashboard.tokenRankingNoData')" />
      </div>

      <!-- Data -->
      <div v-else class="space-y-2">
        <div
          v-for="(item, index) in items"
          :key="item.user_id"
          class="flex items-center justify-between rounded-xl bg-gray-50 p-4 transition-colors hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800"
        >
          <div class="flex items-center gap-4">
            <div class="flex h-8 w-8 items-center justify-center rounded-full text-sm font-bold"
              :class="index < 3 ? 'bg-amber-100 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400' : 'bg-gray-200 text-gray-500 dark:bg-dark-700 dark:text-gray-400'">
              {{ index + 1 }}
            </div>
            <div>
              <p class="text-sm font-medium text-gray-900 dark:text-white">{{ item.username || item.email }}</p>
              <p class="text-xs text-gray-500 dark:text-dark-400">{{ item.email }}</p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-sm font-semibold text-gray-900 dark:text-white">{{ formatTokens(item.tokens) }} tokens</p>
            <p class="text-xs">
              <span class="text-green-600 dark:text-green-400">${{ formatCost(item.actual_cost) }}</span>
              <span class="text-gray-400 dark:text-gray-500"> · </span>
              <span class="text-gray-500 dark:text-gray-400">{{ formatNumber(item.requests) }} {{ t('dashboard.tokenRankingRequests') }}</span>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import { usageAPI } from '@/api/usage'

const { t } = useI18n()

type Tab = 'today' | 'monthly' | 'history'
interface RankingItem {
  user_id: number
  email: string
  username: string
  actual_cost: number
  requests: number
  tokens: number
}

const activeTab = ref<Tab>('today')
const items = ref<RankingItem[]>([])
const loading = ref(false)
const monthInput = ref('')
const loadSeq = ref(0)

const now = new Date()
const maxMonth = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
monthInput.value = maxMonth

const switchTab = (tab: Tab) => {
  activeTab.value = tab
  if (tab === 'today') loadToday()
  else if (tab === 'monthly') loadMonthly()
  else loadHistory()
}

const onMonthChange = () => {
  loadMonthly()
}

const loadToday = async () => {
  const seq = ++loadSeq.value
  loading.value = true
  try {
    const res = await usageAPI.getRankingToday()
    if (seq !== loadSeq.value) return
    items.value = res.ranking || []
  } catch {
    if (seq !== loadSeq.value) return
    items.value = []
  } finally {
    if (seq === loadSeq.value) loading.value = false
  }
}

const loadMonthly = async () => {
  const seq = ++loadSeq.value
  loading.value = true
  try {
    const [year, month] = monthInput.value.split('-').map(Number)
    const res = await usageAPI.getRankingMonthly({ year, month })
    if (seq !== loadSeq.value) return
    items.value = res.ranking || []
  } catch {
    if (seq !== loadSeq.value) return
    items.value = []
  } finally {
    if (seq === loadSeq.value) loading.value = false
  }
}

const loadHistory = async () => {
  const seq = ++loadSeq.value
  loading.value = true
  try {
    const res = await usageAPI.getRankingHistory()
    if (seq !== loadSeq.value) return
    items.value = res.ranking || []
  } catch {
    if (seq !== loadSeq.value) return
    items.value = []
  } finally {
    if (seq === loadSeq.value) loading.value = false
  }
}

// Initial load
loadToday()

const formatTokens = (value: number): string => {
  if (value >= 1_000_000_000) return `${(value / 1_000_000_000).toFixed(2)}B`
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(2)}M`
  if (value >= 1_000) return `${(value / 1_000).toFixed(2)}K`
  return value.toLocaleString()
}

const formatCost = (value: number): string => {
  if (value >= 1000) return (value / 1000).toFixed(2) + 'K'
  if (value >= 1) return value.toFixed(2)
  if (value >= 0.01) return value.toFixed(3)
  return value.toFixed(4)
}

const formatNumber = (value: number): string => {
  return value.toLocaleString()
}
</script>

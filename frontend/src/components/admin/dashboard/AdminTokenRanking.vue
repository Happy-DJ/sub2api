<template>
  <div class="card p-4">
    <div class="mb-4 flex items-center justify-between gap-3">
      <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
        {{ t('admin.dashboard.tokenRankingTitle') }}
      </h3>
      <div class="flex flex-wrap items-center justify-end gap-2">
        <!-- Tabs -->
        <div class="inline-flex rounded-lg bg-gray-100 p-1 dark:bg-dark-800">
          <button
            type="button"
            class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
            :class="activeTab === 'today'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            @click="switchTab('today')"
          >
            {{ t('admin.dashboard.tokenRankingToday') }}
          </button>
          <button
            type="button"
            class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
            :class="activeTab === 'monthly'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            @click="switchTab('monthly')"
          >
            {{ t('admin.dashboard.tokenRankingMonthly') }}
          </button>
          <button
            type="button"
            class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
            :class="activeTab === 'history'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            @click="switchTab('history')"
          >
            {{ t('admin.dashboard.tokenRankingHistory') }}
          </button>
        </div>

        <!-- Month Picker (only for monthly tab) -->
        <div v-if="activeTab === 'monthly'" class="flex items-center gap-1">
          <input
            v-model="monthInput"
            type="month"
            :max="maxMonth"
            class="rounded-md border border-gray-200 bg-white px-2 py-1 text-xs dark:border-gray-700 dark:bg-dark-800 dark:text-white"
            @change="onMonthChange"
          />
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex h-48 items-center justify-center">
      <LoadingSpinner />
    </div>

    <!-- Error -->
    <div v-else-if="error" class="flex h-48 items-center justify-center text-sm text-red-500">
      {{ t('admin.dashboard.failedToLoad') }}
    </div>

    <!-- Data Table -->
    <div v-else-if="items.length > 0">
      <div class="overflow-x-auto">
        <table class="w-full text-xs">
          <thead>
            <tr class="border-b border-gray-100 text-gray-500 dark:border-gray-700">
              <th class="pb-2 pr-4 text-left font-medium">{{ t('admin.dashboard.tokenRankingNo') }}</th>
              <th class="pb-2 pr-4 text-left font-medium">{{ t('admin.dashboard.tokenRankingUsername') }}</th>
              <th class="pb-2 pr-4 text-left font-medium">{{ t('admin.dashboard.tokenRankingNickname') }}</th>
              <th class="pb-2 pr-4 text-right font-medium">{{ t('admin.dashboard.tokenRankingTokens') }}</th>
              <th class="pb-2 pr-4 text-right font-medium">{{ t('admin.dashboard.tokenRankingCost') }}</th>
              <th class="pb-2 text-right font-medium">{{ t('admin.dashboard.tokenRankingRequests') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(item, index) in items"
              :key="item.user_id"
              class="cursor-pointer border-b border-gray-50 transition-colors hover:bg-gray-50 dark:border-gray-800 dark:hover:bg-dark-700/40"
              @click="emit('row-click', item)"
            >
              <td class="py-2 pr-4">
                <span class="inline-flex h-5 w-5 items-center justify-center rounded bg-primary-100 text-[10px] font-semibold text-primary-600 dark:bg-primary-900/30 dark:text-primary-400">
                  {{ index + 1 }}
                </span>
              </td>
              <td class="py-2 pr-4">
                <span class="font-medium text-gray-900 dark:text-white">{{ item.email }}</span>
              </td>
              <td class="py-2 pr-4">
                <span class="text-gray-600 dark:text-gray-400">{{ item.username || '-' }}</span>
              </td>
              <td class="py-2 pr-4 text-right font-medium text-gray-900 dark:text-white">
                {{ formatTokens(item.tokens) }}
              </td>
              <td class="py-2 pr-4 text-right text-green-600 dark:text-green-400">
                ${{ formatCost(item.actual_cost) }}
              </td>
              <td class="py-2 text-right text-gray-600 dark:text-gray-400">
                {{ formatNumber(item.requests) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Empty -->
    <div v-else class="flex h-48 items-center justify-center text-sm text-gray-400 dark:text-gray-500">
      {{ t('admin.dashboard.tokenRankingNoData') }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { getUserSpendingRanking } from '@/api/admin/dashboard'
import type { UserSpendingRankingItem } from '@/types'

const { t } = useI18n()

// Format date in local timezone (not UTC) to avoid off-by-one-day errors
const fmtLocalDate = (d: Date): string => {
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

type Tab = 'today' | 'monthly' | 'history'

const emit = defineEmits<{
  'row-click': [item: UserSpendingRankingItem]
}>()

const activeTab = ref<Tab>('today')
const items = ref<UserSpendingRankingItem[]>([])
const loading = ref(false)
const error = ref(false)
const monthInput = ref('')
const loadSeq = ref(0)

// Max month is current month (disable future)
const now = new Date()
const maxMonth = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
monthInput.value = maxMonth

const switchTab = (tab: Tab) => {
  activeTab.value = tab
  if (tab === 'monthly') {
    loadMonthlyRanking()
  } else if (tab === 'today') {
    loadTodayRanking()
  } else {
    loadHistoryRanking()
  }
}

const onMonthChange = () => {
  loadMonthlyRanking()
}

const loadTodayRanking = async () => {
  const currentSeq = ++loadSeq.value
  loading.value = true
  error.value = false
  try {
    const now = new Date()
    const startOfDay = new Date(now.getFullYear(), now.getMonth(), now.getDate())
    const endOfDay = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1)
    const response = await getUserSpendingRanking({
      start_date: fmtLocalDate(startOfDay),
      end_date: fmtLocalDate(endOfDay),
      sort_by: 'tokens',
      limit: 20
    })
    if (currentSeq !== loadSeq.value) return
    items.value = response.ranking || []
  } catch (e) {
    if (currentSeq !== loadSeq.value) return
    console.error('Error loading today ranking:', e)
    items.value = []
    error.value = true
  } finally {
    if (currentSeq === loadSeq.value) {
      loading.value = false
    }
  }
}

const loadMonthlyRanking = async () => {
  const currentSeq = ++loadSeq.value
  loading.value = true
  error.value = false
  try {
    const [year, month] = monthInput.value.split('-').map(Number)
    const startOfMonth = new Date(year, month - 1, 1)
    const endOfMonth = new Date(year, month, 1)
    const response = await getUserSpendingRanking({
      start_date: fmtLocalDate(startOfMonth),
      end_date: fmtLocalDate(endOfMonth),
      sort_by: 'tokens',
      limit: 20
    })
    if (currentSeq !== loadSeq.value) return
    items.value = response.ranking || []
  } catch (e) {
    if (currentSeq !== loadSeq.value) return
    console.error('Error loading monthly ranking:', e)
    items.value = []
    error.value = true
  } finally {
    if (currentSeq === loadSeq.value) {
      loading.value = false
    }
  }
}

const loadHistoryRanking = async () => {
  const currentSeq = ++loadSeq.value
  loading.value = true
  error.value = false
  try {
    // Use a very old start date for history
    const response = await getUserSpendingRanking({
      start_date: '2020-01-01',
      end_date: fmtLocalDate(new Date()),
      sort_by: 'tokens',
      limit: 20
    })
    if (currentSeq !== loadSeq.value) return
    items.value = response.ranking || []
  } catch (e) {
    if (currentSeq !== loadSeq.value) return
    console.error('Error loading history ranking:', e)
    items.value = []
    error.value = true
  } finally {
    if (currentSeq === loadSeq.value) {
      loading.value = false
    }
  }
}

// Initial load
loadTodayRanking()

const formatTokens = (value: number): string => {
  if (value >= 1_000_000_000) {
    return `${(value / 1_000_000_000).toFixed(2)}B`
  } else if (value >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(2)}M`
  } else if (value >= 1_000) {
    return `${(value / 1_000).toFixed(2)}K`
  }
  return value.toLocaleString()
}

const formatCost = (value: number): string => {
  if (value >= 1000) {
    return (value / 1000).toFixed(2) + 'K'
  } else if (value >= 1) {
    return value.toFixed(2)
  } else if (value >= 0.01) {
    return value.toFixed(3)
  }
  return value.toFixed(4)
}

const formatNumber = (value: number): string => {
  return value.toLocaleString()
}
</script>

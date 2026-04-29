<template>
  <div
    class="group p-5 rounded-2xl min-h-[220px] w-full bg-white/70 backdrop-blur-xl border border-gray-200/80 shadow-card dark:bg-dark-800/60 dark:border-dark-700/70 hover:-translate-y-1 hover:shadow-card-hover dark:hover:border-primary-500/30 hover:border-gray-300 transition-all duration-300 ease-out flex flex-col"
  >
    <!-- Header: platform icon + model name -->
    <div class="flex items-start gap-3">
      <span
        class="w-9 h-9 rounded-xl ring-1 ring-black/5 dark:ring-white/10 grid place-items-center flex-shrink-0"
        :class="[platformGradient, platformTintClass]"
      >
        <component :is="platformIcon" :size="20" />
      </span>
      <div class="flex-1 min-w-0">
        <div class="text-base font-semibold truncate text-gray-900 dark:text-gray-100">
          {{ model.display_name }}
        </div>
        <div class="mt-0.5 flex items-center gap-1.5 min-w-0">
          <span
            class="inline-flex items-center rounded-md px-1.5 py-0.5 text-[10px] font-medium flex-shrink-0"
            :class="platformBadgeClass"
          >
            {{ platformLabel }}
          </span>
          <span class="font-mono text-xs truncate text-gray-500 dark:text-gray-400">
            {{ model.id }}
          </span>
          <!-- Copy model name button -->
          <button
            class="flex-shrink-0 w-5 h-5 rounded flex items-center justify-center text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700 transition-colors"
            :title="t('modelPlaza.copyModelName')"
            @click="copyModelId"
          >
            <!-- Clipboard icon -->
            <svg v-if="!copied" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <!-- Check icon (copied feedback) -->
            <svg v-else class="w-3 h-3 text-emerald-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>
          <span
            v-if="model.mode && model.mode !== 'chat'"
            class="inline-flex items-center rounded-md px-1.5 py-0.5 text-[10px] font-medium bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-gray-300 flex-shrink-0"
          >
            {{ model.mode }}
          </span>
        </div>
      </div>
    </div>

    <!-- Pricing section -->
    <div class="mt-4 space-y-2 flex-1">
      <!-- Input price -->
      <div class="flex items-center justify-between text-sm">
        <span class="text-gray-500 dark:text-gray-400">{{ t('modelPlaza.inputPrice') }}</span>
        <span class="font-semibold text-gray-900 dark:text-gray-100">
          {{ model.input_price_per_mtok != null ? '$' + model.input_price_per_mtok.toFixed(2) : '—' }}
          <span v-if="model.input_price_per_mtok != null" class="text-xs font-normal text-gray-400">{{
            t('modelPlaza.perMillionTokens')
          }}</span>
        </span>
      </div>

      <!-- Output price -->
      <div class="flex items-center justify-between text-sm">
        <span class="text-gray-500 dark:text-gray-400">{{ t('modelPlaza.outputPrice') }}</span>
        <span class="font-semibold text-gray-900 dark:text-gray-100">
          {{ model.output_price_per_mtok != null ? '$' + model.output_price_per_mtok.toFixed(2) : '—' }}
          <span v-if="model.output_price_per_mtok != null" class="text-xs font-normal text-gray-400">{{
            t('modelPlaza.perMillionTokens')
          }}</span>
        </span>
      </div>

      <!-- Cache read price (only if available) -->
      <div v-if="model.cache_read_price_per_mtok != null" class="flex items-center justify-between text-sm">
        <span class="text-gray-500 dark:text-gray-400">{{ t('modelPlaza.cacheReadPrice') }}</span>
        <span class="font-semibold text-gray-900 dark:text-gray-100">
          ${{ model.cache_read_price_per_mtok.toFixed(2) }}
          <span class="text-xs font-normal text-gray-400">{{ t('modelPlaza.perMillionTokens') }}</span>
        </span>
      </div>
    </div>

    <!-- Feature badges -->
    <div class="mt-3 flex items-center gap-1.5 flex-wrap">
      <span
        v-if="model.supports_prompt_caching"
        class="inline-flex items-center rounded-md px-1.5 py-0.5 text-[10px] font-medium bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400"
      >
        {{ t('modelPlaza.supportsCaching') }}
      </span>
      <span
        v-if="model.supports_vision"
        class="inline-flex items-center rounded-md px-1.5 py-0.5 text-[10px] font-medium bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400"
      >
        {{ t('modelPlaza.supportsVision') }}
      </span>
      <span
        v-if="model.input_price_per_mtok == null && model.output_price_per_mtok == null"
        class="inline-flex items-center rounded-md px-1.5 py-0.5 text-[10px] font-medium bg-gray-100 text-gray-500 dark:bg-dark-700 dark:text-gray-400"
      >
        {{ t('modelPlaza.noPricing') }}
      </span>
    </div>

    <!-- API endpoint -->
    <div
      v-if="model.endpoint_path"
      class="mt-3 pt-3 border-t border-gray-100 dark:border-dark-700/70"
    >
      <div class="flex items-center gap-1.5 text-xs text-gray-400 dark:text-gray-500">
        <svg class="w-3 h-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>
        <span class="font-mono text-gray-600 dark:text-gray-300">
          POST {{ model.endpoint_path }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { ModelPlazaEntry } from '@/types/models'

const props = defineProps<{
  model: ModelPlazaEntry
  platform: string
}>()

const { t } = useI18n()

const copied = ref(false)

async function copyModelId() {
  try {
    await navigator.clipboard.writeText(props.model.id)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 1500)
  } catch {
    // 回退方案：使用传统方法
    const textarea = document.createElement('textarea')
    textarea.value = props.model.id
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 1500)
  }
}

// Platform display configuration
const PLATFORM_CONFIG: Record<string, { label: string; gradient: string; badge: string }> = {
  anthropic: {
    label: 'Claude',
    gradient: 'bg-orange-50 dark:bg-orange-900/20',
    badge: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400',
  },
  openai: {
    label: 'OpenAI',
    gradient: 'bg-emerald-50 dark:bg-emerald-900/20',
    badge: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400',
  },
  gemini: {
    label: 'Gemini',
    gradient: 'bg-sky-50 dark:bg-sky-900/20',
    badge: 'bg-sky-100 text-sky-700 dark:bg-sky-900/30 dark:text-sky-400',
  },
  antigravity: {
    label: 'Antigravity',
    gradient: 'bg-indigo-50 dark:bg-indigo-900/20',
    badge: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400',
  },
  deepseek: {
    label: 'DeepSeek',
    gradient: 'bg-blue-50 dark:bg-blue-900/20',
    badge: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
  },
  moonshot: {
    label: 'Kimi',
    gradient: 'bg-violet-50 dark:bg-violet-900/20',
    badge: 'bg-violet-100 text-violet-700 dark:bg-violet-900/30 dark:text-violet-400',
  },
  zhipu: {
    label: 'Zhipu',
    gradient: 'bg-teal-50 dark:bg-teal-900/20',
    badge: 'bg-teal-100 text-teal-700 dark:bg-teal-900/30 dark:text-teal-400',
  },
  minimax: {
    label: 'MiniMax',
    gradient: 'bg-rose-50 dark:bg-rose-900/20',
    badge: 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-400',
  },
}

const config = computed(() => PLATFORM_CONFIG[props.platform] ?? PLATFORM_CONFIG.anthropic)
const platformLabel = computed(() => config.value.label)
const platformGradient = computed(() => config.value.gradient)
const platformBadgeClass = computed(() => config.value.badge)
const platformTintClass = computed(() => {
  const tints: Record<string, string> = {
    anthropic: 'text-orange-600 dark:text-orange-300',
    openai: 'text-emerald-600 dark:text-emerald-300',
    gemini: 'text-sky-600 dark:text-sky-300',
    antigravity: 'text-indigo-600 dark:text-indigo-300',
    deepseek: 'text-blue-600 dark:text-blue-300',
    moonshot: 'text-violet-600 dark:text-violet-300',
    zhipu: 'text-teal-600 dark:text-teal-300',
    minimax: 'text-rose-600 dark:text-rose-300',
  }
  return tints[props.platform] ?? 'text-gray-500 dark:text-gray-300'
})

// Simple SVG platform icon
const platformIcon = computed(() => ({
  render() {
    return h('svg', {
      xmlns: 'http://www.w3.org/2000/svg',
      width: 20,
      height: 20,
      viewBox: '0 0 24 24',
      fill: 'none',
      stroke: 'currentColor',
      'stroke-width': 2,
      'stroke-linecap': 'round' as const,
      'stroke-linejoin': 'round' as const,
    }, [
      h('rect', { x: 3, y: 3, width: 18, height: 18, rx: 2, ry: 2 }),
      h('line', { x1: 3, y1: 9, x2: 21, y2: 9 }),
      h('line', { x1: 9, y1: 21, x2: 9, y2: 9 }),
    ])
  },
}))
</script>

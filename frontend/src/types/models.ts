/** 模型广场相关类型 */

export interface ModelPlazaEntry {
  id: string
  display_name: string
  mode: string
  input_price_per_mtok?: number
  output_price_per_mtok?: number
  cache_read_price_per_mtok?: number
  supports_prompt_caching: boolean
  supports_vision: boolean
  provider: string
  endpoint_path: string
}

export interface PlatformModelGroup {
  platform: string
  display_name: string
  sort_order: number
  models: ModelPlazaEntry[]
}

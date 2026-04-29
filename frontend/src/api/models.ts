/**
 * 模型广场 API
 * GET /api/v1/models — 获取所有平台可用模型及定价信息
 */

import { apiClient } from './client'
import type { PlatformModelGroup } from '@/types/models'

export async function getModelPlaza(): Promise<PlatformModelGroup[]> {
  const { data } = await apiClient.get<{ platforms: PlatformModelGroup[] }>('/models')
  return data.platforms
}

export const modelsAPI = { getModelPlaza }

export default modelsAPI

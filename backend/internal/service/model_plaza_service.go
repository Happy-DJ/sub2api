package service

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/model"
)

// platformMeta 平台展示元数据。
// 模型列表从账号/渠道系统动态聚合，此处仅定义展示名和排序权重。
type platformMeta struct {
	DisplayName string
	SortOrder   int
}

var platformMetaMap = map[string]platformMeta{
	"anthropic":   {"Anthropic Claude", 1},
	"openai":      {"OpenAI", 2},
	"gemini":      {"Google Gemini", 3},
	"antigravity": {"Antigravity", 4},
	"deepseek":    {"DeepSeek", 5},
	"moonshot":    {"Kimi / Moonshot", 6},
	"zhipu":       {"智谱 GLM", 7},
	"minimax":     {"MiniMax", 8},
}

// ModelPlazaService 模型广场服务。
//
// 模型数据来源（双源聚合）：
//  1. 账号（Account）— 通过账号管理模块的 model_mapping 白名单获取平台支持的模型列表
//  2. 渠道（Channel）— 通过渠道的 ModelMapping + ModelPricing 获取模型及其定价
//
// 定价优先级：渠道定价 > LiteLLM 全局定价回落 > 推断。
type ModelPlazaService struct {
	accountRepo    AccountRepository // 账号管理模块：模型白名单来源
	channelService *ChannelService   // 渠道系统：模型定价来源
	pricing        *PricingService   // LiteLLM 全局定价回落
}

// NewModelPlazaService 创建模型广场服务。
func NewModelPlazaService(accountRepo AccountRepository, channelService *ChannelService, pricing *PricingService) *ModelPlazaService {
	return &ModelPlazaService{
		accountRepo:    accountRepo,
		channelService: channelService,
		pricing:        pricing,
	}
}

// modelAgg 模型聚合条目（去重前中间态）
type modelAgg struct {
	name        string
	platform    string
	fromAccount bool
	pricing     *ChannelModelPricing
}

// GetPlatformModels 获取所有平台的模型及其定价信息。
//
// 数据聚合逻辑：
//  1. 从活跃账号（Account）的 model_mapping 白名单中提取平台支持的模型列表
//  2. 从渠道（Channel）的 SupportedModels 补充模型及定价
//  3. 合并去重：同 platform × model name 去重，账号来源的模型名优先作为显示名
//  4. 按平台分组，按排序权重排列
func (s *ModelPlazaService) GetPlatformModels(ctx context.Context) ([]model.PlatformModelGroup, error) {
	// 聚合映射：platform → model name (lowercased) → modelAgg
	type dedupKey struct {
		platform string
		name     string // lowercased
	}
	agg := make(map[dedupKey]*modelAgg)

	// ——— 数据源 1：活跃账号的 model_mapping 白名单 ———
	accounts, accErr := s.accountRepo.ListActive(ctx)
	if accErr != nil {
		slog.Warn("[ModelPlaza] failed to list active accounts, will try channels only", "error", accErr)
	}
	for i := range accounts {
		acc := &accounts[i]
		platform := acc.Platform
		if platform == "" {
			continue
		}
		mapping := acc.GetModelMapping()
		if len(mapping) == 0 {
			continue
		}
		for modelName := range mapping {
			if modelName == "" || strings.HasPrefix(modelName, "*") {
				continue // 跳过空名和通配符模式
			}
			key := dedupKey{platform: platform, name: strings.ToLower(modelName)}
			if _, exists := agg[key]; !exists {
				agg[key] = &modelAgg{
					name:        modelName,
					platform:    platform,
					fromAccount: true,
				}
			}
		}
	}

	// ——— 数据源 2：渠道的 SupportedModels（补充模型 + 定价） ———
	channels, chErr := s.channelService.ListAvailable(ctx)
	if chErr != nil {
		slog.Warn("[ModelPlaza] failed to list available channels", "error", chErr)
	}
	for ci := range channels {
		for mi := range channels[ci].SupportedModels {
			sm := &channels[ci].SupportedModels[mi]
			if sm.Name == "" || sm.Platform == "" {
				continue
			}
			key := dedupKey{platform: sm.Platform, name: strings.ToLower(sm.Name)}
			if existing, exists := agg[key]; exists {
				// 已存在：补充定价（账号来源的模型优先无定价）
				if existing.pricing == nil && sm.Pricing != nil {
					existing.pricing = sm.Pricing
				}
			} else {
				agg[key] = &modelAgg{
					name:     sm.Name,
					platform: sm.Platform,
					pricing:  sm.Pricing,
				}
			}
		}
	}

	if len(agg) == 0 {
		// 两个数据源都失败且无数据时返回错误
		if accErr != nil && chErr != nil {
			return nil, accErr // 优先返回账号错误（更关键）
		}
		return []model.PlatformModelGroup{}, nil
	}

	// ——— 按平台分组 ———
	platformModels := make(map[string][]model.ModelPlazaEntry)
	for _, m := range agg {
		entry := s.aggToEntry(m)
		platformModels[m.platform] = append(platformModels[m.platform], entry)
	}

	// ——— 构建结果 ———
	groups := make([]model.PlatformModelGroup, 0, len(platformModels))
	for platform, entries := range platformModels {
		meta, ok := platformMetaMap[platform]
		if !ok {
			meta = platformMeta{DisplayName: platform, SortOrder: 99}
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].ID < entries[j].ID
		})

		groups = append(groups, model.PlatformModelGroup{
			Platform:    platform,
			DisplayName: meta.DisplayName,
			SortOrder:   meta.SortOrder,
			Models:      entries,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].SortOrder < groups[j].SortOrder
	})

	return groups, nil
}

// aggToEntry 将聚合条目转换为模型广场展示条目。
func (s *ModelPlazaService) aggToEntry(m *modelAgg) model.ModelPlazaEntry {
	entry := model.ModelPlazaEntry{
		ID:           m.name,
		DisplayName:  m.name,
		Provider:     m.platform,
		EndpointPath: endpointPathForPlatform(m.platform),
	}

	pricing := m.pricing
	// LiteLLM 全局定价回落
	if pricing == nil && s.pricing != nil {
		if lp := s.pricing.GetModelPricing(m.name); lp != nil {
			pricing = synthesizePricingFromLiteLLM(lp, pricing)
		}
	}

	if pricing != nil {
		entry = s.applyPricing(entry, pricing)
	}

	// 回退推断
	if entry.Mode == "" {
		entry.Mode = inferMode(m.name)
	}
	if !entry.SupportsVision {
		entry.SupportsVision = hasVisionSupport(m.name)
	}

	return entry
}

// applyPricing 将 ChannelModelPricing 填充到 ModelPlazaEntry。
// Token 定价自动转换为每百万 Token（×1e6）。
func (s *ModelPlazaService) applyPricing(entry model.ModelPlazaEntry, p *ChannelModelPricing) model.ModelPlazaEntry {
	switch p.BillingMode {
	case BillingModeToken:
		entry.Mode = "chat"
	case BillingModeImage:
		entry.Mode = "image"
	case BillingModePerRequest:
		entry.Mode = "chat"
	}

	if p.InputPrice != nil && *p.InputPrice > 0 {
		v := *p.InputPrice * 1e6
		entry.InputPricePerMTok = &v
	}
	if p.OutputPrice != nil && *p.OutputPrice > 0 {
		v := *p.OutputPrice * 1e6
		entry.OutputPricePerMTok = &v
	}
	if p.CacheReadPrice != nil && *p.CacheReadPrice > 0 {
		v := *p.CacheReadPrice * 1e6
		entry.CacheReadPricePerMTok = &v
	}

	if !entry.SupportsPromptCaching {
		entry.SupportsPromptCaching = p.CacheReadPrice != nil || p.CacheWritePrice != nil
	}

	return entry
}

// inferMode 根据模型 ID 推断模式。
func inferMode(id string) string {
	lower := strings.ToLower(id)
	switch {
	case strings.Contains(lower, "image") || strings.Contains(lower, "vision") || strings.Contains(lower, "cogview"):
		return "image"
	case strings.Contains(lower, "embed") || strings.Contains(lower, "text-embedding"):
		return "embedding"
	case strings.Contains(lower, "audio") || strings.Contains(lower, "tts") || strings.Contains(lower, "whisper"):
		return "audio"
	default:
		return "chat"
	}
}

// hasVisionSupport 根据模型 ID 判断是否支持视觉输入。
func hasVisionSupport(id string) bool {
	lower := strings.ToLower(id)
	visionKeywords := []string{"vision", "image", "flash-image", "cogview", "pixtral", "gpt-image", "-v "}
	for _, kw := range visionKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// endpointPathForPlatform 根据平台返回 API 端点路径。
// Anthropic 协议平台（Claude/Antigravity/Gemini）使用 Messages API，
// OpenAI 协议平台使用 Chat Completions API。
func endpointPathForPlatform(platform string) string {
	switch platform {
	case "anthropic", "antigravity", "gemini":
		return "/v1/messages"
	default:
		return "/v1/chat/completions"
	}
}

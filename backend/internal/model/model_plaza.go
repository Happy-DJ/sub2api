// Package model 定义服务层使用的数据模型。
package model

// PlatformModelGroup 模型广场中某个平台下的所有模型分组
type PlatformModelGroup struct {
	Platform    string            `json:"platform"`
	DisplayName string            `json:"display_name"`
	SortOrder   int               `json:"sort_order"`
	Models      []ModelPlazaEntry `json:"models"`
}

// ModelPlazaEntry 模型广场中的单个模型卡片条目
type ModelPlazaEntry struct {
	ID                    string   `json:"id"`
	DisplayName           string   `json:"display_name"`
	Mode                  string   `json:"mode"`
	InputPricePerMTok     *float64 `json:"input_price_per_mtok,omitempty"`
	OutputPricePerMTok    *float64 `json:"output_price_per_mtok,omitempty"`
	CacheReadPricePerMTok *float64 `json:"cache_read_price_per_mtok,omitempty"`
	SupportsPromptCaching bool     `json:"supports_prompt_caching"`
	SupportsVision        bool     `json:"supports_vision"`
	Provider              string   `json:"provider"`
	EndpointPath          string   `json:"endpoint_path"` // API 端点路径，如 "/v1/messages" 或 "/v1/chat/completions"
}

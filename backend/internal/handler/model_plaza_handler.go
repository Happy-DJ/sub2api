package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ModelPlazaHandler 处理模型广场相关请求
type ModelPlazaHandler struct {
	svc *service.ModelPlazaService
}

// NewModelPlazaHandler 创建模型广场 handler
func NewModelPlazaHandler(svc *service.ModelPlazaService) *ModelPlazaHandler {
	return &ModelPlazaHandler{svc: svc}
}

// List 列出所有平台可用模型及定价信息
// GET /api/v1/models
func (h *ModelPlazaHandler) List(c *gin.Context) {
	_, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	platforms, err := h.svc.GetPlatformModels(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"platforms": platforms})
}

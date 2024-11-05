package handler

import (
	"01.alem.school/git/nyeltay/forum/internal/service"
	"01.alem.school/git/nyeltay/forum/internal/template_cache"
	"01.alem.school/git/nyeltay/forum/pkg/logger"
)

type Handler struct {
	service       *service.Service
	templateCache template_cache.TemplateCache
	logger        *logger.Logger
}

func NewHandler(service *service.Service, templateCache template_cache.TemplateCache, logger *logger.Logger) *Handler {
	return &Handler{
		service:       service,
		templateCache: templateCache,
		logger:        logger,
	}
}

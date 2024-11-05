package handler

import (
	"net/http"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/validator"
)

type PageData struct {
	Validator         *validator.Validator
	AuthenticatedUser *models.User
	Posts             []models.Post
	Post              models.Post
	Categories        []models.Category
	Comments          []models.Comment
	Error             string
}

func (h *Handler) Render(w http.ResponseWriter, page string, data *PageData) {
	ts, ok := h.templateCache[page]
	if !ok {
		http.Error(w, "Template Not Found", http.StatusInternalServerError)
		//h.logger.
		return
	}

	err := ts.ExecuteTemplate(w, page, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//h.logger.Errorf("execute template: %v", err)
		return
	}
}

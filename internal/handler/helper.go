package handler

import (
	"net/http"
)

type H map[string]any

func (h *Handler) Render(w http.ResponseWriter, page string, obj any) {
	ts, ok := h.templateCache[page]
	if !ok {
		h.logger.Error("template not found:", page)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	err := ts.ExecuteTemplate(w, page, obj)
	if err != nil {
		h.logger.Errorf("execute template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

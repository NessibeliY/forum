package handler

import (
	"net/http"
)

type H map[string]any

func (h *Handler) Render(w http.ResponseWriter, page string, obj any) {
	ts, ok := h.templateCache[page]
	if !ok {
		http.Error(w, "Template Not Found", http.StatusInternalServerError)
		//h.logger.
		return
	}

	err := ts.ExecuteTemplate(w, page, obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//h.logger.Errorf("execute template: %v", err)
		return
	}
}

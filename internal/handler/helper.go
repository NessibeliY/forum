package handler

import (
	"net/http"
)

type H map[string]any

func (h *Handler) Render(w http.ResponseWriter, page string, status int, obj any) {
	ts, ok := h.templateCache[page]
	if !ok {
		h.logger.Error("template not found:", page)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, page, obj)
	if err != nil {
		h.logger.Errorf("execute template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) clientError(w http.ResponseWriter, status int) {
	h.Render(w, "error.page.html", status, H{
		"error_text": http.StatusText(status),
		"error_code": status,
	})
}

func (h *Handler) serverError(w http.ResponseWriter, err error) {
	h.Render(w, "error.page.html", http.StatusInternalServerError, H{
		"error_text": err.Error(),
		"error_code": http.StatusInternalServerError,
	})
}

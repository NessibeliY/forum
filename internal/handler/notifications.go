package handler

import "net/http"

func (h *Handler) Notification(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/notifications" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	categories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		h.logger.Info("get all categories:", err)
		h.serverError(w, err)
		return
	}

	h.Render(w, "notification.page.html", http.StatusOK, H{
		"categories":         categories,
		"authenticated_user": h.getUserFromContext(r),
	})
}

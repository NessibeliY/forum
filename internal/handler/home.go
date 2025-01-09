package handler

import (
	"net/http"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	posts, err := h.service.PostService.GetAllPosts()
	if err != nil {
		h.logger.Info("get all posts:", err)
		h.serverError(w, err)
		return
	}

	categories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		h.logger.Info("get all categories:", err)
		h.serverError(w, err)
		return
	}

	var countNotification int
	var checkModeratorRequest bool
	if h.getUserFromContext(r) != nil {
		countNotification, err = h.service.NotificationService.GetCountNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get countNotification:", err)
			h.serverError(w, err)
			return
		}

		checkModeratorRequest, err = h.service.UserService.CheckModeratorRequestStatus(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get check moderator request:", err)
			h.serverError(w, err)
			return
		}
	}

	h.Render(w, "index.page.html", http.StatusOK, H{
		"posts":                   posts,
		"categories":              categories,
		"authenticated_user":      h.getUserFromContext(r),
		"count_notification":      countNotification,
		"check_moderator_request": checkModeratorRequest,
	})
}

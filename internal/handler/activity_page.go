package handler

import (
	"net/http"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

func (h *Handler) ActivityPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/activity-page" {
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

	var countNotification int
	var currentNotifications []models.Notification
	var archivedNotifications []models.Notification
	var userPosts []models.Post
	var likedDislikedPosts []models.UserReactionPost
	if h.getUserFromContext(r) != nil {
		countNotification, err = h.service.NotificationService.GetCountNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get countNotification:", err)
			h.serverError(w, err)
			return
		}

		currentNotifications, err = h.service.NotificationService.GetCurrentNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get GetListNotifications:", err)
			h.serverError(w, err)
			return
		}

		archivedNotifications, err = h.service.NotificationService.GetArchivedNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get archived notifications:", err)
			h.serverError(w, err)
			return
		}

		userPosts, err = h.service.PostService.GetPostsByAuthorID(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get archived notifications:", err)
			h.serverError(w, err)
			return
		}

		likedDislikedPosts, err = h.service.PostReactionService.GetUserReactionPosts(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get archived notifications:", err)
			h.serverError(w, err)
			return
		}

	}

	h.Render(w, "activity_page.page.html", http.StatusOK, H{
		"categories":             categories,
		"authenticated_user":     h.getUserFromContext(r),
		"count_notification":     countNotification,
		"current_notifications":  currentNotifications,
		"archived_notifications": archivedNotifications,
		"user_posts":             userPosts,
		"likes_dislikes_post":    likedDislikedPosts,
	})
}

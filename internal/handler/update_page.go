package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"01.alem.school/git/nyeltay/forum/pkg/utils"
)

func (h *Handler) UpdatePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/update" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.UpdatePageMethodGet(w, r)
	default:
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) UpdatePageMethodGet(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		h.logger.Error("get all categories:", err.Error())
		h.serverError(w, err)
		return
	}

	var countNotification int
	if h.getUserFromContext(r) != nil {
		countNotification, err = h.service.NotificationService.GetCountNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get countNotification:", err)
			h.serverError(w, err)
			return
		}
	}

	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		h.logger.Error("parse query:", err.Error())
		h.clientError(w, http.StatusBadRequest)
		return
	}

	postIDStr := query.Get("id")

	if len(query) != 1 || postIDStr == "" {
		h.logger.Error("query or post_id invalid")
		h.clientError(w, http.StatusBadRequest)
		return
	}

	postID, err := utils.ParsePositiveIntID(postIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		h.clientError(w, http.StatusNotFound)
		return
	}

	post, err := h.service.PostService.GetPostByID(postID)
	if err != nil {
		h.logger.Error("get post by id:", err.Error())
		h.serverError(w, err)
		return
	}

	if post == nil {
		h.logger.Error("post nil")
		h.clientError(w, http.StatusNotFound)
		return
	}

	fmt.Println("post", post.Categories)

	h.Render(w, "update_page.page.html", http.StatusOK, H{
		"categories":         categories,
		"authenticated_user": h.getUserFromContext(r),
		"count_notification": countNotification,
		"post":               post,
	})
}

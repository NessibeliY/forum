package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"01.alem.school/git/nyeltay/forum/internal/models"
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
	case http.MethodPost:
		h.UpdatePostMethodPost(w, r)
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

	var countNotification int
	if h.getUserFromContext(r) != nil {
		countNotification, err = h.service.NotificationService.GetCountNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get countNotification:", err)
			h.serverError(w, err)
			return
		}
	}

	h.Render(w, "update_page.page.html", http.StatusOK, H{
		"categories":         categories,
		"authenticated_user": h.getUserFromContext(r),
		"count_notification": countNotification,
		"post":               post,
	})
}

func (h *Handler) UpdatePostMethodPost(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.PostFormValue("title"))
	content := strings.TrimSpace(r.PostFormValue("content"))
	categoryNames := r.PostForm["categories"]

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

	var countNotification int
	if h.getUserFromContext(r) != nil {
		countNotification, err = h.service.NotificationService.GetCountNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get countNotification:", err)
			h.serverError(w, err)
			return
		}

	}

	validationsErrMap := validateCreatePostForm(title, content, categoryNames)
	if len(validationsErrMap) > 0 {
		h.logger.Error("validate create post form:", validationsErrMap)

		categories, err := h.service.CategoryService.GetAllCategories()
		if err != nil {
			h.logger.Error("get all categories:", err.Error())
			h.serverError(w, err)
			return
		}

		h.Render(w, "update_page.page.html", http.StatusBadRequest, H{
			"errors":             validationsErrMap,
			"categories":         categories,
			"authenticated_user": h.getUserFromContext(r),
			"count_notification": countNotification,
			"post":               post,
		})
		return
	}

	categories := make([]*models.Category, 0, len(categoryNames))
	for _, categoryName := range categoryNames {
		c, err := h.service.CategoryService.GetCategoryByName(categoryName)
		if err != nil {
			h.logger.Error("get category by name:", err.Error())
			h.serverError(w, err)
			return
		}
		if c == nil {
			validationsErrMap["categories"] = "category not found"
			h.logger.Error("category not found:", c)

			categories, err := h.service.CategoryService.GetAllCategories()
			if err != nil {
				h.logger.Error("get all categories:", err.Error())
				h.serverError(w, err)
				return
			}

			h.Render(w, "update_page.page.html", http.StatusBadRequest, H{
				"errors":             validationsErrMap,
				"categories":         categories,
				"authenticated_user": h.getUserFromContext(r),
				"count_notification": countNotification,
				"post":               post,
			})
			return
		}
		categories = append(categories, c)
	}

	image, err := h.handleImageUpload(r)
	if err != nil {
		h.logger.Error("handle image upload:", err.Error())
		if strings.Contains(err.Error(), "not an image") || strings.Contains(err.Error(), "image too big") {
			h.clientError(w, http.StatusBadRequest)
			return
		}
		h.serverError(w, err)
		return
	}

	updatePost := &models.UpdatePostRequest{
		Title:      title,
		Content:    content,
		AuthorID:   h.getUserFromContext(r).ID,
		Categories: categories,
		ImageFile:  image,
	}

	id, err := h.service.PostService.UpdatePostWithImage(updatePost)
	if err != nil {
		h.logger.Error("update post:", err.Error())
		h.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
}

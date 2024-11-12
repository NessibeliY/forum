package handler

import (
	"fmt"
	"net/http"
	"strings"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.createPostMethodGet(w, r)
	case http.MethodPost:
		h.createPostMethodPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) createPostMethodGet(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		//h.logger
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Render(w, "create.page.html", H{
		"categories":         categories,
		"authenticated_user": h.getUserFromContext(r),
	})
}

func (h *Handler) createPostMethodPost(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.PostFormValue("title"))
	content := strings.TrimSpace(r.PostFormValue("content"))
	categoryNames := r.PostForm["categories"]

	validationsErrMap := validateCreatePostForm(title, content, categoryNames)
	if validationsErrMap != nil {
		categories, err := h.service.CategoryService.GetAllCategories()
		if err != nil {
			//h.logger
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.Render(w, "create.page.html", H{
			"errors":             validationsErrMap,
			"categories":         categories,
			"authenticated_user": h.getUserFromContext(r),
		})
		return
	}

	categories := make([]*models.Category, 0, len(categoryNames))
	for _, categoryName := range categoryNames {
		c, err := h.service.CategoryService.GetCategoryByName(categoryName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if c == nil {
			validationsErrMap["categories"] = "category not found"
			//h.logger
			w.WriteHeader(http.StatusBadRequest)
			h.Render(w, "create.page.html", H{
				"errors":             validationsErrMap,
				"categories":         categories,
				"authenticated_user": h.getUserFromContext(r),
			})
		}
		categories = append(categories, c)
	}

	createPostRequest := &models.CreatePostRequest{
		Title:      title,
		Content:    content,
		AuthorID:   h.getUserFromContext(r).ID,
		Categories: categories,
	}

	id, err := h.service.PostService.CreatePost(createPostRequest)
	if err != nil {
		//h.logger
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusFound)
}

func validateCreatePostForm(title, content string, categoryNames []string) map[string]string {
	errors := make(map[string]string)

	if title == "" {
		errors["title"] = "title is required"
	}

	if content == "" {
		errors["content"] = "content is required"
	}

	if len(categoryNames) == 0 {
		errors["categories"] = "at least one category must be selected"
	}

	if len(title) > 100 {
		errors["title"] = "title cannot exceed 100 characters"
	}

	if len(content) > 10000 {
		errors["content"] = "content cannot exceed 10,000 characters"
	}

	return errors
}

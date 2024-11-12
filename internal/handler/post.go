package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/post/") {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		//h.logger
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	postIDStr := query.Get("post-id")

	if len(query) != 1 || postIDStr == "" {
		http.Error(w, "query must only contain 'post-id'", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID < 1 {
		http.NotFound(w, r)
		return
	}

	post, err := h.service.PostService.GetPostByID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.NotFound(w, r)
		return
	}

	comments, err := h.service.CommentService.GetAllCommentsByPostID(postID)
	if err != nil {
		//h.logger
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, comment := range comments {
		comment.LikesCount, comment.DislikesCount, err = h.service.CommentReactionService.GetCommentLikesAndDislikesByID(comment.ID)
		if err != nil {
			//h.logger
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	post.LikesCount, post.DislikesCount, err = h.service.PostReactionService.GetPostLikesAndDislikesByID(postID)
	if err != nil {
		//h.logger
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Render(w, "post.page.html", H{
		"post":               post,
		"comments":           comments,
		"authenticated_user": h.getUserFromContext(r),
	})
}

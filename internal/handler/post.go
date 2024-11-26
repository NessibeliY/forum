package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/utils"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.createPostMethodGet(w, r)
	case http.MethodPost:
		h.createPostMethodPost(w, r)
	default:
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) createPostMethodGet(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		h.logger.Error("get all categories:", err.Error())
		h.serverError(w, err)
		return
	}

	h.Render(w, "create_post.page.html", http.StatusOK, H{
		"categories":         categories,
		"authenticated_user": h.getUserFromContext(r),
	})
}

func (h *Handler) createPostMethodPost(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.PostFormValue("title"))
	content := strings.TrimSpace(r.PostFormValue("content"))
	categoryNames := r.PostForm["categories"]

	validationsErrMap := validateCreatePostForm(title, content, categoryNames)
	if len(validationsErrMap) > 0 {
		h.logger.Error("validate create post form:", validationsErrMap)

		categories, err := h.service.CategoryService.GetAllCategories()
		if err != nil {
			h.logger.Error("get all categories:", err.Error())
			h.serverError(w, err)
			return
		}

		h.Render(w, "create_post.page.html", http.StatusBadRequest, H{
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

			h.Render(w, "create_post.page.html", http.StatusBadRequest, H{
				"errors":             validationsErrMap,
				"categories":         categories,
				"authenticated_user": h.getUserFromContext(r),
			})
			return
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
		h.logger.Error("create post:", err.Error())
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
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

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/delete" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.URL.Query().Get("id")
	if postIDStr == "" {
		h.logger.Error("get query for post_id")
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
		h.logger.Info("post nil")
		h.clientError(w, http.StatusNotFound)
		return
	}

	deletePostRequest := &models.DeletePostRequest{
		ID: postID,
	}

	err = h.service.PostService.DeletePost(deletePostRequest)
	if err != nil {
		h.logger.Error("delete post:", err.Error())
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
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

	comments, err := h.service.CommentService.GetAllCommentsByPostID(postID)
	if err != nil {
		h.logger.Error("get all comments by post_id:", err.Error())
		h.serverError(w, err)
		return
	}

	for _, comment := range comments {
		comment.LikesCount, comment.DislikesCount, err = h.service.CommentReactionService.GetCommentLikesAndDislikesByID(comment.ID)
		if err != nil {
			h.logger.Error("get comment likes and dislikes by id:", err.Error())
			h.serverError(w, err)
			return
		}
	}

	post.LikesCount, post.DislikesCount, err = h.service.PostReactionService.GetPostLikesAndDislikesByID(postID)
	if err != nil {
		h.logger.Error("get post likes and dislikes by id:", err.Error())
		h.serverError(w, err)
		return
	}

	categories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		h.logger.Info("get all categories:", err)
		h.serverError(w, err)
		return
	}

	h.Render(w, "post.page.html", http.StatusOK, H{
		"post":               post,
		"comments":           comments,
		"categories":         categories,
		"authenticated_user": h.getUserFromContext(r),
	})
}

func (h *Handler) ShowMyPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/myposts" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	posts, err := h.service.PostService.GetPostsByAuthorID(h.getUserFromContext(r).ID)
	if err != nil {
		h.logger.Error("get posts by author_id:", err.Error())
		h.serverError(w, err)
		return
	}

	categories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		h.logger.Info("get all categories:", err)
		h.serverError(w, err)
		return
	}

	h.Render(w, "index.page.html", http.StatusOK, H{
		"posts":              posts,
		"categories":         categories,
		"authenticated_user": h.getUserFromContext(r),
	})
}

func (h *Handler) ShowLikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/likedposts" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	posts, err := h.service.PostService.GetLikedPosts(h.getUserFromContext(r).ID)
	if err != nil {
		h.logger.Error("get liked posts:", err.Error())
		h.serverError(w, err)
		return
	}

	categories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		h.logger.Info("get all categories:", err)
		h.serverError(w, err)
		return
	}

	h.Render(w, "index.page.html", http.StatusOK, H{
		"posts":              posts,
		"categories":         categories,
		"authenticated_user": h.getUserFromContext(r),
	})
}

func (h *Handler) ShowPostsByCategory(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/showposts" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()
	categories := queryParams["category"]

	if len(categories) == 0 {
		h.logger.Error("categories empty")
		h.clientError(w, http.StatusBadRequest)
		return
	}

	for _, categoryName := range categories {
		c, err := h.service.CategoryService.GetCategoryByName(categoryName)
		if err != nil {
			h.logger.Error("get category by name:", err.Error())
			h.serverError(w, err)
			return
		}
		if c == nil {
			h.clientError(w, http.StatusNotFound)
			return
		}
	}

	posts, err := h.service.PostService.GetPostsByCategories(categories)
	if err != nil {
		h.logger.Error("get posts by categories:", err.Error())
		h.serverError(w, err)
		return
	}

	allCategories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		h.logger.Error("get all categories:", err.Error())
		h.serverError(w, err)
		return
	}

	h.Render(w, "index.page.html", http.StatusOK, H{
		"posts":              posts,
		"categories":         allCategories,
		"authenticated_user": h.getUserFromContext(r),
	})
}

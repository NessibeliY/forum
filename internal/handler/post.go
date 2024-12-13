package handler

import (
	"errors"
	"fmt"
	"io"
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

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		createPostRequest.ImageFile = nil
		if errors.Is(err, http.ErrMissingFile) ||
			err == io.EOF ||
			(err.Error() == "http: no such file") ||
			strings.Contains(err.Error(), "multipart") {
			categories, err := h.service.CategoryService.GetAllCategories()
			if err != nil {
				h.logger.Error("get all categories:", err.Error())
				h.serverError(w, err)
				return
			}

			validationsErrMap["image"] = err.Error()
			h.logger.Error("form file client error:", err.Error())

			h.Render(w, "create_post.page.html", http.StatusBadRequest, H{
				"errors":             validationsErrMap,
				"categories":         categories,
				"authenticated_user": h.getUserFromContext(r),
			})
			return
		}

		h.logger.Error("form file:", err.Error())
		h.serverError(w, err)
		return
	}

	createPostRequest.ImageFile = file
	defer file.Close()

	fileType := fileHeader.Header.Get("Content-Type")
	if !isImage(fileType) {
		categories, err := h.service.CategoryService.GetAllCategories()
		if err != nil {
			h.logger.Error("get all categories:", err.Error())
			h.serverError(w, err)
			return
		}

		validationsErrMap["image"] = "file is not an image"
		h.logger.Error("file is not an image:", fileType)

		h.Render(w, "create_post.page.html", http.StatusBadRequest, H{
			"errors":             validationsErrMap,
			"categories":         categories,
			"authenticated_user": h.getUserFromContext(r),
		})
		return
	}

	if fileHeader.Size > 5*1024*1024 {
		categories, err := h.service.CategoryService.GetAllCategories()
		if err != nil {
			h.logger.Error("get all categories:", err.Error())
			h.serverError(w, err)
			return
		}

		validationsErrMap["image"] = "file is too big"
		h.logger.Error("file is too big:", fileType)

		h.Render(w, "create_post.page.html", http.StatusBadRequest, H{
			"errors":             validationsErrMap,
			"categories":         categories,
			"authenticated_user": h.getUserFromContext(r),
		})
		return
	}

	id, err := h.service.PostService.CreatePostWithImage(createPostRequest)
	if err != nil {
		h.logger.Error("create post:", err.Error())
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
}

func isImage(fileType string) bool {
	allowedTypes := "image/jpeg,image/png,image/gif,image/jpg,image/webp"
	return strings.Contains(allowedTypes, fileType)
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

func (h *Handler) SendReport(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/report" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.sendReportMethodGet(w, r)
	case http.MethodPost:
		h.sendReportMethodPost(w, r)
	default:
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) sendReportMethodGet(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.logger.Error("get query for user_id")
		h.clientError(w, http.StatusBadRequest)
		return
	}

	userID, err := utils.ParsePositiveIntID(userIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		h.clientError(w, http.StatusNotFound)
		return
	}

	user, err := h.service.UserService.GetUserByID(userID)
	if err != nil {
		h.logger.Error("get user by id:", err.Error())
		h.serverError(w, err)
		return
	}
	if user == nil {
		h.logger.Info("user nil")
		h.clientError(w, http.StatusNotFound)
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
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

	h.Render(w, "report.page.html", http.StatusOK, H{
		"reported_user":      user,
		"reported_post":      post,
		"authenticated_user": h.getUserFromContext(r),
	})
}

func (h *Handler) sendReportMethodPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
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

	reason := strings.TrimSpace(r.PostFormValue("report_text"))
	err = validateReportText(reason)
	if err != nil {
		h.logger.Error("validate report text:", err)

		categories, err := h.service.CategoryService.GetAllCategories()
		if err != nil {
			h.logger.Error("get all categories:", err.Error())
			h.serverError(w, err)
			return
		}

		h.Render(w, "report.page.html", http.StatusBadRequest, H{
			"error":              err,
			"categories":         categories,
			"authenticated_user": h.getUserFromContext(r),
		})
		return
	}

	err = h.service.PostService.SendReport(models.SendReportRequest{
		PostID:      postID,
		Reason:      reason,
		ModeratorID: h.getUserFromContext(r).ID,
	})
	if err != nil {
		h.logger.Error("get post by id:", err.Error())
		h.serverError(w, err)
		return
	}
}

func validateReportText(reportText string) error {
	if reportText == "" {
		return errors.New("report text is empty")
	}

	if len(reportText) > 1000 {
		return errors.New("report text is too long")
	}

	return nil
}

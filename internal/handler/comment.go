package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/cookies"
	"01.alem.school/git/nyeltay/forum/pkg/utils"
)

const errorsMapCookieName = "forum_errors_map_cookie"

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/create" {
		h.logger.Error("url path:", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		h.logger.Errorf("method not allowed: %s", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parse form:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := strings.TrimSpace(r.PostFormValue("content"))
	postIDStr := strings.TrimSpace(r.PostFormValue("post_id"))

	validationsErrMap := validateCreateCommentForm(content, postIDStr)
	if len(validationsErrMap) > 0 {
		h.logger.Error("validate create comment form:", validationsErrMap)
		errorsJSON, err := json.Marshal(validationsErrMap)
		if err != nil {
			h.logger.Error("marshal:", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cookies.SetCookie(w, errorsMapCookieName, string(errorsJSON), 300)

		http.Redirect(w, r, fmt.Sprintf("/post?id=%s", postIDStr), http.StatusSeeOther)
		return
	}

	postID, err := utils.ParsePositiveIntID(postIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		http.NotFound(w, r)
		return
	}

	createCommentRequest := &models.CreateCommentRequest{
		Content:  content,
		AuthorID: h.getUserFromContext(r).ID,
		PostID:   postID,
	}

	err = h.service.CommentService.CreateComment(createCommentRequest)
	if err != nil {
		h.logger.Error("create comment:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
}

func validateCreateCommentForm(content string, postIDStr string) map[string]string {
	errors := make(map[string]string)

	if content == "" {
		errors["content"] = "content is required"
	}

	if postIDStr == "" {
		errors["post_id"] = "post id is required"
	}

	if len(content) > 10000 {
		errors["content"] = "content must be less than 10000 characters"
	}

	return errors
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/delete" {
		h.logger.Error("url path:", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodDelete {
		h.logger.Errorf("method not allowed: %s", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	commentIDStr := r.URL.Query().Get("comment_id")
	if commentIDStr == "" {
		h.logger.Error("comment id is required:", commentIDStr)
		http.Error(w, "comment_id is required", http.StatusBadRequest)
		return
	}

	commentID, err := utils.ParsePositiveIntID(commentIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deleteCommentRequest := &models.DeleteCommentRequest{
		ID: commentID,
	}

	err = h.service.CommentService.DeleteComment(deleteCommentRequest)
	if err != nil {
		h.logger.Error("delete comment:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		h.logger.Error("get url query:", err.Error())
		http.Error(w, "post_id is required", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%s", postIDStr), http.StatusSeeOther)
}

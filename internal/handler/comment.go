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
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parse form:", err.Error())
		h.clientError(w, http.StatusBadRequest)
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
			h.serverError(w, err)
			return
		}

		cookies.SetCookie(w, errorsMapCookieName, string(errorsJSON), 300)

		http.Redirect(w, r, fmt.Sprintf("/post?id=%s", postIDStr), http.StatusSeeOther)
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
		h.logger.Error("post doesn't exist: post nil")
		h.clientError(w, http.StatusNotFound)
		return
	}

	createCommentRequest := &models.CreateCommentRequest{
		Content:  content,
		AuthorID: h.getUserFromContext(r).ID,
		PostID:   postID,
	}

	err = h.service.CommentService.CreateComment(createCommentRequest)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			h.logger.Error("post doesn't exist:", err.Error())
			h.clientError(w, http.StatusNotFound)
			return
		}
		h.logger.Error("create comment:", err.Error())
		h.serverError(w, err)
		return
	}

	currentUser := h.getUserFromContext(r)
	if currentUser == nil {
		h.logger.Error("user not authenticated")
		h.clientError(w, http.StatusUnauthorized)
		return
	}

	if post.AuthorID != currentUser.ID {
		notificationRequst := &models.NotificationRequest{
			PostID:  postID,
			Message: content,
		}
		_, err = h.service.NotificationService.CreateNotification(notificationRequst)
		if err != nil {
			h.logger.Error("create notifications:", err.Error())
			h.serverError(w, err)
			return
		}
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
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	commentIDStr := r.URL.Query().Get("comment_id")
	if commentIDStr == "" {
		h.logger.Error("comment id is required:", commentIDStr)
		h.clientError(w, http.StatusBadRequest)
		return
	}

	commentID, err := utils.ParsePositiveIntID(commentIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		h.clientError(w, http.StatusNotFound)
		return
	}

	deleteCommentRequest := &models.DeleteCommentRequest{
		ID: commentID,
	}

	err = h.service.CommentService.DeleteComment(deleteCommentRequest)
	if err != nil {
		h.logger.Error("delete comment:", err.Error())
		h.serverError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.logger.Error("parse form:", err.Error())
		h.clientError(w, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/update" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parse form:", err.Error())
		h.clientError(w, http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.PostFormValue("content"))
	commentIDStr := strings.TrimSpace(r.PostFormValue("comment_id"))

	postIDStr := r.URL.Query().Get("id")
	if postIDStr == "" {
		h.logger.Error("get query for post_id")
		h.clientError(w, http.StatusBadRequest)
		return
	}

	validationsErrMap := validateCreateCommentForm(content, postIDStr)
	if len(validationsErrMap) > 0 {
		h.logger.Error("validate create comment form:", validationsErrMap)
		errorsJSON, err := json.Marshal(validationsErrMap)
		if err != nil {
			h.logger.Error("marshal:", err.Error())
			h.serverError(w, err)
			return
		}

		cookies.SetCookie(w, errorsMapCookieName, string(errorsJSON), 300)

		http.Redirect(w, r, fmt.Sprintf("/post?id=%s", postIDStr), http.StatusSeeOther)
		return
	}

	commentID, err := utils.ParsePositiveIntID(commentIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		h.clientError(w, http.StatusNotFound)
		return
	}

	checkCommentID, err := h.service.CommentService.GetCommentByID(commentID)
	if err != nil {
		h.logger.Error("get comment id:", err.Error())
		h.serverError(w, err)
		return
	}

	if checkCommentID == nil {
		h.logger.Error("check comment id:", "bad request")
		h.clientError(w, http.StatusBadRequest)
		return
	}

	updateCommentReq := &models.UpdateCommentRequest{
		ID:      commentID,
		Content: content,
	}

	err = h.service.CommentService.UpdateComment(updateCommentReq)
	if err != nil {
		h.logger.Error("update comment:", err.Error())
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

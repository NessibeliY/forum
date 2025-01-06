package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/utils"
)

func (h *Handler) CreatePostReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/reaction/create" {
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

	redirectTo := "/"
	referer := r.Header.Get("Referer")
	parsedURL, err := url.Parse(referer)
	if err == nil && parsedURL.RawQuery != "" {
		redirectTo = parsedURL.Path + "?" + parsedURL.RawQuery
	}

	if !h.isValidRedirectTo(redirectTo) {
		h.logger.Error("invalid redirect to:", redirectTo)
		redirectTo = "/"
	}

	postIDStr := strings.TrimSpace(r.PostFormValue("post_id"))
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

	reaction := strings.TrimSpace(r.PostFormValue("reaction"))
	err = validateCreatePostReactionForm(reaction)
	if err != nil {
		h.logger.Error("validate create post reaction form:", err.Error())
		h.clientError(w, http.StatusBadRequest)
		return
	}

	createPostReactionRequest := &models.CreatePostReactionRequest{
		AuthorID: h.getUserFromContext(r).ID,
		PostID:   postID,
		Reaction: reaction,
	}

	err = h.service.PostReactionService.CreatePostReaction(createPostReactionRequest)
	if err != nil {
		h.logger.Error("create post reaction:", err.Error())
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
		notificationRequest := &models.NotificationRequest{
			PostID:  postID,
			Message: reaction,
		}

		_, err = h.service.NotificationService.CreateNotification(notificationRequest)
		if err != nil {
			h.logger.Error("create post notification:", err.Error())
			h.serverError(w, err)
			return
		}
	}

	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}

func validateCreatePostReactionForm(reaction string) error {
	if reaction != "like" && reaction != "dislike" {
		return fmt.Errorf("reaction must either like or dislike")
	}

	return nil
}

func (h *Handler) isValidRedirectTo(redirectTo string) bool {
	if redirectTo == "/" {
		return true
	}

	matched, err := regexp.MatchString(`^\/post\?id=[1-9]\d*$`, redirectTo)
	if err != nil {
		h.logger.Error("match string error:", err.Error())
		return false
	}
	if !matched {
		h.logger.Error("not matched")
		return false
	}

	parsedURL, err := url.Parse(redirectTo)
	if err != nil {
		h.logger.Error("url parse:", err.Error())
		return false
	}

	postIDStr := parsedURL.Query().Get("id")
	postID, err := utils.ParsePositiveIntID(postIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		return false
	}

	post, err := h.service.PostService.GetPostByID(postID)
	if err != nil {
		h.logger.Error("get post by id:", err.Error())
		return false
	}
	if post == nil {
		h.logger.Error("post nil")
		return false
	}

	return true
}

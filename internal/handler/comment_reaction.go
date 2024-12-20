package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/utils"
)

func (h *Handler) CreateCommentReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/reaction/create" {
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

	redirectTo := r.Header.Get("Referer")
	parsedURL, err := url.Parse(redirectTo)
	if err != nil {
		h.logger.Error("url parse:", err.Error())
		redirectTo = "/"
	} else {
		redirectTo = parsedURL.Path + "?" + parsedURL.RawQuery
	}

	if !h.isValidRedirectTo(redirectTo) {
		h.logger.Error("invalid redirect to:", redirectTo)
		redirectTo = "/"
	}

	commentIDStr := strings.TrimSpace(r.PostFormValue("comment_id"))
	commentID, err := utils.ParsePositiveIntID(commentIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		h.clientError(w, http.StatusNotFound)
		return
	}

	comment, err := h.service.CommentService.GetCommentByID(commentID)
	if err != nil {
		h.logger.Error("get comment by id:", err.Error())
		h.serverError(w, err)
		return
	}
	if comment == nil {
		h.logger.Error("comment doesn't exist: comment nil")
		h.clientError(w, http.StatusNotFound)
		return
	}

	reaction := r.PostFormValue("reaction")
	err = validateCreateCommentReactionForm(reaction)
	if err != nil {
		h.logger.Error("validate create comment reaction form:", err.Error())
		h.clientError(w, http.StatusBadRequest)
		return
	}

	createCommentReactionRequest := &models.CreateCommentReactionRequest{
		AuthorID:  h.getUserFromContext(r).ID,
		Reaction:  reaction,
		CommentID: commentID,
	}

	err = h.service.CommentReactionService.CreateCommentReaction(createCommentReactionRequest)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			h.logger.Error("comment doesn't exist:", err.Error())
			h.clientError(w, http.StatusNotFound)
			return
		}
		h.logger.Error("create comment reaction:", err.Error())
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}

func validateCreateCommentReactionForm(reaction string) error {
	if reaction != "like" && reaction != "dislike" {
		return fmt.Errorf("reaction must either like or dislike")
	}

	return nil
}

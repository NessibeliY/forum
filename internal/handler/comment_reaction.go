package handler

import (
	"fmt"
	"net/http"
	"strings"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/utils"
)

func (h *Handler) CreateCommentReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/reaction/create" {
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

	redirectTo := strings.TrimSpace(r.PostFormValue("redirect_to"))
	if !isValidRedirectTo(redirectTo) {
		redirectTo = "/"
	}

	commentIDStr := strings.TrimSpace(r.PostFormValue("comment_id"))
	commentID, err := utils.ParsePositiveIntID(commentIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reaction := r.PostFormValue("reaction")
	err = validateCreateCommentReactionForm(reaction)
	if err != nil {
		h.logger.Error("validate create comment reaction form:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createCommentReactionRequest := &models.CreateCommentReactionRequest{
		AuthorID:  h.getUserFromContext(r).ID,
		Reaction:  reaction,
		CommentID: commentID,
	}

	err = h.service.CommentReactionService.CreateCommentReaction(createCommentReactionRequest)
	if err != nil {
		h.logger.Error("create comment reaction:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectTo, http.StatusFound)
}

func validateCreateCommentReactionForm(reaction string) error {
	if reaction != "like" && reaction != "dislike" {
		return fmt.Errorf("reaction must either like or dislike")
	}

	return nil
}

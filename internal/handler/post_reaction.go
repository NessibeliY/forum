package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/utils"
)

func (h *Handler) CreatePostReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/reaction/create" {
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
		h.logger.Error("redirect_to:", redirectTo)
		redirectTo = "/"
	}

	postIDStr := strings.TrimSpace(r.PostFormValue("post_id"))
	postID, err := utils.ParsePositiveIntID(postIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reaction := strings.TrimSpace(r.PostFormValue("reaction"))
	err = validateCreatePostReactionForm(reaction)
	if err != nil {
		h.logger.Error("validate create post reaction form:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}

func validateCreatePostReactionForm(reaction string) error {
	if reaction != "like" && reaction != "dislike" {
		return fmt.Errorf("reaction must either like or dislike")
	}

	return nil
}

func isValidRedirectTo(redirectTo string) bool {
	if redirectTo == "/" {
		return true
	}

	matched, err := regexp.MatchString(`^\/post\?id=[1-9]\d*$`, redirectTo)
	return err == nil && matched
}

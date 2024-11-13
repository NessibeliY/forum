package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

func (h *Handler) CreatePostReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/reaction/create" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postIDStr := strings.TrimSpace(r.PostFormValue("post_id"))
	reaction := strings.TrimSpace(r.PostFormValue("reaction"))

	err = validateCreatePostReactionForm(postIDStr, reaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createPostReactionRequest := models.CreatePostReactionRequest{
		AuthorID: h.getUserFromContext(r).ID,
		PostID:   postID,
		Reaction: reaction,
	}

	fmt.Println(createPostReactionRequest)
	// err=h.service.PostReactionService.
}

func validateCreatePostReactionForm(postIDStr string, reaction string) error {
	if postIDStr == "" || reaction == "" {
		return fmt.Errorf("postID or reaction cannot be empty")
	}

	if reaction != "like" && reaction != "dislike" {
		return fmt.Errorf("reaction must either like or dislike")
	}

	return nil
}

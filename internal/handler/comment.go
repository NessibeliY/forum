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
		// h.logger.
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		// h.logger.
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		// h.logger
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := strings.TrimSpace(r.PostFormValue("content"))
	postIDStr := strings.TrimSpace(r.PostFormValue("post_id"))

	validationsErrMap := validateCreateCommentForm(content, postIDStr)
	if validationsErrMap != nil {
		// h.logger
		errorsJSON, err := json.Marshal(validationsErrMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cookies.SetCookie(w, errorsMapCookieName, string(errorsJSON), 300)

		http.Redirect(w, r, fmt.Sprintf("/post?id=%s", postIDStr), http.StatusSeeOther)
		return
	}

	postID, err := utils.ParsePositiveIntID(postIDStr)
	if err != nil {
		// h.logger
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

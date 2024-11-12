package handler

import (
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/create" {
		//h.logger.
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		//h.logger.
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		//h.logger
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := strings.TrimSpace(r.PostFormValue("content"))
	postIDStr := strings.TrimSpace(r.PostFormValue("post_id"))

	validationsErrMap := validateCreateCommentForm(content, postIDStr)
	if validationsErrMap != nil {
		//h.logger
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		//h.logger
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func validateCreateCommentForm(content string, postIDStr string) map[string]string {
	errors := make(map[string]string)

	if content == "" || postIDStr == "" {
		return errors.New("content or post_id is empty")
	}

	if len(content) > 10000 {
		return errors.New("content must be less than 10000 characters")
	}

	return nil
}

package handler

import (
	"net/http"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	posts, err := h.service.PostService.GetAllPosts()
	if err != nil {
		h.logger.Info("get categories home page", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		h.logger.Info("get categories home page", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Render(w, "index.page.html", H{
		"authenticated_user": h.getUserFromContext(r),
		"posts":              posts,
		"categories":         categories,
	})
}

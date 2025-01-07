package handler

import (
	"errors"
	"net/http"
	"strings"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/utils"
)

func (h *Handler) SendReport(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/report" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.sendReportMethodGet(w, r)
	case http.MethodPost:
		h.sendReportMethodPost(w, r)
	default:
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) sendReportMethodGet(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.logger.Error("get query for user_id")
		h.clientError(w, http.StatusBadRequest)
		return
	}

	userID, err := utils.ParsePositiveIntID(userIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		h.clientError(w, http.StatusNotFound)
		return
	}

	user, err := h.service.UserService.GetUserByID(userID)
	if err != nil {
		h.logger.Error("get user by id:", err.Error())
		h.serverError(w, err)
		return
	}
	if user == nil {
		h.logger.Info("user nil")
		h.clientError(w, http.StatusNotFound)
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		h.logger.Error("get query for post_id")
		h.clientError(w, http.StatusBadRequest)
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
		h.logger.Info("post nil")
		h.clientError(w, http.StatusNotFound)
		return
	}

	h.Render(w, "report.page.html", http.StatusOK, H{
		"reported_user":      user,
		"reported_post":      post,
		"authenticated_user": h.getUserFromContext(r),
	})
}

func (h *Handler) sendReportMethodPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		h.logger.Error("get query for post_id")
		h.clientError(w, http.StatusBadRequest)
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
		h.logger.Info("post nil")
		h.clientError(w, http.StatusNotFound)
		return
	}

	reason := strings.TrimSpace(r.PostFormValue("report_text"))
	err = validateReportText(reason)
	if err != nil {
		h.logger.Error("validate report text:", err)

		categories, err := h.service.CategoryService.GetAllCategories()
		if err != nil {
			h.logger.Error("get all categories:", err.Error())
			h.serverError(w, err)
			return
		}

		h.Render(w, "report.page.html", http.StatusBadRequest, H{
			"error":              err,
			"categories":         categories,
			"authenticated_user": h.getUserFromContext(r),
		})
		return
	}

	err = h.service.PostService.SendReport(&models.SendReportRequest{
		PostID:      postID,
		Reason:      reason,
		ModeratorID: h.getUserFromContext(r).ID,
	})
	if err != nil {
		h.logger.Error("get post by id:", err.Error())
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func validateReportText(reportText string) error {
	if reportText == "" {
		return errors.New("report text is empty")
	}

	if len(reportText) > 1000 {
		return errors.New("report text is too long")
	}

	return nil
}

func (h *Handler) SendModeratorRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/moderator-request" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := h.service.UserService.SendModeratorRequest(h.getUserFromContext(r).ID)
	if err != nil {
		h.logger.Error("send moderator request:", err)
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) ViewModeratorRequests(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/view/moderator-requests" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	requests, err := h.service.UserService.GetModeratorRequests()
	if err != nil {
		h.logger.Error("get moderator requests:", err)
		h.serverError(w, err)
		return
	}

	h.Render(w, "moderator_requests.page.html", http.StatusOK, H{
		"requests":           requests,
		"authenticated_user": h.getUserFromContext(r),
	})
}

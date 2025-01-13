package handler

import (
	"database/sql"
	"net/http"

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
	case http.MethodPost:
		h.sendReportMethodPost(w, r)
	default:
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) sendReportMethodPost(w http.ResponseWriter, r *http.Request) {
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
		h.logger.Error("get post by id")
		h.clientError(w, http.StatusBadRequest)
		return
	}

	err = h.service.PostService.SendReport(&models.SendReportRequest{
		PostID:      postID,
		ModeratorID: h.getUserFromContext(r).ID,
		Post:        post,
		Moderator:   user,
	})
	if err != nil {
		h.logger.Error("get post by id:", err.Error())
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

	var countNotification int
	if h.getUserFromContext(r) != nil {
		countNotification, err = h.service.NotificationService.GetCountNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get countNotification:", err)
			h.serverError(w, err)
			return
		}
	}

	h.Render(w, "moderator_requests.page.html", http.StatusOK, H{
		"requests":           requests,
		"authenticated_user": h.getUserFromContext(r),
		"count_notification": countNotification,
	})
}

func (h *Handler) SetNewRole(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/moderator-decision" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.logger.Error("user id is required")
		h.clientError(w, http.StatusBadRequest)
		return
	}
	userID, err := utils.ParsePositiveIntID(userIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		h.clientError(w, http.StatusNotFound)
		return
	}

	var decision bool
	switch r.URL.Query().Get("decision") {
	case "1":
		decision = true
	case "0":
		decision = false
	default:
		h.logger.Error("decision must be 0 or 1:", r.URL.Query().Get("decision"))
		h.clientError(w, http.StatusBadRequest)
		return
	}

	request := &models.UpdateRoleRequest{
		UserID:    userID,
		AdminID:   h.getUserFromContext(r).ID,
		Processed: decision,
	}
	err = h.service.UserService.SetNewRole(request)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Error("user not found", err)
			h.clientError(w, http.StatusBadRequest)
			return
		}
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) ReportModeration(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/reports/moderation" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.ReportModerationGet(w, r)
	case http.MethodPost:
		h.ReportModerationPost(w, r)
	default:
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) ReportModerationGet(w http.ResponseWriter, r *http.Request) {
	var countNotification int
	var err error
	if h.getUserFromContext(r) != nil {
		countNotification, err = h.service.NotificationService.GetCountNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get count notification:", err)
			h.serverError(w, err)
			return
		}
	}

	moderatedList, err := h.service.PostService.GetAllModeratedPosts()
	if err != nil {
		h.logger.Error("get moderated posts")
		h.clientError(w, http.StatusBadRequest)
		return
	}

	h.Render(w, "report_moderation.page.html", http.StatusOK, H{
		"authenticated_user": h.getUserFromContext(r),
		"count_notification": countNotification,
		"moderated_list":     moderatedList,
	})
}

func (h *Handler) ReportModerationPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		h.logger.Error("post id is required")
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
		h.logger.Info("get post:", err)
		h.serverError(w, err)
		return
	}

	if post == nil {
		h.logger.Error("post doesn't exist: post nil")
		h.clientError(w, http.StatusNotFound)
		return
	}

	decision := r.URL.Query().Get("decision")
	if decision == "" {
		h.logger.Error("decision is required")
		h.clientError(w, http.StatusNotFound)
		return
	}

	reportRequst := models.ModerationReport{
		IsModerated: true,
		AdminAnswer: decision,
		PostID:      postID,
	}

	err = h.service.PostService.UpdateModerationReport(&reportRequst)
	if err != nil {
		h.logger.Info("update moderation report:", err)
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/reports/moderation", http.StatusSeeOther)
}

func (h *Handler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/manage/users" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.ChangeUserRoleGet(w, r)
	case http.MethodPost:
		h.ChangeUserRolePost(w, r)
	default:
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) ChangeUserRoleGet(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.UserService.GetAllUsers()
	if err != nil {
		h.logger.Info("get all users:", err)
		h.serverError(w, err)
		return
	}

	var countNotification int
	if h.getUserFromContext(r) != nil {
		countNotification, err = h.service.NotificationService.GetCountNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get countNotification:", err)
			h.serverError(w, err)
			return
		}
	}

	h.Render(w, "users_page.page.html", http.StatusOK, H{
		"authenticated_user": h.getUserFromContext(r),
		"count_notification": countNotification,
		"users":              users,
	})
}

func (h *Handler) ChangeUserRolePost(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		h.logger.Error("user id is required")
		h.clientError(w, http.StatusBadRequest)
		return
	}
	userID, err := utils.ParsePositiveIntID(userIDStr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		h.clientError(w, http.StatusNotFound)
		return
	}

	role := r.URL.Query().Get("role")
	if role == "" {
		h.logger.Error("role is required")
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if role != models.UserRole && role != models.ModeratorRole {
		h.logger.Error("role is invalid: ", role)
		h.clientError(w, http.StatusBadRequest)
		return
	}

	err = h.service.UserService.ChangeRole(userID, role)
	if err != nil {
		h.logger.Info("change role:", err)
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/manage/users", http.StatusSeeOther)
}

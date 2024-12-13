package handler

import (
	"net/http"
	"strings"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/utils"
)

func (h *Handler) Notification(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/notifications" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var countNotification int
	var currentNotifications []models.Notification
	var archivedNotifications []models.Notification
	var err error
	if h.getUserFromContext(r) != nil {
		countNotification, err = h.service.NotificationService.GetCountNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get countNotification:", err)
			h.serverError(w, err)
			return
		}

		currentNotifications, err = h.service.NotificationService.GetCurrentNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get GetListNotifications:", err)
			h.serverError(w, err)
			return
		}

		archivedNotifications, err = h.service.NotificationService.GetArchivedNotifications(h.getUserFromContext(r).ID)
		if err != nil {
			h.logger.Info("get archived notifications:", err)
			h.serverError(w, err)
			return
		}
	}

	categories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		h.logger.Info("get all categories:", err)
		h.serverError(w, err)
		return
	}

	h.Render(w, "notification.page.html", http.StatusOK, H{
		"categories":             categories,
		"authenticated_user":     h.getUserFromContext(r),
		"count_notification":     countNotification,
		"current_notifications":  currentNotifications,
		"archived_notifications": archivedNotifications,
	})
}

func (h *Handler) MakeNotificationIsRead(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/notifications/read" {
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

	notificationIDstr := strings.TrimSpace(r.PostFormValue("notification_id"))
	notificationID, err := utils.ParsePositiveIntID(notificationIDstr)
	if err != nil {
		h.logger.Error("parse positive int:", err.Error())
		h.clientError(w, http.StatusBadRequest)
		return
	}

	notification, err := h.service.NotificationService.GetNotificationByID(notificationID)
	if err != nil {
		h.logger.Info("get Notification:", err)
		h.serverError(w, err)
		return
	}

	if notification == nil {
		h.logger.Error("get notification by id", err.Error())
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if h.getUserFromContext(r) != nil {
		err = h.service.NotificationService.MakeNotificationIsRead(h.getUserFromContext(r).ID, notificationID)
		if err != nil {
			h.logger.Errorf("make notifications is read: %s", err.Error())
			h.clientError(w, http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/notifications", http.StatusSeeOther)
}

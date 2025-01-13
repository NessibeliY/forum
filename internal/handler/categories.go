package handler

import (
	"fmt"
	"net/http"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/utils"
)

func (h *Handler) ManageCategories(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/manage/categories" {
		h.logger.Error("url path:", r.URL.Path)
		h.clientError(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.ManageCategoriesGet(w, r)
	case http.MethodPost:
		h.ManageCategoriesPost(w, r)
	default:
		h.logger.Errorf("method not allowed: %s", r.Method)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) ManageCategoriesGet(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.CategoryService.GetAllCategories()
	if err != nil {
		h.logger.Info("get all categories:", err)
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

	h.Render(w, "manage_categories.page.html", http.StatusOK, H{
		"authenticated_user": h.getUserFromContext(r),
		"count_notification": countNotification,
		"categories":         categories,
	})
}

func (h *Handler) ManageCategoriesPost(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.URL.Query().Get("category_id")

	// delete category
	if categoryIDStr != "" {
		categoryID, err := utils.ParsePositiveIntID(categoryIDStr)
		if err != nil {
			h.logger.Error("parse positive int:", err.Error())
			h.clientError(w, http.StatusNotFound)
			return
		}

		category, err := h.service.CategoryService.GetCategoryByID(categoryID)
		if err != nil {
			h.logger.Error("get category by id (check id category):", err.Error())
			h.serverError(w, err)
			return
		}

		if category == nil {
			var countNotification int
			if h.getUserFromContext(r) != nil {
				countNotification, err = h.service.NotificationService.GetCountNotifications(h.getUserFromContext(r).ID)
				if err != nil {
					h.logger.Info("get countNotification:", err)
					h.serverError(w, err)
					return
				}
			}

			h.Render(w, "manage_categories.page.html", http.StatusBadRequest, H{
				"error":              "Category with the given ID does not exist.",
				"count_notification": countNotification,
				"authenticated_user": h.getUserFromContext(r),
			})
			return
		}

		err = h.service.CategoryService.DeleteCategory(categoryID)
		if err != nil {
			h.logger.Error("delete category:", err.Error())
			h.serverError(w, err)
			return
		}

		http.Redirect(w, r, "/manage/categories", http.StatusSeeOther)
		return

	}

	// add new category

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parse form create category:", err.Error())
		h.clientError(w, http.StatusBadRequest)
		return
	}

	categoryName := r.FormValue("category_name")
	if categoryName == "" {
		var countNotification int
		if h.getUserFromContext(r) != nil {
			countNotification, err = h.service.NotificationService.GetCountNotifications(h.getUserFromContext(r).ID)
			if err != nil {
				h.logger.Info("get countNotification:", err)
				h.serverError(w, err)
				return
			}
		}

		h.Render(w, "manage_categories.page.html", http.StatusBadRequest, H{
			"error":              "category name cannot be empty",
			"count_notification": countNotification,
			"authenticated_user": h.getUserFromContext(r),
		})
		return
	}

	categoryRequest := &models.Category{
		Name: categoryName,
	}
	fmt.Println(categoryRequest)

	_, err = h.service.CategoryService.CreateCategory(categoryRequest)
	if err != nil {
		h.logger.Error("create category:", err.Error())
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/manage/categories", http.StatusSeeOther)
}

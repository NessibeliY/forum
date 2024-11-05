package handler

import (
	"net/http"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/validator"
)

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/signup" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.signupGet(w, r)
	case http.MethodPost:
		h.signupPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) signupGet(w http.ResponseWriter, r *http.Request) {
	h.Render(w, "signup.page.html", &PageData{
		AuthenticatedUser: h.getUserFromContext(r),
	})
}

func (h *Handler) getUserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		//h.logger.Info("user is not authenticated")
		return nil
	}
	return user
}

func (h *Handler) signupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		//h.logger.Errorf("parse form: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validator := validator.New(r.PostForm)
	validator.Required("username", "email", "password")
	validator.MinLength("username", 3)
	validator.MaxLength("username", 30)
	validator.MaxLength("email", 30)
	validator.ValidEmail("email")
	validator.MinLength("password", 8)
	validator.MaxLength("password", 30)

	if !validator.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		h.Render(w, "signup.page.html", &PageData{
			Validator: validator,
		})
		return
	}

	signupRequest := &models.SignupRequest{
		Username: r.PostFormValue("username"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	err = h.service.UserService.SignupUser(signupRequest)
	if err != nil {
		switch err {
		case models.ErrDuplicateEmail:
			validator.Errors.Add("email", "email already in use")
			w.WriteHeader(http.StatusBadRequest)
			h.Render(w, "signup.page.html", &PageData{
				Validator: validator,
			})
			return
		case models.ErrDuplicateUsername:
			validator.Errors.Add("username", "username already in use")
			w.WriteHeader(http.StatusBadRequest)
			h.Render(w, "signup.page.html", &PageData{
				Validator: validator,
			})
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/login" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.loginGet(w, r)
	case http.MethodPost:
		h.loginPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

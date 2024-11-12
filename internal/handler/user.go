package handler

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/cookies"
)

var emailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+\\/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+\\/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")

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
	h.Render(w, "signup.page.html", H{
		"authenticated_user": h.getUserFromContext(r),
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
		//h.logger.Errorf("signup post %w", err)
		http.Error(w, "unable to parse form", http.StatusInternalServerError)
		return
	}

	username := strings.TrimSpace(r.PostFormValue("username"))
	email := strings.TrimSpace(r.PostFormValue("email"))
	password := r.PostFormValue("password")

	err = validateSignupForm(username, email, password)
	if err != nil {
		//h.logger
		h.Render(w, "signup.page.html", H{
			"error": err.Error(),
		})
		return
	}

	signupRequest := &models.SignupRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	err = h.service.UserService.SignupUser(signupRequest)
	if err != nil {
		var errorMsg string
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			errorMsg = models.ErrDuplicateEmail.Error()
		case errors.Is(err, models.ErrDuplicateUsername):
			errorMsg = models.ErrDuplicateUsername.Error()
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		h.Render(w, "signup.page.html", H{
			"error": errorMsg,
		})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func validateSignupForm(username string, email string, password string) error {
	if username == "" || email == "" || password == "" {
		return fmt.Errorf("username or email or password is empty")
	}

	if len(username) < 3 || len(username) > 50 {
		return fmt.Errorf("username length must be between 3 and 50 characters")
	}

	if len(email) > 320 || !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format or length exceeds 320 characters")
	}

	if len(password) < 8 || len(password) > 50 {
		return fmt.Errorf("password length must be between 8 and 50 characters")
	}

	return nil
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

func (h *Handler) loginGet(w http.ResponseWriter, r *http.Request) {
	h.Render(w, "login.page.html", H{
		"authenticated_user": h.getUserFromContext(r),
	})
}

func (h *Handler) loginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		//h.logger.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := strings.TrimSpace(r.PostFormValue("email"))
	password := r.PostFormValue("password")

	err = validateLoginForm(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.Render(w, "login.page.html", H{
			"error": err.Error(),
		})
		return
	}

	loginPostRequest := &models.LoginRequest{
		Email:    email,
		Password: password,
	}
	userID, err := h.service.UserService.LoginUser(loginPostRequest)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			//h.logger.
			h.Render(w, "login.page.html", H{
				"Error": err.Error(),
			})
			return
		}
		//h.logger
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	session, err := h.service.SessionService.SetSession(userID)
	if err != nil {
		//h.logger.Errorf("create session: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	cookies.SetCookie(w, session.UUID, int(time.Until(session.ExpiresAt).Seconds()))

	http.Redirect(w, r, "/", http.StatusFound)
}

func validateLoginForm(email string, password string) error {
	if email == "" || password == "" {
		return fmt.Errorf("email or password is empty")
	}

	if len(email) > 320 || !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format or length exceeds 320 characters")
	}

	if len(password) < 8 || len(password) > 50 {
		return fmt.Errorf("password length must be between 8 and 50 characters")
	}

	return nil
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/logout" {
		//h.logger
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		//h.logger
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	cookie, err := cookies.GetCookie(r)
	if err != nil {
		//h.logger
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = h.service.SessionService.DeleteSession(cookie.Value)
	if err != nil {
		//h.logger
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	cookies.DeleteCookie(w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

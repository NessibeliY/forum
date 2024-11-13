package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/cookies"
)

const (
	sessionCookieName = "forum_session_cookie"
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
	h.Render(w, "sign_up.page.html", H{
		"authenticated_user": h.getUserFromContext(r),
	})
}

func (h *Handler) getUserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		// h.logger.Info("user is not authenticated")
		return nil
	}
	return user
}

func (h *Handler) signupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// h.logger.Errorf("signup post %w", err)
		http.Error(w, "unable to parse form", http.StatusInternalServerError)
		return
	}

	username := strings.TrimSpace(r.PostFormValue("username"))
	email := strings.TrimSpace(r.PostFormValue("email"))
	password := r.PostFormValue("password")

	errorsMap := validateSignupForm(username, email, password)

	if len(errorsMap) > 0 {
		// h.logger
		h.Render(w, "sign_up.page.html", H{
			"Username":      username,
			"Email":         email,
			"Password":      password,
			"ErrorMessages": errorsMap,
		})
		return
	}

	signupRequest := &models.SignupRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	err = h.service.UserService.SignupUser(signupRequest)
	var errorMsg string
	if err != nil {
		switch err {
		case models.ErrDuplicateEmail:
			errorMsg = "Email already in use"
			h.Render(w, "sign_up.page.html", H{
				"Username": username,
				"Email":    email,
				"Password": password,
				"error":    errorMsg,
			})
			return
		case models.ErrDuplicateUsername:
			errorMsg = "UserName already in use"
			h.Render(w, "sign_up.page.html", H{
				"Username": username,
				"Email":    email,
				"Password": password,
				"error":    errorMsg,
			})
			return
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func validateSignupForm(username string, email string, password string) map[string]string {
	errors := make(map[string]string)

	if username == "" {
		errors["username"] = "username cannot be empty"
	}

	if email == "" {
		errors["email"] = "email cannot be empty"
	}

	if password == "" {
		errors["password"] = "password cannot be empty"
	}

	if len(username) < 3 || len(username) > 50 {
		errors["username"] = "username length must be between 3 and 50 characters"
	}

	if len(email) > 320 || !emailRegex.MatchString(email) {
		errors["email"] = "invalid email format or length exceeds 320 characters"
	}

	if len(password) < 8 || len(password) > 50 {
		errors["password"] = "password length must be between 8 and 50 characters"
	}

	return errors
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
	h.Render(w, "sign_in.page.html", H{
		"authenticated_user": h.getUserFromContext(r),
	})
}

func (h *Handler) loginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("parse form", err)
		// h.logger.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := strings.TrimSpace(r.PostFormValue("email"))
	password := r.PostFormValue("password")

	errMap := validateLoginForm(email, password)
	if len(errMap) > 0 {
		h.Render(w, "sign_in.page.html", H{
			"Email":         email,
			"Password":      password,
			"ErrorMessages": errMap,
		})
		return
	}

	loginPostRequest := &models.LoginRequest{
		Email:    email,
		Password: password,
	}

	userID, err := h.service.UserService.LoginUser(loginPostRequest)
	fmt.Println("userID", userID)
	fmt.Println("err", err)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			h.Render(w, "sign_in.page.html", H{
				"Email":    email,
				"Password": password,
				"Error":    err.Error(),
			})
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Println("here")

	session, err := h.service.SessionService.SetSession(userID)
	fmt.Println("session", session)
	fmt.Println("err", err)
	if err != nil {
		fmt.Println("SetSession", err)
		// h.logger.Errorf("create session: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	cookies.SetCookie(w, sessionCookieName, session.UUID, int(time.Until(session.ExpiresAt).Seconds()))

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func validateLoginForm(email string, password string) map[string]string {
	errors := make(map[string]string)

	if email == "" {
		errors["email"] = "email  is empty"
	}

	if password == "" {
		errors["password"] = "password  is empty"
	}

	if len(email) > 320 || !emailRegex.MatchString(email) {
		errors["email"] = "invalid email format or length exceeds 320 characters"
	}

	if len(password) < 8 || len(password) > 50 {
		errors["password"] = "password length must be between 8 and 50 characters"
	}

	return errors
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/logout" {
		// h.logger
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		// h.logger
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	cookie, err := cookies.GetCookie(r, sessionCookieName)
	if err != nil {
		// h.logger
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = h.service.SessionService.DeleteSession(cookie.Value)
	if err != nil {
		// h.logger
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	cookies.DeleteCookie(w, sessionCookieName)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

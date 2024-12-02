package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/cookies"
)

var (
	googleOAuthEndpoint    = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenEndpoint    = "https://accounts.google.com/o/oauth2/token"
	googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type tokenResp struct {
	AccessToken string `json:"access_token"`
}

type googleUserInfo struct {
	email string `json:"email"`
	name  string `json:"name"`
	sub   string `json:"sub"`
}

func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=profile email", googleOAuthEndpoint, h.googleConfig.ClientID, h.googleConfig.RedirectURI)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.Error("missing code in callback request")
		h.clientError(w, http.StatusBadRequest)
		return
	}

	data := strings.NewReader(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, h.googleConfig.ClientID, h.googleConfig.ClientSecret, h.googleConfig.RedirectURI))
	resp, err := http.Post(googleTokenEndpoint, "application/x-www-form-urlencoded", data)
	if err != nil {
		h.logger.Error("http post:", err)
		h.serverError(w, err)
		return
	}
	defer resp.Body.Close()

	tokenResp := &tokenResp{}

	err = json.NewDecoder(resp.Body).Decode(tokenResp)
	if err != nil {
		h.logger.Error("decode token response:", err)
		h.serverError(w, err)
		return
	}

	req, err := http.NewRequest(http.MethodGet, googleUserInfoEndpoint, nil)
	if err != nil {
		h.logger.Error("new request:", err)
		h.serverError(w, err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		h.logger.Error("do:", err)
		h.serverError(w, err)
		return
	}
	defer resp.Body.Close()

	googleUserInfo := googleUserInfo{}
	err = json.NewDecoder(resp.Body).Decode(&googleUserInfo)
	if err != nil {
		h.logger.Error("decode google user info:", err)
		h.serverError(w, err)
		return
	}

	user, err := h.service.UserService.GetUserByEmail(googleUserInfo.email)
	if err != nil {
		h.logger.Error("get user by email:", err)
		h.serverError(w, err)
		return
	}
	if user == nil {
		signupRequest := &models.SignupRequest{
			Username: googleUserInfo.name,
			Email:    googleUserInfo.email,
			Password: googleUserInfo.sub,
		}

		err = h.service.UserService.SignupUser(signupRequest)
		if err != nil {
			h.logger.Error("signup user:", err)
			h.serverError(w, err)
			return
		}
	}

	loginRequest := &models.LoginRequest{
		Email:    googleUserInfo.email,
		Password: googleUserInfo.sub,
	}

	userID, err := h.service.UserService.LoginUser(loginRequest)
	if err != nil {
		h.logger.Error("login user:", err.Error())
		h.serverError(w, err)
		return
	}

	session, err := h.service.SessionService.SetSession(userID)
	if err != nil {
		h.logger.Error("set session:", err.Error())
		h.serverError(w, err)
		return
	}
	cookies.SetCookie(w, sessionCookieName, session.UUID, int(time.Until(session.ExpiresAt).Seconds()))

	http.Redirect(w, r, "/", http.StatusFound)
}

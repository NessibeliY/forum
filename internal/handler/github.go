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

var githubOAuthEndpoint = "https://github.com/login/oauth/authorize"
var githubTokenEndpoint = "https://github.com/login/oauth/access_token"
var githubUserInfoEndpoint = "https://api.github.com/user"

type githubUserInfo struct {
	name   string `json:"name"`
	nodeID string `json:"node_id"`
}

func (h *Handler) GithubLogin(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=user:email", githubOAuthEndpoint, h.githubConfig.ClientID, h.githubConfig.RedirectURI)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.Error("missing code in callback request")
		h.clientError(w, http.StatusBadRequest)
		return
	}

	data := strings.NewReader(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s", code, h.githubConfig.ClientID, h.githubConfig.ClientSecret, h.githubConfig.RedirectURI))
	resp, err := http.Post(githubTokenEndpoint, "application/x-www-form-urlencoded", data)
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

	req, err := http.NewRequest(http.MethodGet, githubUserInfoEndpoint, nil)
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

	githubUserInfo := githubUserInfo{}
	err = json.NewDecoder(resp.Body).Decode(&githubUserInfo)
	if err != nil {
		h.logger.Error("decode google user info:", err)
		h.serverError(w, err)
		return
	}

	user, err := h.service.UserService.GetUserByEmail(githubUserInfo.nodeID)
	if err != nil {
		h.logger.Error("get user by email:", err)
		h.serverError(w, err)
		return
	}
	if user == nil {
		signupRequest := &models.SignupRequest{
			Username: githubUserInfo.name,
			Email:    githubUserInfo.nodeID,
			Password: githubUserInfo.nodeID,
		}

		err = h.service.UserService.SignupUser(signupRequest)
		if err != nil {
			h.logger.Error("signup user:", err)
			h.serverError(w, err)
			return
		}
	}

	loginRequest := &models.LoginRequest{
		Email:    githubUserInfo.nodeID,
		Password: githubUserInfo.nodeID,
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

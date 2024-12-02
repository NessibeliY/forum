package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/pkg/cookies"
)

var (
	githubOAuthEndpoint    = "https://github.com/login/oauth/authorize"
	githubTokenEndpoint    = "https://github.com/login/oauth/access_token"
	githubUserInfoEndpoint = "https://api.github.com/user"
)

type githubUserInfo struct {
	Name   string `json:"login"`
	NodeID string `json:"node_id"`
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

	body, err := io.ReadAll(resp.Body)

	params, err := url.ParseQuery(string(body))
	if err != nil {
		h.logger.Error("parse query:", err)
		h.serverError(w, err)
		return
	}
	accessToken := params.Get("access_token")

	req, err := http.NewRequest(http.MethodGet, githubUserInfoEndpoint, nil)
	if err != nil {
		h.logger.Error("new request:", err)
		h.serverError(w, err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		h.logger.Error("do:", err)
		h.serverError(w, err)
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)

	githubUserInfo := githubUserInfo{}
	err = json.Unmarshal(body, &githubUserInfo)
	if err != nil {
		h.logger.Error("unmarshal github user info:", err)
		h.serverError(w, err)
		return
	}

	user, err := h.service.UserService.GetUserByEmail(githubUserInfo.NodeID)
	if err != nil {
		h.logger.Error("get user by email:", err)
		h.serverError(w, err)
		return
	}
	if user == nil {
		signupRequest := &models.SignupRequest{
			Username: githubUserInfo.Name,
			Email:    githubUserInfo.NodeID,
			Password: githubUserInfo.NodeID,
		}

		err = h.service.UserService.SignupUser(signupRequest)
		if err != nil {
			h.logger.Error("signup user:", err)
			h.serverError(w, err)
			return
		}
	}

	loginRequest := &models.LoginRequest{
		Email:    githubUserInfo.NodeID,
		Password: githubUserInfo.NodeID,
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

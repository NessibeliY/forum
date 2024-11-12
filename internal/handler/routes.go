package handler

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/user/signup", h.Signup)
	mux.HandleFunc("/login", h.Login)
	mux.Handle("/user/logout", h.RequireAuthentication(http.HandlerFunc(h.Logout)))

	mux.Handle("/post/create", h.RequireAuthentication(http.HandlerFunc(h.CreatePost)))

	return mux
}

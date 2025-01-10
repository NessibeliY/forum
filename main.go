package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/handler"
	"01.alem.school/git/nyeltay/forum/internal/repository"
	"01.alem.school/git/nyeltay/forum/internal/service"
	"01.alem.school/git/nyeltay/forum/internal/template_cache"
	"01.alem.school/git/nyeltay/forum/pkg/db"
	"01.alem.school/git/nyeltay/forum/pkg/logger"

	"01.alem.school/git/nyeltay/forum/conf"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configPath := filepath.Join(cwd, "config.json")

	config, err := conf.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	l, err := logger.Setup(config)
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.Open(config.DSN)
	if err != nil {
		l.Fatal(err)
	}
	l.Info("connected to db")

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	templateCache, err := template_cache.NewTemplateCache()
	if err != nil {
		l.Fatal(err)
	}

	handler := handler.NewHandler(service, templateCache, l, config.GoogleConfig, config.GithubConfig)

	l.Infof("server is running on https://localhost%s", config.Port)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.Handle("/favicon.ico", http.NotFoundHandler())

	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/signup", handler.Signup)
	mux.HandleFunc("/login", handler.Login)
	mux.Handle("POST /logout", handler.RequireAuthentication(http.HandlerFunc(handler.Logout)))

	mux.HandleFunc("/google/callback", handler.GoogleCallback)
	mux.HandleFunc("/login/google/callback", handler.GoogleLogin)
	mux.HandleFunc("/github/callback", handler.GithubCallback)
	mux.HandleFunc("/login/github/callback", handler.GithubLogin)

	mux.Handle("/post/create", handler.RequireAuthentication(http.HandlerFunc(handler.CreatePost)))
	mux.HandleFunc("/post", handler.ShowPost)
	mux.Handle("/post/delete", handler.RequireAuthentication(http.HandlerFunc(handler.DeletePost)))

	mux.Handle("/comment/create", handler.RequireAuthentication(http.HandlerFunc(handler.CreateComment)))
	mux.Handle("/comment/delete", handler.RequireAuthentication(http.HandlerFunc(handler.DeleteComment)))

	mux.Handle("/post/reaction/create", handler.RequireAuthentication(http.HandlerFunc(handler.CreatePostReaction)))
	mux.Handle("/comment/reaction/create", handler.RequireAuthentication(http.HandlerFunc(handler.CreateCommentReaction)))

	mux.Handle("/myposts", handler.RequireAuthentication(http.HandlerFunc(handler.ShowMyPosts)))
	mux.Handle("/likedposts", handler.RequireAuthentication(http.HandlerFunc(handler.ShowLikedPosts)))
	mux.HandleFunc("/showposts", handler.ShowPostsByCategory)

	mux.Handle("/notifications", handler.RequireAuthentication(http.HandlerFunc(handler.Notification)))
	mux.Handle("/notifications/read", handler.RequireAuthentication(http.HandlerFunc(handler.MakeNotificationIsRead)))

	mux.Handle("/activity-page", handler.RequireAuthentication(http.HandlerFunc(handler.ActivityPage)))
	mux.Handle("/post/update", handler.RequireAuthentication(http.HandlerFunc(handler.UpdatePage)))
	mux.Handle("/comment/update", handler.RequireAuthentication(http.HandlerFunc(handler.UpdateComment)))

	mux.Handle("/report", handler.RequireAuthentication(handler.IsModerator(http.HandlerFunc(handler.SendReport))))
	mux.Handle("/moderator-request", handler.RequireAuthentication(handler.IsUser(http.HandlerFunc(handler.SendModeratorRequest))))
	mux.Handle("/view/moderator-requests", handler.RequireAuthentication(handler.IsAdmin(http.HandlerFunc(handler.ViewModeratorRequests))))
	mux.Handle("/moderator-decision", handler.RequireAuthentication(handler.IsAdmin(http.HandlerFunc(handler.SetNewRole))))
	mux.Handle("/reports/moderation", handler.RequireAuthentication(handler.IsAdmin(http.HandlerFunc(handler.ReportModerationGet))))

	rateLimiter := handler.NewRateLimiter(10, 50, 1*time.Minute)
	finalHandler := rateLimiter.Limit(handler.SecureHeaders(
		handler.RecoverPanic(
			handler.LogRequest(
				handler.Authenticate(mux),
			),
		),
	))

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.CurveP384,
		},
	}

	server := &http.Server{
		Addr:         config.Port,
		Handler:      finalHandler,
		TLSConfig:    tlsConfig,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"); err != nil && err != http.ErrServerClosed {
			l.Fatal(err)
		}
	}()

	<-stop
	l.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		l.Fatalf("server forced to shutdown: %v", err)
	}

	l.Info("server stopped gracefully")
}

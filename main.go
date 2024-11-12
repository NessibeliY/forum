package main

import (
	"log"
	"net/http"
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
	config, err := conf.Load("config.json")
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
	handler := handler.NewHandler(service, templateCache, l)

	l.Infof("server is running on localhost%s", config.Port)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", handler.Home)
	mux.HandleFunc("/user/signup", handler.Signup)
	mux.HandleFunc("/login", handler.Login)
	mux.Handle("/user/logout", handler.RequireAuthentication(http.HandlerFunc(handler.Logout)))

	mux.Handle("/post/create", handler.RequireAuthentication(http.HandlerFunc(handler.CreatePost)))
	mux.HandleFunc("/post/", handler.ShowPost)

	mux.Handle("comment/create", handler.RequireAuthentication(http.HandlerFunc(handler.CreateComment)))

	finalHandler := handler.SecureHeaders(
		handler.RecoverPanic(
			handler.LogRequest(
				handler.Authenticate(mux),
			),
		),
	)

	router := &http.Server{
		Addr:         config.Port,
		Handler:      finalHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		err := router.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			l.Fatal(err)
		} else {
			l.Info("server stopped")
		}
	}()

	select {}
}

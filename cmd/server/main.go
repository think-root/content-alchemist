package main

import (
	"content-alchemist/config"
	"content-alchemist/server/middlewares"
	"content-alchemist/server/routers"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.RateLimit)
	r.Use(middlewares.CheckToken)

	r.Post("/think-root/api/manual-generate/", routers.ManualGenerate)
	r.Post("/think-root/api/auto-generate/", routers.AutoGenerate)
	r.Post("/think-root/api/get-repository/", routers.GetRepository)
	r.Patch("/think-root/api/update-posted/", routers.UpdatePostedStatus)

	log.Printf("Server listen on port %s (app version: %s)\n\n",
		config.SERVER_PORT, config.APP_VERSION)

	http.ListenAndServe(":"+config.SERVER_PORT, r)
}

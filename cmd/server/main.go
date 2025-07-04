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
	r.Use(middlewares.CORSMiddleware)
	r.Use(middlewares.RateLimit)
	r.Use(middlewares.CheckToken)

	r.Post("/think-root/api/manual-generate/", routers.ManualGenerate)
	r.Post("/think-root/api/auto-generate/", routers.AutoGenerate)
	r.Post("/think-root/api/get-repository/", routers.GetRepository)
	r.Patch("/think-root/api/update-posted/", routers.UpdatePostedStatus)
	r.Patch("/think-root/api/update-repository-text/", routers.UpdateRepositoryText)
	r.Delete("/think-root/api/delete-repository/", routers.DeleteRepository)

	log.Printf("Server listen on port %s (app version: %s)\n\n",
		config.SERVER_PORT, config.APP_VERSION)

	err := http.ListenAndServe(":"+config.SERVER_PORT, r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

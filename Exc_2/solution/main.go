package main

import (
	"embed"
	"io/fs"
	"log"
	"log/slog"
	"net/http"

	"ordersystem/repository"
	"ordersystem/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "ordersystem/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

//go:embed frontend/*
var embeddedFrontend embed.FS

// @title       Order System
// @description This system enables drink orders and should not be used for the forbidden Hungover Games.
func main() {
	db := repository.NewDatabaseHandler()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Frontend index
	staticFS, err := fs.Sub(embeddedFrontend, "frontend")
	if err != nil {
		log.Fatal(err)
	}
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// If ServeFileFS isn't available in your Go version, switch to:
		// http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
		http.ServeFileFS(w, r, staticFS, "index.html")
	})

	// API
	r.Get("/api/menu", rest.GetMenu(db))
	r.Get("/api/order/all", rest.GetOrders(db))
	r.Get("/api/order/totalled", rest.GetOrdersTotalled(db))
	r.Post("/api/order", rest.PostOrder(db))

	// OpenAPI
	r.Get("/openapi/*", httpSwagger.WrapHandler)

	slog.Info("⚡ Order System is up and running")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}

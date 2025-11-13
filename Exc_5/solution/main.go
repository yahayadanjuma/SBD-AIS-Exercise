package main

import (
	"log"
	"strings"
	"log/slog"
	"net/http"
	"ordersystem/repository"
	"ordersystem/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	_ "ordersystem/docs"

	// OpenApi
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title				Order System
// @description			This system enables drink orders and should not be used for the forbidden Hungover Games.
func main() {
	db, err := repository.NewDatabaseHandler()
	if err != nil {
		log.Fatalln(err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// allow local cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Origin", "cache-control", "expires", "pragma"},
		ExposedHeaders:   []string{"Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Menu Routes
	// Wrap the API menu handler so that when a browser requests it (Accept: text/html)
	// we redirect to the UI at /menu. API clients still receive JSON.
	r.Get("/api/menu", func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		if strings.Contains(accept, "text/html") {
			http.Redirect(w, r, "/menu", http.StatusFound)
			return
		}
		rest.GetMenu(db)(w, r)
	})
	// Order Routes
	r.Get("/api/order/all", rest.GetOrders(db))
	r.Get("/api/order/totalled", rest.GetOrdersTotal(db))
	r.Post("/api/order", rest.PostOrder(db))
	// OpenAPI Routes
	r.Get("/openapi/*", httpSwagger.WrapHandler)

	// (No root or /menu redirects here — the static web server serves the UI.)

	slog.Info("⚡⚡⚡ Order System is up and running ⚡⚡⚡")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}

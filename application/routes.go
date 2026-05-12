package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/anand-aman/kubepilot/handler"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"kubepilot api running"}`))
	})

	router.Route("/namespaces", a.loadNamespacesRoutes)

	a.router = router
}

func (a *App) loadNamespacesRoutes(router chi.Router) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if a.client == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error":"kubernetes client not available"}`))
			return
		}
		h := handler.NewHandler(a.client)
		h.GetNamespaces(w, r)
	})
}

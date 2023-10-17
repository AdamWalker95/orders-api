package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/AdamWalker95/orders-api/handler"
	"github.com/AdamWalker95/orders-api/repository/client"
	"github.com/AdamWalker95/orders-api/repository/order"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/user", a.loadUserRoutes)
	router.Route("/orders", a.loadOrderRoutes)

	a.router = router
}

func (a *App) loadUserRoutes(router chi.Router) {
	userHandler := &handler.Client{
		Repo: &client.SqlRepo{
			Client:       a.db,
			DatabaseName: a.config.MySqlDatabaseName,
		},
	}

	router.Post("/", userHandler.Create)
	router.Get("/{id}", userHandler.GetByID)
	router.Put("/{id}", userHandler.UpdateByID)
	router.Delete("/{id}", userHandler.DeleteByID)
}

func (a *App) loadOrderRoutes(router chi.Router) {
	orderHandler := &handler.Order{
		Repo: &order.RedisRepo{
			Client: a.rdb,
		},
	}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetByID)
	router.Put("/{id}", orderHandler.UpdateByID)
	router.Delete("/{id}", orderHandler.DeleteByID)
}

package main

import (
	"github.com/go-chi/chi/v5"
	cMiddleware "github.com/go-chi/chi/v5/middleware"
	"gofemart/cmd/gophermart/handlers/login"
	"gofemart/cmd/gophermart/handlers/orders"
	"gofemart/internal/logger"
	"gofemart/internal/middlewares"
	"gofemart/internal/token"
)

// getRouter конфигурация роутинга приложение
func getRouter() chi.Router {
	router := chi.NewRouter()
	// Устанавливаем мидлваре
	router.Use(
		middlewares.JSONHeaders,
		cMiddleware.StripSlashes, // Убираем лишние слеши
		logger.LogRequests,       // Логируем данные запроса
	)

	router.Post("/api/user/register", login.RegistrationHandler)
	router.Post("/api/user/login", login.LoginHandler)

	router.Group(func(r chi.Router) {
		r.Use(token.AuthMiddleware)
		r.Post("/api/user/orders", orders.RegisterOrderHandler)
	})

	return router
}

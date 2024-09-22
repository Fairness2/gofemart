package router

import (
	"github.com/go-chi/chi/v5"
	cMiddleware "github.com/go-chi/chi/v5/middleware"
	//httpSwagger "github.com/swaggo/http-swagger"
	"gofemart/cmd/gophermart/handlers/balance"
	"gofemart/cmd/gophermart/handlers/login"
	"gofemart/cmd/gophermart/handlers/orders"
	//_ "gofemart/docs"
	config "gofemart/internal/configuration"
	database "gofemart/internal/databse"
	"gofemart/internal/logger"
	"gofemart/internal/middlewares"
	"gofemart/internal/token"
)

// NewRouter конфигурация роутинга приложение
func NewRouter(dbPool *database.DBPool, cnf *config.CliConfig) chi.Router {
	lHandlers := login.NewHandlers(dbPool.DBx, cnf.JWTKeys, cnf.TokenExpiration, cnf.HashKey)
	bHandlers := balance.NewHandlers(dbPool.DBx)
	oHandlers := orders.NewHandlers(dbPool.DBx)
	authenticator := token.NewAuthenticator(dbPool.DBx, cnf.JWTKeys, cnf.TokenExpiration)
	router := chi.NewRouter()
	// Устанавливаем мидлваре
	router.Use(
		middlewares.JSONHeaders,
		cMiddleware.StripSlashes, // Убираем лишние слеши
		logger.LogRequests,       // Логируем данные запроса
	)
	// Адрес свагера
	/*router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
	))*/
	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", lHandlers.RegistrationHandler)
		r.Post("/login", lHandlers.LoginHandler)
		r.Group(registerRoutesWithAuth(bHandlers, oHandlers, authenticator))
	})

	return router
}

// registerRoutesWithAuth маршруты с аутентификацией
func registerRoutesWithAuth(bHandlers *balance.Handlers, oHandlers *orders.Handlers, authenticator *token.Authenticator) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(
			authenticator.Middleware,
			cMiddleware.Compress(5, "gzip", "deflate"),
		)
		r.Post("/orders", oHandlers.RegisterOrderHandler)
		r.Post("/balance/withdraw", bHandlers.WithdrawHandler)
		r.Get("/balance", bHandlers.GetBalanceHandler)
		r.Group(registerRoutesWithCompressed(oHandlers))
	}
}

// registerRoutesWithCompressed настраивает маршруты для заказов с применением промежуточного программного обеспечения сжатия.
// Так как ответы данных путей могут быть большими
func registerRoutesWithCompressed(oHandlers *orders.Handlers) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(
			cMiddleware.Compress(5, "gzip", "deflate"),
		)
		r.Get("/orders", oHandlers.GetOrdersHandler)
		r.Get("/withdrawals", oHandlers.GetOrdersWwithdrawalsHandler)
	}
}

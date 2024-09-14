package token

import (
	"context"
	config "gofemart/internal/configuration"
	database "gofemart/internal/databse"
	"gofemart/internal/helpers"
	"gofemart/internal/logger"
	"gofemart/internal/repositories"
	"net/http"
	"strconv"
	"strings"
)

// Key тип ключей в контексте
type Key string

// UserKey ключ авторизованного пользователя в контексте
var UserKey Key = "user"

// AuthMiddleware авторизовываем пользователя по токену и записываем его в контекст
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tknString := r.Header.Get("Authorization")
		if tknString == "" {
			helpers.ProcessErrorWithStatus("Authorization header is not exists", http.StatusUnauthorized, w)
			return
		}
		if !strings.HasPrefix(tknString, "Bearer ") {
			helpers.ProcessErrorWithStatus("Authorization header is not exists", http.StatusUnauthorized, w)
			return
		}
		tknString = strings.TrimPrefix(tknString, "Bearer ")
		generator := NewJWTGenerator(config.Params.JWTKeys.Private, config.Params.JWTKeys.Public, config.Params.TokenExpiration)
		tkn, err := generator.Parse(tknString)
		if err != nil {
			logger.Log.Info(err)
			helpers.ProcessErrorWithStatus("token is not valid", http.StatusUnauthorized, w)
			return
		}
		idStr, err := tkn.Claims.GetSubject()
		if err != nil {
			logger.Log.Info(err)
			helpers.ProcessErrorWithStatus("token doesnt has user id", http.StatusUnauthorized, w)
			return
		}
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			logger.Log.Info(err)
			helpers.ProcessErrorWithStatus("user id is incorrect", http.StatusUnauthorized, w)
			return
		}

		userRepository := repositories.NewUserRepository(r.Context(), database.DBx)
		user, exists, err := userRepository.GetUserByID(userID)
		if err != nil {
			helpers.SetInternalError(err, w)
			return
		}
		if !exists {
			helpers.ProcessErrorWithStatus("user does not exist", http.StatusUnauthorized, w)
			return
		}

		newR := r.WithContext(context.WithValue(r.Context(), UserKey, user))
		next.ServeHTTP(w, newR)
	})
}

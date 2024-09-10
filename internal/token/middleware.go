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

// AuthMiddleware авторизовываем пользователя по токену и записываем его в контекст
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tknString := r.Header.Get("Authorization")
		if tknString == "" {
			processNotExists("Authorization header is not exists", http.StatusUnauthorized, w)
			return
		}
		if !strings.HasPrefix(tknString, "Bearer ") {
			processNotExists("Authorization header is not exists", http.StatusUnauthorized, w)
			return
		}
		tknString = strings.TrimPrefix(tknString, "Bearer ")
		generator := NewJWTGenerator(config.Params.JWTKeys.Private, config.Params.JWTKeys.Public, config.Params.TokenExpiration)
		tkn, err := generator.Parse(tknString)
		if err != nil {
			logger.Log.Info(err)
			processNotExists("token is not valid", http.StatusUnauthorized, w)
			return
		}
		idStr, err := tkn.Claims.GetSubject()
		if err != nil {
			logger.Log.Info(err)
			processNotExists("token doesnt has user id", http.StatusUnauthorized, w)
			return
		}
		userId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			logger.Log.Info(err)
			processNotExists("user id is incorrect", http.StatusUnauthorized, w)
			return
		}

		userRepository := repositories.NewUserRepository(r.Context(), database.DBx)
		user, exists, err := userRepository.GetUserById(userId)
		if err != nil {
			setInternalError(err, w)
			return
		}
		if !exists {
			processNotExists("user does not exist", http.StatusUnauthorized, w)
		}

		newR := r.WithContext(context.WithValue(r.Context(), "user", user))
		next.ServeHTTP(w, newR)
	})
}

func processNotExists(message string, status int, response http.ResponseWriter) {
	errBody, responseErr := helpers.GetErrorJSONBody(message, status)
	if responseErr != nil {
		logger.Log.Error(errBody)
		if rErr := helpers.SetHTTPResponse(response, http.StatusInternalServerError, []byte{}); rErr != nil {
			logger.Log.Error(rErr)
		}
	}
	if rErr := helpers.SetHTTPResponse(response, status, errBody); rErr != nil {
		logger.Log.Error(rErr)
	}
}

func setInternalError(err error, response http.ResponseWriter) {
	logger.Log.Error(err)
	if rErr := helpers.SetHTTPResponse(response, http.StatusInternalServerError, []byte{}); rErr != nil {
		logger.Log.Error(rErr)
	}
}

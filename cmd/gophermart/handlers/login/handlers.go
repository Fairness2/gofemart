package login

import (
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	config "gofemart/internal/configuration"
	"gofemart/internal/gofemarterrors"
	"gofemart/internal/helpers"
	"gofemart/internal/logger"
	"gofemart/internal/models"
	"gofemart/internal/payloads"
	"gofemart/internal/repositories"
	"gofemart/internal/token"
	"io"
	"net/http"
	"time"
)

type Handlers struct {
	dbPool          repositories.SQLExecutor
	jwtKeys         *config.JWTKeys
	tokenExpiration time.Duration
	hashKey         string
}

func NewHandlers(dbPool repositories.SQLExecutor, jwtKeys *config.JWTKeys, tokenExpiration time.Duration, hashKey string) *Handlers {
	return &Handlers{
		dbPool:          dbPool,
		jwtKeys:         jwtKeys,
		tokenExpiration: tokenExpiration,
		hashKey:         hashKey,
	}
}

func (l *Handlers) RegistrationHandler(response http.ResponseWriter, request *http.Request) {
	// Читаем тело запроса
	body, err := l.getBody(request)
	if err != nil {
		helpers.ProcessRequestErrorWithBody(err, response)
		return
	}

	userRepository := repositories.NewUserRepository(request.Context(), l.dbPool)
	// Проверим есть ли пользователь с таким логином
	exists, err := userRepository.UserExists(body.Login)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if exists {
		helpers.ProcessErrorWithStatus("user already exists", http.StatusConflict, response)
		return
	}

	// Создаём и регистрируем пользователя
	user, err := l.createAndSaveUser(body, userRepository)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}

	// Создаём токен для пользователя
	tkn, err := l.createJWTToken(user)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}

	// Устанавливаем токен в ответ
	payload := payloads.Authorization{Token: tkn}
	responseBody, err := json.Marshal(payload)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}

	response.Header().Set("Authorization", "Bearer "+tkn)

	if rErr := helpers.SetHTTPResponse(response, http.StatusOK, responseBody); rErr != nil {
		logger.Log.Error(rErr)
	}
}

// getBody получаем тело для регистрации
func (l *Handlers) getBody(request *http.Request) (*payloads.Register, error) {
	// Читаем тело запроса
	rawBody, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	// Парсим тело в структуру запроса
	var body payloads.Register
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		return nil, &gofemarterrors.RequestError{InternalError: err, HTTPStatus: http.StatusBadRequest}
	}

	result, err := govalidator.ValidateStruct(body)
	if err != nil {
		return nil, err
	}

	if !result {
		return nil, &gofemarterrors.RequestError{InternalError: errors.New("bad request"), HTTPStatus: http.StatusBadRequest}
	}

	return &body, nil
}

// createUser создаём нового пользователя
func (l *Handlers) createUser(body *payloads.Register) (*models.User, error) {
	user := &models.User{
		Login:    body.Login,
		Password: body.Password,
	}
	err := user.GeneratePasswordHash(l.hashKey)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Создаём и сохраняем нового пользователя
func (l *Handlers) createAndSaveUser(body *payloads.Register, repository *repositories.UserRepository) (*models.User, error) {
	user, err := l.createUser(body)
	if err != nil {
		return nil, err
	}
	if err = repository.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

// createJWTToken создаём JWT токен
func (l *Handlers) createJWTToken(user *models.User) (string, error) {
	generator := token.NewJWTGenerator(l.jwtKeys.Private, l.jwtKeys.Public, l.tokenExpiration)
	return generator.Generate(user)
}

func (l *Handlers) LoginHandler(response http.ResponseWriter, request *http.Request) {
	// Читаем тело запроса
	body, err := l.getBody(request)
	if err != nil {
		helpers.ProcessRequestErrorWithBody(err, response)
		return
	}
	requestedUser, err := l.createUser(body)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}

	userRepository := repositories.NewUserRepository(request.Context(), l.dbPool)
	dbUser, exists, err := userRepository.GetUserByLogin(requestedUser.Login)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if !exists {
		helpers.ProcessErrorWithStatus("password and login are incorrect", http.StatusUnauthorized, response)
		return
	}

	ok, err := dbUser.CheckPasswordHash(requestedUser.PasswordHash)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if !ok {
		helpers.ProcessErrorWithStatus("password and login are incorrect", http.StatusUnauthorized, response)
		return
	}

	// Создаём токен для пользователя
	tkn, err := l.createJWTToken(dbUser)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}

	// Устанавливаем токен в ответ
	payload := payloads.Authorization{Token: tkn}
	responseBody, err := json.Marshal(payload)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}

	response.Header().Set("Authorization", "Bearer "+tkn)

	if rErr := helpers.SetHTTPResponse(response, http.StatusOK, responseBody); rErr != nil {
		logger.Log.Error(rErr)
		helpers.SetInternalError(err, response)
	}
}

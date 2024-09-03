package login

import (
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	config "gofemart/internal/configuration"
	database "gofemart/internal/databse"
	"gofemart/internal/gofemarterrors"
	"gofemart/internal/helpers"
	"gofemart/internal/logger"
	"gofemart/internal/models"
	"gofemart/internal/payloads"
	"gofemart/internal/repositories"
	"gofemart/internal/token"
	"io"
	"net/http"
)

func RegistrationHandler(response http.ResponseWriter, request *http.Request) {
	// Читаем тело запроса
	body, err := getBody(request)
	if err != nil {
		processBadRequestError(err, response)
		return
	}

	userRepository := repositories.NewUserRepository(request.Context(), database.DBx)
	// Проверим есть ли пользователь с таким логином
	exists, err := userRepository.UserExists(body.Login)
	if err != nil {
		setInternalError(err, response)
		return
	}
	if exists {
		processNotExists("user already exists", http.StatusConflict, response)
		return
	}

	// Создаём и регистрируем пользователя
	user, err := createAndSaveUser(body, userRepository)
	if err != nil {
		setInternalError(err, response)
		return
	}

	// Создаём токен для пользователя
	tkn, err := createJWTToken(user)
	if err != nil {
		setInternalError(err, response)
		return
	}

	// Устанавливаем токен в ответ
	payload := payloads.Authorization{Token: tkn}
	responseBody, err := json.Marshal(payload)
	if err != nil {
		setInternalError(err, response)
		return
	}

	if rErr := helpers.SetHTTPResponse(response, http.StatusOK, responseBody); rErr != nil {
		logger.Log.Error(rErr)
	}
}

// getBody получаем тело для регистрации
func getBody(request *http.Request) (*payloads.Register, error) {
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
func createUser(body *payloads.Register) (*models.User, error) {
	user := &models.User{
		Login:    body.Login,
		Password: body.Password,
	}
	err := user.GeneratePasswordHash()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Создаём и сохраняем нового пользователя
func createAndSaveUser(body *payloads.Register, repository *repositories.UserRepository) (*models.User, error) {
	user, err := createUser(body)
	if err != nil {
		return nil, err
	}
	if err = repository.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

// createJWTToken создаём JWT токен
func createJWTToken(user *models.User) (string, error) {
	generator := token.NewJWTGenerator(config.Params.JWTKeys.Private, config.Params.JWTKeys.Public, config.Params.TokenExpiration)
	return generator.Generate(user)
}

func LoginHandler(response http.ResponseWriter, request *http.Request) {
	// Читаем тело запроса
	body, err := getBody(request)
	if err != nil {
		processBadRequestError(err, response)
		return
	}
	requestedUser, err := createUser(body)
	if err != nil {
		setInternalError(err, response)
		return
	}

	userRepository := repositories.NewUserRepository(request.Context(), database.DBx)
	dbUser, exists, err := userRepository.GetUserByLogin(requestedUser.Login)
	if err != nil {
		setInternalError(err, response)
		return
	}
	if !exists {
		processNotExists("password and login are incorrect", http.StatusUnauthorized, response)
		return
	}

	ok, err := dbUser.CheckPasswordHash(requestedUser.PasswordHash)
	if err != nil {
		setInternalError(err, response)
		return
	}
	if !ok {
		processNotExists("password and login are incorrect", http.StatusUnauthorized, response)
		return
	}

	// Создаём токен для пользователя
	tkn, err := createJWTToken(dbUser)
	if err != nil {
		setInternalError(err, response)
		return
	}

	// Устанавливаем токен в ответ
	payload := payloads.Authorization{Token: tkn}
	responseBody, err := json.Marshal(payload)
	if err != nil {
		setInternalError(err, response)
		return
	}

	if rErr := helpers.SetHTTPResponse(response, http.StatusOK, responseBody); rErr != nil {
		logger.Log.Error(rErr)
		setInternalError(err, response)
	}
}

func processBadRequestError(err error, response http.ResponseWriter) {
	var errWithStatus *gofemarterrors.RequestError
	var errBody []byte
	var responseErr error
	var httpStatus int
	if errors.As(err, &errWithStatus) {
		logger.Log.Info(err)
		httpStatus = errWithStatus.HTTPStatus
		errBody, responseErr = helpers.GetErrorJSONBody(errWithStatus.Error(), errWithStatus.HTTPStatus)
	} else {
		logger.Log.Error(err)
		httpStatus = http.StatusInternalServerError
		errBody, responseErr = helpers.GetErrorJSONBody(err.Error(), http.StatusInternalServerError)
	}
	if responseErr != nil {
		logger.Log.Error(err)
		if rErr := helpers.SetHTTPResponse(response, http.StatusInternalServerError, []byte{}); rErr != nil {
			logger.Log.Error(rErr)
		}
	} else {
		if rErr := helpers.SetHTTPResponse(response, httpStatus, errBody); rErr != nil {
			logger.Log.Error(rErr)
		}
	}
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

package registration

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

func Handler(response http.ResponseWriter, request *http.Request) {
	// Читаем тело запроса
	body, err := getBody(request)
	if err != nil {
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
		return
	}

	userRepository := repositories.NewUserRepository(request.Context(), database.DBx)
	// Проверим есть ли пользователь с таким логином
	exists, err := userRepository.UserExists(body.Login)
	if err != nil {
		logger.Log.Error(err)
		if rErr := helpers.SetHTTPResponse(response, http.StatusInternalServerError, []byte{}); rErr != nil {
			logger.Log.Error(rErr)
		}
	}
	if exists {
		logger.Log.Infow("user already exists", "login", body.Login)
		errBody, responseErr := helpers.GetErrorJSONBody("user already exists", http.StatusConflict)
		if responseErr != nil {
			logger.Log.Error(err)
			if rErr := helpers.SetHTTPResponse(response, http.StatusInternalServerError, []byte{}); rErr != nil {
				logger.Log.Error(rErr)
			}
		}
		if rErr := helpers.SetHTTPResponse(response, http.StatusConflict, errBody); rErr != nil {
			logger.Log.Error(rErr)
		}
		return
	}

	// Создаём и регистрируем пользователя
	user, err := createAndSaveUser(body, userRepository)
	if err != nil {
		logger.Log.Error(err)
		if rErr := helpers.SetHTTPResponse(response, http.StatusInternalServerError, []byte{}); rErr != nil {
			logger.Log.Error(rErr)
		}
	}

	// Создаём токен для пользователя
	tkn, err := createJWTToken(user)
	if err != nil {
		logger.Log.Error(err)
		if rErr := helpers.SetHTTPResponse(response, http.StatusInternalServerError, []byte{}); rErr != nil {
			logger.Log.Error(rErr)
		}
	}

	// Устанавливаем токен в ответ
	payload := payloads.Authorization{Token: tkn}
	responseBody, err := json.Marshal(payload)
	if err != nil {
		logger.Log.Error(err)
		if rErr := helpers.SetHTTPResponse(response, http.StatusInternalServerError, []byte{}); rErr != nil {
			logger.Log.Error(rErr)
		}
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
		return nil, &gofemarterrors.RequestError{InternalError: errors.New("bad request for registration"), HTTPStatus: http.StatusBadRequest}
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

package helpers

import (
	"encoding/json"
	"errors"
	"gofemart/internal/gofemarterrors"
	"gofemart/internal/logger"
	"gofemart/internal/payloads"
	"net/http"
)

// SetHTTPResponse Отправка ошибки и сообщения ошибки.
// Parameters:
// - response: http.ResponseWriter object containing information about the HTTP response
// - status: the HTTP status code to set in the response
// - message: the message to write to the response
func SetHTTPResponse(response http.ResponseWriter, status int, message []byte) error {
	response.WriteHeader(status)
	_, err := response.Write(message) // TODO подумать, нужно ли
	return err
}

// GetErrorJSONBody Создание тела ответа с json ошибкой
func GetErrorJSONBody(message string, statue int) ([]byte, error) {
	responseBody := payloads.ErrorResponseBody{
		Status:  statue,
		Message: message,
	}
	return json.Marshal(responseBody)
}

func SetInternalError(err error, response http.ResponseWriter) {
	logger.Log.Error(err)
	if rErr := SetHTTPResponse(response, http.StatusInternalServerError, []byte{}); rErr != nil {
		logger.Log.Error(rErr)
	}
}

func ProcessRequestErrorWithBody(err error, response http.ResponseWriter) {
	var errWithStatus *gofemarterrors.RequestError
	var errBody []byte
	var responseErr error
	var httpStatus int
	if errors.As(err, &errWithStatus) {
		logger.Log.Info(err)
		httpStatus = errWithStatus.HTTPStatus
		errBody, responseErr = GetErrorJSONBody(errWithStatus.Error(), errWithStatus.HTTPStatus)
	} else {
		logger.Log.Error(err)
		httpStatus = http.StatusInternalServerError
		errBody, responseErr = GetErrorJSONBody(err.Error(), http.StatusInternalServerError)
	}
	if responseErr != nil {
		SetInternalError(err, response)
	} else {
		if rErr := SetHTTPResponse(response, httpStatus, errBody); rErr != nil {
			logger.Log.Error(rErr)
		}
	}
}

func ProcessErrorWithStatus(message string, status int, response http.ResponseWriter) {
	errBody, responseErr := GetErrorJSONBody(message, status)
	if responseErr != nil {
		SetInternalError(responseErr, response)
		return
	}
	if rErr := SetHTTPResponse(response, status, errBody); rErr != nil {
		logger.Log.Error(rErr)
	}
}

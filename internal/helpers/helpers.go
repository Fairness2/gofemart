package helpers

import (
	"encoding/json"
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

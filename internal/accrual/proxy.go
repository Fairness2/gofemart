package accrual

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"gofemart/internal/logger"
	"gofemart/internal/models"
	"gofemart/internal/payloads"
	"net/http"
	"sync"
	"time"
)

// getOrderURL эндпоинт для получения начислений по заказам
const getOrderURL = "/api/orders/"

// ErrorOrderNotRegistered указывает на то, что заказ не зарегистрирован в сервисе начисления.
var ErrorOrderNotRegistered = errors.New("order not registered in accrual service")

// ErrorInternalAccrual указывает на то, что в сервисе начисления произошла внутренняя ошибка.
var ErrorInternalAccrual = errors.New("accrual service internal error")

// ErrorTooManyRequests указывает, что сервер получил слишком много запросов.
var ErrorTooManyRequests = errors.New("too many requests")

// ErrorUnknownStatusRequests указывает,что сервер получил непредвиденный ответ
var ErrorUnknownStatusRequests = errors.New("unknown accrual error")

// TooManyRequestError представляет собой ошибку, указывающую, что к службе было отправлено слишком много запросов.
// Он включает внутреннюю ошибку и рекомендуемую продолжительность паузы перед повторной попыткой запроса.
type TooManyRequestError struct {
	InternalError error
	PauseDuration time.Duration
}

func (e *TooManyRequestError) Error() string {
	return e.InternalError.Error()
}
func (e *TooManyRequestError) Unwrap() error {
	return e.InternalError
}

// Proxy представляет клиент, который обрабатывает связь с внешней службой с возможностями ограничения скорости и паузы.
type Proxy struct {
	pauseDuration time.Duration
	client        *resty.Client
	senderMutex   sync.RWMutex
}

// NewProxy создает новый экземпляр Proxy с указанной длительностью паузы и URL-адресом службы начисления
func NewProxy(pause time.Duration, accrualURL string) *Proxy {
	client := resty.New()
	client = client.SetBaseURL(accrualURL)
	return &Proxy{
		pauseDuration: pause,
		client:        client,
		senderMutex:   sync.RWMutex{},
	}
}

// Pause Если мы попали в блок от системы, делаем паузу между запросами
func (p *Proxy) Pause(duration time.Duration) {
	logger.Log.Infow("Pause", "duration", duration)
	p.senderMutex.Lock()
	defer p.senderMutex.RUnlock()
	<-time.After(duration)
}

// Accrual запрашиваем статус заказа в системе начислений
func (p *Proxy) Accrual(order *models.Order) (*payloads.Accrual, error) {
	logger.Log.Infow("Accrual", "order", order.Number)
	p.senderMutex.RLock()
	defer p.senderMutex.RUnlock()
	url := getOrderURL + order.Number
	request := p.client.R()
	request.SetHeader("Content-Type", "application/json")
	response, err := request.Get(url)
	if err != nil {
		return nil, err
	}
	switch response.StatusCode() {
	case http.StatusNoContent:
		logger.Log.Infow("Order not registered", "order", order.Number, "status", http.StatusNoContent)
		return nil, ErrorOrderNotRegistered
	case http.StatusInternalServerError:
		logger.Log.Infow("Accrual check failed", "order", order.Number, "status", http.StatusInternalServerError)
		return nil, ErrorInternalAccrual
	case http.StatusTooManyRequests:
		logger.Log.Infow("Too many requests", "order", order.Number, "status", http.StatusTooManyRequests)
		pauseDuration := p.pauseDuration
		if pauseHeader := response.Header().Get("Retry-After"); pauseHeader != "" {
			pauseHeaderValue, err := time.ParseDuration(pauseHeader)
			if err == nil {
				pauseDuration = pauseHeaderValue
			} else {
				logger.Log.Error(err)
			}
		}
		return nil, &TooManyRequestError{InternalError: ErrorTooManyRequests, PauseDuration: pauseDuration}
	case http.StatusOK:
		logger.Log.Infow("Order registered", "order", order.Number, "status", http.StatusOK)
		return p.processAccrualResponse(response)
	default:
		logger.Log.Infow("Unknown status", "order", order.Number, "status", response.StatusCode())
		return nil, ErrorUnknownStatusRequests
	}
}

// processAccrualResponse обрабатывает ответ от службы начисления и преобразует его в структуру начисления.
func (p *Proxy) processAccrualResponse(res *resty.Response) (*payloads.Accrual, error) {
	// Парсим тело в структуру запроса
	var body payloads.Accrual
	err := json.Unmarshal(res.Body(), &body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

package ordercheck

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	database "gofemart/internal/databse"
	"gofemart/internal/logger"
	"gofemart/internal/models"
	"gofemart/internal/payloads"
	"gofemart/internal/repositories"
	"sync"
	"time"
)

var ErrorPoolClosed = errors.New("pool closed")
var ErrorOrderNotRegistered = errors.New("order not registered in accrual service")
var ErrorInternalAccrual = errors.New("accrual service internal error")
var ErrorTooManyRequests = errors.New("too many requests")

type tooManyRequestError struct {
	InternalError error
	pauseDuration time.Duration
}

func (e *tooManyRequestError) Error() string {
	return e.InternalError.Error()
}
func (e *tooManyRequestError) Unwrap() error {
	return e.InternalError
}

type WorkedOrder struct {
	model  *models.Order
	inWork bool
}

type Pool struct {
	closeFlag bool
	orderMap  map[string]*WorkedOrder
	//queue     []string
	//inWork    []string
	mutex             sync.RWMutex
	inChanel          chan string
	ctx               context.Context
	wg                sync.WaitGroup
	cancel            context.CancelFunc
	pauseDuration     time.Duration
	client            *resty.Client
	senderMutex       sync.RWMutex
	olderThenDuration time.Duration
}

var CheckPool *Pool

func NewPool(ctx context.Context, queueSize int, workerCount int, pause time.Duration, accrualURL string) *Pool {
	inChanel := make(chan string, queueSize)
	poolContext, cancel := context.WithCancel(ctx)
	client := resty.New()
	client = client.SetBaseURL(accrualURL)
	pool := &Pool{
		//inWork:   make([]*models.Order, 0, workerCount),
		mutex:             sync.RWMutex{},
		inChanel:          inChanel,
		ctx:               poolContext,
		cancel:            cancel,
		orderMap:          make(map[string]*WorkedOrder),
		wg:                sync.WaitGroup{},
		pauseDuration:     pause,
		client:            client,
		senderMutex:       sync.RWMutex{},
		olderThenDuration: time.Minute * 10,
	}

	// Запускаем проверку закрытия
	go pool.finishWork()

	// Запускаем воркеры
	for i := 0; i < workerCount; i++ {
		pool.wg.Add(1)
		go pool.pushFromQueue()
	}

	// Запускаем проверку базы данных
	pool.wg.Add(1)
	go pool.pushFromDB(time.Minute)

	return pool
}

// Close функция закрытия пула, закрываем локальный контекст, ждём завершения всех воркеров, закрываем канал очереди
func (p *Pool) Close() {
	p.closeFlag = true
	p.cancel()
	p.wg.Wait()
	close(p.inChanel)
}

// Push adds an Order to the Pool's queue if the queue is not full and the Pool is not closed, returning success status and error.
func (p *Pool) Push(order *models.Order) (bool, error) { // TODO ограниченная очередь
	if p.closeFlag {
		return false, ErrorPoolClosed
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if len(p.inChanel) < cap(p.inChanel) {
		p.orderMap[order.Number] = &WorkedOrder{model: order}
		p.inChanel <- order.Number
		return true, nil
	}
	return false, nil
}

// pushFromQueue обработчик очереди заказов
func (p *Pool) pushFromQueue() {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		case number := <-p.inChanel:
			if !p.checkInWork(number) {
				p.processOder(number)
			}
		}
	}
}

// finishWork закрываем пул
func (p *Pool) finishWork() {
	<-p.ctx.Done()
	p.Close()
}

// poolInWork взятие в работу ордера
func (p *Pool) poolInWork(number string) (*WorkedOrder, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if order, ok := p.orderMap[number]; ok {
		if order.inWork {
			return nil, false
		}
		order.inWork = true
		return order, true
	}
	return nil, false
}

// processOder обрабатываем заказ, запрашивая информацию у внешней системы
func (p *Pool) processOder(number string) {
	order, ok := p.poolInWork(number)
	// Если заказ не найден, то пропускаем его
	if order == nil || !ok {
		return
	}
	defer p.removeFromWork(number)
	accrualResponse, err := p.accrual(order.model)
	if err != nil {
		logger.Log.Error(err)
		var tmrErr *tooManyRequestError
		if ok := errors.Is(err, tmrErr); ok {
			p.pause(tmrErr.pauseDuration)
		}
		return
	}
	if err := p.processOrderAccrual(accrualResponse, order.model); err != nil {
		logger.Log.Error(err)
	}
}

// removeFromWork снимаем флаг взятия в работу
func (p *Pool) removeFromWork(number string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if order, ok := p.orderMap[number]; ok {
		order.inWork = false
	}
}

// deleteFromMap удаляем отработанный заказ
func (p *Pool) deleteFromMap(number string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.orderMap, number)
}

// pause Если мы попали в блок от системы, делаем паузу между запросами
func (p *Pool) pause(duration time.Duration) {
	p.senderMutex.Lock()
	defer p.senderMutex.RUnlock()
	<-time.After(duration)
}

// accrual запрашиваем статус заказа в системе начислений
func (p *Pool) accrual(order *models.Order) (*payloads.Accrual, error) {
	p.senderMutex.RLock()
	defer p.senderMutex.RUnlock()
	url := "/api/orders/" + order.Number
	request := p.client.R()
	request.SetHeader("Content-Type", "application/json")
	response, err := request.Get(url)
	if err != nil {
		return nil, err
	}
	switch response.StatusCode() {
	case 204:
		return nil, ErrorOrderNotRegistered
	case 500:
		return nil, ErrorInternalAccrual
	case 429:
		pauseDuration := p.pauseDuration
		if pauseHeader := response.Header().Get("Retry-After"); pauseHeader != "" {
			pauseHeaderValue, err := time.ParseDuration(pauseHeader)
			if err == nil {
				pauseDuration = pauseHeaderValue
			} else {
				logger.Log.Error(err)
			}
		}
		return nil, &tooManyRequestError{InternalError: ErrorTooManyRequests, pauseDuration: pauseDuration}
	case 200:
		return p.processAccrualResponse(response)
	default:
		return nil, errors.New("unknown accrual error")
	}
}

// processAccrualResponse processes the response from the accrual service and parses it into an Accrual structure.
func (p *Pool) processAccrualResponse(res *resty.Response) (*payloads.Accrual, error) {
	// Парсим тело в структуру запроса
	var body payloads.Accrual
	err := json.Unmarshal(res.Body(), &body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

// processOrderAccrual обрабатываем ответ системы начислений, обновляем заказ и создаём запись в счёте пользователя
func (p *Pool) processOrderAccrual(accrual *payloads.Accrual, order *models.Order) error {
	order.LastCheckedAt = sql.NullTime{Time: time.Now(), Valid: true}
	switch accrual.Status {
	case payloads.StatusAccrualProcessing, payloads.StatusAccrualRegistered:
		order.StatusCode = models.StatusProcessing
	case payloads.StatusAccrualInvalid:
		order.StatusCode = models.StatusInvalid
	case payloads.StatusAccrualProcessed:
		if _, err := p.createNewAccount(order.Number, order.UserId, accrual.Accrual); err != nil {
			return err
		}
		order.StatusCode = models.StatusProcessed
	}
	orderRep := p.getOrderRepository()
	return orderRep.UpdateOrder(order)
}

// getAccountRepository создаём репозиторий для начислений
func (p *Pool) getAccountRepository() *repositories.AccountRepository {
	return repositories.NewAccountRepository(p.ctx, database.DBx)
}

// getOrderRepository создаём репозиторий заказов
func (p *Pool) getOrderRepository() *repositories.OrderRepository {
	return repositories.NewOrderRepository(p.ctx, database.DBx)
}

// createNewAccount создаём новую запись о начислении
func (p *Pool) createNewAccount(orderNumber string, userId int64, diff int) (*models.Account, error) {
	repository := p.getAccountRepository()
	account := models.NewAccount(sql.NullString{String: orderNumber, Valid: true}, userId, diff)
	if err := repository.CreateAccount(account); err != nil {
		return nil, err
	}
	return account, nil
}

// checkInWork проверяем находится ли заказ в работе
func (p *Pool) checkInWork(number string) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if order, ok := p.orderMap[number]; ok {
		return order.inWork
	}
	return false
}

// pushFromDB периодическое пополнение очереди из базы данных
func (p *Pool) pushFromDB(dur time.Duration) {
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-p.ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			if err := p.pushDBProcessingOrdersToQueue(); err != nil {
				logger.Log.Error(err)
			}
		}
	}
}

// getCurrentOrdersKeys получаем номера заказов в очереди
func (p *Pool) getCurrentOrdersKeys() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	keys := make([]string, 0, len(p.orderMap))
	for key := range p.orderMap {
		keys = append(keys, key)
	}

	return keys
}

// pushDBProcessingOrdersToQueue получаем из базы данных необработанные заказы и пушим их в очередь
func (p *Pool) pushDBProcessingOrdersToQueue() error {
	limit := cap(p.inChanel) - len(p.inChanel)
	if limit <= 0 {
		return nil
	}
	keys := p.getCurrentOrdersKeys()
	rep := p.getOrderRepository()
	olderThen := time.Now().Add(-p.olderThenDuration)
	orders, err := rep.GetOrdersExcludeOrdersWhereStatusIn(limit, keys, olderThen, models.StatusProcessing, models.StatusNew)
	if err != nil {
		return err
	}
	for _, order := range orders {
		if _, err := p.Push(&order); err != nil {
			return err
		}
	}

	return nil
}

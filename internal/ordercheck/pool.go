package ordercheck

import (
	"context"
	"errors"
	"gofemart/internal/accrual"
	"gofemart/internal/logger"
	"gofemart/internal/models"
	"gofemart/internal/payloads"
	"gofemart/internal/repositories"
	"sync"
	"sync/atomic"
	"time"
)

// ErrorPoolClosed Ошибка, что пул обработки уже закрыт
var ErrorPoolClosed = errors.New("pool closed")

// oRepo определяет методы взаимодействия с заказами в репозитории.
type oRepo interface {
	GetOrdersExcludeOrdersWhereStatusIn(limit int, excludedNumbers []string, olderThen time.Time, statuses ...string) ([]models.Order, error)
	UpdateOrder(order *models.Order) error
}

// aRepo определяет интерфейс для взаимодействия с начислениями в репозитории.
type aRepo interface {
	CreateAccount(account *models.Account) error
}

// Accrual предоставляет методы для проверки начислений и управления паузами.
type Accrual interface {
	Accrual(order *models.Order) (*payloads.Accrual, error)
	Pause(duration time.Duration)
}

// WorkedOrder представляет собой обрабатываемый заказ.
type WorkedOrder struct {
	model  *models.Order
	inWork bool
}

// Pool управляет обработкой заказов, включая организацию очередей,
// параллельную обработку и взаимодействие с внешними системами.
type Pool struct {
	closeFlag         atomic.Bool
	orderMap          map[string]*WorkedOrder
	mutex             sync.RWMutex
	inChanel          chan string
	ctx               context.Context
	wg                sync.WaitGroup
	cancel            context.CancelFunc
	olderThenDuration time.Duration
	orderRepo         oRepo
	accountRepo       aRepo
	accrualProxy      Accrual
}

// CheckPool глобальный инстенс пула обработки заказов.
var CheckPool *Pool

// PoolConfig Конфигурация для пула обработки
type PoolConfig struct {
	CTX         context.Context
	QueueSize   int
	WorkerCount int
	Pause       time.Duration
	AccrualURL  string
	DBExecutor  repositories.SQLExecutor
}

// NewPool инициализирует и возвращает новый экземпляр Pool с указанным контекстом, размером очереди, количеством рабочих процессов, длительностью паузы и URL-адресом накопления.
func NewPool(cnf PoolConfig) *Pool {
	logger.Log.Infow("New pool", "queueSize", cnf.QueueSize, "workerCount", cnf.WorkerCount, "pause", cnf.Pause, "accrualURL", cnf.AccrualURL)
	inChanel := make(chan string, cnf.QueueSize)
	poolContext, cancel := context.WithCancel(cnf.CTX)
	proxy := accrual.NewProxy(cnf.Pause, cnf.AccrualURL)

	pool := &Pool{
		mutex:             sync.RWMutex{},
		inChanel:          inChanel,
		ctx:               poolContext,
		cancel:            cancel,
		orderMap:          make(map[string]*WorkedOrder),
		wg:                sync.WaitGroup{},
		olderThenDuration: time.Second * 5,
		accountRepo:       getAccountRepository(cnf.CTX, cnf.DBExecutor),
		orderRepo:         getOrderRepository(cnf.CTX, cnf.DBExecutor),
		accrualProxy:      proxy,
	}
	initPool(cnf.WorkerCount, pool)

	return pool
}

// initPool инициализирует и запускает пул рабочих процессов,
// запускает рабочие процессы и планирует проверки базы данных.
func initPool(workerCount int, pool *Pool) {
	// Запускаем проверку закрытия
	go pool.finishWork()
	// Запускаем воркеры
	pool.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go pool.pushFromQueue()
	}
	// Запускаем проверку базы данных
	pool.wg.Add(1)
	go pool.pushFromDB(5 * time.Second)
}

// Close функция закрытия пула, закрываем локальный контекст, ждём завершения всех воркеров, закрываем канал очереди
func (p *Pool) Close() {
	logger.Log.Info("Close pool")
	if p.closeFlag.Load() {
		return
	}
	p.closeFlag.Store(true)
	p.cancel()
	p.wg.Wait()
	close(p.inChanel)
}

// finishWork закрываем пул
func (p *Pool) finishWork() {
	logger.Log.Info("Finish work")
	<-p.ctx.Done()
	p.Close()
}

// getAccountRepository создаём репозиторий для начислений
func getAccountRepository(ctx context.Context, executor repositories.SQLExecutor) *repositories.AccountRepository {
	logger.Log.Infow("Get account repository")
	return repositories.NewAccountRepository(ctx, executor)
}

// getOrderRepository создаём репозиторий заказов
func getOrderRepository(ctx context.Context, executor repositories.SQLExecutor) *repositories.OrderRepository {
	logger.Log.Infow("Get order repository")
	return repositories.NewOrderRepository(ctx, executor)
}

package ordercheck

import (
	"gofemart/internal/logger"
	"gofemart/internal/models"
	"time"
)

// pushFromDB периодическое пополнение очереди из базы данных
func (p *Pool) pushFromDB(dur time.Duration) {
	logger.Log.Infow("Push from db", "duration", dur)
	defer p.wg.Done()
	if err := p.pushDBProcessingOrdersToQueue(); err != nil {
		logger.Log.Error(err)
	}
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

// pushDBProcessingOrdersToQueue получаем из базы данных необработанные заказы и пушим их в очередь
func (p *Pool) pushDBProcessingOrdersToQueue() error {
	limit := cap(p.inChanel) - len(p.inChanel)
	logger.Log.Infow("Push db processing orders to queue", "limit", limit)
	if limit <= 0 {
		return nil
	}
	keys := p.getCurrentOrdersKeys()
	rep := p.orderRepo
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

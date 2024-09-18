package ordercheck

import (
	"gofemart/internal/logger"
	"gofemart/internal/models"
)

// Push добавляет заказ в очередь пула,
// если очередь не заполнена и пул не закрыт, возвращая статус успешного выполнения и ошибку.
func (p *Pool) Push(order *models.Order) (bool, error) { // TODO ограниченная очередь
	logger.Log.Infow("Push order to pool", "order", order.Number)
	if p.closeFlag.Load() {
		return false, ErrorPoolClosed
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	lenChanel := len(p.inChanel)
	capChanel := cap(p.inChanel)
	if lenChanel < capChanel {
		logger.Log.Infow("Push order to queue", "order", order.Number, "len", lenChanel, "cap", capChanel)
		p.orderMap[order.Number] = &WorkedOrder{model: order}
		p.inChanel <- order.Number
		return true, nil
	}
	logger.Log.Infow("Order not pushed to queue", "order", order.Number, "len", lenChanel, "cap", capChanel)
	return false, nil
}

// pushFromQueue обработчик очереди заказов
func (p *Pool) pushFromQueue() {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			logger.Log.Info("Pool context closed. Push from queue stopped")
			return
		case number, ok := <-p.inChanel:
			logger.Log.Infow("Push from queue", "number", number, "ok", ok)
			if !ok {
				return
			}
			if !p.checkInWork(number) {
				p.processOder(number)
			}
		}
	}
}

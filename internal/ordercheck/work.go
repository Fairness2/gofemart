package ordercheck

import "gofemart/internal/logger"

// poolInWork взятие в работу ордера
func (p *Pool) poolInWork(number string) (*WorkedOrder, bool) {
	logger.Log.Infow("Pool in work", "number", number)
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

// removeFromWork снимаем флаг взятия в работу
func (p *Pool) removeFromWork(number string) {
	logger.Log.Infow("Remove from work", "number", number)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if order, ok := p.orderMap[number]; ok {
		order.inWork = false
	}
}

// checkInWork проверяем находится ли заказ в работе
func (p *Pool) checkInWork(number string) bool {
	logger.Log.Infow("Check in work", "number", number)
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if order, ok := p.orderMap[number]; ok {
		return order.inWork
	}
	return false
}

package ordercheck

import "gofemart/internal/logger"

// deleteFromMap удаляем отработанный заказ
func (p *Pool) deleteFromMap(number string) {
	logger.Log.Infow("Delete from map", "number", number)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.orderMap, number)
}

// getCurrentOrdersKeys получаем номера заказов в очереди
func (p *Pool) getCurrentOrdersKeys() []string {
	logger.Log.Infow("Get current orders keys")
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	keys := make([]string, 0, len(p.orderMap))
	for key := range p.orderMap {
		keys = append(keys, key)
	}

	return keys
}

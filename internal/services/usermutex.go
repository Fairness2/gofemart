package services

import "sync"

// UserMutex управляет коллекцией мьютексов, связанных с идентификаторами пользователей,
// обеспечивая безопасный одновременный доступ для изменения баланса.
type UserMutex struct {
	usersMap map[int64]*sync.Mutex
	mutex    sync.RWMutex
}

// UserMutexInstance глобальный инстенс сервиса
var UserMutexInstance *UserMutex

func init() {
	UserMutexInstance = NewUserMutex()
}

// NewUserMutex создание сервиса мьютексов пользователей
func NewUserMutex() *UserMutex {
	return &UserMutex{
		usersMap: make(map[int64]*sync.Mutex),
		mutex:    sync.RWMutex{},
	}
}

// GetMutex получение мьютекса пользователя
func (um *UserMutex) GetMutex(userID int64) (*sync.Mutex, bool) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	mutex, ok := um.usersMap[userID]
	return mutex, ok
}

// SetMutex создаём новый мьютекс пользователя
func (um *UserMutex) SetMutex(userID int64) *sync.Mutex {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	if mutex, ok := um.usersMap[userID]; ok {
		return mutex
	}
	mutex := &sync.Mutex{}
	um.usersMap[userID] = mutex
	return mutex
}

// DeleteMutex удаляем использованный мьютекс
func (um *UserMutex) DeleteMutex(userID int64) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	mutex, ok := um.usersMap[userID]
	if !ok {
		return
	}
	if mutex.TryLock() {
		delete(um.usersMap, userID)
	}
}

package services

import (
	"errors"
	"gofemart/internal/logger"
	"sync"
)

// ErrorUserMutexInUse Ошибка, что удалить мьютекс пользователя не получится, так как он используется
var ErrorUserMutexInUse = errors.New("user mutex in use")

// UserMutex управляет коллекцией мьютексов, связанных с идентификаторами пользователей,
// обеспечивая безопасный одновременный доступ для изменения баланса.
type UserMutex struct {
	usersMap map[int64]*sync.Mutex
	mutex    sync.RWMutex
}

func GetUserMutexInstance() *UserMutex {
	return userMutexInstance
}

// userMutexInstance глобальный инстенс сервиса
var userMutexInstance *UserMutex = newUserMutex()

// newUserMutex создание сервиса мьютексов пользователей
func newUserMutex() *UserMutex {
	return &UserMutex{
		usersMap: make(map[int64]*sync.Mutex),
		mutex:    sync.RWMutex{},
	}
}

// GetMutex получение мьютекса пользователя
func (um *UserMutex) GetMutex(userID int64) (*sync.Mutex, bool) {
	logger.Log.Debugf("Get mutex userID: %d", userID)
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	mutex, ok := um.usersMap[userID]
	return mutex, ok
}

// SetMutex создаём новый мьютекс пользователя
func (um *UserMutex) SetMutex(userID int64) *sync.Mutex {
	logger.Log.Debugf("Set mutex userID: %d", userID)
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
func (um *UserMutex) DeleteMutex(userID int64) error {
	logger.Log.Debugf("Delete mutex userID: %d", userID)
	um.mutex.Lock()
	defer um.mutex.Unlock()
	mutex, ok := um.usersMap[userID]
	if !ok {
		return nil
	}
	if mutex.TryLock() {
		delete(um.usersMap, userID)
		return nil
	}
	return ErrorUserMutexInUse
}

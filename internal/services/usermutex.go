package services

import "sync"

type UserMutex struct { // TODO удаление неиспользуемых мьютекосв
	usersMap map[int64]*sync.Mutex
	mutex    sync.RWMutex
}

var UserMutexInstance *UserMutex

func init() {
	UserMutexInstance = NewUserMutex()
}

func NewUserMutex() *UserMutex {
	return &UserMutex{
		usersMap: make(map[int64]*sync.Mutex),
		mutex:    sync.RWMutex{},
	}
}

func (um *UserMutex) GetMutex(userID int64) (*sync.Mutex, bool) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	mutex, ok := um.usersMap[userID]
	return mutex, ok
}

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

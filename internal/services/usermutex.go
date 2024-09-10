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

func (um *UserMutex) GetMutex(userId int64) (*sync.Mutex, bool) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	mutex, ok := um.usersMap[userId]
	return mutex, ok
}

func (um *UserMutex) SetMutex(userId int64) *sync.Mutex {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	if mutex, ok := um.usersMap[userId]; ok {
		return mutex
	}
	mutex := &sync.Mutex{}
	um.usersMap[userId] = mutex
	return mutex
}

func (um *UserMutex) DeleteMutex(userId int64) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	mutex, ok := um.usersMap[userId]
	if !ok {
		return
	}
	if mutex.TryLock() {
		delete(um.usersMap, userId)
	}
}

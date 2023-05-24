package broadcaster

import "sync"

type UserManager struct {
	data map[string]User
	mu   sync.RWMutex
}

func newUserManager() *UserManager {
	return &UserManager{
		data: map[string]User{},
		mu:   sync.RWMutex{},
	}
}

func (manager *UserManager) Add(user User) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	manager.data[user.ID()] = user
}

func (manager *UserManager) Delete(id string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	delete(manager.data, id)
}

func (manager *UserManager) Get(id string) (User, bool) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	user, ok := manager.data[id]
	return user, ok
}

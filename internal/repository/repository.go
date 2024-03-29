package repository

import (
	"sync"

	"manager/internal/domain"
)

type Repository struct {
	storage domain.Storage
	mutex   *sync.RWMutex
}

func New() *Repository {
	return &Repository{
		storage: make(domain.Storage),
		mutex:   new(sync.RWMutex),
	}
}

func (r *Repository) SetStorage(storage domain.Storage) {
	copyStorage := make(domain.Storage, len(storage))

	for name, elem := range storage {
		copyStorage[name] = elem
	}

	r.mutex.RLock()
	r.storage = copyStorage
	r.mutex.RUnlock()
}

func (r *Repository) Get(name string) (domain.Service, bool) {
	r.mutex.RLock()
	service, ok := r.storage[name]
	r.mutex.RUnlock()

	if !ok {
		return domain.Service{}, false
	}

	return service, true
}

func (r *Repository) GetAll() domain.Storage {
	r.mutex.RLock()
	copyStorage := make(domain.Storage, len(r.storage))

	for id, elem := range r.storage {
		copyStorage[id] = elem
	}

	r.mutex.RUnlock()

	return copyStorage
}

func (r *Repository) Reset() {
	r.mutex.Lock()
	r.storage = make(domain.Storage)
	r.mutex.Unlock()
}

// login

func (r *Repository) AppendLogin(name, login string, elem domain.Element) bool {
	r.mutex.Lock()

	_, ok := r.storage[name]
	if ok {
		r.mutex.Unlock()

		return false
	}

	_, ok = r.storage[name].Elements[login]
	if ok {
		r.mutex.Unlock()

		return false
	}

	r.storage[name].Elements[login] = elem
	r.mutex.Unlock()

	return true
}

func (r *Repository) UpdateLogin(name, login string, elem domain.Element) bool {
	r.mutex.Lock()
	_, ok := r.storage[name]
	if !ok {
		r.mutex.Unlock()

		return false
	}

	_, ok = r.storage[name].Elements[login]
	if !ok {
		r.mutex.Unlock()

		return false
	}

	r.storage[name].Elements[login] = elem
	r.mutex.Unlock()

	return true
}

func (r *Repository) DeleteLogin(name, login string) bool {
	r.mutex.Lock()

	_, ok := r.storage[name]
	if !ok {
		r.mutex.Unlock()

		return false
	}

	_, ok = r.storage[name].Elements[login]
	if !ok {
		r.mutex.Unlock()

		return false
	}

	delete(r.storage[name].Elements, login)
	r.mutex.Unlock()

	return true
}

// service

func (r *Repository) AppendService(name, serviceType string, favorite bool) bool {
	r.mutex.Lock()

	_, ok := r.storage[name]
	if ok {
		r.mutex.Unlock()

		return false
	}

	r.storage[name] = domain.Service{
		Type:     serviceType,
		Favorite: favorite,
		Elements: nil,
	}
	r.mutex.Unlock()

	return true
}

func (r *Repository) UpdateService(name, serviceType string, favorite bool) bool {
	r.mutex.Lock()
	service, ok := r.storage[name]
	if !ok {
		r.mutex.Unlock()

		return false
	}

	r.storage[name] = domain.Service{
		Type:     serviceType,
		Favorite: favorite,
		Elements: service.Elements,
	}

	r.mutex.Unlock()

	return true
}

func (r *Repository) DeleteService(name string) bool {
	r.mutex.Lock()

	_, ok := r.storage[name]
	if !ok {
		r.mutex.Unlock()

		return false
	}

	delete(r.storage, name)
	r.mutex.Unlock()

	return true
}

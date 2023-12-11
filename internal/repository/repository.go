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

func (r *Repository) Get(name string) (domain.Element, bool) {
	r.mutex.RLock()
	elem, ok := r.storage[name]
	r.mutex.RUnlock()

	if !ok {
		return domain.Element{}, false
	}

	return elem, true
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

func (r *Repository) Update(name string, elem domain.Element) bool {
	r.mutex.Lock()
	_, ok := r.storage[name]
	if !ok {
		r.mutex.Unlock()

		return false
	}

	r.storage[name] = elem
	r.mutex.Unlock()

	return true
}

func (r *Repository) Append(name string, elem domain.Element) bool {
	r.mutex.Lock()

	_, ok := r.storage[name]
	if ok {
		r.mutex.Unlock()

		return false
	}

	r.storage[name] = elem
	r.mutex.Unlock()

	return true
}

func (r *Repository) DeleteByName(name string) bool {
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

func (r *Repository) Reset() {
	r.mutex.Lock()
	r.storage = make(domain.Storage)
	r.mutex.Unlock()
}

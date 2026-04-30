package repository

import (
	"sync"
)

type MemoryRepo struct{
	mu sync.Mutex
	byURL map[string]string
	byKey map[string]string
}

func NewMemoryRepo() *MemoryRepo{
	return &MemoryRepo{
		byURL: make(map[string]string),
		byKey: make(map[string]string),
	}
}

func (m* MemoryRepo) Save(url, short string) error{
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.byURL[url]; ok{
		return nil
	}

	m.byURL[url] = short
	m.byKey[short] = url

	return nil
}

func (m *MemoryRepo) GetByURL(url string) (string, error){
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.byURL[url], nil
}

func (m *MemoryRepo) GetByShortURL(short string) (string, error){
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.byKey[short], nil
}

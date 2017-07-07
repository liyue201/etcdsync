package etcdsync

import (
	"sync"
)

type MutexFactory struct {
	sync.Mutex
	machines   []string
	mutexMap   map[string]*Mutex
	mutexCount map[string]int
}

func NewMutexFactory(machines []string) *MutexFactory {
	return &MutexFactory{
		machines:   machines,
		mutexMap:   make(map[string]*Mutex, 0),
		mutexCount: make(map[string]int, 0),
	}
}

func (factory *MutexFactory) GetMutex(key string, ttl int) *Mutex {
	factory.Lock()
	defer factory.Unlock()

	m, ok := factory.mutexMap[key]
	if ok {
		factory.mutexCount[key]++
		return m
	}

	m = New(key, ttl, factory.machines)
	factory.mutexMap[key] = m
	factory.mutexCount[key] = 1
	return m
}

func (factory *MutexFactory) ReleaseMutex(m *Mutex) {
	factory.Lock()
	defer factory.Unlock()

	count := factory.mutexCount[m.key]
	count--

	if count == 0 {
		delete(factory.mutexCount, m.key)
		delete(factory.mutexMap, m.key)
	} else {
		factory.mutexCount[m.key] = count
	}
}

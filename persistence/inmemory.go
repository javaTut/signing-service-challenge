package persistence

import (
	"signing-service-challenge/domain"
	"sync"
)

// // TODO: in-memory persistence ...

// Simple thread-safe in-memory data storage based on map data structure.
// Thread safety is achieved through Mutex locking.
// RWMutex was used instead of Mutex with the assumption that
// there will be more read than write operations.
// For the sake of simplicity locking logic is done in repositories.

type InMemoryDB struct {
	Devices        map[string]domain.SignatureDevice
	Signatures     map[string]domain.Signature
	DevicesLock    *sync.RWMutex
	SignaturesLock *sync.RWMutex
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		Devices:        make(map[string]domain.SignatureDevice),
		Signatures:     make(map[string]domain.Signature),
		DevicesLock:    &sync.RWMutex{},
		SignaturesLock: &sync.RWMutex{},
	}
}

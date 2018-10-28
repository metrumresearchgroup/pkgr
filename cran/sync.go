package cran

import (
	"sync"
)

// SyncMap provides a mechanism to use a map in a concurrent context
type SyncMap struct {
	Mutex sync.RWMutex // Read Write Mutex to allow for multiple readers
	Map   map[string]Download
}

// NewSyncMap initializes a sync map
func NewSyncMap() *SyncMap {
	return &SyncMap{Mutex: sync.RWMutex{}, Map: make(map[string]Download)}
}

// Put a value into the map, acquiring a lock for concurrent use
func (sm *SyncMap) Put(key string, value Download) {
	sm.Mutex.Lock()
	sm.Map[key] = value
	sm.Mutex.Unlock()
}

// Delete deletes a value from the map
func (sm *SyncMap) Delete(key string) {
	sm.Mutex.Lock()
	delete(sm.Map, key)
	sm.Mutex.Unlock()
}

// Get gets the value associated with a key, and whether it existed
func (sm *SyncMap) Get(key string) (Download, bool) {
	// Reads can be concurrent so this is a read lock and should provide for great read performance
	sm.Mutex.RLock()
	v, ok := sm.Map[key]
	sm.Mutex.RUnlock()
	return v, ok
}

package cran

import (
	"sync"
)

// PkgMap provides a mechanism to use a map in a concurrent context
type PkgMap struct {
	Mutex sync.RWMutex // Read Write Mutex to allow for multiple readers
	Map   map[string]Download
}

// NewPkgMap initializes a sync map
func NewPkgMap() *PkgMap {
	return &PkgMap{Mutex: sync.RWMutex{}, Map: make(map[string]Download)}
}

// Put a value into the map, acquiring a lock for concurrent use
func (sm *PkgMap) Put(key string, value Download) {
	sm.Mutex.Lock()
	sm.Map[key] = value
	sm.Mutex.Unlock()
}

// Delete deletes a value from the map
func (sm *PkgMap) Delete(key string) {
	sm.Mutex.Lock()
	delete(sm.Map, key)
	sm.Mutex.Unlock()
}

// Get gets the value associated with a key, and whether it existed
func (sm *PkgMap) Get(key string) (Download, bool) {
	// Reads can be concurrent so this is a read lock and should provide for great read performance
	sm.Mutex.RLock()
	v, ok := sm.Map[key]
	sm.Mutex.RUnlock()
	return v, ok
}

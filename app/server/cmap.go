package server

import "sync"

// ClientsMap thread-safe map of Clients
type ClientsMap struct {
	mutex sync.RWMutex
	data  map[string]*Client
}

// NewClientsMap creates a ClientsMap
func NewClientsMap() ClientsMap {
	return ClientsMap{
		mutex: sync.RWMutex{},
		data:  make(map[string]*Client),
	}
}

// Add store client to map
func (cm *ClientsMap) Add(c *Client) {
	defer cm.mutex.Unlock()
	cm.mutex.Lock()
	cm.data[c.id] = c
}

// Get fetch client from map by its unique id
func (cm *ClientsMap) Get(cid string) (*Client, bool) {
	defer cm.mutex.RUnlock()
	cm.mutex.RLock()
	r, ok := cm.data[cid]
	return r, ok
}

// Delete client from map
func (cm *ClientsMap) Delete(cid string) {
	defer cm.mutex.Unlock()
	cm.mutex.Lock()
	delete(cm.data, cid)
}

package server

import "sync"

type ClientsMap struct {
	mutex sync.Mutex
	data  map[string]*Client
}

func NewClientsMap() ClientsMap {
	return ClientsMap{
		mutex: sync.Mutex{},
		data:  make(map[string]*Client),
	}
}

func (cm *ClientsMap) Add(c *Client) {
	defer cm.mutex.Unlock()
	cm.mutex.Lock()
	cm.data[c.id] = c
}

func (cm *ClientsMap) Get(cid string) (*Client, bool) {
	defer cm.mutex.Unlock()
	cm.mutex.Lock()
	r, ok := cm.data[cid]
	return r, ok
}

func (cm *ClientsMap) Delete(cid string) {
	defer cm.mutex.Unlock()
	cm.mutex.Lock()
	delete(cm.data, cid)
}

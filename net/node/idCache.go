package node

import (
	"DNA/common"
	"DNA/net/protocol"
	"sync"
)

type idCache struct {
	sync.RWMutex
	lastid    common.Uint256
	index     int
	idarray   []common.Uint256
	idmaplsit map[common.Uint256]int
}

func (c *idCache) init() {
	c.index = 0
	c.idmaplsit = make(map[common.Uint256]int, protocol.MAXIDCACHED)
	c.idarray = make([]common.Uint256, protocol.MAXIDCACHED)
}

func (c *idCache) add(id common.Uint256) {
	oldid := c.idarray[c.index]
	delete(c.idmaplsit, oldid)
	c.idarray[c.index] = id
	c.idmaplsit[id] = c.index
	c.index++
	c.lastid = id
	c.index = c.index % protocol.MAXIDCACHED

}

func (c *idCache) ExistedID(id common.Uint256) bool {
	c.Lock()
	defer c.Unlock()
	if id == c.lastid {
		return true
	}
	if _, ok := c.idmaplsit[id]; ok {
		return true
	} else {
		c.add(id)
	}
	return false
}

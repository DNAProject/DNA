package node

import (
	"DNA/common/log"
	. "DNA/net/protocol"
	"math/rand"
	"sync"
)

const (
	// needAddressThreshold is the number of addresses under which the
	// address manager will claim to need more addresses.
	needAddressThreshold = 1000
)

type KnownAddress struct {
	srcAddr NodeAddr
}

type KnownAddressList struct {
	sync.RWMutex
	List      map[uint64]*KnownAddress
	addrCount uint64
}

func (ka *KnownAddress) SaveAddr(na NodeAddr) {
	ka.srcAddr.Time = na.Time
	ka.srcAddr.Services = na.Services
	ka.srcAddr.IpAddr = na.IpAddr
	ka.srcAddr.Port = na.Port
	ka.srcAddr.ID = na.ID
}

func (ka *KnownAddress) NetAddress() NodeAddr {
	return ka.srcAddr
}

func (ka *KnownAddress) GetID() uint64 {
	return ka.srcAddr.ID
}

func (al *KnownAddressList) NeedMoreAddresses() bool {
	al.Lock()
	defer al.Unlock()

	return al.addrCount < needAddressThreshold
}

func (al *KnownAddressList) AddressExisted(uid uint64) bool {
	_, ok := al.List[uid]
	return ok
}

func (al *KnownAddressList) AddAddressToKnownAddress(na NodeAddr) {
	al.Lock()
	defer al.Unlock()

	ka := new(KnownAddress)
	ka.SaveAddr(na)
	if al.AddressExisted(ka.GetID()) {
		log.Debug("Insert a existed addr\n")
	} else {
		al.List[ka.GetID()] = ka
		al.addrCount++
	}
}

func (al *KnownAddressList) DelAddressFromList(id uint64) bool {
	al.Lock()
	defer al.Unlock()

	_, ok := al.List[id]
	if ok == false {
		return false
	}
	delete(al.List, id)
	return true
}

func (al *KnownAddressList) GetAddressCnt() uint64 {
	al.RLock()
	defer al.RUnlock()
	if al != nil {
		return al.addrCount
	}
	return 0
}

func (al *KnownAddressList) init() {
	al.List = make(map[uint64]*KnownAddress)
}

func isInNbrList(id uint64, nbrAddrs []NodeAddr) bool {
	for _, na := range nbrAddrs {
		if id == na.ID {
			return true
		}
	}
	return false
}

func (al *KnownAddressList) RandGetAddresses(nbrAddrs []NodeAddr) []NodeAddr {
	al.RLock()
	defer al.RUnlock()
	var keys []uint64
	for k := range al.List {
		isInNbr := isInNbrList(k, nbrAddrs)
		if isInNbr == false {
			keys = append(keys, k)
		}
	}
	addrLen := len(keys)
	var i int
	addrs := []NodeAddr{}
	if MAXOUTBOUNDCNT-len(nbrAddrs) > addrLen {
		for _, v := range keys {
			ka, ok := al.List[v]
			if !ok {
				continue
			}
			addrs = append(addrs, ka.srcAddr)
		}
	} else {
		order := rand.Perm(addrLen)
		var count int
		count = MAXOUTBOUNDCNT - len(nbrAddrs)
		for i = 0; i < count; i++ {
			for j, v := range keys {
				if j == order[j] {
					ka, ok := al.List[v]
					if !ok {
						continue
					}
					addrs = append(addrs, ka.srcAddr)
					keys = append(keys[:j], keys[j+1:]...)
					break
				}
			}
		}
	}

	return addrs
}

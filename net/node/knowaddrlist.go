package node

import (
	"DNA/common/log"
	. "DNA/net/protocol"
	"math/rand"
	"sync"
	"time"
)

const (
	// needAddressThreshold is the number of addresses under which the
	// address manager will claim to need more addresses.
	needAddressThreshold = 1000
	// numMissingDays is the number of days before which we assume an
	// address has vanished if we have not seen it announced  in that long.
	numMissingDays = 30
	// numRetries is the number of tried without a single success before
	// we assume an address is bad.
	numRetries = 10
)

type KnownAddress struct {
	srcAddr        NodeAddr
	lastattempt    time.Time
	lastDisconnect time.Time
	attempts       int
}

type KnownAddressList struct {
	sync.RWMutex
	List      map[uint64]*KnownAddress
	addrCount uint64
}

func (ka *KnownAddress) LastAttempt() time.Time {
	return ka.lastattempt
}

func (ka *KnownAddress) increaseAttempts() {
	ka.attempts++
}

func (ka *KnownAddress) updateLastAttempt() {
	// set last tried time to now
	ka.lastattempt = time.Now()
}

func (ka *KnownAddress) updateLastDisconnect() {
	// set last disconnect time to now
	ka.lastDisconnect = time.Now()
}

// chance returns the selection probability for a known address.  The priority
// depends upon how recently the address has been seen, how recently it was last
// attempted and how often attempts to connect to it have failed.
func (ka *KnownAddress) chance() float64 {
	now := time.Now()
	lastAttempt := now.Sub(ka.lastattempt)

	if lastAttempt < 0 {
		lastAttempt = 0
	}

	c := 1.0

	// Very recent attempts are less likely to be retried.
	if lastAttempt < 10*time.Minute {
		c *= 0.01
	}

	// Failed attempts deprioritise.
	for i := ka.attempts; i > 0; i-- {
		c /= 1.5
	}

	return c
}

// isBad returns true if the address in question has not been tried in the last
// minute and meets one of the following criteria:
// 1) Just tried in one minute
// 2) It hasn't been seen in over a month
// 3) It has failed at least ten times
// 4) It has failed ten times in the last week
// All addresses that meet these criteria are assumed to be worthless and not
// worth keeping hold of.
func (ka *KnownAddress) isBad() bool {
	// just tried in one minute?
	if ka.lastattempt.After(time.Now().Add(-1 * time.Minute)) {
		return false
	}

	// Over a month old?
	if ka.srcAddr.Time < (time.Now().Add(-1 * numMissingDays * time.Hour * 24)).UnixNano() {
		return true
	}

	// Just disconnected in one minute?
	if ka.lastDisconnect.After(time.Now().Add(-1 * time.Minute)) {
		return true
	}

	// tried too many times?
	if ka.attempts >= numRetries {
		return true
	}

	return false
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

func (al *KnownAddressList) UpdateLastDisconn(id uint64) {
	ka := al.List[id]
	ka.updateLastDisconnect()
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
		isBad := al.List[k].isBad()
		if isInNbr == false && isBad == false {
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
					ka.increaseAttempts()
					ka.updateLastAttempt()
					addrs = append(addrs, ka.srcAddr)
					keys = append(keys[:j], keys[j+1:]...)
					break
				}
			}
		}
	}

	return addrs
}

func (al *KnownAddressList) RandSelectAddresses() []NodeAddr {
	al.RLock()
	defer al.RUnlock()
	var keys []uint64
	addrs := []NodeAddr{}
	for k := range al.List {
		keys = append(keys, k)
	}
	addrLen := len(keys)

	var count int
	if MAXOUTBOUNDCNT > addrLen {
		count = addrLen
	} else {
		count = MAXOUTBOUNDCNT
	}
	for i, v := range keys {
		if i < count {
			ka, ok := al.List[v]
			if !ok {
				continue
			}
			addrs = append(addrs, ka.srcAddr)
		}
	}

	return addrs
}

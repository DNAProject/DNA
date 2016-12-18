package net

import (
	"log"
	"strconv"
	"net"
	"sync"
	"time"
	"runtime"
	"sync/atomic"
	"GoOnchain/common"
)

// The node state
const (
	INIT = 0
	HANDSHAKEING = 1
	HANDSHAKED = 2
	ESTABLISH = 3
	INACTIVITY = 4
)

// The node capability flag
const (
	RELAY  = 0x01
	SERVER = 0x02
)

type node struct {
	state		uint		// node status
	id		string		// The nodes's id, MAC or IP?
	addr		string 		// The address of the node
	conn		net.Conn	// Connect socket with the peer node
	cap		uint32  	// The node capability set
	handshakeRetry  uint32		// Handshake retry times
	handshakeTime	time.Time	// Last Handshake trigger time
	height		uint64		// The node latest block height
	time		time.Time	// The latest time the node activity
	// TODO does this channel should be a buffer channel
	chF		chan func()	// Channel used to operate the node without lock
	private		*uint		// Reserver for future using
}

type nodeMap struct {
	node *node
	lock sync.RWMutex
	list map[string]*node
}

var nodes nodeMap

func newNode() (*node) {
	node := node{
		state: INIT,
		chF: make(chan func()),
	}

	runtime.SetFinalizer(&node, rmNode)
	go node.backend()
	return &node
}

func rmNode(node *node) {
	log.Printf("Remove node %s", node.addr)
}

// TODO pass pointer to method only need modify it
func (node *node) backend() {
	common.Trace()
	for f := range node.chF {
		f()
	}
}

func (node *node) getID() string {
	return node.id
}

func (node *node) getState() uint {
	return node.state
}

func (node *node) setState(state uint) {
	node.state = state
}

func (node *node) getHandshakeTime() (time.Time) {
	return node.handshakeTime
}

func (node *node) setHandshakeTime(t time.Time) {
	node.handshakeTime = t
}

func (node *node) getHandshakeRetry() uint32 {
	return atomic.LoadUint32(&(node.handshakeRetry))
}

func (node *node) setHandshakeRetry(r uint32) {
	node.handshakeRetry = r
	atomic.StoreUint32(&(node.handshakeRetry), r)
}

func (node *node) updateTime(t time.Time) {
	node.time = t
}

func (node *node) rx() {
	// TODO using select instead of for loop
	for {
		buf := make([]byte, MAXBUFLEN)
		len, err := node.conn.Read(buf)
		if err != nil {
			log.Println("Error reading", err.Error())
			return
		}

		msg := new(Msg)
		err = msg.deserialization(buf)
		if err != nil {
			log.Println("Deserilization buf to message failure")
			return
		}
		log.Printf("Received data: %v", string(buf[:len]))
		go handleNodeMsg(node, msg)
	}
}

// Init the server port, should be run in another thread
func (node *node) initRx () {
	listener, err := net.Listen("tcp", "localhost:" + strconv.Itoa(NODETESTPORT))
	if err != nil {
		log.Println("Error listening", err.Error())
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting", err.Error())
			return
		}
		node := newNode()
		// Currently we use the address as the ID
		node.id = conn.RemoteAddr().String()
		node.addr = conn.RemoteAddr().String()
		log.Println("Remote node %s connect with %s",
			conn.RemoteAddr(), conn.LocalAddr())
		node.conn = conn
		// TOOD close the conn when erro happened
		// TODO lock the node and assign the connection to Node.
		nodes.add(node)
		go node.rx()
	}
	//TODO When to free the net listen resouce?
}

func (node *node) connect(nodeAddr string)  {
	node.chF <- func() {
		common.Trace()
		conn, err := net.Dial("tcp", nodeAddr)
		if err != nil {
			log.Println("Error dialing", err.Error())
			return
		}

		node := newNode()
		node.conn = conn
		node.id = conn.RemoteAddr().String()
		node.addr = conn.RemoteAddr().String()

		log.Printf("Connect node %s connect with %s with %s",
			conn.LocalAddr().String(), conn.RemoteAddr().String(),
			conn.RemoteAddr().Network())
		// TODO Need lock
		nodes.add(node)
		go node.rx()
	}
}

func (node node) tx(msg *Msg) {
	node.chF <- func() {
		buf, err := msg.serialization()
		if (err != nil) {
			log.Println("Error Convert net message ", err.Error())
			return
		}
		_, err = node.conn.Write(buf)
		if err != nil {
			log.Println("Error sending messge to peer node", err.Error())
		}
		return
	}
}

func (nodes *nodeMap) broadcast(msg *Msg) {
	// TODO lock the map
	// TODO Check whether the node existed or not
	for _, node := range nodes.list {
		if node.state == ESTABLISH {
			go node.tx(msg)
		}
	}
}

func (nodes *nodeMap) add(node *node) {
	//TODO lock the node Map
	// TODO check whether the node existed or not
	// TODO dupicate IP address nodes issue
	nodes.list[node.id] = node
	// Unlock the map
}

func (nodes *nodeMap) delNode(node *node) {
	//TODO lock the node Map
	delete(nodes.list, node.id)
	// Unlock the map
}

func InitNodes() {
	nodes.node = newNode()
	nodes.list = make(map[string]*node)
}

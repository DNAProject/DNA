package message

import (
	"GoOnchain/common"
	"GoOnchain/common/log"
	"GoOnchain/core/ledger"
	"GoOnchain/events"
	. "GoOnchain/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
	"sync"
)

var DefaultBlockPool BlockPool
type BlockPool []ledger.Block

type blockReq struct {
	msgHdr
	//TBD
}

type block struct {
	msgHdr
	blk ledger.Block
	mu           sync.Mutex
	// TBD
	//event *events.Event
}

func (msg block) Handle(node Noder) error {
	common.Trace()
	msg.mu.Lock()
	defer msg.mu.Unlock()
	log.Debug("RX block message and the block height is ",msg.blk.Blockdata.Height)
	if msg.blk.Blockdata.Height == ledger.DefaultLedger.Blockchain.BlockHeight + 1{
		err := ledger.DefaultLedger.Blockchain.AddBlock(&msg.blk)
		if (err != nil) {
			log.Warn("Add block error")
			return errors.New("Add block error before Xmit\n")
		}
		DefaultBlockPool.CheckAndAddBlockFromPool(msg.blk.Blockdata.Height+1)
		node.LocalNode().GetEvent("block").Notify(events.EventNewInventory, &msg.blk)
	}else {
		if msg.blk.Blockdata.Height > (ledger.DefaultLedger.Blockchain.BlockHeight + 1){
			msg.AddBlockToPool()
		}
	}
	return nil
}

func (msg *block) AddBlockToPool() error {
	if exist:=msg.CheckBlockPoolIsExist(); !exist{
		DefaultBlockPool = append(DefaultBlockPool,msg.blk)
	}
	sort.Sort(BlockPool(DefaultBlockPool))
	return nil
}

func (msg *block) CheckBlockPoolIsExist() bool{
	for _, v := range DefaultBlockPool {
	    if v.Blockdata.Height == msg.blk.Blockdata.Height {
		    return  true
	    }
	}
	return false
}

func (bp *BlockPool) CheckAndAddBlockFromPool(height uint32) error {
	for _, v := range *bp {
		if v.Blockdata.Height==height{
			err := ledger.DefaultLedger.Blockchain.AddBlock(&v)
			if (err != nil) {
				log.Warn("Add block error and blockheight is ",v.Blockdata.Height)
				return errors.New("Add block error from BlockPool\n")
			}
		}
		height++
	}

	return nil
}

func (b BlockPool) Len() int           { return len(b) }
func (b BlockPool) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b BlockPool) Less(i, j int) bool { return b[i].Blockdata.Height < b[j].Blockdata.Height}


func (msg dataReq) Handle(node Noder) error {
	common.Trace()
	reqtype := msg.dataType
	hash := msg.hash
	switch reqtype {
	case 0x01:
		block := NewBlockFromHash(hash)
		buf, _ := NewBlock(block)
		go node.Tx(buf)

	case 0x02:
		tx := NewTxFromHash(hash)
		buf, _ := NewTx(tx)
		go node.Tx(buf)
	}
	return nil
}

func NewBlockFromHash(hash common.Uint256) *ledger.Block {
	bk, _ := ledger.DefaultLedger.Store.GetBlock(hash)
	return bk
}

func NewBlock(bk *ledger.Block) ([]byte, error) {
	common.Trace()
	var msg block
	msg.blk = *bk
	msg.msgHdr.Magic = NETMAGIC
	cmd := "block"
	copy(msg.msgHdr.CMD[0:len(cmd)], cmd)
	tmpBuffer := bytes.NewBuffer([]byte{})
	bk.Serialize(tmpBuffer)
	p := new(bytes.Buffer)
	err := binary.Write(p, binary.LittleEndian, tmpBuffer.Bytes())
	if err != nil {
		fmt.Println("Binary Write failed at new Msg")
		return nil, err
	}
	s := sha256.Sum256(p.Bytes())
	s2 := s[:]
	s = sha256.Sum256(s2)
	buf := bytes.NewBuffer(s[:4])
	binary.Read(buf, binary.LittleEndian, &(msg.msgHdr.Checksum))
	msg.msgHdr.Length = uint32(len(p.Bytes()))
	fmt.Printf("The message payload length is %d\n", msg.msgHdr.Length)

	m, err := msg.Serialization()
	if err != nil {
		fmt.Println("Error Convert net message ", err.Error())
		return nil, err
	}

	return m, nil
}

func reqBlkData(node Noder, hash common.Uint256) error {
	var msg dataReq
	msg.dataType = common.BLOCK
	// TODO handle the hash array case
	msg.hash = hash

	buf, _ := msg.Serialization()
	go node.Tx(buf)

	return nil
}

func (msg block) Verify(buf []byte) error {
	err := msg.msgHdr.Verify(buf)
	// TODO verify the message Content
	return err
}

func (msg block) Serialization() ([]byte, error) {
	hdrBuf, err := msg.msgHdr.Serialization()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(hdrBuf)
	msg.blk.Serialize(buf)

	return buf.Bytes(), err
}

func (msg *block) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)

	err := binary.Read(buf, binary.LittleEndian, &(msg.msgHdr))
	if err != nil {
		log.Warn("Parse block message hdr error")
		return errors.New("Parse block message hdr error")
	}

	err = msg.blk.Deserialize(buf)
	if err != nil {
		log.Warn("Parse block message error")
		return errors.New("Parse block message error")
	}

	return err
}

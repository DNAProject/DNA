package ledger

import (
	. "GoOnchain/common"
	"GoOnchain/common/log"
	tx "GoOnchain/core/transaction"
	"GoOnchain/crypto"
	. "GoOnchain/errors"
	"GoOnchain/events"
	"sync"
)

type Blockchain struct {
	//BlockCache  map[Uint256]*Block
	GenesisBlock *Block
	BlockHeight uint32
	BCEvents    *events.Event
	mutex       sync.Mutex
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		BlockHeight: 0,
		//BlockCache: make(map[Uint256]*Block),
		BCEvents:   events.NewEvent(),
	}
}

func (bc *Blockchain) AddBlock(block *Block) error {
	Trace()
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	err := bc.SaveBlock(block)
	if err != nil {
		return err
	}

	// Need atomic oepratoion
	bc.BlockHeight = bc.BlockHeight +1

	return nil
}

//func (bc *Blockchain) AddBlockCache(block *Block) {
//	bc.mutex.Lock()
//	defer bc.mutex.Unlock()
//	if _, ok := bc.BlockCache[block.Hash()]; !ok {
//		bc.BlockCache[block.Hash()] = block
//	}
//}

func (bc *Blockchain) ContainsBlock(hash Uint256) bool {
	//TODO: implement ContainsBlock
	Trace()
	log.Info("bc.GenesisBlock.Hash()=",bc.GenesisBlock.Hash())
	log.Info("Blockhash=",hash)
	if hash == bc.GenesisBlock.Hash(){
		return true
	}
	return false
}

func (bc *Blockchain) GetHeader(hash Uint256) (*Header,error) {
	 header,err:=DefaultLedger.Store.GetHeader(hash)
	if err != nil{
		return nil, NewDetailErr(err, ErrNoCode, "[Blockchain], GetHeader failed.")
	}
	return header,nil
}

func (bc *Blockchain) SaveBlock(block *Block) error {
	Trace()
	err := DefaultLedger.Store.SaveBlock(block,DefaultLedger)
	if err != nil {
		log.Error("Save block failure ,err=",err)
		return err
	}
	bc.BCEvents.Notify(events.EventBlockPersistCompleted, block)

	return nil
}

func (bc *Blockchain) ContainsTransaction(hash Uint256) bool {
	//TODO: implement error catch
	_ ,err := DefaultLedger.Store.GetTransaction(hash)
	if (err!= nil){
		return false
	}
	return true
}

func (bc *Blockchain) GetMinersByTXs(others []*tx.Transaction) []*crypto.PubKey {
	//TODO: GetMiners()
	//TODO: Just for TestUse

	return StandbyMiners
}

func (bc *Blockchain) GetMiners() []*crypto.PubKey {
	//TODO: GetMiners()
	//TODO: Just for TestUse

	return StandbyMiners
}

func (bc *Blockchain) CurrentBlockHash() Uint256 {
	return DefaultLedger.Store.GetCurrentBlockHash()
}


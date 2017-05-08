package ledger

import (
	. "DNA/common"
	"DNA/common/log"
	tx "DNA/core/transaction"
	"DNA/crypto"
	. "DNA/errors"
	"DNA/events"
	"fmt"
	"sync"
)

type Blockchain struct {
	BlockHeight uint32
	BCEvents    *events.Event
	mutex       sync.Mutex

	blockSaveCompletedSubscriber events.Subscriber
}

func NewBlockchain(height uint32) *Blockchain {
	return &Blockchain{
		BlockHeight: height,
		BCEvents:    events.NewEvent(),
	}
}

func NewBlockchainWithGenesisBlock() (*Blockchain, error) {
	genesisBlock, err := GenesisBlockInit()
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Blockchain], NewBlockchainWithGenesisBlock failed.")
	}
	genesisBlock.RebuildMerkleRoot()
	hashx := genesisBlock.Hash()
	genesisBlock.hash = &hashx

	height, err := DefaultLedger.Store.InitLevelDBStoreWithGenesisBlock(genesisBlock)
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Blockchain], InitLevelDBStoreWithGenesisBlock failed.")
	}
	blockchain := NewBlockchain(height)
	blockchain.blockSaveCompletedSubscriber = blockchain.BCEvents.Subscribe(events.EventBlockSaveCompleted, blockchain.BlockSaveCompleted)

	return blockchain, nil
}

func (bc *Blockchain) BlockSaveCompleted(v interface{}) {
	log.Debug()
	if block, ok := v.(*Block); ok {
		bc.BCEvents.Notify(events.EventBlockPersistCompleted, block)
		log.Info(fmt.Sprintf("[BlockSaveCompleted] persist block: %d", block.Hash()))
	}
}

func (bc *Blockchain) AddBlock(block *Block) bool {
	log.Debug()
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if !bc.SaveBlock(block) {
		return false
	}

	// Need atomic oepratoion
	bc.BlockHeight = bc.BlockHeight + 1

	return true
}

//
//func (bc *Blockchain) ContainsBlock(hash Uint256) bool {
//	//TODO: implement ContainsBlock
//	if hash == bc.GenesisBlock.Hash(){
//		return true
//	}
//	return false
//}

func (bc *Blockchain) GetHeader(hash Uint256) (*Header, error) {
	header, err := DefaultLedger.Store.GetHeader(hash)
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Blockchain], GetHeader failed.")
	}
	return header, nil
}

func (bc *Blockchain) SaveBlock(block *Block) bool {
	log.Debug()
	log.Info("block hash ", block.Hash())
	if !DefaultLedger.Store.SaveBlock(block, DefaultLedger) {
		log.Warn("Save block failure.")
		return false
	}

	return true
}

func (bc *Blockchain) ContainsTransaction(hash Uint256) bool {
	//TODO: implement error catch
	_, err := DefaultLedger.Store.GetTransaction(hash)
	if err != nil {
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

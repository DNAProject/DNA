package ledger

import (
	. "DNA/common"
	"DNA/common/log"
	tx "DNA/core/transaction"
	"DNA/crypto"
	. "DNA/errors"
	"DNA/events"
	"sync"
)

type Blockchain struct {
	BlockHeight uint32
	BCEvents    *events.Event
	mutex       sync.Mutex
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		BlockHeight: 0,
		BCEvents:    events.NewEvent(),
	}
}

func NewBlockchainWithGenesisBlock() (*Blockchain, error) {
	blockchain := NewBlockchain()
	genesisBlock, err := GenesisBlockInit()
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Blockchain], NewBlockchainWithGenesisBlock failed.")
	}
	genesisBlock.RebuildMerkleRoot()
	hashx := genesisBlock.Hash()
	genesisBlock.hash = &hashx
	//blockchain.AddBlock(genesisBlock)
	DefaultLedger.Store.InitLevelDBStoreWithGenesisBlock(genesisBlock)
	return blockchain, nil
}

func (bc *Blockchain) AddBlock(block *Block) error {
	log.Trace()
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	err := bc.SaveBlock(block)
	if err != nil {
		return err
	}

	// Need atomic oepratoion
	bc.BlockHeight = bc.BlockHeight + 1

	return nil
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

func (bc *Blockchain) SaveBlock(block *Block) error {
	log.Trace()
	log.Info("block hash ", block.Hash())
	err := DefaultLedger.Store.SaveBlock(block, DefaultLedger)
	if err != nil {
		log.Warn("Save block failure ,err= ", err)
		return err
	}

	return nil
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

package ledger

import (
	. "GoOnchain/common"
	"GoOnchain/common/log"
	tx "GoOnchain/core/transaction"
	"GoOnchain/crypto"
	. "GoOnchain/errors"
	"GoOnchain/events"
	"errors"
	"sync"
	"fmt"
	"bytes"
	"sort"
)

type BlockPool []Block

type Blockchain struct {
	BlockCache  BlockPool
	BlockHeight uint32
	BCEvents    *events.Event
	mutex       sync.Mutex
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		BlockHeight: 0,
		BlockCache: BlockPool{},
		BCEvents:   events.NewEvent(),
	}
}

func NewBlockchainWithGenesisBlock() (*Blockchain,error) {
	blockchain := NewBlockchain()
	genesisBlock,err:=GenesisBlockInit()
	if err != nil{
		return nil,NewDetailErr(err, ErrNoCode, "[Blockchain], NewBlockchainWithGenesisBlock failed.")
	}
	genesisBlock.RebuildMerkleRoot()
	hashx :=genesisBlock.Hash()
	genesisBlock.hash = &hashx
	blockchain.SaveBlock(genesisBlock)
	return blockchain,nil
}

func (bc *Blockchain) AddBlock(block *Block) error {
	Trace()

	//set block cache
	bc.AddBlockCache(block)

	//Block header verfiy
	if ok:=bc.BlockAddVerifyOK(block);ok{
		//save block
		buf := bytes.NewBuffer([]byte{})
		block.Serialize(buf)
		fmt.Printf("***Blockchain Height %d,AddBlock detail %d\n",bc.BlockHeight,buf.Bytes())
		err := bc.SaveBlock(block)
		if err != nil {
			return err
		}
		bc.BlockHeight = bc.BlockHeight+1
		bc.BlockCache.CheckAndAddBlockFromPool(bc.BlockHeight)

	}
	return nil
}

func (bc *Blockchain) AddBlockCache(block *Block) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	bc.BlockCache.AddBlockToPool(block)
}

func (bc *Blockchain) ContainsBlock(hash Uint256) bool {
	_,err:=DefaultLedger.GetBlockWithHash(hash)
	if err!= nil{
		return false
	}
	return true
}

func (bc *Blockchain) GetHeader(hash Uint256) (*Header,error) {
	 header,err:=DefaultLedger.Store.GetHeader(hash)
	if err != nil{
		return nil, NewDetailErr(errors.New("[Blockchain], GetHeader failed."), ErrNoCode, "")
	}
	return header,nil
}

func (bc *Blockchain) SaveBlock(block *Block) error {
	Trace()
	err := DefaultLedger.Store.SaveBlock(block)
	if err != nil {
		log.Error("Save block failure")
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

func (bc *Blockchain) BlockAddVerifyOK(block *Block) bool{
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	fmt.Println(block.Blockdata.Height)
	fmt.Println(DefaultLedger.Blockchain.BlockHeight)
	if block.Blockdata.Height != DefaultLedger.Blockchain.BlockHeight +1{
		return false
	}
	if ok:=bc.ContainsBlock(block.Hash());ok{
		return false
	}
	return true
}

func (bp *BlockPool) CheckAndAddBlockFromPool(height uint32) error {
	for _, v := range *bp {
		if v.Blockdata.Height > height{
			return nil
		}
		if v.Blockdata.Height==height{
			err := DefaultLedger.Blockchain.AddBlock(&v)
			if (err != nil) {
				log.Warn("Add block error and blockheight is ",v.Blockdata.Height)
				return errors.New("Add block error from BlockPool\n")
			}
		}
		height++
		DefaultLedger.Blockchain.BlockHeight = DefaultLedger.Blockchain.BlockHeight + 1
	}
	return nil
}

func (b BlockPool) Len() int           { return len(b) }
func (b BlockPool) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b BlockPool) Less(i, j int) bool { return b[i].Blockdata.Height < b[j].Blockdata.Height}

func (bp *BlockPool) CheckBlockPoolIsExist(bk *Block) bool{
	for _, v := range *bp {
		if v.Blockdata.Height == bk.Blockdata.Height {
			return  true
		}
	}
	return false
}

func (bp *BlockPool) AddBlockToPool(bk *Block) error {
	if exist:=bp.CheckBlockPoolIsExist(bk); !exist{
		*bp = append(*bp,*bk)
	}
	sort.Sort(BlockPool(*bp))
	return nil
}
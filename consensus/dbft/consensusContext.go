package dbft

import (
	cl "DNA/client"
	. "DNA/common"
	"DNA/common/log"
	ser "DNA/common/serialization"
	"DNA/core/ledger"
	tx "DNA/core/transaction"
	"DNA/crypto"
	"DNA/net"
	msg "DNA/net/message"
	"fmt"
	"sort"
	"sync"
)

const ContextVersion uint32 = 0

type ConsensusContext struct {
	State           ConsensusState
	PrevHash        Uint256
	Height          uint32
	ViewNumber      byte
	BookKeepers     []*crypto.PubKey
	Owner           *crypto.PubKey
	BookKeeperIndex int
	PrimaryIndex    uint32
	Timestamp       uint32
	Nonce           uint64
	NextBookKeeper  Uint160
	Transactions    []*tx.Transaction
	Signatures      [][]byte
	ExpectedView    []byte

	header *ledger.Block

	contextMu sync.Mutex
}

func (cxt *ConsensusContext) M() int {
	log.Debug()
	return len(cxt.BookKeepers) - (len(cxt.BookKeepers)-1)/3
}

func NewConsensusContext() *ConsensusContext {
	log.Debug()
	return &ConsensusContext{}
}

func (cxt *ConsensusContext) ChangeView(viewNum byte) {
	log.Debug()
	p := (cxt.Height - uint32(viewNum)) % uint32(len(cxt.BookKeepers))
	cxt.State &= SignatureSent
	cxt.ViewNumber = viewNum
	if p >= 0 {
		cxt.PrimaryIndex = uint32(p)
	} else {
		cxt.PrimaryIndex = uint32(p) + uint32(len(cxt.BookKeepers))
	}

	if cxt.State == Initial {
		cxt.Transactions = nil
		cxt.Signatures = make([][]byte, len(cxt.BookKeepers))
	}
	cxt.header = nil
}

func (cxt *ConsensusContext) MakeChangeView() *msg.ConsensusPayload {
	log.Debug()
	cv := &ChangeView{
		NewViewNumber: cxt.ExpectedView[cxt.BookKeeperIndex],
	}
	cv.msgData.Type = ChangeViewMsg
	return cxt.MakePayload(cv)
}

func (cxt *ConsensusContext) MakeHeader() *ledger.Block {
	log.Debug()
	if cxt.Transactions == nil {
		return nil
	}
	txHash := []Uint256{}
	for _, t := range cxt.Transactions {
		txHash = append(txHash, t.Hash())
	}
	txRoot, err := crypto.ComputeRoot(txHash)
	if err != nil {
		return nil
	}
	if cxt.header == nil {
		blockData := &ledger.Blockdata{
			Version:          ContextVersion,
			PrevBlockHash:    cxt.PrevHash,
			TransactionsRoot: txRoot,
			Timestamp:        cxt.Timestamp,
			Height:           cxt.Height,
			ConsensusData:    cxt.Nonce,
			NextBookKeeper:   cxt.NextBookKeeper,
		}
		cxt.header = &ledger.Block{
			Blockdata:    blockData,
			Transactions: []*tx.Transaction{},
		}
	}
	return cxt.header
}

func (cxt *ConsensusContext) MakePayload(message ConsensusMessage) *msg.ConsensusPayload {
	log.Debug()
	message.ConsensusMessageData().ViewNumber = cxt.ViewNumber
	return &msg.ConsensusPayload{
		Version:         ContextVersion,
		PrevHash:        cxt.PrevHash,
		Height:          cxt.Height,
		BookKeeperIndex: uint16(cxt.BookKeeperIndex),
		Timestamp:       cxt.Timestamp,
		Data:            ser.ToArray(message),
		Owner:           cxt.Owner,
	}
}

func (cxt *ConsensusContext) MakePrepareRequest() *msg.ConsensusPayload {
	log.Debug()
	preReq := &PrepareRequest{
		Nonce:          cxt.Nonce,
		NextBookKeeper: cxt.NextBookKeeper,
		Transactions:   cxt.Transactions,
		Signature:      cxt.Signatures[cxt.BookKeeperIndex],
	}
	preReq.msgData.Type = PrepareRequestMsg
	return cxt.MakePayload(preReq)
}

func (cxt *ConsensusContext) MakePrepareResponse(signature []byte) *msg.ConsensusPayload {
	log.Debug()
	preRes := &PrepareResponse{
		Signature: signature,
	}
	preRes.msgData.Type = PrepareResponseMsg
	return cxt.MakePayload(preRes)
}

func (cxt *ConsensusContext) GetSignaturesCount() (count int) {
	log.Debug()
	count = 0
	for _, sig := range cxt.Signatures {
		if sig != nil {
			count += 1
		}
	}
	return count
}

func (cxt *ConsensusContext) GetStateDetail() string {

	return fmt.Sprintf("Initial: %t, Primary: %t, Backup: %t, RequestSent: %t, RequestReceived: %t, SignatureSent: %t, BlockSent: %t, ",
		cxt.State.HasFlag(Initial),
		cxt.State.HasFlag(Primary),
		cxt.State.HasFlag(Backup),
		cxt.State.HasFlag(RequestSent),
		cxt.State.HasFlag(RequestReceived),
		cxt.State.HasFlag(SignatureSent),
		cxt.State.HasFlag(BlockSent))

}

func (cxt *ConsensusContext) Reset(client cl.Client, localNode net.Neter) {
	log.Debug()
	cxt.State = Initial
	cxt.PrevHash = ledger.DefaultLedger.Blockchain.CurrentBlockHash()
	cxt.Height = ledger.DefaultLedger.Blockchain.BlockHeight + 1
	cxt.ViewNumber = 0
	cxt.BookKeeperIndex = -1

	bookKeepers, _ := localNode.GetBookKeepersAddrs()
	cxt.Owner = bookKeepers[0]
	sort.Sort(crypto.PubKeySlice(bookKeepers))
	cxt.BookKeepers = bookKeepers

	bookKeeperLen := len(cxt.BookKeepers)
	cxt.PrimaryIndex = cxt.Height % uint32(bookKeeperLen)
	cxt.Transactions = nil
	cxt.Signatures = make([][]byte, bookKeeperLen)
	cxt.ExpectedView = make([]byte, bookKeeperLen)

	for i := 0; i < bookKeeperLen; i++ {
		ac, _ := client.GetDefaultAccount()
		if ac.PublicKey.X.Cmp(cxt.BookKeepers[i].X) == 0 {
			cxt.BookKeeperIndex = i
			break
		}
	}

	cxt.header = nil
}

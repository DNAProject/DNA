package smartcontract

import (
	"DNA/vm"
	"DNA/core/ledger"
	"DNA/common"
	"math/big"
	"errors"
)

type StateReader struct {
	serviceMap map[string]func(*vm.ExecutionEngine) (bool, error)
}

func NewStateReader() *StateReader {
	var stateReader StateReader
	stateReader.serviceMap = make(map[string]func(*vm.ExecutionEngine) (bool, error), 0)
	stateReader.Register("DNA.Blockchain.GetHeight", stateReader.BlockChainGetHeight)
	stateReader.Register("DNA.Blockchain.GetHeader", stateReader.BlockChainGetHeader)
	stateReader.Register("DNA.Blockchain.GetBlock", stateReader.BlockChainGetBlock)
	stateReader.Register("DNA.Blockchain.GetTransaction", stateReader.BlockChainGetTransaction)
	stateReader.Register("DNA.Blockchain.GetAsset", stateReader.BlockChainGetAsset)

	stateReader.Register("DNA.Header.GetHash", stateReader.HeaderGetHash);

	return &stateReader
}


func (s *StateReader) Register(methodName string, handler func(*vm.ExecutionEngine) (bool, error)) bool {
	if _, ok := s.serviceMap[methodName]; ok {
		return false
	}
	s.serviceMap[methodName] = handler
	return true
}

func (s *StateReader) GetServiceMap() map[string]func(*vm.ExecutionEngine) (bool, error) {
	return s.serviceMap
}

func (s *StateReader) BlockChainGetHeight(e *vm.ExecutionEngine) (bool, error) {
	var i *big.Int
	if ledger.DefaultLedger == nil {
		i = big.NewInt(0)
	}else {
		i = big.NewInt(int64(ledger.DefaultLedger.Store.GetHeight()))
	}
	err := vm.PushData(e, i)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StateReader) BlockChainGetHeader(e *vm.ExecutionEngine) (bool, error) {
	data := e.GetEvaluationStack().Pop().GetStackItem().GetByteArray()
	var (
		header *ledger.Header
		err error
	)
	l := len(data)
	if l <= 5 {
		b := new(big.Int)
		height := uint32(b.SetBytes(data).Int64())
		if ledger.DefaultLedger != nil {
			hash, err := ledger.DefaultLedger.Store.GetBlockHash(height)
			if err != nil {
				return false, err
			}
			header, err = ledger.DefaultLedger.Store.GetHeader(hash)
			if err != nil {
				return false, err
			}
		}else {
			header = nil
		}
	}else if l == 32 {
		hash, _ := common.Uint256ParseFromBytes(data)
		if ledger.DefaultLedger != nil {
			header, err = ledger.DefaultLedger.Store.GetHeader(hash)
			if err != nil {
				return false, err
			}
		}else {
			header = nil
		}
	}else {
		return false, errors.New("The data length is error in function blockchaningetheader!")
	}
	err = vm.PushData(e, header)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StateReader) BlockChainGetBlock(e *vm.ExecutionEngine) (bool, error) {
	data := e.GetEvaluationStack().Pop().GetStackItem().GetByteArray()
	var (
		block *ledger.Block
		err error
	)
	l := len(data)
	if l <= 5 {
		b := new(big.Int)
		height := uint32(b.SetBytes(data).Int64())
		if ledger.DefaultLedger != nil {
			hash, err := ledger.DefaultLedger.Store.GetBlockHash(height)
			if err != nil {
				return false, err
			}
			block, err = ledger.DefaultLedger.Store.GetBlock(hash)
			if err != nil {
				return false, err
			}
		}else {
			block = nil
		}
	}else if l == 32 {
		hash, err := common.Uint256ParseFromBytes(data)
		if err != nil {
			return false, err
		}
		if ledger.DefaultLedger != nil {
			block, err = ledger.DefaultLedger.Store.GetBlock(hash)
			if err != nil {
				return false, err
			}
		}else {
			block = nil
		}
	}else {
		return false, errors.New("The data length is error in function blockchaningetblock!")
	}
	err = vm.PushData(e, block)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StateReader) BlockChainGetTransaction(e *vm.ExecutionEngine) (bool, error) {
	data := e.GetEvaluationStack().Pop().GetStackItem().GetByteArray()
	hash, err := common.Uint256ParseFromBytes(data)
	if err != nil {
		return false, err
	}
	tx, err := ledger.DefaultLedger.Store.GetTransaction(hash)
	if err != nil {
		return false, err
	}

	err = vm.PushData(e, tx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StateReader) BlockChainGetAsset(e *vm.ExecutionEngine) (bool, error) {
	data := e.GetEvaluationStack().Pop().GetStackItem().GetByteArray()
	hash, err := common.Uint256ParseFromBytes(data)
	if err != nil {
		return false, err
	}
	asset, err := ledger.DefaultLedger.Store.GetAsset(hash)
	if err != nil {
		return false, err
	}
	err = vm.PushData(e, asset)


	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StateReader) HeaderGetHash(e *vm.ExecutionEngine) (bool, error) {
	data := e.GetEvaluationStack().Pop().GetStackItem().GetInterface()
	if data == nil {
		return false, errors.New("Get stack data error in function headergethash!")
	}
	hash := data.(*ledger.Block).Hash()
	err := vm.PushData(e, hash.ToArray())
	if err != nil {
		return false, err
	}
	return true, nil
}
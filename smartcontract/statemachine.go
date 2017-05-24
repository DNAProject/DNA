package smartcontract

import (
	"DNA/vm"
	"DNA/common"
	"DNA/smartcontract/states"
	. "DNA/errors"
	"DNA/core/store"
	"DNA/smartcontract/storage"
	"fmt"
)

type StorageContext byte

const (
	Current StorageContext = 0x01
	CallingContract StorageContext = 0x02
	EntryContract StorageContext = 0x04
)

type StateMachine struct {
	*StateReader
	RWSet *storage.RWSet
}

func NewStateMachine() *StateMachine {
	var stateMachine StateMachine
	stateMachine.RWSet = storage.NewRWSet()
	stateMachine.StateReader = NewStateReader()
	stateMachine.StateReader.Register("DNA.Storage.Get", stateMachine.StorageGet)
	stateMachine.StateReader.Register("DNA.Storage.Put", stateMachine.StoragePut)
	stateMachine.StateReader.Register("DNA.Storage.Delete", stateMachine.StorageDelete)
	return &stateMachine
}

func(s *StateMachine) CheckStorageContext(engine *vm.ExecutionEngine, context StorageContext) (*common.Uint160, error) {
	var hash []byte
	switch context {
		case Current:
			hash = engine.CurrentContext().GetCodeHash()
		case CallingContract:
			hash = engine.CallingContext().GetCodeHash()
		case EntryContract:
			hash = engine.EntryContext().GetCodeHash()
		default:
			return nil, NewErr("Error StorageContext!")

	}
	if hash == nil || len(hash) == 0 {
		return nil, NewErr("Get code hash from context fail!")
	}
	scriptHash, _ := common.Uint160ParseFromBytes(hash)
	return &scriptHash, nil
}


func(s *StateMachine) StorageGet(engine *vm.ExecutionEngine) (bool, error) {
	codeHash, err := s.GetCodeHash(engine)
	if err != nil { return false, err }
	key := engine.GetEvaluationStack().Pop().GetStackItem().GetByteArray()
	storageKey := states.NewStorageKey(codeHash, key)

	item, err := s.RWSet.TryGet(store.ST_Storage, storage.KeyToStr(storageKey))
	if err != nil {
		return false, err
	}
	i, err := vm.NewStackItemInterface(item.(*states.StorageItem).Value)
	if err != nil {
		return false, err
	}
	engine.GetEvaluationStack().Push(vm.NewStackItem(i))
	return true, nil
}

func(s *StateMachine) StoragePut(engine *vm.ExecutionEngine) (bool, error) {
	codeHash, err := s.GetCodeHash(engine)
	if err != nil { return false, err }
	key := engine.GetEvaluationStack().Pop().GetStackItem().GetByteArray()
	value := engine.GetEvaluationStack().Pop().GetStackItem().GetByteArray()
	storageKey := states.NewStorageKey(codeHash, key)
	s.RWSet.GetOrPut(storage.KeyToStr(storageKey), states.NewStorageItem(value))
	return true, nil
}

func(s *StateMachine) StorageDelete(engine *vm.ExecutionEngine) (bool, error) {
	codeHash, err := s.GetCodeHash(engine)
	if err != nil { return false, err }
	key := engine.GetEvaluationStack().Pop().GetStackItem().GetByteArray()
	storageKey := states.NewStorageKey(codeHash, key)
	s.RWSet.Delete(storage.KeyToStr(storageKey))
	return true, nil
}

func(s *StateMachine) GetCodeHash(engine *vm.ExecutionEngine) (*common.Uint160, error) {
	context := StorageContext(byte(engine.GetEvaluationStack().Pop().GetStackItem().GetBigInteger().Int64()))
	codeHash, err := s.CheckStorageContext(engine, StorageContext(context))
	if err != nil{
		return nil, NewErr(fmt.Sprintf("Get Code Hash Err:%v", err))
	}
	return codeHash, nil
}





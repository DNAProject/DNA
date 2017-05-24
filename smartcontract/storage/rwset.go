package storage

import (
	"DNA/smartcontract/states"
	"bytes"
	"DNA/core/store"
	"DNA/core/ledger"
	"errors"
)

type RWSet struct {
	ReadSet map[string]*ReadSet
	WriteSet map[string]*WriteSet
}

type WriteSet struct {
	Key string
	Item states.IStateValueInterface
	isDeleted bool
}

type ReadSet struct {
	Key states.IStateKeyInterface
	Version string
}

func NewRWSet() *RWSet {
	var rwSet RWSet
	rwSet.WriteSet = make(map[string]*WriteSet, 0)
	rwSet.ReadSet = make(map[string]*ReadSet, 0)
	return &rwSet
}

func(rw *RWSet)Put(key string, value states.IStateValueInterface) {
	rw.WriteSet[key] = &WriteSet{
		Key: key,
		Item: value,
		isDeleted: false,
	}
}

func(rw *RWSet) Delete(key string){
	if _, ok := rw.WriteSet[key]; ok {
		rw.WriteSet[key].isDeleted = true
	}else {
		rw.WriteSet[key] = &WriteSet{
			Key: key,
			Item: nil,
			isDeleted: true,
		}
	}
}

func(rw *RWSet) GetOrPut(key string, value states.IStateValueInterface) states.IStateValueInterface{
	var writeSet *WriteSet
	if v, ok := rw.WriteSet[key]; ok {
		writeSet = v
		if writeSet.isDeleted{
			writeSet.Item = value
			writeSet.isDeleted = false
		}
	}else {
		writeSet = &WriteSet{
			Key: key,
			Item: nil,
			isDeleted: false,
		}
		if writeSet.Item == nil {
			writeSet.Item = value

		}
		rw.WriteSet[key] = writeSet
	}
	return writeSet.Item
}

func(rw *RWSet) TryGet(prefix store.DataEntryPrefix, key string) (states.IStateValueInterface, error){
	if v, ok := rw.WriteSet[key]; ok {
		if v.isDeleted {
			return nil, errors.New("the value is deleted!")
		}
		return  v.Item, nil
	}else {
		write := new(bytes.Buffer)
		var (
			value []byte
			err error
		)
		switch prefix {
			case store.ST_Storage: {
				value, err = ledger.DefaultLedger.Store.GetStorage([]byte(key))
				if err != nil {
					return nil, err
				}
				item := &states.StorageItem{}
				write.Write(value)
				item.Deserialize(write)
				return item, nil
			}
			case store.ST_Contract: {
				value, err = ledger.DefaultLedger.Store.GetContract([]byte(key))
				if err != nil {
					return nil, err
				}
				item := &states.ContractState{}
				write.Write(value)
				item.Deserialize(write)
				return item, nil
			}
			default:
				return nil, errors.New("the store prefix not exist!")

		}
	}
}

func(rw *RWSet) GetChangeSet() map[string]map[string]string {
	w := make(map[string]map[string]string, 0)
	m := make(map[string]string, 0)
	for k, v := range rw.WriteSet {
		value := new(bytes.Buffer)
		v.Item.Serialize(value)
		m[k] = string(value.Bytes())
		switch v.Item.(type) {
			case *states.ContractState: {
				w[string(store.ST_Contract)] = m
			}
			case *states.StorageItem: {
				w[string(store.ST_Storage)] = m
			}
		}
	}
	return w
}


func KeyToStr(key states.IStateKeyInterface) string {
	k := new(bytes.Buffer)
	key.Serialize(k)
	return string(k.Bytes())
}



package contract

import (
	. "GoOnchain/common"
	pg "GoOnchain/core/contract/program"
	sig "GoOnchain/core/signature"
	"GoOnchain/crypto"
	. "GoOnchain/errors"
	"errors"
	_ "fmt"
	"math/big"
	"sort"
	"fmt"
)

type ContractContext struct {
	Data          sig.SignableData
	ProgramHashes []Uint160
	Codes         [][]byte
	Parameters    [][][]byte

	MultiPubkeyPara [][]PubkeyParameter
}

func NewContractContext(data sig.SignableData) *ContractContext {

	programHashes, _ := data.GetProgramHashes() //TODO: check error
	fmt.Println("programHashes=",programHashes)
	fmt.Println("hashLen := len(programHashes)",len(programHashes))
	hashLen := len(programHashes)
	return &ContractContext{
		Data:            data,
		ProgramHashes:   programHashes,
		Codes:           make([][]byte, hashLen),
		Parameters:      make([][][]byte, hashLen),
		MultiPubkeyPara: make([][]PubkeyParameter, hashLen),
	}
}

func (cxt *ContractContext) Add(contract *Contract, index int, parameter []byte) error {
	Trace()
	i := cxt.GetIndex(contract.ProgramHash)
	if i < 0 {
		return errors.New("Program Hash is not exist.")
	}
	if cxt.Codes[i] == nil {
		cxt.Codes[i] = contract.Code
	}
	if cxt.Parameters[i] == nil {
		cxt.Parameters[i] = make([][]byte, len(contract.Parameters))
	}
	cxt.Parameters[i][index] = parameter
	return nil
}

func (cxt *ContractContext) AddContract(contract *Contract, pubkey *crypto.PubKey, parameter []byte) error {
	Trace()
	if contract.GetType() == MultiSigContract {
		Trace()
		// add multi sig contract

		index := cxt.GetIndex(contract.ProgramHash)
		if index <= 0 {
			return errors.New("The program hash is not exist.")
		}
		if cxt.Codes[index] == nil {
			cxt.Codes[index] = contract.Code
		}
		if cxt.Parameters[index] == nil {
			cxt.Parameters[index] = make([][]byte, len(contract.Parameters))
		}

		pkParaArray := cxt.MultiPubkeyPara[index]
		temp, err := pubkey.EncodePoint(true)
		if err != nil {
			return NewDetailErr(err, ErrNoCode, "[Contract],AddContract failed.")
		}
		pubkeyPara := PubkeyParameter{
			PubKey:    ToHexString(temp),
			Parameter: ToHexString(parameter),
		}
		pkParaArray = append(pkParaArray, pubkeyPara)

		if len(pkParaArray) == len(contract.Parameters) {
			Trace()
			i := 0
			pubkeys := []*crypto.PubKey{}
			switch contract.Code[i] {
			case 1:
				i += 2
				break
			case 2:
				i += 3
				break
			}
			for contract.Code[i] == 33 {
				i++
				temp,err:=crypto.DecodePoint(contract.Code[i:33])
				if err!=nil{
					return NewDetailErr(err, ErrNoCode, "[Contract],AddContract DecodePoint failed.")
				}
				pubkeys = append(pubkeys,temp )
				i += 33
			}

			//generate Pubkey/Index map by pubkey array
			pkIndexMap := make(map[crypto.PubKey]int)
			for i, pk := range pubkeys {
				pkIndexMap[*pk] = i
			}

			//generate parameter/index map by pubkey parameter arrar
			paraIndexs := make([]ParameterIndex, len(pkParaArray))
			for _, pkPara := range pkParaArray {
				temp,err :=crypto.DecodePoint(HexToBytes(pkPara.PubKey))
				if err!=nil{
					return NewDetailErr(err, ErrNoCode, "[Contract],AddContract DecodePoint failed.")
				}
				paraIndex := ParameterIndex{
					Parameter: HexToBytes(pkPara.Parameter),
					Index:     pkIndexMap[*temp],
				}
				paraIndexs = append(paraIndexs, paraIndex)
			}

			//sort parameter by Index
			sort.Sort(sort.Reverse(ParameterIndexSlice(paraIndexs)))

			//generate sorted parameter list
			paras := make([][]byte, len(pkParaArray))
			for _, paIndex := range paraIndexs {
				paras = append(paras, paIndex.Parameter)
			}

			for i, para := range paras {
				if err := cxt.Add(contract, i, para); err != nil {
					return err
				}
			}

			cxt.MultiPubkeyPara[index] = nil

		} //pkParaArray
	} else {
		//add non multi sig contract
		Trace()
		index := -1
		for i := 0; i < len(contract.Parameters); i++ {
			if contract.Parameters[i] == Signature {
				if index >= 0 {
					return errors.New("Contract Parameters are not supported.")
				} else {
					index = i
				}
			}
		}
		return cxt.Add(contract, index, parameter)
	}
	return nil
}

func (cxt *ContractContext) GetIndex(programHash Uint160) int {
	for i := 0; i < len(cxt.ProgramHashes); i++ {
		if cxt.ProgramHashes[i] == programHash {
			return i
		}
	}
	return -1
}

func (cxt *ContractContext) GetPrograms() []*pg.Program {
	Trace()
	fmt.Println("!cxt.IsCompleted()=",!cxt.IsCompleted())
	fmt.Println(cxt.Codes)
	fmt.Println(cxt.Parameters)
	if !cxt.IsCompleted() {
		return nil
	}
	programs := make([]*pg.Program, len(cxt.Parameters))

	fmt.Println(" len(cxt.Codes)", len(cxt.Codes))

	for i := 0; i < len(cxt.Codes); i++ {
		sb := pg.NewProgramBuilder()

		for _, parameter := range cxt.Parameters[i] {
			if len(parameter) <= 2 {
				sb.PushNumber(new(big.Int).SetBytes(parameter))
			} else {
				sb.PushData(parameter)
			}
		}
		fmt.Println(" cxt.Codes[i])", cxt.Codes[i])
		fmt.Println(" sb.ToArray()", sb.ToArray())
		programs[i] = &pg.Program{
			Code:      cxt.Codes[i],
			Parameter: sb.ToArray(),
		}
	}
	return programs
}

func (cxt *ContractContext) IsCompleted() bool {
	for _, p := range cxt.Parameters {
		if p == nil {
			return false
		}

		for _, pp := range p {
			if pp == nil {
				return false
			}
		}
	}
	return true
}

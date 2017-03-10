package contract

import (
	. "GoOnchain/common"
	pg "GoOnchain/core/contract/program"
	sig "GoOnchain/core/signature"
	"GoOnchain/crypto"
	_ "GoOnchain/errors"
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

	//temp index for multi sig
	tempParaIndex int
}

func NewContractContext(data sig.SignableData) *ContractContext {
	Trace()
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
		tempParaIndex: 0,
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

		fmt.Println("Multi Sig: contract.ProgramHash:",contract.ProgramHash)
		fmt.Println("Multi Sig: cxt.ProgramHashes:",cxt.ProgramHashes)

		index := cxt.GetIndex(contract.ProgramHash)

		fmt.Println("Multi Sig: GetIndex:" ,index)

		if index < 0 {
			return errors.New("The program hash is not exist.")
		}

		fmt.Println("Multi Sig: contract.Code:" ,cxt.Codes[index])

		if cxt.Codes[index] == nil {
			cxt.Codes[index] = contract.Code
		}
		fmt.Println("Multi Sig: cxt.Codes[index]:" ,cxt.Codes[index])

		if cxt.Parameters[index] == nil {
			cxt.Parameters[index] = make([][]byte, len(contract.Parameters))
		}
		fmt.Println("Multi Sig: cxt.Parameters[index]:" ,cxt.Parameters[index])

		if err := cxt.Add(contract, cxt.tempParaIndex, parameter); err != nil {
			return err
		}

		cxt.tempParaIndex++

		//all paramenters added, sort the parameters
		if(cxt.tempParaIndex == len(contract.Parameters)){
			cxt.tempParaIndex = 0
		}

		//TODO: Sort the parameter according contract's PK list sequence
		//if err := cxt.AddSignatureToMultiList(index,contract,pubkey,parameter); err != nil {
		//	return err
		//}
		//
		//if(cxt.tempParaIndex == len(contract.Parameters)){
		//	//all multi sigs added, sort the sigs and add to context
		//	if err := cxt.AddMultiSignatures(index,contract,pubkey,parameter);err != nil {
		//		return err
		//	}
		//}

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

func (cxt *ContractContext) AddSignatureToMultiList(contractIndex int, contract *Contract, pubkey*crypto.PubKey, parameter []byte) error {
	if cxt.MultiPubkeyPara[contractIndex] == nil {
		cxt.MultiPubkeyPara[contractIndex] = make([]PubkeyParameter, len(contract.Parameters))
	}
	pk, err := pubkey.EncodePoint(true)
	if err != nil {
		return err
	}

	pubkeyPara := PubkeyParameter{
		PubKey:    ToHexString(pk),
		Parameter: ToHexString(parameter),
	}
	cxt.MultiPubkeyPara[contractIndex] = append(cxt.MultiPubkeyPara[contractIndex], pubkeyPara)

	return nil
}

func (cxt *ContractContext) AddMultiSignatures(index int,contract *Contract, pubkey *crypto.PubKey, parameter []byte) error{
	pkIndexs,err := cxt.ParseContractPubKeys(contract)
	if err != nil {
		return  errors.New("Contract Parameters are not supported.")
	}

	paraIndexs := []ParameterIndex{}
	for _, pubkeyPara := range cxt.MultiPubkeyPara[index] {
		paraIndex := ParameterIndex{
			Parameter: HexToBytes(pubkeyPara.Parameter),
			Index:    pkIndexs[pubkeyPara.PubKey],
		}
		paraIndexs = append(paraIndexs, paraIndex)
	}

	//sort parameter by Index
	sort.Sort(sort.Reverse(ParameterIndexSlice(paraIndexs)))

	//generate sorted parameter list
	for i, paraIndex := range paraIndexs {
		if err := cxt.Add(contract, i, paraIndex.Parameter); err != nil {
			return err
		}
	}

	cxt.MultiPubkeyPara[index] = nil

	return nil
}


func (cxt *ContractContext) ParseContractPubKeys(contract *Contract) (map[string]int,error) {

	pubkeyIndex := make(map[string]int)

	Index := 0
	//parse contract's pubkeys
	i := 0
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
		//pubkey, err := crypto.DecodePoint(contract.Code[i:33])
		//if err != nil {
		//	return nil, errors.New("[Contract],AddContract DecodePoint failed.")
		//}

		//add to parameter index
		pubkeyIndex[ToHexString(contract.Code[i:33])] = Index

		i += 33
		Index++
	}

	return pubkeyIndex,nil
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
	//fmt.Println("!cxt.IsCompleted()=",!cxt.IsCompleted())
	//fmt.Println(cxt.Codes)
	//fmt.Println(cxt.Parameters)
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
		//fmt.Println(" cxt.Codes[i])", cxt.Codes[i])
		//fmt.Println(" sb.ToArray()", sb.ToArray())
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

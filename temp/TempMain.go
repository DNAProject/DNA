package main

import (
	"GoOnchain/common"
	"GoOnchain/core/contract"
	tx "GoOnchain/core/transaction"
	_ "GoOnchain/core/ledger"
	_ "GoOnchain/core/transaction/payload"
	_ "GoOnchain/core/validation"
	"GoOnchain/core/asset"
	"GoOnchain/core/signature"
	"GoOnchain/crypto"
	"fmt"
)

func main() {
	fmt.Println("temp main:")


}

func RegisterAsset(assetName *string, assetAmount common.Fixed64, Issuer signature.Signer,Controller common.Uint160) error {

	newasset :=  &asset.Asset{
		Name: assetName,
		AssetType: asset.Token,
		RecordType: asset.UTXO,
	}

	TX,_ := tx.NewAssetRegistrationTransaction(*newasset,assetAmount,1,Issuer.PubKey(),Controller)
	cxt := contract.NewContractContext(TX)

	Sign(cxt)


	//context.Signable.Scripts = context.GetScripts();
	//Program.CurrentWallet.SaveTransaction(tx);
	//Program.LocalNode.Relay(tx);
	//InformationBox.Show(tx.Hash.ToString(), Strings.SendTxSucceedMessage, Strings.SendTxSucceedTitle);

	programs,_ := cxt.GetPrograms()
	cxt.Data.SetPrograms(programs)

	//Peer.Relay(TX)

	return nil
}

func Sign(cxt *contract.ContractContext)  {

	for _, programHash := range cxt.ProgramHashes {
		contract,_ := GetContract(programHash)
		account,_ := GetAccountByProgramHash(programHash)
		sig,_ := signature.Sign(cxt.Data,account)
		cxt.AddContract(contract,account.PubKey(),sig)
	}

}

func GetContract(programHash common.Uint160) (*contract.Contract,error) {
	//TODO: GetContract

	return nil,nil

}

func GetAccountByProgramHash(programHash common.Uint160) (*Account,error) {
	//TODO: GetAccountByProgramHash

	return &Account{},nil

}


type Account  struct {

}

func (a *Account) PrivKey() []byte {

	return nil
}

func (a *Account) PubKey() crypto.PubKey {

	return *&crypto.PubKey{}
}

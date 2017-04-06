package main

import (
	_ "DNA/consensus"
	_ "DNA/config"
	_ "DNA/errors"
	_ "DNA/events"
	_ "DNA/net"
	// "DNA/client"
	_ "DNA/crypto"
	_ "DNA/vm"
	_ "DNA/common/log"
	_ "DNA/common/serialization"
	_ "DNA/consensus/dbft"
	_ "DNA/core/contract"
	_ "DNA/net/httpjsonrpc"
	_ "DNA/core/validation"
	_ "DNA/net/message"
	_ "DNA/net/node"
	_ "DNA/net/protocol"
	_ "DNA/core/store/LevelDBStore"
	"fmt"
	"DNA/core/store"
	"crypto/sha256"
	"DNA/core/contract/program"
	"DNA/core/transaction/payload"
	"bytes"
	"io"
	. "DNA/common"
	. "DNA/core/asset"
	"DNA/core/ledger"
	"DNA/core/transaction"
	"DNA/crypto"
	"DNA/core/contract"
	"DNA/client"
	"DNA/core/signature"
	"DNA/core/validation"
	"time"
	"DNA/vm"
	"DNA/common/log"
)

func main() {
	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 0. Client Set                                                      ***")
	fmt.Println("//**************************************************************************")
	ledger.DefaultLedger = new(ledger.Ledger)
	ledger.DefaultLedger.Store = store.NewLedgerStore()
	ledger.DefaultLedger.Store.InitLedgerStore(ledger.DefaultLedger)
	transaction.TxStore = ledger.DefaultLedger.Store
	ledger.DefaultLedger.Blockchain = ledger.NewBlockchain()
	log.CreatePrintLog("./")
	crypto.SetAlg(crypto.P256R1)
	fmt.Println("  Client set completed. Test Start...")
	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 1. Generate [Account]                                              ***")
	fmt.Println("//**************************************************************************")
	//根据私钥创建新的账号
	user,_ := client.NewAccountWithPrivatekey([]byte{0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,0x09,0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01})
	admin,_:= client.NewAccountWithPrivatekey([]byte{0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,0x09,0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x01,0x03})
	userpubkey,err := user.PublicKey.EncodePoint(true)
	adminpubkey,err := admin.PublicKey.EncodePoint(true)
	fmt.Printf( "user.PrivateKey: %x user.PrivateKey Len: %d\n", user.PrivateKey,  len(user.PrivateKey) )
	fmt.Printf( "user.PublicKey: %x user.PublicKey Len: %d\n", userpubkey,  len(userpubkey) )
	fmt.Printf( "admin.PrivateKey: %x admin.PrivateKey Len: %d\n", admin.PrivateKey,  len(admin.PrivateKey) )
	fmt.Printf( "admin.PublicKey: %x admin.PublicKey Len: %d\n", adminpubkey,  len(adminpubkey) )

	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 2. Set Miner                                                     ***")
	fmt.Println("//**************************************************************************")
	miner := []*crypto.PubKey{}
	miner = append(miner, user.PublicKey)
	//miner = append(miner,miner2.PublicKey)
	ledger.StandbyMiners= miner
	fmt.Println("miner1.PublicKey",user.PublicKey)

	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 3. BlockChain init                                                 ***")
	fmt.Println("//**************************************************************************")
	sampleBlockchain := InitBlockChain()
	ledger.DefaultLedger.Blockchain = &sampleBlockchain



	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 3. Generate [Asset] Test                                           ***")
	fmt.Println("//**************************************************************************")
	a1 := SampleAsset()

	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 4. [controllerPGM] Generate Test                                   ***")
	fmt.Println("//**************************************************************************")
	controllerPGM,_ := contract.CreateSignatureContract(admin.PubKey())
	/*x
	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 5. Generate Sample [UTXO] Tx Input                                 ***")
	fmt.Println("//**************************************************************************")
	utxo := SampleUTXOTxInput()
	utxos := []*transaction.UTXOTxInput{}
	utxos = append(utxos, utxo)
	*/

	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 6. Generate [Transaction] Test                                     ***")
	fmt.Println("//**************************************************************************")
	ammount := Fixed64(10)
	tx,_ := transaction.NewAssetRegistrationTransaction(a1,&ammount,user.PubKey(),&controllerPGM.ProgramHash)
	//tx.UTXOInputs = utxos

	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 7. Generate [signature],[sign],set transaction [Program]           ***")
	fmt.Println("//**************************************************************************")

	//1.Transaction [Contract]
	transactionContract,_ := contract.CreateSignatureContract(user.PubKey())
	//2.Transaction Signdate
	signdate,err := signature.SignBySigner(tx, user)
	if err != nil{
		fmt.Println(err,"signdate SignBySigner failed")
	}

	pk := user.PubKey()
	publich,_ := pk.EncodePoint(true)
	message := signature.GetHashForSigning(tx)
	fmt.Printf("11111 signdate: %x\n",signdate)
	fmt.Printf("11111 public: %x\n",publich)
	fmt.Printf("11111 message: %x\n",message)

	//3.Transaction [contractContext]
	fmt.Printf("11111 transactionContract.Code: %x\n",transactionContract.Code)
	fmt.Printf("11111 transactionContract.Parameters: %x\n",transactionContract.Parameters)
	fmt.Printf("11111 transactionContract.ProgramHash: %x\n",transactionContract.ProgramHash)
	transactionContractContext := contract.NewContractContext(tx)
	//4.add  Contract , public key, signdate to ContractContext
	transactionContractContext.AddContract(transactionContract,user.PublicKey, signdate)
	fmt.Println("22222 transactionContract.Code=",transactionContractContext.Codes)
	fmt.Println("22222 ",transactionContractContext.GetPrograms()[0])

	//5.get ContractContext Programs & setinto transaction
	tx.SetPrograms(transactionContractContext.GetPrograms())

	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 8. Transaction [Validation]                                       ***")
	fmt.Println("//**************************************************************************")
	//1.validate transaction content
	err = validation.VerifyTransaction(tx, ledger.DefaultLedger, nil)
	if err !=nil{
		fmt.Println("Transaction Verify error.",err)
	}else{
		fmt.Println("Transaction Verify Normal Completed.")
	}
	//2.validate transaction signdate
	_,err = validation.VerifySignature(tx, user.PubKey(), signdate)
	if err !=nil{
		fmt.Println("Transaction Signature Verify error.",err)
	}else{
		fmt.Println("Transaction Signature Verify Normal Completed.")
	}

	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 9. Generate Sample [BlockDate] Test                                ***")
	fmt.Println("//**************************************************************************")
	//For test display
	genesisBlockHash, _ := ledger.DefaultLedger.Store.GetBlockHash(0)
	fmt.Println("gensisBlockGet =", genesisBlockHash)
	transactionHashes := []Uint256{}
	transactionHashes = append(transactionHashes,tx.Hash())

	genesisBlock ,_ :=ledger.DefaultLedger.GetBlockWithHeight(0)

	fmt.Println("genesisBlock.Hash()=",genesisBlock.Hash())
	fmt.Println("genesisBlockHash=",genesisBlockHash)

	sampleBlockdata := SampleBlockdateWithPreBlockHash(genesisBlockHash,transactionHashes)

	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 10. SampleBlockdate [Sign] Test                                    ***")
	fmt.Println("//**************************************************************************")
	/*
	//fmt.Println("11111",blockContractContext.GetPrograms()[0])
	//1.产生block Contract
	blockContract,_ := contract.CreateSignatureContract(user.PubKey())
	//2.产生签名
	blockSign ,err := signature.SignBySigner(sampleBlockdata,user)
	if err != nil{
		fmt.Println(err,"BlockSign signBySigner failed.")
	}

	fmt.Println("11111 transactionContract.Code",blockContract.Code)
	fmt.Println("11111 transactionContract.Parameters",blockContract.Parameters)
	fmt.Println("11111 transactionContract.ProgramHash",blockContract.ProgramHash)
	//3.产生block blockContractContext
	blockContractContext := contract.NewContractContext(sampleBlockdata)

	//4.将以上2个内容打包加到blockContractContext ***因为本例为单transaction 单签名，所以只加了一个param**
	sampleBlockdata.Hash()
	blockContractContext.AddContract(blockContract,user.PublicKey,blockSign)
	fmt.Println("22222 transactionContract.Code=",blockContractContext.Codes)
	//fmt.Println("22222 ",blockContractContext.GetPrograms()[0])

	//5.设定到blockdate中

	//fmt.Println("blockContractContext.GetPrograms()",blockContractContext.GetPrograms())
	sampleBlockdata.SetPrograms(blockContractContext.GetPrograms())
	//fmt.Println("len GetPrograms()",len(tx.GetPrograms()),"GetPrograms()=",tx.GetPrograms())
	*/
	/*
	//fmt.Println("11111",blockContractContext.GetPrograms()[0])
	//1.产生block Contract
	blockContract,_ := contract.CreateSignatureContract(miner1.PubKey())
	//2.产生签名
	blockSign ,err := signature.SignBySigner(sampleBlockdata,miner1)
	if err != nil{
		fmt.Println(err,"BlockSign signBySigner failed.")
	}

	fmt.Println("11111 transactionContract.Code",blockContract.Code)
	fmt.Println("11111 transactionContract.Parameters",blockContract.Parameters)
	fmt.Println("11111 transactionContract.ProgramHash",blockContract.ProgramHash)
	//3.产生block blockContractContext
	blockContractContext := contract.NewContractContext(sampleBlockdata)

	//4.将以上2个内容打包加到blockContractContext ***因为本例为单transaction 单签名，所以只加了一个param**
	sampleBlockdata.Hash()
	blockContractContext.AddContract(blockContract,user.PublicKey,blockSign)
	fmt.Println("22222 transactionContract.Code=",blockContractContext.Codes)
	//fmt.Println("22222 ",blockContractContext.GetPrograms()[0])

	//5.设定到blockdate中

	//fmt.Println("blockContractContext.GetPrograms()",blockContractContext.GetPrograms())
	sampleBlockdata.SetPrograms(blockContractContext.GetPrograms())
	//fmt.Println("len GetPrograms()",len(tx.GetPrograms()),"GetPrograms()=",tx.GetPrograms())
	//SignByMiner(sampleBlockdata,miner2)
	*/
	//1.产生block Contract
	blockContract,_ := contract.CreateSignatureContract(user.PubKey())
	//2.产生签名
	blockSign ,err := signature.SignBySigner(sampleBlockdata,user)
	if err != nil{
		fmt.Println(err,"BlockSign signBySigner failed.")
	}

	fmt.Println("11111 transactionContract.Code",blockContract.Code)
	fmt.Println("11111 transactionContract.Parameters",blockContract.Parameters)
	fmt.Println("11111 transactionContract.ProgramHash",blockContract.ProgramHash)
	//3.产生block blockContractContext
	blockContractContext := contract.NewContractContext(sampleBlockdata)

	//4.将以上2个内容打包加到blockContractContext ***因为本例为单transaction 单签名，所以只加了一个param**
	sampleBlockdata.Hash()
	blockContractContext.AddContract(blockContract,user.PublicKey,blockSign)
	fmt.Println("22222 transactionContract.Code=",blockContractContext.Codes)
	//fmt.Println("22222 ",blockContractContext.GetPrograms()[0])

	//5.设定到blockdate中

	//fmt.Println("blockContractContext.GetPrograms()",blockContractContext.GetPrograms())
	sampleBlockdata.SetPrograms(blockContractContext.GetPrograms())
	//fmt.Println("len GetPrograms()",len(tx.GetPrograms()),"GetPrograms()=",tx.GetPrograms())

	//SignByMiner(sampleBlockdata,user)

	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 11. Generate Sample [Block] Test                                    ***")
	fmt.Println("//**************************************************************************")
	Transcations := []*transaction.Transaction{}
	Transcations = append(Transcations, tx)
	sampleBlock := SampleBlock(sampleBlockdata, Transcations)
	//sampleBlockHast := sampleBlock.Hash()

	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 12. Block [Verify]                                                 ***")
	fmt.Println("//**************************************************************************")
	err = validation.VerifyBlock(sampleBlock,ledger.DefaultLedger,true)
	if err !=nil{
		fmt.Println("Block Verify error",err)
	}else{
		fmt.Println("Block Verify Normal Completed",err)
	}
	_, err = validation.VerifySignature(sampleBlockdata,user.PubKey(),blockSign)
	if err !=nil{
		fmt.Println("Block Signature Verify error.",err)
	}else{
		fmt.Println("Block Signature Verify Normal Completed.")
	}
	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 13. BlockChain Add Test                                            ***")
	fmt.Println("//**************************************************************************")
	fmt.Printf("Sample blockchain's height is %d befor add\n",sampleBlockchain.BlockHeight)
	sampleBlockchain.AddBlock(sampleBlock)
	//因为没有经过共识做区块+1操作，所以此处 手动+1
	sampleBlockchain.BlockHeight= sampleBlockchain.BlockHeight +1
	fmt.Printf("Sample blockchain's height is %d after add\n",sampleBlockchain.BlockHeight)

	fmt.Println("//**************************************************************************")
	fmt.Println("//*** 14. BlockChain Get Test                                            ***")
	fmt.Println("//**************************************************************************")
	//blockGet,_ := sampleBlockchain.GetBlockWithHash(sampleBlockHast)
	//SampleFuncTestBlock(blockGet)
}

//*********************************************************************************************
//** Init Func
//*********************************************************************************************
func InitBlockChain() ledger.Blockchain {
	blockchain := ledger.NewBlockchain()
	//if err!=nil{
	//	fmt.Println(err,"  BlockChain generate failed")
	//}
	fmt.Println("  BlockChain generate completed. Func test Start...")
	return *blockchain
}

//*********************************************************************************************
//** Sample Struct Test
//*********************************************************************************************
func SampleAsset() *Asset {
	var x string = "Onchain"
	a1 := Asset{Uint256(sha256.Sum256([]byte("a"))), x, byte(0x00), AssetType(Share), UTXO}
	fmt.Println("  Asset generate complete. Func test Start...")
	return &a1
}

func SampleBlockdateWithPreBlockHash(hash Uint256,transactionhash []Uint256) *ledger.Blockdata {
	fmt.Println("ledger.DefaultLedger.Blockchain.BlockHeight",ledger.DefaultLedger.Blockchain.BlockHeight)
	hashRoot,_:=crypto.ComputeRoot(transactionhash)
	var data = new(ledger.Blockdata)
	{
		data.Version = uint32(0)
		data.PrevBlockHash = hash
		data.TransactionsRoot = hashRoot
		tm := time.Now()
		data.Timestamp = uint32(tm.Unix())
		data.Height = ledger.DefaultLedger.Blockchain.BlockHeight + 1
		data.ConsensusData = uint64(0x11)
		pg := new(program.Program)
		pg.Code = []byte{'0'}
		pg.Parameter = []byte{byte(vm.PUSHT)}
		data.Program = pg
	}
	fmt.Println("  Blockdata generate completed. Func test Start...")
	return data
}

func SampletAssetRegistration(asset1 *Asset) payload.AssetRegistration {
	//使用input的asset来生成一个AssetRegistration对象
	ply := new(payload.AssetRegistration)
	ply.Asset = asset1
	fix := Fixed64(int64(0x11))
	ply.Amount = &fix
	pub := new(crypto.PubKey)
	ply.Issuer = pub
	var temp1 [20]uint8
	for i := 0; i < 20; i++ {
		temp1[i] = 0x11
	}
	temp2 := Uint160(temp1)
	ply.Controller = &temp2
	fmt.Println("  AssetRegistration generate completed. Func test Start...")
	return *ply
}

func SampleUTXOTxInput() *transaction.UTXOTxInput {
	utxo := new(transaction.UTXOTxInput)
	utxo.ReferTxID = Uint256(sha256.Sum256([]byte("a")))
	utxo.ReferTxOutputIndex = uint16(0x11)
	return utxo
}

func SampleTransaction(ply payload.AssetRegistration, utxos []*transaction.UTXOTxInput) transaction.Transaction {
	tx := new(transaction.Transaction)
	{
		//tx.TxType = transaction.AssetRegister
		tx.PayloadVersion = byte(0x11)
		tx.Payload = &ply
		tx.Nonce = uint64(0x11)
		tx.Attributes = nil
		tx.UTXOInputs = utxos
		tx.BalanceInputs = nil
		tx.Outputs = nil
		{
			programHashes := []*program.Program{}
			outputHashes := program.Program{nil, nil}
			programHashes = append(programHashes, &outputHashes)
			tx.Programs = programHashes
		}
	}
	fmt.Println("  transaction generate completed. Func test Start...")
	return *tx
}

func SampleBlock(data *ledger.Blockdata, transaction []*transaction.Transaction) *ledger.Block {
	var bk = new(ledger.Block)
	bk.Blockdata = data
	//bk.Transcations = transaction

	fmt.Println("  Block generate completed. Func test Start...")
	return bk
}

//*********************************************************************************************
//** Sample Func Test
//*********************************************************************************************
func SampleFuncTestAsset(a *Asset) {
	b := new(bytes.Buffer)
	SampleFuncTestAssetSerialize(b, a)
	SampleFuncTestAssetDeserialize(b)
}

func SampleFuncTestAssetSerialize(b io.Writer, a1 *Asset) {
	a1.Serialize(b)
	fmt.Println("  >>Serialize() :", a1.ID, a1.Name, a1.Precision, a1.AssetType, a1.RecordType)

}

func SampleFuncTestAssetDeserialize(b io.Reader) {
	a2 := new(Asset)
	a2.Deserialize(b)
	fmt.Println("  >>DeSerialize() :", a2.ID, a2.Name, a2.Precision, a2.AssetType, a2.RecordType)

}

func SampleFuncTestBlockdate(data ledger.Blockdata) {
	//Serialize Test
	b_buf := new(bytes.Buffer)
	data.Serialize(b_buf)
	fmt.Println("  >>Serialize() :", b_buf.Bytes())

	//DeSerialize Test
	var d2 = new(ledger.Blockdata)
	d2.Deserialize(b_buf)
	fmt.Println("  >>DeSerialize() :", d2)

	//GetProgramHashes() Test
	val, _ := data.GetProgramHashes()
	fmt.Println("  >>GetProgramHashes() :", val)

	//Hash() Test
	fmt.Println("  >>Hash() :", data.Hash())

}

func SampleFuncTestTransaction(tx transaction.Transaction) {
	b_buf := new(bytes.Buffer)
	tx.Serialize(b_buf)
	fmt.Println("  >>Serialize() :", b_buf.Bytes())
	var t2 = new(transaction.Transaction)
	t2.Deserialize(b_buf)
	fmt.Println("  >>DeSerialize() :", t2)
}

func SampleFuncTestUTXOTxInput(utxo *transaction.UTXOTxInput) {
	b_buf := new(bytes.Buffer)
	utxo.Serialize(b_buf)
	fmt.Println("  >>Serialize() :", b_buf.Bytes())
	var t2 = new(transaction.UTXOTxInput)
	t2.Deserialize(b_buf)
	fmt.Println("  >>DeSerialize() :", t2)
	fmt.Println("  >>Equals() :", utxo.Equals(t2))

}

func SampleFuncTestBlock(block *ledger.Block) {
	b_buf := new(bytes.Buffer)
	block.Serialize(b_buf)
	fmt.Println("  >>Serialize() :", b_buf.Bytes())
	var t2 = new(ledger.Block)
	t2.Deserialize(b_buf)
	fmt.Println("  >>DeSerialize() :", t2)
}

func SampleFuncTestPayload(ply transaction.Payload) {
	b_buf := new(bytes.Buffer)
	ply.Serialize(b_buf)
	fmt.Println("  >>Serialize() :", b_buf.Bytes())

	t2 := new(payload.AssetRegistration)
	t2.Deserialize(b_buf)
	fmt.Println("  >>DeSerialize() :", t2)
}

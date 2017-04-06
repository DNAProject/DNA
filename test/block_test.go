package test

import (
	"testing"
	"fmt"
	"DNA/client"
	"DNA/crypto"
	"DNA/core/ledger"
)

func TestBlock(t *testing.T) {
	t.Log("测试block开始!")
	t.Log("生成块链开始！")
	t.Log("创建旷工开始!")
	crypto.SetAlg(crypto.P256R1)
	miner := []*crypto.PubKey{}
	fmt.Println("打开a客户端")
	cl := client.OpenClient("/Users/tanyuan/Documents/wallet1.json", []byte{0x13, 0x34, 0x56})
	fmt.Println("获取a账号")
	account, err := cl.GetDefaultAccount()
	fmt.Println("a账号为：", account, err)
	fmt.Println("将a账号加入旷工行列")
	fmt.Println("公钥为：", account.PublicKey)
	miner = append(miner, account.PublicKey)

	fmt.Println("打开b客户端")
	cl = client.OpenClient("/Users/tanyuan/Documents/wallet2.json", []byte{0x13, 0x34, 0x56})
	fmt.Println("获取b账号")
	account, err = cl.GetDefaultAccount()
	fmt.Println("b账号为：", account, err)
	fmt.Println("将b账号加入旷工行列")
	fmt.Println("公钥为：", account.PublicKey)
	miner = append(miner, account.PublicKey)

	fmt.Println("打开c客户端")
	cl = client.OpenClient("/Users/tanyuan/Documents/wallet2.json", []byte{0x13, 0x34, 0x56})
	fmt.Println("获取c账号")
	account, err = cl.GetDefaultAccount()
	fmt.Println("c账号为：", account, err)
	fmt.Println("将c账号加入旷工行列")
	fmt.Println("公钥为：", account.PublicKey)
	miner = append(miner, account.PublicKey)

	fmt.Println("打开d客户端")
	cl = client.OpenClient("/Users/tanyuan/Documents/wallet2.json", []byte{0x13, 0x34, 0x56})
	fmt.Println("获取d账号")
	account, err = cl.GetDefaultAccount()
	fmt.Println("d账号为：", account, err)
	fmt.Println("将d账号加入旷工行列")
	fmt.Println("公钥为：", account.PublicKey)
	miner = append(miner, account.PublicKey)

	fmt.Println("所有旷工集合为：", miner)

	ledger.StandbyMiners = miner

	ledger.NewBlockchainWithGenesisBlock()
}
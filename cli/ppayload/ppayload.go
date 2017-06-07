package ppayload

import (
	. "DNA/cli/common"
	"DNA/client"
	"DNA/core/transaction"
	"DNA/core/transaction/payload"
	"DNA/crypto"
	"DNA/net/httpjsonrpc"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	. "github.com/bitly/go-simplejson"
	"github.com/urfave/cli"
)

func makePrivacyTx(admin *client.Account, toPubkeyStr string, pload string) (string, error) {
	data, _ := hex.DecodeString(pload)
	toPk, _ := hex.DecodeString(toPubkeyStr)
	bytesBuffer := bytes.NewBuffer(toPk)
	toPubkey := new(crypto.PubKey)
	toPubkey.DeSerialize(bytesBuffer)

	tx, _ := transaction.NewPrivacyPayloadTransaction(admin.PrivateKey, admin.PublicKey, toPubkey, payload.RawPayload, data)

	var buffer bytes.Buffer
	if err := tx.Serialize(&buffer); err != nil {
		fmt.Println("serialize registtransaction failed")
		return "", err
	}
	return hex.EncodeToString(buffer.Bytes()), nil
}

func ppayloadAction(c *cli.Context) error {
	if c.NumFlags() == 0 {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	enc := c.Bool("enc")
	dec := c.Bool("dec")
	getpk := c.Bool("getpk")
	if !enc && !dec && !getpk {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	if getpk {
		wallet := client.OpenClient(c.String("wallet"), []byte(c.String("password")))
		admin, _ := wallet.GetDefaultAccount()
		bytesBuffer := bytes.NewBuffer([]byte{})
		admin.PublicKey.Serialize(bytesBuffer)

		encoding, _ := json.Marshal(map[string]string{"pubkey": hex.EncodeToString(bytesBuffer.Bytes())})
		FormatOutput(encoding)
	}

	if enc {
		wallet := client.OpenClient(c.String("wallet"), []byte(c.String("password")))
		admin, _ := wallet.GetDefaultAccount()
		pload := c.String("payload")
		to := c.String("to")

		txHex, err := makePrivacyTx(admin, to, pload)
		resp, err := httpjsonrpc.Call(Address(), "sendrawtransaction", 0, []interface{}{txHex})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		FormatOutput(resp)
	}

	if dec {
		wallet := client.OpenClient(c.String("wallet"), []byte(c.String("password")))
		admin, _ := wallet.GetDefaultAccount()

		txhash := c.String("txhash")
		resp, err := httpjsonrpc.Call(Address(), "getrawtransaction", 0, []interface{}{txhash})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		js, err := NewJson(resp)
		txType, _ := js.Get("result").Get("TxType").Int()
		if transaction.TransactionType(txType) != transaction.PrivacyPayload {
			return errors.New("txType error")
		}

		plDataStr, _ := js.Get("result").Get("Payload").Get("Payload").String()
		plData, _ := hex.DecodeString(plDataStr)

		enType, _ := js.Get("result").Get("Payload").Get("EncryptType").Int()
		switch payload.PayloadEncryptType(enType) {
		case payload.ECDH_AES256:
			enAttr, _ := js.Get("result").Get("Payload").Get("EncryptAttr").String()
			Attr, _ := hex.DecodeString(enAttr)
			bytesBuffer := bytes.NewBuffer(Attr)
			encryptAttr := new(payload.EcdhAes256)
			encryptAttr.Deserialize(bytesBuffer)

			privkey := admin.PrivateKey
			data, _ := encryptAttr.Decrypt(plData, privkey)

			encoding, _ := json.Marshal(map[string]string{"result": hex.EncodeToString(data)})
			FormatOutput(encoding)

		default:
			return errors.New("enType error")
		}
	}

	return nil
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:        "ppayload",
		Usage:       "support encryption for payloads",
		Description: "With nodectl ppayload, you could create privacy payload.",
		ArgsUsage:   "[args]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "enc",
				Usage: "create an privacy  payload",
			},
			cli.BoolFlag{
				Name:  "dec",
				Usage: "decrypt the privacy payload",
			},
			cli.BoolFlag{
				Name:  "getpk",
				Usage: "get a public key form a wallet",
			},
			cli.StringFlag{
				Name:  "to",
				Usage: "payload to whom",
			},
			cli.StringFlag{
				Name:  "payload",
				Usage: "payload to be encrypted",
			},
			cli.StringFlag{
				Name:  "wallet",
				Usage: "wallet name",
				Value: DefaultWalletName,
			},
			cli.StringFlag{
				Name:  "password",
				Usage: "wallet password",
				Value: DefaultWalletPasswd,
			},
			cli.StringFlag{
				Name:  "txhash",
				Usage: "hash of a transaction",
			},
		},
		Action: ppayloadAction,
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			PrintError(c, err, "privacyPayload")
			return cli.NewExitError("", 1)
		},
	}
}

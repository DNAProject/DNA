package record

import (
	. "DNA/cli/common"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"DNA/client"
	"DNA/core/transaction"
	"math/rand"

	"DNA/net/httpjsonrpc"
	"github.com/urfave/cli"
)

func openWallet(name string, passwd []byte) client.Client {
	if name == DefaultWalletName {
		fmt.Println("Using default wallet: ", DefaultWalletName)
	}
	wallet := client.OpenClient(name, passwd)
	if wallet == nil {
		fmt.Println("Failed to open wallet: ", name)
		os.Exit(1)
	}
	return wallet
}

func recordAction(c *cli.Context) error {
	if c.NumFlags() == 0 {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	add := c.Bool("add")
	get := c.Bool("get")
	cat := c.Bool("cat")
	if !add && !get && !cat {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	// wallet := openWallet(c.String("wallet"), []byte(c.String("password")))
	// admin, _ := wallet.GetDefaultAccount()

	var resp []byte
	//var txHex string
	var err error
	if add {
		filepath := c.String("file")

		if filepath == "" {
			cli.ShowSubcommandHelp(c)
			return nil
		}
		name := c.String("name")
		if name == "" {
			rbuf := make([]byte, 4)
			rand.Read(rbuf)
			name = "DNA-" + hex.EncodeToString(rbuf)
		}
		var tx *transaction.Transaction

		if _, err := os.Stat(filepath); err != nil {
			fmt.Printf("invalid file path:%s\n", err)
			return err
		}

		f, err := os.OpenFile(filepath, os.O_RDONLY, 0664)
		defer f.Close()
		if err != nil {
			fmt.Printf("open file error:%s\n", err)
			return err
		}
		//read file
		var payload []byte
		var eof = false
		for {
			if eof {
				break
			}
			buf := make([]byte, 1024)
			nr, err := f.Read(buf[:])

			switch true {
			case nr < 0:
				fmt.Fprintf(os.Stderr, "cat: error reading: %s\n", err.Error())
				os.Exit(1)
			case nr == 0: // EOF
				eof = true
			case nr > 0:
				payload = append(payload, buf...)

			}
		}

		tx, _ = transaction.NewRecordTransaction("file", name, payload)
		tx.Nonce = uint64(rand.Int63())
		var buffer bytes.Buffer
		if err := tx.Serialize(&buffer); err != nil {
			fmt.Println("serialize registtransaction failed")
			return err
		}

		txHex := hex.EncodeToString(buffer.Bytes())

		resp, err = httpjsonrpc.Call(Address(), "addFileRecord", 0, []interface{}{txHex})

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

	}
	if cat {
		txhash := c.String("txhash")
		if txhash == "" {
			cli.ShowSubcommandHelp(c)
			return nil
		}
		if txhash != "" {
			resp, err = httpjsonrpc.Call(Address(), "catRecord", 0, []interface{}{txhash})

		}

	}
	if get {
		txhash := c.String("txhash")
		if txhash == "" {
			cli.ShowSubcommandHelp(c)
			return nil
		}
		if txhash != "" {
			resp, err = httpjsonrpc.Call(Address(), "getRecord", 0, []interface{}{txhash})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}
		}

	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	FormatOutput(resp)
	return nil
}

//NewCommand commands of ipfs and ipfs cluster
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "record",
		Usage: "record  registration and retrieve ",
		UsageText: `
This command can be used to manage record on chain.
`,
		ArgsUsage: "[args]",
		Flags: []cli.Flag{

			cli.BoolFlag{
				Name:  "add,a",
				Usage: "register record",
			},
			cli.BoolFlag{
				Name:  "cat,c",
				Usage: "read record that has been registered",
			},
			cli.BoolFlag{
				Name:  "get,g",
				Usage: "get record that has been registered",
			},
			cli.StringFlag{
				Name:  "name,n",
				Usage: "record name",
			},
			cli.StringFlag{
				Name:  "txhash,t",
				Usage: "transaction hash",
			},
			cli.StringFlag{
				Name:  "file,f",
				Usage: "record file path",
			},
			// cli.StringFlag{
			// 	Name:  "wallet, w",
			// 	Usage: "wallet name",
			// 	Value: DefaultWalletName,
			// },
			// cli.StringFlag{
			// 	Name:  "password, p",
			// 	Usage: "wallet password",
			// 	Value: DefaultWalletPasswd,
			// },
		},
		Action: recordAction,
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			PrintError(c, err, "record")
			return cli.NewExitError("", 1)
		},
	}
}

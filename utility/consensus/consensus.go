package consensus

import (
	"GoOnchain/net/httpjsonrpc"
	"GoOnchain/utility"
	"errors"
	"fmt"
	"os"
)

var usage = `switch of consensus function`

var flags = []string{"ip", "port", "rpcid", "start", "stop"}

func main(args []string, p utility.Param) (err error) {
	var resp []byte
	addr, err := utility.Address(p.Ip, p.Port)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	id := p.RPCID
	if p.Start && p.Stop {
		fmt.Fprintln(os.Stdout, "Which option do you want? (start or stop)")
		return errors.New("ambiguous flag")

	} else if p.Start {
		resp, err = httpjsonrpc.Call(addr, "startconsensus", id, []interface{}{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		fmt.Fprintln(os.Stdout, resp)

	} else if p.Stop {
		resp, err = httpjsonrpc.Call(addr, "stopconsensus", id, []interface{}{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		fmt.Fprintln(os.Stdout, resp)
	} else {
		fmt.Fprintln(os.Stdout, "Do you miss option? (start or stop)")
		return errors.New("missing flag")
	}

	utility.FormatOutput(resp)

	return nil
}

var Command = &utility.Command{UsageText: usage, Flags: flags, Main: main}

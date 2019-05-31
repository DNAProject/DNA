// Copyright 2016 DNA Dev team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"fmt"
	"os"

	. "DNA/cli/common"
	"DNA/net/httpjsonrpc"

	"github.com/urfave/cli"
)

func testAction(c *cli.Context) (err error) {
	if c.NumFlags() == 0 {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	txnType := c.String("tx")
	txnNum := c.Int("num")
	if txnType != "" {
		resp, err := httpjsonrpc.Call(Address(), "sendsampletransaction", 0, []interface{}{txnType, txnNum})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		FormatOutput(resp)
	}
	return nil
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:        "test",
		Usage:       "run test routine",
		Description: "With nodectl test, you could run simple tests.",
		ArgsUsage:   "[args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "tx, t",
				Usage: "send sample transaction",
				Value: "perf",
			},
			cli.IntFlag{
				Name:  "num, n",
				Usage: "sample transaction numbers",
				Value: 1,
			},
		},
		Action: testAction,
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			PrintError(c, err, "test")
			return cli.NewExitError("", 1)
		},
	}
}

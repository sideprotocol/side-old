package main

import (
	"fmt"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/sideprotocol/cmd/sided/cmd"
	"github.com/sideprotocol/side/app"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}

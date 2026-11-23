package cmd

import (
	"fmt"
	"os"

	gosdkcmd "github.com/bizshuk/gosdk/cmd"
	"github.com/bizshuk/gosdk/config"
	_ "github.com/bizshuk/gosdk/log"
	"github.com/bizshuk/upasspdf/config/example"
	"github.com/spf13/cobra"
)

// APP_NAME 決定設定目錄 ~/.config/upasspdf.
const APP_NAME = "upasspdf"

// RootCmd 是 upasspdf 的根命令.
var RootCmd = &cobra.Command{
	Use:           APP_NAME,
	Short:         "Remove the open password from PDF files in place",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute 執行根命令.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		config.Default(config.WithAppName(APP_NAME), config.WithDefaultValue(string(example.Settings)))
	})
	RootCmd.AddCommand(gosdkcmd.ConfigCmd, RunCmd, DecryptCmd)
}

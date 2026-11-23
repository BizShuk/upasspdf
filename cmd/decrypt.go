package cmd

import (
	"fmt"

	"github.com/bizshuk/upasspdf/svc/batch"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	decryptDir      string
	decryptPassword string
	decryptFilename string
	decryptDryRun   bool
)

// DecryptCmd 以手動指定的資料夾與密碼解除 PDF 密碼, 不讀 entities 設定.
var DecryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt PDFs in a given folder with a given password",
	Example: `  upasspdf decrypt --dir ~/Downloads/bank-estatements/中國信託_CTBC --password A123456789
  upasspdf decrypt -d ./statements -p secret -f "*_Estatement_*.pdf"`,
	RunE: func(c *cobra.Command, args []string) error {
		dir, err := homedir.Expand(decryptDir)
		if err != nil {
			return err
		}
		summary := batch.Run(c.Context(), batch.Job{
			Dir:      dir,
			Filename: decryptFilename,
			Password: decryptPassword,
			DryRun:   decryptDryRun,
		})
		fmt.Println(summary)
		if summary.Failed > 0 {
			return fmt.Errorf("%d file(s) failed to decrypt", summary.Failed)
		}
		return nil
	},
}

func init() {
	DecryptCmd.Flags().StringVarP(&decryptDir, "dir", "d", "", "source folder holding the PDFs (required)")
	DecryptCmd.Flags().StringVarP(&decryptPassword, "password", "p", "", "PDF open password")
	DecryptCmd.Flags().StringVarP(&decryptFilename, "filename", "f", "*.pdf", "glob matched against file names in the folder")
	DecryptCmd.Flags().BoolVar(&decryptDryRun, "dry-run", false, "list encrypted files without rewriting them")
	_ = DecryptCmd.MarkFlagRequired("dir")
}

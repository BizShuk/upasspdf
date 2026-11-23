package cmd

import (
	"fmt"

	"github.com/bizshuk/upasspdf/model"
	"github.com/bizshuk/upasspdf/svc/batch"
	"github.com/spf13/cobra"
)

var (
	runEntity string
	runDryRun bool
)

// RunCmd 依 settings.json 的 entities 逐一解除各機構資料夾內 PDF 的密碼.
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Decrypt PDFs for every entity configured in settings.json",
	RunE: func(c *cobra.Command, args []string) error {
		settings, err := model.LoadSettings()
		if err != nil {
			return err
		}
		entities, err := settings.Select(runEntity)
		if err != nil {
			return err
		}
		failed := 0
		for _, e := range entities {
			password, ok := e.Password()
			if !ok {
				// 未設定密碼的 entity 只列出待處理檔, 不算失敗, 排程才不會因尚未填密碼而反覆報錯.
				fmt.Printf("[%s] password env %s is not set, listing only\n", e.Name, e.PasswordEnv)
			}
			summary := batch.Run(c.Context(), batch.Job{
				Dir:      settings.EntityDir(e),
				Filename: e.Filename,
				Password: password,
				DryRun:   runDryRun || !ok,
			})
			fmt.Printf("[%s] %s\n", e.Name, summary)
			failed += summary.Failed
		}
		if failed > 0 {
			return fmt.Errorf("%d file(s) failed to decrypt", failed)
		}
		return nil
	},
}

func init() {
	RunCmd.Flags().StringVarP(&runEntity, "entity", "e", "", "only process the entity with this name")
	RunCmd.Flags().BoolVar(&runDryRun, "dry-run", false, "list encrypted files without rewriting them")
}

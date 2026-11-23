// Package pdf 處理單一 PDF 檔的密碼解除.
package pdf

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// ErrNotEncrypted 表示檔案本來就沒有加密.
var ErrNotEncrypted = pdfcpu.ErrNotEncrypted

// ErrWrongPassword 表示提供的密碼無法開啟檔案.
var ErrWrongPassword = pdfcpu.ErrWrongPassword

func init() {
	// pdfcpu 預設會在 os.UserConfigDir 建立自己的設定目錄; 本工具不需要.
	api.DisableConfigDir()
}

// IsEncrypted 回報檔案是否帶有加密字典.
func IsEncrypted(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	ctx, err := pdfcpu.Read(context.Background(), f, newConfig(""))
	if err != nil {
		// 讀取時就需要密碼, 代表有加密.
		if errors.Is(err, pdfcpu.ErrWrongPassword) {
			return true, nil
		}
		return false, fmt.Errorf("read %s: %w", path, err)
	}
	return ctx.Encrypt != nil, nil
}

// DecryptInPlace 以 password 解除 path 的加密並覆寫同名檔案, 保留原本的修改時間.
// 檔案未加密時回傳 ErrNotEncrypted, 密碼錯誤時回傳 ErrWrongPassword.
func DecryptInPlace(ctx context.Context, path, password string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	err = api.DecryptFile(ctx, path, "", newConfig(password))
	if errors.Is(err, ErrWrongPassword) && password != "" {
		// 只設 owner 權限限制的檔案 user password 為空, 給了密碼反而驗不過.
		err = api.DecryptFile(ctx, path, "", newConfig(""))
	}
	if err != nil {
		return err
	}
	return os.Chtimes(path, info.ModTime(), info.ModTime())
}

func newConfig(password string) *model.Configuration {
	conf := model.NewDefaultConfiguration()
	conf.UserPW = password
	conf.OwnerPW = password
	return conf
}

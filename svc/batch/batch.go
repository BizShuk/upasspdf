// Package batch 對一個資料夾內符合檔名 glob 的 PDF 逐一解除密碼.
package batch

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bizshuk/upasspdf/svc/pdf"
)

// DEFAULT_FILENAME 是未指定檔名 glob 時使用的樣式.
const DEFAULT_FILENAME = "*.pdf"

// Job 描述一次資料夾批次處理.
type Job struct {
	Dir      string
	Filename string
	Password string
	DryRun   bool
}

// Summary 統計一次批次處理的結果.
type Summary struct {
	Decrypted int
	Pending   int // dry-run 時待解密的檔案數
	Skipped   int // 本來就沒有加密
	Failed    int
}

func (s Summary) String() string {
	return fmt.Sprintf("decrypted=%d pending=%d skipped=%d failed=%d", s.Decrypted, s.Pending, s.Skipped, s.Failed)
}

// Run 處理 job.Dir 內 (不遞迴) 檔名符合 job.Filename 的 PDF; 單檔失敗會記錄後繼續.
func Run(ctx context.Context, job Job) Summary {
	var sum Summary
	pattern := job.Filename
	if pattern == "" {
		pattern = DEFAULT_FILENAME
	}
	entries, err := os.ReadDir(job.Dir)
	if err != nil {
		slog.Error("read dir failed", "dir", job.Dir, "filename", pattern, "err", err)
		sum.Failed++
		return sum
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(name), ".pdf") {
			continue
		}
		if ok, err := filepath.Match(pattern, name); err != nil || !ok {
			if err != nil {
				slog.Error("invalid filename glob", "dir", job.Dir, "filename", pattern, "err", err)
				sum.Failed++
				return sum
			}
			continue
		}
		path := filepath.Join(job.Dir, name)
		encrypted, err := pdf.IsEncrypted(path)
		if err != nil {
			slog.Error("inspect pdf failed", "path", path, "dir", job.Dir, "err", err)
			sum.Failed++
			continue
		}
		if !encrypted {
			sum.Skipped++
			continue
		}
		if job.DryRun {
			fmt.Println("encrypted:", path)
			sum.Pending++
			continue
		}
		if err := pdf.DecryptInPlace(ctx, path, job.Password); err != nil {
			if errors.Is(err, pdf.ErrNotEncrypted) {
				sum.Skipped++
				continue
			}
			slog.Error("decrypt pdf failed", "path", path, "dir", job.Dir, "password_set", job.Password != "", "err", err)
			sum.Failed++
			continue
		}
		fmt.Println("decrypted:", path)
		sum.Decrypted++
	}
	return sum
}

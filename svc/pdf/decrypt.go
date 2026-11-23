// Package pdf 處理單一 PDF 檔的密碼解除, 底層呼叫 qpdf.
package pdf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// QPDF_BIN 是解密使用的外部程式, 由 $PATH 解析 (brew install qpdf).
const QPDF_BIN = "qpdf"

// qpdf 的 exit code.
const (
	EXIT_OK       = 0
	EXIT_ERROR    = 2
	EXIT_WARNINGS = 3
)

// ErrNotEncrypted 表示檔案本來就沒有加密.
var ErrNotEncrypted = errors.New("pdf is not encrypted")

// ErrWrongPassword 表示提供的密碼無法開啟檔案.
var ErrWrongPassword = errors.New("invalid pdf password")

// IsEncrypted 回報檔案是否加密 (不需要密碼即可判斷).
func IsEncrypted(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		return false, err
	}
	code, stderr, err := runQPDF(context.Background(), "", "--is-encrypted", "--", path)
	if err != nil {
		return false, err
	}
	switch code {
	case EXIT_OK:
		return true, nil
	case EXIT_ERROR:
		if stderr == "" {
			return false, nil
		}
	}
	return false, fmt.Errorf("qpdf --is-encrypted %s: exit %d: %s", path, code, stderr)
}

// DecryptInPlace 以 password 解除 path 的加密並覆寫同名檔案, 保留原本的修改時間.
// 檔案未加密時回傳 ErrNotEncrypted, 密碼錯誤時回傳 ErrWrongPassword; 失敗時原檔不變.
func DecryptInPlace(ctx context.Context, path, password string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	encrypted, err := IsEncrypted(path)
	if err != nil {
		return err
	}
	if !encrypted {
		return ErrNotEncrypted
	}
	err = decrypt(ctx, path, password)
	if errors.Is(err, ErrWrongPassword) && password != "" {
		// 只設 owner 權限限制的檔案 user password 為空, 給了密碼反而驗不過.
		err = decrypt(ctx, path, "")
	}
	if err != nil {
		return err
	}
	return os.Chtimes(path, info.ModTime(), info.ModTime())
}

func decrypt(ctx context.Context, path, password string) error {
	// 寫到同資料夾的暫存檔再 rename: --replace-input 遇到警告會留下 .~qpdf-orig.
	tmp, err := os.CreateTemp(filepath.Dir(path), ".upasspdf-*.pdf")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpPath)

	// --password-file=- 從 stdin 讀密碼, 避免密碼出現在 process 參數.
	code, stderr, err := runQPDF(ctx, password+"\n", "--password-file=-", "--decrypt", "--", path, tmpPath)
	if err != nil {
		return err
	}
	switch {
	case code == EXIT_OK:
	case code == EXIT_WARNINGS:
		// 例: 元大的 /Perms 與 /EncryptMetadata 不一致; 內容仍可正確解密.
		slog.Warn("qpdf decrypted with warnings", "path", path, "exit_code", code, "stderr", stderr)
	case strings.Contains(stderr, "invalid password"):
		return ErrWrongPassword
	default:
		return fmt.Errorf("qpdf --decrypt %s: exit %d: %s", path, code, stderr)
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	// CreateTemp 建立的檔案是 0600, 改回原檔權限再取代.
	if err := os.Chmod(tmpPath, info.Mode().Perm()); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

// runQPDF 執行 qpdf 並回傳 exit code 與 stderr; 只有無法啟動 qpdf 時才回傳 error.
func runQPDF(ctx context.Context, stdin string, args ...string) (int, string, error) {
	cmd := exec.CommandContext(ctx, QPDF_BIN, args...)
	cmd.Stdin = strings.NewReader(stdin)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	var exitErr *exec.ExitError
	if err != nil && !errors.As(err, &exitErr) {
		return 0, "", fmt.Errorf("run %s (brew install qpdf): %w", QPDF_BIN, err)
	}
	return cmd.ProcessState.ExitCode(), strings.TrimSpace(stderr.String()), nil
}

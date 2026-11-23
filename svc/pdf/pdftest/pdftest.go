// Package pdftest 產生測試用的明文與加密 PDF.
package pdftest

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/stretchr/testify/require"
)

// minimalPDF 產生一頁空白 PDF, xref offset 依實際位元組計算.
func minimalPDF() []byte {
	objs := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 200 200] >>",
	}
	var b bytes.Buffer
	b.WriteString("%PDF-1.7\n")
	offsets := make([]int, len(objs))
	for i, o := range objs {
		offsets[i] = b.Len()
		fmt.Fprintf(&b, "%d 0 obj\n%s\nendobj\n", i+1, o)
	}
	xref := b.Len()
	fmt.Fprintf(&b, "xref\n0 %d\n0000000000 65535 f \n", len(objs)+1)
	for _, off := range offsets {
		fmt.Fprintf(&b, "%010d 00000 n \n", off)
	}
	fmt.Fprintf(&b, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objs)+1, xref)
	return b.Bytes()
}

// WritePlainPDF 寫出未加密 PDF.
func WritePlainPDF(t testing.TB, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, minimalPDF(), 0o644))
	return path
}

// WriteEncryptedPDF 寫出以 userPW 加密的 PDF (owner password 與 user password 不同).
func WriteEncryptedPDF(t testing.TB, dir, name, userPW string) string {
	t.Helper()
	path := WritePlainPDF(t, dir, name)
	conf := model.NewAESConfiguration(userPW, "owner-"+userPW, 256)
	require.NoError(t, api.EncryptFile(context.Background(), path, "", conf))
	return path
}

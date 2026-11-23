package batch

import (
	"context"
	"testing"

	"github.com/bizshuk/upasspdf/svc/pdf"
	"github.com/bizshuk/upasspdf/svc/pdf/pdftest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()
	stmt := pdftest.WriteEncryptedPDF(t, dir, "2026-01_Estatement.pdf", "secret")
	other := pdftest.WriteEncryptedPDF(t, dir, "2026-01_PaymentSlip.pdf", "other")
	pdftest.WritePlainPDF(t, dir, "2026-02_Estatement.pdf")

	sum := Run(context.Background(), Job{Dir: dir, Filename: "*_Estatement.pdf", Password: "secret"})
	assert.Equal(t, Summary{Decrypted: 1, Skipped: 1}, sum)

	enc, err := pdf.IsEncrypted(stmt)
	require.NoError(t, err)
	assert.False(t, enc)
	enc, err = pdf.IsEncrypted(other)
	require.NoError(t, err)
	assert.True(t, enc, "file outside the glob stays encrypted")

	// 再跑一次: 已解密檔視為 skipped, 達成冪等.
	assert.Equal(t, Summary{Skipped: 2}, Run(context.Background(), Job{Dir: dir, Filename: "*_Estatement.pdf", Password: "secret"}))
}

func TestRunDryRunAndWrongPassword(t *testing.T) {
	dir := t.TempDir()
	path := pdftest.WriteEncryptedPDF(t, dir, "a.pdf", "secret")

	assert.Equal(t, Summary{Pending: 1}, Run(context.Background(), Job{Dir: dir, DryRun: true}))
	assert.Equal(t, Summary{Failed: 1}, Run(context.Background(), Job{Dir: dir, Password: "wrong"}))

	enc, err := pdf.IsEncrypted(path)
	require.NoError(t, err)
	assert.True(t, enc)
}

func TestRunMissingDir(t *testing.T) {
	assert.Equal(t, Summary{Failed: 1}, Run(context.Background(), Job{Dir: "/nonexistent/upasspdf"}))
}

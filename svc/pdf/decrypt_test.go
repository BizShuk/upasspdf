package pdf_test

import (
	"context"
	"os"
	"testing"
	"time"

	. "github.com/bizshuk/upasspdf/svc/pdf"
	"github.com/bizshuk/upasspdf/svc/pdf/pdftest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecryptInPlace(t *testing.T) {
	dir := t.TempDir()
	path := pdftest.WriteEncryptedPDF(t, dir, "a.pdf", "secret")
	mtime := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	require.NoError(t, os.Chtimes(path, mtime, mtime))

	enc, err := IsEncrypted(path)
	require.NoError(t, err)
	assert.True(t, enc)

	require.NoError(t, DecryptInPlace(context.Background(), path, "secret"))

	enc, err = IsEncrypted(path)
	require.NoError(t, err)
	assert.False(t, enc)
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.ModTime().Equal(mtime), "mtime kept")
}

func TestDecryptInPlaceWrongPassword(t *testing.T) {
	dir := t.TempDir()
	path := pdftest.WriteEncryptedPDF(t, dir, "a.pdf", "secret")
	before, err := os.ReadFile(path)
	require.NoError(t, err)

	err = DecryptInPlace(context.Background(), path, "wrong")
	assert.ErrorIs(t, err, ErrWrongPassword)

	after, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, before, after, "original left untouched")
}

func TestDecryptInPlaceNotEncrypted(t *testing.T) {
	path := pdftest.WritePlainPDF(t, t.TempDir(), "a.pdf")
	assert.ErrorIs(t, DecryptInPlace(context.Background(), path, "secret"), ErrNotEncrypted)
}

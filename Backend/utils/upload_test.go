package utils

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func multipartUploadContext(t *testing.T, name string, content []byte) *gin.Context {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("lampiran", name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = req
	return ctx
}

func TestUploadFileAcceptsJPGAndJPEG(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("UPLOAD_DIR", dir)
	content := []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00}
	for _, name := range []string{"bukti.jpg", "bukti.jpeg"} {
		ctx := multipartUploadContext(t, name, content)
		result, err := UploadFileMetadata(ctx, "lampiran")
		if err != nil {
			t.Fatalf("%s rejected: %v", name, err)
		}
		if result.MIMEType != "image/jpeg" || result.Original != name || result.Size != int64(len(content)) {
			t.Fatalf("unexpected metadata: %#v", result)
		}
		if filepath.Base(result.Path) != result.Path {
			t.Fatalf("unsafe stored path: %q", result.Path)
		}
		if _, err := os.Stat(filepath.Join(dir, result.Path)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestUploadFileRejectsUnsupportedTypeAndOversize(t *testing.T) {
	t.Setenv("UPLOAD_DIR", t.TempDir())
	if _, err := UploadFileMetadata(multipartUploadContext(t, "bukti.txt", []byte("plain text")), "lampiran"); err == nil {
		t.Fatal("unsupported MIME should be rejected")
	}
	large := make([]byte, MaxUploadSize+1)
	copy(large, []byte("%PDF-1.7\n"))
	if _, err := UploadFileMetadata(multipartUploadContext(t, "bukti.pdf", large), "lampiran"); err == nil {
		t.Fatal("oversize file should be rejected")
	}
}

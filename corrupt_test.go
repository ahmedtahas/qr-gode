package qrgode

import (
	"os"
	"testing"
)

func TestCorruptedInputs(t *testing.T) {
	// Create temp corrupted files
	tmpDir := t.TempDir()

	// Fake PNG (text file with .png extension)
	fakePNG := tmpDir + "/fake.png"
	os.WriteFile(fakePNG, []byte("not a real png file"), 0644)

	// Empty file
	emptyFile := tmpDir + "/empty.png"
	os.WriteFile(emptyFile, []byte{}, 0644)

	// Random bytes
	randomFile := tmpDir + "/random.png"
	os.WriteFile(randomFile, []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0x89, 0x50, 0x4E}, 0644)

	// Truncated PNG header
	truncatedPNG := tmpDir + "/truncated.png"
	os.WriteFile(truncatedPNG, []byte{0x89, 0x50, 0x4E, 0x47}, 0644) // PNG magic but incomplete

	// Valid PNG header but corrupted body
	corruptPNG := tmpDir + "/corrupt.png"
	os.WriteFile(corruptPNG, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0xFF, 0xFF}, 0644)

	tests := []struct {
		name        string
		fn          func() ([]byte, error)
		shouldPanic bool
	}{
		// Corrupted logo files
		{
			name: "fake PNG as logo",
			fn: func() ([]byte, error) {
				return New("test").Logo(fakePNG).SVG()
			},
		},
		{
			name: "empty file as logo",
			fn: func() ([]byte, error) {
				return New("test").Logo(emptyFile).SVG()
			},
		},
		{
			name: "random bytes as logo",
			fn: func() ([]byte, error) {
				return New("test").Logo(randomFile).SVG()
			},
		},
		{
			name: "truncated PNG as logo",
			fn: func() ([]byte, error) {
				return New("test").Logo(truncatedPNG).SVG()
			},
		},
		{
			name: "corrupt PNG as logo",
			fn: func() ([]byte, error) {
				return New("test").Logo(corruptPNG).SVG()
			},
		},

		// Corrupted finder/module images
		{
			name: "fake PNG as finder",
			fn: func() ([]byte, error) {
				return New("test").FinderImage(fakePNG).SVG()
			},
		},
		{
			name: "fake PNG as module",
			fn: func() ([]byte, error) {
				return New("test").ModuleImage(fakePNG).SVG()
			},
		},

		// Weird data inputs
		{
			name: "null bytes in data",
			fn: func() ([]byte, error) {
				return New("test\x00\x00data").SVG()
			},
		},
		{
			name: "only null bytes",
			fn: func() ([]byte, error) {
				return New("\x00\x00\x00").SVG()
			},
		},
		{
			name: "binary data",
			fn: func() ([]byte, error) {
				return New(string([]byte{0x89, 0x50, 0x4E, 0x47, 0xFF, 0xFE})).SVG()
			},
		},
		{
			name: "very long unicode",
			fn: func() ([]byte, error) {
				data := ""
				for i := 0; i < 500; i++ {
					data += "日本語"
				}
				return New(data).SVG()
			},
		},
		{
			name: "mixed encodings",
			fn: func() ([]byte, error) {
				return New("Hello\x80\x81World日本語\xFF").SVG()
			},
		},
		{
			name: "newlines and tabs",
			fn: func() ([]byte, error) {
				return New("line1\nline2\tline3\r\nline4").SVG()
			},
		},
		{
			name: "SQL injection attempt",
			fn: func() ([]byte, error) {
				return New("'; DROP TABLE users; --").SVG()
			},
		},
		{
			name: "XSS attempt",
			fn: func() ([]byte, error) {
				return New("<script>alert('xss')</script>").SVG()
			},
		},
		{
			name: "path traversal in data",
			fn: func() ([]byte, error) {
				return New("../../../etc/passwd").SVG()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC: %v", r)
				}
			}()

			result, err := tt.fn()
			// We don't care if it errors, just that it doesn't panic
			t.Logf("result=%d bytes, err=%v", len(result), err)
		})
	}
}

func TestDirectoryAsFile(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PANIC: %v", r)
		}
	}()

	// Try to use a directory as a logo
	_, err := New("test").Logo("/tmp").SVG()
	t.Logf("directory as logo: %v", err)

	// Try to use /dev/null
	_, err = New("test").Logo("/dev/null").SVG()
	t.Logf("/dev/null as logo: %v", err)
}

func TestSymlinkLoop(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PANIC: %v", r)
		}
	}()

	tmpDir := t.TempDir()

	// Create symlink loop (if possible)
	link1 := tmpDir + "/link1.png"
	link2 := tmpDir + "/link2.png"
	os.Symlink(link2, link1)
	os.Symlink(link1, link2)

	_, err := New("test").Logo(link1).SVG()
	t.Logf("symlink loop: %v", err)
}

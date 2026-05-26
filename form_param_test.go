package ino

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMultipartFormParam(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		boundary  string
		maxMemory int64
		optional  bool
		hasError  bool
	}{
		{
			name:      "valid multipart form",
			body:      "--boundary\r\nContent-Disposition: form-data; name=\"field\"\r\n\r\nvalue\r\n--boundary--",
			boundary:  "boundary",
			maxMemory: 1024,
			optional:  false,
			hasError:  false,
		},
		{
			name:      "invalid multipart form required",
			body:      "invalid multipart data",
			boundary:  "",
			maxMemory: 1024,
			optional:  false,
			hasError:  true,
		},
		{
			name:      "invalid multipart form optional",
			body:      "invalid multipart data",
			boundary:  "",
			maxMemory: 1024,
			optional:  true,
			hasError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.body))
			if tt.boundary != "" {
				req.Header.Set("Content-Type", "multipart/form-data; boundary="+tt.boundary)
			}

			param := MultipartFormParam(tt.maxMemory)
			if tt.optional {
				param = param.Optional()
			}

			result, errs := param.GetParamValue(req)

			if tt.hasError && len(errs) == 0 {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}

			if !tt.hasError && tt.boundary != "" {
				if result == nil {
					t.Errorf("expected non-nil multipart form")
				}
			}
		})
	}
}

// TestFileParamValidFile tests FileParam with a valid file upload
func TestFileParamValidFile(t *testing.T) {
	req, err := createMultipartRequest("file", "testfile.txt", []byte("file content"))
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}

	param := FileParam("file")
	fileHeader, errs := param.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("expected no error but got: %v", errs)
	}
	if fileHeader == nil {
		t.Errorf("expected non-nil file header")
	}
	if fileHeader.Filename != "testfile.txt" {
		t.Errorf("expected filename 'testfile.txt', got %q", fileHeader.Filename)
	}
	if fileHeader.Size != 12 {
		t.Errorf("expected size 12, got %d", fileHeader.Size)
	}
}

// TestFileParamMissingRequiredFile tests FileParam when required file is missing
func TestFileParamMissingRequiredFile(t *testing.T) {
	req, err := createMultipartRequest("upload", "file.txt", []byte("content"))
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}

	param := FileParam("nonexistent")
	fileHeader, errs := param.GetParamValue(req)

	if len(errs) == 0 {
		t.Errorf("expected error for missing required file")
	}
	if fileHeader != nil {
		t.Errorf("expected nil file header for missing file")
	}
	if len(errs) > 0 && !strings.Contains(errs[0], "nonexistent") {
		t.Errorf("error message should mention field name 'nonexistent', got: %v", errs[0])
	}
}

// TestFileParamOptionalMissing tests FileParam.Optional() when file is not provided
func TestFileParamOptionalMissing(t *testing.T) {
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")

	param := FileParam("file").Optional()
	fileHeader, errs := param.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("expected no error for optional missing file but got: %v", errs)
	}
	if fileHeader != nil {
		t.Errorf("expected nil file header for optional missing file")
	}
}

// TestFileParamOptionalProvided tests FileParam.Optional() when file is provided
func TestFileParamOptionalProvided(t *testing.T) {
	req, err := createMultipartRequest("avatar", "profile.jpg", []byte("jpeg data"))
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}

	param := FileParam("avatar").Optional()
	fileHeader, errs := param.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("expected no error but got: %v", errs)
	}
	if fileHeader == nil {
		t.Errorf("expected non-nil file header")
	}
	if fileHeader.Filename != "profile.jpg" {
		t.Errorf("expected filename 'profile.jpg', got %q", fileHeader.Filename)
	}
}

// TestFileParamMultipleFilesReturnsFirst tests that FileParam returns only the first file
// when multiple files are uploaded under the same field name
func TestFileParamMultipleFilesReturnsFirst(t *testing.T) {
	req, err := createMultipartRequestWithMultipleFiles("files", []string{"first.txt", "second.txt"})
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}

	param := FileParam("files")
	fileHeader, errs := param.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("expected no error but got: %v", errs)
	}
	if fileHeader == nil {
		t.Errorf("expected non-nil file header")
	}
	if fileHeader.Filename != "first.txt" {
		t.Errorf("expected first filename 'first.txt', got %q", fileHeader.Filename)
	}
}

// TestFileParamDifferentFieldNames tests FileParam with different field names
func TestFileParamDifferentFieldNames(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		filename  string
	}{
		{
			name:      "document field",
			fieldName: "document",
			filename:  "report.pdf",
		},
		{
			name:      "attachment field",
			fieldName: "attachment",
			filename:  "image.png",
		},
		{
			name:      "upload field",
			fieldName: "upload",
			filename:  "data.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createMultipartRequest(tt.fieldName, tt.filename, []byte("test content"))
			if err != nil {
				t.Fatalf("failed to create multipart request: %v", err)
			}

			param := FileParam(tt.fieldName)
			fileHeader, errs := param.GetParamValue(req)

			if len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
			if fileHeader.Filename != tt.filename {
				t.Errorf("expected filename %q, got %q", tt.filename, fileHeader.Filename)
			}
		})
	}
}

// TestFileParamEmptyFile tests FileParam with an empty file
func TestFileParamEmptyFile(t *testing.T) {
	req, err := createMultipartRequest("file", "empty.txt", []byte{})
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}

	param := FileParam("file")
	fileHeader, errs := param.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("expected no error for empty file but got: %v", errs)
	}
	if fileHeader == nil {
		t.Errorf("expected non-nil file header for empty file")
	}
	if fileHeader.Size != 0 {
		t.Errorf("expected size 0 for empty file, got %d", fileHeader.Size)
	}
}

// TestFileParamLargeFileName tests FileParam with a long filename
func TestFileParamLargeFileName(t *testing.T) {
	longFilename := "very_long_filename_" + strings.Repeat("a", 200) + ".txt"
	req, err := createMultipartRequest("file", longFilename, []byte("content"))
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}

	param := FileParam("file")
	fileHeader, errs := param.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("expected no error but got: %v", errs)
	}
	if fileHeader.Filename != longFilename {
		t.Errorf("expected filename %q, got %q", longFilename, fileHeader.Filename)
	}
}

// TestFileParamSpecialCharactersFilename tests FileParam with special characters in filename
func TestFileParamSpecialCharactersFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{name: "spaces in filename", filename: "file with spaces.txt"},
		{name: "hyphen in filename", filename: "file-name.txt"},
		{name: "underscore in filename", filename: "file_name.txt"},
		{name: "dot in filename", filename: "file.name.backup.txt"},
		{name: "unicode characters", filename: "файл.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createMultipartRequest("file", tt.filename, []byte("content"))
			if err != nil {
				t.Fatalf("failed to create multipart request: %v", err)
			}

			param := FileParam("file")
			fileHeader, errs := param.GetParamValue(req)

			if len(errs) > 0 {
				t.Errorf("expected no error but got: %v", errs)
			}
			if fileHeader.Filename != tt.filename {
				t.Errorf("expected filename %q, got %q", tt.filename, fileHeader.Filename)
			}
		})
	}
}

// TestFileParamWrongFieldName tests FileParam requesting a field that doesn't exist
func TestFileParamWrongFieldName(t *testing.T) {
	req, err := createMultipartRequest("upload", "file.txt", []byte("content"))
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}

	param := FileParam("wrongfield")
	fileHeader, errs := param.GetParamValue(req)

	if len(errs) == 0 {
		t.Errorf("expected error for wrong field name")
	}
	if fileHeader != nil {
		t.Errorf("expected nil file header for wrong field name")
	}
}

// TestFileParamFileContentPreserved tests that file content can be read correctly
func TestFileParamFileContentPreserved(t *testing.T) {
	content := []byte("This is test file content with special characters: \n\r\t\x00")
	req, err := createMultipartRequest("file", "test.bin", content)
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}

	param := FileParam("file")
	fileHeader, errs := param.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("expected no error but got: %v", errs)
	}

	// Open the file and verify content
	file, err := fileHeader.Open()
	if err != nil {
		t.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	readContent, err := io.ReadAll(file)
	if err != nil {
		t.Errorf("failed to read file content: %v", err)
	}

	if !bytes.Equal(readContent, content) {
		t.Errorf("file content mismatch: expected %v, got %v", content, readContent)
	}
}

// TestFileParamCachedMultipartForm tests that FileParam uses cached multipart form if already parsed
func TestFileParamCachedMultipartForm(t *testing.T) {
	req, err := createMultipartRequest("file", "test.txt", []byte("content"))
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}

	// First call to parse the form
	param := FileParam("file")
	fileHeader1, errs1 := param.GetParamValue(req)

	if len(errs1) > 0 {
		t.Errorf("first call failed: %v", errs1)
	}

	// Second call should use cached form
	fileHeader2, errs2 := param.GetParamValue(req)

	if len(errs2) > 0 {
		t.Errorf("second call failed: %v", errs2)
	}

	if fileHeader1.Filename != fileHeader2.Filename {
		t.Errorf("cached multipart form returned different file")
	}
}

// TestFileParamWithOtherFormFields tests FileParam when multipart form also contains other fields
func TestFileParamWithOtherFormFields(t *testing.T) {
	req, err := createMultipartRequestWithFields(map[string]string{
		"name":  "John Doe",
		"email": "john@example.com",
	}, "profile", "avatar.jpg", []byte("jpg data"))
	if err != nil {
		t.Fatalf("failed to create multipart request: %v", err)
	}

	param := FileParam("profile")
	fileHeader, errs := param.GetParamValue(req)

	if len(errs) > 0 {
		t.Errorf("expected no error but got: %v", errs)
	}
	if fileHeader.Filename != "avatar.jpg" {
		t.Errorf("expected filename 'avatar.jpg', got %q", fileHeader.Filename)
	}
}

// Helper function to create a multipart request with a single file
func createMultipartRequest(fieldName, filename string, content []byte) (*http.Request, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	defer writer.Close()

	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return nil, err
	}

	if _, err := part.Write(content); err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", "/test", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

// Helper function to create a multipart request with multiple files
func createMultipartRequestWithMultipleFiles(fieldName string, filenames []string) (*http.Request, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	defer writer.Close()

	for _, filename := range filenames {
		part, err := writer.CreateFormFile(fieldName, filename)
		if err != nil {
			return nil, err
		}
		if _, err := part.Write([]byte("file content")); err != nil {
			return nil, err
		}
	}

	req := httptest.NewRequest("POST", "/test", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

// Helper function to create a multipart request with both form fields and a file
func createMultipartRequestWithFields(fields map[string]string, fileFieldName, filename string, content []byte) (*http.Request, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	defer writer.Close()

	// Add form fields
	for name, value := range fields {
		if err := writer.WriteField(name, value); err != nil {
			return nil, err
		}
	}

	// Add file
	part, err := writer.CreateFormFile(fileFieldName, filename)
	if err != nil {
		return nil, err
	}

	if _, err := part.Write(content); err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", "/test", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

package filehandler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteToFile(t *testing.T) {
	// Test cases for WriteToFile function
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	// Test creating new file
	err := WriteToFile(testFile, true, false, "test content")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify file was created and contains expected content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Error reading test file: %v", err)
	}

	expected := "test content\n"
	if string(content) != expected {
		t.Errorf("Expected %q, got %q", expected, string(content))
	}
}

func TestRotateFile(t *testing.T) {
	// Test cases for RotateFile function
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "source.txt")
	destFile := filepath.Join(tempDir, "dest.txt")

	// Create source file
	err := WriteToFile(sourceFile, true, false, "rotation test")
	if err != nil {
		t.Fatalf("Error creating source file: %v", err)
	}

	// Test rotation
	err = RotateFile(sourceFile, destFile, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify destination file contains source content
	content, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("Error reading destination file: %v", err)
	}

	if !strings.Contains(string(content), "rotation test") {
		t.Errorf("Destination file should contain rotated content")
	}
}

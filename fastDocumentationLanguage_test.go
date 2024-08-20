package main

import (
	"log"
	"os"
	"testing"
)

func TestGetFdlFilePath(t *testing.T) {
	// Setup test directory
	testDir := "testDir"
	if err := os.Mkdir(testDir, 0755); err != nil {
		log.Panic("Error until creating directory: ", err)
	}
	defer os.RemoveAll(testDir)

	// Create test .fdl file
	testFile := testDir + "/test.fdl"
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.Close()

	// Change to test directory
	originalDir, _ := os.Getwd()
	erro := os.Chdir(testDir)
	if erro != nil {
		log.Fatal(erro)
	}

	defer func(dir string) {
		_ = os.Chdir(dir)
	}(originalDir)

	testFile, _ = os.Getwd()
	// Test function
	files := getFdlFilePath()
	if len(files) != 1 {
		t.Errorf("Expected 1 .fdl file, got %d", len(files))
	}

	fileTest := testFile + "/test.fdl"
	if files[0] != fileTest {
		t.Errorf("Expected file path %s, got %s", fileTest, files[0])
	}
}

func TestConvertFileNameToHTMLFile(t *testing.T) {
	input := "example.fdl"
	expected := "example.html"
	result := convertFileNameToHTMLFile(input)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestFormatInfo(t *testing.T) {
	input := "@info This is an info message."
	expected := "<div style='background-color:#e7f3fe;padding:10px;border-left:6px solid #2196F3;'><strong>Info:" +
		"</strong> This is an info message.</div>"
	result := formatInfo(input)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestGenerateTableOfContents(t *testing.T) {
	sections := map[string]string{
		"section-1": "Introduction",
		"section-2": "Details",
	}
	expected := "<h2>Table of Contents</h2><ul><li><a href='#section-1'>Introduction</a></li><li>" +
		"<a href='#section-2'>Details</a></li></ul>"
	result := generateTableOfContents(sections)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestParseLine(t *testing.T) {
	line := "@title Example Title"
	expected := "<h1>Example Title</h1>"
	result, _, _ := parseLine(line, false, false, map[string]string{})
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestEscapeHTML(t *testing.T) {
	input := "<div>Example</div>"
	expected := "&lt;div&gt;Example&lt;/div&gt;"
	result := escapeHTML(input)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

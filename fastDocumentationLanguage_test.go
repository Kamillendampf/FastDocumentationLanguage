package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestGetFdlFilePath(t *testing.T) {
	// Setup test directory
	testDir := "testDir"
	if err := os.Mkdir(testDir, 0755); err != nil {
		log.Panic("Error creating directory: ", err)
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

	// Test function
	files := getFdlFilePath()
	if len(files) != 1 {
		t.Errorf("Expected 1 .fdl file, got %d", len(files))
	}

	expectedFile := "/home/runner/work/FastDocumentationLanguage/FastDocumentationLanguage/testDir/test.fdl" //This Path is correct for the CI Pipeline
	if files[0] != expectedFile {
		t.Errorf("Expected file path %s, got %s", expectedFile, files[0])
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

func TestFormatWarning(t *testing.T) {
	input := "@warning This is a warning message."
	expected := "<div style='background-color:#ffcccb;padding:10px;border-left:6px solid #f44336;'><strong>Warning:" +
		"</strong> This is a warning message.</div>"
	result := formatWarning(input)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	} else {
		fmt.Println(result != expected)
	}
}

func TestProcessSection(t *testing.T) {
	sections := make(map[string]string)
	input := "@section Introduction"
	expected := "<h2 id='introduction'>Introduction</h2>"
	result := processSection(input, sections)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
	if sections["introduction"] != "Introduction" {
		t.Errorf("Expected section 'Introduction', got %s", sections["introduction"])
	}
}

func TestProcessDefaultLine(t *testing.T) {
	input := "This is a regular line."
	expected := "This is a regular line."
	result := processDefaultLine(input, false, false)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test within code block
	input = "    code block line"
	expected = "    code block line\n"
	result = processDefaultLine(input, true, false)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test within table (should return empty)
	input = "| Table | Line |"
	expected = ""
	result = processDefaultLine(input, false, true)
	if result != expected {
		t.Errorf("Expected empty string, got %s", result)
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

		// Test @author
		line = "@author John Doe"
		expected = "<p>Author: John Doe</p>"
		result, _, _ = parseLine(line, false, false, map[string]string{})
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	}
	// Test @date
	line = "@date 2023-08-24"
	expected = "<p>Date: 2023-08-24</p>"
	result, _, _ = parseLine(line, false, false, map[string]string{})
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test @abstract
	line = "@abstract This is an abstract."
	expected = "<h2>Abstract</h2><p>"
	result, _, _ = parseLine(line, false, false, map[string]string{})
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test @note
	line = "@note This is a note."
	expected = "<p><em>Note:</em> This is a note.</p>"
	result, _, _ = parseLine(line, false, false, map[string]string{})
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test @code block
	line = "@code"
	expected = "<pre><code>"
	result, inCodeBlock, _ := parseLine(line, false, false, map[string]string{})
	if result != expected || !inCodeBlock {
		t.Errorf("Expected %s with inCodeBlock=true, got %s with inCodeBlock=%v", expected, result, inCodeBlock)
	}

	// Test @endcode block
	line = "@endcode"
	expected = "</code></pre>"
	result, inCodeBlock, _ = parseLine(line, true, false, map[string]string{})
	if result != expected || inCodeBlock {
		t.Errorf("Expected %s with inCodeBlock=false, got %s with inCodeBlock=%v", expected, result, inCodeBlock)
	}

	// Test @table and @row
	line = "@table"
	expected = "<table border='1'>"
	result, _, inTable := parseLine(line, false, false, map[string]string{})
	if result != expected || !inTable {
		t.Errorf("Expected %s with inTable=true, got %s with inTable=%v", expected, result, inTable)
	}

	line = "@row cell1 | cell2"
	expected = "<tr><td>cell1</td><td>cell2</td></tr>"
	result, _, inTable = parseLine(line, false, true, map[string]string{})
	if result != expected && inTable {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test @endtable
	line = "@endtable"
	expected = "</table>"
	result, _, inTable = parseLine(line, false, true, map[string]string{})
	if result != expected || inTable {
		t.Errorf("Expected %s with inTable=false, got %s with inTable=%v", expected, result, inTable)
	}

	// Test @tbc (to be continued)
	line = "@tbc"
	expected = "to be continued"
	result, _, _ = parseLine(line, false, false, map[string]string{})
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

func TestCreateOrCleanOutputDir(t *testing.T) {
	// Setup test directory
	testDir := "fdlDocumentation"
	if err := os.Mkdir(testDir, 0755); err != nil && !os.IsExist(err) {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	createOrCleanOutputDir()

	// Check if directory exists
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Errorf("Expected directory %s to exist, but it does not", testDir)
	}
}

func TestProcessFileDefaultMode(t *testing.T) {
	// Setup: Create a temporary directory for testing.
	tempDir := t.TempDir()

	// Override the current working directory with the temporary directory.
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory to temp dir: %v", err)
	}

	// Create a mock .fdl file for testing.
	mockFdlFile := "test.fdl"
	f, err := os.Create(mockFdlFile)
	if err != nil {
		t.Fatalf("Failed to create mock .fdl file: %v", err)
	}
	defer f.Close()

	// Write sample content to the mock .fdl file.
	_, err = f.WriteString("@title Sample Title\n@section Sample Section\n@info This is an info message.")
	if err != nil {
		t.Fatalf("Failed to write to mock .fdl file: %v", err)
	}

	// Run the processFileDefaultMode function.
	processFileDefaultMode()

	// Verify the output.
	outputFile := filepath.Join(tempDir, "fdlDocumentation", "test.html")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Expected output file %s to be created, but it does not exist", outputFile)
	}

	// Optionally, read and verify the content of the generated HTML file.
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := "<h1><h2>Table of Contents</h2><ul><li><a href='#sample-section'>Sample Section</a></li></ul>\nSample Title</h1>\n<h2 id='sample-section'>Sample Section</h2>\n<div style='background-color:#e7f3fe;padding:10px;border-left:6px solid #2196F3;'><strong>Info:</strong> This is an info message.</div>\n"
	if string(content) != expectedContent {
		t.Errorf("Expected content:\n%s\nGot:\n%s", expectedContent, content)
	}

	// Verify the index.html file.
	indexFile := filepath.Join(tempDir, "fdlDocumentation", "index.html")
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		t.Errorf("Expected index file %s to be created, but it does not exist", indexFile)
	}
}

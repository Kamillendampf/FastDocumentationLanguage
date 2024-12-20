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
	files := getFilePath(".fdl")
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
	expected := "This is a regular line.<br>"
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

	createOrCleanOutputDir("/fdlDocumentation")

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
	processFiles()

	// Verify the output.
	outputFile := filepath.Join(tempDir, "documentation", "test.html")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Expected output file %s to be created, but it does not exist", outputFile)
	}

	// Optionally, read and verify the content of the generated HTML file.
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := "<h1><h2>Table of Contents</h2><ul><li><a href='#sample-section'>Sample Section</a></li></ul>\nSample Title</h1>\n<h2 id='sample-section'>Sample Section</h2>\n<div style='background-color:#e7f3fe;padding:10px;border-left:6px solid #2196F3;'><strong>Info:</strong> This is an info message.</div>\n<style>.example-box {border: 2px solid black;padding: 10px;margin: 20px 0;border-radius: 5px;background-color: #f9f9f9;position: relative;overflow: hidden;}.example-title {font-weight: bold;margin: 0;padding: 5px 10px;background-color: #e0e0e0;border-bottom: 2px solid black;position: absolute;top: 0;left: 0;width: 100%;box-sizing: border-box;}.example-content {padding-top: 40px;}</style>"
	if string(content) != expectedContent {
		t.Errorf("Expected content:\n%s\nGot:\n%s", expectedContent, content)
	}

	// Verify the index.html file.
	indexFile := filepath.Join(tempDir, "documentation", "index.html")
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		t.Errorf("Expected index file %s to be created, but it does not exist", indexFile)
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		line                       string
		inCodeBlock                bool
		inTable                    bool
		inList                     bool
		isUseCaseORExample         bool
		expectedOutput             string
		expectedInCodeBlock        bool
		expectedInTable            bool
		expectedInList             bool
		expectedIsUseCaseORExample bool
	}{
		// Testcases for different types of lines and states

		// @title
		{"@title My Title", false, false, false, false, "<h1>My Title</h1>", false, false, false, false},

		// @author
		{"@author John Doe", false, false, false, false, "<p>Author: John Doe</p>", false, false, false, false},

		// @date
		{"@date 2024-08-25", false, false, false, false, "<p>Date: 2024-08-25</p>", false, false, false, false},

		// @abstract
		{"@abstract", false, false, false, false, "<h2>Abstract</h2><p>", false, false, false, false},

		// @info
		{"@info Information", false, false, false, false, formatInfo("@info Information"), false, false, false, false},

		// @warning
		{"@warning Warning Message", false, false, false, false, formatWarning("@warning Warning Message"), false, false, false, false},

		// @note
		{"@note This is a note", false, false, false, false, "<p><em>Note:</em> This is a note</p>", false, false, false, false},

		// @code and @endcode
		{"@code", false, false, false, false, "<div class='example-box'><div class='example-title'>Code:</div><div class='example-content'><pre><code>", true, false, false, false},
		{"@endcode", true, false, false, false, "</code></pre></div></div>", false, false, false, false},
		{"@code", false, false, false, true, "<pre><code>", true, false, false, true},
		{"@endcode", true, false, false, true, "</code></pre>", false, false, false, true},

		// @tbc
		{"@tbc", false, false, false, false, "", false, false, false, false},

		// @table and @row
		{"@table", false, false, false, false, "<table border='1'>", false, true, false, false},
		{"@row cell1|cell2|cell3", false, true, false, false, "<tr><td>cell1</td><td>cell2</td><td>cell3</td></tr>", false, true, false, false},
		{"@endtable", false, true, false, false, "</table>", false, false, false, false},

		// @version
		{"@version 1.0.0", false, false, false, false, "<p><em>Version:</em> 1.0.0</p>", false, false, false, false},

		// @since
		{"@since 2024", false, false, false, false, "<p><em>Since:</em> 2024</p>", false, false, false, false},

		// @deprecated
		{"@deprecated", false, false, false, false, "<strong><em style='color:red;'>Deprecated!</em></strong>", false, false, false, false},

		// @param
		{"@param param1|param2", false, false, false, false, "<p><b>Parameters</b></p><p>param1</p><p>param2</p>", false, false, false, false},

		// @return
		{"@return return1|return2", false, false, false, false, "<p><b>Return:</b></p><p>return1</p><p>return2</p>", false, false, false, false},

		// @list
		{"@list -n", false, false, false, false, "<ol>", false, false, true, false},
		{"@list", false, false, false, false, "<ul>", false, false, true, false},

		// @item
		{"@item List item", false, false, true, false, "<li>List item</li>", false, false, true, false},

		// @endlist
		{"@endlist", false, false, true, false, "</ul></ol>", false, false, false, false},

		// @tip
		{"@tip This is a tip", false, false, false, false, formatTip("This is a tip"), false, false, false, false},

		// @todo
		{"@todo This is a todo", false, false, false, false, "<p><em>TODO:</em> This is a todo</p>", false, false, false, false},

		// @example
		{"@example", false, false, false, false, "<div class='example-box'><div class='example-title'>Example:</div><div class='example-content'>", false, false, false, true},
		{"@endexample", false, false, false, true, "</div></div>", false, false, false, false},

		// @usecase
		{"@usecase", false, false, false, false, "<div class='example-box'><div class='example-title'>UseCase:</div><div class='example-content'>", false, false, false, true},
		{"@endusecase", false, false, false, true, "</div></div>", false, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got, gotInCodeBlock, gotInTable, gotInList, gotIsUseCaseORExample := parseLine(tt.line, tt.inCodeBlock, tt.inTable, tt.inList, tt.isUseCaseORExample, nil)
			if got != tt.expectedOutput {
				t.Errorf("parseLine() = %v, want %v", got, tt.expectedOutput)
			}
			if gotInCodeBlock != tt.expectedInCodeBlock {
				t.Errorf("parseLine() inCodeBlock = %v, want %v", gotInCodeBlock, tt.expectedInCodeBlock)
			}
			if gotInTable != tt.expectedInTable {
				t.Errorf("parseLine() inTable = %v, want %v", gotInTable, tt.expectedInTable)
			}
			if gotInList != tt.expectedInList {
				t.Errorf("parseLine() inList = %v, want %v", gotInList, tt.expectedInList)
			}
			if gotIsUseCaseORExample != tt.expectedIsUseCaseORExample {
				t.Errorf("parseLine() isUseCaseORExample = %v, want %v", gotIsUseCaseORExample, tt.expectedIsUseCaseORExample)
			}
		})
	}
}

// TestGetFlagsFromCli testet die getFlagsFromCli Funktion
func TestGetFlagsFromCli(t *testing.T) {
	tests := []struct {
		args        []string
		expected    flag
		description string
	}{
		{
			args:        []string{"cmd", "--file-extension=.txt", "--directory=/docs"},
			expected:    flag{FileExtension: ".txt", Directory: "/docs"},
			description: "Valid flags provided",
		},
		{
			args:        []string{"cmd"},
			expected:    flag{FileExtension: ".fdl", Directory: "/documentation"},
			description: "Default flags used",
		},
		{
			args:        []string{"cmd", "-fe=.md"},
			expected:    flag{FileExtension: ".md", Directory: "/documentation"},
			description: "Short flag for file extension",
		},
		{
			args:        []string{"cmd", "-dir=./custom"},
			expected:    flag{FileExtension: ".fdl", Directory: "./custom"},
			description: "Short flag for directory",
		},
		{
			args:        []string{"cmd", "--file-extension=invalid"},
			expected:    flag{FileExtension: "invalid", Directory: "/documentation"},
			description: "Invalid file extension flag",
		},
	}

	// Speichern des ursprünglichen os.Args
	originalArgs := os.Args

	// Testen der Fälle
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Setze os.Args auf die Testdaten
			os.Args = tt.args

			// Rufe die Funktion auf
			got := getFlagsFromCli()

			// Überprüfe, ob das Ergebnis den Erwartungen entspricht
			if got != tt.expected {
				t.Errorf("getFlagsFromCli() = %+v; want %+v", got, tt.expected)
			}
		})
	}

	// Stelle os.Args auf die ursprünglichen Werte zurück
	os.Args = originalArgs
}

func TestGetFilePath(t *testing.T) {
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

	// Erstelle einige Testdateien
	testFiles := []struct {
		name     string
		content  string
		expected bool // Gibt an, ob die Datei in den Ergebnissen erwartet wird
	}{
		{"test1.fdl", "Inhalt von test1", true},
		{"test2.txt", "Inhalt von test2", false},
		{"test3.fdl", "Inhalt von test3", true},
		{"test4.md", "Inhalt von test4", false},
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file.name)
		err := os.WriteFile(filePath, []byte(file.content), 0644)
		if err != nil {
			t.Fatalf("Could not create test file %s: %v", file.name, err)
		}
	}

	// Teste die getFilePath Funktion mit der .fdl Erweiterung
	got := getFilePath(".fdl")

	// Überprüfe, ob die richtigen Dateipfade zurückgegeben werden
	expectedPaths := []string{
		filepath.Join(tempDir, "test1.fdl"),
		filepath.Join(tempDir, "test3.fdl"),
	}

	// Überprüfe, ob die Rückgabe der Funktion die erwarteten Pfade enthält
	for _, expectedPath := range expectedPaths {
		if !contains(got, expectedPath) {
			t.Errorf("Expected path %s not found in results: %+v", expectedPath, got)
		}
	}

	// Überprüfe, ob keine unerwarteten Pfade zurückgegeben werden
	for _, file := range testFiles {
		if file.expected && !contains(got, filepath.Join(tempDir, file.name)) {
			t.Errorf("Expected path %s not found in results: %+v", filepath.Join(tempDir, file.name), got)
		} else if !file.expected && contains(got, filepath.Join(tempDir, file.name)) {
			t.Errorf("Unexpected path found: %s", filepath.Join(tempDir, file.name))
		}
	}
}

// Hilfsfunktion, um zu überprüfen, ob ein Slice einen bestimmten String enthält
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

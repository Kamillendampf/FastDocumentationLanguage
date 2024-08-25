package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getFdlFilePath() []string {
	var pathSlices []string

	cwd, err := os.Getwd()
	if err != nil {
		log.Panic("Error until reading the directories: ", err)
		return nil
	}

	_ = filepath.WalkDir(cwd, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".fdl") {
			pathSlices = append(pathSlices, path)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return pathSlices
}

func convertFileNameToHTMLFile(fileName string) string {
	filenameSlices := strings.Split(fileName, ".")
	return filenameSlices[0] + ".html"
}

func formatInfo(line string) string {
	return fmt.Sprintf("<div style='background-color:#e7f3fe;padding:10px;border-left:6px solid #2196F3;'>"+
		"<strong>Info:</strong> %s</div>", strings.TrimSpace(line[5:]))
}

func formatWarning(line string) string {
	return fmt.Sprintf("<div style='background-color:#ffcccb;padding:10px;border-left:6px solid #f44336;'>"+
		"<strong>Warning:</strong> %s</div>", strings.TrimSpace(line[8:]))
}

func formatTip(line string) string {
	return fmt.Sprintf("<div style='background-color:#8fbc8f;padding:10px;border-left:6px solid #6e8b3d;'>"+
		"<strong>Tip:</strong> %s</div>", strings.TrimSpace(line[8:]))
}

func processSection(line string, sections map[string]string) string {
	sectionTitle := strings.TrimSpace(line[8:])
	sectionID := strings.ToLower(strings.ReplaceAll(sectionTitle, " ", "-"))
	sections[sectionID] = sectionTitle
	return fmt.Sprintf("<h2 id='%s'>%s</h2>", sectionID, sectionTitle)
}

func processDefaultLine(line string, inCodeBlock bool, inTable bool) string {
	if inCodeBlock {
		return fmt.Sprintf("%s\n", line)
	}
	if inTable {
		return ""
	}
	return line
}

func generateTableOfContents(sections map[string]string) string {
	var tocBuilder strings.Builder
	if len(sections) > 0 {
		tocBuilder.WriteString("<h2>Table of Contents</h2><ul>")
		for id, title := range sections {
			tocBuilder.WriteString(fmt.Sprintf("<li><a href='#%s'>%s</a></li>", id, title))
		}
		tocBuilder.WriteString("</ul>")
	}
	return tocBuilder.String()
}

func escapeHTML(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(input, "<", "&lt;"), ">", "&gt;")
}

func parseLine(line string, inCodeBlock bool, inTable bool, inList bool, isUseCaseORExample bool, sections map[string]string) (string, bool, bool, bool, bool) {
	switch {
	case strings.HasPrefix(line, "@title") && !inCodeBlock:
		return fmt.Sprintf("<h1>%s</h1>", strings.TrimSpace(line[6:])), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@author") && !inCodeBlock:
		return fmt.Sprintf("<p>Author: %s</p>", strings.TrimSpace(line[7:])), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@date") && !inCodeBlock:
		return fmt.Sprintf("<p>Date: %s</p>", strings.TrimSpace(line[5:])), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@abstract") && !inCodeBlock:
		return "<h2>Abstract</h2><p>", inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@info"):
		return formatInfo(line), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@warning"):
		return formatWarning(line), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@section") && !inCodeBlock && !isUseCaseORExample:
		return processSection(line, sections), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@note"):
		return fmt.Sprintf("<p><em>Note:</em> %s</p>", strings.TrimSpace(line[5:])), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@code") && !inCodeBlock:
		if !isUseCaseORExample {
			return "<div class='example-box'><div class='example-title'>Code:</div><div class='example-content'><pre><code>", !inCodeBlock, inTable, inList, isUseCaseORExample
		} else {
			return "<pre><code>", !inCodeBlock, inTable, inList, isUseCaseORExample
		}
	case strings.HasPrefix(line, "@endcode") && inCodeBlock:
		if !isUseCaseORExample {
			return "</code></pre></div></div>", !inCodeBlock, inTable, inList, isUseCaseORExample
		} else {
			return "</code></pre>", !inCodeBlock, inTable, inList, isUseCaseORExample
		}
	case strings.HasPrefix(line, "@tbc"):
		return "", inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@table"):
		return "<table border='1'>", inCodeBlock, !inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@row"):
		cells := strings.Split(strings.TrimSpace(line[4:]), "|")
		var rowBuilder strings.Builder
		rowBuilder.WriteString("<tr>")
		for _, cell := range cells {
			rowBuilder.WriteString(fmt.Sprintf("<td>%s</td>", strings.TrimSpace(cell)))
		}
		rowBuilder.WriteString("</tr>")
		return rowBuilder.String(), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@endtable"):
		return "</table>", inCodeBlock, !inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@version"):
		return fmt.Sprintf("<p><em>Version:</em> %s</p>", strings.TrimSpace(line[8:])), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@since"):
		return fmt.Sprintf("<p><em>Since:</em> %s</p>", strings.TrimSpace(line[6:])), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@deprecated"):
		return "<strong><em style='color:red;'>Deprecated!</em></strong>", inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@param"):
		params := strings.Split(strings.TrimSpace(line[6:]), "|")
		var rowBuilder strings.Builder
		rowBuilder.WriteString("<p><b>Parameters</b></p>")
		for _, param := range params {
			rowBuilder.WriteString(fmt.Sprintf("<p>%s</p>", param))
		}
		return rowBuilder.String(), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@return"):
		params := strings.Split(strings.TrimSpace(line[7:]), "|")
		var rowBuilder strings.Builder
		rowBuilder.WriteString("<p><b>Return:</b></p>")
		for _, param := range params {
			rowBuilder.WriteString(fmt.Sprintf("<p>%s</p>", param))
		}
		return rowBuilder.String(), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@list"):
		line = strings.TrimSpace(line[5:])
		switch {
		case strings.HasPrefix(line, "-n"):
			return "<ol>", inCodeBlock, inTable, !inList, isUseCaseORExample
		default:
			return "<ul>", inCodeBlock, inTable, !inList, isUseCaseORExample
		}
	case strings.HasPrefix(line, "@item") && inList:
		return fmt.Sprintf("<li>%s</li>", strings.TrimSpace(line[5:])), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@endlist"):
		return "</ul></ol>", inCodeBlock, inTable, !inList, isUseCaseORExample
	case strings.HasPrefix(line, "@tip"):
		return formatTip(strings.TrimSpace(line[4:])), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@todo"):
		return fmt.Sprintf("<p><em>TODO:</em> %s</p>", strings.TrimSpace(line[5:])), inCodeBlock, inTable, inList, isUseCaseORExample
	case strings.HasPrefix(line, "@example"):
		return "<div class='example-box'><div class='example-title'>Example:</div><div class='example-content'>", inCodeBlock, inTable, inList, !isUseCaseORExample
	case strings.HasPrefix(line, "@endexample"):
		return "</div></div>", inCodeBlock, inTable, inList, !isUseCaseORExample
	case strings.HasPrefix(line, "@usecase"):
		return "<div class='example-box'><div class='example-title'>UseCase:</div><div class='example-content'>", inCodeBlock, inTable, inList, !isUseCaseORExample
	case strings.HasPrefix(line, "@endusecase"):
		return "</div></div>", inCodeBlock, inTable, inList, !isUseCaseORExample
	default:
		return processDefaultLine(line, inCodeBlock, inTable), inCodeBlock, inTable, inList, isUseCaseORExample
	}
}

func createOrCleanOutputDir() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic("Error until reading the directories: ", err)
	}

	outputPath := cwd + "/fdlDocumentation"

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		log.Println("The directory don't exist, it is created")
		if err := os.Mkdir(outputPath, 0777); err != nil {
			log.Panic("The System couldn't create the directory: ", err)
		}
	} else if err != nil {
		log.Panic("error:", err)
	} else {
		log.Println("The directory is cleaned and created again")
		err := os.RemoveAll(outputPath)
		if err != nil {
			log.Panic("Error until delete :", err)
		}
		if err := os.Mkdir(outputPath, 0777); err != nil {
			return
		}
	}
}

func outputStream(finaleFormattedHTML string, filename string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic("Error until reading the directories: ", err)
	}

	outputPath := cwd + "/fdlDocumentation"

	if err := os.Chdir(outputPath); err != nil {
		fmt.Println(cwd)
		log.Panic("Can't change directory: ", err)
	}

	outputFile, err := os.Create(filename)
	if err != nil {
		log.Panic("Can't write in output file: ", err)
	}

	_, _ = outputFile.WriteString(finaleFormattedHTML)
	if err != nil {
		return
	}
	outputFile.Close()
	if err := os.Chdir(cwd); err != nil {
		fmt.Println(cwd)
		log.Panic("Can't change directory: ", err)
	}
}

func creatIndex(tableofContent []string) {
	table := "<html><body><h1>Documentation <br> Table Of Content</h1><ul>"
	for index, content := range tableofContent {
		chapterName := strings.Split(content, ".")
		chapterNumber := index + 1

		chapterFullName := strconv.Itoa(chapterNumber) + " " + chapterName[0]
		table = table + "<li> <a href='" + content + "'>" + chapterFullName + "</a></li>"
	}
	table = table + "</ul></body></html>"

	outputStream(table, "index.html")

}
func processStyling() string {
	return "<style>.example-box {border: 2px solid black;padding: 10px;margin: 20px 0;" +
		"border-radius: 5px;background-color: #f9f9f9;position: relative;overflow: hidden;}" +
		".example-title {font-weight: bold;margin: 0;padding: 5px 10px;background-color: #e0e0e0;" +
		"border-bottom: 2px solid black;position: absolute;top: 0;left: 0;width: 100%;box-sizing: border-box;}" +
		".example-content {padding-top: 40px;}</style>"
}
func processFileDefaultMode() {
	var mainTableOfContent []string

	createOrCleanOutputDir()

	for _, path := range getFdlFilePath() {
		var output strings.Builder
		sections := make(map[string]string)
		inCodeBlock := false
		inTable := false
		inList := false
		isUsecaseOrExample := false
		mainTableOfContent = append(mainTableOfContent, convertFileNameToHTMLFile(filepath.Base(path)))
		currentFile := mainTableOfContent[len(mainTableOfContent)-1]
		fdlFile, err := os.Open(path)
		if err != nil {
			log.Panic(err)
			return
		}

		scanner := bufio.NewScanner(fdlFile)
		for scanner.Scan() {
			line := scanner.Text()
			line, inCodeBlock, inTable, inList, isUsecaseOrExample = parseLine(escapeHTML(line), inCodeBlock, inTable, inList, isUsecaseOrExample, sections)
			if line != "" {
				output.WriteString(line)
				output.WriteString("\n")
			}
		}

		if err := scanner.Err(); err != nil {
			log.Panic("Error reading file:", err)
		}

		toc := generateTableOfContents(sections)
		finalOutput := strings.Replace(output.String(), "<h1>", "<h1>"+toc+"\n", 1)
		finalOutput = finalOutput + processStyling()

		outputStream(finalOutput, currentFile)
		fdlFile.Close()
	}

	creatIndex(mainTableOfContent)
}

func main() {
	processFileDefaultMode()
}

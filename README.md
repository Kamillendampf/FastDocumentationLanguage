# Custom Markup to HTML Converter

## Status

[![Go](https://github.com/Kamillendampf/FastDocumentationLanguage/actions/workflows/ci.yml/badge.svg)](https://github.com/Kamillendampf/FastDocumentationLanguage/actions/workflows/ci.yml)

## Overview

This program is a simple and efficient tool that converts plain text files using a custom markup language into well-structured HTML documents. The custom markup language is designed to facilitate the creation of HTML content by using intuitive and easy-to-remember tags. The program reads an input text file, processes each line according to the markup commands, and outputs a fully formatted HTML document.

## Features

- **Title and Metadata**: Convert titles, authors, and dates into HTML headers and paragraphs.
- **Sectioning**: Define sections and subsections of the document to create a structured and navigable layout.
- **Alerts and Notes**: Highlight important information and warnings using styled HTML divs.
- **Code Blocks**: Insert and format code snippets, with automatic HTML escaping to preserve the appearance of code.
- **Table of Contents**: Automatically generate a table of contents based on the sections defined in the document.

## Supported Markup Commands

The following commands are supported by the custom markup language:

- `@title <Title>`: Defines the main title of the document.
- `@author <Author>`: Specifies the author of the document.
- `@date <Date>`: Adds a date to the document.
- `@abstract`: Begins an abstract section.
- `@section <Section Title>`: Starts a new section with the specified title.
- `@info <Information>`: Highlights important information with a styled block.
- `@warning <Warning>`: Emphasizes a warning message with a styled block.
- `@note <Note>`: Adds a note in italicized text.
- `@code`: Begins a code block.
- `@endcode`: Ends the current code block.
- `@tbc`: Placeholder for content to be continued (no output).
- `@table` :  Starts the definition of a table. This command creates a <table> element in the HTML output.
- `@row <Header 1> | <Header 2> | <Header 3>` : Defines a new row in the table. The row content should be separated by the | character, which will be converted into <td> (table cell) elements. Each @row creates a <tr> (table row) in the HTML.
- `@row Row 1 Col 1 | Row 1 Col 2 | Row 1 Col 3` :
- `@endtable` :  Ends the table definition. This command closes the <table> element in the HTML output.


    ## How It Works

1. **Input Parsing**: The program reads the input text file line by line.
2. **Line Processing**: Each line is processed based on the markup commands. The program determines if the line should be converted into an HTML tag, a section heading, or part of a code block.
3. **HTML Escaping**: All special HTML characters, such as `<` and `>`, are automatically escaped to prevent them from being interpreted as HTML code.
4. **Output Generation**: The processed content is then written into an HTML structure. If sections are defined, a table of contents is generated and inserted into the document.

## Usage

To use the converter:

1. Write your document, using the custom markup language.
2. Run the program.
3. The program outputs the formatted HTML document.

### Example

**Input File (`input.fdl`):**

```text
@title Example Document
@author John Doe
@date 2024-08-18

@abstract
This is a brief summary of the document.

@section Introduction
This section introduces the topic.

@info This is some important information.

@code
func example() {
    fmt.Println("<Hello, World!>")
}
@endcode

@warning This is a critical warning.
```
## Generated HTML

```html
<h1>Example Document</h1>
<p><strong>Author:</strong> John Doe</p>
<p><strong>Date:</strong> 2024-08-18</p>
<h2>Abstract</h2>
<p>This is a brief summary of the document.</p>
<h2 id='introduction'>Introduction</h2>
<p>This section introduces the topic.</p>
<div style='background-color:#e7f3fe;padding:10px;border-left:6px solid #2196F3;'><strong>Info:</strong> This is some important information.</div>
<pre><code>func example() {
    fmt.Println("&lt;Hello, World!&gt;")
}
</code></pre>
<div style='background-color:#ffcccb;padding:10px;border-left:6px solid #f44336;'><strong>Warning:</strong> This is a critical warning.</div>
```
## Installation

You have to options for the installation. 

1. Build the project by your own:
    ```Git
    git clone https://github.com/Kamillendampf/FastDocumentationLanguage.git
    cd markup-to-html
    go build
    ```
2. Download the .exe file from the release section:

    [Download](https://github.com/Kamillendampf/FastDocumentationLanguage/releases) the latest version.

    **Note:** No matter which option you choose, you must always place the file in the root directory of your project.

    ## Running the Program

    You can run the program by specifying the input file as follows:

    ```CMD
    ./FastDocumentationLanguage.exe 
    ```

    otherwise you could run it by double click. (First option is necessary if you would like to use it in to a pipeline for automated document generation)

    The output is directly printed in to the Files.

    ## Contributing

    Contributions are welcome! Feel free to open issues or submit pull requests to improve the functionality or add new features.

    ## License

    This project is licensed under the MIT License. See the `LICENSE` file for details.
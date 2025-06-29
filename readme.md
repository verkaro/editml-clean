# editml-clean

`editml-clean` is a command-line utility written in Go that parses text formatted with EditML (Editorial Markup Language) and outputs a clean, human-readable version. It processes additions, deletions, comments, highlights, and structural edits (moves/copies) to produce the final "Clean View" of a document.

This tool is a reference implementation based on the `github.com/verkaro/editml-go` library.

---

## Specification Note

This application adheres to the v0.1 version of the EditML specification as implemented by the underlying `verkaro/editml-go` library. It supports all specified core features but may not handle all advanced nesting or error conditions found in the full EditML specification.

---

## Features

-   Processes EditML from files or standard input (`stdin`).
-   Outputs clean plain text to files or standard output (`stdout`).
-   Supports all standard EditML operations:
    -   **Additions:** `{+text+}`
    -   **Deletions:** `{-text-}`
    -   **Highlights:** `{=text=}`
    -   **Comments:** `{>text<}`
    -   **Moves:** `{move~text~TAG}` and `{move:TAG}`
    -   **Copies:** `{copy~text~TAG}` and `{copy:TAG}`
-   Command-line flags for controlling output, debugging, and strictness.

---

## Installation

To use `editml-clean`, you need to have Go (version 1.21 or later) installed on your system.

1.  **Clone the repository (or download the source files):**
    Ensure you have `main.go`, `go.mod`, etc., in a local directory.

2.  **Build the binary:**
    Navigate to the project's root directory in your terminal and run the `go build` command. This will create the `editml-clean` executable in your current directory.

    ```bash
    go build -o editml-clean .
    ```

3.  **Run the tool:**
    You can now run the tool directly from that directory. For system-wide access, you can move the `editml-clean` binary to a directory in your system's `PATH` (e.g., `/usr/local/bin`).

    ```bash
    ./editml-clean --version
    ```

---

## Usage

The tool can read from `stdin` or a file and write to `stdout` or a file.

### Command-Line Interface

```
Usage:
  editml-clean [flags] [input-file]
```

### Flags

| Flag                | Shorthand | Description                                                   |
| ------------------- | --------- | ------------------------------------------------------------- |
| `--version`         |           | Print the application version and exit.                       |
| `--output <file>`   | `-o <file>` | Write output to the specified file instead of stdout.         |
| `--debug`           |           | Emit any parse/transform issues (warnings/errors) to stderr.  |
| `--strict`          |           | Treat any warnings as fatal errors (exits with a non-zero code).|
| `--help`            | `-h`      | Show usage information.                                       |

### Positional Argument

-   `input-file` (optional): The path to an EditML file to process. If this argument is omitted, the tool will read from `stdin`.

---

## Examples

### 1. Cleaning a File

Given a file named `draft.editml` with the content:
```editml
This is a{- really-} {=great=} document.{>Or is it?<}
```

Run the following command:

```bash
./editml-clean draft.editml
```

**Output:**
```
This is a great document.
```

### 2. Using stdin and stdout

You can pipe content directly into the tool.

```bash
echo "Testing a choice{>comment<} and{move~ text to move~t1} here is{move:t1}." | ./editml-clean
```

**Output:**
```
Testing a choice and here is text to move.
```

### 3. Writing to an Output File

To save the cleaned content to a new file:

```bash
./editml-clean --output final.txt draft.editml
```
This will create a file named `final.txt` with the cleaned output.

### 4. Debugging a File

If you have a file that isn't processing correctly, the `--debug` flag can provide more information.

```bash
# Assuming a file with an unclosed tag
echo "Here is an unclosed addition {+" | ./editml-clean --debug
```

**Example Stderr Output:**
```
[Error] Parsing error: unclosed edit tag (L0:C0)
```

---

## Development & Testing

The project includes a full suite of unit and integration tests. To run them, navigate to the project root and execute:

```bash
go test
```

---

## Acknowledgements

This utility was implemented by Google's Gemini based on an initial specification. The final, robust application was achieved through a collaborative process of iterative development, testing, and refinement.


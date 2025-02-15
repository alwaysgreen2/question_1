package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// normalizePath converts a relative path (which may include occurrences of "./", "../",
// and multiple consecutive slashes) to its corresponding normalized absolute path
// based on the given current directory.
// In addition, for non-special segments (not exactly "." or ".."), it trims trailing
// dots and question marks.
func normalizePath(currDir, relPath string) string {
	// Process the current directory: split by "/" and ignore empty or "."
	baseStack := []string{}
	for _, seg := range strings.Split(currDir, "/") {
		if seg == "" || seg == "." {
			continue
		}
		if seg == ".." {
			if len(baseStack) > 0 {
				baseStack = baseStack[:len(baseStack)-1]
			}
			continue
		}
		baseStack = append(baseStack, seg)
	}

	// Remove any leading slashes from the relative path so that even inputs like "////..."
	// are treated as relative.
	relPath = strings.TrimLeft(relPath, "/")
	relStack := []string{}
	for _, seg := range strings.Split(relPath, "/") {
		if seg == "" {
			continue
		}
		// If the segment is exactly "." or "..", handle them
		if seg == "." {
			continue
		}
		if seg == ".." {
			if len(relStack) > 0 {
				relStack = relStack[:len(relStack)-1]
			} else {
				// Ignore ".." if nothing to pop from relative stack.
			}
			continue
		}
		// For any other segment, trim trailing '.' and '?' characters.
		trimmed := strings.TrimRight(seg, ".?")
		// Only add if trimming didn't result in an empty string.
		if trimmed != "" {
			relStack = append(relStack, trimmed)
		}
	}

	// Combine the base directory and the processed relative path.
	finalStack := append(baseStack, relStack...)
	result := strings.Join(finalStack, "/")
	// If the current directory was absolute (i.e. started with '/'), add a leading slash.
	if len(currDir) > 0 && currDir[0] == '/' {
		result = "/" + result
	}
	return result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter current directory (or 'q' to quit): ")
		if !scanner.Scan() {
			break
		}
		currDir := strings.TrimSpace(scanner.Text())
		if currDir == "q" {
			break
		}

		fmt.Print("Enter relative path (or 'q' to quit): ")
		if !scanner.Scan() {
			break
		}
		relPath := strings.TrimSpace(scanner.Text())
		if relPath == "q" {
			break
		}

		absPath := normalizePath(currDir, relPath)
		fmt.Printf("The absolute path is: %s\n\n", absPath)
	}
}

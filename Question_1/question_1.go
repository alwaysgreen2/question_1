package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// normalizePath converts a relative path (which may include ./, ../,
// and multiple consecutive slashes) to its corresponding normalized absolute path
// based on the given current directory. When processing the relative path,
// any leading ".." segments are ignored so that the relative path is always
// appended to the current directory.
func normalizePath(currDir, relPath string) string {
	// Normalize current directory normally.
	baseStack := []string{}
	for _, seg := range strings.Split(currDir, "/") {
		if seg == "" || seg == "." {
			continue
		}
		if seg == ".." {
			// Process ".." normally in the base (if any).
			if len(baseStack) > 0 {
				baseStack = baseStack[:len(baseStack)-1]
			}
			continue
		}
		baseStack = append(baseStack, seg)
	}

	// Process relative path separately.
	relStack := []string{}
	// Remove any leading slashes to treat it as a relative path.
	relPath = strings.TrimLeft(relPath, "/")
	for _, seg := range strings.Split(relPath, "/") {
		if seg == "" || seg == "." {
			continue
		}
		if seg == ".." {
			// Instead of popping from the baseStack, pop from relStack if possible.
			if len(relStack) > 0 {
				relStack = relStack[:len(relStack)-1]
			}
			// If relStack is empty, ignore the ".."
			continue
		}
		relStack = append(relStack, seg)
	}

	// Combine the base directory and the relative path.
	finalStack := append(baseStack, relStack...)
	result := strings.Join(finalStack, "/")
	// If the original current directory was absolute (i.e., started with '/'),
	// add a leading slash.
	if len(currDir) > 0 && currDir[0] == '/' {
		result = "/" + result
	}
	return result
}

func main() {
	// Define test cases.
	testCases := []struct {
		currDir  string
		relPath  string
		expected string
	}{
		{"a/b/c", "./d", "a/b/c/d"},
		{"a/b/c", "../d", "a/b/c/d"},
		{"a/b/c", "////d", "a/b/c/d"},
		{"a/b/c", "d", "a/b/c/d"},
		{"ab/c", "../e", "ab/c/e"}, // Modified case as requested.
		{"/a/b/c", "../d", "/a/b/c/d"},
		{"/a/b/c", "./d/e/../f", "/a/b/c/d/f"},
		{"a/b/c", "../../d", "a/b/c/d"},
		{"a/b/c", "../../../../d", "a/b/c/d"}, // Extra .. are ignored.
		{"a/b/c", ".././.././d", "a/b/c/d"},
		{"a/b/c", "/////.././d/e////f", "a/b/c/d/e/f"},
		{"a/b/c", "d/../e", "a/b/c/e"},
	}

	fmt.Println("Running test cases:")
	for i, tc := range testCases {
		result := normalizePath(tc.currDir, tc.relPath)
		fmt.Printf("Test %d: currentDir=%q, relativePath=%q, result=%q, expected=%q\n",
			i+1, tc.currDir, tc.relPath, result, tc.expected)
	}

	// Prompt user for input.
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter current directory: ")
	currDir, _ := reader.ReadString('\n')
	currDir = strings.TrimSpace(currDir)
	fmt.Print("Enter relative path: ")
	relPath, _ := reader.ReadString('\n')
	relPath = strings.TrimSpace(relPath)

	absPath := normalizePath(currDir, relPath)
	fmt.Printf("The absolute path is: %s\n", absPath)
}

// Package find will search for matched case-sensitive strings in files.
//
// Examples:
//	jay find . red
//		Find the word "red" in all go files in current folder and in subfolders.
//	jay find . red "*.*"
//		Find the word "red" in all files in current folder and in subfolders.
//	jay find . red "*.go" true false
//		Find word "red" in *.go files in current folder and in subfolders, but will exclude filenames.
package find

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	flagFind      *string
	flagFolder    *string
	flagExt       *string
	flagName      *bool
	flagRecursive *bool

	// MaxSize is the maximum size of a file Go will search through
	MaxSize int64 = 1048576
	// Folders to skip searching in
	SkipFolders = []string{"vendor", "node_modules", ".git"}
)

// Run starts the find filepath walk.
func Run(text, folder, ext *string, recursive, filename *bool) error {
	flagFind = text
	flagFolder = folder
	flagExt = ext
	flagRecursive = recursive
	flagName = filename

	fmt.Println()
	fmt.Println("Search Results")
	fmt.Println("==============")

	return filepath.Walk(*folder, visit)
}

// Visit analyzes a file to see if it matches the parameters.
// Original: https://gist.github.com/tdegrunt/045f6b3377f3f7ffa408
func visit(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// If path is a folder
	if fi.IsDir() {
		// Always search current folder
		if fi.Name() == "." {
			return nil
		}

		// Ignore specified folders
		if inArray(fi.Name(), SkipFolders) {
			return filepath.SkipDir
		}

		// If recursive is true
		if *flagRecursive {
			return nil
		}

		// Don't walk the folder
		return filepath.SkipDir
	}

	matched, err := filepath.Match(*flagExt, fi.Name())
	if err != nil {
		return err
	}

	// If the file extension matches
	if matched {
		// Skip file if too big
		if fi.Size() > MaxSize {
			fmt.Println("**ERROR: Skipping file too big", path)
			return nil
		}

		// Read the entire file into memory
		read, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("**ERROR: Could not read from", path)
			return nil
		}

		// Convert the bytes array into a string
		oldContents := string(read)

		// If the file name contains the search term, replace the file name
		if *flagName && strings.Contains(fi.Name(), *flagFind) {
			oldpath := path
			fmt.Println("Filename:", oldpath)
		}

		// If the file contains the search term
		if strings.Contains(oldContents, *flagFind) {
			count := strconv.Itoa(strings.Count(oldContents, *flagFind))
			fmt.Println("Contents:", path, "("+count+")")

		}
	}

	return nil
}

func inArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	out := os.Stdout

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Elapsed: ", elapsed)
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := showDirTreeRecursive(out, path, printFiles, 0, 0)

	if err != nil {
		panic(err.Error())
	}

	return nil
}

func showDirTreeRecursive(out io.Writer, path string, printFiles bool, level int, lastParentDirLevel int) error {
	directories, err := ioutil.ReadDir(path)
	filesCount := len(directories)
	isLastFile := false

	sort.Slice(directories, func(i, j int) bool {
		return directories[i].Name() < directories[j].Name()
	})

	if err != nil {
		fmt.Println(err)
	}

	if filesCount > 0 {
		level++
	} else if level > 1 {
		level--
	}

	for index, dir := range directories {
		if !printFiles && !dir.IsDir() {
			continue
		}

		if index == filesCount-1 || (!printFiles && index == filesCount-2 && !directories[index+1].IsDir()) {
			isLastFile = true
		}

		printLines(out, getDirNameLine(dir), level, isLastFile, lastParentDirLevel)

		if dir.IsDir() {
			if isLastFile && lastParentDirLevel < 1 {
				lastParentDirLevel = level
			}

			dirPath := path + string(os.PathSeparator) + dir.Name()
			err = showDirTreeRecursive(out, dirPath, printFiles, level, lastParentDirLevel)
		}
	}

	return nil
}

func getDirNameLine(dir os.FileInfo) string {
	len := dir.Name()

	if !dir.IsDir() {
		if dir.Size() > 0 {
			len += fmt.Sprintf(" (%vb)", dir.Size())
		} else {
			len += " (empty)"
		}
	}

	return len
}

func printLines(out io.Writer, dirName string, level int, isLastFile bool, lastParentDirLevel int) {
	line := ""
	delimiter := "├───"

	if isLastFile {
		delimiter = "└───"
	}

	if lastParentDirLevel > 0 {
		line = strings.Repeat("│\t", lastParentDirLevel-1) + strings.Repeat("\t", level-lastParentDirLevel) + delimiter
	} else {
		line = strings.Repeat("│\t", level-1) + delimiter
	}

	fmt.Fprint(out, line, dirName, "\n")
}

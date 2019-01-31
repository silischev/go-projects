package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
	err := dirTree(out, path, printFiles, -1, false)
	if err != nil {
		panic(err.Error())
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Elapsed: ", elapsed)
}

func dirTree(out *os.File, path string, printFiles bool, level int, isLastParentDir bool) error {
	directories, err := ioutil.ReadDir(path)
	filesCount := len(directories)
	isLastFile := false

	if err != nil {
		fmt.Println(err)
	}

	if filesCount > 0 {
		level++
	} else {
		level--
	}

	for index, dir := range directories {
		if index == filesCount-1 {
			isLastFile = true
		}

		if !printFiles && !dir.IsDir() {
			continue
		}

		printLines(dir.Name(), level, isLastFile, isLastParentDir)

		if dir.IsDir() {
			if index == filesCount-1 {
				isLastParentDir = true
			}

			dirPath := path + string(os.PathSeparator) + dir.Name()
			err = dirTree(out, dirPath, printFiles, level, isLastParentDir)
		}
	}

	return err
}

func printLines(dirName string, level int, isLastFile bool, isLastParentDir bool) {
	delimiter := "├───"

	if isLastFile {
		delimiter = "└───"
	}

	line := delimiter

	if level > 0 {
		if isLastFile {
			if isLastParentDir {
				line = "│\t" + strings.Repeat("\t", level-1) + delimiter
			} else {
				line = "│\t" + strings.Repeat("│\t", level-1) + delimiter
			}
		} else {
			line = strings.Repeat("│\t", level) + delimiter
		}
	}

	//fmt.Print(line, dirName, " -> ", isLastFile, " -> ", isLastParentDir, "\n")
	fmt.Print(line, dirName, "\n")
}

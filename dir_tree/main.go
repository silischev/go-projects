package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	//fmt.Println("*" + strings.Repeat("\t", 1), "test")
	//panic("end")

	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles, 0)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out *os.File, path string, printFiles bool, level int) error {
	directories, err := ioutil.ReadDir(path)
	filesCount := len(directories)
	currentDirNum := 0

	if err != nil {
		fmt.Println(err)
	}

	if printFiles {
		for _, dir := range directories {
			fmt.Println("├───"+strings.Repeat("	", level), dir.Name())
			//fmt.Println(len(directories), dir.Name())

			if dir.IsDir() {
				level += 1
				currentDirNum += 1
				err = dirTree(out, path+string(os.PathSeparator)+dir.Name(), true, level)
			} else if currentDirNum == filesCount {
				level -= 1
			}
		}
	}

	return err
}

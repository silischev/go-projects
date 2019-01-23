package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
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
}

func dirTree(out *os.File, path string, printFiles bool) error {
	directories, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Println(err)
	}

	if printFiles {
		for _, dir := range directories {
			//fmt.Println("└───", dir.Name())

			if dir.IsDir() {
				fmt.Println("└───", dir.Name())
				dirTree(out, path+string(os.PathSeparator)+dir.Name(), true)
			}
		}
	}

	return err
}

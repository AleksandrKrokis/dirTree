package main

import (
	"strconv"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sort"
	"bytes"
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

func dirTree(out io.Writer, path string, printFiles bool) error {
	path, ok := filepath.Abs(path)
	if (ok != nil) {
		panic("wrong path")
	}
	buffer := bytes.NewBufferString(drawBranch(path, path, printFiles))
	out.Write(buffer.Bytes())

	return nil
}

func drawBranch(branchPath string, basePath string, printFiles bool) (branch string) {
	dir, err := os.Open(branchPath)
	if err != nil {
		panic("can not open the dir")
	}
	dirInfo, err := dir.Readdir(0)
	if err != nil {
		panic("can't read the dir")
	}
	sort.SliceStable(dirInfo, func(i, j int) bool { return dirInfo[i].Name() < dirInfo[j].Name() })

	for i, item := range dirInfo {
		if item.IsDir() {
			branch = branch + generateLineString(filepath.Join(branchPath, item.Name()), basePath, i+1 == len(dirInfo))
			branch = branch + drawBranch(filepath.Join(branchPath, item.Name()), basePath, printFiles)
		} else if printFiles {
			branch = branch + generateLineString(filepath.Join(branchPath, item.Name()), basePath, i+1 == len(dirInfo))
		}
	}

	return
}

func generateLineString (path string, basePath string, last bool) (line string) {
	const levelOffset string = "│\t"
	connectionSymbol := "├───"
	if last {
		connectionSymbol = "└───"
	}
	lineLevel := len(strings.Split(path, "/")) - len(strings.Split(basePath, "/")) - 1

	item, err := os.Open(path)
	if err != nil {
		panic("can not open the dir")
	}
	itemInfo, err := item.Stat()
	if err != nil {
		panic("can not read stat")
	}

	for i := 0; i < lineLevel; i++ {
		line = line + levelOffset
	}

	if itemInfo.IsDir() {
		itemContent, err := item.Readdirnames(0)
		if err != nil {
			panic("can not open the dir")
		}

		if len(itemContent) < 1 {
			line = line + connectionSymbol + itemInfo.Name() + "\n"
		} else {
			line = line + connectionSymbol + itemInfo.Name() + "\n"
		}
	} else {
		size := "(empty)"
		if itemInfo.Size() > 0 {
			size = "(" + strconv.Itoa(int(itemInfo.Size())) + "b)"
		}
		line = line + connectionSymbol + itemInfo.Name() + " " + size + "\n"
	}
	return
}
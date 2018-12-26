package main

import (
	"strings"
	"strconv"
	"io"
	"os"
	"path/filepath"
	"sort"
	"bytes"
)

type fileSystemPoint struct {
	path string
	rootPath string
}

func (f *fileSystemPoint) Name() string {
	return filepath.Base(f.path)
}

func (f *fileSystemPoint) Size() int {
	dir, err := os.Open(f.path)
	if (err != nil) {
		panic("wrong path")
	}
	itemInfo, err := dir.Stat()
	if err != nil {
		panic("can not read stat")
	}
	if itemInfo.IsDir() {
		return int(0)
	}
	return int(itemInfo.Size())
}

func (f *fileSystemPoint) IsDir() bool {
	dir, err := os.Open(f.path)
	if (err != nil) {
		panic("wrong path")
	}
	itemInfo, err := dir.Stat()
	if err != nil {
		panic("can not read stat")
	}
	return itemInfo.IsDir()
}

func (f *fileSystemPoint) ReadDir(includeFiles bool) (content []fileSystemPoint) {
	if !f.IsDir() {
		return content
	}
	dir, err := os.Open(f.path)
	if (err != nil) {
		panic("wrong path")
	}
	dirInfo, err := dir.Readdir(0)
	if err != nil {
		panic("can't read the dir")
	}
	sort.SliceStable(dirInfo, func(i, j int) bool { return dirInfo[i].Name() < dirInfo[j].Name() });

	for _, item := range dirInfo {
		if item.IsDir() {
			content = append(content, fileSystemPoint{filepath.Join(f.path, item.Name()), f.rootPath})
		} else if includeFiles {
			content = append(content, fileSystemPoint{filepath.Join(f.path, item.Name()), f.rootPath})
		}
	}
	return
}

func (f *fileSystemPoint) Draw(includeFiles bool) (line string) {
	const simpleOffset = "│\t"
	const emptyOffset = "\t"
	const connectionOffset = "├───"
	const connectionOffsetLast = "└───"

	pathArr := strings.Split(strings.TrimLeft(f.path, filepath.Dir(f.rootPath)), string(os.PathSeparator))

	content := f.ReadDir(includeFiles)
	for i, item := range content {
		for index, pathPart := range pathArr {
			if pathPart == item.Name() {
				line += connectionOffset
			} else if index != 0 {
				line += simpleOffset
			}
		}
		if (i+1 == len(content)) {
			line += connectionOffsetLast + item.Name()
		} else {
			line += connectionOffset + item.Name()
		}
		if item.IsDir() {
			line += "\n" + item.Draw(includeFiles)
		} else {
			size := "(empty)"
			if item.Size() > int(0) {
				size = "(" + strconv.Itoa(item.Size()) + "b)"
			}
			line += " " + size + "\n"
		}
	}
	return
}

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
	point := fileSystemPoint{path, path}
	buffer := bytes.NewBufferString(point.Draw(printFiles))
	out.Write(buffer.Bytes())

	return nil
}
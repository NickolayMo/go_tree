package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func readDirNames(dirname string,  printFiles bool) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	b := names[:0]
	for _, name := range names {
		filename := filepath.Join(dirname, name)
		fileInfo, _ := os.Lstat(filename)
		if !fileInfo.IsDir() && !printFiles {
			continue
		}
		b = append(b, name)
	}
	sort.Strings(b)
	return b, nil
}

func getTree(out io.Writer, path string, printFiles bool, prefix string)  {
	names, _ := readDirNames(path, printFiles)
	for key, name := range names {
		filename := filepath.Join(path, name)
		s := "├───"
		fileInfo, _ := os.Lstat(filename)
		newPrefix := "";
		if fileInfo.IsDir(){
			newPrefix = prefix + "│\t";
		}
		if !fileInfo.IsDir(){
			size := fileInfo.Size();
			if size == 0 {
				name += " (empty)"
			}else {
				name += fmt.Sprintf(" (%vb)", size)
			}
		}
		if(key == len(names) - 1){
			s = "└───"
			newPrefix = prefix + "\t"
		}
		fmt.Fprintln(out, prefix + s + name)

		getTree(out, filename, printFiles, newPrefix)
		}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	getTree(out, path, printFiles, "")
	return nil
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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/knightso/enameg"
)

var (
	output = flag.String("output", "", "output file name; default srcdir/<filename>_ename.go")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tenameg [flags] [directory]\n")
	fmt.Fprintf(os.Stderr, "\tenameg [flags] files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		return
	}

	var files []string
	if len(args) == 1 && isDirectory(args[0]) {
		files = listFiles(args[0])
	} else {
		files = args
	}

	packageName, generated := enameg.Generate(files)
	if generated == "" {
		return
	}

	if *output == "" {
		*output = newOutputPath(files, packageName)
	}

	err := ioutil.WriteFile(*output, []byte(generated), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func newOutputPath(files []string, packageName string) string {
	if len(files) == 1 {
		baseName := files[0]
		components := strings.Split(baseName, ".")
		return strings.Join(components[0:len(components)-1], ".") + "_ename." + components[len(components)-1]
	}

	return filepath.Join(filepath.Dir(files[0]), packageName+"_ename.go")
}

func isDirectory(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func listFiles(dirname string) []string {
	fs, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	files := make([]string, 0, len(fs))
	for _, f := range fs {
		if !f.IsDir() {
			files = append(files, filepath.Join(dirname, f.Name()))
		}
	}

	return files
}

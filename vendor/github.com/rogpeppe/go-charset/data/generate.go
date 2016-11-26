// +build ignore

// go run generate.go && go fmt

// The generate-charset-data command generates the Go source code
// for code.google.com/p/go-charset/data from the data files
// found in code.google.com/p/go-charset/datafiles.
// It should be run in the go-charset root directory.
// The resulting Go files will need gofmt'ing.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

type info struct {
	Path string
}

var tfuncs = template.FuncMap{
	"basename": func(s string) string {
		return filepath.Base(s)
	},
	"read": func(path string) ([]byte, error) {
		return ioutil.ReadFile(path)
	},
}

var tmpl = template.Must(template.New("").Funcs(tfuncs).Parse(`
	// This file is automatically generated by generate-charset-data.
	// Do not hand-edit.

	package data
	import (
		"code.google.com/p/go-charset/charset"
		"io"
		"io/ioutil"
		"strings"
	)

	func init() {
		charset.RegisterDataFile({{basename .Path | printf "%q"}}, func() (io.ReadCloser, error) {
			r := strings.NewReader({{read .Path | printf "%q"}})
			return ioutil.NopCloser(r), nil
		})
	}
`))

var docTmpl = template.Must(template.New("").Funcs(tfuncs).Parse(`
	// This file is automatically generated by generate-charset-data.
	// Do not hand-edit.

	// The {{basename .Package}} package embeds all the charset
	// data files as Go data. It registers the data with the charset
	// package as a side effect of its import. To use:
	//
	//	import _ "code.google.com/p/go-charset"
	package {{basename .Package}}
`))

func main() {
	dataDir := filepath.Join("..", "datafiles")
	d, err := os.Open(dataDir)
	if err != nil {
		fatalf("%v", err)
	}
	names, err := d.Readdirnames(0)
	if err != nil {
		fatalf("cannot read datafiles dir: %v", err)
	}
	for _, name := range names {
		writeFile("data_"+name+".go", tmpl, info{
			Path: filepath.Join(dataDir, name),
		})
	}
}

func writeFile(name string, t *template.Template, data interface{}) {
	w, err := os.Create(name)
	if err != nil {
		fatalf("cannot create output file: %v", err)
	}
	defer w.Close()
	err = t.Execute(w, data)
	if err != nil {
		fatalf("template execute %q: %v", name, err)
	}
}

func fatalf(f string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s\n", fmt.Sprintf(f, a...))
	os.Exit(2)
}

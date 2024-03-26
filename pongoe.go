package pongoe

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/flosch/pongo2/v6"
)

type Templates struct {
	templates map[string]*pongo2.Template
	globals   pongo2.Context
}

func (views Templates) Add(name string, template *pongo2.Template) {

	if name == "" {
		return
	}

	if template == nil {
		return
	}

	views.templates[name] = template
}

func (views Templates) Render(name string, w io.Writer, c pongo2.Context) error {

	if _, ok := views.templates[name]; !ok {
		return fmt.Errorf("template %s not found", name)
	}

	return views.templates[name].ExecuteWriter(views.globals.Update(c), w)

}

func (views Templates) Dbg() {

	for k := range views.templates {
		fmt.Printf("template=%s\n", k)
	}
}

func LoadTemplatesFS(dir fs.FS) *Templates {
	templs := Templates{}
	templs.globals = pongo2.Context{}
	templs.templates = make(map[string]*pongo2.Template)

	ts := pongo2.NewSet("main", pongo2.NewFSLoader(dir))

	err := fs.WalkDir(dir, ".", func(path string, info fs.DirEntry, err error) error {

		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		np := strings.Replace(path, "./", "", 1)

		templ := pongo2.Must(ts.FromFile(path))

		templs.Add(np, templ)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return &templs
}

func LoadTemplates(dirpath string) *Templates {
	templs := Templates{}
	templs.globals = pongo2.Context{}

	templs.templates = make(map[string]*pongo2.Template)

	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		np := strings.Replace(path, dirpath+"/", "", 1)

		templ := pongo2.Must(pongo2.FromFile(path))

		templs.Add(np, templ)

		return nil
	})

	if err != nil {
		panic(err)
	}

	return &templs
}

func (views Templates) AddGlobal(key string, value any) {
	views.globals[key] = value
}


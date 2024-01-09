package pongoe

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/flosch/pongo2/v6"
)

type Templates map[string]*pongo2.Template

// for now we dont error
func (views Templates) Add(name string, template *pongo2.Template) {

	if name == "" {
		return
	}

	if template == nil {
		return
	}

	views[name] = template
}

func (views Templates) Render(name string, w io.Writer, c pongo2.Context) error {

	if _, ok := views[name]; !ok {
		return fmt.Errorf("template %s not found", name)
	}

	return views[name].ExecuteWriter(c, w)
}

func (views Templates) Dbg() {

	for k := range views {
		fmt.Printf("template=%s\n", k)
	}
}

func LoadTemplates(dirpath string) *Templates {
	templs := make(Templates)

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


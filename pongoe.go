package pongoe

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
)

type templates map[string]*pongo2.template

// for now we dont error
func (views templates) add(name string, template *pongo2.template) {

	if name == "" {
		return
	}

	if template == nil {
		return
	}

	views[name] = template
}

func (views templates) render(name string, w io.writer, c pongo2.context) error {

	if _, ok := views[name]; !ok {
		return fmt.errorf("template %s not found", name)
	}

	return views[name].executewriter(c, w)
}

func (views templates) dbg() {

	for k := range views {
		fmt.printf("template=%s\n", k)
	}
}

func loadtemplates(dirpath string) *templates {
	templs := make(templates)

	err := filepath.walk(dirpath, func(path string, info os.fileinfo, err error) error {

        fmt.printf("path=%s, isDir=%t", path, info.isdir())
		if err != nil {
			fmt.println(err)
			return err
		}

		if info.isdir() {
			return nil
		}

		parent := filepath.base(filepath.dir(path))

		if parent == dirpath || parent == "layout" {
			return nil
		}

		filename := filepath.base(path)

		name := filename[:len(filename)-len(filepath.ext(filename))]

		tname := parent + "_" + name

		templ := pongo2.must(pongo2.fromfile(path))

		templs.add(tname, templ)

		return nil
	})

	if err != nil {
		panic(err)
	}

	return &templs
}


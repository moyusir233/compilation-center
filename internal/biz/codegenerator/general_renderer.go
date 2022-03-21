package codegenerator

import (
	"bytes"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"strings"
	"text/template"
)

var (
	funcMap = template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"indent": func(count int, s string) string {
			return strings.Repeat(" ", count) + s
		},
		"camelCase": func(s string) string {
			return generator.CamelCase(s)
		},
	}
)

type generalRenderer struct {
	tmpl *template.Template
}

// renderMap key为需要渲染的模板名，value为渲染模板时传入的数据
type renderOption struct {
	tmplName string
	data     interface{}
}

type Device struct {
	DeviceClassID int
	Fields        []Field
}

type Field struct {
	Name, Type string
}

func newGeneralRenderer() *generalRenderer {
	return &generalRenderer{tmpl: template.New("generalRenderer").Funcs(funcMap)}
}

func (r *generalRenderer) render(option ...renderOption) (buffers []*bytes.Buffer, err error) {
	buffers = make([]*bytes.Buffer, 0, len(option))
	for _, o := range option {
		buffer := bytes.NewBuffer(make([]byte, 0, 1024))
		err := r.tmpl.ExecuteTemplate(buffer, o.tmplName, o.data)
		if err != nil {
			return nil, err
		}
		buffers = append(buffers, buffer)
	}
	return
}

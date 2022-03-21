package codegenerator

import (
	"bytes"
	"path/filepath"
)

type dataProcessingTmplRenderer struct {
	*generalRenderer
}

func newDataProcessingTmplRenderer(root string) (*dataProcessingTmplRenderer, error) {
	renderer := newGeneralRenderer()
	var err error

	// 实例化模板，解析指定目录下所有模板文件
	renderer.tmpl, err = renderer.tmpl.ParseGlob(filepath.Join(root, "*.template"))
	if err != nil {
		return nil, err
	}

	return &dataProcessingTmplRenderer{generalRenderer: renderer}, nil
}

func (r *dataProcessingTmplRenderer) renderConfigTmpl(configs []Device) (
	Code *bytes.Buffer, Proto *bytes.Buffer, err error) {

	options := []renderOption{
		{
			tmplName: "config.go.template",
			data:     configs,
		},
		{
			tmplName: "config.proto.template",
			data:     configs,
		},
	}

	buffers, err := r.render(options...)
	if err != nil {
		return nil, nil, err
	}

	return buffers[0], buffers[1], nil
}

func (r *dataProcessingTmplRenderer) renderWarningDetectTmpl(states []Device) (
	Code *bytes.Buffer, Proto *bytes.Buffer, err error) {
	options := []renderOption{
		{
			tmplName: "warning_detect.go.template",
			data:     states,
		},
		{
			tmplName: "warningDetect.proto.template",
			data:     states,
		},
	}

	buffers, err := r.render(options...)
	if err != nil {
		return nil, nil, err
	}

	return buffers[0], buffers[1], nil
}

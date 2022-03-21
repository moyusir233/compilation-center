package codegenerator

import (
	"bytes"
	"path/filepath"
)

type dataCollectionTmplRenderer struct {
	*generalRenderer
}

func newDataCollectionTmplRenderer(root string) (*dataCollectionTmplRenderer, error) {
	renderer := newGeneralRenderer()
	var err error

	// 实例化模板，解析指定目录下所有模板文件
	renderer.tmpl, err = renderer.tmpl.ParseGlob(filepath.Join(root, "*.template"))
	if err != nil {
		return nil, err
	}

	return &dataCollectionTmplRenderer{generalRenderer: renderer}, nil
}

// 渲染配置管理相关的go源码模板和protobuf服务定义模板
func (r *dataCollectionTmplRenderer) renderConfigTmpl(configs []Device) (
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

// 渲染故障预警相关的go源码模板和protobuf服务定义模板
func (r *dataCollectionTmplRenderer) renderWarningDetectTmpl(states []Device, warningDetectStates []Device) (
	Code *bytes.Buffer, Proto *bytes.Buffer, err error) {
	options := []renderOption{
		{
			tmplName: "warning_detect.go.template",
			data:     warningDetectStates,
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

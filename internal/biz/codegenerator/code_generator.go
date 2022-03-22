// Package codegenerator 负责生成数据收集服务和数据处理服务的服务端与客户端代码
package codegenerator

import (
	"bytes"
	v1 "gitee.com/moyusir/util/api/util/v1"
	"strings"
)

type CodeGenerator struct {
	dpRenderer *dataProcessingTmplRenderer
	dcRenderer *dataCollectionTmplRenderer
}

func NewCodeGenerator(dpTmplDir, dcTmplDir string) (*CodeGenerator, error) {
	processingTmplRenderer, err := newDataProcessingTmplRenderer(dpTmplDir)
	if err != nil {
		return nil, err
	}

	collectionTmplRenderer, err := newDataCollectionTmplRenderer(dcTmplDir)
	if err != nil {
		return nil, err
	}

	return &CodeGenerator{
		dpRenderer: processingTmplRenderer,
		dcRenderer: collectionTmplRenderer,
	}, nil
}

func (g *CodeGenerator) GetDataProcessingServiceFiles(
	configInfo []*v1.DeviceConfigRegisterInfo, stateInfo []*v1.DeviceStateRegisterInfo) (
	map[string]*bytes.Buffer, error) {
	files := make(map[string]*bytes.Buffer, 4)
	configs, states := transformFields(configInfo, stateInfo)

	// 产生数据处理服务相关的代码与服务定义文件
	dpConfigCode, dpConfigProto, err := g.dpRenderer.renderConfigTmpl(configs)
	if err != nil {
		return nil, err
	}
	files["config.go"] = dpConfigCode
	files["config.proto"] = dpConfigProto

	dpWarningCode, dpWarningProto, err := g.dpRenderer.renderWarningDetectTmpl(states)
	if err != nil {
		return nil, err
	}
	files["warning_detect.go"] = dpWarningCode
	files["warning_detect.proto"] = dpWarningProto

	return files, nil
}

func (g *CodeGenerator) GetDataCollectionServiceFiles(
	configInfo []*v1.DeviceConfigRegisterInfo, stateInfo []*v1.DeviceStateRegisterInfo) (
	map[string]*bytes.Buffer, error) {
	configs, states := transformFields(configInfo, stateInfo)
	files := make(map[string]*bytes.Buffer, 4)

	// 产生数据收集服务相关的代码与服务定义文件
	dcConfigCode, dcConfigProto, err := g.dcRenderer.renderConfigTmpl(configs)
	if err != nil {
		return nil, err
	}
	files["config.go"] = dcConfigCode
	files["config.proto"] = dcConfigProto

	dcWarningCode, dcWarningProto, err := g.dcRenderer.renderWarningDetectTmpl(states)
	if err != nil {
		return nil, err
	}
	files["warning_detect.go"] = dcWarningCode
	files["warning_detect.proto"] = dcWarningProto

	return files, nil
}

// 转换字段形式
func transformFields(configInfo []*v1.DeviceConfigRegisterInfo, stateInfo []*v1.DeviceStateRegisterInfo) (
	configs, states []Device) {
	configs = make([]Device, len(configInfo))
	states = make([]Device, len(stateInfo))

	// 处理配置注册信息
	for i, info := range configInfo {
		configs[i].DeviceClassID = i
		configs[i].Fields = make([]Field, len(info.Fields))

		for j, f := range info.Fields {
			configs[i].Fields[j].Name = f.Name
			// 时间戳字段需要转换声明的类型，不能直接用type的名称
			if f.Type == v1.Type_TIMESTAMP {
				configs[i].Fields[j].Type = "google.protobuf.Timestamp"
			} else {
				configs[i].Fields[j].Type = strings.ToLower(f.Type.String())
			}
		}
	}

	// 处理状态注册信息
	for i, info := range stateInfo {
		states[i].DeviceClassID = i
		states[i].Fields = make([]Field, len(info.Fields))

		for j, f := range info.Fields {
			states[i].Fields[j].Name = f.Name
			// 时间戳字段需要转换声明的类型，不能直接用type的名称
			if f.Type == v1.Type_TIMESTAMP {
				states[i].Fields[j].Type = "google.protobuf.Timestamp"
			} else {
				states[i].Fields[j].Type = strings.ToLower(f.Type.String())
			}
			// 预警规则非空的为预警字段
			states[i].Fields[j].Warning = f.WarningRule != nil
		}
	}

	return
}

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

// GetServiceFiles 生成数据收集服务和数据处理服务的代码
func (g *CodeGenerator) GetServiceFiles(
	configInfo []*v1.DeviceConfigRegisterInfo, stateInfo []*v1.DeviceStateRegisterInfo) (
	dc map[string]*bytes.Buffer, dp map[string]*bytes.Buffer, err error) {
	var (
		configs       = make([]Device, len(configInfo))
		states        = make([]Device, len(stateInfo))
		warningStates = make([]Device, len(stateInfo))
	)
	dc = make(map[string]*bytes.Buffer, 4)
	dp = make(map[string]*bytes.Buffer, 4)

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
		warningStates[i].DeviceClassID = i

		for j, f := range info.Fields {
			states[i].Fields[j].Name = f.Name
			// 时间戳字段需要转换声明的类型，不能直接用type的名称
			if f.Type == v1.Type_TIMESTAMP {
				states[i].Fields[j].Type = "google.protobuf.Timestamp"
			} else {
				states[i].Fields[j].Type = strings.ToLower(f.Type.String())
			}

			// 预警规则不为空即为预警字段
			if f.WarningRule != nil {
				warningStates[i].Fields = append(warningStates[i].Fields, states[i].Fields[j])
			}
		}
	}

	// 产生数据收集服务相关的代码与服务定义文件
	dcConfigCode, dcConfigProto, err := g.dcRenderer.renderConfigTmpl(configs)
	if err != nil {
		return nil, nil, err
	}
	dc["config.go"] = dcConfigCode
	dc["config.proto"] = dcConfigProto

	dcWarningCode, dcWarningProto, err := g.dcRenderer.renderWarningDetectTmpl(states, warningStates)
	if err != nil {
		return nil, nil, err
	}
	dc["warning_detect.go"] = dcWarningCode
	dc["warning_detect.proto"] = dcWarningProto

	// 产生数据处理服务相关的代码与服务定义文件
	dpConfigCode, dpConfigProto, err := g.dpRenderer.renderConfigTmpl(configs)
	if err != nil {
		return nil, nil, err
	}
	dp["config.go"] = dpConfigCode
	dp["config.proto"] = dpConfigProto

	dpWarningCode, dpWarningProto, err := g.dpRenderer.renderWarningDetectTmpl(states)
	if err != nil {
		return nil, nil, err
	}
	dp["warning_detect.go"] = dpWarningCode
	dp["warning_detect.proto"] = dpWarningProto

	return dc, dp, nil
}

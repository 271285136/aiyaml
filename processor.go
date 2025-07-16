package aiyaml

import (
	"context"
)

// Processor 主处理器，整合所有功能模块
type Processor struct {
	eventProcessor *EventProcessor
	yamlParser     *YAMLParser
	stringUtils    *StringUtils
	regexPatterns  *YAMLRegexPatterns
	logger         Logger
}

// NewProcessor 创建新的处理器
func NewProcessor(logger Logger) *Processor {
	return &Processor{
		eventProcessor: NewEventProcessor(logger),
		yamlParser:     NewYAMLParser(logger),
		stringUtils:    NewStringUtils(),
		regexPatterns:  NewYAMLRegexPatterns(),
		logger:         logger,
	}
}

// ProcessAIResponseEvents 处理AI响应事件流（保持向后兼容）
func (p *Processor) ProcessAIResponseEvents(ctx context.Context, eventChan chan SSEvent) (map[string]interface{}, error) {
	return p.eventProcessor.ProcessAIResponseEvents(ctx, eventChan)
}

// ProcessYAMLLines 直接处理YAML行（用于测试或独立使用）
func (p *Processor) ProcessYAMLLines(ctx context.Context, lines []string) (map[string]interface{}, error) {
	return p.yamlParser.LinesToMap(ctx, lines)
}

// GetStringUtils 获取字符串工具
func (p *Processor) GetStringUtils() *StringUtils {
	return p.stringUtils
}

// GetRegexPatterns 获取正则表达式模式
func (p *Processor) GetRegexPatterns() *YAMLRegexPatterns {
	return p.regexPatterns
}

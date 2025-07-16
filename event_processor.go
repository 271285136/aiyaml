package aiyaml

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// EventProcessor 事件处理器
type EventProcessor struct {
	logger Logger
}

// NewEventProcessor 创建新的事件处理器
func NewEventProcessor(logger Logger) *EventProcessor {
	return &EventProcessor{
		logger: logger,
	}
}

// ProcessAIResponseEvents 处理AI响应事件流
func (ep *EventProcessor) ProcessAIResponseEvents(ctx context.Context, eventChan chan SSEvent) (map[string]interface{}, error) {
	logEntry := ep.logger.WithContext(ctx).WithField("module", "yaml")
	result := []string{}
	allContent := ""
	line := ""

	// 正则表达式用于匹配YAML格式
	re := regexp.MustCompile(`^.*[a-zA-Z]+:.*`)
	re1 := regexp.MustCompile(`^.*[a-zA-Z]+: .+`)

	for event := range eventChan {
		// 如果上下文被取消，则退出
		if ctx.Err() != nil {
			logEntry.WithError(ctx.Err()).Error("context error")
			return nil, ctx.Err()
		}

		if event.Err != nil {
			logEntry.WithError(event.Err).Error("event error")
			return nil, fmt.Errorf("event error: %v", event.Err)
		}

		var rawData map[string]interface{}
		if err := json.Unmarshal([]byte(event.Data), &rawData); err != nil {
			logEntry.WithError(err).Error("unmarshal error")
			return nil, fmt.Errorf("unmarshal error: %v", err)
		}

		if choices, ok := rawData["Choices"].([]interface{}); ok && len(choices) > 0 {
			if delta, ok := choices[0].(map[string]interface{})["Delta"].(map[string]interface{}); ok {
				if contentValue, exists := delta["Content"]; exists {
					content := contentValue.(string)
					line += content
					allContent += content

					if strings.HasSuffix(line, "\n") || strings.HasSuffix(line, "\\n") {
						logEntry.Infof("line: %s\n", line)
						preLine := ""
						if len(result) > 0 {
							preLine = result[len(result)-1]
							logEntry.Infof("preLine: %s", preLine)
						}

						// 处理行合并逻辑
						if !re.MatchString(line) && len(result) > 0 && (!strings.HasPrefix(strings.TrimSpace(line), "- ") || re1.MatchString(preLine)) {
							line = preLine + line
							result[len(result)-1] = line
							line = ""
							continue
						}

						tabCount := strings.Count(line, "\t")
						logEntry.Infof("tabCount: %d", tabCount)

						// 清理YAML标记
						line = strings.TrimPrefix(line, "```yaml")
						line = strings.TrimSuffix(line, "```")

						if strings.TrimSpace(line) == "" {
							continue
						}

						result = append(result, line)
						line = ""
					}
				}
			}
		}
	}

	// 处理最后一行
	tabCount := strings.Count(line, "\t")
	logEntry.Infof("tabCount: %d", tabCount)
	logEntry.Infof("最后一行line: %s", line)
	logEntry.Infof("allContent: %s", allContent)

	line = strings.TrimPrefix(line, "```yaml\n")
	line = strings.TrimSuffix(line, "\n```")
	line = strings.TrimSuffix(line, "```")

	if strings.TrimSpace(line) != "" {
		result = append(result, line)
	}

	// 将YAML行转换为map
	yamlParser := NewYAMLParser(ep.logger)
	yamlMap, err := yamlParser.LinesToMap(ctx, result)
	if err != nil {
		logEntry.WithError(err).Error("yamlLinesToMap error")
		return nil, fmt.Errorf("yamlLinesToMap error: %v", err)
	}

	jsonData, err := json.Marshal(yamlMap)
	if err != nil {
		logEntry.WithError(err).Error("json marshal error")
		return nil, fmt.Errorf("json marshal error: %v", err)
	}

	logEntry.Infof("yamlMap: %v", string(jsonData))
	return yamlMap, nil
}

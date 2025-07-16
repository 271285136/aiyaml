package aiyaml

import (
	"context"
	"strings"
)

// 为了保持向后兼容，保留原始函数名
// ProcessAIResponseEvents 处理AI响应事件流
func ProcessAIResponseEvents(ctx context.Context, eventChan chan SSEvent) (map[string]interface{}, error) {
	processor := NewProcessor(NewDefaultLogger().WithContext(ctx))
	var result []string
	line := ""
	allContent := ""
	for event := range eventChan {
		if ctx.Err() != nil {
			processor.logger.Error("context error", ctx.Err())
			return nil, ctx.Err()
		}
		if event.Err != nil {
			processor.logger.Error("event error", event.Err)
			return nil, event.Err
		}
		line += string(event.Data)
		allContent += string(event.Data)
		if strings.HasSuffix(line, "\n") || strings.HasSuffix(line, "\\n") {
			processor.logger.Infof("line: %s\n", line)
			preLine := ""
			if len(result) > 0 {
				preLine = result[len(result)-1]
				processor.logger.Infof("preLine: %s", preLine)
			}
			if !processor.regexPatterns.KeyValuePattern.MatchString(line) && len(result) > 0 && (!strings.HasPrefix(strings.TrimSpace(line), "- ") || processor.regexPatterns.KeyValueWithContent.MatchString(preLine)) {
				line = preLine + line
				line = strings.TrimRight(line, "\n")
				line = strings.TrimRight(line, "\\n")
				result[len(result)-1] = line
				line = ""
				continue
			}
			tabCount := strings.Count(line, "\t")
			processor.logger.Infof("tabCount: %d", tabCount)
			line = processor.stringUtils.CleanYAMLMarkers(line)
			if strings.TrimSpace(line) == "" {
				continue
			}
			line = strings.TrimRight(line, "\n")
			line = strings.TrimRight(line, "\\n")
			result = append(result, line)
			line = ""
		}
	}
	line = processor.stringUtils.CleanYAMLMarkers(line)
	if strings.TrimSpace(line) != "" {
		result = append(result, line)
	}
	processor.logger.Infof("allContent: \n%s", allContent)
	processor.logger.Infof("result: %v", result)
	return processor.ProcessYAMLLines(ctx, result)
}

// YamlLinesToMap 将yaml代码行转换为map（保持向后兼容）
func YamlLinesToMap(ctx context.Context, lines []string) (map[string]interface{}, error) {
	processor := NewProcessor(NewDefaultLogger().WithContext(ctx))
	return processor.ProcessYAMLLines(ctx, lines)
}

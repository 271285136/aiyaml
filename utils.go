package aiyaml

import (
	"regexp"
	"strings"
)

// YAMLRegexPatterns YAML正则表达式模式
type YAMLRegexPatterns struct {
	KeyValuePattern     *regexp.Regexp
	KeyValueWithContent *regexp.Regexp
}

// NewYAMLRegexPatterns 创建YAML正则表达式模式
func NewYAMLRegexPatterns() *YAMLRegexPatterns {
	return &YAMLRegexPatterns{
		KeyValuePattern:     regexp.MustCompile(`^.*[a-zA-Z]+:.*`),
		KeyValueWithContent: regexp.MustCompile(`^.*[a-zA-Z]+: .+`),
	}
}

// StringUtils 字符串工具
type StringUtils struct{}

// NewStringUtils 创建字符串工具
func NewStringUtils() *StringUtils {
	return &StringUtils{}
}

// CleanYAMLMarkers 清理YAML标记
func (su *StringUtils) CleanYAMLMarkers(line string) string {
	line = strings.TrimPrefix(line, "```yaml")
	line = strings.TrimPrefix(line, "```yaml\n")
	line = strings.TrimSuffix(line, "```")
	line = strings.TrimSuffix(line, "\n```")
	// return strings.TrimSpace(line)
	return line
}

// CalculateIndent 计算缩进级别
func (su *StringUtils) CalculateIndent(line string) int {
	indent := strings.Count(line, "\t")
	if indent == 0 {
		indent = strings.Count(line, "  ")
	}
	return indent
}

// IsArrayItem 判断是否为数组项
func (su *StringUtils) IsArrayItem(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "- ")
}

// ParseKeyValue 解析键值对
func (su *StringUtils) ParseKeyValue(line string) (string, string, bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
}

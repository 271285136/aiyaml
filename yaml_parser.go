package aiyaml

import (
	"context"
	"regexp"
	"strings"
)

// YAMLParser YAML解析器
type YAMLParser struct {
	logger Logger
}

// NewYAMLParser 创建新的YAML解析器
func NewYAMLParser(logger Logger) *YAMLParser {
	return &YAMLParser{
		logger: logger,
	}
}

// LinesToMap 将yaml代码行转换为map
func (yp *YAMLParser) LinesToMap(ctx context.Context, lines []string) (map[string]interface{}, error) {
	re := regexp.MustCompile(`^.*[a-zA-Z]+: .*`)

	type node struct {
		value  interface{} // map[string]interface{} 或 []interface{}
		key    string      // 当前key
		indent int
	}
	result := make(map[string]interface{})
	stack := []node{{value: result, key: "", indent: -1}}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		indent := strings.Count(line, "\t")
		if indent == 0 {
			indent = strings.Count(line, "  ")
		}
		trimmedLine := strings.TrimSpace(line)
		for len(stack) > 0 && indent <= stack[len(stack)-1].indent {
			stack = stack[:len(stack)-1]
		}
		parent := stack[len(stack)-1]
		if strings.HasPrefix(trimmedLine, "- ") && len(stack) > 1 {
			// 直接操作父级map的key
			if p, ok := stack[len(stack)-2].value.(map[string]interface{}); ok {
				arr, _ := p[parent.key].([]interface{})
				itemStr := strings.TrimPrefix(trimmedLine, "- ")
				var newItem interface{}
				if re.MatchString(itemStr) {
					parts := strings.SplitN(itemStr, ":", 2)
					k := strings.TrimSpace(parts[0])
					v := strings.TrimSpace(parts[1])
					item := map[string]interface{}{k: v}
					newItem = item
					if i+1 < len(lines) {
						nextLine := lines[i+1]
						nextIndent := strings.Count(nextLine, "\t")
						if nextIndent == 0 {
							nextIndent = strings.Count(nextLine, "  ")
						}
						if nextIndent > indent {
							stack = append(stack, node{value: item, key: "", indent: indent})
						}
					}
				} else {
					newItem = itemStr
				}
				arr = append(arr, newItem)
				p[parent.key] = arr
			}
			continue
		}
		// 处理 key: value 或 key:
		parts := strings.SplitN(trimmedLine, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if value == "" {
			// 判断下一级是数组还是map
			isArray := false
			if i+1 < len(lines) {
				nextLine := lines[i+1]
				nextIndent := strings.Count(nextLine, "\t")
				if nextIndent == 0 {
					nextIndent = strings.Count(nextLine, "  ")
				}
				nextTrimmed := strings.TrimSpace(nextLine)
				if nextTrimmed != "" && strings.HasPrefix(nextTrimmed, "- ") && nextIndent > indent {
					isArray = true
				}
			}
			if isArray {
				newArr := []interface{}{}
				parent.value.(map[string]interface{})[key] = newArr
				stack = append(stack, node{value: newArr, key: key, indent: indent})
			} else {
				newMap := map[string]interface{}{}
				parent.value.(map[string]interface{})[key] = newMap
				stack = append(stack, node{value: newMap, key: key, indent: indent})
			}
		} else {
			parent.value.(map[string]interface{})[key] = value
		}
	}
	return result, nil
}

package aiyaml

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestYAMLParser(t *testing.T) {
	// 创建处理器
	processor := NewProcessor(NewDefaultLogger())

	// 测试YAML行
	lines := []string{
		"name: test",
		"version: 1.0",
		"settings:",
		"  debug: true",
		"  timeout: 30",
		"features:",
		"  - feature1",
		"  - feature2",
		"  - name: feature3",
		"    enabled: true",
	}

	ctx := context.Background()
	result, err := processor.ProcessYAMLLines(ctx, lines)

	if err != nil {
		t.Fatalf("处理YAML行失败: %v", err)
	}

	// 验证结果
	if result["name"] != "test" {
		t.Errorf("期望 name=test, 得到 %v", result["name"])
	}

	if result["version"] != "1.0" {
		t.Errorf("期望 version=1.0, 得到 %v", result["version"])
	}

	settings, ok := result["settings"].(map[string]interface{})
	if !ok {
		t.Fatal("settings 不是 map")
	}

	if settings["debug"] != "true" {
		t.Errorf("期望 settings.debug=true, 得到 %v", settings["debug"])
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json marshal error: %v", err)
	}
	fmt.Printf("解析结果: %+v\n", string(jsonData))
}

func TestYAMLParserComplex(t *testing.T) {
	processor := NewProcessor(NewDefaultLogger())

	// 测试复杂的嵌套结构
	lines := []string{
		"api:",
		"  version: v1",
		"  endpoints:",
		"    - name: users",
		"      path: /api/users",
		"      methods:",
		"        - GET",
		"        - POST",
		"    - name: posts",
		"      path: /api/posts",
		"      methods:",
		"        - GET",
		"        - POST",
		"        - PUT",
		"        - DELETE",
		"database:",
		"  host: localhost",
		"  port: 5432",
		"  credentials:",
		"    username: admin",
		"    password: secret",
	}

	ctx := context.Background()
	result, err := processor.ProcessYAMLLines(ctx, lines)

	if err != nil {
		t.Fatalf("处理复杂YAML失败: %v", err)
	}

	// 验证嵌套结构
	api, ok := result["api"].(map[string]interface{})
	if !ok {
		t.Fatal("api 不是 map")
	}

	if api["version"] != "v1" {
		t.Errorf("期望 api.version=v1, 得到 %v", api["version"])
	}

	endpoints, ok := api["endpoints"].([]interface{})
	if !ok {
		t.Fatal("endpoints 不是数组")
	}

	if len(endpoints) != 2 {
		t.Errorf("期望 2 个端点, 得到 %d", len(endpoints))
	}

	// 验证第一个端点
	firstEndpoint := endpoints[0].(map[string]interface{})
	if firstEndpoint["name"] != "users" {
		t.Errorf("期望第一个端点名称是 users, 得到 %v", firstEndpoint["name"])
	}

	methods := firstEndpoint["methods"].([]interface{})
	if len(methods) != 2 {
		t.Errorf("期望 2 个方法, 得到 %d", len(methods))
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json marshal error: %v", err)
	}
	fmt.Printf("解析结果: %+v\n", string(jsonData))
}

func TestYAMLParserEdgeCases(t *testing.T) {
	processor := NewProcessor(NewDefaultLogger())

	testCases := []struct {
		name     string
		lines    []string
		expected map[string]interface{}
	}{
		{
			name: "空行和注释",
			lines: []string{
				"",
				"# 这是注释",
				"name: test",
				"  # 行内注释",
			},
			expected: map[string]interface{}{
				"name": "test",
			},
		},
		{
			name: "空值",
			lines: []string{
				"empty_value:",
				"  nested:",
				"    deep:",
			},
			expected: map[string]interface{}{
				"empty_value": map[string]interface{}{
					"nested": map[string]interface{}{
						"deep": map[string]interface{}{},
					},
				},
			},
		},
		{
			name: "数组中的对象",
			lines: []string{
				"items:",
				"  - id: 1",
				"    name: item1",
				"  - id: 2",
				"    name: item2",
			},
			expected: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"id": "1", "name": "item1"},
					map[string]interface{}{"id": "2", "name": "item2"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := processor.ProcessYAMLLines(ctx, tc.lines)

			if err != nil {
				t.Fatalf("处理失败: %v", err)
			}

			// 简单验证结果不为空
			if len(result) == 0 {
				t.Error("结果为空")
			}
		})
	}
}

func TestStringUtils(t *testing.T) {
	utils := NewStringUtils()

	// 测试清理YAML标记
	testCases := []struct {
		input    string
		expected string
	}{
		{"```yaml\nname: test\n```", "name: test"},
		{"```yamlname: test```", "name: test"},
		{"name: test", "name: test"},
		{"```yaml\n```", ""},
	}

	for _, tc := range testCases {
		result := utils.CleanYAMLMarkers(tc.input)
		if result != tc.expected {
			t.Errorf("输入 '%s', 期望 '%s', 得到 '%s'", tc.input, tc.expected, result)
		}
	}

	// 测试计算缩进
	indentTests := []struct {
		input    string
		expected int
	}{
		{"\t\tname: test", 2},
		{"  name: test", 1},
		{"name: test", 0},
		{"\t  \tname: test", 2},
	}

	for _, tc := range indentTests {
		result := utils.CalculateIndent(tc.input)
		if result != tc.expected {
			t.Errorf("输入 '%s', 期望缩进 %d, 得到 %d", tc.input, tc.expected, result)
		}
	}

	// 测试解析键值对
	keyValueTests := []struct {
		input    string
		key      string
		value    string
		expected bool
	}{
		{"name: test", "name", "test", true},
		{"version: 1.0", "version", "1.0", true},
		{"empty:", "empty", "", true},
		{"invalid", "", "", false},
		{"", "", "", false},
	}

	for _, tc := range keyValueTests {
		key, value, ok := utils.ParseKeyValue(tc.input)
		if ok != tc.expected {
			t.Errorf("输入 '%s', 期望解析成功 %v, 得到 %v", tc.input, tc.expected, ok)
		}
		if ok && (key != tc.key || value != tc.value) {
			t.Errorf("输入 '%s', 期望 key='%s', value='%s', 得到 key='%s', value='%s'",
				tc.input, tc.key, tc.value, key, value)
		}
	}
}

func TestStringUtilsIsArrayItem(t *testing.T) {
	utils := NewStringUtils()

	testCases := []struct {
		input    string
		expected bool
	}{
		{"- item", true},
		{"  - item", true},
		{"\t- item", true},
		{"- name: value", true},
		{"item", false},
		{"  item", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := utils.IsArrayItem(tc.input)
		if result != tc.expected {
			t.Errorf("输入 '%s', 期望 %v, 得到 %v", tc.input, tc.expected, result)
		}
	}
}

func TestRegexPatterns(t *testing.T) {
	patterns := NewYAMLRegexPatterns()

	// 测试键值对模式
	testCases := []struct {
		input    string
		expected bool
	}{
		{"name: test", true},
		{"version: 1.0", true},
		{"enabled: true", true},
		{"name:", false},
		{"invalid", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := patterns.KeyValuePattern.MatchString(tc.input)
		if result != tc.expected {
			t.Errorf("KeyValuePattern 输入 '%s', 期望 %v, 得到 %v", tc.input, tc.expected, result)
		}
	}

	// 测试带内容的键值对模式
	contentTests := []struct {
		input    string
		expected bool
	}{
		{"name: test", true},
		{"version: 1.0", true},
		{"name:", false},
		{"invalid", false},
		{"", false},
	}

	for _, tc := range contentTests {
		result := patterns.KeyValueWithContent.MatchString(tc.input)
		if result != tc.expected {
			t.Errorf("KeyValueWithContent 输入 '%s', 期望 %v, 得到 %v", tc.input, tc.expected, result)
		}
	}
}

func TestEventProcessor(t *testing.T) {
	processor := NewProcessor(NewDefaultLogger())

	// 测试事件处理器创建
	if processor.eventProcessor == nil {
		t.Error("事件处理器未创建")
	}

	if processor.yamlParser == nil {
		t.Error("YAML解析器未创建")
	}

	if processor.stringUtils == nil {
		t.Error("字符串工具未创建")
	}

	if processor.regexPatterns == nil {
		t.Error("正则表达式模式未创建")
	}
}

func TestDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger()

	// 测试日志器方法不会panic
	ctx := context.Background()

	// 测试WithContext
	loggerWithCtx := logger.WithContext(ctx)
	if loggerWithCtx == nil {
		t.Error("WithContext 返回 nil")
	}

	// 测试WithField
	loggerWithField := logger.WithField("key", "value")
	if loggerWithField == nil {
		t.Error("WithField 返回 nil")
	}

	// 测试WithError
	loggerWithError := logger.WithError(fmt.Errorf("test error"))
	if loggerWithError == nil {
		t.Error("WithError 返回 nil")
	}

	// 测试日志方法（这些应该不会panic）
	logger.Info("test info")
	logger.Infof("test info %s", "format")
	logger.Error("test error")
	logger.Errorf("test error %s", "format")
}

func TestProcessorMethods(t *testing.T) {
	processor := NewProcessor(NewDefaultLogger())

	// 测试获取工具方法
	stringUtils := processor.GetStringUtils()
	if stringUtils == nil {
		t.Error("GetStringUtils 返回 nil")
	}

	regexPatterns := processor.GetRegexPatterns()
	if regexPatterns == nil {
		t.Error("GetRegexPatterns 返回 nil")
	}

	// 测试工具方法功能
	indent := stringUtils.CalculateIndent("\t\tname: test")
	if indent != 2 {
		t.Errorf("期望缩进 2, 得到 %d", indent)
	}

	if !regexPatterns.KeyValuePattern.MatchString("name: test") {
		t.Error("正则表达式模式不匹配")
	}
}

func TestBackwardCompatibility(t *testing.T) {
	// 测试向后兼容的函数
	ctx := context.Background()

	// 测试 YamlLinesToMap
	lines := []string{
		"name: test",
		"version: 1.0",
	}

	result, err := YamlLinesToMap(ctx, lines)
	if err != nil {
		t.Fatalf("YamlLinesToMap 失败: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("期望 name=test, 得到 %v", result["name"])
	}

	// 测试 ProcessAIResponseEvents（模拟事件）
	eventChan := make(chan SSEvent, 1)
	go func() {
		// 将内容分词拆分发送
		for _, ch := range []byte("name: test-\ntea") {
			eventChan <- SSEvent{
				Data: []byte{ch},
				Err:  nil,
			}
		}
		eventChan <- SSEvent{
			Data: []byte(`\n`),
			Err:  nil,
		}
		for _, ch := range []byte("version: 1.0") {
			eventChan <- SSEvent{
				Data: []byte{ch},
				Err:  nil,
			}
		}
		eventChan <- SSEvent{
			Data: []byte(`\n`),
			Err:  nil,
		}
		close(eventChan)
	}()

	result, err = ProcessAIResponseEvents(ctx, eventChan)
	if err != nil {
		t.Fatalf("ProcessAIResponseEvents 失败: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("期望 name=test, 得到 %v", result["name"])
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json marshal error: %v", err)
	}
	fmt.Printf("解析结果: %+v\n", string(jsonData))
}

func TestBackwardCompatibilityComplex(t *testing.T) {
	ctx := context.Background()

	// 测试复杂的嵌套结构
	lines := []string{
		"api:",
		"  version: v1",
		"  endpoints:",
		"    - name: users",
		"      path: /api/users",
		"      methods:",
		"        - GET",
		"        - POST",
		"    - name: posts",
		"      path: /api/posts",
		"      methods:",
		"        - GET",
		"        - POST",
		"        - PUT",
		"        - DELETE",
	}

	result, err := YamlLinesToMap(ctx, lines)
	if err != nil {
		t.Fatalf("YamlLinesToMap 复杂结构失败: %v", err)
	}

	// 验证嵌套结构
	api, ok := result["api"].(map[string]interface{})
	if !ok {
		t.Fatal("api 不是 map")
	}

	if api["version"] != "v1" {
		t.Errorf("期望 api.version=v1, 得到 %v", api["version"])
	}

	endpoints, ok := api["endpoints"].([]interface{})
	if !ok {
		t.Fatal("endpoints 不是数组")
	}

	if len(endpoints) != 2 {
		t.Errorf("期望 2 个端点, 得到 %d", len(endpoints))
	}

	// 验证第一个端点
	firstEndpoint := endpoints[0].(map[string]interface{})
	if firstEndpoint["name"] != "users" {
		t.Errorf("期望第一个端点名称是 users, 得到 %v", firstEndpoint["name"])
	}

	methods := firstEndpoint["methods"].([]interface{})
	if len(methods) != 2 {
		t.Errorf("期望 2 个方法, 得到 %d", len(methods))
	}
}

func TestBackwardCompatibilityWithEvents(t *testing.T) {
	ctx := context.Background()

	// 测试通过事件流处理复杂YAML
	eventChan := make(chan SSEvent, 1)
	go func() {
		complexYAML := `api:\n
  version: v1\n
  endpoints:\n
    - name: users\n
      path: /api/users\n
      methods:\n
        - GET\n
        - POST\n
    - name: posts\n
      path: /api/posts\n
      methods:\n
        - GET\n
        - POST\n
        - PUT\n
        - DELETE\n`

		// 逐字符发送
		for _, ch := range []byte(complexYAML) {
			eventChan <- SSEvent{
				Data: []byte{ch},
				Err:  nil,
			}
		}
		close(eventChan)
	}()

	result, err := ProcessAIResponseEvents(ctx, eventChan)
	if err != nil {
		t.Fatalf("ProcessAIResponseEvents 复杂结构失败: %v", err)
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json marshal error: %v", err)
	}
	fmt.Printf("解析结果: %+v\n", string(jsonData))

	// 验证结果
	api, ok := result["api"].(map[string]interface{})
	if !ok {
		t.Fatal("api 不是 map")
	}

	if api["version"] != "v1" {
		t.Errorf("期望 api.version=v1, 得到 %v", api["version"])
	}

	endpoints, ok := api["endpoints"].([]interface{})
	if !ok {
		t.Fatal("endpoints 不是数组")
	}

	if len(endpoints) != 2 {
		t.Errorf("期望 2 个端点, 得到 %d", len(endpoints))
	}
}

func TestBackwardCompatibilityEdgeCases(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name     string
		lines    []string
		expected map[string]interface{}
	}{
		{
			name: "空值和嵌套",
			lines: []string{
				"empty_value:",
				"  nested:",
				"    deep:",
			},
			expected: map[string]interface{}{
				"empty_value": map[string]interface{}{
					"nested": map[string]interface{}{
						"deep": map[string]interface{}{},
					},
				},
			},
		},
		{
			name: "数组中的对象",
			lines: []string{
				"items:",
				"  - id: 1",
				"    name: item1",
				"  - id: 2",
				"    name: item2",
			},
			expected: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"id": "1", "name": "item1"},
					map[string]interface{}{"id": "2", "name": "item2"},
				},
			},
		},
		{
			name: "混合类型",
			lines: []string{
				"string_value: hello",
				"number_value: 42",
				"boolean_value: true",
				"null_value:",
				"array_value:",
				"  - item1",
				"  - item2",
			},
			expected: map[string]interface{}{
				"string_value":  "hello",
				"number_value":  "42",
				"boolean_value": "true",
				"null_value":    map[string]interface{}{},
				"array_value":   []interface{}{"item1", "item2"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := YamlLinesToMap(ctx, tc.lines)
			if err != nil {
				t.Fatalf("YamlLinesToMap 失败: %v", err)
			}

			// 验证结果不为空
			if len(result) == 0 {
				t.Error("结果为空")
			}

			// 验证基本结构
			for key := range tc.expected {
				if result[key] == nil {
					t.Errorf("期望 key '%s' 存在，但为 nil", key)
				}
			}
		})
	}
}

func TestBackwardCompatibilityWithYAMLMarkers(t *testing.T) {
	ctx := context.Background()

	// 测试包含YAML标记的内容
	lines := []string{
		"```yaml",
		"name: test",
		"version: 1.0",
		"settings:",
		"  debug: true",
		"  timeout: 30",
		"```",
	}

	result, err := YamlLinesToMap(ctx, lines)
	if err != nil {
		t.Fatalf("YamlLinesToMap 带标记失败: %v", err)
	}

	// 验证YAML标记被正确清理
	if result["name"] != "test" {
		t.Errorf("期望 name=test, 得到 %v", result["name"])
	}

	if result["version"] != "1.0" {
		t.Errorf("期望 version=1.0, 得到 %v", result["version"])
	}

	settings, ok := result["settings"].(map[string]interface{})
	if !ok {
		t.Fatal("settings 不是 map")
	}

	if settings["debug"] != "true" {
		t.Errorf("期望 settings.debug=true, 得到 %v", settings["debug"])
	}
}

func TestBackwardCompatibilityWithLargeData(t *testing.T) {
	ctx := context.Background()

	// 测试大数据量
	largeLines := make([]string, 0, 1000)
	largeLines = append(largeLines, "root:")
	for i := 0; i < 100; i++ {
		largeLines = append(largeLines, fmt.Sprintf("  item_%d:", i))
		largeLines = append(largeLines, fmt.Sprintf("    id: %d", i))
		largeLines = append(largeLines, fmt.Sprintf("    name: item_%d", i))
		largeLines = append(largeLines, "    tags:")
		for j := 0; j < 5; j++ {
			largeLines = append(largeLines, fmt.Sprintf("      - tag_%d_%d", i, j))
		}
	}

	result, err := YamlLinesToMap(ctx, largeLines)
	if err != nil {
		t.Fatalf("YamlLinesToMap 大数据量失败: %v", err)
	}

	// 验证基本结构
	root, ok := result["root"].(map[string]interface{})
	if !ok {
		t.Fatal("root 不是 map")
	}

	// 验证有几个item
	itemCount := 0
	for key := range root {
		if strings.HasPrefix(key, "item_") {
			itemCount++
		}
	}

	if itemCount != 100 {
		t.Errorf("期望 100 个 item, 得到 %d", itemCount)
	}
}

func TestBackwardCompatibilityWithSpecialCharacters(t *testing.T) {
	ctx := context.Background()

	// 测试特殊字符
	lines := []string{
		"special_chars:",
		"  unicode: 中文测试",
		"  symbols: !@#$%^&*()",
		"  quotes: \"double quotes\"",
		"  single_quotes: 'single quotes'",
		"  newlines: \"line1\\nline2\"",
		"  tabs: \"col1\\tcol2\"",
	}

	result, err := YamlLinesToMap(ctx, lines)
	if err != nil {
		t.Fatalf("YamlLinesToMap 特殊字符失败: %v", err)
	}

	specialChars, ok := result["special_chars"].(map[string]interface{})
	if !ok {
		t.Fatal("special_chars 不是 map")
	}

	// 验证特殊字符被正确处理
	if specialChars["unicode"] != "中文测试" {
		t.Errorf("期望 unicode=中文测试, 得到 %v", specialChars["unicode"])
	}

	if specialChars["symbols"] != "!@#$%^&*()" {
		t.Errorf("期望 symbols=!@#$%%^&*(), 得到 %v", specialChars["symbols"])
	}
}

func TestBackwardCompatibilityWithNestedArrays(t *testing.T) {
	ctx := context.Background()

	// 测试嵌套数组（当前解析器支持的结构）
	lines := []string{
		"matrix:",
		"  - row: 1",
		"    values:",
		"      - 1",
		"      - 2",
		"      - 3",
		"  - row: 2",
		"    values:",
		"      - 4",
		"      - 5",
		"      - 6",
		"  - row: 3",
		"    values:",
		"      - 7",
		"      - 8",
		"      - 9",
	}

	result, err := YamlLinesToMap(ctx, lines)
	if err != nil {
		t.Fatalf("YamlLinesToMap 嵌套数组失败: %v", err)
	}

	matrix, ok := result["matrix"].([]interface{})
	if !ok {
		t.Fatal("matrix 不是数组")
	}

	if len(matrix) != 3 {
		t.Errorf("期望 3 行, 得到 %d", len(matrix))
	}

	// 验证第一行
	firstRow, ok := matrix[0].(map[string]interface{})
	if !ok {
		t.Fatal("第一行不是map")
	}

	if firstRow["row"] != "1" {
		t.Errorf("期望第一行的row是 1, 得到 %v", firstRow["row"])
	}

	values, ok := firstRow["values"].([]interface{})
	if !ok {
		t.Fatal("values 不是数组")
	}

	if len(values) != 3 {
		t.Errorf("期望第一行有 3 个值, 得到 %d", len(values))
	}

	if values[0] != "1" {
		t.Errorf("期望第一个值是 1, 得到 %v", values[0])
	}
}

func TestBackwardCompatibilityWithErrorHandling(t *testing.T) {
	ctx := context.Background()

	// 测试错误处理
	testCases := []struct {
		name        string
		lines       []string
		expectError bool
	}{
		{
			name:        "空输入",
			lines:       []string{},
			expectError: false,
		},
		{
			name:        "只有空行",
			lines:       []string{"", "", ""},
			expectError: false,
		},
		{
			name:        "只有注释",
			lines:       []string{"# comment", "# another comment"},
			expectError: false,
		},
		{
			name:        "无效格式",
			lines:       []string{"invalid yaml format", "no colons"},
			expectError: false, // 解析器应该跳过无效行而不是报错
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := YamlLinesToMap(ctx, tc.lines)

			if tc.expectError {
				if err == nil {
					t.Error("期望错误但没有得到")
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但得到: %v", err)
				}
				// 验证结果不为nil
				if result == nil {
					t.Error("结果不应该为nil")
				}
			}
		})
	}
}

func TestBackwardCompatibilityWithContextCancellation(t *testing.T) {
	// 测试上下文取消
	ctx, cancel := context.WithCancel(context.Background())

	lines := []string{
		"name: test",
		"version: 1.0",
	}

	// 测试YamlLinesToMap在正常上下文下工作
	result, err := YamlLinesToMap(ctx, lines)
	if err != nil {
		t.Errorf("YamlLinesToMap 应该正常工作: %v", err)
	}
	if result["name"] != "test" {
		t.Errorf("期望 name=test, 得到 %v", result["name"])
	}

	// 测试事件处理器中的上下文取消
	eventChan := make(chan SSEvent, 1)
	go func() {
		eventChan <- SSEvent{
			Data: []byte("name: test\n"),
			Err:  nil,
		}
		close(eventChan)
	}()

	result, err = ProcessAIResponseEvents(ctx, eventChan)
	if err != nil {
		t.Errorf("ProcessAIResponseEvents 应该正常工作: %v", err)
	}
	if result["name"] != "test" {
		t.Errorf("期望 name=test, 得到 %v", result["name"])
	}

	// 测试上下文取消
	cancel()

	// 创建新的上下文用于测试取消
	ctx2, cancel2 := context.WithCancel(context.Background())
	cancel2() // 立即取消

	eventChan2 := make(chan SSEvent, 1)
	go func() {
		eventChan2 <- SSEvent{
			Data: []byte("name: test\n"),
			Err:  nil,
		}
		close(eventChan2)
	}()

	_, err = ProcessAIResponseEvents(ctx2, eventChan2)
	if err != nil && err != context.Canceled {
		t.Errorf("上下文取消应该返回 context.Canceled 或 nil, 得到: %v", err)
	}
}

func TestBackwardCompatibilityPerformance(t *testing.T) {
	ctx := context.Background()

	// 性能测试：大量数据
	largeLines := make([]string, 0, 5000)
	largeLines = append(largeLines, "performance_test:")
	for i := 0; i < 1000; i++ {
		largeLines = append(largeLines, fmt.Sprintf("  item_%d:", i))
		largeLines = append(largeLines, fmt.Sprintf("    id: %d", i))
		largeLines = append(largeLines, fmt.Sprintf("    name: item_%d", i))
		largeLines = append(largeLines, "    metadata:")
		largeLines = append(largeLines, fmt.Sprintf("      created: 2024-01-%02d", (i%12)+1))
		largeLines = append(largeLines, fmt.Sprintf("      updated: 2024-01-%02d", (i%12)+1))
		largeLines = append(largeLines, "    tags:")
		for j := 0; j < 3; j++ {
			largeLines = append(largeLines, fmt.Sprintf("      - tag_%d_%d", i, j))
		}
	}

	start := time.Now()
	result, err := YamlLinesToMap(ctx, largeLines)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("性能测试失败: %v", err)
	}

	// 验证性能
	if duration > 100*time.Millisecond {
		t.Errorf("解析耗时过长: %v", duration)
	}

	// 验证结果
	perfTest, ok := result["performance_test"].(map[string]interface{})
	if !ok {
		t.Fatal("performance_test 不是 map")
	}

	itemCount := 0
	for key := range perfTest {
		if strings.HasPrefix(key, "item_") {
			itemCount++
		}
	}

	if itemCount != 1000 {
		t.Errorf("期望 1000 个 item, 得到 %d", itemCount)
	}

	t.Logf("解析 %d 行数据耗时: %v", len(largeLines), duration)
}

func BenchmarkYAMLParser(b *testing.B) {
	processor := NewProcessor(NewDefaultLogger())
	lines := []string{
		"api:",
		"  version: v1",
		"  endpoints:",
		"    - name: users",
		"      path: /api/users",
		"      methods:",
		"        - GET",
		"        - POST",
		"    - name: posts",
		"      path: /api/posts",
		"      methods:",
		"        - GET",
		"        - POST",
		"        - PUT",
		"        - DELETE",
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := processor.ProcessYAMLLines(ctx, lines)
		if err != nil {
			b.Fatalf("处理失败: %v", err)
		}
		jsonData, err := json.Marshal(result)
		if err != nil {
			b.Fatalf("json marshal error: %v", err)
		}
		fmt.Printf("解析结果: %+v\n", string(jsonData))
	}
}

func BenchmarkStringUtils(b *testing.B) {
	utils := NewStringUtils()
	input := "```yaml\nname: test\nversion: 1.0\n```"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.CleanYAMLMarkers(input)
	}
}

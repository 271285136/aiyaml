# AIYAML SDK

一个用于处理AI响应事件流和YAML内容解析的Go语言SDK。

## 概述

AIYAML SDK 是一个模块化的Go语言库，专门用于处理AI服务返回的YAML格式响应。它提供了完整的YAML解析、事件流处理、配置管理等功能，支持实时处理AI生成的YAML内容。

## 主要功能

- **AI响应事件处理**: 实时处理AI服务返回的SSE（Server-Sent Events）流
- **YAML内容解析**: 将YAML格式的文本转换为结构化数据
- **配置管理**: 支持YAML配置文件的加载、验证和保存
- **模块化架构**: 高度解耦的模块设计，便于维护和扩展
- **向后兼容**: 保持原有API的兼容性

## 安装

```bash
go get github.com/271285136/aiyaml
```

## 快速开始

### 基本使用

```go
package main

import (
    "context"
    "fmt"
    "github.com/271285136/aiyaml"
)

func main() {
    // 创建事件通道
    eventChan := make(chan aiyaml.SSEvent)
    
    // 处理AI响应事件
    result, err := aiyaml.ProcessAIResponseEvents(context.Background(), eventChan)
    if err != nil {
        panic(err)
    }
    
    // 使用解析结果
    fmt.Printf("解析结果: %+v\n", result)
}
```

### 直接处理YAML行

```go
package main

import (
    "context"
    "fmt"
    "github.com/271285136/aiyaml"
)

func main() {
    yamlLines := []string{
        "name: test",
        "version: 1.0",
        "settings:",
        "  debug: true",
        "  timeout: 30",
    }
    
    result, err := aiyaml.YamlLinesToMap(context.Background(), yamlLines)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("解析结果: %+v\n", result)
}
```

### 配置管理

```go
package main

import (
    "github.com/271285136/aiyaml"
)

func main() {
    // 创建SDK实例
    sdk := aiyaml.New()
    
    // 加载配置文件
    config, err := sdk.LoadConfig("config.yaml")
    if err != nil {
        panic(err)
    }
    
    // 验证配置
    if err := sdk.ValidateConfig(config); err != nil {
        panic(err)
    }
    
    // 设置API密钥
    sdk.SetAPIKey("your-api-key")
    
    // 保存配置
    err = sdk.SaveConfig("config.yaml", config)
    if err != nil {
        panic(err)
    }
}
```

## 高级使用

### 模块化处理器

```go
package main

import (
    "context"
    "github.com/271285136/aiyaml"
)

func main() {
    // 创建处理器
    processor := aiyaml.NewProcessor(aiyaml.NewDefaultLogger())
    
    // 处理事件
    eventChan := make(chan aiyaml.SSEvent)
    result, err := processor.ProcessAIResponseEvents(context.Background(), eventChan)
    
    // 直接处理YAML行
    lines := []string{"name: test", "version: 1.0"}
    result, err = processor.ProcessYAMLLines(context.Background(), lines)
    
    // 使用工具函数
    stringUtils := processor.GetStringUtils()
    regexPatterns := processor.GetRegexPatterns()
}
```

### 自定义日志

```go
package main

import (
    "context"
    "github.com/271285136/aiyaml"
)

// 自定义日志实现
type CustomLogger struct{}

func (l *CustomLogger) WithContext(ctx context.Context) aiyaml.Logger { return l }
func (l *CustomLogger) WithField(key string, value interface{}) aiyaml.Logger { return l }
func (l *CustomLogger) WithError(err error) aiyaml.Logger { return l }
func (l *CustomLogger) Info(args ...interface{}) {}
func (l *CustomLogger) Infof(format string, args ...interface{}) {}
func (l *CustomLogger) Error(args ...interface{}) {}
func (l *CustomLogger) Errorf(format string, args ...interface{}) {}

func main() {
    // 使用自定义日志
    processor := aiyaml.NewProcessor(&CustomLogger{})
    
    // 处理YAML内容
    lines := []string{"name: test"}
    result, err := processor.ProcessYAMLLines(context.Background(), lines)
}
```

## 配置结构

```yaml
ai:
  model: "gpt-4"
  temperature: 0.7
  max_tokens: 1000
  api_key: "your-api-key"

settings:
  debug: true
  timeout: 30
```

## 使用示例

### 处理AI生成的YAML响应

```go
package main

import (
    "context"
    "fmt"
    "github.com/271285136/aiyaml"
)

func main() {
    // 模拟AI响应事件
    eventChan := make(chan aiyaml.SSEvent)
    
    go func() {
        // 模拟AI返回的YAML内容
        yamlContent := "name: my-app\nversion: 1.0.0\nsettings:\n  debug: true\n  port: 8080"
        
        // 发送事件
        eventChan <- aiyaml.SSEvent{Data: []byte(yamlContent)}
        close(eventChan)
    }()
    
    // 处理事件
    result, err := aiyaml.ProcessAIResponseEvents(context.Background(), eventChan)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("解析结果: %+v\n", result)
    // 输出: map[name:my-app version:1.0.0 settings:map[debug:true port:8080]]
}
```

### 处理复杂嵌套结构

```go
package main

import (
    "context"
    "fmt"
    "github.com/271285136/aiyaml"
)

func main() {
    complexYAML := []string{
        "api:",
        "  version: v1",
        "  endpoints:",
        "    - name: users",
        "      path: /api/users",
        "      methods:",
        "        - GET",
        "        - POST",
        "database:",
        "  host: localhost",
        "  port: 5432",
    }
    
    result, err := aiyaml.YamlLinesToMap(context.Background(), complexYAML)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("复杂结构解析结果: %+v\n", result)
}
```

## 模块说明

### 核心模块

- **`sdk.go`** - 主SDK结构，提供配置管理功能
- **`yaml.go`** - 主入口文件，保持向后兼容的API
- **`processor.go`** - 主处理器，整合所有功能模块
- **`event_processor.go`** - 事件处理模块，处理AI响应事件流
- **`yaml_parser.go`** - YAML解析模块，将YAML行转换为map结构
- **`utils.go`** - 工具函数模块，包含正则表达式和字符串处理
- **`logger.go`** - 日志接口定义
- **`default_logger.go`** - 默认日志实现
- **`types.go`** - 类型定义

### 功能模块

#### EventProcessor
- 负责处理AI响应事件流
- 解析JSON事件数据
- 处理YAML内容的分行和合并逻辑

#### YAMLParser
- 将YAML行转换为map结构
- 处理缩进和层级关系
- 支持数组和嵌套对象

#### StringUtils
- 提供字符串处理工具函数
- 清理YAML标记
- 计算缩进级别
- 解析键值对

#### YAMLRegexPatterns
- 预编译的正则表达式模式
- 用于匹配YAML格式

## 日志接口

SDK 提供了 `Logger` 接口，支持自定义日志实现：

```go
type Logger interface {
    WithContext(ctx context.Context) Logger
    WithField(key string, value interface{}) Logger
    WithError(err error) Logger
    Info(args ...interface{})
    Infof(format string, args ...interface{})
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
}
```

如果没有提供日志实现，SDK 会使用默认的 `DefaultLogger`。

## 测试

运行测试：

```bash
go test ./...
```

运行基准测试：

```bash
go test -bench=. ./...
```

## 特性

1. **模块化设计**: 每个功能都有专门的模块
2. **可测试性**: 可以独立测试每个模块
3. **可扩展性**: 易于添加新功能或修改现有功能
4. **向后兼容**: 保持原有API不变
5. **高性能**: 优化的YAML解析和事件处理
6. **错误处理**: 完善的错误处理和日志记录
7. **配置管理**: 完整的配置加载、验证和保存功能
8. **实时处理**: 支持实时处理AI响应事件流
9. **YAML解析**: 强大的YAML内容解析和转换功能

## 依赖

- Go 1.20+
- `gopkg.in/yaml.v3` - YAML解析库

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request来改进这个项目。

## 更新日志

### v1.0.0
- 初始版本发布
- 支持AI响应事件流处理
- 支持YAML内容解析
- 模块化架构设计
- 配置管理功能 
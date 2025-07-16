package aiyaml

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Config 表示AI配置结构
type Config struct {
	AI       AIConfig       `yaml:"ai"`
	Settings SettingsConfig `yaml:"settings"`
}

// AIConfig AI相关配置
type AIConfig struct {
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
	APIKey      string  `yaml:"api_key,omitempty"`
}

// SettingsConfig 通用设置配置
type SettingsConfig struct {
	Debug   bool `yaml:"debug"`
	Timeout int  `yaml:"timeout"`
}

// SDK 主要的SDK结构
type SDK struct {
	config *Config
}

// New 创建新的SDK实例
func New() *SDK {
	return &SDK{}
}

// LoadConfig 从文件加载配置
func (s *SDK) LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析YAML配置失败: %w", err)
	}

	s.config = &config
	return &config, nil
}

// SaveConfig 保存配置到文件
func (s *SDK) SaveConfig(filename string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// GetConfig 获取当前配置
func (s *SDK) GetConfig() *Config {
	return s.config
}

// ValidateConfig 验证配置
func (s *SDK) ValidateConfig(config *Config) error {
	if config.AI.Model == "" {
		return fmt.Errorf("AI模型名称不能为空")
	}

	if config.AI.Temperature < 0 || config.AI.Temperature > 2 {
		return fmt.Errorf("温度值必须在0-2之间")
	}

	if config.AI.MaxTokens <= 0 {
		return fmt.Errorf("最大令牌数必须大于0")
	}

	if config.Settings.Timeout <= 0 {
		return fmt.Errorf("超时时间必须大于0")
	}

	return nil
}

// SetAPIKey 设置API密钥
func (s *SDK) SetAPIKey(apiKey string) {
	if s.config != nil {
		s.config.AI.APIKey = apiKey
	}
}

// GetAPIKey 获取API密钥
func (s *SDK) GetAPIKey() string {
	if s.config != nil {
		return s.config.AI.APIKey
	}
	return ""
}

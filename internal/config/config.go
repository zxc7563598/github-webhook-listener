package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ActionConfig struct {
	Type    string `yaml:"type"`
	Command string `yaml:"command"`
}

type Rule struct {
	Event    string         `yaml:"event"`
	Branches []string       `yaml:"branches"`
	Actions  []ActionConfig `yaml:"actions"`
}

type RepoConfig struct {
	Secret string `yaml:"secret"`
	Rules  []Rule `yaml:"rules"`
}

type Config struct {
	Repos map[string]*RepoConfig `yaml:"repos"` // key = 仓库名称
}

// LoadConfig 解析 YAML
func LoadConfig(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := yaml.Unmarshal(raw, c); err != nil {
		return nil, err
	}
	// 验证配置
	if err := ValidateConfig(c); err != nil {
		return nil, err
	}
	return c, nil
}

// ValidateConfig 验证配置的有效性
func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("配置不能为空")
	}
	if len(cfg.Repos) == 0 {
		return fmt.Errorf("配置中至少需要一个仓库配置")
	}
	for repoName, repoCfg := range cfg.Repos {
		if repoName == "" {
			return fmt.Errorf("仓库名称不能为空")
		}
		if repoCfg == nil {
			return fmt.Errorf("仓库 %s 的配置不能为空", repoName)
		}
		if repoCfg.Secret == "" {
			return fmt.Errorf("仓库 %s 的 secret 不能为空（安全要求）", repoName)
		}
		if len(repoCfg.Rules) == 0 {
			return fmt.Errorf("仓库 %s 至少需要配置一个规则", repoName)
		}
		for i, rule := range repoCfg.Rules {
			if rule.Event == "" {
				return fmt.Errorf("仓库 %s 的规则 %d 的 event 不能为空", repoName, i)
			}
			if len(rule.Actions) == 0 {
				return fmt.Errorf("仓库 %s 的规则 %d 至少需要一个 action", repoName, i)
			}
			for j, action := range rule.Actions {
				if action.Type == "" {
					return fmt.Errorf("仓库 %s 的规则 %d 的 action %d 的 type 不能为空", repoName, i, j)
				}
				if action.Type == "shell" && action.Command == "" {
					return fmt.Errorf("仓库 %s 的规则 %d 的 shell action %d 的 command 不能为空", repoName, i, j)
				}
			}
		}
	}
	return nil
}

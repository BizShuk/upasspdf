package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Settings 是 settings.json 的完整內容: 每個機構的解密規則.
type Settings struct {
	Entities []Entity `mapstructure:"entities"`
}

// LoadSettings 從 viper 讀出 Settings 並展開每個 entity root_dir 的 ~.
func LoadSettings() (Settings, error) {
	var s Settings
	if err := viper.Unmarshal(&s); err != nil {
		return s, fmt.Errorf("parse settings: %w", err)
	}
	if len(s.Entities) == 0 {
		return s, errors.New("settings has no entities; run `upasspdf config default` to seed settings.json")
	}
	for i := range s.Entities {
		if s.Entities[i].RootDir == "" {
			return s, fmt.Errorf("entity %q missing root_dir", s.Entities[i].Name)
		}
		root, err := homedir.Expand(s.Entities[i].RootDir)
		if err != nil {
			return s, fmt.Errorf("expand root_dir %q for entity %q: %w", s.Entities[i].RootDir, s.Entities[i].Name, err)
		}
		s.Entities[i].RootDir = root
	}
	return s, nil
}

// Select 回傳名稱相符的 entity; name 為空時回傳全部.
func (s Settings) Select(name string) ([]Entity, error) {
	if name == "" {
		return s.Entities, nil
	}
	for _, e := range s.Entities {
		if strings.EqualFold(e.Name, name) {
			return []Entity{e}, nil
		}
	}
	return nil, fmt.Errorf("entity %q not found in settings", name)
}

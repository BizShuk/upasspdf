package model

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Settings 是 settings.json 的完整內容: 根目錄加上每個機構的解密規則.
type Settings struct {
	RootDir  string   `mapstructure:"root_dir"`
	Entities []Entity `mapstructure:"entities"`
}

// LoadSettings 從 viper 讀出 Settings 並展開根目錄的 ~.
func LoadSettings() (Settings, error) {
	var s Settings
	if err := viper.Unmarshal(&s); err != nil {
		return s, fmt.Errorf("parse settings: %w", err)
	}
	root, err := homedir.Expand(s.RootDir)
	if err != nil {
		return s, fmt.Errorf("expand root_dir %q: %w", s.RootDir, err)
	}
	s.RootDir = root
	if len(s.Entities) == 0 {
		return s, errors.New("settings has no entities; run `upasspdf config default` to seed settings.json")
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

// EntityDir 回傳 entity 的資料夾: 絕對路徑原樣使用, 相對路徑接在 root_dir 之下.
func (s Settings) EntityDir(e Entity) string {
	p, err := homedir.Expand(e.Path)
	if err == nil && filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(s.RootDir, e.Path)
}

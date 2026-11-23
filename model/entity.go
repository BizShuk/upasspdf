package model

import (
	"strings"

	"github.com/spf13/viper"
)

// Entity 是一家機構 (銀行, 保險) 的解密規則: 資料夾, 檔名 glob, 以及存放密碼的環境變數名稱.
type Entity struct {
	Name        string `mapstructure:"name"`
	RootDir     string `mapstructure:"root_dir"`
	Filename    string `mapstructure:"filename"`
	PasswordEnv string `mapstructure:"password_env"`
}

// Password 以 PasswordEnv 為 key 讀取密碼 (OS env 或 ~/.config/upasspdf/.env); 未設定時 ok 為 false.
func (e Entity) Password() (password string, ok bool) {
	if e.PasswordEnv == "" {
		return "", false
	}
	password = viper.GetString(strings.ToLower(e.PasswordEnv))
	return password, password != ""
}

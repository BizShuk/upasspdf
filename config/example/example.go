// Package example 內嵌 settings.json 與 .env 的出廠範本, 供 `upasspdf config default` 種入設定目錄.
package example

import (
	_ "embed"

	cfgcmd "github.com/bizshuk/gosdk/cmd/config"
)

// Settings 是 settings.json 的範本, 同時作為 config.Default 的預設值層.
//
//go:embed settings.example.json
var Settings []byte

// Env 是 .env 的範本, 只列出密碼環境變數的名稱, 值一律留空.
//
//go:embed .env.example
var Env []byte

func init() {
	cfgcmd.MustRegisterDefault("settings.json", Settings)
	cfgcmd.MustRegisterDefault(".env", Env)
}

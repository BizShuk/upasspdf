# upasspdf — 技術脈絡 (Technical Context)

業務定義由 [README.md](README.md) 單一擁有.

## 目錄結構 (Directory Layout)

```tree
upasspdf/
├── main.go                  # cmd.Execute()
├── cmd/
│   ├── root.go              # RootCmd, config.Default(WithAppName("upasspdf")), 掛 config/run/decrypt
│   ├── run.go               # run: 依 settings.json entities 批次解密 (--entity, --dry-run)
│   └── decrypt.go           # decrypt: 手動指定 --dir/--password/--filename
├── config/example/          # settings.example.json + .env.example, embed 給 `config default`
├── model/                   # Settings / Entity: viper 反序列化, 路徑解析, 從 env 取密碼
└── svc/
    ├── pdf/                 # 單檔: IsEncrypted, DecryptInPlace (pdfcpu, 保留 mtime)
    │   └── pdftest/         # 測試用明文 / 加密 PDF 產生器
    └── batch/               # 資料夾: glob 篩選, 逐檔解密, Summary 統計
```

## 慣例 (Conventions)

- PDF 處理只用 `github.com/pdfcpu/pdfcpu` (純 Go, 不依賴 qpdf); `api.DisableConfigDir()` 避免 pdfcpu 建立自己的設定目錄.
- 覆寫採 pdfcpu 的 staged output (先寫暫存再 rename), 失敗時原檔不變.
- 密碼只透過 `password_env` 間接取得: `viper.GetString(lower(password_env))`, 來源為 `.env` 或 OS env.
- 使用者提供密碼卻回 `ErrWrongPassword` 時, 會再以空密碼重試一次 (只設 owner 權限限制的檔案).
- 指令: `npm run test` / `npm run lint` / `npm run build`.

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
    ├── pdf/                 # 單檔: IsEncrypted, DecryptInPlace (呼叫 qpdf, 保留 mtime 與權限)
    │   └── pdftest/         # 測試用明文 / 加密 PDF 產生器
    └── batch/               # 資料夾: glob 篩選, 逐檔解密, Summary 統計
```

## 慣例 (Conventions)

- 解密只透過外部 `qpdf` (`brew install qpdf`), 密碼以 `--password-file=-` 從 stdin 傳入, 不出現在 process 參數.
- 不用 pdfcpu: 元大對帳單的 `/Perms` 標記 metadata 不加密, 加密字典卻缺 `/EncryptMetadata false`, pdfcpu 一律拒絕 (`invalid permissions`) 且無放寬選項; qpdf 只發警告 (exit 3) 並正確解密.
- 覆寫流程: qpdf 寫到同資料夾 `.upasspdf-*.pdf` 暫存檔, 成功後還原權限再 rename; 不用 `--replace-input` (遇警告會留下 `.~qpdf-orig`). 失敗時原檔不變.
- 密碼只透過 `password_env` 間接取得: `viper.GetString(lower(password_env))`, 來源為 `.env` 或 OS env.
- 使用者提供密碼卻回 `ErrWrongPassword` 時, 會再以空密碼重試一次 (只設 owner 權限限制的檔案).
- 測試的加密 PDF 由 `qpdf --encrypt` 產生, 因此跑測試也需要 qpdf.
- 指令: `npm run test` / `npm run lint` / `npm run build`.

# upasspdf

就地移除 PDF 的開啟密碼 (open password): 解密後的檔案以`同名`寫回`同一個資料夾`, 保留原本的修改時間.
主要用途是把各家銀行, 保險的電子對帳單轉成不需密碼即可開啟的 PDF, 方便歸檔與後續 OCR.

## 業務領域 (Business Domains)

### 機構批次解密 (Entity Batch Decrypt)

每家機構 (entity) 一組規則, 寫在 `~/.config/upasspdf/settings.json`:

| 欄位 | 說明 |
| --- | --- |
| `name` | 機構代號, `run -e <name>` 用它挑選 |
| `path` | 資料夾; 相對路徑接在 `root_dir` 之下, 絕對路徑原樣使用 |
| `filename` | 檔名 glob, 預設 `*.pdf`; 同一資料夾不同密碼的檔案拆成多個 entity |
| `password_env` | 存放密碼的環境變數名稱, 密碼本身`不寫`在 settings.json |

密碼值放在 `~/.config/upasspdf/.env` (或 OS env), key 就是 `password_env` 的名稱.

`領域流程 (Domain Flow):`

1. `upasspdf run` 讀 settings.json, 逐一處理每個 entity.
2. 列出資料夾內 (不遞迴) 符合 `filename` 的 PDF; 未加密的檔案略過 (skipped).
3. 以 `password_env` 指向的密碼解密, 覆寫同名檔案.
4. 密碼環境變數未設定的 entity 只列出待解密檔案, 不視為失敗.
5. 任何檔案解密失敗 (密碼錯誤, 檔案損壞) 時 exit code 為 1, 原檔保持不變.

重複執行是冪等的: 已解密的檔案下次會被略過.

### 手動解密 (Manual Decrypt)

不讀 entity 設定, 直接指定資料夾與密碼:

```bash
upasspdf decrypt --dir ~/Downloads/bank-estatements/中國信託_CTBC --password <密碼>
upasspdf decrypt -d ./statements -p <密碼> -f "*_Estatement_*.pdf" --dry-run
```

## 安裝與設定 (Setup)

```bash
npm run install                     # go install + 種入 ~/.config/upasspdf/{settings.json,.env}
upasspdf config --file .env --update UPASSPDF_CTBC_PASSWORD=<密碼>
upasspdf run --dry-run              # 先看每個 entity 待解密的檔案
upasspdf run                        # 實際解密
```

排程由 `~/projects/default/ecosystem.config.js` 的 `Unlock Bank Statement PDF - Daily` 任務執行 (`optional: true`, 預設 paused).

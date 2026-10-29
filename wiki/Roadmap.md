# 🗺️ 迭代路線圖 (Roadmap)

本頁面記載 Gbootcms 的版本演進與功能規劃清單。

---

## 📅 迭代總覽

| 階段 | 規劃模組 | 對應 Issue | 狀態 | 核心目標 |
| :--- | :--- | :--- | :--- | :--- |
| **Iteration 1** | 媒體庫路徑標準化與 WebP 轉碼 | [#1](https://github.com/waicun/Gbootcms/issues/1) | **已驗收 (Done)** | 消除資料庫與物理路徑歧異，純 Go 伺服器端 WebP 自動轉碼 |
| **Iteration 2** | 邊緣 CDN 快取聯動與分級快取 | [#2](https://github.com/waicun/Gbootcms/issues/2) | **規劃中 (Backlog)** | Cloudflare Cache-Tag 聯動 Purge，支援 Redis / 記憶體快取驅動 |
| **Iteration 3** | 全站動態 Markdown 鏡像與 AI 檢索端點 | [#3](https://github.com/waicun/Gbootcms/issues/3) | **規劃中 (Backlog)** | 支援 `/content/{id}.md` 端點，維護 `/llms.txt` 索引與 GSC 推送 |
| **Iteration 4** | 雙因子驗證 (2FA / TOTP) 與暴力破解防護 | [#4](https://github.com/waicun/Gbootcms/issues/4) | **規劃中 (Backlog)** | RFC 6238 TOTP 認證，登入失敗次數超限自動封鎖客戶端 IP |
| **Iteration 5** | SQLite 線上快照與物件儲存備份 | [#5](https://github.com/waicun/Gbootcms/issues/5) | **規劃中 (Backlog)** | 使用 `VACUUM INTO` 進行無鎖熱快照，同步至 S3 / R2 儲存桶 |

---

## 📌 模組技術規範

### Iteration 1：媒體庫健全度與 WebP 轉碼（已完成）
- 統一資料庫儲存路徑（`upload/...`）與磁碟讀取路徑（`static/upload/...`）。
- 導入純 Go `gowebp` 模組，於上傳時自動進行格式轉換與下採樣。
- 支援 WebP 格式尺寸與 EXIF 資訊解析。

### Iteration 2：邊緣 CDN 快取聯動與分級快取
- 響應標頭附加 `Cache-Tag: content-{id}, sort-{scode}, site-{acode}`。
- 後台內容更新時非阻塞非同步調用 Cloudflare Purge API。
- 抽象標準 `CacheDriver` 介面，支援記憶體快取與 Redis 連線池。

### Iteration 3：全站動態 Markdown 鏡像與 AI 檢索
- 詳情頁路由支援 `.md` 後綴或 `Accept: text/markdown` 請求頭輸出純文字。
- 自動生成 `/llms.txt` 與 `/llms-full.txt` 結構。
- 整合 Search Console Indexing API 主動推送 URL 更新。

### Iteration 4：管理員雙因子驗證 (2FA / TOTP)
- 基於 RFC 6238 規範實現 TOTP 動態碼校驗（相容 Authenticator / 1Password）。
- 連續 5 次密碼錯誤自動將 IP 寫入封鎖表 30 分鐘，並記錄安全日誌。

### Iteration 5：SQLite 線上快照備份
- 呼叫 SQLite 3 `VACUUM INTO` 語法產出一致性檔案快照。
- 支援標準 S3 協定（Cloudflare R2 / AWS S3 / MinIO）自動定時上傳。
- 後台提供快照清單瀏覽、手動下載與狀態檢視。

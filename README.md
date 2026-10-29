# 🚀 Gbootcms

基於 Go 語言重構的高性能通用內容管理系統（CMS），相容 PbootCMS 模板標籤語法與 `ay_` 資料表結構。

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Gin Framework](https://img.shields.io/badge/Gin-v1.12-00ADD8?style=flat-square&logo=gin)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-v1.31-00ADD8?style=flat-square)](https://gorm.io/)
[![SQLite](https://img.shields.io/badge/SQLite-WAL_Mode-003B57?style=flat-square&logo=sqlite)](https://github.com/glebarez/sqlite)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)

---

## ✨ 核心特性

- 📦 **單一二進制分發**：`CGO_ENABLED=0` 純 Go 靜態編譯，零外部 C 函式庫或 PHP-FPM 執行環境依賴，跨平台開箱即用。
- ⚡ **SQLite WAL 高併發儲存**：採用 Write-Ahead Logging 模式，實現讀寫分離，高併發讀取不鎖表，鎖競爭自動重試。
- 🖼️ **伺服器端原生圖片轉碼**：上傳時自動進行雙三次平滑縮放與 WebP 轉碼，有效降低 60%~85% 磁碟儲存與 CDN 流量成本。
- 🌐 **原生多語言前綴路由**：提供 `/`（繁體中文）、`/en/`（英文）、`/sc/`（簡體中文）獨立路徑劃分，資料庫分區隔離與專屬模板支援。
- 🔍 **SEO & AI 檢索原生配套**：
  - 🛡️ **動態環境感知 Robots**：測試環境自動輸出 `Disallow: /` 與 `noindex` 標頭；正式環境自動放行。
  - 🗺️ **多語系 XML Sitemap**：提供 `/sitemap.xml` 總索引及各語言獨立分區子索引。
  - 📋 **Schema.org JSON-LD**：自動輸出標準結構化資料，`@id` 與 `canonical` 嚴格等同。
  - 🔗 **精準雙向 hreflang**：基於資料庫關聯代碼（`rel_id`）精確映射跨語言內容。
  - 🤖 **大模型檢索端點**：標準提供 `/llms.txt` 與 `/llms-full.txt` 端點供 AI Agent 索引。
- 🎛️ **開箱即用後台與 RESTful API**：內建欄目、文章、媒體庫、自訂表單、留言與日誌模組；提供前後端分離 API（支援 JWT 與 API Key 雙重認證）。

---

## ⚡ 快速開始

### 1. 預編譯版本執行

至 [Releases](https://github.com/waicun/Gbootcms/releases) 下載適用作業系統之發行包：

```bash
# Linux (amd64)
tar -xzf gbootcms-linux.tar.gz
cd gbootcms-linux
chmod +x start.sh
./start.sh

# Windows (x64)
# 解壓 gbootcms-windows.zip 後點擊執行 start.bat
```

### 2. 源碼本機編譯

```bash
# 取得源碼
git clone https://github.com/waicun/Gbootcms.git
cd Gbootcms

# 下載依賴並編譯
go mod download
go build -trimpath -ldflags="-s -w" -o bin/gbootcms .

# 啟動服務
./bin/gbootcms
```

### 3. 系統訪問入口

啟動完成後預設監聽連接埠 `8080`：

| 介面 | 訪問 URL | 預設憑據 |
| :--- | :--- | :--- |
| **前台首頁** | `http://localhost:8080/` | - |
| **管理後台** | `http://localhost:8080/admin` | 帳號：`admin`<br>密碼：`123456` |

---

## 🏛️ 技術架構

| 組件層次 | 選型方案 | 說明 |
| :--- | :--- | :--- |
| **開發語言** | Go 1.25 | 單一靜態二進制部署，記憶體開銷低 |
| **Web 核心** | Gin v1.12 | 高效 HTTP 路由分發與中介軟體管道 |
| **資料庫引擎** | SQLite 3 (`glebarez/sqlite`) | WAL 併發模式，純 Go 驅動無須 GCC |
| **ORM 映射** | GORM v1.31 | 關係映射，相容原版 `ay_` 前綴表結構 |
| **圖片處理** | `gowebp` + `x/image` | 原生 WebP 編碼與雙三次等比縮放 |
| **前台視圖** | TagParser 解析引擎 | 相容 `{gboot:xxx}` 標籤語法體系 |
| **後台視圖** | Pongo2 v6.1 | Django 風格語義 AST 模板引擎 |

---

## 📂 目錄結構

```
Gbootcms/
├── apps/               # 業務模組 (admin, home, api, common)
│   ├── admin/          # 後台管理 (Controller, Model, View)
│   ├── home/           # 前台展示 (FrontController)
│   ├── api/            # RESTful API 端點
│   └── common/         # 圖片壓縮、標籤解析、中介軟體
├── bin/                # 編譯二進制產物
├── config/             # 系統設定檔 (config.json)
├── core/               # 核心底層 (db 連線池, 系統遙測, 外掛)
├── data/               # SQLite 資料庫儲存目錄
├── static/             # 後台樣式與上傳資源
├── template/           # 前台模板目錄 (default, en, sc)
├── wiki/               # 本地技術文檔庫
├── main.go             # 程式進入點
└── LICENSE             # MIT 授權條款
```

---

## 📚 技術手冊與支援

- 完整開發與標籤手冊請參閱 [Gbootcms Wiki](https://github.com/waicun/Gbootcms/wiki)。
- 缺陷反饋與功能規劃請至 [Issues](https://github.com/waicun/Gbootcms/issues)。

---

## 📄 授權條款

本項目採用 [MIT License](LICENSE) 授權分發。

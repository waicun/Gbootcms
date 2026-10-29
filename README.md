# Gbootcms

基於 Go 語言開發的通用內容管理系統（CMS），相容 PbootCMS 模板標籤語法與 `ay_` 資料表結構。

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-v1.12-00ADD8)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-v1.31-00ADD8)](https://gorm.io/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## 產品特性

- **單一二進制分發**：原生編譯，`CGO_ENABLED=0` 零動態連結庫依賴，跨平台開箱即用。
- **高併發資料存儲**：採用 SQLite WAL 模式，支援併發讀取與安全寫入。
- **自動化圖片轉碼**：上傳時自動進行雙三次平滑縮放與 WebP 轉碼，降低磁碟占用與傳輸頻寬。
- **原生多語言架構**：提供路徑前綴路由（`/`、`/en/`、`/sc/`），資料分區隔離與模板字典渲染。
- **搜尋引擎與 AI 檢索**：
  - 環境感知 `robots.txt` 與多語系 XML `sitemap.xml`。
  - 標準 Schema.org JSON-LD 結構化資料。
  - 雙向精確 `hreflang` 多語言對齊。
  - AI 爬蟲檢索端點（`/llms.txt`、`/llms-full.txt`）。
- **完整後台與 API**：包含欄目、內容、媒體庫、自訂表單、日誌管理，提供 RESTful API（支援 JWT 與 API Key 驗證）。

---

## 快速開始

### 1. 下載預編譯版本

至 [Releases](https://github.com/waicun/Gbootcms/releases) 下載對應作業系統的壓縮包：

```bash
# Linux
tar -xzf gbootcms-linux.tar.gz
cd gbootcms-linux
chmod +x start.sh
./start.sh

# Windows
# 解壓 gbootcms-windows.zip 後執行 start.bat
```

### 2. 源碼編譯

```bash
# 取得源碼
git clone https://github.com/waicun/Gbootcms.git
cd Gbootcms

# 編譯
go build -trimpath -ldflags="-s -w" -o bin/gbootcms .

# 運行
./bin/gbootcms
```

### 3. 系統訪問

- 前台首頁：`http://localhost:8080/`
- 後台管理：`http://localhost:8080/admin`
- 預設憑據：帳號 `admin` / 密碼 `123456`

---

## 技術架構

| 組件 | 選型 | 說明 |
| :--- | :--- | :--- |
| **開發語言** | Go 1.25 | 單一靜態二進制部署 |
| **Web 框架** | Gin v1.12 | HTTP 路由與中介軟體處理 |
| **資料庫引擎** | SQLite 3 (`glebarez/sqlite`) | WAL 併發模式，純 Go 驅動 |
| **ORM** | GORM v1.31 | 關係映射，相容 `ay_` 資料表 |
| **圖片處理** | `gowebp` + `x/image` | 原生 WebP 編碼與圖像縮放 |
| **前台模板** | TagParser | 相容 `{gboot:xxx}` 標籤語法 |
| **後台模板** | Pongo2 v6.1 | 後台視圖渲染 |

---

## 目錄結構

```
Gbootcms/
├── apps/               # 業務模組 (admin, home, api, common)
├── bin/                # 編譯產物
├── config/             # 系統設定檔
├── core/               # 核心底層 (db, version, plugins)
├── data/               # SQLite 資料庫目錄
├── static/             # 靜態資源 (css, js, upload)
├── template/           # 前台視圖模板
├── wiki/               # 技術文檔
├── main.go             # 程式進入點
└── LICENSE             # MIT 授權條款
```

---

## 相關文檔

- 完整技術文檔請參閱 [Gbootcms Wiki](https://github.com/waicun/Gbootcms/wiki)。
- 問題反饋與功能需求請提交至 [Issues](https://github.com/waicun/Gbootcms/issues)。

---

## 授權條款

本項目採用 [MIT License](LICENSE) 授權。

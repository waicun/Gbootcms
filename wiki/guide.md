# 🧭 Gbootcms 全景技術與開發手冊 (Full Guide)

本手冊提供 Gbootcms 的系統全景技術說明，涵蓋快速入門、底層架構、模板標籤、多語言、SEO 與部署規範。

---

## 📑 目錄導覽

1. [⚡ 快速入門 (Quick Start)](#1-⚡-快速入門-quick-start)
2. [🏛️ 核心架構 (Architecture)](#2-🏛️-核心架構-architecture)
3. [🖼️ 原生圖片 WebP 轉碼 (Image Processing)](#3-🖼️-原生圖片-webp-轉碼-image-processing)
4. [🏷️ 模板標籤解析 (TagParser)](#4-🏷️-模板標籤解析-tagparser)
5. [🌐 多語言路由劃分 (Multilingual)](#5-🌐-多語言路由劃分-multilingual)
6. [🤖 SEO 與 AI 檢索規範 (SEO & GEO)](#6-🤖-seo-與-ai-檢索規範-seo--geo)
7. [🚢 生產部署指引 (Deployment)](#7-🚢-生產部署指引-deployment)
8. [📋 開發規範與約束 (Standards)](#8-📋-開發規範與約束-standards)

---

## 1. ⚡ 快速入門 (Quick Start)

### 1.1 環境準備
- **Go 編譯環境**：Go 1.25 或以上。
- **CGO**：`CGO_ENABLED=0`，純 Go SQLite 驅動，不依賴 GCC。

### 1.2 源碼編譯與啟動
```powershell
# 編譯二進制產物至 bin/ 目錄
go build -trimpath -ldflags="-s -w" -o bin/gbootcms.exe .

# 啟動服務
.\bin\gbootcms.exe
```

### 1.3 存取入口
- 前台首頁：`http://localhost:8080/`
- 後台管理：`http://localhost:8080/admin`
- 預設管理員帳號：`admin` / 密碼：`123456`

---

## 2. 🏛️ 核心架構 (Architecture)

### 2.1 資料庫連線池 (`core/db/db.go`)
全系統維護單一 DB 連線池，採用 SQLite WAL 模式：

```go
db.Exec("PRAGMA journal_mode = WAL;")     // WAL 模式：併發讀與寫入互不阻塞
db.Exec("PRAGMA synchronous = NORMAL;")  // 平衡 I/O 速度與資料安全
db.Exec("PRAGMA busy_timeout = 5000;")    // 鎖競爭等待上限 5 秒
db.Exec("PRAGMA cache_size = -64000;")   // 64MB 記憶體頁快取
```

### 2.2 結構遷移收斂 (`apps/admin/model/migrate.go`)
所有資料模型透過 `MigrateAll()` 執行自動遷移：

```go
func MigrateAll() error {
    if err := system.AutoMigrate(DB); err != nil { return err }
    if err := content.AutoMigrate(DB); err != nil { return err }
    if err := member.AutoMigrate(DB); err != nil { return err }
    return nil
}
```

---

## 3. 🖼️ 原生圖片 WebP 轉碼 (Image Processing)

伺服器端原生處理管道（`apps/common/imgcompress`）於上傳時自動執行轉碼：
- **全入口適用**：後台單圖、多圖相簿、富文本編輯器與 API 上傳均自動轉碼。
- **純 Go 實現**：基於 `gowebp` 與 `golang.org/x/image/draw.BiCubic`。
- **尺寸約束**：超過最大寬度限制時自動進行雙三次平滑縮放。

```go
newPath, err := imgcompress.CompressToWebP(savePath, imgcompress.Options{
    Quality:  80,   // 壓縮品質 1-100
    MaxWidth: 1920, // 超出寬度等比縮放
})
```

---

## 4. 🏷️ 模板標籤解析 (TagParser)

### 4.1 導航標籤 (`{gboot:nav}`)
```html
<ul class="navbar">
    <li><a href="/">首頁</a></li>
    {gboot:nav parent=0}
    <li class="{if('[nav:scode]'=='{gboot:sort:tcode}')}active{/if}">
        <a href="[nav:link]">[nav:name]</a>
    </li>
    {/gboot:nav}
</ul>
```

### 4.2 內容列表標籤 (`{gboot:list}`)
```html
{gboot:list scode=2 num=10 order='date desc'}
<div class="article-card">
    <a href="[list:link]">
        <img src="[list:ico]" alt="[list:title]">
        <h3>[list:title]</h3>
        <p>[list:description]</p>
        <span class="time">[list:date style='Y-m-d']</span>
    </a>
</div>
{/gboot:list}

<div class="page-bar">
    {gboot:pagebar}
</div>
```

---

## 5. 🌐 多語言路由劃分 (Multilingual)

前台根據 URL 路徑前綴識別當前語言代碼（`acode`）：
- `/`：繁體中文站點 (`tc`)
- `/en/`：英文站點 (`en`)
- `/sc/`：簡體中文站點 (`sc`)

由 `acodeplugin` 自動解析前綴並注入上下文，GORM 查詢自動加入語言條件。各語言使用專屬模板目錄（`template/default`、`template/en`、`template/sc`），靜態資源由 `template/default/static/` 統一提供。

---

## 6. 🤖 SEO 與 AI 檢索規範 (SEO & GEO)

1. **動態環境感知 Robots.txt**：
   - 測試環境自動返回 `Disallow: /` 與 `X-Robots-Tag: noindex, nofollow` 響應頭。
   - 正式生產環境放行並關聯 `Sitemap: https://{Host}/sitemap.xml`。
2. **結構化資料 JSON-LD (`{gboot:schema}`)**：
   - 輸出 Schema.org 標準 `WebSite` 與 `Article` 結構。
   - `mainEntityOfPage.@id` 數值等同於 `<link rel="canonical">`。
3. **雙向跨語言關聯 (`{gboot:hreflang}`)**：
   - 根據資料庫 `rel_id` 精準輸出各語言 alternate 標籤。
4. **AI 檢索端點**：
   - 提供 `/llms.txt` 與 `/llms-full.txt` 端點支援大模型爬蟲索引。

---

## 7. 🚢 生產部署指引 (Deployment)

### 7.1 Linux 一鍵打包 (`build-linux.bat`)
```powershell
.\build-linux.bat
# 產物：dist/gbootcms-linux.tar.gz（體積約 12MB，已執行 UPX 壓縮）
```

### 7.2 Systemd 服務管理
將二進制部署於 `/opt/gbootcms`，使用 systemd 管理行程：
```ini
[Unit]
Description=Gbootcms Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/gbootcms
ExecStart=/opt/gbootcms/bin/gbootcms
Restart=always

[Install]
WantedBy=multi-user.target
```

---

## 8. 📋 開發規範與約束 (Standards)

1. **保留原版表結構**：嚴格相容 `ay_` 資料表前綴與既有欄位名稱。
2. **錯誤檢查閉環**：所有資料庫讀寫操作必須顯式檢查 `.Error`。
3. **多語言查詢綁定 Context**：模型查詢時使用 `db.DB.WithContext(c.Request.Context())`。
4. **標準繁體中文文檔**：文檔與代碼註釋保持客觀、簡潔，杜絕誇飾修辭。

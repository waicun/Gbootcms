# 🌐 多語言路由架構 (Multilingual)

Gbootcms 提供多語言分區架構，透過「路徑前綴劃分」與「語言獨立模板」進行內容隔離與多語系展示。

---

## 🧭 路由前綴規格

前台路由器根據 HTTP 請求路徑前綴識別當前語言代碼（`acode`）：

| 區域代碼 (`acode`) | 語言標籤 | URL 路徑規則 | 範例請求 |
| :--- | :--- | :--- | :--- |
| `tc` (預設) | 繁體中文 | `/` | `/about.html`, `/news/12.html` |
| `en` | English | `/en/` | `/en/about.html`, `/en/news/12.html` |
| `sc` | 簡體中文 | `/sc/` | `/sc/about.html`, `/sc/news/12.html` |

底層由 `acodeplugin` 中介軟體解析前綴並寫入請求上下文（`gin.Context`），GORM 查詢自動加入 `acode = ?` 資料分區過濾條件。

---

## 📁 模板目錄配置

各語言配置對應之視圖模板目錄，支援獨立 HTML 文案與標籤佈局：

```
template/
├── default/                 # 預設模板（繁體中文）
│   ├── index.html
│   ├── product.html
│   └── comm/
├── en/                      # 英文專屬模板
│   ├── index.html
│   ├── product.html
│   └── comm/
└── sc/                      # 簡體中文專屬模板
    ├── index.html
    ├── product.html
    └── comm/
```

> [!NOTE]
> 公共 CSS、JavaScript 及字型資源由 `template/default/static/` 統一維護，系統於各語言路徑（如 `/template/en/static/...`）自動進行虛擬路徑映射。

---

## 🏷️ 模板套用語法範例

### 繁體中文模板 (`template/default/product.html`)

```html
<div class="meta-item">上架時間：{content:date style=Y-m-d}</div>
<div class="meta-item">瀏覽次數：{content:visits}</div>
<div class="nav-links">
    <p>上一篇：{content:precontent}</p>
    <p>下一篇：{content:nextcontent}</p>
</div>
```

### 英文模板 (`template/en/product.html`)

```html
<div class="meta-item">Release Date: {content:date style=Y-m-d}</div>
<div class="meta-item">Views: {content:visits}</div>
<div class="nav-links">
    <p>Previous: {content:precontent}</p>
    <p>Next: {content:nextcontent}</p>
</div>
```

---

## 🔗 跨語言關聯 (`hreflang`)

多語言內容透過資料庫 `rel_id` 欄位建立多對多關聯：

1. 詳情頁調用 `{gboot:hreflang}` 標籤。
2. 系統依據當前文章之 `rel_id` 檢索其他語言區域之對應文章。
3. 輸出標準 HTML 關聯標註：

```html
<link rel="alternate" hreflang="zh-Hant" href="https://example.com/product/10.html" />
<link rel="alternate" hreflang="en" href="https://example.com/en/product/11.html" />
<link rel="alternate" hreflang="x-default" href="https://example.com/product/10.html" />
```

若某語言無關聯資料，則不輸出該語言標註，避免產生無效鏈結。

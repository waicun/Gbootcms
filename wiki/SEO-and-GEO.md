# 🤖 SEO 與 AI 檢索規範 (SEO & GEO)

Gbootcms 內建搜尋引擎最佳化（SEO）與生成式引擎檢索（GEO）支援機制，提供動態環境感知、結構化資料輸出與大模型端點。

---

## 🛡️ 動態環境感知 Robots.txt

系統根據連線請求的 `Host` 與後台配置的正式網域自動判定環境狀態：

```go
func RobotsHandler(c *gin.Context) {
    if IsStagingHost(c) {
        // 測試/開發環境：阻擋所有爬蟲抓取，並附加 HTTP 響應頭
        c.Header("X-Robots-Tag", "noindex, nofollow")
        c.String(http.StatusOK, "User-agent: *\nDisallow: /\n")
        return
    }

    // 正式生產環境：放行爬蟲並輸出 Sitemap 索引位址
    host := GetRequestBaseURL(c)
    c.String(http.StatusOK, fmt.Sprintf(
        "User-agent: *\nAllow: /\nDisallow: /admin/\nDisallow: /api/\n\nSitemap: %s/sitemap.xml\n", 
        host,
    ))
}
```

---

## 🗺️ 多語系 XML Sitemap

系統自動註冊並動態輸出符合 Sitemaps.org 規範的 XML 索引：

| 端點 URL | 功能說明 |
| :--- | :--- |
| `/sitemap.xml` | 站點 Sitemap 總索引（Sitemap Index），分發各語言子索引 |
| `/sitemap-tc.xml` | 繁體中文區域欄目與文章 URL 索引 |
| `/sitemap-en.xml` | 英文區域 URL 索引 |
| `/sitemap-sc.xml` | 簡體中文區域 URL 索引 |

---

## 📋 JSON-LD 結構化資料 (`Schema.org`)

前台模板引用 `{gboot:schema}` 標籤時，自動依據頁面類型生成結構化 JSON-LD：

```html
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@graph": [
    {
      "@type": "WebSite",
      "@id": "https://example.com/#website",
      "url": "https://example.com/",
      "name": "站點名稱"
    },
    {
      "@type": "Article",
      "@id": "https://example.com/news/15.html",
      "isPartOf": { "@id": "https://example.com/#website" },
      "headline": "文章標題",
      "datePublished": "2026-10-23T08:00:00+08:00",
      "mainEntityOfPage": {
        "@type": "WebPage",
        "@id": "https://example.com/news/15.html"
      }
    }
  ]
}
</script>
```

> [!NOTE]
> JSON-LD 之 `mainEntityOfPage.@id` 與頁面標籤 `<link rel="canonical" href="...">` 保持相同數值。

---

## 🤖 大模型檢索端點 (`/llms.txt`)

為 AI Agent 與 LLM 爬蟲提供低資源消耗之純文字索引：

- `/llms.txt`：精簡版全站知識庫導航與重點欄目描述。
- `/llms-full.txt`：完整 Markdown 格式文檔聚合輸出。

---

## 🧹 訪客歸因路徑過濾

訪客著陸頁面與 UTM 參數追蹤模組（`apps/common/attribution.go`）自動忽略以下系統端點，避免爬蟲抓取污染歸因統計：

- `/robots.txt`、`/sitemap*.xml`、`/llms*.txt`
- `/api/*`、`/admin/*`、`/static/*`、`/upload/*`、`/favicon.ico`

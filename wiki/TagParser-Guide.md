# 🏷️ 標籤解析語法 (TagParser Guide)

Gbootcms 前台模板引擎使用 `TagParser` 解析器，相容 `{gboot:xxx}` 標籤語法，支援欄目遞迴、內容列表分頁與通用變量解析。

---

## 🌐 站點與公司全域標籤

在全域模板中存取系統基礎設定值：

| 標籤 | 輸出內容 | 範例 |
| :--- | :--- | :--- |
| `{gboot:site:title}` | 網站主標題 | `Gbootcms` |
| `{gboot:site:subtitle}` | 網站副標題 | `企業內容管理系統` |
| `{gboot:site:keywords}` | 站點關鍵字 | `Go CMS, 高性能` |
| `{gboot:site:description}` | 站點描述說明 | `基於 Go 語言的內容管理系統` |
| `{gboot:site:logo}` | 站點 Logo 圖片路徑 | `/static/upload/...` |
| `{gboot:company:name}` | 公司名稱 | `企業網絡科技` |
| `{gboot:company:phone}` | 聯絡電話 | `0755-88888888` |
| `{gboot:company:address}` | 公司地址 | `資訊科技園區` |

### 範例：HTML Head 區塊

```html
<head>
    <title>{gboot:site:title} - {gboot:site:subtitle}</title>
    <meta name="keywords" content="{gboot:site:keywords}">
    <meta name="description" content="{gboot:site:description}">
    <link rel="icon" href="{gboot:site:logo}">
</head>
```

---

## 🧭 導航選單標籤 (`{gboot:nav}`)

### 標籤參數

| 參數 | 類型 | 預設值 | 說明 |
| :--- | :--- | :--- | :--- |
| `parent` | 數字 | `0` | 父欄目編號，`0` 表示頂級導航 |
| `num` | 數字 | 全部 | 提取導航筆數限制 |

### 程式碼範例

```html
<ul class="nav">
    <li class="{if('{gboot:sort:scode}'=='')}active{/if}"><a href="/">首頁</a></li>
    {gboot:nav parent=0}
    <li class="{if('[nav:scode]'=='{gboot:sort:tcode}')}active{/if}">
        <a href="[nav:link]">[nav:name]</a>
        {if([nav:soncount]>0)}
        <ul class="sub-nav">
            {gboot:nav parent=[nav:scode]}
            <li><a href="[nav:link]">[nav:name]</a></li>
            {/gboot:nav}
        </ul>
        {/if}
    </li>
    {/gboot:nav}
</ul>
```

---

## 📰 內容列表標籤 (`{gboot:list}`)

用於輸出指定欄目或全站的文章與產品清單，支援分頁與多欄位排序。

### 標籤參數

| 參數 | 類型 | 預設值 | 說明 |
| :--- | :--- | :--- | :--- |
| `scode` | 數字/逗號字串 | 當前欄目 | 目標欄目代碼，支援 `scode=2,3` |
| `num` | 數字 | `10` | 每一頁或單批提取筆數 |
| `order` | 字串 | `sorting asc, id desc` | 排序規則，如 `date desc`、`visits desc` |
| `page` | 數字 | `1` | 是否啟用分頁計算（`0` 關閉，`1` 開啟） |

### 程式碼範例

```html
{gboot:list scode=2 num=10 order='date desc'}
<div class="news-item">
    <a href="[list:link]">
        <img src="[list:ico]" alt="[list:title]">
        <h3>[list:title]</h3>
        <p>[list:description]</p>
        <span class="date">[list:date style='Y-m-d']</span>
    </a>
</div>
{/gboot:list}

<!-- 分頁導航條 -->
<div class="pagination">
    {gboot:pagebar}
</div>
```

---

## 📄 內容詳情頁標籤 (`{gboot:content:xxx}`)

於文章詳情頁（`content.html`）提取當前紀錄數值：

```html
<article class="article-detail">
    <h1>{gboot:content:title}</h1>
    <div class="meta">
        <span>發布時間：{gboot:content:date style='Y-m-d H:i'}</span>
        <span>瀏覽次數：{gboot:content:visits}</span>
        <span>作者：{gboot:content:author}</span>
    </div>
    <div class="content-body">
        {gboot:content:content}
    </div>
    
    <div class="prev-next">
        <p>上一篇：{gboot:content:pre}</p>
        <p>下一篇：{gboot:content:next}</p>
    </div>
</article>
```

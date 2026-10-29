# 📚 Gbootcms 技術文檔庫

歡迎查閱 **Gbootcms** 官方文檔。Gbootcms 為基於 Go 語言開發的高性能通用內容管理系統，支援單一靜態二進制分發、SQLite WAL 併發儲存與 PbootCMS 原生標籤相容。

---

## ⚡ 核心優勢

| 特性模組 | 技術實現 | 交付價值 |
| :--- | :--- | :--- |
| **部署模型** | 純 Go 靜態編譯 (`CGO_ENABLED=0`) | 零動態鏈結庫依賴，單檔案跨平台部署 |
| **資料存儲** | SQLite 3 WAL 模式 (`glebarez/sqlite`) | 支援讀寫分離與高併發存取，無需安裝獨立資料庫服務 |
| **圖片處理** | 伺服器端原生 WebP 編碼與平滑縮放 | 自動壓縮上傳圖片，降低存儲空間與 CDN 流量成本 |
| **多語言支援** | 路由前綴隔離 (`/`, `/en/`, `/sc/`) | 支援多語系獨立內容分區與動態字典替換 |
| **檢索最佳化** | 內建 Schema.org JSON-LD 與 hreflang | 提高搜尋引擎爬蟲解析精確度，提供標準 `/llms.txt` 接口 |

---

## 🧭 文檔目錄

### 1. 入門與部署
- [⚡ 快速上手 (Quick Start)](Quick-Start)：環境需求、二進制編譯與基本運行流程。
- [🧭 全景指南 (Full Guide)](guide)：涵蓋從開發、配置到維運的一站式技術指南。
- [🚢 生產部署 (Deployment)](Deployment-Guide)：Linux 服務配置、極限壓縮打包與自動化發布。

### 2. 系統架構與特性
- [📦 系統架構 (Architecture)](Architecture)：連線池架構、資料遷移機制與系統遙測。
- [🌐 多語言路由 (Multilingual)](Multilingual)：語言分區路由策略、資料庫外鍵關聯與模板解析。
- [🖼️ 圖片轉碼管線 (Image Processing)](Image-Compression)：圖像上傳自動化轉碼、規格約束與演算法說明。

### 3. 視圖與標籤
- [🏷️ 標籤解析語法 (TagParser)](TagParser-Guide)：前台動態標籤、分頁、文章列表與欄目解析規範。
- [🤖 SEO 與 AI 檢索 (SEO & GEO)](SEO-and-GEO)：動態 `robots.txt`、結構化資料輸出與大模型索引配置。

---

## 💻 快速呼叫範例

在自訂業務控制器中存取全域資料庫連線池：

```go
package hello

import (
	"gbootcms/apps/common"
	"gbootcms/core/db"
	"github.com/gin-gonic/gin"
)

type HelloController struct {
	common.BaseController
}

func (ctrl *HelloController) Index(c *gin.Context) {
	var total int64
	db.DB.Table("ay_content").Where("status = ?", 1).Count(&total)

	ctrl.JSONOK(c, gin.H{
		"code":    1,
		"message": "success",
		"total":   total,
	})
}
```

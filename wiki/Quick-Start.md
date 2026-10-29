# ⚡ 快速上手 (Quick Start)

本指南說明如何在本地環境完成 Gbootcms 的源碼獲取、依賴安裝、編譯與初次運行。

---

## 🛠️ 環境需求

- **Go 執行環境**：Go 1.25 或以上版本。
- **CGO 需求**：`CGO_ENABLED=0`（系統使用純 Go SQLite 驅動，不依賴 GCC 或本機 C 編譯器）。
- **支援平台**：Windows (x64)、Linux (amd64 / arm64)、macOS (amd64 / arm64)。

---

## 📥 源碼獲取

```bash
git clone https://github.com/waicun/Gbootcms.git
cd Gbootcms
go mod download
```

---

## 🏗️ 編譯與運行

### Windows (PowerShell)

```powershell
# 靜態編譯二進制檔
go build -trimpath -ldflags="-s -w" -o bin/gbootcms.exe .

# 啟動服務
.\bin\gbootcms.exe

# 或使用整合腳本自動編譯並運行
.\build-run.bat
```

### Linux / macOS (Bash)

```bash
# 交叉編譯或本機編譯
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/gbootcms .

# 設置權限並啟動
chmod +x bin/gbootcms
./bin/gbootcms
```

---

## 🌐 系統訪問入口

啟動完成後，預設監聽連接埠為 `8080`：

| 介面 | 訪問 URL | 預設憑據 |
| :--- | :--- | :--- |
| **前台首頁** | `http://localhost:8080/` | - |
| **管理後台** | `http://localhost:8080/admin` | 帳號：`admin`<br>密碼：`123456` |

> [!IMPORTANT]
> 首次登入系統後，請立即至「系統管理」->「系統用戶」修改預設管理員密碼。

---

## ⚙️ 核心開發與測試指令

```bash
# 執行單元測試
go test -v ./...

# 檢查資料庫連接
go run . --check-db
```

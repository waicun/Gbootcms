# GbootCMS 後台在線一鍵升級 (Self-Update) 實施計畫

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 為 GbootCMS 實裝面向使用者的後台「檢查更新」與「一鍵在線自升級」能力，自動自 GitHub Releases 獲取最新版本、校驗 SHA-256、原子原地替換二進制並觸發平滑重啟。

**Architecture:** 借鑑 `minio/selfupdate` 與 `restic self-update`，構建自更新業務服務層。服務層自 GitHub API 獲取最新 Release 資訊，比對當前版本，匹配當前平臺架構下載資產並解壓二進制，透過 `selfupdate.Apply` 完成跨平臺原地替換，最後觸發容器/進程平滑重啟；在控制器層集成 `DEMO_MODE` 保護阻斷演示環境破壞。

**Tech Stack:** Go 1.25, `github.com/minio/selfupdate`, GitHub Releases API, Gin, Layui.

---

### Task 1: 實裝自更新服務層 (`upgrade.go`) 與單元測試

**Files:**
- Create: `apps/admin/service/upgrade/upgrade.go`
- Create: `apps/admin/service/upgrade/upgrade_test.go`

- [ ] **Step 1: 編寫失敗測試案例**
測試 GitHub Release 格式解析、平臺資產匹配邏輯（`linux/amd64` 匹配 `gbootcms-linux.tar.gz`，`windows/amd64` 匹配 `gbootcms-windows.zip`）與版本比對邏輯。

- [ ] **Step 2: 執行測試確認失敗**
`go test ./apps/admin/service/upgrade/...`

- [ ] **Step 3: 編寫核心升級服務代碼**
  - `CheckUpdate(repo string, currentVersion string) (*ReleaseInfo, bool, error)`
  - `DownloadAndExtractBinary(assetUrl string) (io.Reader, error)`
  - `ApplyUpdate(binaryStream io.Reader, expectedSha256 string) error`

- [ ] **Step 4: 執行測試確認通過**
`go test ./apps/admin/service/upgrade/... -v`

- [ ] **Step 5: 提交代碼**
`git add apps/admin/service/upgrade/ go.mod go.sum && git commit -m "feat(upgrade): implement core selfupdate service and asset extractor"`

---

### Task 2: 實裝後台控制器 `UpgradeController.go`

**Files:**
- Create: `apps/admin/controller/system/UpgradeController.go`
- Modify: `apps/admin/controller/IndexController.go`

- [ ] **Step 1: 實現控制器方法**
  - `Check(c *gin.Context)`：返回最新版本、更新日誌、發布時間與是否有更新的 JSON。
  - `Apply(c *gin.Context)`：
    - 檢查 `common.IsDemoMode()`，若處於演示環境返回「演示環境已開啟保護模式，禁止在線升級！」。
    - 調用服務層執行下載、校驗與 `selfupdate.Apply`。
    - 替換成功後排程延遲 1 秒退出進程（以供 Docker / systemd 自動拉起新進程）。

- [ ] **Step 2: 驗證編譯與邏輯**
`go build -v .`

- [ ] **Step 3: 提交代碼**
`git add apps/admin/controller/system/UpgradeController.go && git commit -m "feat(upgrade): add UpgradeController with DemoMode guard"`

---

### Task 3: 註冊後台路由

**Files:**
- Modify: `apps/route/route.go`

- [ ] **Step 1: 在管理後台系統模組註冊路由**
  - `/admin/upgrade/check` -> `upgradeController.Check`
  - `/admin/upgrade/apply` -> `upgradeController.Apply`

- [ ] **Step 2: 驗證編譯**
`go test ./apps/route/...`

- [ ] **Step 3: 提交代碼**
`git add apps/route/route.go && git commit -m "feat(route): register admin upgrade routes"`

---

### Task 4: 後台視圖整合 Layui 檢查更新與一鍵升級彈窗

**Files:**
- Modify: `apps/admin/view/system/home.html`
- Modify: `apps/admin/view/system/server.html`

- [ ] **Step 1: 在 `home.html` 核心架構項添加「檢查更新」按鈕與互動彈窗**
  - 點擊「檢查更新」向 `/admin/upgrade/check` 發送 AJAX 請求。
  - 若有新版本，彈出 Layui 提示框展示版本與 Changelog，並提供「立即升級」按鈕。
  - 點擊「立即升級」發送 POST `/admin/upgrade/apply`，展示載入遮罩，成功後提示重載。

- [ ] **Step 2: 提交代碼**
`git add apps/admin/view/system/ && git commit -m "feat(view): integrate one-click upgrade dialog in admin dashboard"`

---

### Task 5: 全量編譯與端到端驗證

- [ ] **Step 1: 執行單元測試**
`go test ./...`

- [ ] **Step 2: 驗證 Windows 與 Linux 交叉編譯**
`go build -trimpath -o bin/test_build .`
`$env:GOOS='linux'; go build -trimpath -o bin/test_linux .`

- [ ] **Step 3: 驗證 DemoMode 攔截與安全閉環**
驗證當 `DEMO_MODE=1` 時，升級介面安全返回攔截錯誤。

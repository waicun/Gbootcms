package system

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"gbootcms/apps/admin/service/upgrade"
	"gbootcms/apps/common"
	"gbootcms/core/version"

	"github.com/gin-gonic/gin"
)

type UpgradeController struct {
	common.BaseController
}

// Check 檢查線上最新發行版本
func (ctrl *UpgradeController) Check(c *gin.Context) {
	res, err := upgrade.CheckForUpdate(upgrade.DefaultRepo, version.Version)
	if err != nil {
		slog.Error("檢查線上更新失敗", "error", err)
		ctrl.JSONFail(c, fmt.Sprintf("檢查更新失敗: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "獲取更新資訊成功",
		"data": res,
	})
}

// Apply 執行在線自更新
func (ctrl *UpgradeController) Apply(c *gin.Context) {
	if common.IsDemoMode() {
		ctrl.JSONFail(c, "演示環境已開啟保護模式，禁止在線升級！")
		return
	}

	checkRes, err := upgrade.CheckForUpdate(upgrade.DefaultRepo, version.Version)
	if err != nil {
		ctrl.JSONFail(c, fmt.Sprintf("獲取版本資訊失敗: %v", err))
		return
	}

	if !checkRes.HasUpdate {
		ctrl.JSONFail(c, "當前已是最新版本，無需升級")
		return
	}

	if checkRes.TargetAsset == nil {
		ctrl.JSONFail(c, fmt.Sprintf("未找到適配當前平臺 (%s/%s) 的發行資產", runtime.GOOS, runtime.GOARCH))
		return
	}

	slog.Info("開始下載新版本二進制資產", "url", checkRes.TargetAsset.BrowserDownloadURL, "size", checkRes.TargetAsset.Size)

	assetData, err := upgrade.DownloadAsset(checkRes.TargetAsset.BrowserDownloadURL)
	if err != nil {
		slog.Error("下載升級資產失敗", "error", err)
		ctrl.JSONFail(c, fmt.Sprintf("下載資產失敗: %v", err))
		return
	}

	var binaryBytes []byte
	if runtime.GOOS == "windows" {
		binaryBytes, err = upgrade.ExtractBinaryFromZip(assetData)
	} else {
		binaryBytes, err = upgrade.ExtractBinaryFromTarGz(assetData)
	}

	if err != nil {
		slog.Error("解壓二進制檔案失敗", "error", err)
		ctrl.JSONFail(c, fmt.Sprintf("解壓二進制失敗: %v", err))
		return
	}

	if err := upgrade.ApplySelfUpdate(binaryBytes); err != nil {
		slog.Error("執行二進制原地替換失敗", "error", err)
		ctrl.JSONFail(c, fmt.Sprintf("應用更新失敗: %v", err))
		return
	}

	ctrl.LogAction(c, fmt.Sprintf("在線升級系統至 %s 成功", checkRes.LatestVersion))

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "二進制原地替換成功！服務將在 2 秒後重啟生效，請稍候刷新頁面。",
	})

	// 排程優雅重啟：延時給予 HTTP 響應足夠時間發送至前端，隨後進程退出由 Docker/1Panel/systemd 自動重拉
	go func() {
		time.Sleep(1500 * time.Millisecond)
		slog.Info("升級完成，觸發進程重啟...")
		os.Exit(0)
	}()
}

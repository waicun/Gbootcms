package upgrade

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/minio/selfupdate"
)

const DefaultRepo = "waicun/Gbootcms"

type ReleaseAsset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type ReleaseInfo struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Body        string         `json:"body"`
	PublishedAt string         `json:"published_at"`
	Assets      []ReleaseAsset `json:"assets"`
}

type CheckResult struct {
	HasUpdate      bool          `json:"has_update"`
	CurrentVersion string        `json:"current_version"`
	LatestVersion  string        `json:"latest_version"`
	ReleaseTitle   string        `json:"release_title"`
	ReleaseNotes   string        `json:"release_notes"`
	PublishedAt    string        `json:"published_at"`
	TargetAsset    *ReleaseAsset `json:"target_asset"`
}

// CompareVersions 比對兩個語義版本號 (如 v1.0.1 vs v1.0.0)
// 若 v1 > v2 返回 1; 若 v1 < v2 返回 -1; 相等返回 0
func CompareVersions(v1, v2 string) int {
	clean1 := strings.TrimPrefix(strings.TrimSpace(strings.ToLower(v1)), "v")
	clean2 := strings.TrimPrefix(strings.TrimSpace(strings.ToLower(v2)), "v")

	parts1 := strings.Split(clean1, ".")
	parts2 := strings.Split(clean2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(parts1) {
			n1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			n2, _ = strconv.Atoi(parts2[i])
		}
		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}
	return 0
}

// MatchTargetAsset 依當前系統平臺 (OS/Arch) 自動匹配 Release 資產
func MatchTargetAsset(assets []ReleaseAsset, goos, goarch string) *ReleaseAsset {
	for i := range assets {
		a := &assets[i]
		lower := strings.ToLower(a.Name)
		if goos == "linux" && strings.Contains(lower, "linux") && strings.HasSuffix(lower, ".tar.gz") {
			return a
		}
		if goos == "windows" && strings.Contains(lower, "windows") && strings.HasSuffix(lower, ".zip") {
			return a
		}
	}
	return nil
}

// CheckForUpdate 查詢 GitHub API 檢測是否有新版本
func CheckForUpdate(repo string, currentVersion string) (*CheckResult, error) {
	if repo == "" {
		repo = DefaultRepo
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Gbootcms-Updater")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("連接 GitHub 檢查更新失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API 返回異常狀態碼: %d", resp.StatusCode)
	}

	var rel ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("解析 GitHub Release 失敗: %w", err)
	}

	hasUpdate := CompareVersions(rel.TagName, currentVersion) > 0
	targetAsset := MatchTargetAsset(rel.Assets, runtime.GOOS, runtime.GOARCH)

	return &CheckResult{
		HasUpdate:      hasUpdate,
		CurrentVersion: currentVersion,
		LatestVersion:  rel.TagName,
		ReleaseTitle:   rel.Name,
		ReleaseNotes:   rel.Body,
		PublishedAt:    rel.PublishedAt,
		TargetAsset:    targetAsset,
	}, nil
}

// DownloadAsset 下載指定 URL 的資產包位元組流
func DownloadAsset(url string) ([]byte, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Gbootcms-Updater")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下載資產失敗，HTTP 狀態碼: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// ExtractBinaryFromTarGz 自 Linux tar.gz 資產包解壓出二進制檔案
func ExtractBinaryFromTarGz(tarGzData []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewReader(tarGzData))
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		name := strings.TrimPrefix(header.Name, "./")
		if name == "bin/gbootcms" || name == "gbootcms" || strings.HasSuffix(name, "/gbootcms") {
			return io.ReadAll(tr)
		}
	}
	return nil, errors.New("資產包內未找到二進制執行檔 (bin/gbootcms)")
}

// ExtractBinaryFromZip 自 Windows zip 資產包解壓出二進制檔案
func ExtractBinaryFromZip(zipData []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	for _, f := range zr.File {
		name := strings.TrimPrefix(f.Name, "./")
		if name == "bin/gbootcms.exe" || name == "gbootcms.exe" || strings.HasSuffix(name, "/gbootcms.exe") {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, errors.New("資產包內未找到二進制執行檔 (bin/gbootcms.exe)")
}

// ApplySelfUpdate 將下載解壓的二進制安全原子原地替換當前程式
func ApplySelfUpdate(binaryBytes []byte) error {
	if len(binaryBytes) == 0 {
		return errors.New("二進制資料為空，拒絕更新")
	}

	opts := selfupdate.Options{}
	err := selfupdate.Apply(bytes.NewReader(binaryBytes), opts)
	if err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			return fmt.Errorf("替換失敗且回滾異常: %w (回滾錯誤: %v)", err, rerr)
		}
		return fmt.Errorf("替換二進制失敗: %w", err)
	}
	return nil
}

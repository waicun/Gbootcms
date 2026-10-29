# 🚢 生產部署指引 (Deployment Guide)

Gbootcms 為純 Go 靜態編譯架構，不依賴 PHP-FPM、動態 C 函式庫或外部 Web 伺服器，可直接以二進制檔或 Docker 容器獨立運作。

---

## 📦 發行包建構

### Linux amd64 打包

於專案根目錄執行發行打包腳本：

```powershell
# Windows PowerShell
.\build-linux.bat
```

產物規格：
- **發行檔案**：`dist/gbootcms-linux.tar.gz`
- **檔案大小**：約 12 MB（內含 UPX 壓縮執行檔、靜態資源、基礎模板）
- **編譯參數**：`-trimpath -ldflags="-s -w -buildid="`（剝離偵錯符號與路徑資訊）

---

## ⚙️ Systemd 守護行程配置

在 Linux (Ubuntu / Debian / RHEL) 環境中建立 `/etc/systemd/system/gbootcms.service`：

```ini
[Unit]
Description=Gbootcms Service
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/gbootcms
ExecStart=/opt/gbootcms/bin/gbootcms
Restart=always
RestartSec=5
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

### 行程管理指令

```bash
# 重新載入設定
sudo systemctl daemon-reload

# 啟動與自啟動
sudo systemctl enable --now gbootcms

# 檢查運行狀態
sudo systemctl status gbootcms
```

---

## 🐳 Docker 容器化部署

### 容器建構檔 (`Dockerfile`)

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o gbootcms .

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Taipei
WORKDIR /app
COPY --from=builder /app/gbootcms .
COPY --from=builder /app/config ./config
COPY --from=builder /app/template ./template
COPY --from=builder /app/static ./static
COPY --from=builder /app/apps/admin/view ./apps/admin/view

EXPOSE 8080
VOLUME ["/app/data", "/app/static/upload"]
CMD ["./gbootcms"]
```

---

## 🌐 Nginx 反向代理配置

若需要於前端附加 TLS 終結或既有網域路由：

```nginx
server {
    listen 80;
    server_name example.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name example.com;

    ssl_certificate /etc/nginx/ssl/example.com.crt;
    ssl_certificate_key /etc/nginx/ssl/example.com.key;

    client_max_body_size 50M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location ~* \.(jpg|jpeg|png|gif|webp|css|js|woff2)$ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        expires 30d;
        add_header Cache-Control "public, no-transform";
    }
}
```

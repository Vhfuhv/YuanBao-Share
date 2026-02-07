# 快速启动指南

## 最简单的启动方式（不使用爬虫）

如果你只想快速体验平台功能，不需要自动爬虫：

```bash
# 1. 安装Go依赖
go mod download

# 2. 直接运行
go run main.go

# 3. 访问
# 打开浏览器访问 http://localhost:18080
```

数据库会自动创建，无需任何配置！

## 完整启动方式（包含爬虫）

如果需要自动爬虫功能：

### 步骤1：安装Python依赖

```bash
# 创建虚拟环境
python -m venv venv

# 激活虚拟环境
# Windows:
venv\Scripts\activate
# Linux/Mac:
source venv/bin/activate

# 安装依赖
pip install -r python_test/requirements.txt
```

### 步骤2：配置爬虫

```bash
# 进入爬虫目录
cd python_test

# 复制配置模板
cp config.example.json config.json

# 编辑配置文件（填入你的Cookie）
# Windows: notepad config.json
# Linux/Mac: nano config.json
```

**获取Cookie的方法**：
1. 打开浏览器，登录百度贴吧
2. 按 F12 打开开发者工具
3. 切换到 Network 标签
4. 刷新页面，找到任意请求
5. 在请求头中找到 Cookie，复制完整内容
6. 粘贴到 config.json 的 cookies 数组中

### 步骤3：启动程序

```bash
# 返回项目根目录
cd ..

# 运行程序
go run main.go
```

程序会自动：
- 启动时立即执行一次爬虫
- 每30分钟执行方案1（单帖子爬虫）
- 每1小时执行方案2（首页爬虫）
- 每1小时清理过期的爬虫口令
- 每天0点清空所有数据

## 常见问题

### Q: 爬虫报错怎么办？

A: 如果爬虫配置有问题，程序仍会正常运行，只是没有自动采集功能。用户仍可以手动上传和获取口令。

### Q: 如何只使用用户上传功能？

A: 不配置爬虫即可。程序会检测到没有config.json，跳过爬虫功能。

### Q: 数据存储在哪里？

A: 所有数据存储在项目根目录的 `yuanbao.db` 文件中（SQLite数据库）。

### Q: 如何清空所有数据？

A: 删除 `yuanbao.db` 文件，重启程序会自动创建新的空数据库。

### Q: 如何修改端口？

A: 编辑 `main.go` 文件，修改最后一行的端口号：
```go
r.Run(":18080")  // 改为你想要的端口
```

## 部署到生产环境

### 编译

```bash
# 编译为可执行文件
go build -o yuanbao

# 跨平台编译（Linux）
GOOS=linux GOARCH=amd64 go build -o yuanbao-linux

# 优化编译（减小文件大小）
go build -ldflags="-s -w" -o yuanbao
```

### 后台运行

**Linux (使用 systemd)**:

创建服务文件 `/etc/systemd/system/yuanbao.service`:

```ini
[Unit]
Description=YuanBao Share Service
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/YuanBao-Share
ExecStart=/path/to/YuanBao-Share/yuanbao
Restart=always

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl enable yuanbao
sudo systemctl start yuanbao
sudo systemctl status yuanbao
```

**Windows (使用 NSSM)**:

1. 下载 NSSM: https://nssm.cc/download
2. 安装服务：
```cmd
nssm install YuanBao "E:\YuanBao-Share\yuanbao.exe"
nssm set YuanBao AppDirectory "E:\YuanBao-Share"
nssm start YuanBao
```

## 技术支持

如有问题，请查看：
- 主文档：README.md
- 爬虫文档：python_test/README.md

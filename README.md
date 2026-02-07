# 腾讯元宝口令红包分享平台 - Go 版本

一个使用 Go 语言实现的口令红包分享平台，用户可以上传自己的元宝口令，也可以随机获取他人的口令。

**🚀 新用户？查看 [快速启动指南](QUICKSTART.md) 快速上手！**

## 技术栈

- **后端**: Go 1.21 + Gin + GORM
- **数据库**: SQLite 3
- **爬虫**: Python 3.x + requests + BeautifulSoup4
- **前端**: 纯 HTML + CSS + JavaScript

## 功能特性

- ✅ 完全匿名，无需登录
- ✅ 上传口令
- ✅ 随机获取他人口令
- ✅ 一键复制口令
- ✅ 实时统计可用口令总数
- ✅ 每个口令最多被展示 3 次（符合元宝红包规则）
- ✅ 悲观锁机制，防止并发超发
- ✅ 自动爬虫系统，从百度贴吧自动采集口令
- ✅ 双来源优先级：优先展示用户上传的口令
- ✅ IP过滤：用户不会获取到自己上传的口令
- ✅ 定时清理：每小时清理过期爬虫口令，每天0点清空所有数据

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- Python 3.7+ （用于爬虫功能）

### 爬虫配置（可选）

如果需要使用自动爬虫功能：

1. 创建Python虚拟环境
```bash
python -m venv venv
```

2. 激活虚拟环境并安装依赖
```bash
# Windows
venv\Scripts\activate
pip install -r python_test\requirements.txt

# Linux/Mac
source venv/bin/activate
pip install -r python_test/requirements.txt
```

3. 配置爬虫
```bash
cd python_test
cp config.example.json config.json
# 编辑 config.json，填入你的百度贴吧Cookie
```

**注意**：如果不配置爬虫，程序仍可正常运行，只是没有自动采集功能。

### 运行步骤

1. 安装依赖
```bash
go mod download
```

2. 运行应用
```bash
go run main.go
```

或者编译后运行：
```bash
go build -o yuanbao.exe
yuanbao.exe
```

3. 访问应用

打开浏览器访问：http://localhost:18080

## 项目结构

```
YuanBao-Share/
├── main.go                      # 主程序入口
├── config/
│   └── database.go             # 数据库配置（SQLite）
├── models/
│   └── command.go              # 数据模型
├── repositories/
│   └── command_repository.go   # 数据访问层
├── services/
│   ├── command_service.go      # 业务逻辑层
│   └── crawler_service.go      # 爬虫服务
├── controllers/
│   └── command_controller.go   # 控制器层
├── middleware/
│   └── rate_limiter.go         # 限流中间件
├── static/                      # 前端静态文件
│   ├── index.html
│   ├── style.css
│   └── app.js
├── python_test/                 # Python爬虫脚本
│   ├── tieba_crawler.py        # 单帖子爬虫
│   ├── tieba_crawler_v2.py     # 首页爬虫
│   ├── config.json             # 爬虫配置（需自行创建）
│   ├── config.example.json     # 配置模板
│   ├── requirements.txt        # Python依赖
│   └── README.md               # 爬虫说明
├── yuanbao.db                   # SQLite数据库文件（自动创建）
├── go.mod                       # Go 模块文件
└── README.md                    # 项目说明
```

## API 接口

### 上传口令
```
POST /api/commands
Content-Type: application/json

{
  "content": "你的口令内容"
}
```

### 随机获取口令
```
GET /api/commands/random
```

### 获取口令总数
```
GET /api/commands/count
```

### 报告无效口令
```
POST /api/commands/report
Content-Type: application/json

{
  "content": "无效的口令内容"
}
```

## 并发控制

项目使用悲观锁机制防止并发超发：

```go
// 使用 SELECT ... FOR UPDATE 锁定行
config.DB.Clauses(clause.Locking{Strength: "UPDATE"}).
    Where("display_count < ?", 3).
    Order("RAND()").
    First(&command)
```

**工作原理：**
1. 事务开始时锁定一行数据
2. 其他并发请求等待
3. 更新 display_count + 1
4. 提交事务，释放锁

这样可以保证同一个口令不会被并发获取超过 3 次。

## 爬虫系统

项目集成了自动爬虫系统，可从百度贴吧自动采集口令：

### 工作原理

1. **方案1（单帖子爬虫）**：每30分钟执行一次，爬取指定帖子的最新20分钟内的口令
2. **方案2（首页爬虫）**：每1小时执行一次，爬取元宝吧首页前10个帖子的最新口令
3. **自动清理**：每1小时清理1小时前的爬虫口令，每天0点清空所有数据

### 口令优先级

- **用户上传的口令**：优先展示，且用户不会获取到自己上传的口令（通过IP过滤）
- **爬虫采集的口令**：作为备用，当没有用户口令时展示

### 配置说明

详见 `python_test/README.md`


## 部署建议

### 生产环境

1. 使用环境变量配置数据库连接
2. 启用 Gin 的 Release 模式
3. 配置 HTTPS
4. 添加日志监控
5. 使用 systemd 或 supervisor 管理进程

### 编译优化

```bash
# 编译优化版本
go build -ldflags="-s -w" -o yuanbao

# 跨平台编译
GOOS=linux GOARCH=amd64 go build -o yuanbao-linux
```

## 许可证

MIT License

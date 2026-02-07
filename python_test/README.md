# 百度贴吧爬虫测试脚本（基于时间的增量爬取）

## 功能说明

这是一个用于测试百度贴吧爬虫逻辑的 Python 脚本，支持基于时间的智能增量爬取。

### 核心逻辑

1. **定位最后一页**：获取帖子的最后一页页码
2. **倒序爬取楼层**：从最后一页的最后一楼开始，逐楼检查发帖时间
3. **时间过滤**：如果楼层发帖时间与脚本执行时间相差 <= 20分钟，加入结果
4. **内容验证**：验证口令长度和格式
5. **跨页爬取**：如果该页所有楼层都满足时间条件，继续爬取上一页
6. **智能停止**：遇到超过20分钟的楼层，立即停止爬取

### 主要特性

- ✅ 基于时间的增量爬取（不依赖页数变化）
- ✅ 支持多种时间格式解析（"X分钟前"、"X小时前"、"18:30"、"2026-02-06 18:30"）
- ✅ 倒序遍历（从最新到最旧）
- ✅ 内容验证（长度、链接过滤）
- ✅ 随机间隔（10-20秒）
- ✅ 状态记录

## 安装依赖

```bash
cd python_test
pip install -r requirements.txt
```

或者手动安装：

```bash
pip install requests beautifulsoup4 lxml
```

## 使用方法

### 1. 配置爬虫

**重要**：首次使用需要创建配置文件

```bash
# 复制配置模板
cp config.example.json config.json

# 编辑配置文件
# Windows: notepad config.json
# Linux/Mac: nano config.json
```

配置文件说明：
```json
{
  "cookies": [
    "你的第一个Cookie",
    "你的第二个Cookie"
  ],
  "tieba_url": "https://tieba.baidu.com/p/10449473531",
  "tieba_homepage": "https://tieba.baidu.com/f?ie=utf-8&kw=%E8%85%BE%E8%AE%AF%E5%85%83%E5%AE%9D&fr=search",
  "time_threshold_minutes": 20
}
```

**获取Cookie方法**：
1. 打开浏览器，登录百度贴吧
2. 按 F12 打开开发者工具
3. 切换到 Network 标签
4. 刷新页面，找到任意请求
5. 在请求头中找到 Cookie，复制完整内容
6. 粘贴到 config.json 的 cookies 数组中

**多Cookie说明**：
- 支持配置多个Cookie，脚本会随机选择使用
- 可以降低被检测为爬虫的风险
- 至少配置1个Cookie即可运行

### 2. 运行脚本

```bash
python tieba_crawler.py
```

或使用虚拟环境：

```bash
E:\YuanBao-Share\venv\Scripts\python.exe tieba_crawler.py
```

### 3. 查看结果

- 口令保存在 `commands.json` 和 `commands_v2.json`（JSON格式）
- 状态保存在 `crawler_state.json` 和 `crawler_state_v2.json`

**注意**：这些文件会被Go程序自动读取并导入数据库，无需手动处理。

## 工作流程示例

### 场景 1：首次运行（帖子最后一页有新内容）

```
爬取开始时间: 2026-02-06 19:00:00
帖子最后一页: 1003

正在爬取第 1003 页...
  ✓ 加入: 18:55 - 口令内容1...
  ✓ 加入: 18:50 - 口令内容2...
  ✓ 加入: 18:45 - 口令内容3...
  ⊗ 超时: 18:30 (相差 30.0 分钟)

爬取完成！共获取 3 个有效口令
```

### 场景 2：跨页爬取（最后一页所有楼层都在20分钟内）

```
爬取开始时间: 2026-02-06 19:00:00
帖子最后一页: 1003

正在爬取第 1003 页...
  ✓ 加入: 18:55 - 口令内容1...
  ✓ 加入: 18:50 - 口令内容2...
  ✓ 加入: 18:45 - 口令内容3...
第 1003 页所有楼层都在时间范围内，继续爬取上一页

等待 15.3 秒...

正在爬取第 1002 页...
  ✓ 加入: 18:42 - 口令内容4...
  ⊗ 超时: 18:35 (相差 25.0 分钟)

爬取完成！共获取 4 个有效口令
```

### 场景 3：无新内容

```
爬取开始时间: 2026-02-06 19:00:00
帖子最后一页: 1003

正在爬取第 1003 页...
  ⊗ 超时: 17:30 (相差 90.0 分钟)

爬取完成！共获取 0 个有效口令
```

## 时间格式支持

脚本支持解析以下时间格式：

| 格式 | 示例 | 说明 |
|------|------|------|
| 完整日期时间 | `2026-02-06 18:30` | 精确到分钟 |
| 今天时间 | `今天18:30` 或 `18:30` | 当天的时间 |
| 相对时间（分钟） | `5分钟前` | 相对于当前时间 |
| 相对时间（小时） | `2小时前` | 相对于当前时间 |

## 验证规则

脚本会过滤以下内容：

- ❌ 长度小于 10 字符或大于 500 字符
- ❌ 包含 http:// 或 https:// 链接

## 输出格式

### commands.json / commands_v2.json

```json
{
  "crawl_time": "2026-02-06 19:00:00",
  "source": "single_thread",
  "thread_url": "https://tieba.baidu.com/p/10449473531",
  "commands": [
    {
      "content": "口令内容1",
      "post_time": "18:55"
    },
    {
      "content": "口令内容2",
      "post_time": "18:50"
    }
  ]
}
```

### crawler_state.json

```json
{
  "last_crawl_time": "2026-02-06 19:00:00",
  "commands_count": 3,
  "last_page": 1003
}
```

## 定时任务配置

**注意**：如果使用Go程序，定时任务已自动集成，无需手动配置。

以下配置仅用于独立运行Python脚本的场景：

### Windows 任务计划程序

1. 打开"任务计划程序"
2. 创建基本任务
3. 触发器：每 20 分钟执行一次
4. 操作：启动程序
   - 程序：`E:\YuanBao-Share\venv\Scripts\python.exe`
   - 参数：`E:\YuanBao-Share\python_test\tieba_crawler.py`
   - 起始于：`E:\YuanBao-Share\python_test`

### Linux Cron

```bash
# 每20分钟执行一次
*/20 * * * * cd /path/to/python_test && /path/to/python tieba_crawler.py
```

## 注意事项

1. **Cookie 时效性**：Cookie 会过期，如果爬取失败，需要更新 config.json 中的 cookies
2. **配置文件安全**：config.json 包含敏感信息，已添加到 .gitignore，不会被提交到Git
3. **时间阈值**：默认 20 分钟，可在 config.json 中调整 `time_threshold_minutes`
4. **随机间隔**：跨页爬取时，每页间隔 10-20 秒，避免被反爬虫检测
5. **时间解析**：如果遇到新的时间格式，需要在 `parse_post_time()` 函数中添加支持
6. **首次运行**：必须先创建 config.json 文件，否则脚本会报错

## 优势对比

### 旧版本（基于页数）
- ❌ 依赖页数变化，不够灵活
- ❌ 可能爬取大量旧内容
- ❌ 无法精确控制时间范围

### 新版本（基于时间）
- ✅ 精确控制时间范围（20分钟内）
- ✅ 自动跨页爬取，直到遇到超时楼层
- ✅ 避免重复爬取旧内容
- ✅ 适合定时任务，每次只爬取新增内容

## 故障排查

### 问题1：无法获取时间

**现象**：输出 "⊗ 无法解析时间"

**解决**：
1. 检查贴吧页面结构是否变化
2. 在 `crawl_page_with_time()` 函数中添加调试输出
3. 查看 `data-field` 或 `tail-info` 元素是否存在

### 问题2：Cookie 失效

**现象**：爬取失败，返回登录页面或403错误

**解决**：
1. 打开浏览器，登录百度贴吧
2. 按 F12 打开开发者工具
3. 切换到 Network 标签
4. 刷新页面，找到请求头中的 Cookie
5. 复制完整 Cookie 到 config.json 的 cookies 数组中
6. 可以配置多个Cookie以提高稳定性

### 问题3：找不到config.json

**现象**：FileNotFoundError: config.json

**解决**：
```bash
cd python_test
cp config.example.json config.json
# 然后编辑 config.json 填入你的Cookie
```

### 问题3：时间解析错误

**现象**：时间差计算不正确

**解决**：
1. 检查系统时间是否正确
2. 在 `parse_post_time()` 函数中添加调试输出
3. 确认时间格式是否在支持列表中

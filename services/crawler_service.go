package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"yuanbao/repositories"
)

// CrawlerResult 爬虫结果结构
type CrawlerResult struct {
	CrawlTime  string    `json:"crawl_time"`
	Source     string    `json:"source"`
	ThreadURL  string    `json:"thread_url,omitempty"`
	Commands   []Command `json:"commands,omitempty"`
	Threads    []Thread  `json:"threads,omitempty"`
}

type Command struct {
	Content  string `json:"content"`
	PostTime string `json:"post_time"`
}

type Thread struct {
	Title    string    `json:"title"`
	URL      string    `json:"url"`
	Commands []Command `json:"commands"`
}

// RunCrawler 执行爬虫任务
func RunCrawler() error {
	return RunCrawlerV1()
}

// RunCrawlerV1 执行第一套方案（单个帖子）
func RunCrawlerV1() error {
	log.Println("========================================")
	log.Println("开始执行爬虫任务（方案1：单个帖子）")
	log.Println("========================================")

	// 获取项目根目录
	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}

	pythonDir := filepath.Join(rootDir, "python_test")

	// 跨平台Python路径检测
	var pythonExe string
	if _, err := os.Stat(filepath.Join(rootDir, "venv", "Scripts", "python.exe")); err == nil {
		// Windows
		pythonExe = filepath.Join(rootDir, "venv", "Scripts", "python.exe")
	} else if _, err := os.Stat(filepath.Join(rootDir, "venv", "bin", "python")); err == nil {
		// Linux/Mac
		pythonExe = filepath.Join(rootDir, "venv", "bin", "python")
	} else {
		// 使用系统Python
		pythonExe = "python"
	}

	// 执行第一套方案
	log.Println("执行脚本: tieba_crawler.py")
	err = runPythonScript(pythonExe, pythonDir, "tieba_crawler.py")
	if err != nil {
		log.Printf("脚本执行失败: %v", err)
		return err
	}

	// 读取并处理结果
	jsonFile := filepath.Join(pythonDir, "commands.json")
	err = processJSONFile(jsonFile, "方案1")
	if err != nil {
		log.Printf("处理结果失败: %v", err)
		return err
	}

	log.Println("========================================")
	return nil
}

// RunCrawlerV2 执行第二套方案（元宝吧首页）
func RunCrawlerV2() error {
	log.Println("========================================")
	log.Println("开始执行爬虫任务（方案2：元宝吧首页）")
	log.Println("========================================")

	// 获取项目根目录
	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}

	pythonDir := filepath.Join(rootDir, "python_test")

	// 跨平台Python路径检测
	var pythonExe string
	if _, err := os.Stat(filepath.Join(rootDir, "venv", "Scripts", "python.exe")); err == nil {
		// Windows
		pythonExe = filepath.Join(rootDir, "venv", "Scripts", "python.exe")
	} else if _, err := os.Stat(filepath.Join(rootDir, "venv", "bin", "python")); err == nil {
		// Linux/Mac
		pythonExe = filepath.Join(rootDir, "venv", "bin", "python")
	} else {
		// 使用系统Python
		pythonExe = "python"
	}

	// 执行第二套方案
	log.Println("执行脚本: tieba_crawler_v2.py")
	err = runPythonScript(pythonExe, pythonDir, "tieba_crawler_v2.py")
	if err != nil {
		log.Printf("脚本执行失败: %v", err)
		return err
	}

	// 读取并处理结果
	jsonFile := filepath.Join(pythonDir, "commands_v2.json")
	err = processJSONFile(jsonFile, "方案2")
	if err != nil {
		log.Printf("处理结果失败: %v", err)
		return err
	}

	log.Println("========================================")
	return nil
}

// runPythonScript 执行Python脚本
func runPythonScript(pythonExe, workDir, scriptName string) error {
	scriptPath := filepath.Join(workDir, scriptName)

	log.Printf("执行脚本: %s", scriptName)
	log.Printf("Python路径: %s", pythonExe)
	log.Printf("工作目录: %s", workDir)

	cmd := exec.Command(pythonExe, scriptPath)
	cmd.Dir = workDir

	// 捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("脚本执行失败: %v", err)
		log.Printf("输出: %s", string(output))
		return err
	}

	log.Printf("脚本执行成功")
	// 可选：打印部分输出
	// log.Printf("输出: %s", string(output))

	return nil
}

// processJSONFile 处理JSON文件并保存到数据库
func processJSONFile(jsonFile, source string) error {
	log.Printf("读取文件: %s", jsonFile)

	// 检查文件是否存在
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", jsonFile)
	}

	// 读取JSON文件
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	// 解析JSON
	var result CrawlerResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("解析JSON失败: %v", err)
	}

	log.Printf("爬取时间: %s", result.CrawlTime)
	log.Printf("数据源: %s", result.Source)

	// 统计数据
	totalCommands := 0
	successCount := 0
	duplicateCount := 0
	errorCount := 0

	// 处理命令
	if result.Source == "single_thread" {
		// 第一套方案：直接处理commands
		totalCommands = len(result.Commands)
		log.Printf("共获取 %d 个口令", totalCommands)

		for _, cmd := range result.Commands {
			_, err := SaveCrawlerCommand(cmd.Content)
			if err != nil {
				if err.Error() == "该口令已存在" {
					duplicateCount++
				} else {
					errorCount++
					// 安全截断内容
					preview := cmd.Content
					if len(preview) > 30 {
						preview = preview[:30]
					}
					log.Printf("保存失败: %s - %v", preview, err)
				}
			} else {
				successCount++
			}
		}
	} else if result.Source == "homepage_threads" {
		// 第二套方案：遍历threads
		for _, thread := range result.Threads {
			totalCommands += len(thread.Commands)
			for _, cmd := range thread.Commands {
				_, err := SaveCrawlerCommand(cmd.Content)
				if err != nil {
					if err.Error() == "该口令已存在" {
						duplicateCount++
					} else {
						errorCount++
						// 安全截断内容
						preview := cmd.Content
						if len(preview) > 30 {
							preview = preview[:30]
						}
						log.Printf("保存失败: %s - %v", preview, err)
					}
				} else {
					successCount++
				}
			}
		}
		log.Printf("共爬取 %d 个帖子，获取 %d 个口令", len(result.Threads), totalCommands)
	}

	// 输出统计
	log.Printf("----------------------------------------")
	log.Printf("总口令数: %d", totalCommands)
	log.Printf("成功保存: %d", successCount)
	log.Printf("重复跳过: %d", duplicateCount)
	log.Printf("保存失败: %d", errorCount)
	log.Printf("----------------------------------------")

	return nil
}

// StartCrawlerScheduler 启动爬虫定时任务
func StartCrawlerScheduler() {
	log.Println("========================================")
	log.Println("启动定时任务系统")
	log.Println("========================================")
	log.Println("- 方案1（单个帖子）：启动时立即执行，之后每30分钟执行")
	log.Println("- 方案2（元宝吧首页）：启动时立即执行，之后每1小时执行")
	log.Println("- 清理爬虫token：每1小时执行")
	log.Println("- 清空所有数据：每天0点执行")
	log.Println("========================================")

	// 1. 启动时立即执行两个爬虫方案
	go func() {
		time.Sleep(5 * time.Second) // 等待服务器启动完成
		log.Println("\n[启动任务] 执行方案1...")
		if err := RunCrawlerV1(); err != nil {
			log.Printf("方案1执行失败: %v", err)
		}

		time.Sleep(10 * time.Second) // 两个方案间隔10秒

		log.Println("\n[启动任务] 执行方案2...")
		if err := RunCrawlerV2(); err != nil {
			log.Printf("方案2执行失败: %v", err)
		}
	}()

	// 2. 方案1定时任务：每30分钟执行
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("\n[定时任务] 执行方案1...")
			if err := RunCrawlerV1(); err != nil {
				log.Printf("方案1执行失败: %v", err)
			}
		}
	}()

	// 3. 方案2定时任务：每1小时执行
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("\n[定时任务] 执行方案2...")
			if err := RunCrawlerV2(); err != nil {
				log.Printf("方案2执行失败: %v", err)
			}
		}
	}()

	// 4. 清理爬虫token：每1小时执行
	StartCleanupScheduler()

	// 5. 每天0点清空所有数据
	StartDailyCleanupScheduler()
}

// StartCleanupScheduler 启动清理定时任务（每小时清理爬虫token）
func StartCleanupScheduler() {
	log.Println("启动清理任务：每1小时清理1小时前的爬虫token")

	// 立即执行一次
	go func() {
		time.Sleep(15 * time.Second) // 等待服务器启动完成
		CleanOldCrawlerCommands()
	}()

	// 定时执行清理
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			CleanOldCrawlerCommands()
		}
	}()
}

// StartDailyCleanupScheduler 启动每日清空任务（每天0点清空所有数据）
func StartDailyCleanupScheduler() {
	log.Println("启动每日清空任务：每天0点清空所有token")

	go func() {
		for {
			// 计算到下一个0点的时间
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			duration := next.Sub(now)

			log.Printf("下次清空时间: %s (还有 %.1f 小时)", next.Format("2006-01-02 15:04:05"), duration.Hours())

			// 等待到0点
			time.Sleep(duration)

			// 执行清空
			log.Println("========================================")
			log.Println("执行每日清空任务（0点）")
			log.Println("========================================")

			count, err := repositories.CleanAllCommands()
			if err != nil {
				log.Printf("清空失败: %v", err)
			} else {
				log.Printf("成功清空 %d 条token（新的一天开始）", count)
			}
			log.Println("========================================")
		}
	}()
}

// CleanOldCrawlerCommands 清理1小时前的爬虫口令
func CleanOldCrawlerCommands() {
	log.Println("========================================")
	log.Println("开始清理旧的爬虫口令")
	log.Println("========================================")

	count, err := repositories.CleanOldCrawlerCommands()
	if err != nil {
		log.Printf("清理失败: %v", err)
		return
	}

	log.Printf("成功清理 %d 条1小时前的爬虫口令", count)
	log.Println("========================================")
}

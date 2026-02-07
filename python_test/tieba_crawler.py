"""
百度贴吧爬虫测试脚本（基于时间的增量爬取）
功能：从最后一页开始，爬取20分钟内的新楼层
逻辑：
1. 定位到最后一页
2. 从最后一页倒数第一楼开始，逐楼检查发帖时间
3. 如果时间差 <= 20分钟，加入结果
4. 如果该页所有楼层都满足，继续爬取倒数第二页
5. 遇到超过20分钟的楼层，停止爬取
"""

import requests
from bs4 import BeautifulSoup
import time
import re
import random
import json
import os
import sys
from datetime import datetime, timedelta

# 加载配置文件
def load_config():
    config_path = os.path.join(os.path.dirname(__file__), 'config.json')
    if not os.path.exists(config_path):
        print("=" * 60)
        print("错误：找不到配置文件 config.json")
        print("=" * 60)
        print("请按照以下步骤创建配置文件：")
        print("1. 复制配置模板：")
        print("   cp config.example.json config.json")
        print("")
        print("2. 编辑 config.json，填入你的百度贴吧Cookie")
        print("")
        print("3. 获取Cookie的方法：")
        print("   - 打开浏览器，登录百度贴吧")
        print("   - 按F12打开开发者工具")
        print("   - 切换到Network标签")
        print("   - 刷新页面，找到请求头中的Cookie")
        print("   - 复制完整Cookie到config.json")
        print("=" * 60)
        sys.exit(1)

    try:
        with open(config_path, 'r', encoding='utf-8') as f:
            return json.load(f)
    except json.JSONDecodeError as e:
        print(f"配置文件格式错误: {e}")
        print("请检查 config.json 是否为有效的JSON格式")
        sys.exit(1)

# 加载配置
config = load_config()
TIEBA_URL = config['tieba_url']
COOKIES = config['cookies']  # Cookie列表
TIME_THRESHOLD_MINUTES = config['time_threshold_minutes']

OUTPUT_FILE = "commands.json"
STATE_FILE = "crawler_state.json"


def get_random_cookie():
    """随机选择一个Cookie"""
    return random.choice(COOKIES)


def get_last_page_number(url):
    """获取帖子的最后一页页码"""
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
        "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
        "Accept-Encoding": "gzip, deflate, br",
        "Referer": "https://tieba.baidu.com/",
        "Connection": "keep-alive",
        "Upgrade-Insecure-Requests": "1",
        "Cookie": get_random_cookie()
    }

    try:
        response = requests.get(url, headers=headers, timeout=10)
        response.raise_for_status()

        soup = BeautifulSoup(response.text, 'html.parser')

        # 查找分页信息
        pager = soup.find('li', class_='l_pager')
        if not pager:
            return 1

        # 查找最后一页的链接
        last_page = 1
        for link in pager.find_all('a'):
            href = link.get('href', '')
            if 'pn=' in href:
                match = re.search(r'pn=(\d+)', href)
                if match:
                    page_num = int(match.group(1))
                    if page_num > last_page:
                        last_page = page_num

        return last_page

    except Exception as e:
        print(f"获取最后页码失败: {e}")
        return 1


def parse_post_time(time_str):
    """
    解析贴吧时间字符串
    支持格式：
    - "2026-02-06 18:30"
    - "今天18:30"
    - "18:30"
    - "1分钟前"
    - "1小时前"
    """
    now = datetime.now()

    try:
        # 格式1: "2026-02-06 18:30"
        if re.match(r'\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}', time_str):
            return datetime.strptime(time_str, "%Y-%m-%d %H:%M")

        # 格式2: "今天18:30" 或 "18:30"
        time_match = re.search(r'(\d{1,2}):(\d{2})', time_str)
        if time_match:
            hour = int(time_match.group(1))
            minute = int(time_match.group(2))
            post_time = now.replace(hour=hour, minute=minute, second=0, microsecond=0)

            # 如果时间在未来，说明是昨天的
            if post_time > now:
                post_time -= timedelta(days=1)

            return post_time

        # 格式3: "X分钟前"
        minute_match = re.search(r'(\d+)\s*分钟前', time_str)
        if minute_match:
            minutes = int(minute_match.group(1))
            return now - timedelta(minutes=minutes)

        # 格式4: "X小时前"
        hour_match = re.search(r'(\d+)\s*小时前', time_str)
        if hour_match:
            hours = int(hour_match.group(1))
            return now - timedelta(hours=hours)

        # 无法解析，返回None
        return None

    except Exception as e:
        print(f"解析时间失败: {time_str}, 错误: {e}")
        return None


def crawl_page_with_time(url, page_num):
    """
    爬取指定页的所有楼层，返回 (内容, 时间) 列表
    返回格式: [(content, post_time), ...]
    """
    # 构建URL
    if page_num > 1:
        page_url = f"{url}?pn={page_num}"
    else:
        page_url = url

    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
        "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
        "Accept-Encoding": "gzip, deflate, br",
        "Referer": "https://tieba.baidu.com/",
        "Connection": "keep-alive",
        "Upgrade-Insecure-Requests": "1",
        "Cookie": get_random_cookie()
    }

    try:
        response = requests.get(page_url, headers=headers, timeout=10)
        response.raise_for_status()

        soup = BeautifulSoup(response.text, 'html.parser')

        # 查找所有楼层
        posts = soup.find_all('div', class_='l_post')

        results = []
        for post in posts:
            try:
                # 解析 data-field JSON
                data_field = post.get('data-field', '{}')
                field_data = json.loads(data_field)

                # 获取楼层内容
                content_div = post.find('div', class_='d_post_content')
                if not content_div:
                    continue

                content = content_div.get_text(strip=True)

                # 获取发帖时间
                # 方法1: 从 data-field 中获取
                post_time_str = field_data.get('content', {}).get('date', '')

                # 方法2: 从页面元素中获取
                if not post_time_str:
                    # 查找 post-tail-wrap 容器
                    tail_wrap = post.find('div', class_='post-tail-wrap')
                    if tail_wrap:
                        # 获取所有 tail-info span
                        tail_infos = tail_wrap.find_all('span', class_='tail-info')
                        # 遍历找到时间格式的 span
                        for span in tail_infos:
                            text = span.get_text(strip=True)
                            # 检查是否包含时间格式（包含冒号或"前"字）
                            if ':' in text or '前' in text or '-' in text:
                                post_time_str = text
                                break

                # 解析时间
                post_time = parse_post_time(post_time_str) if post_time_str else None

                results.append({
                    'content': content,
                    'time': post_time,
                    'time_str': post_time_str
                })

            except Exception as e:
                print(f"解析楼层失败: {e}")
                continue

        return results

    except Exception as e:
        print(f"爬取第 {page_num} 页失败: {e}")
        return []


def is_valid_command(content):
    """验证口令是否有效"""
    # 长度检查
    if len(content) < 10 or len(content) > 500:
        return False

    # 不允许包含链接
    if 'http://' in content.lower() or 'https://' in content.lower():
        return False

    return True


def is_within_time_threshold(post_time, current_time, threshold_minutes):
    """判断发帖时间是否在阈值内"""
    if post_time is None:
        return False

    time_diff = current_time - post_time
    return time_diff.total_seconds() <= threshold_minutes * 60


def save_to_file(commands, filename):
    """保存口令到JSON文件"""
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # 构建JSON数据
    output_data = {
        "crawl_time": timestamp,
        "source": "single_thread",
        "thread_url": TIEBA_URL,
        "commands": []
    }

    for item in commands:
        output_data["commands"].append({
            "content": item['content'],
            "post_time": item['time_str']
        })

    # 保存为JSON
    try:
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(output_data, f, ensure_ascii=False, indent=2)
        print(f"已保存 {len(commands)} 个口令到 {filename}")
    except Exception as e:
        print(f"保存文件失败: {e}")


def load_state():
    """加载上次爬取的状态"""
    if os.path.exists(STATE_FILE):
        try:
            with open(STATE_FILE, 'r', encoding='utf-8') as f:
                return json.load(f)
        except Exception as e:
            print(f"加载状态文件失败: {e}")
    return {}


def save_state(state):
    """保存爬取状态"""
    try:
        with open(STATE_FILE, 'w', encoding='utf-8') as f:
            json.dump(state, f, ensure_ascii=False, indent=2)
        print(f"状态已保存到 {STATE_FILE}")
    except Exception as e:
        print(f"保存状态失败: {e}")


def main():
    print("=" * 60)
    print("百度贴吧爬虫（基于时间的增量爬取）")
    print("=" * 60)
    print(f"目标帖子: {TIEBA_URL}")
    print(f"输出文件: {OUTPUT_FILE}")
    print(f"时间阈值: {TIME_THRESHOLD_MINUTES} 分钟")
    print("=" * 60)

    # 记录爬取开始时间
    crawl_start_time = datetime.now()
    print(f"\n爬取开始时间: {crawl_start_time.strftime('%Y-%m-%d %H:%M:%S')}")

    # 获取最后一页页码
    last_page = get_last_page_number(TIEBA_URL)
    print(f"帖子最后一页: {last_page}")

    all_commands = []
    should_continue = True
    current_page = last_page

    # 从最后一页开始，向前爬取
    while should_continue and current_page >= 1:
        print(f"\n正在爬取第 {current_page} 页...")

        posts = crawl_page_with_time(TIEBA_URL, current_page)

        if not posts:
            print(f"第 {current_page} 页无数据，停止爬取")
            break

        page_all_valid = True  # 该页是否所有楼层都满足条件

        # 倒序遍历（从最后一楼到第一楼）
        for post in reversed(posts):
            content = post['content']
            post_time = post['time']
            time_str = post['time_str']

            # 检查时间是否在阈值内
            if is_within_time_threshold(post_time, crawl_start_time, TIME_THRESHOLD_MINUTES):
                # 验证内容
                if is_valid_command(content):
                    all_commands.append({
                        'content': content,
                        'time_str': time_str
                    })
                    print(f"  [OK] 加入: {time_str} - {content[:30]}...")
                else:
                    print(f"  [X] 内容不符合要求: {content[:30]}...")
            else:
                # 遇到超时的楼层
                if post_time:
                    time_diff = (crawl_start_time - post_time).total_seconds() / 60
                    print(f"  [SKIP] 超时: {time_str} (相差 {time_diff:.1f} 分钟)")
                else:
                    print(f"  [SKIP] 无法解析时间: {time_str}")

                page_all_valid = False
                should_continue = False
                break

        # 如果该页所有楼层都满足条件，继续爬取上一页
        if page_all_valid and current_page > 1:
            print(f"第 {current_page} 页所有楼层都在时间范围内，继续爬取上一页")
            current_page -= 1

            # 随机间隔
            wait_time = 10 + random.random() * 10
            print(f"等待 {wait_time:.1f} 秒...")
            time.sleep(wait_time)
        else:
            should_continue = False

    print("\n" + "=" * 60)
    print(f"爬取完成！共获取 {len(all_commands)} 个有效口令")
    print("=" * 60)

    # 保存到文件
    if all_commands:
        save_to_file(all_commands, OUTPUT_FILE)
        print(f"\n结果已保存到: {OUTPUT_FILE}")
    else:
        print("\n未获取到任何口令")

    # 更新状态
    state = load_state()
    state['last_crawl_time'] = crawl_start_time.strftime("%Y-%m-%d %H:%M:%S")
    state['commands_count'] = len(all_commands)
    state['last_page'] = last_page
    save_state(state)


if __name__ == "__main__":
    main()

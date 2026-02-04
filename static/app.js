const API_BASE = '/api/commands';

// 页面加载时获取统计数据
document.addEventListener('DOMContentLoaded', () => {
    loadStats();
});

// 加载统计数据
async function loadStats() {
    try {
        const response = await fetch(`${API_BASE}/count`);
        const data = await response.json();
        document.getElementById('totalCount').textContent = data.count;
    } catch (error) {
        console.error('加载统计数据失败:', error);
    }
}

// 上传口令
document.getElementById('uploadBtn').addEventListener('click', async () => {
    const input = document.getElementById('commandInput');
    const content = input.value.trim();
    const messageEl = document.getElementById('uploadMessage');
    const btn = document.getElementById('uploadBtn');

    if (!content) {
        showMessage(messageEl, '请输入口令内容', 'error');
        return;
    }

    btn.disabled = true;
    btn.textContent = '上传中...';

    try {
        const response = await fetch(API_BASE, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ content })
        });

        const data = await response.json();

        if (data.success) {
            showMessage(messageEl, '✅ 口令上传成功！现在可以去获取他人的口令了', 'success');
            input.value = '';
            loadStats();
        } else {
            showMessage(messageEl, data.message || '上传失败', 'error');
        }
    } catch (error) {
        showMessage(messageEl, '❌ 网络错误，请稍后重试', 'error');
        console.error('上传失败:', error);
    } finally {
        btn.disabled = false;
        btn.textContent = '上传口令';
    }
});

// 获取随机口令
document.getElementById('getBtn').addEventListener('click', async () => {
    const resultEl = document.getElementById('commandResult');
    const messageEl = document.getElementById('getMessage');
    const btn = document.getElementById('getBtn');

    btn.disabled = true;
    btn.textContent = '获取中...';
    resultEl.classList.remove('show');
    messageEl.style.display = 'none';

    try {
        const response = await fetch(`${API_BASE}/random`);
        const data = await response.json();

        if (data.success) {
            showCommand(resultEl, data.content);
        } else {
            showMessage(messageEl, '暂无可用口令，请先上传你的口令', 'error');
        }
    } catch (error) {
        showMessage(messageEl, '❌ 网络错误，请稍后重试', 'error');
        console.error('获取失败:', error);
    } finally {
        btn.disabled = false;
        btn.textContent = '随机获取一个口令';
    }
});

// 显示口令
function showCommand(element, content) {
    element.innerHTML = `
        <div class="command-text">${escapeHtml(content)}</div>
        <button class="copy-btn" onclick="copyCommand('${escapeHtml(content).replace(/'/g, "\\'")}')">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
            </svg>
            复制口令
        </button>
    `;
    element.classList.add('show');
}

// 复制口令
function copyCommand(text) {
    navigator.clipboard.writeText(text).then(() => {
        alert('✅ 口令已复制！现在打开腾讯元宝APP即可领取红包');
    }).catch(() => {
        // 降级方案
        const textarea = document.createElement('textarea');
        textarea.value = text;
        document.body.appendChild(textarea);
        textarea.select();
        document.execCommand('copy');
        document.body.removeChild(textarea);
        alert('✅ 口令已复制！现在打开腾讯元宝APP即可领取红包');
    });
}

// 显示消息
function showMessage(element, message, type) {
    element.textContent = message;
    element.className = `message ${type}`;
    element.style.display = 'block';
}

// HTML转义
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

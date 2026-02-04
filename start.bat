@echo off
echo 正在启动 Go 版本的元宝口令分享平台...
echo.

REM 检查是否已编译
if not exist yuanbao.exe (
    echo 首次运行，正在编译...
    go build -o yuanbao.exe main.go
    if errorlevel 1 (
        echo 编译失败！请检查 Go 环境是否正确安装。
        pause
        exit /b 1
    )
    echo 编译成功！
    echo.
)

echo 正在启动应用...
yuanbao.exe

pause

# yuluwallpaper

一款轻量级的壁纸自动切换工具，支持定时更换壁纸、系统托盘控制和开机自启动功能，专为Windows系统设计。

## 功能特点

- **自动切换壁纸**：定时从指定来源获取并更新桌面壁纸
- **系统托盘集成**：通过通知栏图标快速访问核心功能
- **开机自启动**：支持设置应用随系统启动
- **多平台适配**：针对Windows系统优化的壁纸设置逻辑
- **日志记录**：详细的运行日志便于问题排查

## 安装方法

### 前提条件
- Windows 10/11 操作系统
- Go 1.16+ 环境（用于源码编译）

### 二进制安装
1. 从 [Releases](https://github.com/frogchou/yuluwallpaper/releases) 页面下载最新版本
2. 解压到任意目录
3. 双击 `yuluwallpaper.exe` 运行

### 源码编译
```powershell
# 克隆仓库
git clone https://github.com/frogchou/yuluwallpaper.git
cd yuluwallpaper

# 设置代理（可选，国内用户）
.et_proxy.ps1

# 编译项目
.uild_windows.ps1

# 运行程序
cd bin
.uluwallpaper.exe
```

## 使用说明

### 基本操作
1. 运行程序后，系统托盘会出现应用图标
2. 右键点击图标打开菜单：
   - **刷新壁纸**：立即更新当前壁纸
   - **显示设置**：打开设置界面（开发中）
   - **退出**：关闭应用程序

### 配置文件
配置文件位于 `%APPDATA%\yuluwallpaper\config.json`，可自定义以下参数：
- `update_interval`：壁纸更新间隔（分钟）
- `wallpaper_source`：壁纸来源配置
- `startup`：是否开机自启动

## 开发指南

### 项目结构
```
├── cmd/yuluwallpaper       # 主程序入口
├── internal/app            # 应用核心逻辑
├── internal/autostart      # 自启动功能实现
├── internal/config         # 配置管理
├── internal/logger         # 日志系统
└── internal/wallpaper      # 壁纸设置功能
```

### 依赖管理
```powershell
# 安装依赖
.et_proxy.ps1

# 运行开发版本
.
un_with_logs.ps1
```

### 构建命令
```powershell
# 构建Windows可执行文件
.uild_windows.ps1
```

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

# danmaku2ass

[English](#english) | [中文](#chinese)

## English

danmaku2ass is a command-line tool written in Go that converts danmaku (comment) files from various streaming platforms into ASS subtitle format. It supports popular platforms including Bilibili, Niconico, and AcFun.

### Features

- Convert danmaku files to ASS subtitle format
- Support multiple streaming platforms:
  - Bilibili
  - Niconico
  - AcFun
- Automatic format detection
- Customizable font settings and display parameters
- Batch processing of multiple input files

### Installation

```bash
# Using go install
go install github.com/m13253/danmaku2ass@latest

# Or clone and build from source
git clone https://github.com/m13253/danmaku2ass.git
cd danmaku2ass
go build
```

### Usage

Basic usage:
```bash
danmaku2ass -s 1920x1080 input.xml
```

With all available options:
```bash
danmaku2ass [options] input_file [input_file...]

Options:
  -o string
        Output file path (default: input_name.ass)
  -s string
        Screen size in the format WIDTHxHEIGHT (default: "320x240")
  -fn string
        Font name (default: "MS PGothic")
  -fs float
        Font size (default: 48)
  -a float
        Alpha value (default: 0.8)
  -dm float
        Duration margin (default: 5)
  -ds float
        Duration start (default: 5)
```

### Example

Convert a Bilibili XML file to ASS format with custom settings:
```bash
danmaku2ass -s 1920x1080 -fn "Microsoft YaHei" -fs 36 -a 0.7 input.xml
```

![Screenshot](screenshot.jpg)

## Chinese

danmaku2ass 是一个用 Go 语言编写的命令行工具，可以将各大视频平台的弹幕文件转换为 ASS 字幕格式。支持包括哔哩哔哩、Niconico、AcFun 等平台。

### 功能特点

- 将弹幕文件转换为 ASS 字幕格式
- 支持多个视频平台：
  - 哔哩哔哩（Bilibili）
  - Niconico
  - AcFun
- 自动检测弹幕格式
- 可自定义字体设置和显示参数
- 支持批量处理多个输入文件

### 安装方法

```bash
# 使用 go install 安装
go install github.com/m13253/danmaku2ass@latest

# 或者克隆源码编译
git clone https://github.com/m13253/danmaku2ass.git
cd danmaku2ass
go build
```

### 使用方法

基本用法：
```bash
danmaku2ass -s 1920x1080 input.xml
```

所有可用选项：
```bash
danmaku2ass [选项] 输入文件 [输入文件...]

选项说明：
  -o string
        输出文件路径（默认：输入文件名.ass）
  -s string
        屏幕尺寸，格式为 宽x高（默认："320x240"）
  -fn string
        字体名称（默认："MS PGothic"）
  -fs float
        字体大小（默认：48）
  -a float
        透明度（默认：0.8）
  -dm float
        弹幕持续时间边界值（默认：5）
  -ds float
        弹幕开始时间偏移（默认：5）
```

### 使用示例

将哔哩哔哩的 XML 弹幕文件转换为 ASS 格式，并自定义设置：
```bash
danmaku2ass -s 1920x1080 -fn "Microsoft YaHei" -fs 36 -a 0.7 input.xml
```

### 许可证

本项目基于 GPL-3.0 许可证开源。

### 贡献

欢迎提交 Issue 和 Pull Request！如果你制作了更好的转换效果，请通过 [Issue](https://github.com/m13253/danmaku2ass/issues) 提交你的作品。


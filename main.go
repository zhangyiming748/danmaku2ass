// Package main 实现了一个弹幕转ASS字幕的命令行工具
// 支持从Bilibili、Niconico和AcFun等平台的弹幕文件转换为ASS字幕格式
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/m13253/danmaku2ass/ass"
	"github.com/m13253/danmaku2ass/parser"
)

const (
	// DefaultSizeWidth 定义默认视频宽度
	DefaultSizeWidth = 320
	// DefaultSizeHeight 定义默认视频高度
	DefaultSizeHeight = 240
)

// Config 存储程序运行所需的所有配置参数
type Config struct {
	OutputFile     string   // 输出ASS文件的路径
	ScreenSize     string   // 视频尺寸，格式为"宽x高"
	FontName       string   // 字幕字体名称
	FontSize       float64  // 字幕字体大小
	Alpha          float64  // 字幕透明度(0-1)
	DurationMargin float64  // 弹幕持续时间边界值
	DurationStart  float64  // 弹幕开始时间偏移
	InputFiles     []string // 输入的弹幕文件列表
	Width          int      // 解析后的视频宽度
	Height         int      // 解析后的视频高度
}

// parseArgs 解析命令行参数并返回配置对象
// 支持的参数包括：
// -o: 输出文件路径
// -s: 屏幕尺寸(宽x高)
// -fn: 字体名称
// -fs: 字体大小
// -a: 透明度
// -dm: 持续时间边界
// -ds: 开始时间偏移
func parseArgs() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.OutputFile, "o", "", "Output file path")
	flag.StringVar(&cfg.ScreenSize, "s", fmt.Sprintf("%dx%d", DefaultSizeWidth, DefaultSizeHeight), "Screen size in the format WIDTHxHEIGHT")
	flag.StringVar(&cfg.FontName, "fn", "MS PGothic", "Font name")
	flag.Float64Var(&cfg.FontSize, "fs", 48, "Font size")
	flag.Float64Var(&cfg.Alpha, "a", 0.8, "Alpha value")
	flag.Float64Var(&cfg.DurationMargin, "dm", 5, "Duration margin")
	flag.Float64Var(&cfg.DurationStart, "ds", 5, "Duration start")

	flag.Parse()

	// Get input files from remaining arguments
	cfg.InputFiles = flag.Args()
	if len(cfg.InputFiles) == 0 {
		return nil, fmt.Errorf("no input files specified")
	}

	// If output file is not specified, use the first input file name with .ass extension
	if cfg.OutputFile == "" {
		base := filepath.Base(cfg.InputFiles[0])
		ext := filepath.Ext(base)
		cfg.OutputFile = base[:len(base)-len(ext)] + ".ass"
	}

	// Parse screen size
	parts := strings.Split(cfg.ScreenSize, "x")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid screen size format: %s", cfg.ScreenSize)
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid screen width: %s", parts[0])
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid screen height: %s", parts[1])
	}

	cfg.Width = width
	cfg.Height = height

	return cfg, nil
}

// main 程序入口函数
// 主要流程：
// 1. 解析命令行参数
// 2. 创建ASS生成器
// 3. 处理所有输入文件
// 4. 生成最终的ASS文件
func main() {
	cfg, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// Create ASS generator
	generator := ass.NewGenerator(
		cfg.Width,
		cfg.Height,
		cfg.FontName,
		cfg.FontSize,
		cfg.Alpha,
		cfg.DurationStart,
		cfg.DurationMargin,
	)

	// Process all input files
	var allComments []parser.Comment
	for _, inputFile := range cfg.InputFiles {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", inputFile, err)
			continue
		}
		defer file.Close()

		// Detect format
		format, err := parser.ProbeFormat(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error detecting format of %s: %v\n", inputFile, err)
			continue
		}

		// Parse comments
		comments, err := parser.ParseComments(file, format, cfg.FontSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", inputFile, err)
			continue
		}

		allComments = append(allComments, comments...)
	}

	// Generate ASS file
	if err := generator.GenerateASS(allComments, cfg.OutputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating ASS file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted to %s\n", cfg.OutputFile)
}

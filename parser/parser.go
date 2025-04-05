// Package parser 提供了弹幕文件的解析功能
// 支持Bilibili、Niconico和AcFun三种主流弹幕格式的解析
// 将不同格式的弹幕文件统一转换为标准的Comment结构
package parser

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Comment 表示单条弹幕的结构体
// 包含弹幕的所有基本属性，如显示时间、位置、颜色等
type Comment struct {
	Timeline  float64 // 弹幕在视频中的显示时间点（秒）
	Timestamp int64   // 弹幕发送时的UNIX时间戳
	No        int     // 弹幕的序号
	Text      string  // 弹幕文本内容
	Position  int     // 弹幕位置类型：0=滚动弹幕，1=顶部固定，2=底部固定，3=逆向滚动
	Color     int     // 弹幕颜色，格式为0xRRGGBB
	Size      float64 // 弹幕字体大小
	Height    float64 // 弹幕预估高度（像素）
	Width     float64 // 弹幕预估宽度（像素）
}

// Format 表示弹幕文件的格式类型
type Format string

// 支持的弹幕格式常量定义
const (
	FormatBilibili Format = "Bilibili" // B站弹幕格式
	FormatNiconico Format = "Niconico" // N站弹幕格式
	FormatAcfun    Format = "Acfun"    // A站弹幕格式
)

// ProbeFormat 检测弹幕文件的格式类型
// 通过读取文件开头的内容来判断是哪种弹幕格式
// 支持检测Bilibili(XML格式)、Niconico(XML格式)和AcFun(JSON格式)三种格式
//
// 参数：
//   - file: 要检测格式的弹幕文件
//
// 返回值：
//   - Format: 检测到的弹幕格式
//   - error: 如果发生错误或无法识别格式则返回错误
func ProbeFormat(file *os.File) (Format, error) {
	// 保存当前文件位置
	curPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}
	defer file.Seek(curPos, io.SeekStart)

	// 读取文件开头部分用于判断格式
	buf := make([]byte, 100)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	content := string(buf[:n])

	// 根据文件内容特征判断格式
	if strings.HasPrefix(content, "<?xml") {
		if strings.Contains(content, "<i>") {
			return FormatBilibili, nil // B站XML格式
		} else if strings.Contains(content, "<chat>") {
			return FormatNiconico, nil // N站XML格式
		}
	} else if strings.HasPrefix(content, "[") {
		return FormatAcfun, nil // A站JSON格式
	}

	return "", fmt.Errorf("unknown format")
}

// ParseComments 解析弹幕文件中的所有弹幕
// 根据指定的格式类型调用相应的解析函数
//
// 参数：
//   - file: 要解析的弹幕文件
//   - format: 弹幕文件的格式类型
//   - fontSize: 基准字体大小，用于计算弹幕实际显示大小
//
// 返回值：
//   - []Comment: 解析出的所有弹幕列表
//   - error: 如果解析过程中发生错误则返回错误
func ParseComments(file *os.File, format Format, fontSize float64) ([]Comment, error) {
	switch format {
	case FormatBilibili:
		return parseBilibili(file, fontSize)
	case FormatNiconico:
		return parseNiconico(file, fontSize)
	case FormatAcfun:
		return parseAcfun(file, fontSize)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// calculateLength 计算文本宽度的辅助函数
// 目前使用简化版本：按字符数计算
// TODO: 实现更准确的文本宽度计算，考虑：
// 1. 不同字符的实际宽度（中文、英文、符号等）
// 2. 字体特性（比如等宽字体vs比例字体）
// 3. 字体大小的影响
//
// 参数：
//   - text: 要计算宽度的文本
//
// 返回值：
//   - float64: 文本的预估宽度
func calculateLength(text string) float64 {
	// TODO: 实现更准确的文本宽度计算
	return float64(len([]rune(text)))
}

// Package ass 实现了ASS字幕文件的生成功能
// 将解析后的弹幕数据转换为ASS字幕格式，支持多种弹幕样式和位置
package ass

import (
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/m13253/danmaku2ass/parser"
)

// Style 定义ASS字幕样式
// 包含字体、颜色、大小等样式属性
type Style struct {
	Name         string  // 样式名称
	FontName     string  // 字体名称
	FontSize     float64 // 字体大小
	PrimaryColor int     // 主要颜色(0xRRGGBB格式)
	Alpha        float64 // 透明度(0-1)
}

// Event 表示ASS对话事件
// 包含字幕的时间、样式和显示内容等信息
type Event struct {
	Start   float64 // 开始时间(秒)
	End     float64 // 结束时间(秒)
	Style   string  // 使用的样式名称
	Text    string  // 显示文本
	MarginL int     // 左边距
	MarginR int     // 右边距
	MarginV int     // 垂直边距
	Effect  string  // 特效名称
}

// Generator 处理ASS字幕的生成
// 包含所有必要的配置参数和生成方法
type Generator struct {
	Width         int     // 视频宽度
	Height        int     // 视频高度
	FontName      string  // 字体名称
	FontSize      float64 // 字体大小
	Alpha         float64 // 透明度
	DurationStart float64 // 弹幕持续时间
	MarginStart   float64 // 边距起始值
}

// NewGenerator 创建一个新的ASS生成器
// 参数：
//   - width: 视频宽度
//   - height: 视频高度
//   - fontName: 字体名称
//   - fontSize: 字体大小
//   - alpha: 透明度(0-1)
//   - durationStart: 弹幕持续时间
//   - marginStart: 边距起始值
func NewGenerator(width, height int, fontName string, fontSize, alpha, durationStart, marginStart float64) *Generator {
	return &Generator{
		Width:         width,
		Height:        height,
		FontName:      fontName,
		FontSize:      fontSize,
		Alpha:         alpha,
		DurationStart: durationStart,
		MarginStart:   marginStart,
	}
}

// GenerateASS 从弹幕评论生成ASS字幕文件
// 主要步骤：
// 1. 按时间线对弹幕进行排序
// 2. 创建输出文件
// 3. 写入ASS文件头部信息
// 4. 生成并写入字幕事件
//
// 参数：
//   - comments: 解析后的弹幕列表
//   - output: 输出ASS文件的路径
//
// 返回值：
//   - error: 如果生成过程中发生错误则返回错误
func (g *Generator) GenerateASS(comments []parser.Comment, output string) error {
	// 按时间线对弹幕进行排序
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].Timeline < comments[j].Timeline
	})

	// 创建输出文件
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入ASS文件头部
	g.writeHeader(file)

	// 生成并写入事件
	events := g.generateEvents(comments)
	g.writeEvents(file, events)

	return nil
}

// writeHeader 写入ASS文件的头部信息
// 包括脚本信息和样式定义
// 主要写入：
// 1. 脚本基本信息（分辨率、比例等）
// 2. 样式格式定义
// 3. 默认样式配置
func (g *Generator) writeHeader(file *os.File) {
	// 生成脚本信息部分
	header := fmt.Sprintf(`[Script Info]
ScriptType: v4.00+
PlayResX: %d
PlayResY: %d
Aspect Ratio: %f
Collisions: Normal
WrapStyle: 2
ScaledBorderAndShadow: yes

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
`, g.Width, g.Height, float64(g.Width)/float64(g.Height))

	// Write default styles
	styles := []Style{
		{Name: "R2L", FontName: g.FontName, FontSize: g.FontSize},
		{Name: "Top", FontName: g.FontName, FontSize: g.FontSize},
		{Name: "Bottom", FontName: g.FontName, FontSize: g.FontSize},
	}

	for _, style := range styles {
		header += fmt.Sprintf("Style: %s,%s,%f,&H%X,&H%X,&H000000,&H000000,0,0,0,0,100,100,0,0,1,2,0,2,20,20,2,0\n",
			style.Name, style.FontName, style.FontSize,
			int(g.Alpha*255)<<24, int(g.Alpha*255)<<24)
	}

	header += "\n[Events]\nFormat: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n"
	file.WriteString(header)
}

// generateEvents 从弹幕列表生成ASS事件列表
// 将每条弹幕转换为对应的ASS字幕事件
//
// 参数：
//   - comments: 解析后的弹幕列表
//
// 返回值：
//   - []Event: 生成的ASS事件列表
func (g *Generator) generateEvents(comments []parser.Comment) []Event {
	events := make([]Event, 0, len(comments))

	for _, comment := range comments {
		// 转换时间线为ASS时间格式
		start := comment.Timeline
		end := start + g.DurationStart

		// 根据弹幕位置确定样式
		var style string
		switch comment.Position {
		case 0: // 从右到左滚动
			style = "R2L"
		case 1: // 顶部固定
			style = "Top"
		case 2: // 底部固定
			style = "Bottom"
		default:
			continue
		}

		// 创建事件
		events = append(events, Event{
			Start:   start,
			End:     end,
			Style:   style,
			Text:    comment.Text,
			MarginL: 0,
			MarginR: 0,
			MarginV: 0,
		})
	}

	return events
}

// writeEvents 将ASS事件列表写入文件
// 将每个事件转换为ASS对话行格式并写入
//
// 参数：
//   - file: 要写入的文件
//   - events: 要写入的事件列表
func (g *Generator) writeEvents(file *os.File, events []Event) {
	for _, event := range events {
		// 将时间转换为ASS格式 (H:MM:SS.cc)
		start := formatTime(event.Start)
		end := formatTime(event.End)

		// 写入事件行
		line := fmt.Sprintf("Dialogue: 0,%s,%s,%s,,0,0,0,,%s\n",
			start, end, event.Style, event.Text)
		file.WriteString(line)
	}
}

// formatTime 将秒数转换为ASS时间格式 (H:MM:SS.cc)
// 例如：123.45秒会被转换为0:02:03.45
//
// 参数：
//   - seconds: 要转换的秒数
//
// 返回值：
//   - string: ASS格式的时间字符串
func formatTime(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60
	centisecs := int(math.Floor((seconds - math.Floor(seconds)) * 100))

	return fmt.Sprintf("%d:%02d:%02d.%02d", hours, minutes, secs, centisecs)
}

// Package parser 实现弹幕解析功能
package parser

import (
	"encoding/json"
	"os"
	"strings"
)

// AcfunComment 表示A站弹幕的JSON结构
// A站弹幕使用JSON数组格式，每条弹幕包含以下字段：
// {
//   "time": 12.34,     // 出现时间（秒）
//   "mode": 1,        // 弹幕模式
//   "size": 25,       // 字体大小
//   "color": 16777215,// 颜色值（十进制RGB）
//   "content": "text" // 弹幕内容
// }
type AcfunComment struct {
	Time    float64 `json:"time"`    // 弹幕出现时间（秒）
	Mode    int     `json:"mode"`    // 弹幕模式（1=滚动，4=底部，5=顶部，6=逆向）
	Size    int     `json:"size"`    // 字体大小（25为标准大小）
	Color   int     `json:"color"`   // 字体颜色（十进制RGB值）
	Content string  `json:"content"` // 弹幕文本内容
}
// parseAcfun 解析A站格式的弹幕文件
// A站弹幕使用JSON格式，将JSON数组解析为统一的Comment结构
//
// 参数：
//   - file: 要解析的弹幕文件
//   - fontSize: 基准字体大小
//
// 返回值：
//   - []Comment: 解析出的弹幕列表
//   - error: 解析错误
func parseAcfun(file *os.File, fontSize float64) ([]Comment, error) {
	// 解析JSON数组
	var acComments []AcfunComment
	if err := json.NewDecoder(file).Decode(&acComments); err != nil {
		return nil, err
	}

	comments := make([]Comment, 0, len(acComments))
	for i, c := range acComments {
		// 将A站的弹幕模式转换为统一的位置类型
		var position int
		switch c.Mode {
		case 1:
			position = 0 // 从右到左滚动弹幕
		case 4:
			position = 2 // 底部固定弹幕
		case 5:
			position = 1 // 顶部固定弹幕
		case 6:
			position = 3 // 从左到右滚动弹幕
		default:
			continue // 跳过不支持的模式
		}

		// 计算弹幕文本尺寸
		// A站字体大小以25为基准，需要根据fontSize进行缩放
		textSize := float64(c.Size) * fontSize / 25.0
		// 处理换行符
		text := strings.Replace(c.Content, "/n", "\n", -1)
		// 计算文本高度（考虑换行）
		height := float64(strings.Count(text, "\n")+1) * textSize
		// 计算文本宽度
		width := calculateLength(text) * textSize

		comments = append(comments, Comment{
			Timeline:  c.Time,
			Timestamp: 0, // Acfun format doesn't include timestamp
			No:        i,
			Text:      text,
			Position:  position,
			Color:     c.Color,
			Size:      textSize,
			Height:    height,
			Width:     width,
		})
	}

	return comments, nil
}

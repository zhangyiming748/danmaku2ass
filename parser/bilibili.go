// Package parser 实现弹幕解析功能
package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// BilibiliComment 表示B站弹幕的XML结构
// B站弹幕XML格式示例：
// <d p="时间,模式,字体大小,颜色,时间戳,弹幕池,用户ID,弹幕ID">弹幕内容</d>
type BilibiliComment struct {
	XMLName xml.Name `xml:"d"`         // XML标签名为d
	P       string   `xml:"p,attr"`    // p属性包含弹幕信息
	Content string   `xml:",chardata"` // 弹幕文本内容
}

// BilibiliXML 表示B站弹幕文件的根XML结构
type BilibiliXML struct {
	XMLName  xml.Name          `xml:"i"` // 根节点标签名为i
	Comments []BilibiliComment `xml:"d"` // 所有弹幕评论
}

// parseBilibili 解析B站格式的弹幕文件
// B站弹幕文件使用XML格式，每条弹幕包含详细的属性信息
func parseBilibili(file *os.File, fontSize float64) ([]Comment, error) {
	var biliXML BilibiliXML
	if err := xml.NewDecoder(file).Decode(&biliXML); err != nil {
		return nil, err
	}

	comments := make([]Comment, 0, len(biliXML.Comments))
	for i, c := range biliXML.Comments {
		// 解析p属性（格式：时间,模式,字体大小,颜色,时间戳,弹幕池,用户ID,弹幕ID）
		var (
			timeline  float64
			mode      string
			size      int
			color     int
			timestamp int64
		)

		_, err := fmt.Sscanf(c.P, "%f,%s,%d,%d,%d", &timeline, &mode, &size, &color, &timestamp)
		if err != nil {
			continue // Skip invalid comments
		}

		// 将B站的弹幕模式转换为统一的位置类型
		var position int
		switch mode {
		case "1":
			position = 0 // 从右到左滚动弹幕
		case "4":
			position = 2 // 底部固定弹幕
		case "5":
			position = 1 // 顶部固定弹幕
		case "6":
			position = 3 // 从左到右滚动弹幕
		default:
			continue // Skip unsupported modes
		}

		// 计算弹幕文本尺寸
		textSize := float64(size) * fontSize / 25.0
		text := strings.Replace(c.Content, "/n", "\n", -1)
		height := float64(strings.Count(text, "\n")+1) * textSize
		width := calculateLength(text) * textSize

		comments = append(comments, Comment{
			Timeline:  timeline,
			Timestamp: timestamp,
			No:        i,
			Text:      text,
			Position:  position,
			Color:     color,
			Size:      textSize,
			Height:    height,
			Width:     width,
		})
	}

	return comments, nil
}

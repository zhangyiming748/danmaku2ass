// Package parser 实现弹幕解析功能
package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// NiconicoComment 表示N站弹幕的XML结构
// N站弹幕XML格式示例：
// <chat vpos="100" no="1" date="1234567890" user_id="user1" mail="184">弹幕内容</chat>
type NiconicoComment struct {
	XMLName xml.Name `xml:"chat"`         // XML标签名为chat
	VPos    int      `xml:"vpos,attr"`    // 视频位置（1/100秒）
	No      int      `xml:"no,attr"`      // 弹幕序号
	Date    int64    `xml:"date,attr"`    // 发送时间戳
	UserID  string   `xml:"user_id,attr"` // 用户ID
	Mail    string   `xml:"mail,attr"`    // 命令字符串
	Content string   `xml:",chardata"`    // 弹幕内容
}

// NiconicoXML 表示N站弹幕文件的根XML结构
type NiconicoXML struct {
	XMLName  xml.Name          `xml:"packet"` // 根节点标签名为packet
	Comments []NiconicoComment `xml:"chat"`   // 所有弹幕评论
}

// parseNiconico 解析N站格式的弹幕文件
// N站弹幕使用XML格式，每条弹幕包含位置、颜色等命令信息
//
// mail属性包含以空格分隔的命令，常见命令：
// - ue: 顶部固定弹幕
// - shita: 底部固定弹幕
// - big: 大号字体
// - small: 小号字体
// - 颜色值: 6位16进制颜色值
func parseNiconico(file *os.File, fontSize float64) ([]Comment, error) {
	var nicoXML NiconicoXML
	if err := xml.NewDecoder(file).Decode(&nicoXML); err != nil {
		return nil, err
	}

	comments := make([]Comment, 0, len(nicoXML.Comments))
	for _, c := range nicoXML.Comments {
		// 解析mail命令
		var position int
		var color int = 0xFFFFFF // 默认颜色为白色
		var size float64 = fontSize

		commands := strings.Split(c.Mail, " ")
		for _, cmd := range commands {
			switch cmd {
			case "ue":
				position = 1 // 顶部固定
			case "shita":
				position = 2 // 底部固定
			case "big":
				size = fontSize * 1.5 // 1.5倍字体大小
			case "small":
				size = fontSize * 0.5 // 0.5倍字体大小
			default:
				// 尝试解析颜色值
				if len(cmd) == 6 {
					if _, err := fmt.Sscanf(cmd, "%x", &color); err == nil {
						continue
					}
				}
			}
		}

		// Calculate text dimensions
		text := strings.Replace(c.Content, "/n", "\n", -1)
		height := float64(strings.Count(text, "\n")+1) * size
		width := calculateLength(text) * size

		// Convert vpos (1/100 seconds) to timeline (seconds)
		timeline := float64(c.VPos) / 100.0

		comments = append(comments, Comment{
			Timeline:  timeline,
			Timestamp: c.Date,
			No:        c.No,
			Text:      text,
			Position:  position,
			Color:     color,
			Size:      size,
			Height:    height,
			Width:     width,
		})
	}

	return comments, nil
}

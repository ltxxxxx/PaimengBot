package images

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
)

// PasteStringDefault 以默认方式（默认字体、黑色、Wrapped、左上角定位、文字居左）贴文字
func (img *ImageCtx) PasteStringDefault(str string, fontSize, lineSpace float64, x, y, width float64) error {
	img.Push()
	defer img.Pop()
	if err := img.UseDefaultFont(fontSize); err != nil {
		return err
	} // 默认字体
	img.SetRGB(0, 0, 0) // 纯黑色
	img.DrawStringWrapped(str, x, y, 0, 0, width, lineSpace, gg.AlignLeft)
	return nil
}

// PasteLine 画线
func (img *ImageCtx) PasteLine(x1, y1, x2, y2, lineWidth float64, colorStr string) {
	img.Push()
	defer img.Pop()
	img.SetColorAuto(colorStr)
	img.DrawLine(x1, y1, x2, y2)
	img.SetLineWidth(lineWidth)
	img.Stroke()
}

type Point struct {
	X, Y float64
}

// DrawStar 绘制星星 n: 角数; (x, y): 圆心坐标; r: 圆半径
func (img *ImageCtx) DrawStar(n int, x, y, r float64) {
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		points[i] = Point{x + r*math.Cos(a), y + r*math.Sin(a)}
	}
	for i := 0; i < n+1; i++ {
		index := (i * 2) % n
		p := points[index]
		img.LineTo(p.X, p.Y)
	}
}

var colorMap map[string]string = map[string]string{
	"white":  "#ffffff",
	"black":  "#000000",
	"gray":   "#a4b0be",
	"red":    "#e74c3c",
	"blue":   "#3498db",
	"green":  "#2ecc71",
	"yellow": "#ffd43b",
}

func (img *ImageCtx) SetColorAuto(colorStr string) {
	if res, ok := colorMap[colorStr]; ok {
		img.SetHexColor(res)
		return
	}
	if strings.HasPrefix(colorStr, "#") {
		img.SetHexColor(colorStr)
		return
	}
	colorStr = strings.ToLower(colorStr)
	if strings.HasPrefix(colorStr, "rgb") {
		colorStr = strings.ReplaceAll(strings.ReplaceAll(colorStr, " ", ""), "\t", "")
		reg := regexp.MustCompile("rgba?\\((\\d{1,3}),(\\d{1,3}),(\\d{1,3})(,\\d{1,3}\\.?\\d*)?\\)")
		sub := reg.FindStringSubmatch(colorStr)
		if len(sub) <= 4 {
			return
		}
		r, _ := strconv.ParseInt(sub[1], 10, 32)
		g, _ := strconv.ParseInt(sub[2], 10, 32)
		b, _ := strconv.ParseInt(sub[3], 10, 32)
		var a int64 = 255
		if len(sub[4]) > 0 {
			sub[4] = sub[4][1:]
		}
		if strings.Contains(sub[4], ".") {
			sub[4] += "0"
			fa, _ := strconv.ParseFloat(sub[4], 32)
			a = int64(255.0 * fa)
		} else if len(sub[4]) > 0 {
			a, _ = strconv.ParseInt(sub[4], 10, 32)
		}
		img.SetRGBA255(int(r), int(g), int(b), int(a))
		return
	}
	img.SetHexColor("#ffffff") // 兜底 纯白
}

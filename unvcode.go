package unvcode

import (
	"fmt"
	"github.com/grd/statistics"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/unicode/norm"
	"image"
	"image/draw"
)

var d = make(map[string][]string)

// Unvcode struct
type Unvcode struct {
	Face font.Face
	// SkipAscii 开启则会跳过ascii字符
	SkipAscii bool
	// Mse 字符相似度的阈值
	Mse float64
}

func init() {
	for i := 0; i < 65536; i++ {
		字 := string(i)
		新字 := norm.NFKC.String(字)
		if 字 != 新字 {
			d[新字] = append(d[新字], 字)
		}
	}
}

// New 使用字体文件创建一个 Unvcode 对象
func New(fontByte []byte) (*Unvcode, error) {
	var u Unvcode
	if len(fontByte) != 0 {
		f, err := opentype.Parse(fontByte)
		if err != nil {
			return nil, fmt.Errorf("Parse: %v ", err)
		}
		u.Face, err = opentype.NewFace(f, &opentype.FaceOptions{
			Size: 64,
			DPI:  72,
		})
		if err != nil {
			return nil, fmt.Errorf("NewFace: %v", err)
		}
	} else {
		return nil, fmt.Errorf("Empty font file ")
	}
	u.SkipAscii = true
	u.Mse = 0.1
	return &u, nil
}

func (u *Unvcode) 画皮(字 string) []uint8 {
	img := image.NewGray(image.Rect(0, 0, 100, 100))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)
	d := font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: u.Face,
		Dot:  fixed.P(u.Face.Metrics().Descent.Round(), u.Face.Metrics().Ascent.Round()),
	}
	d.DrawString(字)
	for i, p := range img.Pix {
		img.Pix[i] = p / 255
	}
	return img.Pix
}

func (u *Unvcode) 比较(字1, 字2 string) float64 {
	var 皮1, 皮2 = u.画皮(字1), u.画皮(字2)
	var 差 statistics.Int64
	for i := 0; i < 10000; i++ {
		差 = append(差, int64(皮1[i])-int64(皮2[i]))
	}
	return statistics.Variance(&差)
}

func (u *Unvcode) 假面(字 int32, skipAscii bool, mse float64) (float64, string) {
	if 字 < 128 && skipAscii {
		return 0.0, string(字)
	}
	候选组 := d[string(字)]
	if len(候选组) == 0 {
		return -1.0, string(字)
	}
	差异组 := make(map[string]float64)
	for _, s := range 候选组 {
		差异组[s] = u.比较(string(字), s)
	}
	差异, 新字 := 1.0, ""
	for k, v := range 差异组 {
		if v < 差异 || 差异 == 0 {
			差异 = v
			新字 = k
		}
	}
	return 差异, 新字
}

// Parse 解析字符串
func (u *Unvcode) Parse(str string) (string, []float64) {
	差异, 串 := []float64{}, ""
	for _, s := range str {
		f, c := u.假面(s, u.SkipAscii, u.Mse)
		差异 = append(差异, f)
		串 += c
	}
	return 串, 差异
}

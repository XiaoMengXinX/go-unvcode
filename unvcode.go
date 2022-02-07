package unvcode

import (
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/unicode/norm"
	"image"
)

var d = make(map[string][]string)

// Unv struct
type Unv struct {
	face font.Face
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

// New 使用字体文件创建一个 Unv 对象
func New(fontByte []byte) (*Unv, error) {
	var u Unv
	if len(fontByte) != 0 {
		f, err := opentype.Parse(fontByte)
		if err != nil {
			return nil, fmt.Errorf("Parse: %v ", err)
		}
		u.face, err = opentype.NewFace(f, &opentype.FaceOptions{
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

func (u *Unv) 画皮(字 string) []uint8 {
	img := image.NewGray(image.Rect(0, 0, 100, 100))
	d := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: u.face,
		Dot:  fixed.P(u.face.Metrics().Descent.Round(), u.face.Metrics().Ascent.Round()),
	}
	d.DrawString(字)
	for i, p := range img.Pix {
		img.Pix[i] = (255 - p) / 255
	}
	return img.Pix
}

func (u *Unv) 比较(字1, 字2 string) float64 {
	var 皮1, 皮2 = u.画皮(字1), u.画皮(字2)
	var 差 []int32
	for i, p := range 皮1 {
		差 = append(差, int32(p)-int32(皮2[i]))
	}
	return variance(差)
}

func variance(a []int32) float64 {
	Len := len(a)
	var 和 int32
	for i := 0; i < Len; i++ {
		和 += a[i]
	}
	平均值 := float64(和) / float64(Len)
	var 平方和 float64
	for i := 0; i < Len; i++ {
		平方和 += (float64(a[i]) - 平均值) * (float64(a[i]) - 平均值)
	}
	return 平方和 / float64(Len)
}

func (u *Unv) 假面(字 int32) (float64, string) {
	候选组 := d[string(字)]
	if (字 < 128 && u.SkipAscii) || len(候选组) == 0 {
		return -1.0, string(字)
	}
	差异组 := make(map[string]float64)
	for _, s := range 候选组 {
		差异组[s] = u.比较(string(字), s)
	}
	差异, 新字 := 1.0, ""
	for k, v := range 差异组 {
		if v < 差异 && v < u.Mse {
			差异 = v
			新字 = k
		}
		if 差异 == 0.0 {
			差异 = v
			新字 = k
			break
		}
	}
	if 新字 == "" || 差异 > u.Mse {
		return -1.0, string(字)
	}
	return 差异, 新字
}

// Unvcode 解析字符串
func (u *Unv) Unvcode(str string) (string, []float64) {
	差异, 串 := []float64{}, ""
	for _, s := range str {
		f, c := u.假面(s)
		差异 = append(差异, f)
		串 += c
	}
	return 串, 差异
}

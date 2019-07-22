package tencent_code

import (
	"fmt"
	"funny/tencent_code/try_chromedp"
	"github.com/pkg/errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
)

const lineHeight = 15
const gap = 80

var line *image.RGBA

func init() {
	line = image.NewRGBA(image.Rect(0, 0, 0, lineHeight-1))
	for y := 0; y < lineHeight; y++ {
		line.Set(0, y, color.NRGBA{R: 255, A: 255})
	}
	//log.Println(colorDistance(color.RGBA{255,255,255, 255}, color.RGBA{255,255,255, 255}))
}

func Init() {
	try_chromedp.Test()
	return
	for i := 1; i <= 11; i++ {
		test(i)
	}
}

func test(i int) {
	log.Println("here")
	src, err := filepath.Abs(fmt.Sprintf("tencent_code/%d.jpeg", i))
	if err != nil {
		log.Fatal(err)
	}
	dest, err := filepath.Abs(fmt.Sprintf("tencent_code/out/%d.png", i))
	if err != nil {
		log.Fatal(err)
	}
	err = getPos1(src, dest)
	log.Println(err)
}

func getPos(src, dest string) (err error) {
	f, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "打开图片错误")
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		return errors.Wrap(err, "解码图片错误")
	}
	r := img.Bounds()
	newRgba := image.NewRGBA(r)
	draw.Draw(newRgba, r, img, image.ZP, draw.Src)
	for y := 0; y <= r.Max.Y; y += 5 {
		last := getLine(img, 0, y)
		for x := 1; x <= r.Max.X; x++ {
			cur := getLine(img, x, y)
			gapEnough := true
			sum := float64(0)
			for i, c := range last {
				sum += colorDistance(c, cur[i])
				if c.A < 230 {
					gapEnough = false
					break
				}
			}
			if gapEnough {
				log.Println(gapEnough)
			}
			if sum/lineHeight < 600 {
				gapEnough = false
			}
			if gapEnough {
				log.Println(x, y)
				//newRgba.Set(x, y, color.NRGBA{R: 255, A: 255})
				//draw.Draw(newRgba, line.Bounds(), line, image.Pt(x, y), draw.Src)
				for dy := 0; dy < lineHeight; dy++ {
					newRgba.Set(x-1, y+dy, color.NRGBA{R: 255, A: 255})
				}
			}
			last = cur
		}

	}
	p, _ := filepath.Abs(dest)
	f, err = os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "打开存输出件错误")
	}
	err = png.Encode(f, newRgba)
	if err != nil {
		return errors.Wrap(err, "图片编码错误")
	}
	f.Close()
	return
}
func getPos1(src, dest string) (err error) {
	f, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "打开图片错误")
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		return errors.Wrap(err, "解码图片错误")
	}
	r := img.Bounds()
	newRgba := image.NewRGBA(r)
	draw.Draw(newRgba, r, img, image.ZP, draw.Src)
	for y := 0; y <= r.Max.Y; y += lineHeight / 1 {
		last := getLine1(img, 0, y)
		for x := 1; x <= r.Max.X; x++ {
			cur := getLine1(img, x, y)
			gapEnough := true
			for i, c := range last {
				if c-cur[i] < gap || c < 250 {
					gapEnough = false
					break
				}
			}
			if gapEnough {
				log.Println(x, y)
				//newRgba.Set(x, y, color.NRGBA{R: 255, A: 255})
				//draw.Draw(newRgba, line.Bounds(), line, image.Pt(x, y), draw.Src)
				for dy := 0; dy < lineHeight; dy++ {
					newRgba.Set(x-1, y+dy, color.NRGBA{R: 255, A: 255})
				}
			}
			last = cur
		}

	}
	p, _ := filepath.Abs(dest)
	f, err = os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "打开存输出件错误")
	}
	err = png.Encode(f, newRgba)
	if err != nil {
		return errors.Wrap(err, "图片编码错误")
	}
	f.Close()
	return
}

func getLine(m image.Image, x, y int) []color.RGBA {
	arr := make([]color.RGBA, lineHeight)
	for dy := 0; dy < lineHeight; dy++ {
		c := m.At(x, y+dy)
		r, g, b, _ := c.RGBA()
		arr[dy] = color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8((r>>8 + g>>8 + b>>8) / 3)}
	}
	return arr
}
func getLine1(m image.Image, x, y int) []byte {
	arr := make([]byte, lineHeight)
	for dy := 0; dy < lineHeight; dy++ {
		c := m.At(x, y+dy)
		r, g, b, _ := c.RGBA()
		arr[dy] = uint8((r>>8 + g>>8 + b>>8) / 3)
	}
	return arr
}

func colorDistance(c1, c2 color.RGBA) float64 {
	redMean := float64(c1.R+c2.R) / 2
	r := float64(c1.R - c2.R)
	g := float64(c1.G - c2.G)
	b := float64(c1.B - c2.B)
	return math.Sqrt(float64(
		(int((512+redMean)*r*r) >> 8) +
			int(4*g*g) +
			(int((767-redMean)*b*b) >> 8)))
}

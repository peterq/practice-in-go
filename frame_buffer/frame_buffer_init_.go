package frame_buffer

import (
	"github.com/bfanger/framebuffer"
	"image"
	"image/png"
	"log"
	"os"
	"time"
)

func Init() {
	glTest()
	return
	fb, err := framebuffer.Open("/dev/video0")
	if err != nil {
		log.Fatalf("could not open framebuffer: %v", err)
	}
	defer fb.Close()
	info, err := fb.VarScreenInfo()
	if err != nil {
		log.Fatalf("could not read screen info: %v", err)
	}
	log.Printf(`
Fixed:
%+v

Variable:
%+v
`, fb.FixScreenInfo, info)
	//rand.Read(fb.Buffer) // fill the buffer with noise

	img, err := os.OpenFile("1.png", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	//img.Write(fb.Buffer)
	//img.Close()

	m := make([]byte, len(fb.Buffer))
	copy(m, fb.Buffer)
	log.Println(m)
	nrgba := image.NewNRGBA(image.Rect(0, 0, 1920, 1080))
	copy(nrgba.Pix, m)
	png.Encode(img, nrgba)
	img.Close()
	copy(fb.Buffer, m)
	//log.Println(nrgba.Pix)
	go colorful(fb)
	<-time.After(30 * time.Second)
}
func colorful(fb *framebuffer.Device) {
	time.Sleep(1 * time.Second)
	m := make([]byte, len(fb.Buffer))
	i := 0
	for {
		r, g, b, _ := RGBA(i)
		for i := range fb.Buffer {
			switch (i + 1) % 4 {
			case 0: // R
				m[i] = byte(r)
			case 1:
				m[i] = byte(g)
			case 2:
				m[i] = byte(b)
			case 3:
				m[i] = 255
			}
		}
		copy(fb.Buffer, m)
		i += 1
		if i == 360 {
			i = 0
		}
		//log.Println(m[:40], r, g, b)
		time.Sleep(3 * time.Second / 360)
		//log.Println(r, g, b)
	}
}

func RGBA(H int) (r, g, b, a int) {
	S, V := 255, 255
	// Direct implementation of the graph in this image:
	// https://en.wikipedia.org/wiki/HSL_and_HSV#/media/File:HSV-RGB-comparison.svg
	max := V
	min := V * (255 - S)

	H %= 360
	segment := H / 60
	offset := H % 60
	mid := ((max - min) * offset) / 60

	//log.Println(H, max, min, mid)
	switch segment {
	case 0:
		return max, min + mid, min, 0xff
	case 1:
		return max - mid, max, min, 0xff
	case 2:
		return min, max, min + mid, 0xff
	case 3:
		return min, max - mid, max, 0xff
	case 4:
		return min + mid, min, max, 0xff
	case 5:
		return max, min, max - mid, 0xff
	}

	return 0, 0, 0, 0xff
}

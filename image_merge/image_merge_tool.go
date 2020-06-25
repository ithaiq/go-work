package image_merge

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"

	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
)

type Point struct {
	x, y int
}

//合并头像
func Merge(dst io.Writer, src []io.Reader) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Merge.recover:", r)
		}
	}()
	var err error
	imagePoints, scaleSize := getXYPoint(len(src))
	width := 500

	//创建背景大图
	background := image.NewRGBA(image.Rect(0, 0, width, width))
	//设置背景为灰色
	for m := 0; m < width; m++ {
		for n := 0; n < width; n++ {
			//rgba := GetRGBA(0xC8CED4)
			background.SetRGBA(m, n, color.RGBA{233, 233, 233, 0})
		}
	}
	//背景图矩形圆角
	newBg, err := CreateRoundRect(background, 10, &color.RGBA{255, 255, 255, 0})
	if err != nil {
		fmt.Errorf("CreateRoundRect has error:", err)
		return err
	}
	//开始合成
	for i, v := range imagePoints {
		x := v.x
		y := v.y

		fOut := bytes.NewBuffer([]byte{})
		//先缩略
		err = scale(src[i], fOut, scaleSize, scaleSize, 100)
		if err != nil {
			fmt.Errorf("scale has error:", err)
			return err
		}
		//矩形圆角
		rgba, err := CreateRoundRectWithoutColor(fOut, 0)
		if err != nil {
			fmt.Errorf("CreateRoundRectWithoutColor has error:", err)
			return err
		}

		draw.Draw(newBg, newBg.Bounds(), rgba, rgba.Bounds().Min.Sub(image.Pt(x, y)), draw.Src)
	}
	//return png.Encode(dst, newBg)
	return jpeg.Encode(dst, newBg, nil)
}

//获取每张图片位置和缩放大小
func getXYPoint(size int) ([]*Point, int) {
	s := make([]*Point, size)
	X := 500
	scaleSize := 0
	e := 20
	if size == 1 {
		scaleSize = X - e*2
		s[0] = &Point{e, e}
	}
	if size == 2 {
		scaleSize = (X - e*3) / 2
		temp := (X - scaleSize) / 2
		s[0] = &Point{e, temp}
		s[1] = &Point{e + scaleSize + e, temp}
	}
	if size == 3 {
		scaleSize = (X - e*3) / 2
		temp := e*2 + scaleSize
		s[0] = &Point{(X - scaleSize) / 2, e}
		s[1] = &Point{e, temp}
		s[2] = &Point{temp, temp}
	}
	if size == 4 {
		scaleSize = (X - e*3) / 2
		temp := e*2 + scaleSize
		s[0] = &Point{e, e}
		s[1] = &Point{temp, e}
		s[2] = &Point{e, temp}
		s[3] = &Point{temp, temp}
	}
	if size == 5 {
		scaleSize = (X - e*4) / 3
		temp := X - e - 2*scaleSize
		s[0] = &Point{temp / 2, temp / 2}
		s[1] = &Point{temp/2 + e + scaleSize, temp / 2}
		s[2] = &Point{e, temp/2 + e + scaleSize}
		s[3] = &Point{e + scaleSize + e, temp/2 + e + scaleSize}
		s[4] = &Point{2*scaleSize + 3*e, temp/2 + e + scaleSize}
	}
	if size == 6 {
		scaleSize = (X - e*4) / 3
		temp := X - e - 2*scaleSize
		s[0] = &Point{e, temp / 2}
		s[1] = &Point{e + scaleSize + e, temp / 2}
		s[2] = &Point{2*scaleSize + 3*e, temp / 2}
		s[3] = &Point{e, temp/2 + scaleSize + e}
		s[4] = &Point{e + scaleSize + e, temp/2 + scaleSize + e}
		s[5] = &Point{2*scaleSize + 3*e, temp/2 + scaleSize + e}
	}

	if size == 7 {
		scaleSize = (X - e*4) / 3
		temp := X - e*2 - scaleSize*3
		s[0] = &Point{(X - scaleSize) / 2, temp / 2}
		s[1] = &Point{e, temp/2 + e + scaleSize}
		s[2] = &Point{e + scaleSize + e, temp/2 + e + scaleSize}
		s[3] = &Point{2*scaleSize + 3*e, temp/2 + e + scaleSize}

		s[4] = &Point{e, temp/2 + 2*e + 2*scaleSize}
		s[5] = &Point{e + scaleSize + e, temp/2 + 2*e + 2*scaleSize}
		s[6] = &Point{2*scaleSize + 3*e, temp/2 + 2*e + 2*scaleSize}
	}
	if size == 8 {

		scaleSize = (X - e*4) / 3
		temp := X - e*2 - scaleSize*3
		s[0] = &Point{(X - e - 2*scaleSize) / 2, temp / 2} //
		s[1] = &Point{(X-e-2*scaleSize)/2 + e + scaleSize, temp / 2}

		s[2] = &Point{e, temp/2 + e + scaleSize}
		s[3] = &Point{e + scaleSize + e, temp/2 + e + scaleSize}
		s[4] = &Point{2*scaleSize + 3*e, temp/2 + e + scaleSize}

		s[5] = &Point{e, temp/2 + 2*e + 2*scaleSize}
		s[6] = &Point{e + scaleSize + e, temp/2 + 2*e + 2*scaleSize}
		s[7] = &Point{2*scaleSize + 3*e, temp/2 + 2*e + 2*scaleSize}
	}
	if size == 9 {
		scaleSize = (X - e*4) / 3
		temp := X - e*2 - scaleSize*3
		s[0] = &Point{e, temp / 2}
		s[1] = &Point{e + scaleSize + e, temp / 2}
		s[2] = &Point{2*scaleSize + 3*e, temp / 2}

		s[3] = &Point{e, temp/2 + e + scaleSize}
		s[4] = &Point{e + scaleSize + e, temp/2 + e + scaleSize}
		s[5] = &Point{2*scaleSize + 3*e, temp/2 + e + scaleSize}

		s[6] = &Point{e, temp/2 + 2*e + 2*scaleSize}
		s[7] = &Point{e + scaleSize + e, temp/2 + 2*e + 2*scaleSize}
		s[8] = &Point{2*scaleSize + 3*e, temp/2 + 2*e + 2*scaleSize}

	}
	return s, scaleSize
}

//图片缩略
func scale(in io.Reader, out io.Writer, width, height, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
	}
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Resize(uint(width), uint(height), origin, resize.Lanczos3)

	switch fm {
	case "jpeg":
		return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	case "png":
		return png.Encode(out, canvas)
	case "gif":
		return gif.Encode(out, canvas, &gif.Options{})
	case "bmp":
		return bmp.Encode(out, canvas)
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}

//创建矩形圆角
func CreateRoundRectWithoutColor(rd io.Reader, r int) (*image.RGBA, error) {
	src, _, err := image.Decode(rd)
	if err != nil {
		return nil, err
	}

	b := src.Bounds()
	x := b.Dx()
	y := b.Dy()
	dst := image.NewRGBA(b)

	p1 := image.Point{r, r}
	p2 := image.Point{x - r, r}
	p3 := image.Point{r, y - r}
	p4 := image.Point{x - r, y - r}

	for m := 0; m < x; m++ {
		for n := 0; n < y; n++ {
			if (p1.X-m)*(p1.X-m)+(p1.Y-n)*(p1.Y-n) > r*r && m <= p1.X && n <= p1.Y {
			} else if (p2.X-m)*(p2.X-m)+(p2.Y-n)*(p2.Y-n) > r*r && m > p2.X && n <= p2.Y {
			} else if (p3.X-m)*(p3.X-m)+(p3.Y-n)*(p3.Y-n) > r*r && m <= p3.X && n > p3.Y {
			} else if (p4.X-m)*(p4.X-m)+(p4.Y-n)*(p4.Y-n) > r*r && m > p4.X && n > p4.Y {
			} else {
				dst.Set(m, n, src.At(m, n))
			}
		}
	}
	return dst, nil
}

//创建矩形圆角
func CreateRoundRect(src *image.RGBA, r int, c *color.RGBA) (*image.RGBA, error) {
	b := src.Bounds()
	x := b.Dx()
	y := b.Dy()
	dst := image.NewRGBA(b)
	draw.Draw(dst, b, src, src.Bounds().Min, draw.Src)

	p1 := image.Point{r, r}
	p2 := image.Point{x - r, r}
	p3 := image.Point{r, y - r}
	p4 := image.Point{x - r, y - r}

	for m := 0; m < x; m++ {
		for n := 0; n < y; n++ {
			if (p1.X-m)*(p1.X-m)+(p1.Y-n)*(p1.Y-n) > r*r && m <= p1.X && n <= p1.Y {
				dst.Set(m, n, c)
			} else if (p2.X-m)*(p2.X-m)+(p2.Y-n)*(p2.Y-n) > r*r && m > p2.X && n <= p2.Y {
				dst.Set(m, n, c)
			} else if (p3.X-m)*(p3.X-m)+(p3.Y-n)*(p3.Y-n) > r*r && m <= p3.X && n > p3.Y {
				dst.Set(m, n, c)
			} else if (p4.X-m)*(p4.X-m)+(p4.Y-n)*(p4.Y-n) > r*r && m > p4.X && n > p4.Y {
				dst.Set(m, n, c)
			}
		}
	}
	return dst, nil
}

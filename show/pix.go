package show

import (
	"image"
	"image/color"
	"image/draw"
	"regexp"
)

var (
	ClrLine = color.RGBA{80, 80, 80, 255}		//	Цвет соединительных линий
	ClrPlus = color.RGBA{60, 60, 100, 255}		//	Цвет плюса и минуса
	ClrArr = color.RGBA{80, 80, 80, 255}
	ClrArrPass = color.RGBA{180, 180, 180, 255}
	ClrArrAct = color.RGBA{64, 64, 255, 255}
	ClrArrBord = color.RGBA{0, 0, 0, 255}		//	Цыет края
)

//	Построить изображение пиктограммы
//	Первый символ:
//		o - тянущаяся вертикальная линия
//		a - ?
//		s - символ переключения
//		p - плюс
//		m - минус
//		z - точка
//	Второй символ:
//	Далее вверх, вправо, вниз, влево:
//		0 - нет линии
//		1 - есть линия
func BuildPixRast(pix string) *image.RGBA {
	p1,p2,pe := "","",""
	if len(pix) > 0 { p1 = pix[:1] }
	if len(pix) > 1 { p2 = pix[1:2] }
	if match := regexp.MustCompile("(\\w)\\.").FindStringSubmatch(pix); match != nil { pe = match[1] }
	sx,sy := 20,20
	if p1 == "o" {
		sy = 1
	} else if p1 == "a" {
		sx = 12
	} else if p1 == "s" {
		sx = 10
	}
	img := image.NewRGBA(image.Rect(0, 0, sx, sy))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.ZP, draw.Src)
	if p1 == "p" || p1 == "m" || p1 == "z" {
		hf := 3
		if p1 == "p" || p1 == "m" {
			hf = 6
			for p := -hf + 2; p < hf - 2; p++ {
				img.Set(sx / 2 + p, sy / 2 - 1, ClrPlus)
				img.Set(sx / 2 + p, sy / 2, ClrPlus)
				if p1 == "p" {
					img.Set(sx / 2 - 1, sy / 2 + p, ClrPlus)
					img.Set(sx / 2, sy / 2 + p, ClrPlus)
				}
			}
		}
		for p := -hf; p < hf; p++ {
			img.Set(sx / 2 + p, sy / 2 - hf, ClrLine)
			img.Set(sx / 2 + p, sy / 2 + hf - 1, ClrLine)
			img.Set(sx / 2 - hf, sy / 2 + p, ClrLine)
			img.Set(sx / 2 + hf - 1, sy / 2 + p, ClrLine)
		}
		for dir := 0; dir < 4 && 2+dir < len(pix); dir++ {
			if pix[1+dir:2+dir] == "1" {
				for p := hf; p <= sx / 2 || p <= sy / 2; p++ {
					if dir == 0 && p <= sy / 2 {
						img.Set(sx / 2 - 1, sy / 2 - p, ClrLine)
						img.Set(sx / 2, sy / 2 - p, ClrLine)
					} else if dir == 1 && p < sx / 2 {
						img.Set(sx / 2 + p, sy / 2 - 1, ClrLine)
						img.Set(sx / 2 + p, sy / 2, ClrLine)
					} else if dir == 2 && p < sy / 2 {
						img.Set(sx / 2 - 1, sy / 2 + p, ClrLine)
						img.Set(sx / 2, sy / 2 + p, ClrLine)
					} else if dir == 3 && p <= sx / 2{
						img.Set(sx / 2 - p, sy / 2 - 1, ClrLine)
						img.Set(sx / 2 - p, sy / 2, ClrLine)
					}
				}
			}
		}
	} else if p1 == "a" {
		clr := ClrArrPass
		if pe == "a" {
			clr = ClrArr
		}
		yb := 4
		for y := yb; y < sy / 2; y++ {
			xl := 1
			xr := 1 + (y - yb) * (sx - 2) / (sy / 2)
			if p2 == "s" || p2 == "b" {
				xl, xr = sx - 1 - xr, sx - 1 - xl
			}
			for x := xl; x <= xr; x++ {
				img.Set(x, y, clr)
				img.Set(x, sy - 1 - y, clr)
			}
			if p2 == "s" {
				img.Set(2, y, clr)
				img.Set(2, sy - 1 - y, clr)
			} else if p2 == "e" {
				img.Set(sx - 3, y, clr)
				img.Set(sx - 3, sy - 1 - y, clr)
			}
		}
	} else if p1 == "o" {
		clr := ClrLine
		img.Set(sx / 2 - 1, 0, clr)
		img.Set(sx / 2, 0, clr)
	} else if p1 == "s" {
		for b := 0; b <= sy / 2; b++ {
			y := sy / 4 + b
			if p2 == "d" { y = sy * 3 / 4 - b }
			xl := sx / 2 - b * (sx / 2 - 1) / (sy / 2)
			xr := sx / 2 + b * (sx / 2 - 1) / (sy / 2)
			for x := xl; x <= xr; x++ {
				clr := ClrArr
				if pe == "a" {
					if x == xl || x == xr || b == sy / 2 {
						clr = ClrArrBord
					} else {
						clr = ClrArrAct
					}
				}
				img.Set(x, y, clr)
			}
		}
	}
	return img
}


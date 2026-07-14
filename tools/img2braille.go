//go:build ignore

// img2braille converts a portrait photo to Unicode Braille art for the SSH portfolio.
// Usage: go run tools/img2braille.go <image.png|jpg>
//
// The output is a Go []string literal ready to paste into views/home.go.
package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

// Target dimensions in braille characters
const (
	brailleCols = 50
	brailleRows = 26
	imgW        = brailleCols * 2 // 100 px — each braille char = 2 cols
	imgH        = brailleRows * 4 // 104 px — each braille char = 4 rows
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run tools/img2braille.go <image>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "open:", err)
		os.Exit(1)
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "decode:", err)
		os.Exit(1)
	}

	// 1. Resize + grayscale with bilinear interpolation
	gray := resizeGray(src, imgW, imgH)

	// 2. Contrast stretch + gamma boost (pulls out details in dark photos)
	autoContrast(gray, imgW, imgH)
	gammaBoost(gray, imgW, imgH, 0.75) // 0.75 = brighten mid-tones

	// 3. Atkinson dithering for crisp dot rendering
	dots := atkinsonDither(gray, imgW, imgH)

	// 4. Encode to braille and print Go literal
	fmt.Printf("// Generated from %s\n", os.Args[1])
	fmt.Printf("// Size: %d cols × %d rows (braille characters)\n", brailleCols, brailleRows)
	fmt.Println("var portraitArt = []string{")
	for row := 0; row < brailleRows; row++ {
		fmt.Print("\t\"")
		for col := 0; col < brailleCols; col++ {
			offset := 0
			// Unicode Braille 8-dot bit layout:
			//  dot1(0x01) dot4(0x08)
			//  dot2(0x02) dot5(0x10)
			//  dot3(0x04) dot6(0x20)
			//  dot7(0x40) dot8(0x80)
			px := col * 2
			py := row * 4
			if dots[px][py]   { offset |= 0x01 }
			if dots[px][py+1] { offset |= 0x02 }
			if dots[px][py+2] { offset |= 0x04 }
			if px+1 < imgW {
				if dots[px+1][py]   { offset |= 0x08 }
				if dots[px+1][py+1] { offset |= 0x10 }
				if dots[px+1][py+2] { offset |= 0x20 }
			}
			if dots[px][py+3]   { offset |= 0x40 }
			if px+1 < imgW && dots[px+1][py+3] { offset |= 0x80 }
			fmt.Printf("%c", rune(0x2800+offset))
		}
		fmt.Println("\",")
	}
	fmt.Println("}")
	fmt.Printf("\nconst portraitWidth = %d\n", brailleCols)
}

// resizeGray converts src to a [w][h] float64 grayscale grid using bilinear interpolation.
func resizeGray(src image.Image, w, h int) [][]float64 {
	b := src.Bounds()
	sw := b.Max.X - b.Min.X
	sh := b.Max.Y - b.Min.Y

	out := make([][]float64, w)
	for x := 0; x < w; x++ {
		out[x] = make([]float64, h)
		for y := 0; y < h; y++ {
			sx := float64(x) * float64(sw) / float64(w)
			sy := float64(y) * float64(sh) / float64(h)

			x0, y0 := int(math.Floor(sx)), int(math.Floor(sy))
			x1, y1 := x0+1, y0+1
			if x1 >= sw { x1 = sw - 1 }
			if y1 >= sh { y1 = sh - 1 }

			fx, fy := sx-math.Floor(sx), sy-math.Floor(sy)

			c00 := lum(src.At(b.Min.X+x0, b.Min.Y+y0))
			c10 := lum(src.At(b.Min.X+x1, b.Min.Y+y0))
			c01 := lum(src.At(b.Min.X+x0, b.Min.Y+y1))
			c11 := lum(src.At(b.Min.X+x1, b.Min.Y+y1))

			out[x][y] = c00*(1-fx)*(1-fy) + c10*fx*(1-fy) +
				c01*(1-fx)*fy + c11*fx*fy
		}
	}
	return out
}

// lum computes the perceptual luminance (0..1) of a colour.
func lum(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	return (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535.0
}

// autoContrast stretches the histogram to span [0, 1].
func autoContrast(gray [][]float64, w, h int) {
	lo, hi := 1.0, 0.0
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			v := gray[x][y]
			if v < lo { lo = v }
			if v > hi { hi = v }
		}
	}
	span := hi - lo
	if span < 1e-9 { return }
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			gray[x][y] = (gray[x][y] - lo) / span
		}
	}
}

// gammaBoost applies gamma correction: v → v^gamma. gamma<1 brightens mid-tones.
func gammaBoost(gray [][]float64, w, h int, gamma float64) {
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			gray[x][y] = math.Pow(gray[x][y], gamma)
		}
	}
}

// atkinsonDither applies Atkinson dithering and returns a boolean pixel grid.
// true = bright (dot is set), false = dark.
func atkinsonDither(gray [][]float64, w, h int) [][]bool {
	buf := make([][]float64, w)
	for x := 0; x < w; x++ {
		buf[x] = make([]float64, h)
		copy(buf[x], gray[x])
	}
	out := make([][]bool, w)
	for x := 0; x < w; x++ {
		out[x] = make([]bool, h)
	}

	push := func(x, y int, e float64) {
		if x >= 0 && x < w && y >= 0 && y < h {
			buf[x][y] += e
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			old := buf[x][y]
			var nv float64
			if old >= 0.5 {
				nv = 1.0
				out[x][y] = true
			}
			e := (old - nv) / 8.0
			push(x+1, y, e)
			push(x+2, y, e)
			push(x-1, y+1, e)
			push(x, y+1, e)
			push(x+1, y+1, e)
			push(x, y+2, e)
		}
	}
	return out
}

//go:build ignore

// img2braille converts a portrait photo to Unicode Braille art for the SSH portfolio.
//
// Usage:
//
//	go run tools/img2braille.go <image.png|jpg> [x1,y1,x2,y2] [gamma]
//
// Crop values are 0.0–1.0 fractions of the original image dimensions.
// Example (face-only crop, 75% height):
//
//	go run tools/img2braille.go photo.png 0.08,0.02,0.92,0.78 0.65
package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strconv"
	"strings"
)

// Target dimensions in braille characters
const (
	brailleCols = 40
	brailleRows = 22
	imgW        = brailleCols * 2 // each braille char = 2 pixel columns
	imgH        = brailleRows * 4 // each braille char = 4 pixel rows
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run tools/img2braille.go <image> [x1,y1,x2,y2] [gamma]")
		os.Exit(1)
	}

	// Crop region (fractions 0.0–1.0), default = full image
	cropX1, cropY1, cropX2, cropY2 := 0.0, 0.0, 1.0, 1.0
	if len(os.Args) >= 3 {
		parts := strings.Split(os.Args[2], ",")
		if len(parts) == 4 {
			cropX1, _ = strconv.ParseFloat(parts[0], 64)
			cropY1, _ = strconv.ParseFloat(parts[1], 64)
			cropX2, _ = strconv.ParseFloat(parts[2], 64)
			cropY2, _ = strconv.ParseFloat(parts[3], 64)
		}
	}

	// Gamma: <1 brightens mid-tones, >1 darkens; default 0.7
	gamma := 0.70
	if len(os.Args) >= 4 {
		gamma, _ = strconv.ParseFloat(os.Args[3], 64)
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

	// 1. Resize + grayscale with bilinear interpolation (respecting crop window)
	gray := resizeGrayCrop(src, imgW, imgH, cropX1, cropY1, cropX2, cropY2)

	// 2. Contrast stretch (pulls histogram to [0,1])
	autoContrast(gray, imgW, imgH)

	// 3. Background suppression — push near-uniform pixels toward black
	//    (eliminates gray-wall noise that shouldn't show as dots)
	suppressBackground(gray, imgW, imgH, 0.18)

	// 4. Gamma boost for facial detail
	gammaBoost(gray, imgW, imgH, gamma)

	// 5. Atkinson dithering for crisp dot rendering
	dots := atkinsonDither(gray, imgW, imgH)

	// 6. Output Go literal
	fmt.Printf("// Generated from %s  crop=[%.2f,%.2f,%.2f,%.2f]  gamma=%.2f\n",
		os.Args[1], cropX1, cropY1, cropX2, cropY2, gamma)
	fmt.Printf("// Size: %d cols × %d rows\n", brailleCols, brailleRows)
	fmt.Println("var portraitArt = []string{")
	for row := 0; row < brailleRows; row++ {
		fmt.Print("\t\"")
		for col := 0; col < brailleCols; col++ {
			offset := 0
			px := col * 2
			py := row * 4
			if safeGet(dots, px, py, imgW, imgH)     { offset |= 0x01 }
			if safeGet(dots, px, py+1, imgW, imgH)   { offset |= 0x02 }
			if safeGet(dots, px, py+2, imgW, imgH)   { offset |= 0x04 }
			if safeGet(dots, px+1, py, imgW, imgH)   { offset |= 0x08 }
			if safeGet(dots, px+1, py+1, imgW, imgH) { offset |= 0x10 }
			if safeGet(dots, px+1, py+2, imgW, imgH) { offset |= 0x20 }
			if safeGet(dots, px, py+3, imgW, imgH)   { offset |= 0x40 }
			if safeGet(dots, px+1, py+3, imgW, imgH) { offset |= 0x80 }
			fmt.Printf("%c", rune(0x2800+offset))
		}
		fmt.Println("\",")
	}
	fmt.Println("}")
	fmt.Printf("\nconst portraitWidth = %d\n", brailleCols)
}

func safeGet(dots [][]bool, x, y, w, h int) bool {
	if x < 0 || x >= w || y < 0 || y >= h { return false }
	return dots[x][y]
}

// resizeGrayCrop converts src to a [w][h] float64 grayscale grid, reading only
// the crop window (x1,y1)→(x2,y2) expressed as image fractions.
func resizeGrayCrop(src image.Image, w, h int, x1, y1, x2, y2 float64) [][]float64 {
	b := src.Bounds()
	sw := float64(b.Max.X - b.Min.X)
	sh := float64(b.Max.Y - b.Min.Y)

	// Crop in pixel space
	cx0 := x1 * sw
	cy0 := y1 * sh
	cw  := (x2 - x1) * sw
	ch  := (y2 - y1) * sh

	out := make([][]float64, w)
	for x := 0; x < w; x++ {
		out[x] = make([]float64, h)
		for y := 0; y < h; y++ {
			sx := cx0 + float64(x)*cw/float64(w)
			sy := cy0 + float64(y)*ch/float64(h)

			x0, y0 := int(math.Floor(sx)), int(math.Floor(sy))
			x1p, y1p := x0+1, y0+1
			if x1p >= int(sw) { x1p = int(sw) - 1 }
			if y1p >= int(sh) { y1p = int(sh) - 1 }

			fx, fy := sx-math.Floor(sx), sy-math.Floor(sy)

			c00 := lum(src.At(b.Min.X+x0, b.Min.Y+y0))
			c10 := lum(src.At(b.Min.X+x1p, b.Min.Y+y0))
			c01 := lum(src.At(b.Min.X+x0, b.Min.Y+y1p))
			c11 := lum(src.At(b.Min.X+x1p, b.Min.Y+y1p))

			out[x][y] = c00*(1-fx)*(1-fy) + c10*fx*(1-fy) +
				c01*(1-fx)*fy + c11*fx*fy
		}
	}
	return out
}

func lum(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	return (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535.0
}

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

// suppressBackground computes local variance in a small window and pushes
// low-variance (flat background) regions toward 0 (black = no dots).
// sigma controls suppression strength (0.1–0.25 is good).
func suppressBackground(gray [][]float64, w, h int, sigma float64) {
	const r = 3 // neighbourhood radius
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			// Compute mean in window
			sum, n := 0.0, 0.0
			for dx := -r; dx <= r; dx++ {
				for dy := -r; dy <= r; dy++ {
					nx, ny := x+dx, y+dy
					if nx >= 0 && nx < w && ny >= 0 && ny < h {
						sum += gray[nx][ny]
						n++
					}
				}
			}
			mean := sum / n
			// Compute variance
			vsum := 0.0
			for dx := -r; dx <= r; dx++ {
				for dy := -r; dy <= r; dy++ {
					nx, ny := x+dx, y+dy
					if nx >= 0 && nx < w && ny >= 0 && ny < h {
						d := gray[nx][ny] - mean
						vsum += d * d
					}
				}
			}
			variance := vsum / n
			// Low-variance pixels (uniform background) → darken toward 0
			if variance < sigma*sigma {
				// Weight: 0 = uniform background, 1 = highly textured
				weight := math.Sqrt(variance) / sigma
				gray[x][y] *= weight
			}
		}
	}
}

func gammaBoost(gray [][]float64, w, h int, gamma float64) {
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			gray[x][y] = math.Pow(gray[x][y], gamma)
		}
	}
}

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
			push(x+1, y, e); push(x+2, y, e)
			push(x-1, y+1, e); push(x, y+1, e); push(x+1, y+1, e)
			push(x, y+2, e)
		}
	}
	return out
}

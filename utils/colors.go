package utils

import (
	"fmt"
	"math"
	"math/rand"
)


type ColorA struct {
	R, G, B, A float32
}

type Color struct {
	R, G, B float64
}


type HSL [3]float64

func RandColor() ColorA {
	return ColorA{rand.Float32(), rand.Float32(), rand.Float32(), 1.0}
}

// returns float32 colorA for sketches
func StepColor(c1, c2 Color, t, i int) ColorA{
	factorStep := 1 / (float64(t) - 1.0)
	c := InterpolateColor(c1, c2, factorStep * float64(i))
	return ColorA{
		R: float32(c.R),
		G: float32(c.G),
		B: float32(c.B),
	}
}

func Rgb2Hsl(c Color) HSL {
	max := math.Max(math.Max(c.R, c.G), c.B)
	min := math.Min(math.Min(c.R, c.G), c.B)
	h := (max + min) / 2
	s := h
	l := h

	if max == min {
		h = 0
		s = 0 // achromatic
	} else {
		d := max - min
		s = d / (max + min)
		if l > 0.5 {
			s = d / (2 - max - min)
		}
		switch max {
		case c.R:
			t := 0.0
			if c.G < c.B {
				t = 6.0
			}
			h = (c.G-c.B)/d + t
		case c.G:
			h = (c.B-c.R)/d + 2.0
		case c.B:
			h = (c.R-c.G)/d + 4.0
		}
		h /= 6
	}

	return HSL{h, s, l}
}

func InterpolateColor(c1, c2 Color, factor float64) Color {
	result := new(Color)
	result.R = c1.R + factor*(c2.R-c1.R)
	result.G = c1.G + factor*(c2.G-c1.G)
	result.B = c1.B + factor*(c2.B-c1.B)
	return *result
}

func InterpolateHSL(h1, h2 HSL, factor float64) HSL {
	o := h1
	for i := 0; i < 3; i++ {
		o[i] += factor * (h1[i] - h2[i])
	}
	return o
}

func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
		fmt.Print(t)
	}
	if t > 1 {
		t -= 1
	}
	if t < 1/6 {
		return p + (q-p)*6*t
	}
	if t < 1/2 {
		return q
	}
	if t < 2/3 {
		return p + (q-p)*(2/3-t)*6
	}
	return p
}

func r2h(c Color) string {
	rgb := []int{
		int(round(c.R*255, 0.5, 0)),
		int(round(c.G*255, 0.5, 0)),
		int(round(c.B*255, 0.5, 0)),
	}
	t := (1 << 24) + (rgb[0] << 16) + (rgb[1] << 8) + rgb[2]
	s := fmt.Sprintf("%x", t)
	return "#" + s
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func Hsl2Rgb(h HSL) Color {
	l := h[2]
	if h[1] == 0 {
		return Color{1.0, 1.0, 1.0}
	} else {
		s := h[1]
		q := l + s - l*s
		if l < 0.5 {
			q = l * (1 + s)
		}
		p := 2*l - q
		r := hue2rgb(p, q, h[0]+1/3)
		g := hue2rgb(p, q, h[0])
		b := hue2rgb(p, q, h[0]-1/3)
		return Color{r, g, b}
	}
}
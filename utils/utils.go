package utils

import "image/color"

func DimColor(c color.Color, factor float64) color.Color {
	r, g, b, a := c.RGBA()
	return color.NRGBA{
		R: uint8(float64(r>>8) * factor),
		G: uint8(float64(g>>8) * factor),
		B: uint8(float64(b>>8) * factor),
		A: uint8(a >> 8),
	}
}

func DatBoiHandler() {

}

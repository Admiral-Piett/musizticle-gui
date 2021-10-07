package main

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"image/color"
)

func songLineWrapMarginsBorder(gtx layout.Context, d layout.Dimensions) layout.Dimensions {
	margins := layout.Inset{
		Top:    unit.Dp(1),
		Right:  unit.Dp(1),
		Bottom: unit.Dp(1),
		Left:   unit.Dp(1),
	}
	border := widget.Border{
		Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
		CornerRadius: unit.Dp(3),
		Width:        unit.Dp(2),
	}
	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return d
		})
	})
}

func songFieldsMarginsBorder(gtx layout.Context, d layout.Dimensions) layout.Dimensions {
	margins := layout.Inset{
		Top:    unit.Dp(0.5),
		Right:  unit.Dp(0),
		Bottom: unit.Dp(0.5),
		Left:   unit.Dp(0),
	}
	border := widget.Border{
		Color:        color.NRGBA{R: 192, G: 192, B: 192, A: 200},
		CornerRadius: unit.Dp(0),
		Width:        unit.Dp(1),
	}
	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return d
		})
	})
}

func headerFieldsMarginsBorder(gtx layout.Context, d layout.Dimensions) layout.Dimensions {
	margins := layout.Inset{
		Top:    unit.Dp(5),
		Right:  unit.Dp(0),
		Bottom: unit.Dp(10),
		Left:   unit.Dp(0),
	}
	border := widget.Border{
		Color:        color.NRGBA{R: 192, G: 192, B: 192, A: 200},
		CornerRadius: unit.Dp(0),
		Width:        unit.Dp(1),
	}
	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return d
		})
	})
}

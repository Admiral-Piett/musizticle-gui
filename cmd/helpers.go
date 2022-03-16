package main

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
)

// ----- Window Stuff -----
func outerSongListWrapper(gtx layout.Context, f func(songList []*Song) []layout.FlexChild, songList []*Song) layout.Dimensions {
	l := layout.Flex{
		// Vertical alignment, from top to bottom
		Axis: layout.Vertical,
		// Empty space is left at the start, i.e. at the top
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		f(songList)...,
	)
	return l
}

// ----- Theme Stuff ------
// Copied from widget.material.theme.NewTheme
func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func CreateTheme(fontCollection []text.FontFace) *material.Theme {
	t := &material.Theme{
		Shaper: text.NewCache(fontCollection),
	}
	t.Palette = material.Palette{
		Fg:         rgb(0x000000),
		Bg:         rgb(0x000000),
		ContrastBg: rgb(0x717EDD),
		ContrastFg: rgb(0xffffff),
	}
	t.TextSize = unit.Sp(16)

	t.Icon.CheckBoxChecked = mustIcon(widget.NewIcon(icons.ToggleCheckBox))
	t.Icon.CheckBoxUnchecked = mustIcon(widget.NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.Icon.RadioChecked = mustIcon(widget.NewIcon(icons.ToggleRadioButtonChecked))
	t.Icon.RadioUnchecked = mustIcon(widget.NewIcon(icons.ToggleRadioButtonUnchecked))

	// 38dp is on the lower end of possible finger size.
	t.FingerSize = unit.Dp(38)

	return t
}

// -------- End Theme -----------

// --------- Styling --------------
func songLineMargins(gtx layout.Context, d layout.Dimensions) layout.Dimensions {
	margins := layout.Inset{
		Top:    unit.Dp(2),
		Right:  unit.Dp(1),
		Bottom: unit.Dp(3),
		Left:   unit.Dp(1),
	}
	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return d
	})
}

func songFieldsMargins(gtx layout.Context, d layout.Dimensions) layout.Dimensions {
	margins := layout.Inset{
		Top:    unit.Dp(0.5),
		Right:  unit.Dp(3),
		Bottom: unit.Dp(0.5),
		Left:   unit.Dp(3),
	}
	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return d
	})
}

func headerFieldsMargins(gtx layout.Context, d layout.Dimensions) layout.Dimensions {
	margins := layout.Inset{
		Top:    unit.Dp(5),
		Right:  unit.Dp(0),
		Bottom: unit.Dp(10),
		Left:   unit.Dp(0),
	}
	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return d
	})
}

// --------- Auth --------------
func generateHeaders() {
	// TODO - HERE - generate headers for all the HTTP requests
}

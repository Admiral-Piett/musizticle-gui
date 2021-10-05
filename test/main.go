package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		// create new window
		w := app.NewWindow()
		th := material.NewTheme(gofont.Collection())
		var ops op.Ops

		for e := range w.Events() {

			// detect what type of event
			switch e := e.(type) {

			// this is sent when the application should re-render.
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				const listLen = 1e6

				var list = layout.List{Axis: layout.Vertical}
				list.Layout(gtx, listLen, func(gtx layout.Context, i int) layout.Dimensions {
					text := fmt.Sprintf("Item %d", i)
					l := material.Label(th, unit.Dp(float32(20)), text)
					return l.Layout(gtx)
				})
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}

package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"log"
	"os"
)

// TODO - Explore other font collections
var th = material.NewTheme(gofont.Collection())

// NOTE: This has to be set out here, otherwise the scroll context can't be maintained across events.
//  https://gioui.org/doc/architecture#my-list-list-won-t-scroll
// QUESTION: Set and reset this to filter by things?
var displayList = &layout.List{Axis: layout.Vertical}


func draw(w *app.Window) error {
	var ops op.Ops
	var startButton widget.Clickable

	var songs Songs

	for {
		//log.Println("loop")
		select {
		case <-displayChange:
			log.Println("update display")
			w.Invalidate()
		case e := <-w.Events():
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)

				if startButton.Clicked() {
					log.Println("I'm clicked")
					return nil
				}
				layout.Flex{
					// Vertical alignment, from top to bottom
					Axis: layout.Vertical,
					// Empty space is left at the start, i.e. at the top
					Spacing: layout.SpaceStart,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return songs.ShowSongs(gtx)
					}),
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &startButton, "Button me up")
							return btn.Layout(gtx)
						},
					),
				)
				e.Frame(gtx.Ops)
			}
		}

	}
	return nil
}

func main() {
	go func() {
		// create new window
		w := app.NewWindow()

		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

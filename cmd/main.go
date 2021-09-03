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

func draw(w *app.Window) error {
	var ops op.Ops
	var startButton widget.Clickable

	// TODO - Explore other font collections
	th := material.NewTheme(gofont.Collection())
	for {
		select {
			case e := <-w.Events():
				switch e := e.(type) {
				case system.FrameEvent:
					if startButton.Clicked() {
						log.Println("I'm clicked")
					}
					gtx := layout.NewContext(&ops, e)
					layout.Flex{
						// Vertical alignment, from top to bottom
						Axis: layout.Vertical,
						// Empty space is left at the start, i.e. at the top
						Spacing: layout.SpaceStart,
					}.Layout(gtx, layout.Rigid(
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

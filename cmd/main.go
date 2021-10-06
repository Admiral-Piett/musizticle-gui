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


func (a *App) draw() error {
	var ops op.Ops
	var startButton widget.Clickable

	for {
		//log.Println("loop")
		select {
		case <-displayChange:
			log.Println("update display")
			a.window.Invalidate()
		case e := <-a.window.Events():
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)

				if startButton.Clicked() {
					log.Println("I'm clicked")
				}
				if a.songs.reload.Clicked() {
					log.Println("Reload clicked")
					go a.songs.initSongs()
				}
				layout.Flex{
					// Vertical alignment, from top to bottom
					Axis: layout.Vertical,
					// Empty space is left at the start, i.e. at the top
					Spacing: layout.SpaceStart,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.ShowSongs(gtx)
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
		s := Songs{}
		a := App{
			displayList: &layout.List{Axis: layout.Vertical},
			songs: s,
			window: w,
		}
		go a.songs.initSongs()

		if err := a.draw(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

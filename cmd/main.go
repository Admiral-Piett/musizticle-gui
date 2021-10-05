package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"log"
	"os"
)

type Songs struct {
	song_list []Song
	selected  int
}

type Song struct {
	line  widget.Clickable
	Title string
}
var songs Songs
// TODO - Explore other font collections
var th = material.NewTheme(gofont.Collection())
// NOTE: This has to be set out here, otherwise the scroll context can't be maintained across events.
//  https://gioui.org/doc/architecture#my-list-list-won-t-scroll
var trackList = &layout.List{Axis: layout.Vertical}

func (s *Songs) setList() {
	for i := 1; i <= 100; i++ {
		songs.song_list = append(songs.song_list, Song{Title: fmt.Sprintf("Number %d", i)})
	}
}

// FIXME - need this?
func generateListEntry(gtx layout.Context, i int) layout.Dimensions {
	text := fmt.Sprintf("Item %d", i)
	l := material.Label(th, unit.Dp(float32(20)), text)
	return l.Layout(gtx)
}


// TODO - NEXT - Wire up the api calls and populate this list for real
func setSongs(gtx layout.Context) layout.Dimensions {
	l_dimensions := trackList.Layout(gtx, len(songs.song_list), func(gtx layout.Context, index int) layout.Dimensions {
		s := &songs.song_list[index]
		if s.line.Clicked() {
			songs.selected = index
			fmt.Printf("Clicked - %d", index)
		}
		dims := material.Clickable(gtx, &s.line, func(gtx layout.Context) layout.Dimensions {
			return material.Label(th, unit.Dp(float32(20)), s.Title).Layout(gtx)
		})
		return dims
	})
	return l_dimensions
}

func draw(w *app.Window) error {
	var ops op.Ops
	var startButton widget.Clickable

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.FrameEvent:
				if startButton.Clicked() {
					log.Println("I'm clicked")
					return nil
				}
				gtx := layout.NewContext(&ops, e)
				layout.Flex{
					// Vertical alignment, from top to bottom
					Axis: layout.Vertical,
					// Empty space is left at the start, i.e. at the top
					Spacing: layout.SpaceStart,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return setSongs(gtx)
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
	songs.setList()
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

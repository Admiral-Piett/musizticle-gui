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

func (s *Songs) setList() {
	for i := 1; i <= 100; i++ {
		songs.song_list = append(songs.song_list, Song{Title: fmt.Sprintf("Number %d", i)})
	}
}

func setSongs(gtx layout.Context, th *material.Theme) layout.Dimensions {
	//p := layout.Position{
	//	First:      0,
	//	Offset:     0,
	//	OffsetLast: 0,
	//	Count:      len(songs.song_list),
	//	Length:     0,
	//}

	l := layout.List{Axis: layout.Vertical}
	l_dimensions := l.Layout(gtx, len(songs.song_list), func(gtx layout.Context, index int) layout.Dimensions {
		s := &songs.song_list[index]
		if s.line.Clicked() {
			songs.selected = index
			fmt.Printf("Clicked - %d", index)
		}
		dims := material.Clickable(gtx, &s.line, func(gtx layout.Context) layout.Dimensions {
			//return widget.Label{}.Layout(gtx, th.Shaper, th., th.TextSize, s.Title)
			//// FIXME - Do I need this uniform thing?  I don't think so, could likely just do a text.
			return layout.UniformInset(unit.Sp(12)).Layout(gtx,
				material.H6(th, s.Title).Layout,
			)
		})
		return dims
		//return layout.Stack{Alignment: layout.S}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		//	//dimensions
		//	dims := material.Clickable(gtx, &s.line, func(gtx layout.Context) layout.Dimensions {
		//		return layout.UniformInset(unit.Sp(12)).Layout(gtx,
		//			material.H6(th, s.Title).Layout,
		//		)
		//	})
		//	return dims
		//}))
	})
	return l_dimensions
}

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
				//log.Printf("event happened %v", e)
				if startButton.Clicked() {
					log.Println("I'm clicked")
					return nil
				}
				gtx := layout.NewContext(&ops, e)
				//setSongs(gtx, th)
				layout.Flex{
					// Vertical alignment, from top to bottom
					Axis: layout.Vertical,
					// Empty space is left at the start, i.e. at the top
					Spacing: layout.SpaceStart,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return setSongs(gtx, th)
					}),
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &startButton, "Button me up")
							return btn.Layout(gtx)
						},
					),
				)
				w.Invalidate()
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

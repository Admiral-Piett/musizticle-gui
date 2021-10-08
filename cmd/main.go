package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
	"time"
)

func (a *App) draw() error {
	var ops op.Ops

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

				if a.songs.reload.Clicked() {
					log.Println("Reload clicked")
					go a.songs.initSongs()
				}
				if a.play.Clicked() {
					log.Println("Play clicked")
					go a.playSong()
				}
				if a.stop.Clicked() {
					log.Println("Stop clicked")
					go a.clickStop()
				}
				if a.next.Clicked() {
					log.Println("Next clicked")
					go a.clickNext()
				}
				if a.previous.Clicked() {
					log.Println("Previous clicked")
					go a.clickPrevious()
				}
				layout.Flex{
					// Vertical alignment, from top to bottom
					Axis: layout.Vertical,
					// Empty space is left at the start, i.e. at the top
					Spacing: layout.SpaceStart,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return SongsHeader(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.SongsList(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.MediaToolBar(gtx)
					}),
				)
				e.Frame(gtx.Ops)
			}
		}

	}
	return nil
}

// Initialize the speaker with the settings from the first song we have.
func (a *App) SetUpSpeaker() {
	log.Println("SettingUpSpeakerStart")
	_, format, err := a.getSong(1)
	if err != nil {
		log.Printf("SettingUpSpeakerFailure - %+v\n", err)
		panic(err)
	}
	a.SampleRate = format.SampleRate
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Printf("SettingUpSpeakerFailure - %+v\n", err)
		panic(err)
	}
	log.Println("SettingUpSpeakerFinish")
}

func main() {
	go func() {
		// create new window
		w := app.NewWindow(app.Title("Media Gui"), app.Size(unit.Dp(1500), unit.Dp(900)))
		s := Songs{}
		a := App{
			displayList:    &layout.List{Axis: layout.Vertical},
			songs:          s,
			window:         w,
			SelectedSongIndex: 0,
		}
		go a.songs.initSongs()

		//Put an invalid song id on the playing queue to start with
		playing <- -1
		a.SetUpSpeaker()

		if err := a.draw(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

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
        select {
        case <-displayChange:
            log.Println("update display")
            a.window.Invalidate()
        case e := <-a.window.Events():
            switch e := e.(type) {
            case system.FrameEvent:
                gtx := layout.NewContext(&ops, e)

                if loginButton.Clicked() {
                    login()
                    op.InvalidateOp{}.Add(gtx.Ops)
                }
                if a.songs.reload.Clicked() {
                    log.Println("Reload clicked")
                    go a.initSongs()
                }
                if a.playBtn.Clicked() {
                    log.Println("Play clicked")
                    a.clickPlay()
                }
                if a.stopBtn.Clicked() {
                    log.Println("Stop clicked")
                    go a.clickStop()
                }
                if a.nextBtn.Clicked() {
                    log.Println("Next clicked")
                    go a.clickNext()
                }
                if a.previousBtn.Clicked() {
                    log.Println("Previous clicked")
                    go a.clickPrevious()
                }
                if a.homeTab.Clicked() {
                    log.Println("Home Tab clicked")
                    a.selectedTab = HOME_TAB
                    op.InvalidateOp{}.Add(gtx.Ops)
                }
                if a.nextTab.Clicked() {
                    log.Println("Next Tab clicked")
                    a.selectedTab = NEXT_TAB
                    op.InvalidateOp{}.Add(gtx.Ops)
                }
                if a.previousTab.Clicked() {
                    log.Println("Previous Tab clicked")
                    a.selectedTab = PREVIOUS_TAB
                    op.InvalidateOp{}.Add(gtx.Ops)
                }
                if loginRequired {
                    a.displayLoginWindow(gtx)
                } else if a.selectedTab == NEXT_TAB {
                    if len(a.navQueueNext) == 0 {
                        outerSongListWrapper(gtx, a.tabDisplay, a.songList)
                    } else {
                        outerSongListWrapper(gtx, a.tabDisplay, a.navQueueNext)
                    }
                } else if a.selectedTab == PREVIOUS_TAB {
                    if len(a.navQueuePrevious) == 0 {
                        outerSongListWrapper(gtx, a.tabDisplay, a.songList)
                    } else {
                        outerSongListWrapper(gtx, a.tabDisplay, a.navQueuePrevious)
                    }
                } else {
                    outerSongListWrapper(gtx, a.tabDisplay, a.songList)
                }
                e.Frame(gtx.Ops)
            }
        case p := <-progressIncrementer:
            if a.playing && progress < 1 {
                progress += p
                a.window.Invalidate()
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

func resetProgress(duration int) {
    // Empty the channel
    for len(progressIncrementer) > 0 {
        <-progressIncrementer
    }
    // The goal is to find the value to increment every 25th of a second in order to reach one when the song ends.
    // So start by multiplying the duration by 25 (how many updates we'll have), and divide 1 by that value to see
    //  what value would go into 1 by our total times.
    incrementCount := duration * 25
    incrementValue := 1 / (float32(incrementCount))
    go func() {
        for i := 0; i <= (incrementCount); i++ {
            time.Sleep(time.Second / 25)
            progressIncrementer <- incrementValue
        }
    }()
}

func main() {
    // TODO - NEXT - figure out how to calculate this to hit one as the song finishes.  Need to find the length of
    //  the song somewhere - from the file maybe?
    progressIncrementer = make(chan float32)
    go func() {
        // create new window
        w := app.NewWindow(app.Title("Media Gui"), app.Size(unit.Dp(1500), unit.Dp(900)))
        s := Songs{}
        a := App{
            displayList: &layout.List{Axis: layout.Vertical},
            songs:       s,
            window:      w,
            selectedTab: HOME_TAB,
        }

        loginRequired = true
        //go a.initSongs()
        go a.startUp()

        if err := a.draw(); err != nil {
            log.Fatal(err)
        }
        os.Exit(0)
    }()
    app.Main()
}

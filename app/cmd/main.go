package main

import (
    "fmt"
    "github.com/Admiral-Piett/musizticle-gui/app/config"
)

func main() {
    settings := config.NewSettings()
    defer func() {
        settings.SetConfigFile()
    }()
    fmt.Println(settings)
}

// Initialize the speaker with the settings from the first song we have.
//func (a *App) SetUpSpeaker() {
//    log.Println("SettingUpSpeakerStart")
//    _, format, err := a.getSong(1)
//    if err != nil {
//        log.Printf("SettingUpSpeakerFailure - %+v\n", err)
//        panic(err)
//    }
//    a.SampleRate = format.SampleRate
//    err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
//    if err != nil {
//        log.Printf("SettingUpSpeakerFailure - %+v\n", err)
//        panic(err)
//    }
//    log.Println("SettingUpSpeakerFinish")
//}
